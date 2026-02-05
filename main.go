package main

import (
	"embed"
	_ "embed"
	"flag"
	"log"

	"github.com/maicek/laluer/core/apps"
	"github.com/maicek/laluer/core/handler"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	go apps.AppServiceInstance.Discover()
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	// application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// read environment variables

	floating := flag.Bool("floating", false, "float window")
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

	app := application.New(application.Options{
		Name:        "laluer",
		Description: "Laluer",
		Services: []application.Service{
			application.NewService(&handler.HandlerService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.

	isDev := *floating
	disableResize := !isDev
	frameless := !isDev
	alwaysOnTop := !isDev

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Laluer",
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		Width:            500,
		Height:           340,
		AlwaysOnTop:      alwaysOnTop,
		DisableResize:    disableResize,
		Frameless:        frameless,
		InitialPosition:  application.WindowCentered,
		URL:              "/",
		BackgroundType:   application.BackgroundTypeSolid,
		Linux: application.LinuxWindow{
			WindowIsTranslucent: false,
		},
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
