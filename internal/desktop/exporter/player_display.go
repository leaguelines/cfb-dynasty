package exporter

import (
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// WearRow is one body-part wear value for display.
type WearRow struct {
	BodyPart string
	Level    int
}

// WearRows returns sorted wear-and-tear rows with non-zero values.
func WearRows(p *dynasty.PlayerExport) []WearRow {
	if p == nil || len(p.WearAndTear) == 0 {
		return nil
	}
	rows := make([]WearRow, 0, len(p.WearAndTear))
	for part, level := range p.WearAndTear {
		if level <= 0 {
			continue
		}
		rows = append(rows, WearRow{BodyPart: formatBodyPart(part), Level: level})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Level != rows[j].Level {
			return rows[i].Level > rows[j].Level
		}
		return rows[i].BodyPart < rows[j].BodyPart
	})
	return rows
}

func formatBodyPart(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range raw {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// SchoolInterestRow is one school's interest in a recruit.
type SchoolInterestRow struct {
	TeamName   string
	TeamID     int
	TeamLogo   string
	Influence  int
	Href       string
	BarPercent int
}

// SchoolInterestRows formats school interest for UI tables.
func SchoolInterestRows(schools []dynasty.RecruitingSchoolInterestExport, shortNames map[int]string) []SchoolInterestRow {
	if len(schools) == 0 {
		return nil
	}
	sorted := append([]dynasty.RecruitingSchoolInterestExport(nil), schools...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Influence != sorted[j].Influence {
			return sorted[i].Influence > sorted[j].Influence
		}
		return sorted[i].TeamName < sorted[j].TeamName
	})
	max := sorted[0].Influence
	rows := make([]SchoolInterestRow, len(sorted))
	for i, s := range sorted {
		row := SchoolInterestRow{
			TeamName:  s.TeamName,
			TeamID:    s.TeamID,
			TeamLogo:  TeamLogoForTeamID(shortNames, s.TeamID),
			Influence: s.Influence,
		}
		if s.TeamID > 0 {
			row.Href = SchoolClassHref(s.TeamID)
		}
		if max > 0 {
			row.BarPercent = s.Influence * 100 / max
		}
		rows[i] = row
	}
	return rows
}

// PitchRow is one active recruiting pitch.
type PitchRow struct {
	Pitch     string
	Intensity string
}

// ActivePitchRows formats active pitches for display.
func ActivePitchRows(pitches []dynasty.RecruitingPitchExport) []PitchRow {
	if len(pitches) == 0 {
		return nil
	}
	rows := make([]PitchRow, len(pitches))
	for i, p := range pitches {
		label := FormatSchoolPitchName(p.Pitch)
		if label == "" {
			label = strings.TrimSpace(p.Pitch)
		}
		rows[i] = PitchRow{Pitch: label, Intensity: p.Intensity}
	}
	return rows
}

// FindRecruitingTarget returns pursuit state for a recruit id.
func FindRecruitingTarget(e dynasty.Export, recruitID int) (*dynasty.RecruitingTargetExport, bool) {
	for i := range e.Recruiting {
		if e.Recruiting[i].RecruitID == recruitID {
			return &e.Recruiting[i], true
		}
	}
	return nil, false
}

// FormatRecruitStageAdvance turns enum values into readable labels.
func FormatRecruitStageAdvance(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "Invalid") {
		return ""
	}
	key := strings.ToLower(raw)
	switch key {
	case "none":
		return ""
	case "top10":
		return "Top 10"
	case "top5":
		return "Top 5"
	case "top3":
		return "Top 3"
	case "battle":
		return "Battle"
	case "signed":
		return "Signed"
	default:
		return FormatRecruitStage(raw)
	}
}

// SeasonBanner summarizes calendar context for the dashboard.
type SeasonBanner struct {
	HasSeason     bool
	Year          int
	Week          int
	WeekType      string
	Phase         string
	ActivePeriods []string
}

// SeasonBannerView builds dashboard season context from an export.
func SeasonBannerView(e dynasty.Export) SeasonBanner {
	if e.Season == nil {
		return SeasonBanner{}
	}
	s := e.Season
	banner := SeasonBanner{
		HasSeason: true,
		Year:      s.Year,
		Week:      s.Week,
		WeekType:  s.WeekType,
		Phase:     s.Phase,
	}
	if s.Periods != nil {
		banner.ActivePeriods = activeSeasonPeriods(s.Periods)
	}
	return banner
}

func activeSeasonPeriods(p *dynasty.SeasonPeriodsExport) []string {
	if p == nil {
		return nil
	}
	type period struct {
		label string
		on    bool
	}
	periods := []period{
		{"Recruiting", p.IsRecruitingPeriodActive},
		{"Scouting", p.IsScoutingPeriodActive},
		{"Visiting", p.IsVisitingPeriodActive},
		{"Pitching", p.IsPitchingPeriodActive},
		{"Scholarships", p.IsScholarshipPeriodActive},
		{"Signing", p.IsSigningPeriodActive},
		{"Transfer portal", p.IsTransferPortalNewlyAvailable},
		{"Transfer signing", p.IsTransferSignPeriodActive},
		{"Draft", p.IsDraftPeriodActive},
		{"Draft scouting", p.IsDraftScoutingActive},
		{"Goals", p.IsGoalsPeriodActive},
		{"Coaching carousel", p.IsCarouselPeriodActive},
		{"Staff hiring", p.IsStaffHiringPeriodActive},
		{"Weekly awards", p.IsWeeklyAwardPeriodActive},
		{"Annual awards", p.IsAnnualAwardPeriodActive},
	}
	out := make([]string, 0, len(periods))
	for _, item := range periods {
		if item.on {
			out = append(out, item.label)
		}
	}
	return out
}
