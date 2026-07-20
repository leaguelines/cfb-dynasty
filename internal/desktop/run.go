package desktop

import (
	"fmt"
	"io/fs"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Options configures the desktop GUI launch.
type Options struct {
	SchemaDir string
	SavePath  string
	Assets    fs.FS
}

// Run starts the Wails desktop application.
func Run(opts Options) error {
	app := NewApp().WithFlags(opts.SchemaDir, opts.SavePath)
	err := wails.Run(&options.App{
		Title:  "CFB Dynasty Explorer",
		Width:  1280,
		Height: 860,
		AssetServer: &assetserver.Options{
			Assets: opts.Assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 20, B: 32, A: 255},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		return fmt.Errorf("desktop gui: %w", err)
	}
	return nil
}
