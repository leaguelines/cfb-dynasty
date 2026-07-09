package dynasty

import (
	"fmt"
	"sort"
)

func gradeField(record Record, name string) string {
	return normalizeEnum(stringField(record, name))
}

func (f *File) buildSchoolGradesExports() ([]SchoolGradesExport, error) {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil, nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()
	exports := make([]SchoolGradesExport, 0, teamTable.ActiveRecordCount())
	for _, team := range teamTable.Records {
		teamID, ok := teams.exportID(team.Index)
		if !ok {
			continue
		}
		ref, ok := team.Get("MySchoolTrackingTable")
		if !ok || ref.Reference == nil {
			continue
		}
		gradesRecord, ok := f.RecordByReference("MySchoolTrackingTable", ref.Reference)
		if !ok {
			continue
		}

		export := SchoolGradesExport{
			TeamID:               teamID,
			TeamName:             teams.nameFromID(teamID),
			AcademicPrestige:     gradeField(gradesRecord, "AcademicPrestigeGrade"),
			AthleticFacilities:   gradeField(gradesRecord, "AthleticFacilitiesGrade"),
			BrandExposure:        gradeField(gradesRecord, "BrandExposureGrade"),
			CampusLifestyle:      gradeField(gradesRecord, "CampusLifestyleGrade"),
			ChampionshipContender: gradeField(gradesRecord, "ChampionshipContenderGrade"),
			CoachPrestige:        gradeField(gradesRecord, "CoachPrestigeGrade"),
			CoachStability:       gradeField(gradesRecord, "CoachStabilityGrade"),
			ConferencePrestige:   gradeField(gradesRecord, "ConferencePrestigeGrade"),
			ProgramTradition:     gradeField(gradesRecord, "ProgramTraditionGrade"),
			StadiumAtmosphere:    gradeField(gradesRecord, "StadiumAtmosphereGrade"),
			ProPotentialQB:       gradeField(gradesRecord, "ProPotentialGradeQB"),
			ProPotentialRB:       gradeField(gradesRecord, "ProPotentialGradeRB"),
			ProPotentialWR:       gradeField(gradesRecord, "ProPotentialGradeWR"),
			ProPotentialTE:       gradeField(gradesRecord, "ProPotentialGradeTE"),
			ProPotentialOL:       gradeField(gradesRecord, "ProPotentialGradeOL"),
			ProPotentialDL:       gradeField(gradesRecord, "ProPotentialGradeDL"),
			ProPotentialLB:       gradeField(gradesRecord, "ProPotentialGradeLB"),
			ProPotentialDB:       gradeField(gradesRecord, "ProPotentialGradeDB"),
			ProPotentialK:        gradeField(gradesRecord, "ProPotentialGradeK"),
			ProPotentialP:        gradeField(gradesRecord, "ProPotentialGradeP"),
		}
		if !hasSchoolGrades(export) {
			continue
		}
		exports = append(exports, export)
	}

	sort.Slice(exports, func(i, j int) bool {
		return exports[i].TeamID < exports[j].TeamID
	})
	return exports, nil
}

func hasSchoolGrades(export SchoolGradesExport) bool {
	return export.AcademicPrestige != "" || export.AthleticFacilities != "" ||
		export.BrandExposure != "" || export.CampusLifestyle != "" ||
		export.ChampionshipContender != "" || export.CoachPrestige != "" ||
		export.ConferencePrestige != "" || export.ProgramTradition != ""
}

func (f *File) buildPipelineInfluenceExports() ([]PipelineInfluenceExport, error) {
	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil, nil
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil, err
	}

	pipelineTable, ok := f.PrimaryTableByName("SchoolPipelineInfluence")
	if !ok {
		return nil, nil
	}
	if err := pipelineTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()
	pipelineID := pipelineTable.Header.TableID
	exports := make([]PipelineInfluenceExport, 0, teamTable.ActiveRecordCount()*8)

	for _, team := range teamTable.Records {
		teamID, ok := teams.exportID(team.Index)
		if !ok {
			continue
		}
		ref, ok := team.Get("SchoolPipelineInfluenceList")
		if !ok || isEmptyReference(team, "SchoolPipelineInfluenceList") {
			continue
		}
		seen := make(map[string]struct{})
		for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
			if memberRef.TableID != 0 && memberRef.TableID != pipelineID {
				continue
			}
			row, ok := f.RecordByReference("SchoolPipelineInfluence", memberRef)
			if !ok {
				continue
			}
			pipeline := normalizeEnum(stringField(row, "Pipeline"))
			if pipeline == "" || pipeline == "Invalid" || pipeline == "None" {
				continue
			}
			level := normalizeEnum(stringField(row, "InfluenceLevel"))
			if level == "" || level == "None" || level == "Invalid" {
				continue
			}
			key := fmt.Sprintf("%d:%s", teamID, pipeline)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			export := PipelineInfluenceExport{
				TeamID:         teamID,
				TeamName:       teams.nameFromID(teamID),
				Pipeline:       pipeline,
				InfluenceLevel: level,
			}
			setOptionalPositiveInt(row, "InfluenceValue", &export.InfluenceValue)
			exports = append(exports, export)
		}
	}

	sort.Slice(exports, func(i, j int) bool {
		if exports[i].TeamID != exports[j].TeamID {
			return exports[i].TeamID < exports[j].TeamID
		}
		return exports[i].Pipeline < exports[j].Pipeline
	})
	return exports, nil
}
