package exporter

import (
	"fmt"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// RosterGroupingError reports when parsed rosters look collapsed into one or a
// few teams — usually a schema bundle mismatch (wrong TeamIndex field mapping).
func RosterGroupingError(e *dynasty.Export) string {
	teams := len(e.Teams)
	rosters := len(e.Rosters)
	if teams < 10 || rosters == 0 {
		return ""
	}

	totalPlayers := 0
	maxRoster := 0
	for _, r := range e.Rosters {
		n := len(r.Players)
		totalPlayers += n
		if n > maxRoster {
			maxRoster = n
		}
	}
	if totalPlayers == 0 {
		return ""
	}

	// Healthy dynasty saves spread players across most FBS teams (~85 per roster).
	if rosters >= teams/3 {
		return ""
	}

	// One mega-roster holding most of the league is the classic bad-schema symptom.
	if rosters <= 2 && maxRoster > 500 && maxRoster*100/totalPlayers > 80 {
		return fmt.Sprintf(
			"team rosters could not be split correctly (%d teams but only %d roster groups with %d players lumped together); "+
				"the schema bundle in SCHEMA_DIR likely does not match this save — add a current C27_*.gz franchise schema matching the game patch and try again",
			teams, rosters, maxRoster,
		)
	}

	if rosters < teams/10 {
		return fmt.Sprintf(
			"only %d of %d teams have roster data; check that SCHEMA_DIR contains a schema bundle matching this save",
			rosters, teams,
		)
	}

	return ""
}
