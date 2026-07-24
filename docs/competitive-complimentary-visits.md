---
title: "CFB 27 Competitive & Complimentary Visits"
geometry: margin=0.75in
---

Related: `docs/recruiting-tunables-math.md`, `docs/coaching-abilities-recruiting.md` (Making Friends / Hospitality).

## Executive summary

- **Main finding:** each position has a **competitive visit threshold**. The UI marks any same-position parallel visit as competitive, but it never shows the `CompetitiveRecruitCountThreshold` that gates the −5 penalty.
- **Assumption:** the threshold is how many **other** same-position visitors you can bring **without** the −5. The penalty is assumed to apply only when that count is **exceeded** (`count > threshold`), not at equality.
- **In-game CB test (threshold 2):** solo, one other CB, and two other CBs all produced the same influence gain (no −5). That matches the FTC value of 2 and rules out `count >= threshold`.
- Complimentary visits are a separate +5 when parallel visitors match that position’s complimentary list (e.g. QB ↔ WR/TE/T/C).
- **Confirmed in-game:** +5 and −5 can both apply on the same visit; a solo visit with no parallel visitors gets neither.
- Making Friends stacks on the complimentary +5.

**Status:** threshold table confirmed from tuning. At-or-below-threshold visits skip the −5 (CB A/B). Over-threshold penalty (`count > threshold`) is the working assumption; not yet shown for `count = threshold + 1`.

In-game CB visit comparison courtesy of **Agusting** (MaxPlaysCFB Discord). Site asset: `/cfb-data/visit-tests/results.png`.

---

## Why this matters

The recruiting UI marks parallel visitors at the same position as competitive. It does not show the threshold associated with the −5 `competitiveVisitPenalty`. You can therefore see a competitive label while still being under (or at) the free allowance, and take no influence hit until you go over the threshold.

---

## Scalars

From `VisitTunables` (`dynasty-tuning-binary.FTC`):

| Field | Value | Effect |
|-------|------:|--------|
| `complimentaryVisitBonus` | **+5** | Influence when the visit counts as complimentary |
| `competitiveVisitPenalty` | **−5** | Influence when the visit counts as competitive |

These are flat additives (not tier-scaled). Coaching ability **Making Friends** stacks on top of the complimentary bonus.

---

## How it works

For a recruit on a visit, the game looks up that recruit’s position in `VisitTunables.PositionVisitTable` → `PositionVisitInfo`:

```
info = PositionVisitTable[recruit.position]
competitiveCount = # of other parallel visitors whose position ∈ info.CompetitivePositions

# Assumption: strict greater-than (threshold = max free competitive others)
if competitiveCount > info.CompetitiveRecruitCountThreshold:
    apply CompetitiveVisitPenalty (−5)

if any parallel visitor position ∈ info.ComplimentaryPositions:
    apply ComplimentaryVisitBonus (+5)
```

The two checks are **independent**. Both modifiers can apply to the same visit (confirmed in-game). A visit with **no** parallel visitors gets neither modifier (confirmed in-game).

---

## In-game CB test

CB `CompetitiveRecruitCountThreshold` in FTC is **2**. Agusting (MaxPlaysCFB Discord) ran matched Week 4 Washington visits for CB **George Dogbe**, varying only how many other CBs were on the trip:

| Visit setup | Other CBs | UI “Competitive Visits” | Measured influence bar |
|-------------|----------:|-------------------------|------------------------|
| Solo | 0 | None | Same (~47 px) |
| + Jameson Clinton | 1 | Lists Clinton | Same (~47 px) |
| + Clinton and Kevin Meehan | 2 | Lists both | Same (~47 px) |

Same game context each time (clear win vs Minnesota, Attend Team Meeting). Influence did not drop at one or two other CBs, even though the UI labeled those visitors competitive. That is exactly at the FTC threshold of 2, so the −5 is **not** applied at `count == threshold`.

**Working assumption:** penalty when `count > threshold` (a third other CB for CB recruits). Over-threshold case not yet captured in this test set.

---

## PositionVisitTable

`CompetitivePositions` is always the recruit’s own position group. The **threshold** is treated as the max number of **other** visitors at those positions allowed before the −5 (see assumption above). The UI does not show this value. Re-dumped from patch 2 and patch 3 FTC (identical).

| Visiting as | Competitive threshold | Competitive positions | Complimentary positions |
|-------------|----------------------:|-----------------------|-------------------------|
| QB | 1 | QB | WR, TE, T, C |
| HB | 1 | HB | C, T, G, FB |
| FB | 2 | FB | HB |
| WR | 3 | WR | QB |
| TE | 2 | TE | QB |
| T | 2 | T | C, G, QB, HB |
| G | 2 | G | T, C, QB, HB |
| C | 2 | C | G, T, QB, HB |
| DE | 2 | DE | DT, OLB, MLB |
| DT | 2 | DT | DE, OLB, MLB |
| OLB | 2 | OLB | MLB, DE, DT, CB, FS, SS |
| MLB | 2 | MLB | DE, DT, OLB, CB, FS, SS |
| CB | 2 | CB | OLB, MLB, FS, SS |
| FS | 2 | FS | OLB, MLB, CB, SS |
| SS | 2 | SS | OLB, MLB, CB, FS |
| K | 1 | K | None |
| P | 1 | P | None |
| LS | 1 | LS | None |

Source: `VisitTunables.PositionVisitTable` → `PlayerAcquisitionPositionLookupTable` → `PositionVisitInfo`.

### Reading examples

- **CB (threshold 2):** up to **two** other CBs → no −5 in A/B; a third other CB is assumed to trigger −5. UI may still label them competitive either way.
- **WR (threshold 3):** up to three other WRs should be free; four are assumed to take −5. A QB on the trip is complimentary (+5).
- **QB (threshold 1):** one other QB should still be free; two other QBs are assumed to take −5. A WR/TE/T/C on the same trip → +5. Both can apply together.
- **K / P / LS:** threshold 1 with no complimentary positions. One other same-position specialist should be free; two are assumed to go competitive for the −5.
- **Solo visit:** no parallel visitors → no ±5 from these rules.

---

## Evidence status

| Claim | Status |
|-------|--------|
| Hidden per-position competitive threshold exists in tuning | Confirmed (FTC) |
| CB threshold = 2 (and full table above) | Confirmed (FTC patch 2 & 3) |
| ±5 scalars | Confirmed (FTC) |
| Per-position competitive/complimentary lists | Confirmed (FTC) |
| UI shows same-position visits as competitive without showing the threshold | Observed (A/B screenshots) |
| Complimentary + competitive can both apply | **Confirmed (A/B)** |
| Lone visit (no parallel visitors) → neither modifier | **Confirmed (A/B)** |
| At `count == threshold`, no −5 | **Confirmed (A/B):** CB + 2 other CBs, same influence as solo |
| Penalty when `count > threshold` | **Assumption** (fits FTC + A/B; over-threshold case not yet shown) |
| Native apply bytecode | Blocked (Denuvo); not required given A/B |

**Credit:** in-game CB visit screenshots and test runs by **Agusting** (MaxPlaysCFB Discord).

---

## Related coaching abilities

Making Friends / Hospitality amplify complimentary visit value; see `docs/coaching-abilities-recruiting.md`.
