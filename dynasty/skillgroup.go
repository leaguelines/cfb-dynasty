package dynasty

import (
	"fmt"
	"strings"
)

// skillGroupCapCount is the number of SkillGroupCap fields on a Player row.
const skillGroupCapCount = 6

// skillGroupBucketMax is the per-bucket upgrade capacity stored in SkillGroupCap1..6.
// Greyed-out cap slots in the UI equal skillGroupBucketMax minus the saved value.
const skillGroupBucketMax = 20

var (
	skillGroupLevelMinimumsFBTE = [...]int{
		59, 61, 63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 84, 85, 86, 87, 88, 89,
	}
	skillGroupLevelMinimumsWRMIKEKP = [...]int{
		61, 64, 66, 69, 71, 74, 76, 79, 81, 84, 86, 89, 91, 93, 95, 96, 97, 98, 99,
	}
	skillGroupLevelMinimumsDefault = [...]int{
		63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 85, 87, 89, 91, 93, 95, 97, 99,
	}
)

type skillGroupSlots struct {
	unlocked []int
	capped   []int
}

// skillGroupSlotsFromRecord reads SkillGroupCap1..6 from a Player row. Each value is
// the number of unlocked upgrade slots in that bucket (0..skillGroupBucketMax).
// Capped/greyed UI slots are skillGroupBucketMax minus the saved value.
func skillGroupSlotsFromRecord(record Record) (slots skillGroupSlots, ok bool) {
	unlocked := make([]int, skillGroupCapCount)
	for i := 0; i < skillGroupCapCount; i++ {
		value, present := record.Get(fmt.Sprintf("SkillGroupCap%d", i+1))
		if !present {
			return skillGroupSlots{}, false
		}
		unlocked[i] = int(value.Int)
	}
	totalUnlocked := sumInts(unlocked)
	if totalUnlocked == 0 {
		return skillGroupSlots{}, false
	}
	capped := make([]int, skillGroupCapCount)
	for i, value := range unlocked {
		capped[i] = skillGroupBucketMax - value
	}
	return skillGroupSlots{unlocked: unlocked, capped: capped}, true
}

func sumInts(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func totalCappedSlots(slots skillGroupSlots) int {
	return skillGroupBucketMax*skillGroupCapCount - sumInts(slots.unlocked)
}

// SkillGroupOverall returns the weighted OVR for a complete group attribute set.
// Tier metadata must be complete to apply 15/3/1 weights; otherwise all known
// group attributes receive equal weight.
func SkillGroupOverall(attributes []SkillGroupAttributeExport) (int, bool) {
	if len(attributes) == 0 {
		return 0, false
	}

	sum := 0
	weightedSum := 0
	totalWeight := 0
	tiersComplete := true
	for _, attribute := range attributes {
		if attribute.Rating == nil {
			return 0, false
		}
		rating := *attribute.Rating
		sum += rating

		weight := 0
		switch strings.ToLower(strings.TrimSpace(attribute.Tier)) {
		case "primary":
			weight = 15
		case "secondary":
			weight = 3
		case "tertiary":
			weight = 1
		default:
			tiersComplete = false
		}
		weightedSum += rating * weight
		totalWeight += weight
	}
	if !tiersComplete {
		return sum / len(attributes), true
	}
	return weightedSum / totalWeight, true
}

func skillGroupCurrentLevel(position string, overall, ceiling int) int {
	if ceiling <= 0 {
		return 0
	}

	minimums := skillGroupLevelMinimumsDefault
	switch position {
	case "FB", "TE":
		minimums = skillGroupLevelMinimumsFBTE
	case "WR", "MIKE", "MLB", "K", "P":
		minimums = skillGroupLevelMinimumsWRMIKEKP
	}

	level := 1
	for i, minimum := range minimums {
		if overall < minimum {
			break
		}
		level = i + 2
	}
	if level > ceiling {
		return ceiling
	}
	return level
}

func skillGroupCurrentLevels(position string, groups []SkillGroupExport, ceilings []int) []int {
	if len(groups) != skillGroupCapCount || len(ceilings) != skillGroupCapCount {
		return nil
	}
	levels := make([]int, skillGroupCapCount)
	for i, group := range groups {
		overall, ok := SkillGroupOverall(group.Attributes)
		if !ok {
			return nil
		}
		levels[i] = skillGroupCurrentLevel(position, overall, ceilings[i])
	}
	return levels
}
