---
title: "CFB 27 Player XP (Dynasty)"
geometry: margin=0.75in
---

Patch 3 tuning (`dynasty-tuning-binary.FTC`) + save `DYNASTY-2026OFFLINEFINAL`.

## Executive summary

- There is no player `Level` field. Flow is: **earn XP → hit XP threshold → gain Skill Points → spend SP on ratings**.
- Offseason training XP uses a random position/year range, the development trait multiplier, coach bonuses, and the facilities bonus.
- Difficulty / Wear & Tear XP modifiers in tuning do **not** meaningfully change dynasty offseason XP (A/B tested).
- Dev traits scale XP via a spline (%); higher traits progress faster.
- Coach percent and flat XP talents: [Coaching abilities: Player XP](./coaching-abilities-player-xp.md).

---

## 1. How is XP calculated?

Base XP comes from offseason progression, games, coach bonuses, and other progression events. Most grant logic is native code (`IncreaseExperiencePoints`) rather than readable FranTk.

What *is* readable and confirmed for dynasty:

### Flat / percent add-ons

| Source | Effect |
|--------|--------|
| Athletic facilities | +0 to +20 XP by letter grade (A+ = 20) |
| Coach XP talents | percent (e.g. Everybody Eats up to +10%) or flat; see [Player XP talents](coaching-abilities-player-xp.md) |

### Difficulty / Wear & Tear XP modifiers (not used for dynasty offseason)

`PlayerProgressionTuning` still contains:

| Field | Stored value |
|-------|-------------:|
| SkillLevelXPModifierFreshman | 0.50 |
| SkillLevelXPModifierVarsity | 0.75 |
| SkillLevelXPModifierAllAmerican | 1.00 |
| SkillLevelXPModifierHeisman | 1.25 |
| WearAndTearXPModifierOn | 1.00 |
| WearAndTearXPModifierOff | 0.40 |

Native code null-checks `skillLevelModRef` / `wearAndTear*ModRef` near `ProgressPlayers`, so the fields are loaded. But **in-game A/B testing across multiple dynasty players/positions showed no meaningful offseason XP difference** when difficulty or Wear & Tear was changed. Treat these as unused for dynasty training-results XP (likely RTG / leftover shared tuning), not live dynasty multipliers.

### Offseason training results (decoded)

Training-results XP is awarded by native `PlayerProgressionEval::AwardOffSeasonXPToPlayer` (also tied to `IssuePostSeasonTrainingRequest`). Each player gets one roll:

```
base = random(MinXP[year][pos], MaxXP[year][pos])
     // redshirt players use RedshirtMinXP / RedshirtMaxXP instead

xp   = base
     × DevTraitSpline / 100
     × (100 + coachPercent) / 100
     + facilitiesFlat
```

School years: Freshman=0, Sophomore=1, Junior=2, Senior=3.  
Tables live in tuning `PositionValueTable` (rows 11–20 normal max/min, 22–31 redshirt max/min), referenced from save `PlayerProgressionEval.MinXPTables` / `MaxXPTables` / redshirt variants.

#### Base XP ranges (non-redshirt, before multipliers)

| Pos | FR | SO | JR | SR |
|-----|---:|---:|---:|---:|
| QB | 35,150–52,250 | 31,350–48,450 | 27,550–44,650 | 25,650–40,850 |
| HB | 37,050–47,500 | 33,250–43,700 | 29,450–39,900 | 25,650–36,100 |
| FB | 15,600–31,200 | 14,625–29,250 | 12,675–25,350 | 11,700–23,400 |
| WR | 30,400–51,300 | 29,450–46,550 | 28,500–44,650 | 27,550–42,750 |
| TE | 25,350–40,950 | 24,375–39,000 | 22,425–35,100 | 21,450–33,150 |
| LT/RT | 33,150–52,650 | 32,175–50,700 | 31,200–48,750 | 30,225–46,800 |
| LG/C/RG | 32,175–50,700 | 31,200–48,750 | 30,225–46,800 | 29,250–44,850 |
| LE/RE | 30,225–44,850 | 29,250–42,900 | 27,300–39,000 | 26,325–37,050 |
| DT | 29,250–44,850 | 28,275–42,900 | 27,300–40,950 | 25,350–37,050 |
| LOLB/ROLB | 26,325–41,925 | 24,375–39,975 | 22,425–36,075 | 20,475–32,175 |
| MLB | 21,450–35,100 | 19,500–33,150 | 17,550–29,250 | 16,575–25,350 |
| CB | 31,350–49,400 | 30,400–47,500 | 28,500–45,600 | 28,500–43,700 |
| FS | 29,450–49,400 | 28,500–47,500 | 26,600–43,700 | 25,650–41,800 |
| SS | 29,450–45,600 | 28,500–43,700 | 26,600–39,900 | 25,650–38,000 |
| K/P | 4,750–9,500 | 4,500–9,000 | 4,250–8,500 | 4,000–8,000 |
| LS | 0 | 0 | 0 | 0 |

Redshirt Min/Max bands are about **5–8% lower** than the matching non-redshirt band.

Relevant coach % talents on this grant (stack additively; see [Player XP talents](coaching-abilities-player-xp.md)):

| Talent | Typical max stack |
|--------|------------------:|
| Put in Work (offseason) | +10% |
| Everybody Eats (all XP) | +10% |
| Star Maker (starters) | +17% |
| Whisperer (FR/SO) | +17% |

#### Worked examples (no coach talents)

| Player | Approx mid base | After trait |
|--------|----------------:|------------:|
| FR QB, Normal | ~43,700 | **~43,700** |
| FR QB, Elite | ~43,700 | **~65,550** |
| FR QB, Elite, max roll | 52,250 | **~78,375** |
| FR K, Normal | ~7,125 | **~7,125** |

---

## 2. How do Dev Traits affect XP gain?

The live CFB offseason progression path uses `PlayerDevTraitsSpline`, divided by `DevTraitDivisor = 100`:

| Dev trait | XP multiplier |
|-----------|--------------:|
| Normal (0) | **1.00x** |
| Impact (1) | **1.20x** |
| Star (2) | **1.35x** |
| Elite (3) | **1.50x** |

Ghidra confirms `AwardOffSeasonXPToPlayer` reads a valid `playerDevTraitInt`, and startup validation requires `mDevTraitSpline` to be attached. The matching trait-indexed spline in the CFB data is:

```
X = [0,   1,   2,   3]
Y = [100, 120, 135, 150]
```

This confirms the modifiers for **offseason progression XP**. It does not prove that every other XP source applies the same trait multiplier.

### Do not use the Madden weekly-goal values

`PlayerWeeklyGoal.DevelopmentSpline` contains 0.75x / 1.00x / 2.00x / 5.00x plus an AgeSpline. CFB dynasty's `PlayerGoalTable` is empty for every position, so that weekly-goal path is not the live dynasty progression formula.

---

## 3. How much XP to earn a level?

Use `PlayerSPThresholdSpline` (OVR → XP needed). This is the only populated XP-cost curve for dynasty players. (`PlayerXPLevelSpline` is all zeros / unused.)

| Player OVR | XP to level |
|-----------:|------------:|
| 50–60 | **1500** |
| 65–70 | **1750** |
| 75–80 | **2000** |
| 85–90 | **2250** |
| 95–99 | **2500** |

Higher-OVR players need more XP per step.

---

## 4. How many Skill Points per level-up?

**1 Skill Point per level-up** in dynasty.

Evidence: `PlayerProgressionTuning.IncrementAmount = 1`.

(RTG is different: `RTGPlayerSPPerLevel = 40`. It does not apply to dynasty.)

You spend those SP on attribute/ability upgrades (`AbilityProgressionTunable.UpgradeCostSpline`). Cost scales with current rating (often 5 SP at low ratings, rising to 30–45+ near 90–95). Cap note: `AbilityUpgradeCapPerSkillPoint = 5`.

---

## Quick examples

- **FR QB, Impact, mid roll (~43,700):** ×1.20 → **~52,440 XP**. At 60 OVR (1500 XP/level) that is ~35 SP before coach/facilities add-ons.
- **Difficulty / Wear & Tear:** changing these does **not** move dynasty offseason training XP in A/B tests.
- **90 OVR:** still needs **2250 XP** per level → **1 SP** each time the threshold is hit.

---

## Sources

| Question | Primary field / table |
|----------|----------------------|
| Offseason base XP | `PlayerProgressionEval` Min/Max/Redshirt XP tables → `PositionValueTable` |
| Coach XP add-ons | coach `Talent` / `TalentDataInt` |
| Difficulty / W&T fields | present on `PlayerProgressionTuning`, **not live for dynasty offseason** (A/B tested) |
| XP to level | `PlayerSPThresholdSpline` |
| SP per level | `IncrementAmount` (= 1) |
| DevTrait offseason XP | Native `AwardOffSeasonXPToPlayer` + `PlayerDevTraitsSpline` |
