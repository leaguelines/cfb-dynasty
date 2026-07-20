package exporter

import (
	"strconv"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// FormatRatingName turns export rating keys into scouting labels.
func FormatRatingName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = strings.TrimSuffix(name, "Rating")
	return FormatPipelineName(name)
}

// FormatSchoolPitchName turns enum pitch keys into readable school pitch labels.
func FormatSchoolPitchName(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "Invalid") {
		return ""
	}
	key := strings.ToLower(strings.TrimSuffix(raw, "_"))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")

	if label, ok := schoolPitchLabels[key]; ok {
		return label
	}
	return FormatPipelineName(raw)
}

// FormatRecruitMotivations returns readable motivation/pitch labels for a recruit.
func FormatRecruitMotivations(p *dynasty.PlayerExport) string {
	if p == nil {
		return ""
	}
	if parts := formatRawMotivations(p.Motivations); len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	pitch := normalizePitchKey(p.IdealRecruitingPitch)
	if pitch == "" || strings.EqualFold(pitch, "invalid") {
		return ""
	}
	if mots := idealPitchMotivations[pitch]; len(mots) > 0 {
		return strings.Join(mots, ", ")
	}
	return FormatPipelineName(p.IdealRecruitingPitch)
}

// FormatIdealPitch returns a readable ideal pitch name.
func FormatIdealPitch(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "Invalid") {
		return ""
	}
	return FormatPipelineName(raw)
}

// SchoolClassHref returns the recruiting class URL for a team.
func SchoolClassHref(teamID int) string {
	if teamID <= 0 {
		return ""
	}
	return "/schools/" + fmtInt(teamID) + "/class"
}

// RecruitCommittedSchool returns the committed school for signed or verbal recruits.
// When recruit schoolInterest is empty, target.TopSchool is used as a fallback (common in exports).
func RecruitCommittedSchool(rec *dynasty.RecruitExport, target *dynasty.RecruitingTargetExport) (name string, teamID int, ok bool) {
	if rec == nil {
		return "", 0, false
	}
	stage := strings.ToLower(strings.TrimSpace(rec.RecruitStage))
	if stage != "signed" && stage != "softcommitted" && stage != "hardcommitted" {
		return "", 0, false
	}
	if len(rec.SchoolInterest) > 0 {
		top := rec.SchoolInterest[0]
		for _, school := range rec.SchoolInterest[1:] {
			if school.Influence > top.Influence {
				top = school
			}
		}
		if top.TeamName != "" {
			return top.TeamName, top.TeamID, true
		}
	}
	if target != nil && target.TopSchool != nil && target.TopSchool.TeamName != "" {
		return target.TopSchool.TeamName, target.TopSchool.TeamID, true
	}
	return "", 0, false
}

// RecruitCommitGauge reports top-school progress toward the commit threshold (commitScore).
func RecruitCommitGauge(rec *dynasty.RecruitExport, target *dynasty.RecruitingTargetExport) (barPercent int, label string, ok bool) {
	if rec == nil || rec.CommitScore == nil || *rec.CommitScore <= 0 {
		return 0, "", false
	}
	stage := strings.ToLower(strings.TrimSpace(rec.RecruitStage))
	if stage == "signed" || stage == "softcommitted" || stage == "hardcommitted" {
		return 100, "100%", true
	}
	if target == nil || target.TopSchool == nil || target.TopSchool.Influence <= 0 {
		return 0, "", false
	}
	pct := target.TopSchool.Influence * 100 / *rec.CommitScore
	bar := pct
	if bar > 100 {
		bar = 100
	}
	return bar, fmtInt(pct) + "%", true
}

// RecruitBelongsToSchool reports whether a prospect is signed or verbally committed to teamID.
func RecruitBelongsToSchool(rec *dynasty.RecruitExport, target *dynasty.RecruitingTargetExport, teamID int) bool {
	if rec == nil || teamID <= 0 {
		return false
	}
	_, id, ok := RecruitCommittedSchool(rec, target)
	return ok && id == teamID
}

func formatRawMotivations(mots []string) []string {
	out := make([]string, 0, len(mots))
	for _, raw := range mots {
		raw = strings.TrimSpace(raw)
		if raw == "" || raw == "0" {
			continue
		}
		if strings.EqualFold(raw, "None") || strings.EqualFold(raw, "Invalid") {
			continue
		}
		if _, err := strconv.Atoi(raw); err == nil {
			continue
		}
		if label := FormatSchoolPitchName(raw); label != "" {
			out = append(out, label)
		} else if label := formatPlayerMotivationEnum(raw); label != "" {
			out = append(out, label)
		}
	}
	return out
}

func formatPlayerMotivationEnum(raw string) string {
	key := strings.ToLower(strings.TrimSuffix(raw, "_"))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")
	if label, ok := playerMotivationLabels[key]; ok {
		return label
	}
	return FormatPipelineName(raw)
}

func normalizePitchKey(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	key := strings.ToLower(strings.TrimSuffix(raw, "_"))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")
	return key
}

var schoolPitchLabels = map[string]string{
	"academicprestige":      "Academic Prestige",
	"athleticfacilities":    "Athletic Facilities",
	"brandexposure":         "Brand Exposure",
	"campuslifestyle":       "Campus Lifestyle",
	"championshipcontender": "Championship Contender",
	"coachprestige":         "Coach Prestige",
	"coachstability":        "Coach Stability",
	"conferenceprestige":    "Conference Prestige",
	"playingstyle":          "Playing Style",
	"playingtime":           "Playing Time",
	"propotential":          "Pro Potential",
	"programtradition":      "Program Tradition",
	"proximitytohome":       "Proximity to Home",
	"stadiumatmosphere":     "Stadium Atmosphere",
}

var playerMotivationLabels = map[string]string{
	"closetohome":             "Proximity to Home",
	"championshipcontender":   "Championship Contender",
	"teamprestige":            "Program Tradition",
	"schemefit":               "Playing Style",
	"topthedepthchart":        "Playing Time",
	"headcoachhistoricrecord": "Coach Prestige",
	"bigmarket":               "Brand Exposure",
	"warmweatherstate":        "Campus Lifestyle",
	"noincometax":             "Low Taxes",
	"mentoratposition":        "Coach Prestige",
	"teamhasfranchiseqb":      "Playing Time",
}

// idealPitchMotivations maps ideal pitch enums to their three school pitch motivations.
var idealPitchMotivations = map[string][]string{
	"sundaybound":         {"Championship Contender", "Conference Prestige", "Pro Potential"},
	"hometownhero":        {"Campus Lifestyle", "Proximity to Home", "Program Tradition"},
	"collegeexperience":   {"Academic Prestige", "Campus Lifestyle", "Stadium Atmosphere"},
	"workhorse":           {"Playing Time", "Playing Style", "Athletic Facilities"},
	"grassroots":          {"Proximity to Home", "Program Tradition", "Playing Time"},
	"tothehouse":          {"Championship Contender", "Stadium Atmosphere", "Program Tradition"},
	"prestigious":         {"Academic Prestige", "Championship Contender", "Coach Prestige"},
	"starter":             {"Playing Time", "Pro Potential", "Playing Style"},
	"tvtime":              {"Brand Exposure", "Stadium Atmosphere", "Championship Contender"},
	"conferencespotlight": {"Conference Prestige", "Brand Exposure", "Championship Contender"},
	"coachsfavorite":      {"Coach Prestige", "Coach Stability", "Playing Time"},
	"teamplayer":          {"Playing Style", "Program Tradition", "Coach Stability"},
	"aspirational":        {"Championship Contender", "Brand Exposure", "Pro Potential"},
	"proveyourself":       {"Playing Time", "Pro Potential", "Playing Style"},
	"itsgametime":         {"Stadium Atmosphere", "Brand Exposure", "Championship Contender"},
	"timetogettowork":     {"Playing Time", "Playing Style", "Athletic Facilities"},
	"studentofthegame":    {"Academic Prestige", "Coach Prestige", "Program Tradition"},
	"footballinfluencer":  {"Brand Exposure", "Pro Potential", "Championship Contender"},
	"campuspersonality":   {"Campus Lifestyle", "Academic Prestige", "Stadium Atmosphere"},
	"theclutch":           {"Championship Contender", "Coach Prestige", "Program Tradition"},
}
