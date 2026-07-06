# cfb-dynasty

Go library and CLI for reading EA Sports College Football dynasty save files on PC and exporting structured league data as JSON.

**Status: pre-alpha.** Parsing works against CFB 27 PC saves when paired with a schema bundle (for example `C27_441_0.gz`). Coverage is read-only and improves incrementally; some references and array tables are still incomplete.

## Install

```bash
go install github.com/leaguelines/cfb-dynasty/cmd/cfb-dynasty@latest
```

Or add the library to your project:

```bash
go get github.com/leaguelines/cfb-dynasty/dynasty
```

## Schema bundles

Table decoding requires a gzip-compressed schema file from the game assets (naming pattern like `C27_441_0.gz`). Point `--schema-dir` at the directory containing that file.

The schema major/minor version is read from the save header when possible; you can also place multiple bundles in the same directory and the closest match is selected.

## CLI

### Inspect a save

```bash
cfb-dynasty inspect /path/to/Dynasty1.sav
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1.sav
cfb-dynasty inspect -json /path/to/Dynasty1.sav
```

Shows compression, format, size, SHA-256, and table marker counts without full parsing.

### Export to JSON

```bash
# Full export (all sections)
cfb-dynasty export -schema-dir ./schemas /path/to/Dynasty1.sav -o dynasty.json

# Selective export — only the listed sections are included
cfb-dynasty export -schema-dir ./schemas --teams --rosters /path/to/Dynasty1.sav
cfb-dynasty export -schema-dir ./schemas --games --no-game-stats /path/to/Dynasty1.sav
```

#### Export sections

| Flag | JSON fields | Contents |
|------|-------------|----------|
| *(default)* | all below | Everything when no section flags are set |
| `--season` | `season` | Current year, week, phase |
| `--teams` | `teams` | Schools, records, poll ranks |
| `--rosters` | `rosters` | Active rosters with player ratings |
| `--games` | `games` | Schedule, scores, optional per-game stats |
| `--recruits` | `recruits` | Recruiting board + nested `player` attributes |
| `--recruiting` | `recruiting` | Pursuit state, NIL, visits, top school interest |
| `--season-stats` | `seasonPlayerStats`, `seasonTeamStats` | Season stat totals |
| `--coaches` | `coaches` | Staff, contracts, career records |
| `--leaving-players` | `leavingPlayers` | Graduation / exit pipeline |
| `--injuries` | `injuries` | Active injuries |
| `--depth-charts` | `depthCharts` | Depth chart slots by team |
| `--history` | `playerAwards`, `leagueAwards`, `conferenceChampions`, `statRecords` | Awards and record book |

Additional flags:

- `--no-game-stats` — omit team/player stat lines from game exports
- `-o path` — write JSON to a file (stdout if omitted)
- `--pretty=false` — compact JSON

Run `cfb-dynasty -h` or `cfb-dynasty export -h` for full usage.

### Example: recruits with `jq`

```bash
cfb-dynasty export -schema-dir ./schemas --recruits /path/to/Dynasty1.sav | \
  jq '.recruits[] | select(.nationalRank != null and .player != null) |
      {rank: .nationalRank, name: "\(.player.firstName) \(.player.lastName)",
       position: .player.position, overall: .player.overall, ratings: .player.ratings}'
```

## Library

```go
package main

import (
	"fmt"
	"log"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

func main() {
	settings := dynasty.DefaultSettings()
	settings.SchemaDir = "/path/to/schemas"
	settings.AutoParse = true

	file, err := dynasty.Open("/path/to/Dynasty1.sav", &settings)
	if err != nil {
		log.Fatal(err)
	}

	info, err := file.Inspect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("size=%d compressed=%v format=%q\n", info.Size, info.Compressed, info.Format)

	// Parse is called automatically when AutoParse is true.
	// Otherwise: if err := file.Parse(); err != nil { ... }

	export, err := file.ExportWithOptions(dynasty.ExportOptions{
		Sections: dynasty.ExportSections{
			Teams:   true,
			Rosters: true,
			Games:   true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	data, err := export.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exported %d bytes\n", len(data))
}
```

Lower-level access to parsed tables:

```go
team, ok := file.PrimaryTableByName("Team")
if ok {
	_ = team.ReadRecords()
	for _, row := range team.Records {
		_ = row.Get("LongName")
	}
}
```

### Package layout

| Package | Purpose |
|---------|---------|
| [`dynasty`](./dynasty/) | Public API — open saves, parse tables, export data |
| [`internal/binary`](./internal/binary/) | Low-level byte scanning helpers |
| [`internal/bitview`](./internal/bitview/) | Bitfield read helpers for record decoding |
| [`internal/compress`](./internal/compress/) | Decompression (zlib) |
| [`cmd/cfb-dynasty`](./cmd/cfb-dynasty/) | Command-line tool |

## Expected save location (PC)

```
%USERPROFILE%\Documents\College Football 27\Saves\
```

Dynasty saves are typically named like `Dynasty1.sav`.

## Known limitations

- **Schema required** — export and record decoding need a matching `C27_*_*.gz` bundle.
- **Read-only** — no save writing or editing.
- **Array tables** — some nested lists (for example full multi-school recruiting interest) do not decode yet.
- **Unplayed saves** — game and season player stats may be empty or sentinel values until games are simmed.
- **Large exports** — full exports with `--recruits` and `--rosters` can be tens of MB; use section flags to trim output.

## Goals

1. **Read-only** parsing of local PC dynasty saves.
2. Export season state, schedule, rosters, recruiting, stats, and league history as JSON.
3. Support headless pipelines: copy save → parse → JSON → sync to a web admin or bot.

## Non-goals (v1)

- Console save extraction
- Writing or repairing save files

## Related projects

| Project | Notes |
|---------|-------|
| [madden-franchise](https://github.com/bep713/madden-franchise) | Reference parser for Madden FranTk-style table DBs |
| [Madden Franchise Editor](https://github.com/bep713/madden-franchise-editor) | Desktop editor built on `madden-franchise` |

CFB shares the broad EA engine family with Madden, but dynasty saves are a distinct format with their own schema bundles.

## Development

```bash
go build ./...
go test ./...
go run ./cmd/cfb-dynasty export -schema-dir ./schemas --help
```

Place test saves and schema bundles in a local `data/` directory (not committed).

### Contributing

Issues and PRs welcome. Useful contributions include schema version coverage, additional export fields, array-table decoding, and regression tests against real saves.

## License

[MIT](./LICENSE)
