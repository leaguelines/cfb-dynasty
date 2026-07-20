// Package exporter defines the canonical set of dynasty collections and how to
// render each one as JSON, CSV, or an HTML table. The registry is the single
// source of truth shared by the UI and the download endpoints.
package exporter

import "github.com/leaguelines/cfb-dynasty/dynasty"

// Collection describes one exportable dataset within an Export.
type Collection struct {
	// Name is the URL-safe identifier (e.g. "teams").
	Name string
	// Title is the human-readable label.
	Title string
	// Count returns how many rows the collection has for a given export.
	Count func(dynasty.Export) int
	// JSON returns the value to marshal for this section alone.
	JSON func(dynasty.Export) any
	// Header is the CSV / table column header.
	Header []string
	// Rows returns the flattened CSV / table rows.
	Rows func(dynasty.Export) [][]string
	// LinkPrefix, when non-empty, turns the cell at LinkColumn into a link to
	// LinkPrefix+cellValue in the HTML table (e.g. link games to detail pages).
	LinkPrefix string
	// LinkColumn is the row index made into a link when LinkPrefix is set.
	LinkColumn int
}

// playerValues holds a nil-safe, pre-formatted view of a player's attributes so
// collections can order columns however is most useful.
type playerValues struct {
	ID                 string
	FirstName          string
	LastName           string
	Position           string
	IsAth              string
	Overall            string
	StarRating         string
	SchoolYear         string
	Archetype          string
	ArchetypeLabel     string
	Age                string
	Height             string
	Weight             string
	Jersey             string
	HomeState          string
	HomeTown           string
	SkillGroupCapTotal string
	TeamIndex          string
}

func playerVals(p *dynasty.PlayerExport) playerValues {
	if p == nil {
		return playerValues{}
	}
	return playerValues{
		ID:                 fmtInt(p.ID),
		FirstName:          p.FirstName,
		LastName:           p.LastName,
		Position:           p.Position,
		IsAth:              fmtBool(p.IsAth),
		Overall:            fmtIntPtr(p.Overall),
		StarRating:         p.StarRating,
		SchoolYear:         p.SchoolYear,
		Archetype:          p.Archetype,
		ArchetypeLabel:     p.ArchetypeLabel,
		Age:                fmtIntPtr(p.Age),
		Height:             fmtIntPtr(p.Height),
		Weight:             fmtIntPtr(p.Weight),
		Jersey:             fmtIntPtr(p.Jersey),
		HomeState:          p.HomeState,
		HomeTown:           p.HomeTown,
		SkillGroupCapTotal: fmtInt(p.SkillGroupCapTotal),
		TeamIndex:          fmtIntPtr(p.TeamIndex),
	}
}

// PlayerValues returns formatted display strings for a player.
func PlayerValues(p *dynasty.PlayerExport) playerValues {
	return playerVals(p)
}

// Registry is the ordered list of all collections.
var Registry = []Collection{
	seasonCollection,
	teamsCollection,
	rostersCollection,
	gamesCollection,
	recruitsCollection,
	recruitingCollection,
	seasonPlayerStatsCollection,
	seasonTeamStatsCollection,
	coachesCollection,
	leavingPlayersCollection,
	injuriesCollection,
	depthChartsCollection,
	playerAwardsCollection,
	leagueAwardsCollection,
	conferenceChampionsCollection,
	recordBookCollection,
	schoolGradesCollection,
	pipelineInfluenceCollection,
	rivalriesCollection,
	positionChangesCollection,
	draftPicksCollection,
	bowlGamesCollection,
}

var byName = func() map[string]Collection {
	m := make(map[string]Collection, len(Registry))
	for _, c := range Registry {
		m[c.Name] = c
	}
	return m
}()

// Lookup returns the collection with the given name.
func Lookup(name string) (Collection, bool) {
	c, ok := byName[name]
	return c, ok
}
