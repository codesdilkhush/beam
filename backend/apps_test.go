package backend

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDesktopFile(t *testing.T) {
	tempDir := t.TempDir()
	desktopFile := filepath.Join(tempDir, "testapp.desktop")

	content := `[Desktop Entry]
Type=Application
Name=Test Application
Comment=A sample application for testing
Exec=testapp --launch %u
Icon=test-icon
Terminal=false
Categories=Utility;
`
	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test desktop file: %v", err)
	}

	app, err := parseDesktopFile(desktopFile)
	if err != nil {
		t.Fatalf("parseDesktopFile failed: %v", err)
	}

	if app.Name != "Test Application" {
		t.Errorf("expected Name 'Test Application', got %q", app.Name)
	}
	if app.Exec != "testapp --launch" {
		t.Errorf("expected Exec 'testapp --launch', got %q", app.Exec)
	}
	if app.Icon != "test-icon" {
		t.Errorf("expected Icon 'test-icon', got %q", app.Icon)
	}
	if app.Terminal {
		t.Errorf("expected Terminal false, got true")
	}
}

func TestParseDesktopFile_Hidden(t *testing.T) {
	tempDir := t.TempDir()
	desktopFile := filepath.Join(tempDir, "hidden.desktop")

	content := `[Desktop Entry]
Type=Application
Name=Hidden App
Exec=hidden
NoDisplay=true
`
	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test desktop file: %v", err)
	}

	_, err := parseDesktopFile(desktopFile)
	if err == nil {
		t.Fatalf("expected error for NoDisplay=true, got nil")
	}
}

func TestGetDesktopDirectories(t *testing.T) {
	dirs := GetDesktopDirectories()
	if len(dirs) == 0 {
		t.Fatalf("expected at least one desktop directory")
	}
}

func TestSearchApplications_EmptyAndQuery(t *testing.T) {
	// Populate cache with dummy apps
	appsMutex.Lock()
	appsCache = []AppInfo{
		{ID: "app1.desktop", Name: "Google Chrome", Exec: "google-chrome-stable", Comment: "Web Browser"},
		{ID: "app2.desktop", Name: "Visual Studio Code", Exec: "code", Comment: "Code Editor"},
		{ID: "app3.desktop", Name: "Terminal", Exec: "gnome-terminal", Comment: "Command Line"},
	}
	appsMutex.Unlock()

	// 1. Search with empty query returns all
	all := SearchApplications("")
	if len(all) < 3 {
		t.Errorf("expected at least 3 apps, got %d", len(all))
	}

	// 2. Search with acronym "gc" matches Google Chrome first
	results := SearchApplications("gc")
	if len(results) == 0 {
		t.Fatalf("expected results for query 'gc'")
	}
	if results[0].Name != "Google Chrome" {
		t.Errorf("expected top result to be 'Google Chrome', got %q", results[0].Name)
	}

	// 3. Search with "vsc" matches Visual Studio Code
	resultsVSC := SearchApplications("vsc")
	if len(resultsVSC) == 0 {
		t.Fatalf("expected results for query 'vsc'")
	}
	if resultsVSC[0].Name != "Visual Studio Code" {
		t.Errorf("expected top result to be 'Visual Studio Code', got %q", resultsVSC[0].Name)
	}
}
