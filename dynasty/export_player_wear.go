package dynasty

import "strings"

var wearAndTearFields = []string{
	"WearAndTear_Back",
	"WearAndTear_LAnkle", "WearAndTear_LArm", "WearAndTear_LElbow", "WearAndTear_LFoot",
	"WearAndTear_LHand", "WearAndTear_LHip", "WearAndTear_LKnee", "WearAndTear_LLeg", "WearAndTear_LShoulder",
	"WearAndTear_RAnkle", "WearAndTear_RArm", "WearAndTear_RElbow", "WearAndTear_RFoot",
	"WearAndTear_RHand", "WearAndTear_RHip", "WearAndTear_Rib", "WearAndTear_RKnee", "WearAndTear_RLeg", "WearAndTear_RShoulder",
}

func wearAndTearFromRecord(record Record) map[string]int {
	out := make(map[string]int)
	for _, field := range wearAndTearFields {
		v, ok := intFieldOK(record, field)
		if !ok || v >= 10 {
			continue
		}
		key := strings.TrimPrefix(field, "WearAndTear_")
		out[key] = v
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func applyPlayerRosterFields(player *PlayerExport, record Record) {
	if player == nil {
		return
	}
	player.Personality = normalizeEnum(stringField(record, "Personality"))
	player.PracticePlan = normalizeEnum(stringField(record, "PracticePlan"))
	setOptionalPositiveInt(record, "Fatigue", &player.Fatigue)
	setOptionalNonNegativeInt(record, "AbsoluteTransferChance", &player.TransferChance)
	setOptionalPositiveInt(record, "PLYR_CONSECYEARSWITHTEAM", &player.ConsecutiveYearsWithTeam)
	setOptionalPositiveInt(record, "PLYR_DRAFTPICK", &player.DraftPick)
	setOptionalIntMax(record, "PLYR_DRAFTROUND", &player.DraftRound, 63)
	player.WearAndTear = wearAndTearFromRecord(record)
}

func setOptionalNonNegativeInt(record Record, name string, dst **int) {
	v, ok := intFieldOK(record, name)
	if !ok || v < 0 {
		return
	}
	*dst = &v
}

func setOptionalIntMax(record Record, name string, dst **int, max int) {
	v, ok := intFieldOK(record, name)
	if !ok || v < 0 || v > max {
		return
	}
	*dst = &v
}
