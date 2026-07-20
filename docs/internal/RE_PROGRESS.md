# CFB Dynasty Save — Reverse Engineering Progress

Living document for debugging and mapping the College Football dynasty file format. Update this as you learn; link hex dumps, sample files, and schema versions here.

**Game target:** EA Sports College Football 27 (PC)  
**Reference:** Madden franchise format (`madden-franchise`, FranTk / SPBF table DB)  
**Last updated:** _fill in as you go_

---

## Investigation log

| Date | Investigator | Summary |
|------|--------------|---------|
| | | |

_Add rows chronologically (newest first)._

---

## Sample files

| Label | Path | Game version | Schema major/minor | Notes |
|-------|------|--------------|-------------------|-------|
| | | | | |

Store copies outside the repo if they contain personal dynasty data. Note file size and SHA256 only in git when possible.

---

## Phase 0 — Container format

Check the raw save before assuming Madden parity.

| Check | Madden (expected) | CFB 27 (observed) | Done |
|-------|-------------------|-------------------|------|
| File extension / naming | `CAREER-*` (no ext) | `.sav` in `Documents\College Football 27\Saves` | ☐ |
| Starts with `FrTk` (`46 72 54 6b`) when uncompressed | Sometimes | | ☐ |
| Zlib header `78 9c` when compressed | Common | | ☐ |
| Skip header before inflate (Madden uses `0x52`) | Yes for franchise | | ☐ |
| OLE/CFB compound file (`D0 CF 11 E0 …`) | Unlikely for inner payload | | ☐ |
| Extra encryption / checksum wrapper | Unknown | | ☐ |
| File size vs dynasty progress correlation | Yes | | ☐ |
| Safe read-only copy while game closed | Yes | | ☐ |

**Notes:**

```
(paste hex of first 64–128 bytes here)
```

---

## Phase 1 — Table discovery

Madden scans unpacked bytes for markers: `SPBF`, `ASTO`, `SPEX`.

| Marker | Found? | Offset(s) | Notes |
|--------|--------|-----------|-------|
| `SPBF` | ☐ | | |
| `ASTO` | ☐ | | |
| `SPEX` | ☐ | | |
| Other | ☐ | | |

**Table count in test save:** _  
**Asset table offset (Madden: uint32 BE @ 0x04, count @ 0x24):** _

---

## Phase 2 — Schemas

Schemas define field names and types; the save file carries per-table index lists for bit offsets.

| Task | Status | Notes |
|------|--------|-------|
| Locate schema assets in game install (Frosty / CAS) | ☐ | |
| Find `FranTkData` / `e-Schemas` blocks | ☐ | |
| Extract and version schemas (major/minor) | ☐ | |
| Match schema version embedded in save | ☐ | |
| Compare table names to Madden (`SeasonGame`, `Team`, …) | ☐ | |
| Document CFB-specific tables (recruiting, dynasty, etc.) | ☐ | |

**Schema files on disk:**

| Game patch | Major | Minor | Path / source |
|------------|-------|-------|---------------|
| | | | |

---

## Phase 3 — High-value tables (export targets)

Map CFB table names to Commissioner's Office needs. Rename columns as discovered.

### Season / calendar

| Concept | Madden analogue | CFB table name | Key fields | Parsed? |
|---------|-----------------|----------------|------------|---------|
| Current week/year | `SeasonInfo` | | | ☐ |
| Game schedule | `SeasonGame` | | | ☐ |
| Standings | `Team` + standings logic | | | ☐ |

### Games (MVP for ELO / results)

| Field | Madden `SeasonGame` | CFB field | Type | Notes |
|-------|---------------------|-----------|------|-------|
| Home team | `HomeTeam` | | ref | |
| Away team | `AwayTeam` | | ref | |
| Home score | `HomeScore` | | int | |
| Away score | `AwayScore` | | int | |
| Season year | `SeasonYear` | | int | |
| Week | `SeasonWeek` | | int | |
| Week type | `SeasonWeekType` | | enum | PreSeason / RegularSeason / … |
| Game status | `GameStatus` | | enum | |

### Teams

| Field | Madden `Team` | CFB field | Parsed? |
|-------|---------------|-----------|---------|
| Short name | `ShortName` | | ☐ |
| Conference standing | `CurSeasonConfStanding` | | ☐ |
| W/L | `ConfWin`, etc. | | ☐ |

### Optional (later)

- [ ] Player rosters / ratings  
- [ ] Season stats (`SeasonOffensiveStats`, etc.)  
- [ ] Recruiting board  
- [ ] Coaching staff / carousel  

---

## Phase 4 — Parser implementation (Go)

| Milestone | Status | PR / commit |
|-----------|--------|-------------|
| Read file from disk | ☐ | |
| Decompress / unpack | ☐ | |
| List tables (name, id, offset) | ☐ | |
| Load schema JSON/XML | ☐ | |
| Build offset table for one table | ☐ | |
| `readRecords` for `SeasonGame` | ☐ | |
| Resolve team references | ☐ | |
| JSON export CLI | ☐ | |
| Tests with fixture save (redacted) | ☐ | |

---

## Phase 5 — Integration pipeline

| Step | Status | Notes |
|------|--------|-------|
| PC agent copies `.sav` to shared folder / S3 | ☐ | |
| Cron runs extractor on new/changed file | ☐ | |
| Normalized output → LeagueLines API / DB | ☐ | |
| Idempotent sync (hash/mtime) | ☐ | |
| Schema version mismatch alerting | ☐ | |

---

## Open questions

1. _Does CFB 27 use the same zlib skip offset as Madden franchise saves?_  
2. _Are dynasty saves encrypted beyond compression?_  
3. _What schema version does launch-day 27 ship with?_  
4. _Table names: `Dynasty` root vs `Franchise`?_  
5. _Online dynasty: is the commish save the only authoritative file?_  

---

## Dead ends / ruled out

_Document approaches that failed so we don't retry them._

| Hypothesis | Result | Date |
|------------|--------|------|
| | | |

---

## Tools & commands

```bash
# First bytes (replace path)
xxd -l 128 /path/to/Dynasty1.sav

# Entropy / size
wc -c /path/to/Dynasty1.sav
sha256sum /path/to/Dynasty1.sav

# Search for Madden table signatures in unpacked blob (after manual inflate trial)
# rg -a SPBF unpacked.bin
```

---

## References

- [madden-franchise `FranchiseFile.js`](https://github.com/bep713/madden-franchise/blob/master/src/FranchiseFile.js) — compress/decompress, table scan, schema pick  
- Madden franchise editor README — zlib unpack, `SeasonGame` table, schema from Frosty  
- [CFB 27 FAQ](https://www.ea.com/games/ea-sports-college-football/college-football-27/faq) — PC release, dynasty mode
