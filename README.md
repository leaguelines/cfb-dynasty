# cfb-dynasty

Go library and CLI for reading EA Sports College Football dynasty save files on PC and exporting structured league data as JSON.

**Status: feature-complete for read-only export.** Against CFB 27 PC dynasty saves (with a matching schema bundle such as `C27_468_2.gz`), the tool covers the main dynasty surfaces: season/teams/rosters, schedule and per-game stats, recruiting, coaches, injuries, depth charts, awards, and the full record book. Remaining gaps are narrow (see [Known limitations](#known-limitations)); writing or editing saves is intentionally out of scope.

## Install

```bash
go install github.com/leaguelines/cfb-dynasty/cmd/cfb-dynasty@latest
```

Or add the library to your project:

```bash
go get github.com/leaguelines/cfb-dynasty/dynasty
```

## Desktop GUI

A cross-platform desktop explorer (`cfb-dynasty-gui`) browses dynasty saves with the same collections as the web app and exports **JSON** or **CSV** (including a full CSV zip). It uses [Wails](https://wails.io) (Go + HTML/CSS bindings) — no local web server.

Schema bundles are **not** shipped in releases. On first launch, choose the folder that contains your `C27_*.gz` schema, then open a save. The app remembers the schema path. On Windows, saves under `Documents\EA SPORTS CFB27\saves` are listed automatically.

```bash
# Dev / local run (Wails requires the production or dev build tag)
go run -tags production ./cmd/cfb-dynasty-gui --schema-dir ./data/schemas

# Install GUI binary
go build -tags production -o cfb-dynasty-gui ./cmd/cfb-dynasty-gui

# Optional: launch via the CLI helper once the GUI binary is on PATH
cfb-dynasty gui --schema-dir ./data/schemas

# Cross-platform release packages (install the Wails CLI first)
cd cmd/cfb-dynasty-gui && wails build
```

Requires platform WebView dependencies (WebView2 on Windows; WebKitGTK on Linux; macOS includes WebKit). See the [Wails docs](https://wails.io/docs/gettingstarted/installation).

### Releases

GitHub Actions publishes builds automatically:

| Trigger | Release |
|---------|---------|
| Push to `main` | Updates the **nightly** prerelease (overwritten each time) |
| Tag `v*` (e.g. `v0.3.0`) | Creates a versioned release |

Each release includes GUI packages for Linux, Windows, and macOS plus matching CLI binaries. Schema bundles are never attached.

## What exports today

| Area | Coverage |
|------|----------|
| Season / teams | Year, week, phase; schools with conference, W–L, OFF/DEF/OVR ratings, poll ranks, recruiting class ranks/score |
| Rosters | Active players with ratings, archetypes, skill-group caps and labels (when tuning data is available) |
| Games | Schedule, scores, team box scores, attributed player lines (offense / defense / special teams) |
| Season stats | Player and team season totals |
| Recruiting | Board + player attributes; pursuit state, NIL, visits, full school interest, active pitches |
| Staff / roster mgmt | Coaches, injuries, depth charts, leaving / graduation pipeline |
| History | Player awards, league awards, conference champions, full record book (league / conference / team × career / season / game) |

Stable `Team.TeamIndex` IDs are used throughout for joins (see [Team IDs](#team-ids)).

## Schema bundles

Table decoding requires a **gzip-compressed JSON schema bundle** derived from the game install. Each bundle lists thousands of table definitions (field names, types, enums) used by the FranTk-style save format.

### File format and naming

Bundles use the same JSON layout as [madden-franchise](https://github.com/bep713/madden-franchise) `.gz` schemas:

```json
{
  "meta": { "major": 468, "minor": 2, "gameYear": 27 },
  "schemas": [ ... ],
  "schemaMap": { ... }
}
```

Place one or more files in a directory and pass it to `--schema-dir`. Recognized filenames:

| Pattern | Example | Meaning |
|---------|---------|---------|
| `C{year}_{major}_{minor}.gz` | `C27_468_2.gz` | CFB, game year 27, bundle major 468, minor 2 |
| `{major}_{minor}.gz` | `468_2.gz` | Major/minor only (game year optional) |
| `M{year}_{major}_{minor}.gz` | `M27_468_2.gz` | Madden-style prefix also accepted |

`inspect` prints the schema version embedded in your save and which bundle was loaded:

```bash
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1
# schema:         major=809 minor=1      ← from save header
# loaded schema:  major=468 minor=2 ...  ← picked bundle
```

The save header version and the bundle `meta` version use **different numbering** (for example save `809.1` vs bundle `468.2`). When several bundles are present, the loader picks the closest major/minor match; with only one file in the directory, that file is used. Prefer the newest CFB 27 bundle you have — older majors (for example `441.0`) mis-decode team and roster fields.

### How to obtain a bundle

You need a **legal copy of the game** on PC. Schema data lives inside the install assets; it is not shipped with this repository.

**1. Extract FranTk schema assets (MMC Frosty)**

1. Install [MMC Frosty Modding Tools](https://github.com/bphit4/MMC-Frosty-Modding-Tools/releases) and point it at your College Football 27 install.
2. Open **Legacy Explorer** and export the dynasty schema tree (the layout this tool expects looks like `cfb27-db-data/<patch>/` with `core/`, `football/`, and `franchise/` full of `.FTX` files, plus optional `.FTC` / `.FTB` siblings).
3. Keep the numbered patch folders intact (`0`, `1`, `2`, …). Folder `0` and `2` can be identical; content revision is what matters (see below).

**2. Build the gzip bundle with this tool**

```bash
# Point at the parent extract — picks the newest dataRevisionVersion patch
cfb-dynasty schema-build -o ./schemas ./data/cfb27-db-data

# Or build a specific patch folder
cfb-dynasty schema-build -o ./schemas ./data/cfb27-db-data/2
```

That writes `C{year}_{major}_{minor}.gz` (for example `C27_468_2.gz`) by evaluating every `.FTX` in the tree (plus a small set of FranTk core extras, same idea as [madden-franchise](https://github.com/bep713/madden-franchise)).

| Flag | Purpose |
|------|---------|
| `-o` | Output directory for the `.gz` (default: current directory) |
| `-major` | Override schema major (see version notes) |
| `-minor` | Override schema minor |
| `-year` | Override game year (default: from `cfbNN` in the path, else `27`) |
| `-no-extras` | Skip embedded FranTk core extra tables |

**Version metadata**

| Field | How `schema-build` chooses it |
|-------|-------------------------------|
| **Major / minor** | From the **College/franchise** root FTX when present (`franchise-schemas.FTX`, namespace `FranTk.College`) — its `dataMajorVersion` / `dataMinorVersion`. Ignores Core (`55.x`) and Football roots, which ship their own majors. |
| **Major (fallback)** | Split extracts without root meta: for game year 27 defaults to **468**. Pass `-major` to override. |
| **Minor (fallback)** | Majority `dataRevisionVersion` on per-table `.FTX` headers, then a numeric patch directory name, then `-minor`. |
| **Game year** | `databaseName` on the franchise root (e.g. `CollegeFB27_…`), else `cfbNN` in the source path, else `-year`, else `27`. |

The save header still uses a **different** numbering scheme (for example `809.1`); that is unrelated to the bundle `meta` major/minor.

**3. Verify**

```bash
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1
cfb-dynasty export -schema-dir ./schemas --teams /path/to/Dynasty1 | head
```

Rebuild after EA patches that change table layouts. If you already have a working `C27_*.gz`, you can keep using it without re-extracting.

**Alternatives:** [`madden-franchise`](https://github.com/bep713/madden-franchise) `schemaGenerator` / Franchise Editor schema search can also produce compatible `.gz` files.

**4. Optional: tuning FTC for skill group data**

Skill group bucket names, attribute membership, tier metadata used for current
levels, and recruiting formula constants live in `dynasty-tuning-binary.FTC`
from the same MMC Frosty extract. Keep it under your schema directory as:

```
schemas/cfb27-db-data/<patch>/dynasty-tuning-binary.FTC
```

The exporter auto-discovers the newest patch folder. Pass `--tuning-path` to
override. Without this file, skill-group capped/unlocked counts still export in
save-slot order, but labels, attributes, and `skillGroupCurrentLevels` are
omitted because group membership cannot be determined.

### Should schema bundles be committed to git?

**No — do not commit them to a public repository.**

| Concern | Notes |
|---------|-------|
| **Copyright** | Bundles are derived from EA game assets (table/field names, types, enums). Redistributing them is likely outside EA's terms, even though the files are not executable game content. |
| **Size** | ~3 MB compressed, ~30 MB inflated — poor fit for git history. |
| **Patch churn** | EA title updates can change schema major/minor; bundles go stale quickly. |

**Recommended approach:** keep bundles in a local `data/` or `schemas/` directory (gitignored), document extraction steps (this section), and let each developer generate or copy their own from an owned game install. Integration tests already **skip** when a local schema bundle / test save is absent.

For private teams, a shared drive or internal artifact bucket is fine; just avoid publishing the files in the open-source repo or release tarballs.

## CLI

### Inspect a save

```bash
cfb-dynasty inspect /path/to/Dynasty1
cfb-dynasty inspect -schema-dir ./schemas /path/to/Dynasty1
cfb-dynasty inspect -json /path/to/Dynasty1
```

Shows compression, format, size, SHA-256, and table marker counts without full parsing.

### Build a schema bundle from FTX extracts

```bash
cfb-dynasty schema-build -o ./schemas ./data/cfb27-db-data
cfb-dynasty schema-build -o ./schemas -major 468 ./data/cfb27-db-data/2
```

See [How to obtain a bundle](#how-to-obtain-a-bundle) for MMC Frosty extraction layout and version detection.

### Export to JSON

```bash
# Full export (all sections)
cfb-dynasty export -schema-dir ./schemas /path/to/Dynasty1 -o dynasty.json

# Selective export — only the listed sections are included
cfb-dynasty export -schema-dir ./schemas --teams --rosters /path/to/Dynasty1
cfb-dynasty export -schema-dir ./schemas --games --no-game-stats /path/to/Dynasty1
```

#### Export sections

| Flag | JSON fields | Contents |
|------|-------------|----------|
| *(default)* | all below | Everything when no section flags are set |
| `--season` | `season` | Current year, week, phase |
| `--teams` | `teams` | Schools, records, poll ranks, recruiting class ranks/score |
| `--rosters` | `rosters` | Active rosters with player ratings |
| `--games` | `games` | Schedule, scores, optional per-game stats |
| `--recruits` | `recruits` | Recruiting board + nested `player` attributes |
| `--recruiting` | `recruiting` | Pursuit state, NIL, visits, full school interest, active pitches |
| `--season-stats` | `seasonPlayerStats`, `seasonTeamStats` | Season stat totals |
| `--coaches` | `coaches` | Staff, contracts, career records |
| `--leaving-players` | `leavingPlayers` | Graduation / exit pipeline |
| `--injuries` | `injuries` | Active injuries |
| `--depth-charts` | `depthCharts` | Depth chart slots by team |
| `--history` | `playerAwards`, `leagueAwards`, `conferenceChampions`, `recordBook` | Awards and the full stat record book |
| `--school-grades` | `schoolGrades` | Per-school recruiting pitch grades |
| `--pipelines` | `pipelineInfluence` | Per-school pipeline influence by region |
| `--rivalries` | `rivalries` | Head-to-head rivalry records |
| `--position-changes` | `positionChanges` | Player position/archetype change history |
| `--draft` | `draftPicks` | Draft pick slot assignments |
| `--bowls` | `bowlGames` | Bowl game metadata |

Player and recruit exports also include enriched attributes when `--rosters` or
`--recruits` is selected: NIL, redshirt status, motivations, abilities, career
summary stats, archetype traits, and more. Game exports include quarter scores,
weather, kicking, and offensive line stats when `--games` is selected.

Additional flags:

- `--no-game-stats` — omit team/player stat lines from game exports
- `--tuning-path` — path to `dynasty-tuning-binary.FTC` for skill group labels, attributes, and current levels (auto-discovered under `--schema-dir` when omitted)
- `-o path` — write JSON to a file (stdout if omitted)
- `--pretty=false` — compact JSON

Run `cfb-dynasty -h` or `cfb-dynasty export -h` for full usage.

### Recruiting formula constants

Dump recruiting tunables from the game install's tuning FTC (no save required):

```bash
cfb-dynasty recruiting-tunables -schema-dir ./data
cfb-dynasty recruiting-tunables -schema-dir ./data -o recruiting-tunables.json
```

Place `cfb27-db-data/<patch>/dynasty-tuning-binary.FTC` under your schema directory
(or pass `-tuning` with an explicit path). The export includes scalar thresholds,
lookup arrays, pitch definitions, visit/action costs, and high-school generation
weights.

### Team IDs

Exported `teamId` / `teams[].id` / `player.teamIndex` values are the game's
stable `Team.TeamIndex` IDs (Akron = 1, Alabama = 2, … Sacramento State = 137),
not the Team table row order. Newer programs keep high TeamIndex values while
sitting earlier in the table (for example Appalachian State is row 3 with id
125), so row numbers must not be used as join keys. Air Force is id `0`. FCS
placeholder slots (`TeamIndex` 255) are omitted.

### Recruiting class ranks

Each `teams[]` entry can include `recruitingClass` after signing day:

- `nationalRank` / `conferenceRank` — stored on the `Team` row as `TopClassRank` and `TopClassConferenceRank`.
- `score` — derived from the team's `CommittedPlayers` list: each commit's `CommitScore` is weighted by national rank using `TopClassesRankWeightPercentageTable` from `dynasty-tuning-binary.FTC` (`score += commitScore * weight / 100`).
- `commitCount` — number of committed recruits included in the score.

Ranks are populated once the game evaluates top classes (typically post-signing). The score is not stored directly in the save; it is recomputed at export time from commit data plus tuning weights. Without tuning FTC, ranks still export but `score` is omitted.

### Example: recruits with `jq`

```bash
cfb-dynasty export -schema-dir ./schemas --recruits /path/to/Dynasty1 | \
  jq '.recruits[] | select(.nationalRank != null and .player != null) |
      {rank: .nationalRank, name: "\(.player.firstName) \(.player.lastName)",
       position: .player.position, archetype: .player.archetypeLabel, overall: .player.overall,
       skillCaps: .player.skillGroupCaps, currentLevels: .player.skillGroupCurrentLevels,
       skillGroups: .player.skillGroups, skillCapTotal: .player.skillGroupCapTotal, unlockedTotal: .player.skillGroupUnlockedTotal}'
```

### Skill group caps and current levels

Every player (including a recruit's linked `player`) exports its six skill-group
cap slots read straight from the save:

- `skillGroupCaps` — greyed-out/capped upgrade slots per bucket (`SkillGroupCapMax - saved value`, 0..20 each).
- `skillGroupUnlockedSlots` — unlocked upgrade slots still available in each bucket (raw `SkillGroupCap1..6` save values).
- `skillGroupCurrentLevels` — current developed levels per bucket. The exporter computes each group OVR with per-attribute weights of 15 for Primary, 3 for Secondary, and 1 for Tertiary, then uses the player's position to select the minimum OVR for levels 2 through 20:
  - `FB`, `TE`: `59, 61, 63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 84, 85, 86, 87, 88, 89`
  - `WR`, `MIKE`, `K`, `P`: `61, 64, 66, 69, 71, 74, 76, 79, 81, 84, 86, 89, 91, 93, 95, 96, 97, 98, 99`
  - All other positions: `63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 85, 87, 89, 91, 93, 95, 97, 99`
  Level 1 applies when no minimum is reached. The result is clamped to the corresponding `skillGroupUnlockedSlots` ceiling.
- `skillGroupCapTotal` — total greyed-out upgrade slots across all six buckets (sum of `skillGroupCaps`).
- `skillGroupUnlockedTotal` — total unlocked upgrade slots (sum of save values; useful recruit ceiling proxy).
- `skillGroups` — when tuning data is available, each bucket includes its UI label, attributes and tiers, capped/unlocked slot counts, and `attributeCount` (number of individual attributes in that bucket from tuning).

Each bucket has up to 20 upgrade slots in the save (`SkillGroupCapMax` from tuning). The UI shows 10 segments per bucket, but the save tracks capacity on a 0..20 scale. Bucket labels, attribute membership, and tiers come from `dynasty-tuning-binary.FTC` in the game install. Place that file under your `--schema-dir` as `cfb27-db-data/<patch>/dynasty-tuning-binary.FTC`, or pass `--tuning-path` on export.

When group attributes and ratings are known but tier metadata is incomplete, the exporter uses a simple average and the result may be inaccurate. When group membership or any required rating is unavailable, it omits the entire six-value array rather than exporting partial values.

### Per-game player stats

When game stats are included (`--games`, on by default), each game carries a
`playerGameStats` list. Each entry is one player's line for that game:

- `playerId` — the player's row index (join to `rosters` for the full record).
- `player` — a lightweight identity (`firstName`, `lastName`, `position`,
  `jersey`, `teamIndex`) so lines are readable without a join.
- `offense` / `defense` — the stat line(s); a player with both merges into one entry.
- `specialTeams` — kick/punt return line (attempts, yards, longest, TDs) when the
  player returned kicks or punts that game.

The game-stat rows themselves store no player reference — ownership lives on the
`Player` side via a `GameStats[]` array store that points at each player's rows.
The exporter inverts those arrays to attribute every stat line, and rows are
bucketed into games by their direct `SeasonGame` record index (stale references
to other seasons are dropped).

Kick/punt returns live in separate `KPReturnStats` tables (the box-score
`TeamStats` table has no special-teams TD field, so team totals cover yardage
only). The same `GameStats[]` / `SeasonStats[]` array stores link those rows back
to their players, and the `specialTeams` block appears on both per-game
(`playerGameStats`) and per-season (`seasonPlayerStats`) lines.

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
cfb-dynasty export -schema-dir ./schemas --history /path/to/Dynasty1 | \
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

	file, err := dynasty.Open("/path/to/Dynasty1", &settings)
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
| [`internal/desktop`](./internal/desktop/) | Desktop GUI app, CSV exporter port, save discovery |
| [`internal/binary`](./internal/binary/) | Low-level byte scanning helpers |
| [`internal/bitview`](./internal/bitview/) | Bitfield read helpers for record decoding |
| [`internal/compress`](./internal/compress/) | Decompression (zlib) |
| [`cmd/cfb-dynasty`](./cmd/cfb-dynasty/) | Command-line tool |
| [`cmd/cfb-dynasty-gui`](./cmd/cfb-dynasty-gui/) | Cross-platform desktop explorer (Wails) |

## Expected save location (PC)

```
%USERPROFILE%\Documents\EA SPORTS CFB27\saves\
```

Dynasty saves are typically extensionless (or `.sav`). The desktop GUI auto-lists files in this folder on Windows.

## Known limitations

- **Schema required** — export and record decoding need a matching `C27_*_*.gz` bundle (not shipped here; build one with `schema-build` from an MMC Frosty extract).
- **Tuning data optional** — skill group labels and `recruiting-tunables` need `dynasty-tuning-binary.FTC` from the game install (not shipped here).
- **Read-only** — no save writing or editing.
- **Record-book team names** — league ranks below #1 often omit a stored team name in the save (the exporter does not invent one).
- **Unplayed saves** — game and season player stats may be empty until games are simmed.
- **Large exports** — full exports with `--recruits` and `--rosters` can be tens of MB; use section flags to trim output.

## Goals

1. Reliable **read-only** parsing of local PC dynasty saves.
2. Export the full set of league surfaces listed above as stable JSON for apps, bots, and pipelines.
3. Stay headless-friendly: copy save → parse → JSON → sync elsewhere.

## Non-goals

- Console save extraction
- Writing, repairing, or live-editing save files
- Redistributing EA schema bundles or other game assets

## Related projects

| Project | Notes |
|---------|-------|
| [madden-franchise](https://github.com/bep713/madden-franchise) | Reference parser for Madden FranTk-style table DBs |
| [Madden Franchise Editor](https://github.com/bep713/madden-franchise-editor) | Desktop editor built on `madden-franchise` |

CFB shares the broad EA engine family with Madden, but dynasty saves are a distinct format with their own schema bundles.

## Development

```bash
go build ./...
go test ./dynasty -short -count=1          # unit tests only (low memory)
go test ./dynasty -parallel 1 -count=1     # full integration suite (serial)
go run ./cmd/cfb-dynasty export -schema-dir ./schemas --help
```

Place test saves and schema bundles in a local `data/` directory. Dynasty saves may contain personal league data; schema bundles are game-derived assets — **neither should be committed** (see [Schema bundles](#schema-bundles)).

### Contributing

Issues and PRs welcome. Useful contributions include schema version coverage after EA patches, filling the remaining limitations above, better docs/examples, and regression tests against real saves.

## License

[MIT](./LICENSE)
