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
			setIntPtr(&player.TeamIndex, value.Int)
		}
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

func setIntPtr(dst **int, value int64) {
	if value == 0 {
		return
	}
	v := int(value)
	*dst = &v
}
