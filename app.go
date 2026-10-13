package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"beam/backend"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	clipboard *backend.ClipboardManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		clipboard: backend.GetClipboardManager(),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Scan applications in background at startup
	go backend.ScanApplications()
	// Start clipboard watcher
	a.clipboard.StartWatcher(ctx)
}

// GetApplications returns all detected Linux desktop applications
func (a *App) GetApplications() []backend.AppInfo {
	return backend.GetCachedApplications()
}

// SearchApps filters applications using fzf-style fuzzy scoring and frequency boost
func (a *App) SearchApps(query string) []backend.AppInfo {
	return backend.SearchApplications(query)
}

// LaunchApp launches the chosen application and tracks launch stats
func (a *App) LaunchApp(execCmd string, terminal bool, appID string) error {
	err := backend.Launch(execCmd, terminal, appID)
	if err == nil {
		// Hide window after successful launch
		runtime.WindowHide(a.ctx)
	}
	return err
}

// GetClipboard returns clipboard history
func (a *App) GetClipboard() []backend.ClipboardItem {
	return a.clipboard.GetHistory()
}

// CopyToClipboard writes text to the clipboard
func (a *App) CopyToClipboard(text string) error {
	err := backend.WriteClipboard(text)
	if err == nil {
		runtime.WindowHide(a.ctx)
	}
	return err
}

// EvaluateMath calculates math expressions inline
func (a *App) EvaluateMath(expr string) backend.CalcResult {
	return backend.Evaluate(expr)
}

// OpenURL opens a URL in the default system browser
func (a *App) OpenURL(url string) error {
	runtime.BrowserOpenURL(a.ctx, url)
	runtime.WindowHide(a.ctx)
	return nil
}

// HideWindow hides the Beam launcher window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}

// QuitApp closes Beam completely
func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

// LoadIconAsDataURL reads a local image/svg and converts to data URI
func (a *App) LoadIconAsDataURL(iconPath string) string {
	if iconPath == "" {
		return ""
	}

	data, err := os.ReadFile(iconPath)
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

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}
