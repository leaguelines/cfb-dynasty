package dynasty

import "testing"

func TestRecruitSchoolInterests(t *testing.T) {
	file := openTestSave(t)
	recruitTable, ok := file.PrimaryTableByName("Recruit")
	if !ok {
		t.Fatal("Recruit table missing")
	}
	if err := recruitTable.ReadRecords(); err != nil {
		t.Fatal(err)
	}
	schoolTable, ok := file.PrimaryTableByName("ProspectTargetSchool")
	if !ok {
		t.Fatal("ProspectTargetSchool table missing")
	}
	if err := schoolTable.ReadRecords(); err != nil {
		t.Fatal(err)
	}
	teams := file.teamMaps()

	withInterest := 0
	for _, record := range recruitTable.Records {
		interests := recruitSchoolInterests(file, record, schoolTable, teams)
		if len(interests) == 0 {
			continue
		}
		withInterest++
		for i := 1; i < len(interests); i++ {
			if interests[i].Influence > interests[i-1].Influence {
				t.Fatalf("interests not sorted by influence: %#v", interests)
			}
		}
		if withInterest >= 3 {
			break
		}
	}
	if withInterest == 0 {
		t.Fatal("expected recruits with school interest")
	}
}

func TestExportSchoolGradesAndPipelines(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{SchoolGrades: true, Pipelines: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(export.SchoolGrades) == 0 {
		t.Fatal("expected school grades")
	}
	if export.SchoolGrades[0].AcademicPrestige == "" && export.SchoolGrades[0].AthleticFacilities == "" {
		t.Fatalf("expected grade fields on first school: %#v", export.SchoolGrades[0])
	}
	if len(export.PipelineInfluence) == 0 {
		t.Fatal("expected pipeline influence")
	}
	if export.PipelineInfluence[0].Pipeline == "" {
		t.Fatalf("expected pipeline name: %#v", export.PipelineInfluence[0])
	}
}

func TestExportRecruits_SchoolInterest(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Recruits: true}})
	if err != nil {
		t.Fatal(err)
	}
	withInterest := 0
	for _, recruit := range export.Recruits {
		if len(recruit.SchoolInterest) > 0 {
			withInterest++
		}
	}
	if withInterest == 0 {
		t.Fatal("expected recruits with schoolInterest")
	}
}
