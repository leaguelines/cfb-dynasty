# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.0] - 2026-07-20

### Added

- Cross-platform desktop GUI (`cfb-dynasty-gui`) built with Wails (HTML/CSS + Go bindings): browse dynasty data and export JSON or CSV without a local web server.
- GUI first-run setup: choose schema directory (`C27_*.gz`), persist preference, open saves via file picker; auto-discover Windows saves under `Documents\EA SPORTS CFB27\saves`.
- Desktop CSV/JSON export helpers ported from the web explorer (per-collection CSV and full CSV zip).
- `cfb-dynasty gui` helper that launches `cfb-dynasty-gui` when it is on `PATH`.
- GitHub Actions release workflow: pushes to `main` update a **nightly** prerelease; `v*` tags create versioned releases with GUI + CLI packages for Linux, Windows, and macOS.
- Recruiting class exports on team records.
- `schema-build` CLI command to generate `C27_*.gz` schema bundles from Frosty FTX extracts.
- In-game team ratings (OFF / DEF / OVR) on team exports.
- This changelog.

### Fixed

- Schema generation versioning for built bundles.

## [v0.2.0] - 2026-07-10

### Added

- Load `dynasty-tuning-binary.FTC` from `--schema-dir` to label skill-group buckets.
- `recruiting-tunables` CLI and `ExportRecruitingTunables` API for formula constants from tuning data.
- Skill group export fields: capped/unlocked slots per bucket, totals, and per-bucket attribute counts from tuning.
- Full recruiting school interest lists and active pitches on pursuit exports.
- Player wear/fatigue, team poll/standings/recruiting board, coach traits, season period flags, and school grade scores.
- Team prestige export fields and physical ability names on player exports.

### Changed

- Skill group caps now report greyed UI slots per bucket (`SkillGroupCapMax` minus unlocked capacity), with unlocked slots exposed separately.
- Improved career and season stat aggregation.
- Pipeline influence rows with zero influence are omitted to reduce noise.

### Fixed

- Array-store member ordering so skill-group bucket labels align with save slots across archetypes.
- Tuning array int decoding (`s_int` bias).
- Team ID mismatches across several export sections.
- Heavy integration tests gated behind `-short` to reduce memory use.

## [v0.1.0] - 2026-07-08

### Added

- Richer player exports: NIL, redshirt, motivations, abilities, career summary, dev trait, impact player, captain.
- Recruit enrichments: full school interest lists and recruit stage advance.
- School grades and pipeline influence sections.
- Game detail: quarter scores, weather, kicking stats, o-line stats.
- Team schemes and playoff context; coach job security and position ratings.
- New top-level sections: rivalries, draft picks, bowl games, position change history (with matching CLI flags).

### Fixed

- Player weight decoding (+160 offset).

## [v0.0.2] - 2026-07-08

### Added

- `isAth` on player exports, derived from linked recruit alternate positions (`AlternatePosition1` / `AlternatePosition2`), for both roster and recruit exports.

## [v0.0.1] - 2026-07-08

Initial public release of the library and CLI.

### Added

- Dynasty save parsing: table discovery, gzip schema loading, bitfield record decoding.
- Structured JSON export API and `cfb-dynasty` CLI with selective section flags (games, teams, rosters, recruits, recruiting, season stats, coaches, leaving players, injuries, depth charts, history, and more).
- Player archetypes and labels; first-pass skill group caps.
- Full record book export (`recordBook`) across league / conference / team × career / season / game.
- Special teams return stats on player game/season exports; yardage on team stats.
- Stable team ID mapping for joins across tables.

### Fixed

- Recruit-to-player linkage (active rows only; reject empty player refs).
- Per-game player stat attribution (invert `GameStats[]`); off-by-one game bucketing.
- History export skipping stale rows beyond `NextRecordToUse`.
- Record-book team refs and default QB position inference for holders.
- Array-store team tables and conference membership by exact reference indexing.
- Archetype mapping corrections.

[Unreleased]: https://github.com/leaguelines/cfb-dynasty/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/leaguelines/cfb-dynasty/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/leaguelines/cfb-dynasty/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/leaguelines/cfb-dynasty/compare/v0.0.2...v0.1.0
[v0.0.2]: https://github.com/leaguelines/cfb-dynasty/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/leaguelines/cfb-dynasty/releases/tag/v0.0.1
