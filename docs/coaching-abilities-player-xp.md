---
title: "CFB 27 Coaching Abilities: Player XP"
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

**Full ability list (all trees):** [Coaching Abilities: Complete Guide](coaching-abilities.md).

## Executive summary

- Coach XP talents stack by **adding** HC + OC + DC, then applying once.
- Percent talents use `xp = xp * (100 + totalBoost) / 100`; flat talents add XP to the grant.
- Biggest percent stacks: **Star Maker / Whisperer** up to **+17%**; **Everybody Eats / Put in Work / Field Study** up to **+10%**.
- Notable flats: **Can't Stop, Won't Stop** up to **+600 XP** per 3-game win streak; **Show Don't Tell** **+1000 XP** per visiting prospect on a win (per staff role that owns it).
- Offseason training XP formula and base ranges: [Player XP, Levels & Dev Traits](./player-xp-levels-devtraits.md).

---

## How stacking works

Each talent stores three tier values:

| Field | Role |
|-------|------|
| `Data` | Head coach |
| `OCData` | Offensive coordinator |
| `DCData` | Defensive coordinator |

Staff contributions **add**, then apply once. They are not chained multipliers.

For position-split abilities, the coordinator that owns that side of the ball gets the larger tier (DC for DB/LB/DL; OC for QB/RB/WR-TE/OL/K-P). The off-side coordinator still contributes the smaller value when present.

**Max stack** = HC + OC + DC when all three trees own the ability.

Percent boosts apply as a single multiplier:

```
xp = xp * (100 + totalBoost) / 100
```

Flat XP bonuses add to the grant.

---

## Player XP abilities

| Ability | Effect | HC | OC | DC | Max stack |
|---------|--------|---:|---:|---:|----------:|
| Everybody Eats | all XP gains | 5% | 2% | 3% | **10%** |
| Put in Work | offseason training | 5% | 2% | 3% | **10%** |
| Star Maker | starters gain XP faster | 10% | 2% | 5% | **17%** |
| Whisperer | underclassmen progress | 10% | 2% | 5% | **17%** |
| Field Study | in-game goal XP | 5% | 2% | 3% | **10%** |
| Can't Stop, Won't Stop | per 3-game win streak | 300 XP | 100 XP | 200 XP | **600 XP** |
| Max Gains | skill-group max | 400 XP | 400 XP | 400 XP | flat |
| Let's Grow | player levels up | 15 XP | 15 XP | 15 XP | flat |
| Grow Together | player levels up | 7 XP | 7 XP | 7 XP | flat |
| Pay it Forward | same-pos player drafted | 4000 XP | varies | varies | flat |
| Show Don't Tell | visiting prospect on a win | 1000 XP | 1000 XP | 1000 XP | flat |
| Limitless | chance to raise skill cap on level-up | 10 | 2 | 5 | **17** |
| Put a Ring on It | chance to raise best cap on title | 50 | 25 | 25 | **100** |
| Senior Superlatives | +1 all skill-group caps (rising SR) | 1 | 0 | 0 | HC only |

OC/DC values swap by side of ball. Defense gets the larger DC tier; offense gets the larger OC tier. Numbers above show the defense-side OC/DC split for the percent abilities (Everybody Eats / Put in Work / Field Study: OC=2 DC=3; Star Maker / Whisperer: OC=2 DC=5; Can't Stop: OC=100 DC=200). Flip OC and DC for offense groups.

Pay it Forward coordinator values vary by position (e.g. DB: OC 1000 / DC 1500; QB: OC 4000 / DC 1000). See [Coaching Abilities: Complete Guide](coaching-abilities.md).

---

## Examples

**Everybody Eats on all three trees for QBs**  
Total = 5 + 3 + 2 = **+10%**. A 1000 XP grant becomes 1100.

**Star Maker HC only for starters**  
Total = **+10%**. A 1000 XP grant becomes 1100. Full staff stack (17%) becomes 1170.

**Can't Stop HC + both coordinators**  
Total = 300 + 100 + 200 = **+600 XP** flat per three-game win streak (not a percent).

---

## Related: Gainz Getter (CEO)

`Training_ChanceBoost_DevTrait` = **+20** (HC only; schema max is 20).

Readable FranTk formula:

```
chanceBoost = 100 + Training_ChanceBoost_DevTrait
increaseChance = increaseChance * chanceBoost / 100
```

Gainz Getter is a **1.20x** relative multiplier on the player's computed development-trait increase chance. It adds 20 to the 100-base `chanceBoost` scale rather than 20 percentage points to the final roll.
