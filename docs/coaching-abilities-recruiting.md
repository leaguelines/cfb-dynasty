---
title: "CFB 27 Coaching Abilities: Recruiting"
geometry: margin=0.75in
header-includes:
  - \usepackage{longtable}
  - \usepackage{booktabs}
  - \usepackage{array}
  - \usepackage{etoolbox}
  - \AtBeginEnvironment{longtable}{\small}
  - \sloppy
---

Derived from `dynasty-tuning-binary.FTC` (CFB 27).  
Talent values are identical in patch 2 and patch 3.

Source: `Talent` / `TalentDataInt`. Runtime: `Team.CoachTalentEffects`.  
Baseline recruiting constants from `RecruitingActionInfo`, `RecruitingTunables`, and `VisitTunables` (see also [Recruiting tunables](recruiting-tunables-math.md)).

**Full ability list (all trees):** [Coaching Abilities: Complete Guide](coaching-abilities.md).

## Executive summary

- Recruiting talents are mostly **flat influence**, **chance points**, or **letter-grade steps**. The style table below identifies each type.
- HC + matching coordinator (and often the off-side) **add**; max stack is usually HC + OC + DC.
- High-impact examples: **Making Friends** multiplies the base complimentary visit +5; **Hospitality** bumps visit activity impact; **Persuasive Personality** adds sway chance on top of motivation/pipeline sway math.
- Baseline hours/influence for SoftSell, HardSell, Sway, and visit bonuses are listed under “How to read the numbers.”
- Visit competitive/complimentary rules: [Competitive & Complimentary Visits](./competitive-complimentary-visits.md).

---

## How stacking works

| Field | Role |
|-------|------|
| `Data` | Head coach |
| `OCData` | Offensive coordinator |
| `DCData` | Defensive coordinator |

Staff contributions **add**. For position-split abilities, the coordinator that owns that side of the ball gets the larger tier (DC for DB/LB/DL; OC for QB/RB/WR-TE/OL/K-P).

**Max stack** = HC + OC + DC when all three trees own the ability.

### How to read the numbers

| Style | Meaning | Examples |
|-------|---------|----------|
| Flat influence | Added directly to weekly influence / starting interest | Magnetic, Most Influential, Making Friends, Hospitality, Upsell, pipeline abilities, Lasting Impression |
| Chance points | Added to a percent chance | Persuasive Personality, Mind Reader, Dream School |
| Grade steps | Bumps effective letter grade on the F..A+ ladder | Ideal Situation |

Letter grades used for pitch/visit bonuses (low to high):

```
F  D-  D  D+  C-  C  C+  B-  B  B+  A-  A  A+
```

Useful baselines for the examples below:

| Action | Hours | Base influence |
|--------|------:|---------------:|
| SoftSell pitch | 20 | 20 |
| HardSell pitch | 40 | 40 |
| Sway pitch | 30 | 15 |
| Contact Friends and Family | 25 | 20 |
| Complimentary visit bonus | None | 5 |
| Visit activity (Interested, A) | None | +30 |
| Visit activity (Dealbreaker, A+) | None | +48 |

SoftSell pitch grade bonus examples (matching motivation slot): A ~ +4 to +8, A+ ~ +6 to +10 depending on slot. HardSell A+ ~ +12 to +20.

---

## Recruiter

| Ability | HC | Own coord | Off coord | Max |
|---------|---:|----------:|----------:|----:|
| Persuasive Personality | 30 | 15 | 0 | **45** |
| Magnetic Personality | 10 | 4 | 2 | **16** |
| Portal King | 10 | 4 | 2 | **16** |

### Persuasive Personality (`Recruiting_SwayBoost`)

Increased chance to sway for that position group.

- HC alone: **+30** chance points.
- Matching coordinator: **+15** more (off-side contributes 0).
- Max with HC + matching coord: **+45**.

Baseline: each matching motivation on a sway already adds `ChanceToSwayBoostPerMatchingMotivation` = **15**, and pipeline level can add up to **25** more (`ChanceToSwayBoostFromPipelineInfluenceLevel`).

**Example:** You SoftSell-sway a QB with 2 matching motivations (+30 from motivations) and pipeline L3 (+15). Without the ability, that is already a sizable sway package. Persuasive Personality HC for QBs adds **+30** on top of that. With an OC who also has it, add another **+15** (total ability **+45**).

### Magnetic Personality (`Recruiting_InfluenceBoost_Start`)

Flat starting-interest boost for HS recruits at that position.

- HC: **+10**. Full staff: **+16**.

Baseline primary-pipeline starting interest ladder (`InitialInterest_Pipeline`): **0 / 5 / 10 / 15 / 25 / 35** by pipeline level.

**Example:** A non-pipeline recruit that would open at 0 starting interest opens at **+10** with HC Magnetic alone, or **+16** with full staff. For comparison, L2 is 10 and L3 is 15 on the starting-interest ladder.

### Portal King (`PlayersLeaving_InterestBoost_Transfers`)

Same numbers as Magnetic (**10 / 4 / 2**, max **16**), applied to transfer starting interest instead of HS.

**Example:** As with Magnetic, a transfer QB opens **+10** more interested with HC Portal King and **+16** with full staff.

---

## Elite Recruiter

| Ability | HC | DC | OC | Max |
|---------|---:|---:|---:|----:|
| Ideal Situation | 3 | 3 | 3 | **9** |
| Upsell | 10 | 5 | 5 | **20** |
| Most Influential | 5 | 3/2 | 2/3 | **10** |

### Ideal Situation (`Recruiting_SchoolGradeBoost_IdealPitch`)

Boosts the effective My School letter grade used for ideal-pitch evaluation by **3 grade steps** per staff tier (max **+9**).

**Example:** Your Playing Time grade is a **B**.  
- HC Ideal Situation (+3 steps): treated as **A-** for ideal-pitch purposes.  
- Full staff (+9 steps): a **B** is treated as **A+** (and anything already A- or better saturates at the top of the ladder).

Ideal pitches use that motivation's letter grade for the grade bonus. Moving from B to A- changes SoftSell slot bonuses from approximately +1/+1/+2 to +3/+4/+6 (exact slot tables in [Recruiting tunables](recruiting-tunables-math.md)).

### Upsell (`Recruiting_PointsBoost_SchoolGrade`)

Adds flat influence on top of the My School grade contribution when grades affect a recruiting action (pitches, and anything else that uses school-grade point bonuses).

- HC: **+10**. Full staff: **+20**.

**Example: SoftSell pitch with an A+ matching motivation:**  
Base SoftSell influence is **20**. A+ matching-motivation grade bonus is about **+6 to +10** depending on slot (call it **+10** for the strongest SoftSell slot).

| Setup | Grade contribution | SoftSell total (base 20 + grade) |
|-------|-------------------:|--------------------------------:|
| No Upsell | +10 | **30** |
| HC Upsell (+10) | +20 | **40** |
| Full staff Upsell (+20) | +30 | **50** |

Upsell does not change the letter grade shown on the HUD. It increases the points earned from that grade. An A+ SoftSell worth about 30 influence becomes about 50 with a full Upsell stack, before pipeline or other bonuses.

**Example: HardSell with A+ on the biggest HardSell slot (+20 grade bonus):**  
Base HardSell **40** + grade **20** = **60**. Full Upsell (+20) makes the grade piece **40**, for **80** total before pipeline.

### Most Influential (`Recruiting_Action_InfluenceBoost`)

Flat influence added to recruiting actions for that position.

- HC: **+5**. Coordinators: +3 own-side / +2 off-side. Max **+10**.

**Example:** SoftSell base **20**.  
- HC Most Influential: SoftSell becomes **25**.  
- Full staff: SoftSell becomes **30**.  
Same +5/+10 also applies on top of Friends and Family (**20** base), HardSell (**40**), Send the House (**50**), etc.

---

## Strategist

| Ability | HC | DC | OC | Max |
|---------|---:|---:|---:|----:|
| Making Friends | 20 | 15 | 15 | **50** |
| Hospitality | 15 | 5 | 5 | **25** |
| Mind Reader | 30 | 30/15 | 15/30 | **75** |

### Making Friends (`Recruiting_PointsBoost_CompVisit`)

Flat boost to complimentary-visit influence.

- HC: **+20**. Full staff: **+50**.

Baseline `complimentaryVisitBonus` = **5**, and competitive visits take a **-5** penalty. Per-position thresholds and complimentary groups: [Competitive & Complimentary Visits](./competitive-complimentary-visits.md).

**Example:** A complimentary visit worth **5** influence becomes **25** with HC Making Friends or **55** with full staff. The talent increases the base complimentary bonus by 5x to 11x.

### Hospitality (`Recruiting_PointsBoost_Visit`)

Flat boost to normal visit impact for that position.

- HC: **+15**. Full staff: **+25**.

Visit activities already scale hard with grade and interest level (Interested A = **+30**, Dealbreaker A+ = **+48**).

**Example:** An Interested visit activity graded A (**+30**) becomes **+45** with HC Hospitality or **+55** with full staff, before stacking multiple activities on the same visit.

### Mind Reader (`Recruiting_UnlockTrait_Visit`)

Chance points to learn a recruit's development trait while scouting that position.

- HC: **+30**. Matching coord: **+30**. Off-side: **+15**. Max **+75**.

**Example:** HC Mind Reader alone is a **+30** chance bump on the scouting unlock roll. With HC + matching coordinator you are at **+60**; full staff **+75**. This is an added chance, not a guaranteed reveal.

---

## Talent Developer

| Ability | HC | DC | OC | Max |
|---------|---:|---:|---:|----:|
| Home Sweet Home | 100 | 50 | 50 | **200** |

### Home Sweet Home (`Recruiting_InfluenceBoost_Pipeline`)

Large boost to the pipeline bonus on recruiting actions for that position.

- HC: **+100**. Full staff: **+200**.

Baseline pipeline bonuses on SoftSell are **0 / 2 / 4 / 6 / 8 / 10** by level; HardSell uses **0 / 4 / 8 / 12 / 16 / 20**. Home Sweet Home adds a much larger flat value to pipeline influence for that position group.

**Example:** SoftSell on an L5 pipeline recruit normally gets **+10** pipeline influence. HC Home Sweet Home adds **+100**, producing **+110** from pipeline alone. Full staff raises that to **+210** on top of the SoftSell base of 20.

---

## Program Builder (HC only)

| Ability | HC |
|---------|---:|
| Household Name | 10 |
| Hometown Discount | 25 |

### Household Name (`Recruiting_InfluenceBoost_PipelineStart`)

**+10** starting interest for recruits from your **primary** pipeline.

**Example:** The primary-pipeline starting-interest ladder is **0 / 5 / 10 / 15 / 25 / 35**. Household Name adds a flat **+10** after the pipeline value.

### Hometown Discount (`Recruiting_Action_InfluenceBoost_Pipeline`)

**+25** flat influence on recruiting actions against primary-pipeline recruits.

**Example:** SoftSell base **20** becomes **45** on a primary-pipeline target. Stacked with HC Most Influential (+5) that SoftSell is **50**. This is 5x the HC Most Influential bonus, but only for primary-pipeline recruits.

---

## CEO (HC only)

| Ability | HC |
|---------|---:|
| Dream School | 10 |
| Lasting Impression | 5 |

### Dream School (`Recruiting_BoostCommitChanceOn1st`)

**+10** chance points to instant-commit when you offer a scholarship as the recruit's top school.

Baseline `InstantCommitOddsPerStarLevel` = **10 / 8 / 6 / 4 / 2** for 1-star through 5-star.

**Example:** Offering as #1 to a 5-star without Dream School is a **2%** instant-commit roll. With Dream School that becomes **12%** if the talent adds on top of the star-level table (additive boost reading). A 3-star goes from **6%** to **16%**.

### Lasting Impression (`Recruiting_PointsBoost_Accelerate`)

**+5** additional interest for every **10 hours** spent on a recruit (`RecruitingPointsBoostAccelerate_Divisor` = 10).

**Example:** Cap is `MaxTotalHoursOnRecruitPerWeek` = **50**.  
- 10 hours → **+5**  
- 30 hours → **+15**  
- 50 hours → **+25**  

That is passive interest on top of whatever the hours actions themselves granted.

---

## Also recruiting-adjacent

See the [complete coaching guide](coaching-abilities.md) for full stacks and Program Builder / CEO value units.

| Ability | HC | Own / off coord | What it does |
|---------|---:|----------------:|--------------|
| Always Be Crootin' | 15 | 10 / 5 | weekly recruiting hours |
| Advanced Look | 1 | 1 / 0 | faster scouting |
| Reach Your Potential | 1 | 1 / 0 | interest when roster player hits skill cap |
| More the Merrier | 4 | None | +4 visit slots (base 4 → **8**) |
| Last Dance | 100 | None | +100 interest per prior visit elsewhere (UI); effect name is chance-style |
| Portal Preview | bool | None | see other schools' At Risk |

---

## Quick reference (HC values)

| Tree | Ability | HC value | What it does |
|------|---------|---------:|--------------|
| Recruiter | Persuasive Personality | +30 | sway chance |
| Recruiter | Magnetic Personality | +10 | HS starting interest |
| Recruiter | Portal King | +10 | transfer starting interest |
| Elite Recruiter | Ideal Situation | +3 grades | ideal-pitch letter grade |
| Elite Recruiter | Upsell | +10 | points from My School grades |
| Elite Recruiter | Most Influential | +5 | flat action influence |
| Strategist | Making Friends | +20 | complimentary visits |
| Strategist | Hospitality | +15 | visit impact |
| Strategist | Mind Reader | +30 | scouting trait reveal chance |
| Talent Dev | Home Sweet Home | +100 | pipeline action bonus |
| Program Builder | Household Name | +10 | primary-pipeline start interest |
| Program Builder | Hometown Discount | +25 | primary-pipeline action influence |
| CEO | Dream School | +10 | instant commit when #1 |
| CEO | Lasting Impression | +5 / 10 hrs | interest from hours spent |
| None | Always Be Crootin' | +15 hours | weekly recruiting hours |

---

## Caveats

1. Talent definitions are unchanged patch 2 to patch 3.
2. Stacking is additive across HC/OC/DC.
3. Flat vs chance vs grade-step readings above are inferred from field names, schema caps, and comparison to baseline recruiting tunables. Native apply bytecode is not in FranTk source. Ideal Situation as letter-grade steps, and PointsBoost_* as flat influence, are the readings that best match the stored values against known baselines (especially Making Friends vs a base complimentary bonus of 5).
