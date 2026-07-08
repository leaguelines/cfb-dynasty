package dynasty

// playerCareerStatsFromRecord builds career playing-time summary from a CareerStats row.
func playerCareerStatsFromRecord(record Record) *PlayerCareerStatsExport {
	stats := &PlayerCareerStatsExport{}
	hasData := false
	set := func(dst **int, name string) {
		if v, ok := careerStatIntOK(record, name); ok {
			*dst = &v
			hasData = true
		}
	}

	set(&stats.GamesPlayed, "GAMESPLAYED")
	set(&stats.GamesStarted, "GAMESSTARTED")
	set(&stats.DownsPlayed, "DOWNSPLAYED")
	set(&stats.GameRating, "GAMERATING")

	if !hasData {
		return nil
	}
	return stats
}

func (f *File) playerCareerStatsExport(record Record) *PlayerCareerStatsExport {
	ref, ok := record.Get("CareerStats")
	if !ok || ref.Reference == nil {
		return nil
	}
	if ref.Reference.TableID == 0 && ref.Reference.RowNumber == 0 {
		return nil
	}
	careerRecord, ok := f.RecordByReference("CareerStats", ref.Reference)
	if !ok {
		return nil
	}
	return playerCareerStatsFromRecord(careerRecord)
}
