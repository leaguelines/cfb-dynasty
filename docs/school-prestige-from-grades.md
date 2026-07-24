---
title: "CFB 27 School Prestige From My School Grades"
excerpt: "Program star rating is a weighted blend of five My School grades, not an average of every motivation. Coach Prestige and Academics have zero weight."
publishDate: "Jul 22 2026"
updatedDate: "Jul 22 2026"
tags:
  - cfb27
  - recruiting
seo:
  title: "CFB 27 School Prestige From My School Grades"
  description: "CFB 27 program stars are a weighted blend of five My School grades. Coach Prestige and Academics do not count."
  pageType: article
---

Related: [Recruiting tunables](./recruiting-tunables-math.md) · [NIL Offer → Recruiting Influence](./nil-offer-recruiting-influence.md).

## Executive summary

- **Main finding:** program prestige (`TeamPrestige`, half-star units where **10 = 5★**) comes from `MySchoolTeamPrestigeTuning`, not from averaging every My School grade.
- Only **five** grades have nonzero weight: Championship Contender, Pro Potential, Brand Exposure (**1.0** each), plus Program Tradition and Conference Prestige (**0.75** each).
- Coach Prestige, Academic Prestige, Athletic Facilities, Campus Lifestyle, Stadium Atmosphere, Playing Style, Proximity, Coach Stability, and Playing Time all weigh **0**.
- Each letter grade maps to points via `GradeToPointsConversion` (A+ = **5.25** through F = **0**). Prestige is a weighted average of those points, then mapped to half-stars.
- Best-fit star mapping from a fresh dynasty save: **`floor(score × 2)`** (clamped 0-10). About **68%** of teams match exactly; **100%** are within ±0.5★. Remaining misses line up with prestige lagging live grades.
- My School grades have per-grade update frequencies (`Weekly` / `Yearly` / `Static` / `ByPlayer`). Working assumption: **program stars refresh on a seasonal path**, not every weekly grade tick.

**Status:** weights, grade-to-points, and grade `ChangeFrequency` values confirmed from `dynasty-tuning-binary.FTC`. Pro Potential aggregation, `floor(score × 2)`, and seasonal star refresh timing are save-correlated / inferred, not yet decompiled.

Requested by **Doug** (MaxPlaysCFB Discord).

---

## Why this matters

The UI shows many My School letter grades plus a separate prestige star rating. Only five of those grades feed the stars. Raising Coach Prestige or Academics can still help recruiting pitches, but it does not change program prestige under this tuning.

---

## Core formula

From `MySchoolTuning.TeamPrestigeTuning` -> `MySchoolTeamPrestigeTuning`:

```
points(grade) = GradeToPointsConversion[LetterGrade]   # A+ through F

score = (
    points(ChampionshipContender) * 1.00 +
    points(ProgramTradition)      * 0.75 +
    points(ProPotential)          * 1.00 +
    points(BrandExposure)         * 1.00 +
    points(ConferencePrestige)    * 0.75
) / 4.5

TeamPrestige ~= floor(score * 2)   # half-star units; 10 = 5.0★
```

`4.5` is the sum of the nonzero weights (`mTeamPrestigeWeightsTotal` in the binary).

**Pro Potential** is stored per position on `MySchoolTrackingTable`. The best save fit averages **QB, RB, WR, TE, OL, DL, LB, DB** (exclude K and P), converts each letter to points, then averages.

---

## Grade weights

`MySchoolTeamPrestigeTuning.GradeWeights` (indexed by `MySchoolGrade`):

| My School grade | Weight |
|-----------------|-------:|
| Championship Contender | **1.00** |
| Pro Potential | **1.00** |
| Brand Exposure | **1.00** |
| Program Tradition | **0.75** |
| Conference Prestige | **0.75** |
| Playing Style | 0 |
| Proximity to Home | 0 |
| Campus Lifestyle | 0 |
| Stadium Atmosphere | 0 |
| Academic Prestige | 0 |
| Coach Stability | 0 |
| Coach Prestige | 0 |
| Athletic Facilities | 0 |
| Playing Time | 0 |

---

## Letter grade to points

`GradeToPointsConversion` (LetterGrade order A+ through F):

| Grade | Points |
|-------|-------:|
| A+ | **5.25** |
| A | 4.833 |
| A- | 4.417 |
| B+ | 4.00 |
| B | 3.58 |
| B- | 3.167 |
| C+ | 2.75 |
| C | 2.333 |
| C- | 1.917 |
| D+ | 1.50 |
| D | 1.083 |
| D- | 0.667 |
| F | **0** |

---

## Worked examples

Grades and saved `TeamPrestige` from a fresh CFB 27 dynasty save. Predicted stars use `floor(score × 2) / 2`.

### Ohio State: 5.0★ (match)

| Input | Grade | Points | Weight | Contribution |
|-------|------:|-------:|-------:|-------------:|
| Championship Contender | A+ | 5.250 | 1.00 | 5.250 |
| Program Tradition | A+ | 5.250 | 0.75 | 3.938 |
| Brand Exposure | A+ | 5.250 | 1.00 | 5.250 |
| Conference Prestige | A+ | 5.250 | 0.75 | 3.938 |
| Pro Potential (8-pos avg) | | 5.041 | 1.00 | 5.041 |

Pro Potential grades: QB B, RB/WR/TE/OL/DL/LB/DB A+.

```
sum = 23.416
score = 23.416 / 4.5 = 5.204
floor(5.204 × 2) = 10  ->  5.0★
```

Saved prestige: **5.0★**.

### Oregon: 4.5★ (match)

| Input | Grade | Points | Weight | Contribution |
|-------|------:|-------:|-------:|-------------:|
| Championship Contender | A+ | 5.250 | 1.00 | 5.250 |
| Program Tradition | B+ | 4.000 | 0.75 | 3.000 |
| Brand Exposure | A | 4.833 | 1.00 | 4.833 |
| Conference Prestige | A+ | 5.250 | 0.75 | 3.938 |
| Pro Potential (8-pos avg) | | 4.937 | 1.00 | 4.937 |

```
sum = 21.958
score = 4.880
floor(4.880 × 2) = 9  ->  4.5★
```

Saved prestige: **4.5★**. Tradition at B+ is the main gap vs a full 5★ profile.

### Rutgers: 2.0★ (match)

| Input | Grade | Points | Weight | Contribution |
|-------|------:|-------:|-------:|-------------:|
| Championship Contender | C- | 1.917 | 1.00 | 1.917 |
| Program Tradition | D | 1.083 | 0.75 | 0.812 |
| Brand Exposure | C- | 1.917 | 1.00 | 1.917 |
| Conference Prestige | A+ | 5.250 | 0.75 | 3.938 |
| Pro Potential (8-pos avg) | | 2.228 | 1.00 | 2.228 |

```
sum = 10.812
score = 2.403
floor(2.403 × 2) = 4  ->  2.0★
```

Saved prestige: **2.0★**. Conference prestige is strong; contender, brand, and tradition keep the stars down.

### Oregon State: 2.0★ (match)

| Input | Grade | Points | Weight | Contribution |
|-------|------:|-------:|-------:|-------------:|
| Championship Contender | D | 1.083 | 1.00 | 1.083 |
| Program Tradition | C | 2.333 | 0.75 | 1.750 |
| Brand Exposure | D+ | 1.500 | 1.00 | 1.500 |
| Conference Prestige | B | 3.580 | 0.75 | 2.685 |
| Pro Potential (8-pos avg) | | 2.281 | 1.00 | 2.281 |

```
sum = 9.299
score = 2.066
floor(2.066 × 2) = 4  ->  2.0★
```

Saved prestige: **2.0★**.

### Alabama: predicted 4.5★, saved 5.0★

| Input | Grade | Points | Weight | Contribution |
|-------|------:|-------:|-------:|-------------:|
| Championship Contender | A- | 4.417 | 1.00 | 4.417 |
| Program Tradition | A+ | 5.250 | 0.75 | 3.938 |
| Brand Exposure | A+ | 5.250 | 1.00 | 5.250 |
| Conference Prestige | A+ | 5.250 | 0.75 | 3.938 |
| Pro Potential (8-pos avg) | | 4.625 | 1.00 | 4.625 |

```
sum = 22.167
score = 4.926
floor(4.926 × 2) = 9  ->  4.5★
```

Saved prestige: **5.0★**. Contender at A- instead of A+ puts the live score just under the 5.0★ cut. See [prestige lag](#prestige-lag-vs-live-grades) for why the HUD can still show five stars.

---

## How often grades update

From `MySchoolGradeInfo.ChangeFrequency`:

| Grade | Frequency |
|-------|-----------|
| Playing Style | Weekly |
| Championship Contender | Weekly |
| Program Tradition | Weekly |
| Stadium Atmosphere | Weekly |
| Brand Exposure | Weekly |
| Coach Stability | Weekly |
| Athletic Facilities | Weekly |
| Playing Time | Weekly |
| Conference Prestige | Yearly |
| Coach Prestige | Yearly |
| Pro Potential | Yearly |
| Campus Lifestyle | Static |
| Academic Prestige | Static |
| Proximity to Home | ByPlayer |

Of the five prestige inputs, Contender / Tradition / Brand are **weekly**; Conference Prestige and Pro Potential are **yearly**.

`MySchoolGradeEval` exposes both `UpdateSeasonalGrades` and `UpdateTeamPrestige`, which is consistent with stars being rewritten on a seasonal pass rather than every weekly grade update.

---

## Prestige lag vs live grades

Examples above come from `DYNASTY-Fresh-for-hard-sell-testing` (**2026 PreSeason**, week 0).

**Working assumption:** `TeamPrestige` is calculated when preseason starts (or on another seasonal My School pass), then held. Live My School grades keep updating on their own schedules, so the star rating on the team can disagree with a live recompute from current grades until the next seasonal refresh.

That fits Alabama here:

- Live formula with current grades -> **4.5★**
- Saved `TeamPrestige` -> **5.0★**
- Championship Contender (a **weekly** prestige input) is already **A-**, enough to miss the 5.0★ cut
- On a brand-new preseason save, stars may still be the dynasty-start seed (Alabama as a default 5★ program) or the last seasonal write, while Contender has already moved

We have not yet proven the exact calendar hook or whether blue bloods are seeded to 5★ at dynasty create. The grade frequencies plus this preseason mismatch are the evidence so far.

---

## Grades that do not move prestige stars

Improving these alone does not change `TeamPrestige`:

- Coach Prestige / Coach Stability
- Academic Prestige
- Athletic Facilities / Campus Lifestyle / Stadium Atmosphere
- Playing Style / Proximity to Home / Playing Time

They still matter for pitch bonuses, dealbreakers, and other systems. They just are not inputs to the program star rating.

---

## Open questions

- Exact Pro Potential rollup (average-of-positions shape fits; K/P exclusion is best-fit).
- Exact score-to-half-star quantization (`floor(score × 2)` vs a threshold table).
- Confirm when `TeamPrestige` is written (preseason start vs other seasonal event), and whether dynasty create hard-seeds some schools.

---

## Sources

- Schema: `MySchoolTeamPrestigeTuning` (`GradeWeights`, `GradeToPointsConversion`), referenced from `MySchoolTuning.TeamPrestigeTuning`
- Grade schedules: `MySchoolGradeInfo.ChangeFrequency`; seasonal hooks on `MySchoolGradeEval` (`UpdateSeasonalGrades`, `UpdateTeamPrestige`)
- Enum order: `MySchoolGrade`, `LetterGrade` (A+ through F)
- Values: `dynasty-tuning-binary.FTC` (CFB 27)
- Cross-check: `data/DYNASTY-Fresh-for-hard-sell-testing` school grades vs `Team.TeamPrestige`
