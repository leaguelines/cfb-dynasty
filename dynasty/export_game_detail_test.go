package dynasty

import (
	"path/filepath"
	"testing"
)

func TestExportGames_DetailFields(t *testing.T) {
	savePath := filepath.Join("..", "data", "DYNASTY-2026OFFLINEFINAL")
	settings := DefaultSettings()
	settings.SchemaDir = filepath.Join("..", "data")
	settings.AutoParse = true
	file, err := Open(savePath, &settings)
	if err != nil {
		t.Skip("offline save not available:", savePath)
	}
	if err := file.Parse(); err != nil {
		t.Fatal(err)
	}

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
