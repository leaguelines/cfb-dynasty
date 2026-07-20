package exporter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

func fmtSchoolInterestList(schools []dynasty.RecruitingSchoolInterestExport) string {
	if len(schools) == 0 {
		return ""
	}
	parts := make([]string, len(schools))
	for i, s := range schools {
		name := strings.TrimSpace(s.TeamName)
		if name == "" {
			name = fmt.Sprintf("team %d", s.TeamID)
		}
		parts[i] = fmt.Sprintf("%s:%d", name, s.Influence)
	}
	return strings.Join(parts, "; ")
}

func fmtActivePitches(pitches []dynasty.RecruitingPitchExport) string {
	if len(pitches) == 0 {
		return ""
	}
	parts := make([]string, len(pitches))
	for i, p := range pitches {
		label := FormatSchoolPitchName(p.Pitch)
		if label == "" {
			label = strings.TrimSpace(p.Pitch)
		}
		if p.Intensity != "" {
			parts[i] = label + " (" + p.Intensity + ")"
		} else {
			parts[i] = label
		}
	}
	return strings.Join(parts, "; ")
}

func fmtScheduledVisit(v *dynasty.RecruitingVisitExport) string {
	if v == nil {
		return ""
	}
	parts := make([]string, 0, 3)
	if v.Activity != "" {
		parts = append(parts, v.Activity)
	}
	if v.Week > 0 {
		parts = append(parts, fmt.Sprintf("week %d", v.Week))
	}
	if v.WeekType != "" {
		parts = append(parts, v.WeekType)
	}
	return strings.Join(parts, ", ")
}

func fmtMapInt(m map[string]int) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s:%d", k, m[k])
	}
	return strings.Join(parts, "; ")
}

func fmtSkillGroupsSummary(groups []dynasty.SkillGroupExport) string {
	if len(groups) == 0 {
		return ""
	}
	parts := make([]string, len(groups))
	for i, g := range groups {
		label := strings.TrimSpace(g.Label)
		if label == "" {
			label = fmt.Sprintf("slot %d", g.Slot)
		}
		parts[i] = fmt.Sprintf("%s capped=%d unlocked=%d", label, g.CappedSlots, g.UnlockedSlots)
	}
	return strings.Join(parts, "; ")
}

func fmtPositionRatings(ratings map[string]int) string {
	if len(ratings) == 0 {
		return ""
	}
	keys := make([]string, 0, len(ratings))
	for k := range ratings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s:%d", k, ratings[k])
	}
	return strings.Join(parts, "; ")
}
