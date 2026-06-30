package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

const usage = `cfb-dynasty — read College Football dynasty saves and export league data

Usage:
  cfb-dynasty inspect [flags] <save-file>
  cfb-dynasty export  [flags] <save-file>

Commands:
  inspect   Show container metadata without full parsing
  export    Export parsed dynasty data as JSON

Global flags:
  -h        Show help
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "inspect":
		os.Exit(runInspect(os.Args[2:]))
	case "export":
		os.Exit(runExport(os.Args[2:]))
	case "-h", "--help", "help":
		fmt.Print(usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runInspect(args []string) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	jsonOut := fs.Bool("json", false, "print inspection as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path := fs.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "inspect: save file path required")
		return 2
	}

	file, err := dynasty.Open(path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "inspect: %v\n", err)
		return 1
	}

	info, err := file.Inspect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "inspect: %v\n", err)
		return 1
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(info); err != nil {
			fmt.Fprintf(os.Stderr, "inspect: encode: %v\n", err)
			return 1
		}
		return 0
	}

	printInspect(info)
	return 0
}

func runExport(args []string) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	output := fs.String("o", "", "write JSON to file (default: stdout)")
	pretty := fs.Bool("pretty", true, "indent JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path := fs.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "export: save file path required")
		return 2
	}

	settings := dynasty.DefaultSettings()
	settings.AutoParse = true
	file, err := dynasty.Open(path, &settings)
	if err != nil {
		if dynasty.IsNotImplemented(err) {
			fmt.Fprintf(os.Stderr, "export: parser not ready yet: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "export: %v\n", err)
		return 1
	}

	export, err := file.Export()
	if err != nil {
		if dynasty.IsNotImplemented(err) {
			fmt.Fprintf(os.Stderr, "export: parser not ready yet: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "export: %v\n", err)
		return 1
	}

	var data []byte
	if *pretty {
		data, err = export.ToJSON()
	} else {
		data, err = export.MarshalJSON()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "export: encode: %v\n", err)
		return 1
	}

	if *output == "" {
		os.Stdout.Write(data)
		if len(data) == 0 || data[len(data)-1] != '\n' {
			fmt.Println()
		}
		return 0
	}

	if err := os.WriteFile(*output, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "export: write: %v\n", err)
		return 1
	}
	return 0
}

func printInspect(info dynasty.FileInfo) {
	fmt.Printf("path:           %s\n", info.Path)
	fmt.Printf("size:           %d bytes\n", info.Size)
	fmt.Printf("sha256:         %s\n", info.SHA256)
	fmt.Printf("compressed:     %v\n", info.Compressed)
	fmt.Printf("format:         %s\n", info.Format)
	fmt.Printf("game year:      %d\n", info.GameYear)
	if info.ExpectedSchema != nil {
		fmt.Printf("schema:         major=%d minor=%d\n", info.ExpectedSchema.Major, info.ExpectedSchema.Minor)
	} else {
		fmt.Printf("schema:         (unknown)\n")
	}
	fmt.Printf("unpacked size:  %d bytes\n", info.UnpackedSize)
	fmt.Printf("header (hex):   %s\n", formatHex(info.HeaderHex))
	if len(info.TableMarkerCounts) > 0 {
		fmt.Println("table markers:")
		for _, marker := range dynasty.AllTableMarkers {
			fmt.Printf("  %s: %d\n", marker, info.TableMarkerCounts[marker])
		}
	}
}

func formatHex(raw string) string {
	if raw == "" {
		return ""
	}
	var b strings.Builder
	for i := 0; i < len(raw); i += 2 {
		if i > 0 {
			b.WriteByte(' ')
		}
		end := i + 2
		if end > len(raw) {
			end = len(raw)
		}
		b.WriteString(raw[i:end])
	}
	return b.String()
}
