package exporter

import (
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// PhysicalAbilityRow is one ability icon row for player detail views.
type PhysicalAbilityRow struct {
	Name      string
	Slot      int
	Tier      int
	TierLabel string
	TierClass string
	ImagePath string
}

// MentalAbilityTierLabel maps mental ability rank (1–3) to a tier name.
func MentalAbilityTierLabel(rank int) string {
	switch rank {
	case 3:
		return "Gold"
	case 2:
		return "Silver"
	case 1:
		return "Bronze"
	default:
		return ""
	}
}

// MentalAbilityTierLevel maps mental ability rank to display tier levels (0–3, no platinum).
func MentalAbilityTierLevel(rank int) int {
	switch rank {
	case 3:
		return 3
	case 2:
		return 2
	case 1:
		return 1
	default:
		return 0
	}
}

// FormatMentalAbilityName turns export enum values into display labels.
func FormatMentalAbilityName(enum string) string {
	if asset, ok := mentalAbilityAssets[enum]; ok && asset.Label != "" {
		return asset.Label
	}
	return enum
}

// MentalAbilityImagePath returns the static URL for a mental ability icon.
func MentalAbilityImagePath(enum string) string {
	asset, ok := mentalAbilityAssets[enum]
	if !ok || asset.Slug == "" {
		return ""
	}
	return "/static/abilities/" + asset.Slug + "." + asset.Ext
}

// PhysicalAbilityTierLevel maps export tier names to display levels (0–4).
func PhysicalAbilityTierLevel(tier string) int {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case "platinum":
		return 4
	case "gold":
		return 3
	case "silver":
		return 2
	case "bronze":
		return 1
	default:
		return 0
	}
}

// PhysicalAbilityTierClass returns the CSS class for a tier level.
func PhysicalAbilityTierClass(level int) string {
	switch level {
	case 4:
		return "ability-tier-platinum"
	case 3:
		return "ability-tier-gold"
	case 2:
		return "ability-tier-silver"
	case 1:
		return "ability-tier-bronze"
	default:
		return "ability-tier-locked"
	}
}

// PhysicalAbilityImagePath returns the static URL for an ability icon, or "" if unknown.
func PhysicalAbilityImagePath(name string) string {
	asset, ok := physicalAbilityAssets[name]
	if !ok || asset.Slug == "" {
		return ""
	}
	return "/static/abilities/" + asset.Slug + "." + asset.Ext
}

func physicalAbilityRow(name string, slot int, tier string) PhysicalAbilityRow {
	level := PhysicalAbilityTierLevel(tier)
	return PhysicalAbilityRow{
		Name:      name,
		Slot:      slot,
		Tier:      level,
		TierLabel: tier,
		TierClass: PhysicalAbilityTierClass(level),
		ImagePath: PhysicalAbilityImagePath(name),
	}
}

// BuildPhysicalAbilityRows formats unlocked physical abilities from export data.
func BuildPhysicalAbilityRows(ab []dynasty.PlayerPhysicalAbilityExport) []PhysicalAbilityRow {
	if len(ab) == 0 {
		return nil
	}
	rows := make([]PhysicalAbilityRow, 0, len(ab))
	for _, a := range ab {
		if a.Name == "" {
			continue
		}
		rows = append(rows, physicalAbilityRow(a.Name, a.Slot, a.Tier))
	}
	return rows
}

// BuildArchetypePhysicalAbilityRows returns all archetype ability slots in order,
// applying unlocked tiers from export data where available.
func BuildArchetypePhysicalAbilityRows(position, archetype string, unlocked []dynasty.PlayerPhysicalAbilityExport) []PhysicalAbilityRow {
	if archetype == "" {
		return nil
	}
	tierBySlot := make(map[int]string, len(unlocked))
	for _, a := range unlocked {
		tierBySlot[a.Slot] = a.Tier
	}
	rows := make([]PhysicalAbilityRow, 0, 5)
	for slot := 1; slot <= 5; slot++ {
		name := ArchetypePhysicalAbilityName(position, archetype, slot)
		if name == "" {
			continue
		}
		rows = append(rows, physicalAbilityRow(name, slot, tierBySlot[slot]))
	}
	if len(rows) == 0 {
		return nil
	}
	return rows
}

func mentalAbilityRow(enum string, slot int, rank int) PhysicalAbilityRow {
	level := MentalAbilityTierLevel(rank)
	return PhysicalAbilityRow{
		Name:      FormatMentalAbilityName(enum),
		Slot:      slot,
		Tier:      level,
		TierLabel: MentalAbilityTierLabel(rank),
		TierClass: PhysicalAbilityTierClass(level),
		ImagePath: MentalAbilityImagePath(enum),
	}
}

// BuildMentalAbilityRows formats unlocked mental abilities from export data.
func BuildMentalAbilityRows(ab []dynasty.PlayerMentalAbilityExport) []PhysicalAbilityRow {
	if len(ab) == 0 {
		return nil
	}
	rows := make([]PhysicalAbilityRow, 0, len(ab))
	for i, a := range ab {
		if a.Name == "" {
			continue
		}
		rank := 0
		if a.Rank != nil {
			rank = *a.Rank
		}
		rows = append(rows, mentalAbilityRow(a.Name, i+1, rank))
	}
	if len(rows) == 0 {
		return nil
	}
	return rows
}
