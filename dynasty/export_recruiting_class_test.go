package dynasty

import "testing"

func TestRecruitingClassScoreForTeam(t *testing.T) {
	weights := []int{100, 99, 98}
	recruit := Record{Fields: map[string]FieldValue{
		"CommitScore":  {Int: 400},
		"NationalRank": {Int: 1},
	}}
	if got := topClassRankWeight(weights, recruit, 3); got != 100 {
		t.Fatalf("weight = %d, want 100", got)
	}
	recruit.Fields["NationalRank"] = FieldValue{Int: 0}
	if got := topClassRankWeight(weights, recruit, 3); got != 3 {
		t.Fatalf("fallback weight = %d, want 3", got)
	}
}

func TestRecruitingClassExportSeasonSave(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Teams: true}})
	if err != nil {
		t.Fatal(err)
	}

	var georgia, ohioState, alabama *TeamExport
	for i := range export.Teams {
		team := &export.Teams[i]
		switch team.LongName {
		case "Georgia":
			georgia = team
		case "Ohio State":
			ohioState = team
		case "Alabama":
			alabama = team
		}
	}
	if georgia == nil || georgia.RecruitingClass == nil {
		t.Fatal("Georgia recruiting class missing")
	}
	if georgia.RecruitingClass.NationalRank == nil || *georgia.RecruitingClass.NationalRank != 1 {
		t.Fatalf("Georgia national rank = %v, want 1", georgia.RecruitingClass.NationalRank)
	}
	if georgia.RecruitingClass.Score != 2394 {
		t.Fatalf("Georgia class score = %d, want 2394", georgia.RecruitingClass.Score)
	}
	if georgia.RecruitingClass.CommitCount != 31 {
		t.Fatalf("Georgia commit count = %d, want 31", georgia.RecruitingClass.CommitCount)
	}
	if ohioState == nil || ohioState.RecruitingClass == nil || ohioState.RecruitingClass.Score != 2037 {
		t.Fatalf("Ohio State class score = %v, want 2037", safeClassScore(ohioState))
	}
	if alabama == nil || alabama.RecruitingClass == nil || alabama.RecruitingClass.Score != 442 {
		t.Fatalf("Alabama class score = %v, want 442", safeClassScore(alabama))
	}
}

func safeClassScore(team *TeamExport) int {
	if team == nil || team.RecruitingClass == nil {
		return -1
	}
	return team.RecruitingClass.Score
}
