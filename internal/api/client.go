package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/nanoloop/cli/internal/sourcemap"
)

const defaultBaseURL = "https://api.nanoloop.app"

type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

func NewClient(token string) *Client {
	baseURL := os.Getenv("NANOLOOP_API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		token:   token,
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

type UploadedFile struct {
	Filename string `json:"filename"`
	Release  string `json:"release"`
}

type UploadResult struct {
	Uploaded []UploadedFile `json:"uploaded"`
	Release  string         `json:"release"`
}

func (c *Client) UploadSourceMaps(appID, release string, maps []sourcemap.File, urlPrefix string) (*UploadResult, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if err := w.WriteField("app_id", appID); err != nil {
		return nil, err
	}
	if err := w.WriteField("release", release); err != nil {
		return nil, err
	}
	if urlPrefix != "" {
		if err := w.WriteField("url_prefix", urlPrefix); err != nil {
			return nil, err
		}
	}

	for _, m := range maps {
		content, err := os.ReadFile(m.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", m.Path, err)
		}

		part, err := w.CreateFormFile("files", m.Filename)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(content); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/sourcemaps/upload", &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed: %s - %s", resp.Status, string(body))
	}

	var result UploadResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
