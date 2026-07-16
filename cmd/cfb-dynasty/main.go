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
  cfb-dynasty recruiting-tunables [flags]
  cfb-dynasty schema-build [flags] <ftx-dir>

Commands:
  inspect   Show container metadata without full parsing
  export    Export parsed dynasty data as JSON
  recruiting-tunables  Dump recruiting formula constants from game tuning data
  schema-build  Build a C27_*.gz schema bundle from Frosty FTX extracts

Export flags (default: all sections):
  --games          Include schedule/scores
  --recruits       Include recruiting board + player attributes
  --season         Include current season metadata
  --teams          Include team records (record, polls, ranks)
  --rosters        Include team rosters + player attributes
  --recruiting     Include recruiting pursuit state + top school interest
  --season-stats   Include player and team season stat totals
  --coaches        Include coaching staff records
  --leaving-players Include offseason exit/graduation players
  --injuries       Include active player injuries
  --depth-charts   Include team depth charts
  --history        Include awards, conference champions, and stat records
  --school-grades  Include per-school recruiting pitch grades
  --pipelines      Include per-school pipeline influence
  --rivalries      Include rivalry records
  --position-changes Include player position change history
  --draft          Include draft pick assignments
  --bowls          Include bowl game metadata
  --no-game-stats  Omit team/player stat lines from game exports

When one or more section flags is set, only those sections are exported.
With no section flags, everything is included.

Export also accepts:
  --tuning-path    Path to dynasty-tuning-binary.FTC (skill group labels; auto-discovered under --schema-dir when omitted)

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
	case "recruiting-tunables":
		os.Exit(runRecruitingTunables(os.Args[2:]))
	case "schema-build":
		os.Exit(runSchemaBuild(os.Args[2:]))
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
	schemaDir := fs.String("schema-dir", "", "directory containing gzip schema bundles")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	path := fs.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "inspect: save file path required")
		return 2
	}

	settings := dynasty.DefaultSettings()
	settings.SchemaDir = *schemaDir
	file, err := dynasty.Open(path, &settings)
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
	if err := file.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "inspect: parse: %v\n", err)
		return 1
	}
	fmt.Printf("parsed tables:  %d\n", len(file.Tables()))
	if *schemaDir != "" {
		if schema := file.Schema(); schema != nil {
			fmt.Printf("loaded schema:  major=%d minor=%d gameYear=%d tables=%d\n",
				schema.Version.Major, schema.Version.Minor, schema.GameYear, len(schema.Tables))
			fmt.Printf("schema file:    %s\n", schema.Version.Path)
		}
	}
	return 0
}

func runExport(args []string) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	output := fs.String("o", "", "write JSON to file (default: stdout)")
	pretty := fs.Bool("pretty", true, "indent JSON output")
	schemaDir := fs.String("schema-dir", "", "directory containing gzip schema bundles")
	tuningPath := fs.String("tuning-path", "", "path to dynasty-tuning-binary.FTC for skill group labels (auto-discovered under schema-dir when omitted)")
	includeGames := fs.Bool("games", false, "include games (schedule/scores)")
	includeRecruits := fs.Bool("recruits", false, "include recruits and player attributes")
	includeSeason := fs.Bool("season", false, "include season metadata")
	includeTeams := fs.Bool("teams", false, "include team records")
	includeRosters := fs.Bool("rosters", false, "include team rosters")
	includeRecruiting := fs.Bool("recruiting", false, "include recruiting pursuit data")
	includeSeasonStats := fs.Bool("season-stats", false, "include season stat totals")
	includeCoaches := fs.Bool("coaches", false, "include coaching staff")
	includeLeavingPlayers := fs.Bool("leaving-players", false, "include offseason exit players")
	includeInjuries := fs.Bool("injuries", false, "include active injuries")
	includeDepthCharts := fs.Bool("depth-charts", false, "include team depth charts")
	includeHistory := fs.Bool("history", false, "include awards and league history")
	includeSchoolGrades := fs.Bool("school-grades", false, "include school recruiting grades")
	includePipelines := fs.Bool("pipelines", false, "include school pipeline influence")
	includeRivalries := fs.Bool("rivalries", false, "include rivalry records")
	includePositionChanges := fs.Bool("position-changes", false, "include player position change history")
	includeDraft := fs.Bool("draft", false, "include draft pick assignments")
	includeBowls := fs.Bool("bowls", false, "include bowl game metadata")
	omitGameStats := fs.Bool("no-game-stats", false, "omit team/player stat lines from games")
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
	settings.SchemaDir = *schemaDir
	settings.TuningPath = *tuningPath
	file, err := dynasty.Open(path, &settings)
	if err != nil {
		if dynasty.IsNotImplemented(err) {
			fmt.Fprintf(os.Stderr, "export: parser not ready yet: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "export: %v\n", err)
		return 1
	}

	exportOpts := dynasty.ExportOptions{
		OmitGameStats: *omitGameStats,
	}
	if *includeGames || *includeRecruits || *includeSeason || *includeTeams || *includeRosters ||
		*includeRecruiting || *includeSeasonStats || *includeCoaches || *includeLeavingPlayers ||
		*includeInjuries || *includeDepthCharts || *includeHistory || *includeSchoolGrades ||
		*includePipelines || *includeRivalries || *includePositionChanges || *includeDraft ||
		*includeBowls {
		exportOpts.Sections = dynasty.ExportSections{
			Games:            *includeGames,
			Recruits:         *includeRecruits,
			Season:           *includeSeason,
			Teams:            *includeTeams,
			Rosters:          *includeRosters,
			Recruiting:       *includeRecruiting,
			SeasonStats:      *includeSeasonStats,
			Coaches:          *includeCoaches,
			LeavingPlayers:   *includeLeavingPlayers,
			Injuries:         *includeInjuries,
			DepthCharts:      *includeDepthCharts,
			History:          *includeHistory,
			SchoolGrades:     *includeSchoolGrades,
			Pipelines:        *includePipelines,
			Rivalries:        *includeRivalries,
			PositionChanges:  *includePositionChanges,
			Draft:            *includeDraft,
			Bowls:            *includeBowls,
		}
	}

	export, err := file.ExportWithOptions(exportOpts)
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

func runRecruitingTunables(args []string) int {
	fs := flag.NewFlagSet("recruiting-tunables", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	output := fs.String("o", "", "write JSON to file (default: stdout)")
	schemaDir := fs.String("schema-dir", "", "directory containing schema bundle and optional cfb27-db-data tuning FTC")
	tuningPath := fs.String("tuning", "", "path to dynasty-tuning-binary.FTC (auto-discovered under schema-dir when omitted)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	data, err := dynasty.RecruitingTunablesJSON(*schemaDir, *tuningPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recruiting-tunables: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "recruiting-tunables: write: %v\n", err)
		return 1
	}
	return 0
}

func runSchemaBuild(args []string) int {
	fs := flag.NewFlagSet("schema-build", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	outDir := fs.String("o", "", "output directory for C{year}_{major}_{minor}.gz (default: current directory)")
	major := fs.Int("major", -1, "schema major override (default: franchise-schemas dataMajorVersion, else 468 for CFB 27)")
	minor := fs.Int("minor", -1, "schema minor override (default: franchise-schemas dataMinorVersion, else dataRevisionVersion)")
	year := fs.Int("year", -1, "game year override (default: from franchise databaseName / path cfbNN, else 27)")
	noExtras := fs.Bool("no-extras", false, "omit madden-franchise core extra schemas")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	source := fs.Arg(0)
	if source == "" {
		fmt.Fprintln(os.Stderr, "schema-build: path to cfb27-db-data or a patch folder required")
		fmt.Fprintln(os.Stderr, "example: cfb-dynasty schema-build -o ./data ./data/cfb27-db-data/2")
		return 2
	}

	opts := dynasty.SchemaBuildOptions{Source: source, OutDir: *outDir}
	if *major >= 0 {
		opts.Major = major
	}
	if *minor >= 0 {
		opts.Minor = minor
	}
	if *year >= 0 {
		opts.GameYear = year
	}
	includeExtras := !*noExtras
	opts.IncludeExtras = &includeExtras

	result, err := dynasty.BuildSchemaBundle(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "schema-build: %v\n", err)
		return 1
	}
	fmt.Printf("wrote %s\n", result.Path)
	fmt.Printf("meta:   major=%d (%s) minor=%d (%s) gameYear=%d\n",
		result.Major, result.MajorSource, result.Minor, result.MinorSource, result.GameYear)
	fmt.Printf("tables: %d  enums: %d\n", result.Tables, result.Enums)
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
