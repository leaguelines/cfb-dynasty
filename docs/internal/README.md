# CFB Dynasty Extractor

Go tool for reading EA Sports College Football dynasty save files on PC and exporting structured data for external apps (e.g. [LeagueLines Commissioner's Office](https://commissioners-office.leaguelineshq.com/)).

This project is a sibling effort to the [Madden Franchise Editor](https://github.com/bep713/madden-franchise-editor) stack, which parses Madden career saves via the [`madden-franchise`](https://github.com/bep713/madden-franchise) Node library. CFB uses the same broad engine family, but dynasty saves are a separate format and require their own parser.

## Status

**Pre-alpha / reverse engineering.** No parser yet. See [RE_PROGRESS.md](./RE_PROGRESS.md) for investigation notes and checklists.

## Goals

1. **Read-only** parsing of local PC dynasty saves (no writing back to game files in v1).
2. Export league-relevant data: season state, schedule, scores, team records, and optionally rosters/stats.
3. Support a headless pipeline: copy save → parse → JSON/API → sync to Discord bot or web admin.

## Expected save location (PC)

```
%USERPROFILE%\Documents\College Football 27\Saves\
```

CFB 27 is the first College Football release on PC. Save files are commonly named like `Dynasty1.sav` (confirm after launch).

## Non-goals (for now)

- Console save extraction
- Modifying or repairing corrupted saves

## Related work

| Resource | Notes |
|----------|--------|
| [madden-franchise](https://github.com/bep713/madden-franchise) | Reference implementation for FranTk-style table DBs |
| Madden franchise editor `schemaSearchService` | CAS/LZ4 schema extraction from game install |
| [RE_PROGRESS.md](./RE_PROGRESS.md) | RE log, byte signatures, table mapping |
| [EXPORT_ROADMAP.md](./EXPORT_ROADMAP.md) | Planned export scope and deferred tables |

## Development

```bash
# Once code exists:
go build ./...
go test ./...
```

## License

TBD — align with parent project if this ships alongside the franchise editor.
