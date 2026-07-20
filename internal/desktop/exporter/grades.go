package exporter

import (
	"strconv"
	"strings"
)

// FormatLetterGrade converts enum-style grades (e.g. "Aplus") to display form ("A+").
func FormatLetterGrade(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.ContainsAny(raw, "+-") {
		return strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
	}

	key := strings.ToLower(strings.TrimSuffix(raw, "_"))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")

	switch key {
	case "aplus":
		return "A+"
	case "aminus":
		return "A-"
	case "a":
		return "A"
	case "bplus":
		return "B+"
	case "bminus":
		return "B-"
	case "b":
		return "B"
	case "cplus":
		return "C+"
	case "cminus":
		return "C-"
	case "c":
		return "C"
	case "dplus":
		return "D+"
	case "dminus":
		return "D-"
	case "d":
		return "D"
	case "fplus":
		return "F+"
	case "fminus":
		return "F-"
	case "f":
		return "F"
	default:
		if len(key) >= 2 {
			if strings.HasSuffix(key, "plus") {
				letter := strings.ToUpper(key[:1])
				if letter >= "A" && letter <= "F" {
					return letter + "+"
				}
			}
			if strings.HasSuffix(key, "minus") {
				letter := strings.ToUpper(key[:1])
				if letter >= "A" && letter <= "F" {
					return letter + "-"
				}
			}
		}
		return strings.ToUpper(raw)
	}
}

// LetterGradeBadgeClass returns a CSS class for a letter grade badge.
func LetterGradeBadgeClass(grade string) string {
	switch FormatLetterGrade(grade) {
	case "A+":
		return "grade-aplus"
	case "A", "A-":
		return "grade-a"
	case "B+", "B", "B-":
		return "grade-b"
	case "C+", "C", "C-":
		return "grade-c"
	case "D+", "D", "D-":
		return "grade-d"
	case "F", "F+", "F-":
		return "grade-f"
	default:
		return "grade-other"
	}
}

// ParsePipelineTier maps influence level/value to a 1–5 tier (5 is best).
func ParsePipelineTier(level string, influenceValue *int) (int, bool) {
	if n, ok := parseTierFromLevel(level); ok {
		return n, true
	}
	return tierFromInfluenceValue(influenceValue)
}

func parseTierFromLevel(level string) (int, bool) {
	if n, ok := parseTierNumber(level); ok {
		return n, true
	}

	key := strings.ToLower(strings.TrimSpace(level))
	key = strings.TrimSuffix(key, "_")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")
	key = strings.TrimPrefix(key, "influence")
	key = strings.TrimPrefix(key, "level")
	key = strings.TrimPrefix(key, "tier")

	if n, err := strconv.Atoi(key); err == nil && n >= 1 && n <= 5 {
		return n, true
	}

	switch key {
	case "five", "5":
		return 5, true
	case "four", "4":
		return 4, true
	case "three", "3":
		return 3, true
	case "two", "2":
		return 2, true
	case "one", "1":
		return 1, true
	case "pink":
		return 5, true
	case "blue":
		return 4, true
	case "gold":
		return 3, true
	case "silver":
		return 2, true
	case "bronze":
		return 1, true
	case "culturalpillar":
		return 5, true
	case "householdname":
		return 4, true
	case "popular":
		return 3, true
	case "respected":
		return 2, true
	case "nicheinterest":
		return 1, true
	case "unrecognized":
		return 1, true
	case "dominant", "elite", "max":
		return 5, true
	case "strong", "high":
		return 4, true
	case "medium", "average", "mid":
		return 3, true
	case "low", "weak":
		return 2, true
	case "minimal", "none", "invalid":
		return 0, false
	default:
		return 0, false
	}
}

func tierFromInfluenceValue(v *int) (int, bool) {
	if v == nil {
		return 0, false
	}
	n := *v
	if n >= 1 && n <= 5 {
		return n, true
	}
	if n >= 0 && n <= 4 {
		return n + 1, true
	}
	switch {
	case n >= 300:
		return 5, true
	case n >= 175:
		return 4, true
	case n >= 90:
		return 3, true
	case n >= 30:
		return 2, true
	case n >= 1:
		return 1, true
	default:
		return 0, false
	}
}

// FormatPipelineName turns enum pipeline ids into readable labels.
func FormatPipelineName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	var b strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := runes[i-1]
			if prev >= 'a' && prev <= 'z' {
				b.WriteByte(' ')
			} else if prev >= 'A' && prev <= 'Z' && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				b.WriteByte(' ')
			}
		}
		b.WriteRune(r)
	}
	return b.String()
}

func parseTierNumber(level string) (int, bool) {
	level = strings.TrimSpace(level)
	if level == "" {
		return 0, false
	}
	level = strings.TrimPrefix(strings.ToLower(level), "tier")
	level = strings.TrimSpace(level)
	n, err := strconv.Atoi(level)
	if err != nil || n < 1 || n > 5 {
		return 0, false
	}
	return n, true
}

// PipelineTierPercent returns bar fill percent for a 1–5 pipeline tier.
func PipelineTierPercent(tier int) int {
	if tier < 1 {
		return 0
	}
	if tier > 5 {
		return 100
	}
	return tier * 100 / 5
}

// PipelineTierBadgeClass returns a CSS class for a pipeline tier badge.
func PipelineTierBadgeClass(tier int) string {
	if tier < 1 || tier > 5 {
		return "pipeline-tier-unknown"
	}
	return "pipeline-tier-" + strconv.Itoa(tier)
}
