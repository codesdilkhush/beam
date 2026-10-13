package main

import (
	"context"

	"beam/backend"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	clipboard *backend.ClipboardManager
	isVisible bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		clipboard: backend.GetClipboardManager(),
		isVisible: true,
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.isVisible = true

	// Scan applications in background at startup
	go backend.ScanApplications()
	// Start clipboard watcher
	a.clipboard.StartWatcher(ctx)

	// Start single-instance socket daemon
	sockPath := backend.GetSocketPath()
	if err := backend.StartDaemonServer(ctx, sockPath, func(cmd string) {
		if cmd == "toggle" {
			if a.isVisible {
				runtime.WindowHide(a.ctx)
				a.isVisible = false
			} else {
				runtime.WindowShow(a.ctx)
				runtime.WindowCenter(a.ctx)
				a.isVisible = true
			}
		} else if cmd == "show" {
			runtime.WindowShow(a.ctx)
			runtime.WindowCenter(a.ctx)
			a.isVisible = true
		} else if cmd == "hide" {
			runtime.WindowHide(a.ctx)
			a.isVisible = false
		}
	}); err != nil {
		runtime.LogErrorf(ctx, "daemon server failed: %s", err.Error())
	}
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
		a.isVisible = false
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
		a.isVisible = false
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
	a.isVisible = false
	runtime.WindowHide(a.ctx)
	return nil
}

// GetScripts discovers Raycast script commands in ~/.config/beam/scripts/
func (a *App) GetScripts() []backend.ScriptCommand {
	return backend.DiscoverScripts()
}

// ExecuteScript runs a chosen script command
func (a *App) ExecuteScript(scriptPath string) (string, error) {
	return backend.RunScript(scriptPath)
}

// HideWindow hides the Beam launcher window
func (a *App) HideWindow() {
	a.isVisible = false
	runtime.WindowHide(a.ctx)
}

// QuitApp closes Beam completely
func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

// LoadIconAsDataURL reads a local image/svg and converts to data URI using the high-speed icon cache
func (a *App) LoadIconAsDataURL(iconPath string) string {
	return backend.GetIconCache().GetOrConvert(iconPath)
}

// GetTodos returns all To-Do items from SQLite
func (a *App) GetTodos() []backend.TodoItem {
	return backend.GetTodoManager().GetTodos()
}

// AddTodo creates a new To-Do item in SQLite
func (a *App) AddTodo(title string) (*backend.TodoItem, error) {
	return backend.GetTodoManager().AddTodo(title)
}

// ToggleTodo toggles completion state of a To-Do item
func (a *App) ToggleTodo(id string) error {
	return backend.GetTodoManager().ToggleTodo(id)
}

// DeleteTodo removes a To-Do item
func (a *App) DeleteTodo(id string) error {
	return backend.GetTodoManager().DeleteTodo(id)
}

// ClearCompletedTodos removes all completed tasks
func (a *App) ClearCompletedTodos() error {
	return backend.GetTodoManager().ClearCompleted()
}

// GetEmojis returns the full emoji catalog
func (a *App) GetEmojis() []backend.EmojiItem {
	return backend.GetEmojiManager().GetAll()
}

// SearchEmojis filters emojis with fuzzy matching
func (a *App) SearchEmojis(query string) []backend.EmojiItem {
	return backend.GetEmojiManager().SearchEmojis(query)
}

// ConvertQuery evaluates unit, currency, color, time, and CSS conversions
func (a *App) ConvertQuery(query string) backend.ConvertResult {
	return backend.Convert(query)
}

// GetGuideTopics returns the comprehensive in-app user guide catalog
func (a *App) GetGuideTopics() []backend.GuideTopic {
	return backend.GetGuideTopics()
}

// SearchGuideTopics filters user guide topics by keyword
func (a *App) SearchGuideTopics(query string) []backend.GuideTopic {
	return backend.SearchGuideTopics(query)
}
