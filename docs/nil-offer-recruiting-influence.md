---
title: "CFB 27 NIL Offer → Recruiting Influence"
geometry: margin=0.75in
---

Related: `docs/recruiting-tunables-math.md`, `docs/coaching-abilities.md` (Deal Maker).

## Executive summary

- **Main finding:** NIL offer influence is driven by **offer ÷ expectation**, looked up on `ScholarshipBonusPercentageSpline`.
- Meeting expectation (`100%`) grants **0** bonus influence. Underpaying is harshly penalized (e.g. **−900** at 50%). Overpaying is capped by `NILOfferRelativeMultMax = 2` at **+50** influence for a 200% offer.
- Recruits with `NILExpectation ≤ 10` use a separate absolute curve (`ScholarshipBonusRawSpline`) instead of the ratio.
- Flat extras: `ScholarshipOfferWeeklyBonus = +5`, and coaching talent **Deal Maker** adds **+10** (`ProgramPoints_RecruitingNILIncreaseInfluence`).
- Low / Medium / High “offer interest” labels are a separate 0–100 UI mapping, not the influence spline itself.

**Status:** spline values and scalars confirmed from `dynasty-tuning-binary.FTC`. Exact weekly stacking order vs other influence sources is inferred from field names and schema, not decompiled eval code.

---

## Why this matters

The recruiting UI shows NIL expectation and offer amounts, but not how those dollars convert to influence. The tuning data shows a steep asymmetric curve: falling short of expectation costs far more influence than overpaying can gain.

---

## Core formula

From `RecruitingTunables` + `ProgramPointsTuning`:

```
ratioPct = CurrentNILOffer / NILExpectation * 100
ratioPct = min(ratioPct, NILOfferRelativeMultMax * 100)   # max 200%

if NILExpectation <= ScholarshipBonusLowExpectationCutoff:  # 10
    influence = ScholarshipBonusRawSpline(offerAmount)
else:
    influence = ScholarshipBonusPercentageSpline(ratioPct)

influence += ScholarshipOfferWeeklyBonus                  # +5
influence += DealMaker_RecruitingNILIncreaseInfluence     # +10 if talent
```

Save fields on the recruit/target side: `CurrentNILOffer`, `NILExpectation`, `OriginalNILExpectation`.

Minimum offer floor in tuning: `MinimumNILToOfferPercentage = 0.80` (80% of expectation).

---

## Percentage spline (normal path)

`RecruitingTunables.ScholarshipBonusPercentageSpline`

X = offer / expectation × 100. Y = recruiting influence.

| Offer / Expected | Influence |
|------------------|----------:|
| 50% | **−900** |
| 60% | **−500** |
| 70% | **−200** |
| 80% | **−75** |
| 90% | **−15** |
| 100% | **0** |
| 120% | **+10** |
| 140% | **+20** |
| 160% | **+30** |
| 180% | **+40** |
| 200% | **+50** |

Intermediate ratios interpolate along the spline. Values above 200% are capped by `NILOfferRelativeMultMax = 2`.

---

## Raw spline (low expectation)

When `NILExpectation ≤ ScholarshipBonusLowExpectationCutoff` (**10**), the game uses `ScholarshipBonusRawSpline` (absolute offer units, roughly ±2 influence per unit):

| Offer (X) | Influence |
|----------:|----------:|
| 0 | −20 |
| 2 | −16 |
| 4 | −12 |
| 6 | −8 |
| 8 | −4 |
| 10 | 0 |
| 12 | +4 |
| 14 | +8 |
| 16 | +12 |
| 18 | +16 |
| 20 | +20 |

This avoids unstable percentages when expectation is a tiny denominator.

---

## Flat bonuses

| Source | Value | Notes |
|--------|------:|-------|
| `ScholarshipOfferWeeklyBonus` | **+5** | Flat; name implies weekly with an active offer |
| Deal Maker (`ProgramPoints_RecruitingNILIncreaseInfluence`) | **+10** | Flat coaching talent; same family as other `*Influence*` CTEs |

---

## UI: offer interest levels (not influence)

`ProgramPointsTuning.OfferInterestThresholdInfoList` maps a **0–100** score to labels:

| Label | Min | Max |
|-------|----:|----:|
| Low | 0 | 33 |
| Medium | 34 | 66 |
| High | 67 | 100 |

`OfferFeedbackThresholdInfoList` drives reaction text (Very Interested / Interested / Neutral / Disinterested / Very Disinterested) from threshold values **50 / 49 / 24 / −1 / −25**. That display scale is separate from the influence spline Y values above.

---

## Related expectation tuning (context)

These feed how large `NILExpectation` is, not the offer→influence conversion itself:

| Field | Role |
|-------|------|
| `BaseNILValueSpline` | OVR → base NIL value (ProgramPointsTuning) |
| `NILPositionValueTable` | Per-position NIL weight (QB 50, DE/T 25, WR 15, …; K/P/FB negative) |
| `NILDealbreakerModifierSpline` | Dealbreaker-side NIL modifier (offset input via `NILDealbreakerModInputOffsetForSpline = 12`) |

---

## Sources

- `RecruitingTunables`: `ScholarshipBonusPercentageSpline`, `ScholarshipBonusRawSpline`, `ScholarshipBonusLowExpectationCutoff`, `ScholarshipOfferWeeklyBonus`, `MinimumNILToOfferPercentage`
- `ProgramPointsTuning`: `NILOfferRelativeMultMax`, `OfferInterestThresholdInfoList`, `OfferFeedbackThresholdInfoList`, `BaseNILValueSpline`, `NILPositionValueTable`, `NILDealbreakerModifierSpline`
- Coaching CTE: `ProgramPoints_RecruitingNILIncreaseInfluence` (Deal Maker)

Spline Y values use FranTk excess-2³¹ signed encoding in the FTC array store; decoded values above.
