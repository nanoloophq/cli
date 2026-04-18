package api

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nanoloop/cli/internal/sourcemap"
)

func TestUploadSourceMaps(t *testing.T) {
	dir := t.TempDir()
	mapContent := `{"version":3,"sources":["src/index.ts"],"mappings":"AAAA"}`

	files := []string{
		"index-abc123.js.map",
		"vendor-def456.js.map",
	}

	var testMaps []sourcemap.File
	for _, f := range files {
		path := filepath.Join(dir, f)
		os.WriteFile(path, []byte(mapContent), 0644)
		testMaps = append(testMaps, sourcemap.File{
			Path:     path,
			Filename: f,
		})
	}

	var receivedAppID, receivedRelease string
	var receivedFiles []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/v1/sourcemaps/upload" {
			t.Errorf("expected /v1/sourcemaps/upload, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected 'Bearer test-token', got '%s'", auth)
		}

		contentType := r.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			t.Fatalf("failed to parse content type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Errorf("expected multipart/form-data, got %s", mediaType)
		}

		reader := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("failed to read part: %v", err)
			}

			name := part.FormName()
			switch name {
			case "app_id":
				data, _ := io.ReadAll(part)
				receivedAppID = string(data)
			case "release":
				data, _ := io.ReadAll(part)
				receivedRelease = string(data)
			case "files":
				receivedFiles = append(receivedFiles, part.FileName())
			}
		}

		resp := UploadResult{
			Release: "v1.0.0",
			Uploaded: []UploadedFile{
				{Filename: "index-abc123.js.map", Release: "v1.0.0"},
				{Filename: "vendor-def456.js.map", Release: "v1.0.0"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	os.Setenv("NANOLOOP_API_URL", server.URL)
	defer os.Unsetenv("NANOLOOP_API_URL")

	client := NewClient("test-token")
	result, err := client.UploadSourceMaps("app-123", "v1.0.0", testMaps, "")
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	if receivedAppID != "app-123" {
		t.Errorf("expected app_id 'app-123', got '%s'", receivedAppID)
	}
	if receivedRelease != "v1.0.0" {
		t.Errorf("expected release 'v1.0.0', got '%s'", receivedRelease)
	}
	if len(receivedFiles) != 2 {
		t.Errorf("expected 2 files, got %d", len(receivedFiles))
	}

	expectedFiles := map[string]bool{
		"index-abc123.js.map":  false,
		"vendor-def456.js.map": false,
	}
	for _, f := range receivedFiles {
		if _, ok := expectedFiles[f]; ok {
			expectedFiles[f] = true
		} else {
			t.Errorf("unexpected file: %s", f)
		}
	}
	for f, found := range expectedFiles {
		if !found {
			t.Errorf("expected file not found: %s", f)
		}
	}

	if len(result.Uploaded) != 2 {
		t.Errorf("expected 2 uploaded, got %d", len(result.Uploaded))
	}
}

func TestUploadSourceMapsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid token"}`))
	}))
	defer server.Close()

	os.Setenv("NANOLOOP_API_URL", server.URL)
	defer os.Unsetenv("NANOLOOP_API_URL")

	dir := t.TempDir()
	mapPath := filepath.Join(dir, "test.js.map")
	os.WriteFile(mapPath, []byte(`{}`), 0644)

	client := NewClient("bad-token")
	_, err := client.UploadSourceMaps("app-123", "v1.0.0", []sourcemap.File{
		{Path: mapPath, Filename: "test.js.map"},
	}, "")

	if err == nil {
		t.Error("expected error for unauthorized request")
	}
}

func TestClientUsesEnvURL(t *testing.T) {
	os.Setenv("NANOLOOP_API_URL", "https://custom.api.com")
	defer os.Unsetenv("NANOLOOP_API_URL")

	client := NewClient("token")
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom URL, got %s", client.baseURL)
	}
}

func TestClientDefaultURL(t *testing.T) {
	os.Unsetenv("NANOLOOP_API_URL")

	client := NewClient("token")
	if client.baseURL != "https://api.nanoloop.app" {
		t.Errorf("expected default URL, got %s", client.baseURL)
	}
}
