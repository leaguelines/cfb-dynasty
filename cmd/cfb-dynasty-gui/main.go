package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/leaguelines/cfb-dynasty/internal/desktop"
)

//go:embed all:frontend/dist
var embeddedAssets embed.FS

func main() {
	fsFlags := flag.NewFlagSet("cfb-dynasty-gui", flag.ExitOnError)
	schemaDir := fsFlags.String("schema-dir", "", "directory containing C27_*.gz schema bundles")
	savePath := fsFlags.String("save", "", "optional dynasty save to open on launch")
	_ = fsFlags.Parse(os.Args[1:])

	assets, err := fs.Sub(embeddedAssets, "frontend/dist")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cfb-dynasty-gui: embed assets: %v\n", err)
		os.Exit(1)
	}

	if err := desktop.Run(desktop.Options{
		SchemaDir: *schemaDir,
		SavePath:  *savePath,
		Assets:    assets,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "cfb-dynasty-gui: %v\n", err)
		os.Exit(1)
	}
}
