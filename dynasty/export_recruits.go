package dynasty

// buildRecruitExports decodes all Recruit rows and their linked Player records.
func (f *File) buildRecruitExports() ([]RecruitExport, error) {
	recruitTable, ok := f.PrimaryTableByName("Recruit")
	if !ok {
		return nil, nil
	}
	if err := recruitTable.ReadRecords(); err != nil {
		return nil, err
	}

	exports := make([]RecruitExport, 0, recruitTable.ActiveRecordCount())
	active := int(recruitTable.ActiveRecordCount())
	teams := f.teamMaps()
	for _, record := range recruitTable.Records {
		if record.Index >= active {
			continue
		}
		export := RecruitExport{ID: record.Index}
		export.Class = stringField(record, "Class")
		export.RecruitStage = stringField(record, "RecruitStage")
		export.QualityModifier = stringField(record, "QualityModifier")
		export.AlternatePosition1 = stringField(record, "AlternatePosition1")
		export.AlternatePosition2 = stringField(record, "AlternatePosition2")

		setOptionalPositiveInt(record, "NationalRank", &export.NationalRank)
		setOptionalPositiveInt(record, "StateRank", &export.StateRank)
		setOptionalPositiveInt(record, "PositionRank", &export.PositionRank)
		setOptionalPositiveInt(record, "CommitScore", &export.CommitScore)
		setOptionalPositiveInt(record, "ProductionGrade", &export.ProductionGrade)
		setOptionalPositiveInt(record, "TotalScholarshipOffers", &export.TotalScholarshipOffers)

		if playerRef, ok := record.Get("Player"); ok && playerRef.Reference != nil {
			if playerRecord, _, ok := f.playerRecordByReference(playerRef.Reference); ok {
				export.Player = buildPlayerExport(playerRecord)
				applyCanonicalTeamIndex(export.Player, teams)
				applyPlayerAth(export.Player, export.AlternatePosition1, export.AlternatePosition2)
			}
		}

		exports = append(exports, export)
	}
	return exports, nil
}

func setOptionalInt(record Record, name string, dst **int) {
	if v, ok := intFieldOK(record, name); ok {
		*dst = &v
	}
}
