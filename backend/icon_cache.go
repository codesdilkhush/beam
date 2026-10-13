package backend

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type IconCache struct {
	mu       sync.RWMutex
	cacheDir string
	memory   map[string]string
}

var (
	iconCacheInstance *IconCache
	iconCacheOnce     sync.Once
)

func GetIconCache() *IconCache {
	iconCacheOnce.Do(func() {
		dir, err := os.UserHomeDir()
		if err != nil {
			dir = "."
		}
		cacheDir := filepath.Join(dir, ".cache/beam/icons")
		if mkErr := os.MkdirAll(cacheDir, 0755); mkErr != nil {
			cacheDir = "icon_cache"
		}

		iconCacheInstance = &IconCache{
			cacheDir: cacheDir,
			memory:   make(map[string]string),
		}
	})
	return iconCacheInstance
}

func (ic *IconCache) GetOrConvert(iconPath string) string {
	if iconPath == "" {
		return ""
	}

	ic.mu.RLock()
	if cached, ok := ic.memory[iconPath]; ok {
		ic.mu.RUnlock()
		return cached
	}
	ic.mu.RUnlock()

	// Check disk cache
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(iconPath)))
	diskFile := filepath.Join(ic.cacheDir, hash+".data")

	if data, err := os.ReadFile(diskFile); err == nil {
		dataURL := string(data)
		ic.mu.Lock()
		ic.memory[iconPath] = dataURL
		ic.mu.Unlock()
		return dataURL
	}

	// Read and convert original file
	fileBytes, err := os.ReadFile(iconPath)
	if err != nil {
		return ""
	}

	mimeType := "image/png"
	lower := strings.ToLower(iconPath)
	if strings.HasSuffix(lower, ".svg") {
		mimeType = "image/svg+xml"
	} else if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		mimeType = "image/jpeg"
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	// Save to memory and disk cache asynchronously
	ic.mu.Lock()
	ic.memory[iconPath] = dataURL
	ic.mu.Unlock()

	go func() {
		if writeErr := os.WriteFile(diskFile, []byte(dataURL), 0644); writeErr != nil {
			return
		}
	}()

	return dataURL
}
