package dynasty

import "strings"

// statRecordPosition returns the position group for a record-book entry. The
// save only populates position on the leader row for each stat category; lower
// ranks leave the PositionGroupE field at its schema default (QB). When the
// stored value is still that default, infer the group from the stat type.
func statRecordPosition(record Record, statType string) string {
	pos := stringField(record, "position")
	if pos != "" && pos != "QB" {
		return pos
	}
	if inferred, ok := positionGroupFromStatType(statType); ok {
		return inferred
	}
	return pos
}

// positionGroupFromStatType maps a StatType enum name to the PositionGroupE
// bucket that typically holds that record (QB, RB, WR, DL, LB, DB, etc.).
func positionGroupFromStatType(statType string) (string, bool) {
	if statType == "" {
		return "", false
	}
	switch {
	case strings.HasPrefix(statType, "Receive"), strings.HasPrefix(statType, "Receiving"):
		return "WR", true
	case strings.HasPrefix(statType, "Rush"), strings.HasPrefix(statType, "Rushing"):
		return "RB", true
	case strings.HasPrefix(statType, "Pass"), strings.HasPrefix(statType, "Passing"):
		return "QB", true
	case statType == "DefensiveInts" || strings.HasPrefix(statType, "DefensiveInts"),
		statType == "DefensivePassDeflections", statType == "DefensiveCatchesAllowed",
		strings.HasPrefix(statType, "InterceptionReturn"):
		return "DB", true
	case statType == "DefensiveSacks" || statType == "DefensiveHalfSack",
		strings.HasPrefix(statType, "DLine"):
		return "DL", true
	case strings.HasPrefix(statType, "DefensiveTackle"), statType == "DefensiveTotalTackles",
		statType == "DefensiveBigHits", statType == "DefensivePassYardsAllowed",
		statType == "DefensiveRushYardsAllowed":
		return "LB", true
	case strings.HasPrefix(statType, "KickReturn"), strings.HasPrefix(statType, "PuntReturn"):
		return "WR", true
	case strings.HasPrefix(statType, "FieldGoal"), strings.HasPrefix(statType, "ExtraPoint"),
		strings.HasPrefix(statType, "Punt"), strings.HasPrefix(statType, "KickNumber"),
		statType == "KickTouchbacks", statType == "PuntTouchbacks", statType == "Touchbacks":
		return "K", true
	case statType == "Pancakes", statType == "SacksAllowed":
		return "OL", true
	default:
		return "", false
	}
}
