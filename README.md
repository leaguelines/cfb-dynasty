# cfb-dynasty

Go library and CLI for reading EA Sports College Football dynasty save files on PC and exporting structured league data.

**Status: pre-alpha.** The save format is still being investigated. The library API and CLI are scaffolded so parsing work can land incrementally before [College Football 27](https://www.ea.com/games/ea-sports-college-football/college-football-27/faq) releases on PC.

## Install

```bash
go install github.com/leaguelines/cfb-dynasty/cmd/cfb-dynasty@latest
```

Or add the library to your project:

```bash
go get github.com/leaguelines/cfb-dynasty/dynasty
```

## CLI

```bash
# Show container metadata (compression, magic bytes, size)
cfb-dynasty inspect /path/to/Dynasty1.sav

# Export parsed dynasty data as JSON (requires a completed parser)
cfb-dynasty export /path/to/Dynasty1.sav -o dynasty.json
```

Run `cfb-dynasty -h` for full usage.

## Library

```go
package main

import (
	"fmt"
	"log"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

func main() {
	file, err := dynasty.Open("/path/to/Dynasty1.sav", nil)
	if err != nil {
		log.Fatal(err)
	}

	info, err := file.Inspect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("size=%d compressed=%v format=%q\n", info.Size, info.Compressed, info.Format)

	if err := file.Parse(); err != nil {
		log.Fatal(err) // ErrNotImplemented until the format is mapped
	}

	export, err := file.Export()
	if err != nil {
		log.Fatal(err)
	}
	_ = export
}
```

### Package layout

| Package | Purpose |
|---------|---------|
| [`dynasty`](./dynasty/) | Public API — open saves, parse tables, export data |
| [`internal/binary`](./internal/binary/) | Low-level byte scanning helpers |
| [`internal/compress`](./internal/compress/) | Decompression (zlib and future codecs) |
| [`cmd/cfb-dynasty`](./cmd/cfb-dynasty/) | Command-line tool |

## Expected save location (PC)

```
%USERPROFILE%\Documents\College Football 27\Saves\
```

Dynasty saves are typically named like `Dynasty1.sav` (confirm after launch).

## Goals

1. **Read-only** parsing of local PC dynasty saves.
2. Export season state, schedule, scores, team records, and optionally rosters/stats.
3. Support headless pipelines: copy save → parse → JSON → sync to web admin or bots.

## Non-goals (v1)

- Console save extraction
- Writing or repairing save files

## Related projects

| Project | Notes |
|---------|-------|
| [madden-franchise](https://github.com/bep713/madden-franchise) | Reference parser for Madden FranTk-style table DBs |
| [Madden Franchise Editor](https://github.com/bep713/madden-franchise-editor) | Desktop editor built on `madden-franchise` |

CFB shares the broad EA engine family with Madden, but dynasty saves are a distinct format and need their own schema work.

## Development

```bash
go build ./...
go test ./...
go run ./cmd/cfb-dynasty inspect -h
```

### Contributing

Issues and PRs welcome once the format is partially understood. Until then, the focus is mapping container layout, table markers, and schemas from real save files.

## License

[MIT](./LICENSE)
