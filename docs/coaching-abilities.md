---
title: "CFB 27 Coaching Abilities: Complete Guide"
---

Derived from `dynasty-tuning-binary.FTC` (CFB 27), Head Coach talent tree order.  
Values identical in patch 2 and patch 3.

Source: `Talent` / `TalentDataInt` / `TalentDataTeamRelativeAbility`.  
Runtime: `Team.CoachTalentEffects`.  
Win-streak length: `PlayerXPBoostWinStreakGameCount` = **3**.

## Executive summary

- Full HC talent tree with **CP costs** and **HC / OC / DC** values in unlock order (coordinators: trees 1–11 only).
- Staff stacks **add** (HC + OC + DC), then apply once. The multipliers are not chained.
- Position-split trees are bought per group; matching coordinator gets the larger tier (flip OC↔DC for defense).
- Program Builder / CEO store relative ability IDs (0–4 style tiers), not raw percents. Their sections explain each value.
- Focused deep dives: [Player XP](./coaching-abilities-player-xp.md) · [Recruiting](./coaching-abilities-recruiting.md).

---

## How to read this guide

Trees are listed in **Head Coach unlock order**. Coordinators get trees 1–11 (no Program Builder / CEO).

| Column | Meaning |
|--------|---------|
| CP | Coach Points to purchase |
| HC / OC / DC | Values from that staff role |

**Stacking:** HC + OC + DC **add**, then apply once (not chained multipliers).

**Position-split trees** (Motivator through Talent Developer, except Scheme Guru): each ability is bought per position group (QB, RB/FB, WR/TE, OL, DL, LB, DB, K/P). Tables show the **offense-group** OC/DC split (matching OC larger). **Flip OC↔DC for defense groups.** Notes identify cases where K/P follows the defense-side split.

**Percent XP:** `xp = xp * (100 + totalBoost) / 100`  
**Flat XP:** added to the grant.

**Scheme Guru / composure gameplay abilities** store `TalentDataInt = 0`. They do **not** expose a percent in FranTk. At runtime they write an integer tier (**0–7**) onto `Team.CoachTalentEffects` (`Offense_*` / `Defense_*` fields). The gameplay meaning of each tier (fatigue %, shed chance, etc.) is handled in native code; see §6.

**Deep dives:** [Player XP talents](coaching-abilities-player-xp.md) · [Recruiting talents](coaching-abilities-recruiting.md) · [Offseason XP math](player-xp-levels-devtraits.md) · [Recruiting tunables](recruiting-tunables-math.md)

---

## 1. Motivator

*Culture setter who is strong in player development*

**Archetype:** Grow Together. 25 CP gives **+7 XP** whenever a player levels up (HC/OC/DC all 7).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Extra Gear | 15 | Tough+2, Stam+2, Inj+2 | see note | see note | rating boosts while talent owned |
| Hot Hand | 20 | 0 | 0 | 0 | stay Hot through timeouts / quarter / games |
| Put in Work | 25 | 5% | 3% | 2% | offseason training XP |
| Locked In | 30 | 0 | 0 | 0 | start 4th quarter Hot in close games |

Extra Gear matching coord ≈ Tough+1, Stam+2, Inj+1; off-side ≈ Tough+0, Stam+1, Inj+0 (defense flips).

K/P costs are lower for Hot Hand (10) / Locked In (15).

---

## 2. Tactician

*Puts players in position to succeed on gameday*

**Archetype:** Winner. 25 CP gives **+450 coach XP per win** (HC/OC/DC all 450).

Each position branch has four **20 CP** rating-boost nodes (`TeamAbilityBoost`). Typical pattern: **HC +2 / matching coord +2 / off-side +1** (K/P uses +1 across the board).

### Passing Game (QB)

| Ability | HC ratings |
|---------|------------|
| Mobile QB Ratings Boost | Throw on Run +2, COD +2 |
| Passing QB Ratings Boost | Break Sack +2, Throw Under Pressure +1 |
| Mobile QB Ratings Boost 2 | Spin +2, Juke +2, BC Vision +2 |
| Passing QB Ratings Boost 2 | Short Acc +2, Mid Acc +2 |

### Running Game (RB/FB)

| Ability | HC ratings |
|---------|------------|
| RB/FB Ratings Boost | Pass Block +2, Run Block +2 |
| RB/FB Ratings Boost 2 | Trucking +2, Break Tackle +2 |
| RB/FB Ratings Boost 3 | BC Vision +2, Juke +2, Spin +2 |
| RB/FB Ratings Boost 4 | Catching +2, Short Route +2 |

### Receiving Game (WR/TE)

| Ability | HC ratings |
|---------|------------|
| WR/TE Ratings Boost | Run Block Power +2, Run Block Finesse +2 |
| WR/TE Ratings Boost 2 | Juke +2, Spin +2 |
| WR/TE Ratings Boost 3 | Short Route +2, Medium Route +2 |
| WR/TE Ratings Boost 4 | Catching +2, Catch in Traffic +2 |

### Blocking (OL)

| Ability | HC ratings |
|---------|------------|
| Pass Blocking OL | Pass Block +2, Impact Block +2 |
| Run Blocking OL | Run Block +2, Lead Block +2 |
| Finesse OL | Pass Block Finesse +2, Run Block Finesse +2 |
| Power OL | Run Block Power +2, Pass Block Power +2 |

### Defensive Line / Linebackers / Secondary

Same +2 HC pattern on paired ratings (moves, shed, coverage, tackle, pursuit, play recog, etc.). Matching DC gets +2; OC gets +1.

### Specialists (K/P)

| Ability | HC / OC / DC |
|---------|--------------|
| K/P Rating Boost 1–2 | Kick Accuracy **+1** |
| K/P Rating Boost 3–4 | Kick Power **+1** |

---

## 3. Recruiter

*Good at talent acquisition*

**Archetype:** Firm Handshakes. 25 CP gives **+100 coach XP** per signed recruit (HC/OC/DC all 100).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Advanced Look | 15 | 1 | 1 | 0 | faster scouting (steps) |
| Persuasive Personality | 20 | 30 | 15 | 0 | sway chance points |
| Magnetic Personality | 25 | 10 | 4 | 2 | HS starting interest |
| Portal King | 30 | 10 | 4 | 2 | transfer starting interest |

Max stacks (offense): Persuasive **45**, Magnetic/Portal King **16**.

---

## 4. Master Motivator

*Can get anyone to run through a brick wall*

**Archetype:** Let's Grow. 15 CP gives **+15 XP** on every player level-up (HC/OC/DC all 15).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Quick Recovery | 25 | 15 | 5 | 5 | faster wear-and-tear recovery |
| Icy Veins | 30 | 0 | 0 | 0 | start games with composure boost |
| Family Atmosphere | 35 | 20 | 10 | 5 | less likely to transfer |
| Everybody Eats | 40 | 5% | 3% | 2% | all XP gains |

Max: Quick Recovery **25**, Family Atmosphere **35**, Everybody Eats **10%**.

---

## 5. Architect

*Master at getting the most out of their players within a scheme*

**Archetype:** Max Gains. 15 CP gives **+400 XP** when a player maxes a skill group (HC/OC/DC all 400).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Can't Stop, Won't Stop | 25 | 300 | 200 | 100 | flat XP every 3-game win streak |
| Star Maker | 30 | 10% | 5% | 2% | starters gain XP faster |
| Limitless | 35 | 10 | 5 | 2 | chance to raise random skill cap on level-up |
| Put a Ring on It | 40 | 50 | 25 | 25 | chance to raise best skill cap on conf/nat title |

Max: Can't Stop **600 XP**, Star Maker **17%**, Limitless **17**, Put a Ring **100**.

---

## 6. Scheme Guru

*A master of the chess game*

**Archetype:** Winners Win. 15 CP gives **+900 coach XP per win** (HC/OC/DC all 900).

Scheme branches are **not** position-split. Abilities cost **30 CP**.

### Why there are no hard percentages here

Unlike XP / recruiting talents (where HC/OC/DC are readable ints like `5%` or `+100`), every Scheme Guru gameplay node stores **`TalentDataInt = 0`** for HC, OC, and DC. FranTk does not contain “Battery Pack = −X% fatigue.”

What *is* readable:

1. **Effect field:** which `CoachTalentEffects` meter the ability hits (e.g. `Offense_ReducedTeamFatigue`).
2. **Runtime tier:** owning the ability updates that meter on `Team.CoachTalentEffects` to an integer **0–7** (schema max). Live saves show tiers commonly at **1 / 2 / 4**, sometimes stacked up to **7** when staff share related effects.
3. **UI copy:** “slightly” vs “significantly” pairs use **different** meters (e.g. `Offense_HFAProtect` vs `Offense_HFAProtect2`), not a larger int on the same talent.
4. **Shared meters:** some abilities write the **same** field (e.g. Cardio Kings and He Shed both use `Defense_DLineWinChance`), so they can stack on one tier.

What is **not** in tuning: the conversion from tier → in-game magnitude. That mapping is native. Getting real percents needs binary RE or controlled A/B tests.

Tables below list each ability with its effect meter and the in-game description.

### Fast Tempo Offense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Battery Pack | `Offense_ReducedTeamFatigue` | offense fatigues slower in hurry-up |
| Caught Napping | `Offense_DefNotLookingBonus` | delay before defenders look to sideline at snap |
| On Their Heels | `Offense_QuickSnapBonus` | composure boost on 1st downs in hurry-up |
| Tipped Your Hand | `Offense_SeeDefCoverageInc` | chance to see coverage shell in hurry-up |

### Pass Game Offense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Take what's there | `Offense_ShortPassBoost` | catch boost ≤10 yards |
| Chew Up Yards | `Offense_MedPassBoost` | catch boost ≤20 yards |
| Stretch the Field | `Offense_LongPassBoost` | catch boost anywhere |
| We See You | `Offense_OLBlitzDetection` | post-snap blitz detection |

### Ground & Pound Offense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Wear 'Em Down | `Offense_DefFatigueDLOnCarry` | DL fatigue hit every 4+ yard rush |
| Demoralizing | `Offense_DefFatigueBoxDefOnCarry` | DL & LB fatigue hit every 4+ yard rush |
| Keep Blocks Longer | `Offense_ImprovedDisengage` | defenders disengage from run blocks less |
| To the Whistle | `Offense_HoldBlocksLonger` | defenders shed run blocks less |

### Discipline: Offense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Hater Blockers | `Offense_HFAProtect` | slight road crowd-noise reduction |
| Teflon | `Offense_HFAProtect2` | large road crowd-noise reduction |
| Polished | `Offense_PenaltyProtect` | slightly fewer offensive penalties |
| Clean Sheet Offense | `Offense_PenaltyProtect2` | significantly fewer offensive penalties |

### Fast Tempo Defense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Gas Tank | `Defense_ReducedTeamFatigue` | defenders fatigue slower |
| Not So Fast | `Defense_QuickSnapProtection` | faster defender sideline look at snap |
| Cardio Kings | `Defense_DLineWinChance` | DL keep shed ability vs hurry-up |
| Camouflaged Coverage | `Defense_RunDisguisedBoost` | harder to see shell in hurry-up |

### Pass Game Defense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Puncher's chance | `Defense_ShortPassDef` | knockout chance ≤10 yards |
| No Mid Here | `Defense_MedPassDef` | knockout chance ≤20 yards |
| Nothing Deep | `Defense_LongPassDef` | knockout chance anywhere |
| Master of Disguise | `Defense_OLBlitzDef` | protection from post-snap blitz detection |

### Ground & Pound Defense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Keep 'Em Coming | `Defense_CntrDefFatigueDLOnCarry` | DL fatigue slower on runs |
| Fresh Legs | `Defense_CntrDefFatigueBoxDefOnCarry` | DL & LB fatigue slower on runs |
| Hands Off | `Defense_OffensiveDisengageLimit` | better disengage on runs |
| He Shed | `Defense_DLineWinChance` | less time between shed attempts (same meter as Cardio Kings) |

### Discipline: Defense

| Ability | Effect meter | What it does (UI) |
|---------|--------------|-------------------|
| Unfriendly Confines | `Defense_HFABoost` | slight home crowd-noise boost |
| Too Loud??? | `Defense_HFABoost2` | large home crowd-noise boost |
| No Free Yards | `Defense_PenaltyProtect` | slightly fewer defensive penalties |
| Clean Sheet Defense | `Defense_PenaltyProtect2` | significantly fewer defensive penalties |

---

## 7. Strategist

*Perfect blend of X's and O's and talent acquisition*

**Archetype:** Show Don't Tell. 15 CP gives **+1000 player XP** per visiting prospect when you win (HC/OC/DC all 1000).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Making Friends | 25 | 20 | 15 | 15 | complimentary-visit influence |
| Hospitality | 30 | 15 | 5 | 5 | visit impact |
| Lower the Bar | 35 | 2 | 1 | 0 | lower dealbreaker threshold |
| Mind Reader | 40 | 30 | 30 | 15 | scouting trait-reveal chance |

Max: Making Friends **50**, Hospitality **25**, Mind Reader **75**.

---

## 8. Elite Recruiter

*The best of the best at talent acquisition*

**Archetype:** Firmer Handshakes. 15 CP gives **+200 coach XP** per signed recruit (HC/OC/DC all 200).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Ideal Situation | 25 | 3 | 3 | 3 | ideal-pitch letter-grade steps |
| Upsell | 30 | 10 | 5 | 5 | flat points from My School grades |
| Most Influential | 35 | 5 | 3 | 2 | flat action influence |
| Always Be Crootin' | 40 | 15 | 10 | 5 | weekly recruiting hours |

Max: Ideal Situation **9** grades, Upsell **20**, Most Influential **10**, Always Be Crootin' **30** hours.  
(Defense / K/P Always Be Crootin' uses OC **5** / DC **10**.)

K/P Always Be Crootin' costs **20 CP**.

---

## 9. Talent Developer

*Great at acquiring and developing talent*

**Archetype:** Draft Dividends. 15 CP gives **+2000 coach XP** when a player is drafted (HC/OC/DC all 2000).

| Ability | CP | HC | OC | DC | What it does |
|---------|---:|---:|---:|---:|--------------|
| Field Study | 25 | 5% | 3% | 2% | in-game goal XP |
| Whisperer | 30 | 10% | 5% | 2% | FR/SO progress faster |
| Home Sweet Home | 35 | 100 | 50 | 50 | pipeline action bonus |
| Pay it Forward | 40 | *varies* | *varies* | *varies* | player XP when same-pos drafted |

Max: Field Study **10%**, Whisperer **17%**, Home Sweet Home **200**.

K/P Whisperer / Field Study / Pay it Forward cost **15 / 25 / 20 CP**.

### Pay it Forward (flat player XP)

| Pos | HC | OC | DC |
|-----|---:|---:|---:|
| QB | 4000 | 4000 | 1000 |
| RB | 5000 | 2000 | 1000 |
| WR/TE | 4000 | 1500 | 1000 |
| OL | 4000 | 1500 | 1000 |
| DL / LB / DB | 4000 | 1000 | 1500 |
| K/P | 7000 | 3000 | 1000 |

---

## 10. Rainmaker

*Exclusively for MVP+ members.* Available to HC and coordinators.

Not position-split: every role stores the **same** int. Staff contributions **add** into one `CoachTalentEffects` field (all of these are `0–100` CTE ints). Max stack = HC + OC + DC when all three own the node.

| Ability | CP | HC | OC | DC | Max stack | What each role's value is | Effect |
|---------|---:|---:|---:|---:|----------:|---------------------------|--------|
| Deal Maker | 15 | **10** | 10 | 10 | **30** | **+10 flat recruiting influence** from NIL offers | `ProgramPoints_RecruitingNILIncreaseInfluence` |
| Budget Booster | 15 | **5** | 5 | 5 | **15** | **+5** on Dynasty Points earned from My School grades (0–100 CTE; best reading **+5%**, same scale as XP % talents) | `ProgramPoints_IncreasePointsFromGrades` |
| Staying Power | 15 | **5** | 5 | 5 | **15** | **+5** weight on how much a NIL offer reduces leave/transfer risk (UI: “increase NIL impact on Risk of Leaving”) | `ProgramPoints_DecreaseRiskOfTransfer` |
| Contract Incentives | 15 | **10** | 10 | 10 | **30** | **+10** on Dynasty Points from AD expectation goals (0–100 CTE; best reading **+10%**) | `ProgramPoints_IncreaseRewardContract` |

**Deal Maker** is the clearest unit: same family as other `*Influence*` recruiting boosts (flat interest/influence points). SoftSell base **20** + HC Deal Maker alone → **+10** NIL-piece influence on top of whatever the offer already granted. Offer÷expectation spline: [NIL Offer → Recruiting Influence](./nil-offer-recruiting-influence.md).

**Budget Booster / Contract Incentives:** FranTk does not expose the Dynasty Point apply formula. Stored values sit on the same `0–100` CTE scale as percent XP talents (`Practice_IncreaseXPEarned`, `Coach_XPBoost_GoalComplete`), so percent-of-award is the reading that best matches the schema; treat exact % vs flat DP as inferred until A/B confirmed.

**Staying Power:** despite the `DecreaseRiskOfTransfer` effect name, in-game text describes **amplifying NIL’s effect** on leave risk rather than a standalone leave-chance cut like Roster Retention abilities.

---

## 11. Visionary

*Awarded for creating a coach in Madden 27.* Available to HC and coordinators.

Same pattern as Rainmaker: equal HC/OC/DC ints, additive stack, `0–100` CTE fields.

| Ability | CP | HC | OC | DC | Max stack | What each role's value is | Effect |
|---------|---:|---:|---:|---:|----------:|---------------------------|--------|
| Pro Pipeline | 15 | **10** | 10 | 10 | **30** | **+10 draft-stock points** (0–100 CTE; round mapping is native) | `PlayersLeaving_IncreaseDraftStock` |
| Practice Makes Perfect | 15 | **5** | 5 | 5 | **15** | **+5% practice XP** (`xp * (100 + total) / 100`) | `Practice_IncreaseXPEarned` |
| Hot Start | 15 | **5** | 5 | 5 | **15** | **+5 chance points** to start Hot from practice | `Practice_IncreaseHotChance` |
| Signing Bonus | 15 | **10** | 10 | 10 | **30** | **+10 annual coach points** awarded by the program / on contract cycle | `CoachCarousel_IncreaseCoachPoints_NewContract` |

**Practice Makes Perfect** mirrors other percent XP talents (Everybody Eats, Put in Work): HC alone is a small bump; full staff is **+15%** on practice XP grants.

**Hot Start** adds chance points to the Hot composure roll from practice. It does not guarantee Hot and is separate from CEO Gasoline (big-game Hot-start flag with `TalentDataInt = 0`).

**Pro Pipeline** raises NFL draft stock on the `0–100` CTE field; how many stock points move a player across draft rounds is native.

**Signing Bonus** is a flat CP add to the annual / new-contract coach-point award (description: “annual Coach Points award from your program”), not a % discount like Bundle Discount / Friends & Family.

---

## 12. Program Builder (HC only)

*Head coach who can oversee everything*

**Archetype:** Winning Together. 70 CP gives **+75 flat coach XP** when a school grade or prestige increases (`School_XPBoost_Prestige`).

HC-only tree. There is **no single unit** for the HC column. Each ability stores a raw `TalentDataInt` used differently by native code. Readings below come from effect names, schema caps, in-game text, and comparison to known recruiting/XP baselines.

### Coaching Family

Effect family: coordinator contracts / hire pool.

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Forever Home | 40 | **50** | chance points reducing how often coordinators get poached | `Contracts_Coordinators_LessPoached` |
| Deal Sweetener | 40 | **4** | chance points increasing coordinator offer acceptance | `Contracts_Coordinators_AcceptMoreOften` |
| Cream of the Crop | 40 | **4** | hire-pool quality weight (not a literal candidate count) | `Contracts_Coordinators_BetterCandidates` |

Exact % conversion for the chance/weight ints is native (not in FranTk).

### Set The Tone

HFA momentum fields are readable **0–100** ints (unlike Scheme Guru composure nodes that store `0` and only write CTE tiers).

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Home Field Advantage | 40 | **10** | +10 HFA momentum vs away team in home games | `HFA_BoostMomentum_HomeGames` |
| Road Warriors | 40 | **10** | +10 counter to opponent HFA momentum on the road | `HFA_CounterMomentum` |
| Not Intimidated | 40 | **10** | +10 HFA momentum help in rivalry games | `HFA_BoostMomentum_RivalGames` |

### Roster Retention

Effect family: NFL-leave chance reduction by draft band. HC ints are **chance points** subtracted from leave likelihood (exact % mapping native).

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Best Advice | 40 | **10** | −10 leave-chance points for **3rd–7th** round projections | `PlayersLeaving_LeaveForProsChanceLate` |
| Let's Run It Back | 40 | **10** | −10 leave-chance points for **2nd** round | `PlayersLeaving_LeaveForProsChance2nd` |
| Delay Sunday | 40 | **20** | −20 leave-chance points for **1st** round | `PlayersLeaving_LeaveForProsChance1st` |

### High Integrity

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Gift of Gab | 40 | **3** | **+3 persuade attempts** per season (literal count) | `PlayersLeaving_AdditionalPersuadeAttempts` |
| Roster Retainer | 40 | **10** | **+10 chance points** on persuade success rate | `PlayersLeaving_BoostPersuadeRate` |
| Full Refund | 40 | **1** | **bool on**; successful persuade refunds the attempt | `PlayersLeaving_FreePersuade` |

### Inside Track

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Smart Goals | 40 | **25** | **+25% coach XP** from completed goals (`xp * (100+25)/100`) | `Coach_XPBoost_GoalComplete` |
| Portal Preview | 40 | **1** | **bool on**; see other schools' At Risk players | `Recruiting_SeeAtRiskPlayers_OtherTeams` |
| Friends & Family Discount | 40 | **50** | stored weight **50**; in-game text: archetypes cost **−20 CP** per staff coach who owns that archetype | `Talents_CoachPointDiscount_NewArcheType` |

For Friends & Family, treat **−20 CP** (UI) as the effective discount per matching coach. The stored `50` does not equal that CP amount 1:1.

### Relationship Builder

Here the HC value **is the pipeline count / tier amount**.

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Invest | 40 | **2** | your **top 2** pipelines each get **+1 tier** | `PIpeline_LevelBoost_Top2` |
| Making Inroads | 40 | **5** | your **top 5** pipelines each get **+1 tier** | `PIpeline_LevelBoost_Top5` |
| Share the Wealth | 40 | **5** | your **bottom 5** pipelines each get **+1 tier** | `PIpeline_LevelBoost_Bottom5` |

### Strong Roots

| Ability | CP | HC | What the HC value is | Effect |
|---------|---:|---:|----------------------|--------|
| Giving Back | 40 | **1** | **+1 tier** to primary (alma mater) pipeline | `PIpeline_LevelBoost_AlmaMater` |
| Household Name | 40 | **10** | **+10 flat starting interest** for primary-pipeline recruits | `Recruiting_InfluenceBoost_PipelineStart` |
| Hometown Discount | 40 | **25** | **+25 flat influence** on recruiting actions vs primary-pipeline recruits | `Recruiting_Action_InfluenceBoost_Pipeline` |

Worked examples for Household Name / Hometown Discount: [Recruiting talents](coaching-abilities-recruiting.md#program-builder-hc-only).

---

## 13. CEO (HC only)

*Elite head coach who has won a National Title*

**Archetype:** Big Game Bonus. 90 CP gives **+12,500 flat coach XP** for playoff wins (`GameResults_XPBoost_BigWin`).

All ability nodes **40 CP**. Same rule as Program Builder: HC is raw `TalentDataInt`; units differ per ability.

| Ability | HC | What the HC value is | Effect |
|---------|---:|----------------------|--------|
| Gainz Getter | **20** | added into `chanceBoost = 100 + 20` → **1.20×** relative multiplier on dev-trait upgrade chance | `Training_ChanceBoost_DevTrait` |
| Second Chance Keeper | **1** | **bool on**; second persuade if the first fails | `PlayersLeaving_SecondPersuade` |
| Lasting Impression | **5** | **+5 flat interest** per **10 hours** spent on a recruit | `Recruiting_PointsBoost_Accelerate` |
| More the Merrier | **4** | **+4 weekly visit slots** on top of base `MaxRecruitVisitsPerWeek` (**4**) → **8** | `Recruiting_Visits_Increase` |
| Gasoline | **0** | composure Hot-start flag (magnitude not in FranTk; CTE/composure pattern) | `Composure_HotStart_BigGame` |
| Last Dance | **100** | UI: **+100 interest** per prior visit elsewhere before yours; effect name is `Recruiting_BoostChance_DevTraitUnlock` (chance-style), while the interest reading matches UI text | `Recruiting_BoostChance_DevTraitUnlock` |
| Dream School | **10** | **+10 chance points** to instant-commit when you offer as their #1 school | `Recruiting_BoostCommitChanceOn1st` |
| Bundle Discount | **5** | stored weight **5**; in-game text: abilities cost **−2 CP** per staff coach who already owns that ability | `Talents_CoachPointDiscount_SameTalents` |
| Senior Superlatives | **1** | **+1 to every skill-group cap** for rising seniors | `Training_SkillGroupIncrease_Seniors` |

**Gainz Getter** (FranTk-readable):

```
chanceBoost = 100 + 20  // = 120
increaseChance = increaseChance * chanceBoost / 100
```

This is +20 on the 100-base `chanceBoost` scale, not 20 percentage points on the final roll.

**Lasting Impression** (`RecruitingPointsBoostAccelerate_Divisor` = 10, weekly hour cap 50):  
10h → +5 · 30h → +15 · 50h → +25 interest on top of the hours actions themselves.

**Dream School / Lasting Impression** examples: [Recruiting talents](coaching-abilities-recruiting.md#ceo-hc-only).  
**Gainz Getter** detail: [Player XP talents](coaching-abilities-player-xp.md#related-gainz-getter-ceo).

---

## Quick XP stack (player progression)

Best trees for raw player XP rate:

| Tree | Ability | Max contribution |
|------|---------|------------------|
| Master Motivator | Everybody Eats | +10% all XP |
| Motivator | Put in Work | +10% offseason |
| Architect | Star Maker | +17% starters |
| Talent Developer | Whisperer | +17% FR/SO |
| Architect | Can't Stop | +600 XP / 3-win streak |
| CEO | Gainz Getter | 1.20× trait-upgrade chance |
| Architect / Motivator / Master Mot. | Max Gains / Grow Together / Let's Grow | flat milestone XP |

Additive % ceiling on an underclassman starter offseason grant:  
Everybody Eats 10 + Put in Work 10 + Star Maker 17 + Whisperer 17 = **+54%**, then × DevTrait.

---

## Caveats

1. Talent definitions unchanged patch 2 → patch 3.
2. Program Builder and CEO are **HC-only**; coordinators stop at Visionary.
3. Scheme Guru gameplay abilities store `TalentDataInt = 0`. They update `CoachTalentEffects` tiers (0–7); FranTk does not expose the tier→gameplay magnitude. Motivator composure nodes (Hot Hand, Locked In, Icy Veins, Gasoline) are the same pattern.
4. Program Builder HFA abilities (`Home Field Advantage`, etc.) are different: those use 0–100 fields with readable talent values (e.g. **10**).
5. Some `Talent` rows exist in tuning but are **not** on any OrderedTalentNodeList (e.g. Cross train, Film Junkies, Reach Your Potential, Legacy), so they are omitted here.
6. Deeper baselines: [Recruiting talents](coaching-abilities-recruiting.md), [Recruiting tunables](recruiting-tunables-math.md), [Offseason XP math](player-xp-levels-devtraits.md), [Player XP talents](coaching-abilities-player-xp.md).

