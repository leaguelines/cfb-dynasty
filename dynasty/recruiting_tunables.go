package dynasty

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// RecruitingTunablesExport is the game's recruiting formula constants from tuning data.
type RecruitingTunablesExport struct {
	Scalars          map[string]any                     `json:"scalars,omitempty"`
	Arrays           map[string][]int                   `json:"arrays,omitempty"`
	Strings          map[string]string                  `json:"strings,omitempty"`
	Pitches          []RecruitingPitchTunableExport     `json:"pitches,omitempty"`
	Visit            *VisitTunablesExport               `json:"visit,omitempty"`
	VisitActivities  []VisitActivityTunableExport       `json:"visitActivities,omitempty"`
	Actions          []RecruitingActionTunableExport    `json:"actions,omitempty"`
	HighSchool       *HighSchoolRecruitingTunablesExport `json:"highSchool,omitempty"`
	Stages           []RecruitingStageTunableExport     `json:"stages,omitempty"`
}

// RecruitingPitchTunableExport is one recruiting pitch definition from tuning data.
type RecruitingPitchTunableExport struct {
	Pitch                    string `json:"pitch,omitempty"`
	ShortName                string `json:"shortName,omitempty"`
	LongName                 string `json:"longName,omitempty"`
	HasDealbreakerMotivation bool   `json:"hasDealbreakerMotivation,omitempty"`
	AssociatedMotivation1    string `json:"associatedMotivation1,omitempty"`
	AssociatedMotivation2    string `json:"associatedMotivation2,omitempty"`
	AssociatedMotivation3    string `json:"associatedMotivation3,omitempty"`
	SwappableImage           *int   `json:"swappableImage,omitempty"`
}

// VisitTunablesExport holds visit-related recruiting constants.
type VisitTunablesExport struct {
	CompetitiveVisitPenalty  *int `json:"competitiveVisitPenalty,omitempty"`
	ComplimentaryVisitBonus  *int `json:"complimentaryVisitBonus,omitempty"`
}

// VisitActivityTunableExport is one campus visit activity from tuning data.
type VisitActivityTunableExport struct {
	ActivityType string `json:"activityType,omitempty"`
	IconID       *int   `json:"iconId,omitempty"`
}

// RecruitingActionTunableExport is one recruiting action (hours, contact, etc.).
type RecruitingActionTunableExport struct {
	ActionType           string `json:"actionType,omitempty"`
	Intensity            string `json:"intensity,omitempty"`
	Cost                 *int   `json:"cost,omitempty"`
	BaseInfluenceGranted *int   `json:"baseInfluenceGranted,omitempty"`
	IconID               *int   `json:"iconId,omitempty"`
	IsEnabled            bool   `json:"isEnabled,omitempty"`
}

// HighSchoolRecruitingTunablesExport holds high-school recruiting generation tunables.
type HighSchoolRecruitingTunablesExport struct {
	Scalars map[string]any `json:"scalars,omitempty"`
}

// RecruitingStageTunableExport names a recruiting pipeline stage.
type RecruitingStageTunableExport struct {
	RecruitStage string `json:"recruitStage,omitempty"`
}

// ExportRecruitingTunables reads recruiting constants from the tuning FTC bundle.
func ExportRecruitingTunables(schemaDir, tuningPath string) (RecruitingTunablesExport, error) {
	tf, err := openTuningFile(schemaDir, tuningPath)
	if err != nil {
		return RecruitingTunablesExport{}, err
	}
	table, ok := tf.PrimaryTableByName("RecruitingTunables")
	if !ok {
		return RecruitingTunablesExport{}, fmt.Errorf("cfb-dynasty: RecruitingTunables table not found")
	}
	if err := table.ReadRecords(); err != nil {
		return RecruitingTunablesExport{}, err
	}
	if len(table.Records) == 0 {
		return RecruitingTunablesExport{}, fmt.Errorf("cfb-dynasty: RecruitingTunables table empty")
	}
	out := recruitingTunablesFromRecord(tf, table.Records[0])
	out.Pitches = recruitingPitchTunables(tf)
	out.Visit = visitTunables(tf)
	out.VisitActivities = visitActivityTunables(tf)
	out.Actions = recruitingActionTunables(tf)
	out.HighSchool = highSchoolRecruitingTunables(tf)
	out.Stages = recruitingStageTunables(tf)
	return out, nil
}

func openTuningFile(schemaDir, tuningPath string) (*File, error) {
	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.TuningPath = tuningPath
	settings.AutoParse = true

	path := tuningPath
	if path == "" {
		path = discoverTuningPath(schemaDir)
	}
	if path == "" {
		return nil, fmt.Errorf("cfb-dynasty: tuning data not found")
	}
	return Open(path, &settings)
}

func recruitingTunablesFromRecord(f *File, record Record) RecruitingTunablesExport {
	out := RecruitingTunablesExport{
		Scalars: make(map[string]any),
		Arrays:  make(map[string][]int),
		Strings: make(map[string]string),
	}
	scalarInts := []string{
		"BattleAdditionalCommitScorePercentage",
		"ChanceToSwayBoostPerMatchingMotivation",
		"CPU_AI_MaxPitchCount",
		"CPU_AI_SpendScoutingWeeklyMaxRecruitsToScout",
		"EventRankingCutoff",
		"HardCommitAdditionalCommitScorePercentage",
		"MaxCollegeYearsRemaining",
		"MaxFiveStarsInClass",
		"MaxJumpBacks",
		"MaxRecruitingBoardTargets",
		"MaxRecruitVisitsPerWeek",
		"MaxTeamScholarshipOffers",
		"MaxTimesScouted",
		"MaxTotalHoursOnRecruitPerWeek",
		"OVRProjection",
		"PointsPercentageThresholdCommit",
		"PointsPercentageThresholdTop3",
		"PointsPercentageThresholdTop5",
		"RecruitStageBattleSchoolCount",
		"RecruitStageHardCommittedSchoolCount",
		"RecruitStageSignedSchoolCount",
		"RecruitStageSoftCommittedSchoolCount",
		"RecruitStageTop10SchoolCount",
		"RecruitStageTop3SchoolCount",
		"RecruitStageTop5SchoolCount",
		"RecruitsToGenerateGroupMax",
		"RecruitsToGenerateMaximum",
		"RecruitsToGenerateMinimum",
		"ScholarshipBonusLowExpectationCutoff",
		"ScholarshipOfferWeeklyBonus",
		"ScoutingAttributeLongSliderWidth",
		"ScoutingAttributeShortSliderWidth",
		"SoftCommitInfluenceVariance",
		"Stage2WarningThreshold",
		"TransferDefaultCommitScore",
		"TriggerRecruitingBattleWithinCommitThreshold",
	}
	for _, name := range scalarInts {
		if v, ok := intFieldOK(record, name); ok {
			out.Scalars[name] = v
		}
	}
	if value, ok := record.Get("MinimumNILToOfferPercentage"); ok && value.Float != 0 {
		out.Scalars["MinimumNILToOfferPercentage"] = value.Float
	}
	for _, name := range []string{"AIAggressivenessSlider", "AINILAggressivenessSlider"} {
		if value, ok := record.Get(name); ok && value.Float != 0 {
			out.Scalars[name] = value.Float
		}
	}
	arrayFields := []string{
		"ChanceToSwayBoostFromPipelineInfluenceLevel",
		"InfluenceRequiredPerPipelineLevel",
		"InitialInterest_AlmaMater_Starlevel",
		"InitialInterest_Pipeline",
		"InstantCommitBonusPrestige",
		"InstantCommitOddsPerStarLevel",
		"InterestBoost_Start_Awards_Starlevel",
		"RecruitingPointsAbstractionTable",
		"RecruitInitialInfluencePerPipelineLevel",
		"ScoutingAttributesUnlockPercentage",
		"ScoutingPhysicalAbilityUnlockPercentage",
		"TopClassesRankWeightPercentageTable",
		"TopClassesStarRatingTable",
		"TransferValueCutOffs",
	}
	for _, name := range arrayFields {
		if arr := intArrayFromRecord(f, record, name); len(arr) > 0 {
			out.Arrays[name] = arr
		}
	}
	for _, name := range []string{
		"EarlySigningDayActionItemString",
		"PreseasonRecruitingActionItemString",
		"TransferPortalActionItemString",
	} {
		if s := stringField(record, name); s != "" {
			out.Strings[name] = s
		}
	}
	return out
}

func recruitingPitchTunables(tf *File) []RecruitingPitchTunableExport {
	table, ok := tf.PrimaryTableByName("RecruitingPitchInfo")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil
	}
	out := make([]RecruitingPitchTunableExport, 0, len(table.Records))
	for _, record := range table.Records {
		entry := RecruitingPitchTunableExport{
			Pitch:                    stringField(record, "Pitch"),
			ShortName:                stringField(record, "ShortName"),
			LongName:                 stringField(record, "LongName"),
			HasDealbreakerMotivation: boolField(record, "HasDealbreakerMotivation"),
			AssociatedMotivation1:    stringField(record, "AssociatedMotivation1"),
			AssociatedMotivation2:    stringField(record, "AssociatedMotivation2"),
			AssociatedMotivation3:    stringField(record, "AssociatedMotivation3"),
		}
		setOptionalNonNegativeInt(record, "SwappableImage", &entry.SwappableImage)
		if entry.Pitch == "" && entry.ShortName == "" {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func visitTunables(tf *File) *VisitTunablesExport {
	table, ok := tf.PrimaryTableByName("VisitTunables")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil || len(table.Records) == 0 {
		return nil
	}
	record := table.Records[0]
	out := &VisitTunablesExport{}
	setOptionalNonNegativeInt(record, "CompetitiveVisitPenalty", &out.CompetitiveVisitPenalty)
	setOptionalNonNegativeInt(record, "ComplimentaryVisitBonus", &out.ComplimentaryVisitBonus)
	if out.CompetitiveVisitPenalty == nil && out.ComplimentaryVisitBonus == nil {
		return nil
	}
	return out
}

func visitActivityTunables(tf *File) []VisitActivityTunableExport {
	table, ok := tf.PrimaryTableByName("VisitActivityInfo")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil
	}
	out := make([]VisitActivityTunableExport, 0, len(table.Records))
	for _, record := range table.Records {
		entry := VisitActivityTunableExport{
			ActivityType: stringField(record, "ActivityType"),
		}
		setOptionalNonNegativeInt(record, "IconId", &entry.IconID)
		if entry.ActivityType == "" {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func recruitingActionTunables(tf *File) []RecruitingActionTunableExport {
	table, ok := tf.PrimaryTableByName("RecruitingActionInfo")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil
	}
	out := make([]RecruitingActionTunableExport, 0, len(table.Records))
	for _, record := range table.Records {
		entry := RecruitingActionTunableExport{
			ActionType: normalizeEnum(stringField(record, "ActionType")),
			Intensity:  normalizeEnum(stringField(record, "Intensity")),
			IsEnabled:  boolField(record, "IsEnabled"),
		}
		setOptionalNonNegativeInt(record, "Cost", &entry.Cost)
		setOptionalNonNegativeInt(record, "BaseInfluenceGranted", &entry.BaseInfluenceGranted)
		setOptionalNonNegativeInt(record, "IconId", &entry.IconID)
		if entry.ActionType == "" || entry.ActionType == "Invalid" {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func highSchoolRecruitingTunables(tf *File) *HighSchoolRecruitingTunablesExport {
	table, ok := tf.PrimaryTableByName("HighSchoolRecruitingTunables")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil || len(table.Records) == 0 {
		return nil
	}
	record := table.Records[0]
	scalars := make(map[string]any)
	for name, value := range record.Fields {
		if value.Float != 0 {
			scalars[name] = value.Float
			continue
		}
		if v, ok := tunableIntFromField(value); ok && v != 0 {
			scalars[name] = v
		}
	}
	if len(scalars) == 0 {
		return nil
	}
	return &HighSchoolRecruitingTunablesExport{Scalars: scalars}
}

func recruitingStageTunables(tf *File) []RecruitingStageTunableExport {
	table, ok := tf.PrimaryTableByName("RecruitingStageDetails")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil
	}
	out := make([]RecruitingStageTunableExport, 0, len(table.Records))
	seen := make(map[string]struct{})
	for _, record := range table.Records {
		stage := stringField(record, "RecruitStage")
		if stage == "" {
			continue
		}
		if _, exists := seen[stage]; exists {
			continue
		}
		seen[stage] = struct{}{}
		out = append(out, RecruitingStageTunableExport{RecruitStage: stage})
	}
	return out
}

func intArrayFromRecord(f *File, record Record, name string) []int {
	value, ok := record.Get(name)
	if !ok || value.Reference == nil {
		return nil
	}
	arrTable, ok := f.GetTableByID(value.Reference.TableID)
	if !ok || arrTable == nil {
		return nil
	}
	if err := arrTable.ReadRecords(); err != nil {
		return nil
	}
	rowIdx := int(value.Reference.RowNumber)
	if rowIdx < 0 || rowIdx >= len(arrTable.Records) {
		return nil
	}
	row := arrTable.Records[rowIdx]
	var out []int
	for i := 0; ; i++ {
		fv, ok := row.Get(strconv.Itoa(i))
		if !ok {
			break
		}
		v, ok := tunableIntFromField(fv)
		if !ok {
			break
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// RecruitingTunablesJSON returns recruiting formula constants as indented JSON.
func RecruitingTunablesJSON(schemaDir, tuningPath string) ([]byte, error) {
	export, err := ExportRecruitingTunables(schemaDir, tuningPath)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(export, "", "  ")
}
