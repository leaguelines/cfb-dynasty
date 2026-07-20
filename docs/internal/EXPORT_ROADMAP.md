# Dynasty Export Roadmap

Living scope doc for planned export expansion. The detailed phased plan lives in Cursor at `Full Dynasty Export-e96e37c9.plan.md`.

## Deferred / out of scope (CFB)

These tables exist in the save schema but are **not** planned for export:

### Player contracts (`PlayerBaseContract`, `ContractYearSummary`)

- `PlayerBaseContract` fields: `Length`, `Bonus`, `Salary`.
- Reachable via `SalaryInfo.DraftedPlayerContractTable` — Madden/pro-style salary progression, not CFB roster/NIL.
- `PlayerPersonnel` also carries contract fields for personnel/staff rows; that is separate from on-field player export.

**Staff contracts** for coaches are exported today on `CoachExport` (`contractSalary`, `contractLength`, `contractYearsRemaining`, `contractStatus`) from the `Coach` table via [`export_coaches.go`](../../dynasty/export_coaches.go).

### `HistoryEntry` event log

- ~38k active rows in the test save.
- Fields include `CurrentStage`, `CurrentWeek`, `CurrentYear`, `ProgressionValue`, `MiscValue`, `ExperienceValue`, `IsSchemeFit`, `Person`, and `Source` (opaque record ref).
- Raw export would be large and low-signal without mapping `Source` refs to human-readable event types.
- **Future option:** curated/mapped history export if `Source` resolution is understood.

## In scope (phased)

See the Full Dynasty Export plan for phases 1–6: player/recruit/team/game/coach enrichments, school grades, pipeline influence, rivalries, position changes, draft picks, bowl games, CLI flags, and tests.
