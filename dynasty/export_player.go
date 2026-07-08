package dynasty

import "strings"

var playerIdentityFields = []string{
	"FirstName", "LastName", "Position", "PlayerType", "SchoolYear", "Age", "Height", "Weight",
	"OverallRating", "ProspectStarRating", "PLYR_HOME_TOWN", "PLYR_HOME_STATE",
	"JerseyNum", "TeamIndex",
}

func buildPlayerExport(record Record) *PlayerExport {
	if len(record.Fields) == 0 {
		return nil
	}
	player := &PlayerExport{ID: record.Index}
	for _, name := range playerIdentityFields {
		value, ok := record.Get(name)
		if !ok {
			continue
		}
		switch name {
		case "FirstName":
			player.FirstName = value.String
		case "LastName":
			player.LastName = value.String
		case "Position":
			player.Position = value.String
		case "PlayerType":
			player.Archetype = value.String
		case "SchoolYear":
			player.SchoolYear = value.String
		case "ProspectStarRating":
			player.StarRating = value.String
		case "PLYR_HOME_TOWN":
			player.HomeTown = value.String
		case "PLYR_HOME_STATE":
			player.HomeState = value.String
		case "Age":
			setIntPtr(&player.Age, value.Int)
		case "Height":
			setIntPtr(&player.Height, value.Int)
		case "Weight":
			if value.Int > 0 {
				setIntPtr(&player.Weight, value.Int)
			}
		case "OverallRating":
			setIntPtr(&player.Overall, value.Int)
		case "JerseyNum":
			setIntPtr(&player.Jersey, value.Int)
		case "TeamIndex":
			// Raw save value is a Team table row (including 0 for Air Force);
			// callers remap via teamMaps so consumers see the stable Team.TeamIndex.
			v := int(value.Int)
			player.TeamIndex = &v
		}
	}

	if label := archetypeLabelFromRecord(record); label != "" {
		player.ArchetypeLabel = label
	}

	if caps, total, ok := skillGroupCapsFromRecord(record); ok {
		player.SkillGroupCaps = caps
		player.SkillGroupCapTotal = total
	}

	ratings := make(map[string]int)
	for key, value := range record.Fields {
		if value.String != "" || value.Reference != nil {
			continue
		}
		if value.Int <= 0 {
			continue
		}
		if strings.HasSuffix(key, "Rating") || strings.HasSuffix(key, "Yards") {
			ratings[key] = int(value.Int)
		}
	}
	if len(ratings) > 0 {
		player.Ratings = ratings
	}

	if player.FirstName == "" && player.LastName == "" && player.Overall == nil && len(player.Ratings) == 0 {
		return nil
	}
	return player
}

// applyCanonicalTeamIndex rewrites player.TeamIndex from a Team table row
// number to the stable Team.TeamIndex ID. Clears the field when the row is
// missing or an FCS placeholder.
func applyCanonicalTeamIndex(player *PlayerExport, teams teamIndexMaps) {
	if player == nil || player.TeamIndex == nil {
		return
	}
	if id, ok := teams.exportID(*player.TeamIndex); ok {
		player.TeamIndex = &id
		return
	}
	player.TeamIndex = nil
}

// buildStatPlayerIdentity returns a lightweight PlayerExport used to label a
// stat line. It carries only identity fields (name, position, jersey, team) so
// per-game stats stay attributable without repeating the player's full ratings
// on every line — join on the player id for the full record.
func buildStatPlayerIdentity(record Record, teams teamIndexMaps) *PlayerExport {
	if len(record.Fields) == 0 {
		return nil
	}
	player := &PlayerExport{ID: record.Index}
	player.FirstName = stringField(record, "FirstName")
	player.LastName = stringField(record, "LastName")
	player.Position = stringField(record, "Position")
	if v, ok := record.Get("JerseyNum"); ok {
		setIntPtr(&player.Jersey, v.Int)
	}
	if v, ok := record.Get("TeamIndex"); ok {
		row := int(v.Int)
		player.TeamIndex = &row
		applyCanonicalTeamIndex(player, teams)
	}
	if player.FirstName == "" && player.LastName == "" {
		return nil
	}
	return player
}

// applyPlayerAth sets IsAth when the linked recruit lists alternate positions.
func applyPlayerAth(player *PlayerExport, alt1, alt2 string) {
	if player == nil {
		return
	}
	player.IsAth = playerIsAth(alt1, alt2)
}

// playerIsAth reports whether a player is an athlete (ATH) recruit — one who can
// play multiple positions. The game stores that as AlternatePosition1/2 on the
// recruit row, not as Position on the player.
func playerIsAth(alt1, alt2 string) bool {
	return isSetAlternatePosition(alt1) || isSetAlternatePosition(alt2)
}

func isSetAlternatePosition(pos string) bool {
	pos = strings.TrimSuffix(pos, "_")
	return pos != "" && pos != "Invalid"
}

func setIntPtr(dst **int, value int64) {
	if value == 0 {
		return
	}
	v := int(value)
	*dst = &v
}
