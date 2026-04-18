package sourcemap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover(t *testing.T) {
	dir := t.TempDir()

	files := []string{
		"index-DIPY5MEo.js.map",
		"vendor-abc123.js.map",
		"assets/chunk-xyz789.js.map",
	}

	for _, f := range files {
		path := filepath.Join(dir, f)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte(`{"version":3}`), 0644)
	}

	os.WriteFile(filepath.Join(dir, "index.js"), []byte("code"), 0644)
	os.WriteFile(filepath.Join(dir, "style.css"), []byte("css"), 0644)

	nodeModules := filepath.Join(dir, "node_modules", "pkg")
	os.MkdirAll(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "lib.js.map"), []byte(`{}`), 0644)

	maps, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(maps) != 3 {
		t.Errorf("expected 3 maps, got %d", len(maps))
	}

	foundFiles := make(map[string]bool)
	for _, m := range maps {
		foundFiles[m.Filename] = true
	}

	expected := []string{
		"index-DIPY5MEo.js.map",
		"vendor-abc123.js.map",
		"chunk-xyz789.js.map",
	}

	for _, e := range expected {
		if !foundFiles[e] {
			t.Errorf("expected to find %s", e)
		}
	}
}

func TestDiscoverEmpty(t *testing.T) {
	dir := t.TempDir()

	maps, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(maps) != 0 {
		t.Errorf("expected 0 maps, got %d", len(maps))
	}
}

func TestDiscoverSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()

	nodeModules := filepath.Join(dir, "node_modules", "some-pkg")
	os.MkdirAll(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "index.js.map"), []byte(`{}`), 0644)

	maps, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(maps) != 0 {
		t.Errorf("expected 0 maps, got %d", len(maps))
	}
}

func TestDiscoverSkipsGitDir(t *testing.T) {
	dir := t.TempDir()

	gitDir := filepath.Join(dir, ".git", "objects")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "pack.js.map"), []byte(`{}`), 0644)

	maps, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(maps) != 0 {
		t.Errorf("expected 0 maps, got %d", len(maps))
	}
}

func TestFilenameIsBasename(t *testing.T) {
	dir := t.TempDir()

	nested := filepath.Join(dir, "assets", "js")
	os.MkdirAll(nested, 0755)
	os.WriteFile(filepath.Join(nested, "app-hash123.js.map"), []byte(`{}`), 0644)

	maps, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(maps) != 1 {
		t.Fatalf("expected 1 map, got %d", len(maps))
	}

	if maps[0].Filename != "app-hash123.js.map" {
		t.Errorf("expected filename 'app-hash123.js.map', got '%s'", maps[0].Filename)
	}

	if !filepath.IsAbs(maps[0].Path) {
		t.Errorf("expected absolute path, got '%s'", maps[0].Path)
	}
}
