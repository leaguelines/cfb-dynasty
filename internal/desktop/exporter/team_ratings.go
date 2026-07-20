package exporter

import "github.com/leaguelines/cfb-dynasty/dynasty"

// TeamRatings is the in-game team OVR/OFF/DEF rating display payload.
type TeamRatings struct {
	Overall string
	Offense string
	Defense string
}

// HasAny reports whether any rating is available for display.
func (r TeamRatings) HasAny() bool {
	return r.Overall != "" || r.Offense != "" || r.Defense != ""
}

// TeamRatingsByTeamID indexes exported team ratings by team id.
func TeamRatingsByTeamID(e dynasty.Export) map[int]TeamRatings {
	out := make(map[int]TeamRatings, len(e.Teams))
	for _, t := range e.Teams {
		if ratings := TeamRatingsFromExport(t); ratings.HasAny() {
			out[t.ID] = ratings
		}
	}
	return out
}

// TeamRatingsForTeamID returns exported ratings for a team id.
func TeamRatingsForTeamID(e dynasty.Export, teamID int) TeamRatings {
	for _, t := range e.Teams {
		if t.ID == teamID {
			return TeamRatingsFromExport(t)
		}
	}
	return TeamRatings{}
}

// TeamRatingsFromExport maps a team export's rating fields into display strings.
func TeamRatingsFromExport(t dynasty.TeamExport) TeamRatings {
	return TeamRatings{
		Overall: fmtIntPtr(t.OverallRating),
		Offense: fmtIntPtr(t.OffensiveRating),
		Defense: fmtIntPtr(t.DefensiveRating),
	}
}
