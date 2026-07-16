---
title: "CFB 27 Recruiting Tunables -- Weights and Math Reference"
geometry: margin=0.75in
header-includes:
  - \usepackage{longtable}
  - \usepackage{booktabs}
  - \usepackage{array}
  - \usepackage{etoolbox}
  - \AtBeginEnvironment{longtable}{\small}
  - \sloppy
---

Derived from `dynasty-tuning-binary.FTC` (CFB 27, patch `cfb27-db-data/2/`).  
Source command: `cfb-dynasty recruiting-tunables -schema-dir ./data`

This document records tuning constants and inferred influence arithmetic only.  
Compiled formula bytecode is not available; evaluation order below is inferred from schema field structure (`RecruitingActionInfo`, `RecruitTarget`, `ProspectTargetSchool`).

---

## Save fields vs tuning constants

| Role | Save export field | Table / notes |
|------|-------------------|---------------|
| Per-school influence | `recruiting[].schoolInterest[].influence` | `ProspectTargetSchool.TeamInfluence` |
| Total influence pool | `recruits[].commitScore` | `Recruit.CommitScore` |
| Recruit stage | `recruits[].recruitStage` | `Recruit.RecruitStage` |
| Weekly influence change | `recruiting[].prospectInfluenceDelta` | `RecruitTarget.ProspectInfluenceDelta` |
| Cumulative influence | `recruiting[].prospectInfluenceTotal` | `RecruitTarget.ProspectInfluenceTotal` |
| Hours this week | `recruiting[].prospectHoursSpentCurrent` | `RecruitTarget.ProspectHoursSpentCurrent` |
| School letter grades | `schoolGrades[]` | `MySchoolTrackingTable` |
| Recruit motivations | `recruits[].player.motivations[]` | `Recruit` motivation fields |
| Active pitch + intensity | `recruiting[].activePitches[]` | `ActiveRecruitingPitch` |

School share of commit score is computed in-game by `RecruitingStageDetails_GetSchoolPercentageOfCommitScore(recruit, schoolRank)` (expression asset; bytecode not decoded).

---

## Influence model (inferred)

Weekly influence gain for an hour-based action:

```
influenceGain = BaseInfluenceGranted
              + PipelineInfluenceBonusTable[pipelineLevel]
              + sum of MotivationGradeBonus(pitchMotivationSlot_i, schoolLetterGrade)
              + matchingMotivationBoosts (see ChanceToSwayBoostPerMatchingMotivation)
              + NonMatchingMotivationPenalty   (PITCH only; negative when pitch mismatches)
```

- `pipelineLevel` is indexed **0-5** (six entries in every `PipelineInfluenceBonusTable`).
- `MotivationGradeBonus` uses `LetterGradeValueTable` rows attached to the action (three rows for PITCH; one row per motivation slot).
- `NonMatchingMotivationPenalty` is **0** for non-pitch actions and for `PITCH` / `Sway`.

---

## Stage visibility thresholds

Percentage of commit score required for stage visibility (from `RecruitingTunables` scalars):

| Threshold | Value |
|-----------|------:|
| `PointsPercentageThresholdTop5` | 35 |
| `PointsPercentageThresholdTop3` | 75 |
| `PointsPercentageThresholdCommit` | 100 |

Stage school-count caps:

| Field | Value |
|-------|------:|
| `RecruitStageTop10SchoolCount` | 10 |
| `RecruitStageTop5SchoolCount` | 5 |
| `RecruitStageTop3SchoolCount` | 3 |
| `RecruitStageBattleSchoolCount` | 3 |
| `RecruitStageSoftCommittedSchoolCount` | 3 |
| `RecruitStageHardCommittedSchoolCount` | 1 |
| `RecruitStageSignedSchoolCount` | 1 |

Recruit stage enum values present in tuning (`RecruitingStageDetails`):  
`Top10`, `Top5`, `Top3`, `Battle`, `SoftCommitted`, `HardCommitted`

Battle proximity threshold: `TriggerRecruitingBattleWithinCommitThreshold` = **10**

Commit-score modifiers:

| Field | Value |
|-------|------:|
| `BattleAdditionalCommitScorePercentage` | 25 |
| `HardCommitAdditionalCommitScorePercentage` | 50 |
| `TransferDefaultCommitScore` | 1000 |

---

## Weekly limits

| Field | Value |
|-------|------:|
| `MaxTotalHoursOnRecruitPerWeek` | 50 |
| `MaxRecruitingBoardTargets` | 35 |
| `MaxTeamScholarshipOffers` | 35 |
| `MaxRecruitVisitsPerWeek` | 4 |
| `MaxTimesScouted` | 5 |
| `CPU_AI_MaxPitchCount` | 2 |

---

## Recruiting actions

From `RecruitingActionInfo` tuning table.

| Action | Intensity | Cost (hrs) | Base influence | Penalty |
|--------|-----------|----------:|---------------:|--------:|
| `SearchSocialMedia` | n/a | 5 | 4 | 0 |
| `SCOUTING` | n/a | 10 | 0 | 0 |
| `ContactHighSchoolCoaches` | n/a | 10 | 8 | 0 |
| `ContactFriendsAndFamily` | n/a | 25 | 20 | 0 |
| `OFFERSCHOLARSHIP` | n/a | 5 | 0 | 0 |
| `PITCH` | `SoftSell` | 20 | 20 | -20 |
| `PITCH` | `HardSell` | 40 | 40 | -40 |
| `PITCH` | `Sway` | 30 | 15 | 0 |
| `SCHEDULEVISIT` | n/a | 40 | 0 | 0 |
| `VisitRecruitsSchool` | n/a | 50 | 40 | 0 |
| `SendTheHouse` | n/a | 50 | 50 | 0 |

`SearchSocialMedia` aliases `GenericHoursAction_First`. `SendTheHouse` aliases `GenericHoursAction_Last`.

Base efficiency (`BaseInfluenceGranted / Cost`):

| Action | Ratio |
|--------|------:|
| `GenericHoursAction_First` | 0.80 |
| `ContactHighSchoolCoaches` | 0.80 |
| `ContactFriendsAndFamily` | 0.80 |
| `VisitRecruitsSchool` | 0.80 |
| `GenericHoursAction_Last` | 1.00 |
| `PITCH` SoftSell | 1.00 |
| `PITCH` HardSell | 1.00 |
| `PITCH` Sway | 0.50 |

### Pipeline influence bonus by action

Add `PipelineInfluenceBonusTable[pipelineLevel]` where `pipelineLevel` is in {0, 1, 2, 3, 4, 5}.

| Action | L0 | L1 | L2 | L3 | L4 | L5 |
|--------|---:|---:|---:|---:|---:|---:|
| `ContactHighSchoolCoaches` | 0 | 1 | 1 | 2 | 3 | 4 |
| `GenericHoursAction_First` | 0 | 1 | 1 | 1 | 1 | 2 |
| `ContactFriendsAndFamily` | 0 | 2 | 4 | 6 | 8 | 10 |
| `VisitRecruitsSchool` | 0 | 4 | 8 | 12 | 16 | 20 |
| `PITCH` SoftSell | 0 | 2 | 4 | 6 | 8 | 10 |
| `PITCH` HardSell | 0 | 4 | 8 | 12 | 16 | 20 |
| `PITCH` Sway | 0 | 0 | 0 | 0 | 0 | 0 |
| `GenericHoursAction_Last` | 0 | 0 | 0 | 0 | 0 | 0 |
| `OFFERSCHOLARSHIP`, `SCHEDULEVISIT`, `SCOUTING` | 0 | 0 | 0 | 0 | 0 | 0 |

---

## Pitch motivation grade bonuses

Each `PITCH` row has `MotivationGradesInfluenceBonusTable` with **three** `LetterGradeValueTable` rows (motivation slots 1-3, aligned with `AssociatedMotivation1..3` on `RecruitingPitchInfo`).  
When a recruit motivation matches a pitch motivation slot, the school's letter grade for the mapped `MySchoolGrade` (see below) selects a bonus from that slot's row.

`PITCH` / `Sway` has no motivation grade table attached in tuning.

Grade column order: F, D-, D, D+, C-, C, C+, B-, B, B+, A-, A, A+

### SoftSell slot 1

```
-6  -5  -4  -3  -2  -1   0   0   1   2   3   4   6
```

### SoftSell slot 2

```
-7  -6  -5  -4  -3  -2  -1   0   1   2   4   5   7
```

### SoftSell slot 3

```
-9  -8  -6  -5  -4  -2  -1   0   2   4   6   8  10
```

### HardSell slot 1

```
-14 -12 -10  -8  -6  -4  -1   1   2   4   8  10  14
```

### HardSell slot 2

```
-20 -18 -16 -14 -10  -6  -2   2   4   6  10  14  20
```

### HardSell slot 3

```
-12 -10  -8  -6  -4  -2  -1   0   2   4   6   8  12
```

Matching-motivation scalar: `ChanceToSwayBoostPerMatchingMotivation` = **15**

---

## Motivation to school grade mapping

From `RecruitingMotivationToMySchoolGradeMapping`. The school grade field used for pitch bonus lookup is the recruit motivation's mapped `MySchoolGrade` (exported under `schoolGrades`).

| Recruiting motivation | MySchoolGrade key |
|-----------------------|-------------------|
| AcademicPrestige | AcademicPrestige |
| AthleticFacilities | AthleticFacilities |
| BrandExposure | BrandExposure |
| CampusLifestyle | CampusLifestyle |
| ChampionshipContender | ChampionshipContender |
| CoachPrestige | CoachPrestige |
| CoachStability | CoachStability |
| ConferencePrestige | ConferencePrestige |
| PlayingStyle | PlayingStyle |
| PlayingTime | PlayingTime |
| ProgramTradition | ProgramTradition |
| ProPotential | ProPotential |
| ProximityToHome | ProximityToHome |
| StadiumAtmosphere | StadiumAtmosphere |

---

## Pitch catalog

From `RecruitingPitchInfo` (20 pitches). `hasDealbreakerMotivation` is stored on pitches that declare it in tuning.

| Pitch ID | Display name |
|----------|--------------|
| Aspirational | Aspirational Goals |
| CampusPersonality | Campus Personality |
| TimeToGetToWork | Clocked In |
| CoachsFavorite | Coach Connection |
| CollegeExperience | College Experience |
| ConferenceSpotlight | Conference Legend |
| FootballInfluencer | Football Influencer |
| ItsGameTime | Gamer |
| Grassroots | Grassroots Traditionalist |
| WorkHorse | Gym Rat |
| HometownHero | Hometown Hero |
| ToTheHouse | House Call |
| TVTime | Primetime Player |
| Prestigious | Standard Bearer |
| Starter | Star Search |
| ProveYourself | Status Seeker |
| StudentOfTheGame | Student of the Game |
| SundayBound | Sunday Player |
| TeamPlayer | Team Player |
| TheClutch | The Clutch |

Motivations by pitch (slots 1, 2, 3), part 1:

| Pitch ID | Motivation 1 | Motivation 2 | Motivation 3 |
|----------|--------------|--------------|--------------|
| Aspirational | ChampionshipContender | ConferencePrestige | CoachPrestige |
| CampusPersonality | CampusLifestyle | PlayingTime | PlayingStyle |
| TimeToGetToWork | PlayingStyle | PlayingTime | ProPotential |
| CoachsFavorite | AthleticFacilities | CoachPrestige | ProximityToHome |
| CollegeExperience | AcademicPrestige | CampusLifestyle | StadiumAtmosphere |
| ConferenceSpotlight | ChampionshipContender | ConferencePrestige | ProximityToHome |
| FootballInfluencer | BrandExposure | PlayingTime | ProPotential |
| ItsGameTime | ConferencePrestige | PlayingStyle | ProPotential |
| Grassroots | ProgramTradition | ProximityToHome | StadiumAtmosphere |
| WorkHorse | AthleticFacilities | BrandExposure | ProPotential |

Motivations by pitch, part 2:

| Pitch ID | Motivation 1 | Motivation 2 | Motivation 3 |
|----------|--------------|--------------|--------------|
| HometownHero | ChampionshipContender | ProgramTradition | ProximityToHome |
| ToTheHouse | BrandExposure | CoachPrestige | ChampionshipContender |
| TVTime | BrandExposure | ChampionshipContender | PlayingTime |
| Prestigious | CoachPrestige | ConferencePrestige | PlayingStyle |
| Starter | BrandExposure | PlayingTime | ProximityToHome |
| ProveYourself | BrandExposure | CoachPrestige | ConferencePrestige |
| StudentOfTheGame | AcademicPrestige | CoachPrestige | ProximityToHome |
| SundayBound | ChampionshipContender | ConferencePrestige | ProPotential |
| TeamPlayer | CoachStability | PlayingStyle | ProximityToHome |
| TheClutch | CoachStability | PlayingStyle | PlayingTime |

---

## Visit tunables

Scalars (`VisitTunables`):

| Field | Value |
|-------|------:|
| `competitiveVisitPenalty` | 5 |
| `complimentaryVisitBonus` | 5 |

`SCHEDULEVISIT` costs **40** hours with **0** base influence; influence from visit activities is applied separately via `VisitActivityInfluenceTable`.

### Visit activity types

From `VisitActivityInfo`:

AttendLecture, AttendPositionMeeting, AttendPractice, AttendTeamMeeting, CampusTour, FamilyVisit, MeetAlumni, OneOnOneCoaching, PodcastInterview, Tailgate, TeamDinner, TeamHistory, TeamWorkout, TrophyTour

### Visit influence by recruit interest level

Each activity maps to an interest level (`Dealbreaker`, `Interested`, `NotInterested`). The selected `LetterGradeValueTable` row provides a bonus from the school's grade for the activity's associated motivation.

Grade column order: F, D-, D, D+, C-, C, C+, B-, B, B+, A-, A, A+

**Dealbreaker**

```
-48 -40 -32 -24 -16  -8   0   8  16  24  32  40  48
```

**Interested**

```
-36 -30 -24 -18 -12  -6   0   6  12  18  24  30  36
```

**NotInterested**

```
-24 -22 -20 -18 -16 -14 -12 -10  -8  -6  -4  -2   0
```

---

## Scholarship and NIL

| Field | Value |
|-------|------:|
| `ScholarshipOfferWeeklyBonus` | 5 |
| `ScholarshipBonusLowExpectationCutoff` | 10 |
| `MinimumNILToOfferPercentage` | 0.80 |
| `SoftCommitInfluenceVariance` | 10 |

---

## Pipeline and initial interest arrays

Indexed arrays from `RecruitingTunables` (index meaning is game-defined; pipeline tables use indices **0-5**).

| Array | Values |
|-------|--------|
| `RecruitInitialInfluencePerPipelineLevel` | 0, 5, 15, 25, 75, 100 |
| `InitialInterest_Pipeline` | 0, 5, 10, 15, 25, 35 |
| `InfluenceRequiredPerPipelineLevel` | 0, 1, 30, 90, 175, 300 |
| `ChanceToSwayBoostFromPipelineInfluenceLevel` | 0, 5, 10, 15, 20, 25 |
| `InitialInterest_AlmaMater_Starlevel` | 10, 10, 10, 10, 10 |
| `InterestBoost_Start_Awards_Starlevel` | 10, 10, 10, 10, 10 |

---

## Scouting arrays

| Array | Values |
|-------|--------|
| `ScoutingAttributesUnlockPercentage` | 10 entries of 10 |
| `ScoutingPhysicalAbilityUnlockPercentage` | 5 entries of 20 |
| `ScoutingAttributeLongSliderWidth` (scalar) | 28 |
| `ScoutingAttributeShortSliderWidth` (scalar) | 13 |

---

## Instant commit and class ranking

| Array | Values |
|-------|--------|
| `InstantCommitOddsPerStarLevel` | 10, 8, 6, 4, 2 |
| `InstantCommitBonusPrestige` | 0, 1, 2, 5, 15, 25, 40, 70, 100 |
| `TopClassesStarRatingTable` | 5, 10, 15, 20, 25 |
| `TopClassesRankWeightPercentageTable` | see below |

`TopClassesRankWeightPercentageTable` (35 entries):

```
100, 99, 98, 95, 91, 86, 80, 74, 67, 61, 54, 47, 41, 35, 30, 25, 21, 17, 14, 11,
8, 7, 5, 4, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3
```

Transfer portal `TransferValueCutOffs`:

```
76, 76, 78, 78, 80, 80, 82, 82, 84, 84, 86
```

---

## Recruiting points abstraction

`RecruitingPointsAbstractionTable` maps raw influence totals to UI tiers (15 entries):

```
141, 100, 61, 41, 21, 11, 6, 1, 0, -5, -10, -20, -40, -60, -99
```

---

## High school recruit generation

From `HighSchoolRecruitingTunables` (single tuning row):

| Field | Value |
|-------|------:|
| AcademicsMeterBonusMax | 30 |
| BrandMeterBonusMax | 75000 |
| ChallengePerformanceThreshold | 10 |
| CoachTrustBonusMax | 1000 |
| FitnessMeterBonusMax | 16 |
| InitialTapeScoreBoostDecrementRate | 20 |
| LeadershipMeterBonusMax | 25 |
| NumberOfRecommendedSchoolScholarshipOffers | 3 |
| RankChangeExpectedPerformanceMax | 10 |
| RankChangeOverPerformanceMax | 25 |
| RankChangeUnderPerformanceMax | 25 |
| ScholarshipBonusDecommitPenalty | 1 |
| ScholarshipBonusVariance | 10 |
| ScholarshipScoreDecommitPenalty | 1500 |
| SchoolPrestigeFactorWeight | 1.0 |
| SkillPointsBonusMax | 600 |
| TeamNeedFactorWeight | 1.0 |

---

## Remaining scalars

| Field | Value |
|-------|------:|
| `AIAggressivenessSlider` | 1.15 |
| `AINILAggressivenessSlider` | 2.25 |
| `CPU_AI_SpendScoutingWeeklyMaxRecruitsToScout` | 35 |
| `ChanceToSwayBoostPerMatchingMotivation` | 15 |
| `EventRankingCutoff` | 10 |
| `MaxCollegeYearsRemaining` | 4 |
| `MaxFiveStarsInClass` | 32 |
| `MaxJumpBacks` | 0 |
| `OVRProjection` | 2 |
| `RecruitsToGenerateGroupMax` | 20 |
| `RecruitsToGenerateMinimum` | 4100 |
| `RecruitsToGenerateMaximum` | 4100 |
| `Stage2WarningThreshold` | 50 |

### UI strings

| Field | Text |
|-------|------|
| `EarlySigningDayActionItemString` | Early National Signing Day |
| `PreseasonRecruitingActionItemString` | Set Up Recruiting Board and Scout |
| `TransferPortalActionItemString` | Transfer Portal and Offseason Recruiting |

---

## Encoding note

Tuning FTC array elements use EA `s_int` storage with `minValue = 0x80000000`. Exported values in `cfb-dynasty recruiting-tunables` decode this bias automatically.

---

## Not decoded

- Compiled expression bodies (`RecruitingStageDetails_GetSchoolPercentageOfCommitScore`, `RecruitingEval_ApplyWeeklyRecruitingHours`, etc.)
- Exact stacking order of additive bonuses vs multipliers
- Per-activity motivation to interest-level mapping (which visit activity triggers which `VisitActivityInterestLevel`)
- `RecruitingActionFeedbackEntry` / `RecruitingActionBonus` runtime bonus breakdown on executed actions
