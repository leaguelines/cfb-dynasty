package dynasty

import "fmt"

// skillGroupCapCount is the number of SkillGroupCap fields on a Player row.
const skillGroupCapCount = 6

// skillGroupCapsFromRecord reads SkillGroupCap1..6 from a Player row. It returns
// the positional caps, their total, and whether any caps were present.
//
// Bucket names come from the tuning FTC (PlayerSkillGroup / PlayerSkillGroupBucket)
// and are attached separately via applySkillGroupLabels when tuning data is available.
func skillGroupCapsFromRecord(record Record) (caps []int, total int, ok bool) {
	caps = make([]int, skillGroupCapCount)
	for i := 0; i < skillGroupCapCount; i++ {
		value, present := record.Get(fmt.Sprintf("SkillGroupCap%d", i+1))
		if !present {
			return nil, 0, false
		}
		caps[i] = int(value.Int)
		total += caps[i]
	}
	if total == 0 {
		return nil, 0, false
	}
	return caps, total, true
}
