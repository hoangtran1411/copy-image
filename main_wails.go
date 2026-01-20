package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// Embed the frontend/dist directory into the binary.
// This allows us to ship a single executable file that contains
// all HTML, CSS, and JavaScript assets without needing external files.
//
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure.
	// This will be bound to the frontend, allowing JavaScript to call Go methods.
	app := NewApp()

	// Configure and run the Wails application.
	// These options control window appearance, behavior, and bindings.
	err := wails.Run(&options.App{
		Title:  "Copy Image Tool",
		Width:  900,
		Height: 700,

		// Prevent the window from being resized smaller than this.
		// This ensures the UI remains usable on smaller displays.
		MinWidth:  700,
		MinHeight: 550,

		// Asset server configuration - serves embedded frontend files.
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		// Set a dark background color to match our theme.
		// This prevents white flash during app startup.
		BackgroundColour: &options.RGBA{R: 15, G: 20, B: 25, A: 1},

		// Lifecycle hooks
		OnStartup: app.startup,

		// Bind Go structs to make their methods callable from JavaScript.
		// The App struct's exported methods become available as window.go.main.App.*
		Bind: []interface{}{
			app,
		},

		// Windows-specific options
		Windows: &windows.Options{
			// Disable transparency for better performance.
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,

			// Show the app icon in the title bar.
			DisableWindowIcon: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
