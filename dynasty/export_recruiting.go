package dynasty

// buildRecruitingExports assembles per-recruit pursuit state and school interest.
func (f *File) buildRecruitingExports() ([]RecruitingTargetExport, error) {
	recruitTable, ok := f.PrimaryTableByName("Recruit")
	if !ok {
		return nil, nil
	}
	if err := recruitTable.ReadRecords(); err != nil {
		return nil, err
	}

	targetTable, _ := f.PrimaryTableByName("RecruitTarget")
	if targetTable != nil {
		_ = targetTable.ReadRecords()
	}
	userTargetTable, _ := f.PrimaryTableByName("UserRecruitTarget")
	if userTargetTable != nil {
		_ = userTargetTable.ReadRecords()
	}

	schoolTable, _ := f.PrimaryTableByName("ProspectTargetSchool")
	if schoolTable != nil {
		_ = schoolTable.ReadRecords()
	}
	teams := f.teamMaps()

	exports := make([]RecruitingTargetExport, 0, recruitTable.ActiveRecordCount())
	active := int(recruitTable.ActiveRecordCount())
	for _, record := range recruitTable.Records {
		if record.Index >= active {
			continue
		}
		if stringField(record, "RecruitStage") == "Invalid" {
			continue
		}

		export := RecruitingTargetExport{RecruitID: record.Index}
		interests := recruitSchoolInterests(f, record, schoolTable, teams)
		export.SchoolInterest = interests
		export.TopSchool = topSchoolFromInterests(interests)

		if targetTable != nil && record.Index < len(targetTable.Records) {
			targetRow := targetTable.Records[record.Index]
			applyRecruitTargetFields(&export, targetRow, f)
			export.ActivePitches = activePitchesFromTarget(f, targetRow)
			setOptionalNonNegativeInt(targetRow, "UnlockedIntelBitfield", &export.UnlockedIntelBitfield)
		}
		if userTargetTable != nil {
			if userRow, ok := userRecruitTargetForRecruit(userTargetTable, record.Index); ok {
				applyUserRecruitTargetFields(f, &export, userRow)
			}
		}

		if !hasRecruitingData(export) {
			continue
		}
		exports = append(exports, export)
	}
	return exports, nil
}

// userRecruitTargetForRecruit finds the UserRecruitTarget whose Recruit ref
// points at recruitIndex. URT capacity is much smaller than Recruit and is not
// 1:1 by row index.
func userRecruitTargetForRecruit(userTargetTable *Table, recruitIndex int) (Record, bool) {
	if userTargetTable == nil {
		return Record{}, false
	}
	for _, row := range userTargetTable.Records {
		ref, ok := row.Get("Recruit")
		if !ok || ref.Reference == nil {
			continue
		}
		if int(ref.Reference.RowNumber) == recruitIndex {
			return row, true
		}
	}
	// Legacy fallback: some older saves aligned URT index with Recruit index.
	if recruitIndex >= 0 && recruitIndex < len(userTargetTable.Records) {
		row := userTargetTable.Records[recruitIndex]
		if ref, ok := row.Get("Recruit"); ok && ref.Reference != nil {
			if int(ref.Reference.RowNumber) == recruitIndex {
				return row, true
			}
			// Indexed row belongs to a different recruit; do not misuse it.
			return Record{}, false
		}
	}
	return Record{}, false
}

func applyRecruitTargetFields(export *RecruitingTargetExport, record Record, f *File) {
	if status := stringField(record, "ScholarshipStatus"); status != "" && status != "None" {
		export.ScholarshipStatus = status
	}
	export.SwayPitch = stringField(record, "SwayPitch")
	export.ContactFriendsAndFamily = boolField(record, "ContactFriendsAndFamily")
	export.ContactHighSchoolCoaches = boolField(record, "ContactHighSchoolCoaches")
	export.SearchSocialMedia = boolField(record, "SearchSocialMedia")
	export.SendTheHouse = boolField(record, "SendTheHouse")
	export.VisitRecruitsSchool = boolField(record, "VisitRecruitsSchool")

	setOptionalPositiveInt(record, "CurrentNILOffer", &export.CurrentNILOffer)
	setOptionalPositiveInt(record, "NILExpectation", &export.NILExpectation)
	setOptionalPositiveInt(record, "OriginalNILExpectation", &export.OriginalNILExpectation)
	setOptionalRecruitingInt(record, "CurrentScholarshipBonus", &export.CurrentScholarshipBonus)
	setOptionalPositiveInt(record, "ProspectInfluenceTotal", &export.ProspectInfluenceTotal)
	setOptionalRecruitingInt(record, "ProspectInfluenceDelta", &export.ProspectInfluenceDelta)
	setOptionalPositiveInt(record, "ProspectInfluenceTotalLastWeek", &export.ProspectInfluenceTotalLastWeek)
	setOptionalPositiveInt(record, "ProspectHoursSpentCurrent", &export.ProspectHoursSpentCurrent)
	setOptionalPositiveInt(record, "CommittedWeekNumber", &export.CommittedWeekNumber)

	if visit, ok := record.Get("ScheduledVisit"); ok && visit.Reference != nil {
		if visitRecord, ok := f.RecordByReference("ActiveVisitInfo", visit.Reference); ok {
			export.ScheduledVisit = buildRecruitingVisitExport(visitRecord)
		}
	}
}

func applyUserRecruitTargetFields(f *File, export *RecruitingTargetExport, record Record) {
	export.IsFavorite = boolField(record, "IsFavorite")
	if status := stringField(record, "ScholarshipStatus"); status != "" && status != "None" {
		export.ScholarshipStatus = status
	}
	if sway := stringField(record, "SwayPitch"); sway != "" {
		export.SwayPitch = sway
	}

	// User-board influence/hours/NIL are the values that match the UI and
	// ProspectTargetSchool deltas; prefer them over RecruitTarget copies.
	setOptionalPositiveInt(record, "CurrentNILOffer", &export.CurrentNILOffer)
	setOptionalPositiveInt(record, "NILExpectation", &export.NILExpectation)
	setOptionalPositiveInt(record, "OriginalNILExpectation", &export.OriginalNILExpectation)
	setOptionalRecruitingInt(record, "CurrentScholarshipBonus", &export.CurrentScholarshipBonus)
	setOptionalPositiveInt(record, "ProspectInfluenceTotal", &export.ProspectInfluenceTotal)
	setOptionalRecruitingInt(record, "ProspectInfluenceDelta", &export.ProspectInfluenceDelta)
	setOptionalPositiveInt(record, "ProspectInfluenceTotalLastWeek", &export.ProspectInfluenceTotalLastWeek)
	setOptionalPositiveInt(record, "ProspectHoursSpentCurrent", &export.ProspectHoursSpentCurrent)
	setOptionalPositiveInt(record, "CommittedWeekNumber", &export.CommittedWeekNumber)

	if len(export.ActivePitches) == 0 {
		export.ActivePitches = activePitchesFromTarget(f, record)
	}
	if export.ScheduledVisit == nil {
		if visit, ok := record.Get("ScheduledVisit"); ok && visit.Reference != nil {
			if visitRecord, ok := f.RecordByReference("ActiveVisitInfo", visit.Reference); ok {
				export.ScheduledVisit = buildRecruitingVisitExport(visitRecord)
			}
		}
	}
	export.RecruitingFeedback = recruitingFeedbackFromTarget(f, record, "RecruitingFeedback")
	export.ImmediateRecruitingFeedback = recruitingFeedbackFromTarget(f, record, "ImmediateRecruitingFeedback")
}

func recruitingFeedbackFromTarget(f *File, record Record, field string) []RecruitingActionFeedbackExport {
	ref, ok := record.Get(field)
	if !ok || ref.Reference == nil {
		return nil
	}
	var out []RecruitingActionFeedbackExport
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("RecruitingActionFeedbackEntry", memberRef)
		if !ok {
			continue
		}
		entry := RecruitingActionFeedbackExport{
			ActionType: normalizeEnum(stringField(row, "RecruitingActionType")),
			Intensity:  normalizeEnum(stringField(row, "RecruitingActionIntensity")),
			Bonuses:    recruitingBonusesFromFeedback(f, row),
		}
		setOptionalNonNegativeInt(row, "HoursSpent", &entry.HoursSpent)
		setOptionalRecruitingInt(row, "InfluenceGained", &entry.InfluenceGained)
		setOptionalRecruitingInt(row, "MinInfluenceGain", &entry.MinInfluenceGain)
		setOptionalRecruitingInt(row, "MaxInfluenceGain", &entry.MaxInfluenceGain)
		setOptionalNonNegativeInt(row, "IntelUnlocked", &entry.IntelUnlocked)
		if entry.ActionType == "" && entry.Intensity == "" &&
			entry.InfluenceGained == nil && len(entry.Bonuses) == 0 {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func recruitingBonusesFromFeedback(f *File, feedback Record) []RecruitingActionBonusExport {
	ref, ok := feedback.Get("BonusList")
	if !ok || ref.Reference == nil {
		return nil
	}
	var out []RecruitingActionBonusExport
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("RecruitingActionBonus", memberRef)
		if !ok {
			continue
		}
		bonus := RecruitingActionBonusExport{
			BonusType:      normalizeEnum(stringField(row, "BonusType")),
			BonusValueType: normalizeEnum(stringField(row, "BonusValueType")),
		}
		if v, ok := intFieldOK(row, "BonusValue"); ok {
			bonus.BonusValue = v
		}
		if bonus.BonusType == "" && bonus.BonusValue == 0 {
			continue
		}
		out = append(out, bonus)
	}
	return out
}

func buildRecruitingVisitExport(record Record) *RecruitingVisitExport {
	visit := &RecruitingVisitExport{
		WeekType: stringField(record, "WeekType"),
		Activity: stringField(record, "Activity"),
	}
	if week, ok := intFieldOK(record, "WeekNumber"); ok && week > 0 {
		visit.Week = week
	}
	if visit.Week == 0 && visit.WeekType == "" && visit.Activity == "" {
		return nil
	}
	return visit
}

func hasRecruitingData(export RecruitingTargetExport) bool {
	if len(export.SchoolInterest) > 0 {
		return true
	}
	if export.TopSchool != nil {
		return true
	}
	if len(export.ActivePitches) > 0 {
		return true
	}
	if export.UnlockedIntelBitfield != nil {
		return true
	}
	if export.ScholarshipStatus != "" && export.ScholarshipStatus != "None" {
		return true
	}
	if export.SwayPitch != "" {
		return true
	}
	if export.ScheduledVisit != nil {
		return true
	}
	if export.IsFavorite {
		return true
	}
	if export.ContactFriendsAndFamily || export.ContactHighSchoolCoaches ||
		export.SearchSocialMedia || export.SendTheHouse || export.VisitRecruitsSchool {
		return true
	}
	return export.CurrentNILOffer != nil || export.NILExpectation != nil ||
		export.OriginalNILExpectation != nil || export.CurrentScholarshipBonus != nil ||
		export.ProspectInfluenceTotal != nil || export.ProspectInfluenceDelta != nil ||
		export.ProspectInfluenceTotalLastWeek != nil || export.ProspectHoursSpentCurrent != nil ||
		export.CommittedWeekNumber != nil ||
		len(export.RecruitingFeedback) > 0 || len(export.ImmediateRecruitingFeedback) > 0
}

func activePitchesFromTarget(f *File, record Record) []RecruitingPitchExport {
	ref, ok := record.Get("ActivePitches")
	if !ok || ref.Reference == nil {
		return nil
	}
	var pitches []RecruitingPitchExport
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("ActiveRecruitingPitch", memberRef)
		if !ok {
			continue
		}
		pitch := RecruitingPitchExport{
			Pitch:     normalizeEnum(stringField(row, "Pitch")),
			Intensity: normalizeEnum(stringField(row, "Intensity")),
		}
		if pitch.Pitch == "" && pitch.Intensity == "" {
			continue
		}
		pitches = append(pitches, pitch)
	}
	return pitches
}

func boolField(record Record, name string) bool {
	value, ok := record.Get(name)
	if !ok {
		return false
	}
	return value.Bool
}

func setOptionalPositiveInt(record Record, name string, dst **int) {
	v, ok := intFieldOK(record, name)
	if !ok || v <= 0 {
		return
	}
	*dst = &v
}

func setOptionalRecruitingInt(record Record, name string, dst **int) {
	v, ok := intFieldOK(record, name)
	if !ok || v <= -100 {
		return
	}
	*dst = &v
}
