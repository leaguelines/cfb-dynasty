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
		if userTargetTable != nil && record.Index < len(userTargetTable.Records) {
			applyUserRecruitTargetFields(&export, userTargetTable.Records[record.Index])
		}

		if !hasRecruitingData(export) {
			continue
		}
		exports = append(exports, export)
	}
	return exports, nil
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

func applyUserRecruitTargetFields(export *RecruitingTargetExport, record Record) {
	export.IsFavorite = boolField(record, "IsFavorite")
	if export.ScholarshipStatus == "" {
		export.ScholarshipStatus = stringField(record, "ScholarshipStatus")
	}
	if export.SwayPitch == "" {
		export.SwayPitch = stringField(record, "SwayPitch")
	}
	if export.ScheduledVisit == nil {
		// UserRecruitTarget duplicates pursuit fields; only fill visit when target row lacked it.
	}
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
		export.CommittedWeekNumber != nil
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
