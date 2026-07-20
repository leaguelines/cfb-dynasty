package exporter

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// RecruitListRow is one row on the prospect list page.
type RecruitListRow struct {
	ID                int
	Href              string
	NationalRank      string
	Position          string
	IsAth             bool
	DisplayName       string
	StarRating        int
	StarsText         string
	QualityLabel      string
	QualityClass      string
	NIL               string
	Stage             string
	HomeState         string
	Hometown          string
	HomePipeline      string
	HomePipelineLabel string
	Class             string
	SchoolName        string
	SchoolLogo        string
	SchoolHref        string
	HasCommitGauge    bool
	CommitGaugeBar    int
	CommitGaugeLabel  string
}

// FormatRecruitDisplayName formats a recruit name like "K. Smith".
func FormatRecruitDisplayName(first, last string) string {
	first = strings.TrimSpace(first)
	last = strings.TrimSpace(last)
	if first == "" && last == "" {
		return ""
	}
	if first == "" {
		return last
	}
	if last == "" {
		return first
	}
	r, _ := utf8.DecodeRuneInString(first)
	return string(r) + ". " + last
}

// FormatRecruitClass turns enum class values into compact labels.
func FormatRecruitClass(class string) string {
	class = strings.TrimSpace(class)
	if class == "" {
		return ""
	}
	key := strings.ToLower(strings.ReplaceAll(class, "_", ""))
	switch key {
	case "highschool":
		return "HS"
	case "juniorcollegesophomore":
		return "JC (SO)"
	case "juniorcollegejunior":
		return "JC (JR)"
	case "juniorcollegesenior":
		return "JC (SR)"
	case "juniorcollegefreshman":
		return "JC (FR)"
	default:
		return strings.ReplaceAll(class, "_", " ")
	}
}

// FormatRecruitStage turns enum stage values into readable labels.
func FormatRecruitStage(stage string) string {
	stage = strings.TrimSpace(stage)
	if stage == "" || strings.EqualFold(stage, "Invalid") {
		return ""
	}
	key := strings.ToLower(stage)
	switch key {
	case "open":
		return "Open"
	case "softcommitted":
		return "Verbal"
	case "top3":
		return "Top 3"
	case "top5":
		return "Top 5"
	case "top10":
		return "Top 10"
	case "battle":
		return "Battle"
	case "signed":
		return "Signed"
	default:
		var b strings.Builder
		for i, r := range stage {
			if i > 0 && r >= 'A' && r <= 'Z' {
				b.WriteByte(' ')
			}
			b.WriteRune(r)
		}
		return b.String()
	}
}

// FormatQualityModifier returns a short label and CSS class for gem/bust/normal.
func FormatQualityModifier(raw string) (label, class string) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "GEM":
		return "Gem", "quality-gem"
	case "BUST":
		return "Bust", "quality-bust"
	default:
		return "", ""
	}
}

// ParseStarRating maps enum or numeric star rating strings to a 0–5 count.
func ParseStarRating(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if n, err := strconv.Atoi(raw); err == nil {
		if n < 0 {
			return 0
		}
		if n > 5 {
			return 5
		}
		return n
	}

	key := strings.ToLower(strings.TrimSuffix(raw, "_"))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")
	key = strings.TrimSuffix(key, "star")
	key = strings.TrimSuffix(key, "stars")

	switch key {
	case "5", "five":
		return 5
	case "4", "four":
		return 4
	case "3", "three":
		return 3
	case "2", "two":
		return 2
	case "1", "one":
		return 1
	default:
		return 0
	}
}

// StarRatingDisplay returns the numeric star count and a five-star text display.
func StarRatingDisplay(raw string) (count int, text string) {
	n := ParseStarRating(raw)
	return n, strings.Repeat("★", n) + strings.Repeat("☆", 5-n)
}

// FormatPlayerHeight formats height in inches as feet and inches.
func FormatPlayerHeight(inches *int) string {
	if inches == nil || *inches <= 0 {
		return ""
	}
	h := *inches
	return strconv.Itoa(h/12) + "' " + strconv.Itoa(h%12) + "\""
}

// FormatPlayerWeight formats weight in pounds.
func FormatPlayerWeight(weight *int) string {
	if weight == nil || *weight <= 0 {
		return ""
	}
	return strconv.Itoa(*weight) + " lbs"
}

// FormatHeightWeight combines height and weight for detail headers.
func FormatHeightWeight(height, weight *int) string {
	h := FormatPlayerHeight(height)
	w := FormatPlayerWeight(weight)
	switch {
	case h != "" && w != "":
		return h + " | " + w
	case h != "":
		return h
	default:
		return w
	}
}

// FormatHometown combines town and state.
func FormatHometown(town, state string) string {
	town = strings.TrimSpace(town)
	state = strings.TrimSpace(state)
	switch {
	case town != "" && state != "":
		return town + ", " + state
	case town != "":
		return town
	default:
		return state
	}
}

// FormatNILDisplay returns expected NIL for recruits when exported.
func FormatNILDisplay(p *dynasty.PlayerExport) string {
	if p == nil {
		return ""
	}
	if p.NILBaseValue != nil && *p.NILBaseValue > 0 {
		return strconv.Itoa(*p.NILBaseValue)
	}
	if p.NILCompensation != nil && *p.NILCompensation > 0 {
		return strconv.Itoa(*p.NILCompensation)
	}
	return ""
}

// RecruitStageRank orders recruit stages for sorting (higher is later in process).
func RecruitStageRank(stage string) int {
	switch strings.ToLower(strings.TrimSpace(stage)) {
	case "signed":
		return 7
	case "softcommitted":
		return 6
	case "battle":
		return 5
	case "top3":
		return 4
	case "top5":
		return 3
	case "top10":
		return 2
	case "open":
		return 1
	default:
		return 0
	}
}
