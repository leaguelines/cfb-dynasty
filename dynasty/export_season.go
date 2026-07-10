package dynasty

func buildSeasonPeriodsExport(record Record) *SeasonPeriodsExport {
	periods := &SeasonPeriodsExport{
		IsRecruitingPeriodActive:       boolField(record, "IsRecruitingPeriodActive"),
		IsSigningPeriodActive:          boolField(record, "IsSigningPeriodActive"),
		IsVisitingPeriodActive:         boolField(record, "IsVisitingPeriodActive"),
		IsPitchingPeriodActive:         boolField(record, "IsPitchingPeriodActive"),
		IsScholarshipPeriodActive:      boolField(record, "IsScholarshipPeriodActive"),
		IsScoutingPeriodActive:         boolField(record, "IsScoutingPeriodActive"),
		IsTransferPortalNewlyAvailable: boolField(record, "IsTransferPortalNewlyAvailable"),
		IsTransferSignPeriodActive:     boolField(record, "IsTransferSignPeriodActive"),
		IsDraftPeriodActive:            boolField(record, "IsDraftPeriodActive"),
		IsDraftScoutingActive:          boolField(record, "IsDraftScoutingActive"),
		IsGoalsPeriodActive:            boolField(record, "IsGoalsPeriodActive"),
		IsCarouselPeriodActive:         boolField(record, "IsCarouselPeriodActive"),
		IsStaffHiringPeriodActive:      boolField(record, "IsStaffHiringPeriodActive"),
		IsWeeklyAwardPeriodActive:      boolField(record, "IsWeeklyAwardPeriodActive"),
		IsAnnualAwardPeriodActive:      boolField(record, "IsAnnualAwardPeriodActive"),
	}
	setOptionalPositiveInt(record, "RegularSeasonLastWeekScheduled", &periods.RegularSeasonLastWeekScheduled)
	setOptionalPositiveInt(record, "PostSeasonNumWeeks", &periods.PostSeasonNumWeeks)
	if !seasonPeriodsHasData(*periods) {
		return nil
	}
	return periods
}

func seasonPeriodsHasData(p SeasonPeriodsExport) bool {
	if p.IsRecruitingPeriodActive || p.IsSigningPeriodActive || p.IsVisitingPeriodActive ||
		p.IsPitchingPeriodActive || p.IsScholarshipPeriodActive || p.IsScoutingPeriodActive ||
		p.IsTransferPortalNewlyAvailable || p.IsTransferSignPeriodActive || p.IsDraftPeriodActive ||
		p.IsDraftScoutingActive || p.IsGoalsPeriodActive || p.IsCarouselPeriodActive ||
		p.IsStaffHiringPeriodActive || p.IsWeeklyAwardPeriodActive || p.IsAnnualAwardPeriodActive {
		return true
	}
	return p.RegularSeasonLastWeekScheduled != nil || p.PostSeasonNumWeeks != nil
}
