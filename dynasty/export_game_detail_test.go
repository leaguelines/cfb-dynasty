package dynasty

import "testing"

func TestExportGames_DetailFields(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)

	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Games: true}})
	if err != nil {
		t.Fatal(err)
	}

	withQuarters, withWeather, withKicking := 0, 0, 0
	for _, game := range export.Games {
		if len(game.HomeQuarterScores) > 0 {
			withQuarters++
		}
		if game.Weather != "" {
			withWeather++
		}
		if len(game.KickingStats) > 0 {
			withKicking++
		}
	}
	if withQuarters == 0 {
		t.Fatal("expected games with quarter scores")
	}
	if withWeather == 0 {
		t.Fatal("expected games with weather")
	}
	if withKicking == 0 {
		t.Fatal("expected games with kicking stats")
	}
}
