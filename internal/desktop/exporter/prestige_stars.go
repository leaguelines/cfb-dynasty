package exporter

import (
	"strconv"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// PrestigeStar is one slot in a five-star prestige display.
type PrestigeStar struct {
	State string // full, half, empty
}

// BuildPrestigeStars returns five star slots from half-star units (10 = 5★, 9 = 4.5★).
func BuildPrestigeStars(halfUnits int) []PrestigeStar {
	if halfUnits <= 0 {
		return nil
	}
	if halfUnits > 10 {
		halfUnits = 10
	}
	full := halfUnits / 2
	hasHalf := halfUnits%2 == 1
	stars := make([]PrestigeStar, 5)
	for i := range stars {
		switch {
		case i < full:
			stars[i].State = "full"
		case i == full && hasHalf:
			stars[i].State = "half"
		default:
			stars[i].State = "empty"
		}
	}
	return stars
}

func prestigeHalfStars(t dynasty.TeamExport) int {
	if t.TeamPrestige != nil && *t.TeamPrestige > 0 {
		return *t.TeamPrestige
	}
	return 0
}

func prestigeStarsLabel(halfUnits int) string {
	if halfUnits <= 0 {
		return ""
	}
	if halfUnits > 10 {
		halfUnits = 10
	}
	full := halfUnits / 2
	if halfUnits%2 == 1 {
		if full == 0 {
			return "0.5"
		}
		return strconv.Itoa(full) + ".5"
	}
	return strconv.Itoa(full)
}
