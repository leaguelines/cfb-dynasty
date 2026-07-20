package exporter

import (
	"math"
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

const skillGroupBucketMax = 20

// SkillGroupSegmentKind is one segment in the upgrade bar.
type SkillGroupSegmentKind string

const (
	SkillSegFilled    SkillGroupSegmentKind = "filled"
	SkillSegAvailable SkillGroupSegmentKind = "available"
	SkillSegCapped    SkillGroupSegmentKind = "capped"
)

// SkillGroupSegment is one block in the game-style cap bar.
type SkillGroupSegment struct {
	Kind SkillGroupSegmentKind
}

// SkillGroupSubRating is one attribute within a skill group bucket.
type SkillGroupSubRating struct {
	Name       string
	Value      int
	BarPercent int
}

// SkillGroupRow is one skill bucket for the player detail page.
type SkillGroupRow struct {
	Label          string
	Rating         int
	HasRating      bool
	Segments       []SkillGroupSegment
	SubRatings     []SkillGroupSubRating
	HasSubRatings  bool
	CappedSlots    int
	UnlockedSlots  int
	AttributeCount int
}

// SkillGroupRows builds game-style skill group rows from a player export.
func SkillGroupRows(p *dynasty.PlayerExport) []SkillGroupRow {
	if p == nil {
		return nil
	}
	if len(p.SkillGroups) > 0 {
		rows := make([]SkillGroupRow, 0, len(p.SkillGroups))
		for _, g := range p.SkillGroups {
			rows = append(rows, buildSkillGroupRow(g, p.Ratings))
		}
		return rows
	}
	if len(p.SkillGroupLabels) == 0 && len(p.SkillGroupCaps) == 0 && len(p.SkillGroupUnlockedSlots) == 0 {
		return nil
	}
	n := len(p.SkillGroupLabels)
	if len(p.SkillGroupCaps) > n {
		n = len(p.SkillGroupCaps)
	}
	if len(p.SkillGroupUnlockedSlots) > n {
		n = len(p.SkillGroupUnlockedSlots)
	}
	if len(p.SkillGroupAttributeCounts) > n {
		n = len(p.SkillGroupAttributeCounts)
	}
	rows := make([]SkillGroupRow, 0, n)
	for i := 0; i < n; i++ {
		label := ""
		if i < len(p.SkillGroupLabels) {
			label = strings.TrimSpace(p.SkillGroupLabels[i])
		}
		if label == "" {
			label = "Slot " + fmtInt(i+1)
		}
		g := dynasty.SkillGroupExport{Slot: i + 1, Label: label}
		if i < len(p.SkillGroupCaps) {
			g.CappedSlots = p.SkillGroupCaps[i]
		}
		if i < len(p.SkillGroupUnlockedSlots) {
			g.UnlockedSlots = p.SkillGroupUnlockedSlots[i]
		}
		if i < len(p.SkillGroupAttributeCounts) {
			g.AttributeCount = p.SkillGroupAttributeCounts[i]
		}
		rows = append(rows, buildSkillGroupRow(g, p.Ratings))
	}
	return rows
}

func buildSkillGroupRow(g dynasty.SkillGroupExport, ratings map[string]int) SkillGroupRow {
	label := strings.TrimSpace(g.Label)
	if label == "" {
		label = "Slot " + fmtInt(g.Slot)
	}
	sub := skillGroupSubRatings(g.Attributes, label, ratings)
	rating, hasRating := skillGroupAverage(sub)
	row := SkillGroupRow{
		Label:          label,
		Rating:         rating,
		HasRating:      hasRating,
		Segments:       buildSkillGroupSegments(rating, g.UnlockedSlots, g.CappedSlots),
		SubRatings:     sub,
		HasSubRatings:  len(sub) > 0,
		CappedSlots:    g.CappedSlots,
		UnlockedSlots:  g.UnlockedSlots,
		AttributeCount: g.AttributeCount,
	}
	return row
}

func buildSkillGroupSegments(rating, unlocked, capped int) []SkillGroupSegment {
	if capped < 0 {
		capped = 0
	}
	if unlocked < 0 {
		unlocked = 0
	}
	playable := skillGroupBucketMax - capped
	if playable < 0 {
		playable = 0
	}
	if unlocked > playable {
		unlocked = playable
	}

	filled := 0
	if rating > 0 && playable > 0 {
		filled = int(math.Round(float64(rating) / 99.0 * float64(playable)))
		if filled > playable {
			filled = playable
		}
	}
	available := playable - filled
	if available < 0 {
		available = 0
	}

	segments := make([]SkillGroupSegment, 0, skillGroupBucketMax)
	for i := 0; i < filled; i++ {
		segments = append(segments, SkillGroupSegment{Kind: SkillSegFilled})
	}
	for i := 0; i < available; i++ {
		segments = append(segments, SkillGroupSegment{Kind: SkillSegAvailable})
	}
	for i := 0; i < capped; i++ {
		segments = append(segments, SkillGroupSegment{Kind: SkillSegCapped})
	}
	return segments
}

func skillGroupSubRatings(attrs []dynasty.SkillGroupAttributeExport, label string, ratings map[string]int) []SkillGroupSubRating {
	if len(attrs) > 0 {
		out := make([]SkillGroupSubRating, 0, len(attrs))
		for _, attr := range attrs {
			val, ok := skillGroupAttributeValue(attr, ratings)
			if !ok {
				continue
			}
			name := strings.TrimSpace(attr.Name)
			if name == "" {
				name = FormatRatingName(attr.RatingKey)
			}
			out = append(out, SkillGroupSubRating{
				Name:       name,
				Value:      val,
				BarPercent: ratingBarPercent(val),
			})
		}
		return out
	}
	return skillGroupSubRatingsHeuristic(label, ratings)
}

func skillGroupAttributeValue(attr dynasty.SkillGroupAttributeExport, ratings map[string]int) (int, bool) {
	if attr.Rating != nil {
		return *attr.Rating, true
	}
	if attr.RatingKey != "" && ratings != nil {
		if val, ok := ratings[attr.RatingKey]; ok {
			return val, true
		}
	}
	if attr.PlayerAbility != "" && ratings != nil {
		key := attr.PlayerAbility + "Rating"
		if val, ok := ratings[key]; ok {
			return val, true
		}
	}
	return 0, false
}

func skillGroupAverage(sub []SkillGroupSubRating) (int, bool) {
	if len(sub) == 0 {
		return 0, false
	}
	sum := 0
	for _, r := range sub {
		sum += r.Value
	}
	return sum / len(sub), true
}

func ratingBarPercent(value int) int {
	if value <= 0 {
		return 0
	}
	pct := int(math.Round(float64(value) / 99.0 * 100))
	if pct > 100 {
		return 100
	}
	return pct
}

// skillGroupSubRatingsHeuristic is a fallback when tuning attributes are unavailable.
func skillGroupSubRatingsHeuristic(label string, ratings map[string]int) []SkillGroupSubRating {
	if len(ratings) == 0 {
		return nil
	}
	terms := skillGroupMatchTerms(label)
	type match struct {
		key string
		val int
	}
	var matches []match
	for key, val := range ratings {
		if key == "OverallRating" {
			continue
		}
		lower := strings.ToLower(key)
		for _, term := range terms {
			if strings.Contains(lower, term) {
				matches = append(matches, match{key: key, val: val})
				break
			}
		}
	}
	if len(matches) == 0 {
		return nil
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].key < matches[j].key
	})
	out := make([]SkillGroupSubRating, len(matches))
	for i, m := range matches {
		out[i] = SkillGroupSubRating{
			Name:       FormatRatingName(m.key),
			Value:      m.val,
			BarPercent: ratingBarPercent(m.val),
		}
	}
	return out
}

func skillGroupMatchTerms(label string) []string {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil
	}
	key := strings.ToLower(strings.ReplaceAll(label, " ", ""))
	if terms, ok := skillGroupSynonyms[key]; ok {
		return terms
	}
	return []string{key}
}

// skillGroupSynonyms maps bucket labels to rating key fragments for legacy exports.
var skillGroupSynonyms = map[string][]string{
	"accuracy":    {"throwaccuracy", "throwunderpressure", "throwontherun"},
	"power":       {"throwpower"},
	"iq":          {"awareness", "playrecognition", "playaction", "carrying"},
	"elusiveness": {"breaksack", "agility", "changeofdirection"},
	"quickness":   {"speed", "acceleration"},
	"health":      {"injury", "stamina"},
}
