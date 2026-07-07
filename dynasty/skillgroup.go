package dynasty

import "fmt"

// skillGroupCapCount is the number of SkillGroupCap fields on a Player row.
const skillGroupCapCount = 6

// skillGroupCapsFromRecord reads SkillGroupCap1..6 from a Player row. It returns
// the positional caps, their total, and whether any caps were present.
//
// The caps are exported as opaque, positional values. The game groups these into
// six position-specific skill groups, but the group *names* are not stored in the
// dynasty save (they live in the tuning FTC, which is not yet extractable for
// CFB 27), so we deliberately do not label the slots until definitive names are
// available. Likewise, only the per-group cap is present here — the individual
// ratings inside each group are tuning-driven and not in the save.
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
