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

Table decoding requires a **gzip-compressed JSON schema bundle** derived from the game install. Each bundle lists thousands of table definitions (field names, types, enums) used by the FranTk-style save format.

### File format and naming

Bundles use the same JSON layout as [madden-franchise](https://github.com/bep713/madden-franchise) `.gz` schemas:

```json
{
  "meta": { "major": 441, "minor": 0, "gameYear": 27 },
  "schemas": [ ... ],
  "schemaMap": { ... }
}
```

Place one or more files in a directory and pass it to `--schema-dir`. Recognized filenames:

| Pattern | Example | Meaning |
|---------|---------|---------|
| `C{year}_{major}_{minor}.gz` | `C27_441_0.gz` | CFB, game year 27, bundle major 441, minor 0 |
| `{major}_{minor}.gz` | `441_0.gz` | Major/minor only (game year optional) |
| `M{year}_{major}_{minor}.gz` | `M27_441_0.gz` | Madden-style prefix also accepted |

`inspect` prints the schema version embedded in your save and which bundle was loaded:

```bash
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1.sav
# schema:         major=809 minor=1      ← from save header
# loaded schema:  major=441 minor=0 ...  ← picked bundle
```

The save header version and the bundle `meta` version use **different numbering** (for example save `809.1` vs bundle `441.0`). When several bundles are present, the loader picks the closest major/minor match; with only one file in the directory, that file is used.

### How to obtain a bundle

You need a **legal copy of the game** on PC. Schema data lives inside the install assets; it is not shipped with this repository.

**1. Extract raw schema assets (Frosty)**

The Madden modding workflow applies to CFB as well:

1. Install [Frosty Editor](https://frostytoolsuite.com/) and point it at your College Football 27 install.
2. Open **Legacy Explorer** and search for franchise / FranTk schema assets (Madden uses `franchise-schemas.ftx` under `common → franchise`; CFB paths are still being mapped — see internal RE notes).
3. Export the schema `.ftx` / `.xml` files to a folder on disk.

**2. Convert to gzip JSON**

The CLI expects the **evaluated** `.gz` bundle, not raw `.ftx`. The usual path is the madden-franchise schema generator (same engine family):

- Use [`madden-franchise`](https://github.com/bep713/madden-franchise) `FranchiseSchema` / `schemaGenerator` to load the extracted `.ftx` and emit a `.gz` JSON bundle, **or**
- Use tooling from the [Madden Franchise Editor](https://github.com/bep713/madden-franchise-editor) ecosystem (`schemaSearchService` can pull schemas from CAS/LZ4 game assets).

**3. Name and verify**

Rename the output to match your game year and bundle `meta` (for example `C27_441_0.gz`), then confirm decoding works:

```bash
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1.sav
cfb-dynasty export -schema-dir ./schemas --teams /path/to/Dynasty1.sav | head
```

If you already have a working bundle (like `C27_441_0.gz` from local RE), you can use it directly — no re-extraction needed until EA ships a patch that changes table layouts.

### Should schema bundles be committed to git?

**No — do not commit them to a public repository.**

| Concern | Notes |
|---------|-------|
| **Copyright** | Bundles are derived from EA game assets (table/field names, types, enums). Redistributing them is likely outside EA's terms, even though the files are not executable game content. |
| **Size** | ~3 MB compressed, ~30 MB inflated — poor fit for git history. |
| **Patch churn** | EA title updates can change schema major/minor; bundles go stale quickly. |

**Recommended approach:** keep bundles in a local `data/` or `schemas/` directory (gitignored), document extraction steps (this section), and let each developer generate or copy their own from an owned game install. Integration tests already **skip** when `data/C27_441_0.gz` is absent.

For private teams, a shared drive or internal artifact bucket is fine; just avoid publishing the files in the open-source repo or release tarballs.

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
| `--history` | `playerAwards`, `leagueAwards`, `conferenceChampions`, `recordBook` | Awards and the full stat record book |

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
       position: .player.position, archetype: .player.archetypeLabel, overall: .player.overall,
       skillCaps: .player.skillGroupCaps, skillCapTotal: .player.skillGroupCapTotal}'
```

### Skill group caps

Every player (including a recruit's linked `player`) exports its six skill-group
caps read straight from the save:

- `skillGroupCaps` — the six positional caps (`SkillGroupCap1..6`).
- `skillGroupCapTotal` — their sum (a strong proxy for a recruit's ceiling / star tier).

The caps are exported as **opaque, positional values** — the array is ordered by
`SkillGroupCap1..6` as stored in the save. The game buckets these into six
position-specific skill groups, but the group **names** are not present in the
dynasty save (they live in the tuning FTC, which is not yet extractable for
CFB 27). Rather than guess at labels, the export intentionally leaves the slots
unnamed until definitive names are available. Note also that only the per-group
*cap* is known — the individual ratings inside each group are tuning-driven and
not in the save.

### Per-game player stats

When game stats are included (`--games`, on by default), each game carries a
`playerGameStats` list. Each entry is one player's line for that game:

- `playerId` — the player's row index (join to `rosters` for the full record).
- `player` — a lightweight identity (`firstName`, `lastName`, `position`,
  `jersey`, `teamIndex`) so lines are readable without a join.
- `offense` / `defense` — the stat line(s); a player with both merges into one entry.

The game-stat rows themselves store no player reference — ownership lives on the
`Player` side via a `GameStats[]` array store that points at each player's rows.
The exporter inverts those arrays to attribute every stat line, and rows are
bucketed into games by their direct `SeasonGame` record index (stale references
to other seasons are dropped).

### Record book

`--history` exports the complete stat record book as a flat `recordBook` list.
The game keeps a record book for three **scopes** — the whole FBS (`league`),
each `conference`, and every `team` — and three **periods** (`career`, `season`,
`game`). Each entry carries:

- `scope` — `league`, `conference`, or `team`.
- `scopeName` — the conference or team the board belongs to (empty for league).
- `period` — `career`, `season`, or `game`.
- `statType`, `statValue`, `rank` — the record and its rank within the board.
- `firstName`, `lastName`, `position`, `teamName`, `calendarYear` — the holder.

League boards are ranked top-N per category; conference and team boards store a
single holder (`rank` 1) per category. Positions come straight from the record
row when present and are otherwise inferred from the stat category (the game
leaves lower league ranks at the schema default).

```bash
# Every team's career passing-yards record holder
cfb-dynasty export -schema-dir ./schemas --history /path/to/Dynasty1.sav | \
  jq '.recordBook[] | select(.scope=="team" and .period=="career" and .statType=="PassYards") |
      {team: .scopeName, holder: "\(.firstName) \(.lastName)", yards: .statValue}'
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

Place test saves and schema bundles in a local `data/` directory. Dynasty saves may contain personal league data; schema bundles are game-derived assets — **neither should be committed** (see [Schema bundles](#schema-bundles)).

### Contributing

Issues and PRs welcome. Useful contributions include schema version coverage, additional export fields, array-table decoding, and regression tests against real saves.

## License

[MIT](./LICENSE)
