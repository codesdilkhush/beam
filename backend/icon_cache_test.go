package backend

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIconCache_GetOrConvert(t *testing.T) {
	tempDir := t.TempDir()
	iconFile := filepath.Join(tempDir, "test.png")

	// Dummy image bytes
	dummyData := []byte("fake-png-binary-data")
	if err := os.WriteFile(iconFile, dummyData, 0644); err != nil {
		t.Fatalf("failed to write dummy icon: %v", err)
	}

	cacheDir := filepath.Join(tempDir, "cache")
	ic := &IconCache{
		cacheDir: cacheDir,
		memory:   make(map[string]string),
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}

	// 1. Initial read and conversion
	dataURL := ic.GetOrConvert(iconFile)
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Fatalf("expected data URL to start with png prefix, got %q", dataURL)
	}

	// 2. Second read should hit memory cache
	cachedDataURL := ic.GetOrConvert(iconFile)
	if cachedDataURL != dataURL {
		t.Fatalf("expected cached URL to match original, got %q", cachedDataURL)
	}

	// 3. Non-existent file should return empty string
	emptyURL := ic.GetOrConvert(filepath.Join(tempDir, "nonexistent.png"))
	if emptyURL != "" {
		t.Fatalf("expected empty string for missing file, got %q", emptyURL)
	}
}
