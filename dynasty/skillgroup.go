package dynasty

import "fmt"

// skillGroupCapCount is the number of SkillGroupCap fields on a Player row.
const skillGroupCapCount = 6

// skillGroupBucketMax is the per-bucket upgrade capacity stored in SkillGroupCap1..6.
// Greyed-out cap slots in the UI equal skillGroupBucketMax minus the saved value.
const skillGroupBucketMax = 20

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
