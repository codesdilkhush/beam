package main

import (
	"embed"

	"beam/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Single-instance check: If already running, signal existing instance to toggle window and exit
	sockPath := backend.GetSocketPath()
	if backend.TrySignalExistingInstance(sockPath) {
		return
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "Beam",
		Width:             760,
		Height:            500,
		Frameless:         true,
		AlwaysOnTop:       true,
		DisableResize:     true,
		StartHidden:       false,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 16, G: 18, B: 24, A: 215},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon:                nil,
			WindowIsTranslucent: true,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
