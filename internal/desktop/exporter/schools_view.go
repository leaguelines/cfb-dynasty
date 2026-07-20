package exporter

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// gradeCell is a compact grade for table columns.
type gradeCell struct {
	Grade string
	Class string
}

// SchoolPitchGrade is one recruiting pitch grade for display.
type SchoolPitchGrade struct {
	Label string
	Grade string
	Class string
}

// SchoolPipelineRow is one pipeline influence row for the school page.
type SchoolPipelineRow struct {
	Name      string
	Tier      int
	TierClass string
	Percent   int
}

// SchoolRivalLink is a rival school link.
type SchoolRivalLink struct {
	Name       string
	ID         int
	SchoolLogo string
}

// SchoolSummary is one school row on the schools index.
type SchoolSummary struct {
	ID                 int
	Name               string
	SchoolLogo         string
	Conference         string
	ConferenceLogo     string
	PrestigeRank       string
	GPA                string
	GPARank            int
	PrestigeHalfStars  int
	PrestigeStars      []PrestigeStar
	HasPrestigeStars   bool
	PrestigeStarsLabel string
	Ratings            TeamRatings
	OffensiveRank      string
	DefensiveRank      string
	CompactGrades      []gradeCell
	Grades             []SchoolPitchGrade
	HasGrades          bool
	Href               string
}

// SchoolDetail is the full school profile page payload.
type SchoolDetail struct {
	SchoolSummary
	Pipelines    []SchoolPipelineRow
	AllPipelines []SchoolPipelineRow
	Rivals       []SchoolRivalLink
	ProPotential []SchoolPitchGrade
	RosterHref   string
	HasRoster    bool
}

var schoolPitchFields = []struct {
	label string
	get   func(g dynasty.SchoolGradesExport) string
}{
	{"Academic Prestige", func(g dynasty.SchoolGradesExport) string { return g.AcademicPrestige }},
	{"Athletic Facilities", func(g dynasty.SchoolGradesExport) string { return g.AthleticFacilities }},
	{"Brand Exposure", func(g dynasty.SchoolGradesExport) string { return g.BrandExposure }},
	{"Campus Lifestyle", func(g dynasty.SchoolGradesExport) string { return g.CampusLifestyle }},
	{"Championship Contender", func(g dynasty.SchoolGradesExport) string { return g.ChampionshipContender }},
	{"Coach Prestige", func(g dynasty.SchoolGradesExport) string { return g.CoachPrestige }},
	{"Coach Stability", func(g dynasty.SchoolGradesExport) string { return g.CoachStability }},
	{"Conference Prestige", func(g dynasty.SchoolGradesExport) string { return g.ConferencePrestige }},
	{"Program Tradition", func(g dynasty.SchoolGradesExport) string { return g.ProgramTradition }},
	{"Stadium Atmosphere", func(g dynasty.SchoolGradesExport) string { return g.StadiumAtmosphere }},
}

var schoolProFields = []struct {
	label string
	get   func(g dynasty.SchoolGradesExport) string
}{
	{"QB", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialQB }},
	{"RB", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialRB }},
	{"WR", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialWR }},
	{"TE", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialTE }},
	{"OL", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialOL }},
	{"DL", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialDL }},
	{"LB", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialLB }},
	{"DB", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialDB }},
	{"K", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialK }},
	{"P", func(g dynasty.SchoolGradesExport) string { return g.ProPotentialP }},
}

// SchoolsIndex returns school summaries for the browser, optionally filtered by query.
func SchoolsIndex(e dynasty.Export, query string) []SchoolSummary {
	summaries := buildSchoolSummaries(e)
	if q := strings.TrimSpace(strings.ToLower(query)); q != "" {
		filtered := make([]SchoolSummary, 0, len(summaries))
		for _, s := range summaries {
			if schoolMatchesQuery(s, q) {
				filtered = append(filtered, s)
			}
		}
		summaries = filtered
	}
	return summaries
}

// SchoolDetailForTeam returns the detail view for a team id.
func SchoolDetailForTeam(e dynasty.Export, teamID int) (SchoolDetail, bool) {
	summaries := buildSchoolSummaries(e)
	var base *SchoolSummary
	for i := range summaries {
		if summaries[i].ID == teamID {
			base = &summaries[i]
			break
		}
	}
	if base == nil {
		return SchoolDetail{}, false
	}

	g, _ := gradesForTeam(e, teamID)
	all := schoolPipelineRows(e, teamID, 0)
	detail := SchoolDetail{
		SchoolSummary: *base,
		AllPipelines:  all,
		Rivals:        schoolRivals(e, teamID),
		ProPotential:  schoolProGrades(g),
	}
	detail.Pipelines = all
	if len(detail.Pipelines) > 3 {
		detail.Pipelines = detail.Pipelines[:3]
	}
	if _, ok := FindRoster(e, teamID); ok {
		detail.HasRoster = true
		detail.RosterHref = "/rosters/" + strconv.Itoa(teamID)
	}
	return detail, true
}

func buildSchoolSummaries(e dynasty.Export) []SchoolSummary {
	gradesByTeam := gradesByTeamID(e)
	gpaRanks := computeGPARanks(e, gradesByTeam)
	ratingsByTeam := TeamRatingsByTeamID(e)

	out := make([]SchoolSummary, 0, len(e.Teams))
	for _, t := range e.Teams {
		name := teamExportName(t)
		if !IsExportableTeamName(name) {
			continue
		}
		g, hasGrades := gradesByTeam[t.ID]
		s := SchoolSummary{
			ID:             t.ID,
			Name:           name,
			SchoolLogo:     TeamLogoForTeam(t),
			Conference:     t.Conference,
			ConferenceLogo: ConferenceLogoPath(t.Conference),
			PrestigeRank:   fmtIntPtr(t.PrestigeRank),
			Ratings:        ratingsByTeam[t.ID],
			OffensiveRank:  fmtIntPtr(t.OffensiveRank),
			DefensiveRank:  fmtIntPtr(t.DefensiveRank),
			Href:           "/schools/" + strconv.Itoa(t.ID),
		}
		if half := prestigeHalfStars(t); half > 0 {
			s.PrestigeHalfStars = half
			s.PrestigeStars = BuildPrestigeStars(half)
			s.HasPrestigeStars = len(s.PrestigeStars) > 0
			s.PrestigeStarsLabel = prestigeStarsLabel(half)
		}
		if hasGrades {
			s.GPA = formatGPA(programGPA(g))
			s.GPARank = gpaRanks[t.ID]
			s.Grades = schoolPitchGrades(g)
			s.HasGrades = len(s.Grades) > 0
		}
		s.CompactGrades = compactPitchGrades(g)
		out = append(out, s)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func teamExportName(t dynasty.TeamExport) string {
	if t.DisplayName != "" {
		return t.DisplayName
	}
	if t.LongName != "" {
		return t.LongName
	}
	return t.ShortName
}

func gradesByTeamID(e dynasty.Export) map[int]dynasty.SchoolGradesExport {
	m := make(map[int]dynasty.SchoolGradesExport, len(e.SchoolGrades))
	for _, g := range e.SchoolGrades {
		m[g.TeamID] = g
	}
	return m
}

func gradesForTeam(e dynasty.Export, teamID int) (dynasty.SchoolGradesExport, bool) {
	g, ok := gradesByTeamID(e)[teamID]
	return g, ok
}

func schoolPitchGrades(g dynasty.SchoolGradesExport) []SchoolPitchGrade {
	out := make([]SchoolPitchGrade, 0, len(schoolPitchFields))
	for _, f := range schoolPitchFields {
		grade := strings.TrimSpace(f.get(g))
		if grade == "" {
			continue
		}
		out = append(out, SchoolPitchGrade{
			Label: f.label,
			Grade: FormatLetterGrade(grade),
			Class: LetterGradeBadgeClass(grade),
		})
	}
	return out
}

func compactPitchGrades(g dynasty.SchoolGradesExport) []gradeCell {
	out := make([]gradeCell, len(schoolPitchFields))
	for i, f := range schoolPitchFields {
		grade := strings.TrimSpace(f.get(g))
		if grade == "" {
			continue
		}
		out[i] = gradeCell{Grade: FormatLetterGrade(grade), Class: LetterGradeBadgeClass(grade)}
	}
	return out
}

func schoolProGrades(g dynasty.SchoolGradesExport) []SchoolPitchGrade {
	out := make([]SchoolPitchGrade, 0, len(schoolProFields))
	for _, f := range schoolProFields {
		grade := strings.TrimSpace(f.get(g))
		if grade == "" {
			continue
		}
		out = append(out, SchoolPitchGrade{
			Label: f.label,
			Grade: FormatLetterGrade(grade),
			Class: LetterGradeBadgeClass(grade),
		})
	}
	return out
}

func pipelinesForTeam(e dynasty.Export, teamID int) []dynasty.PipelineInfluenceExport {
	out := make([]dynasty.PipelineInfluenceExport, 0)
	for _, p := range e.PipelineInfluence {
		if p.TeamID == teamID {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ti, _ := ParsePipelineTier(out[i].InfluenceLevel, out[i].InfluenceValue)
		tj, _ := ParsePipelineTier(out[j].InfluenceLevel, out[j].InfluenceValue)
		if ti != tj {
			return ti > tj
		}
		return out[i].Pipeline < out[j].Pipeline
	})
	return out
}

func schoolPipelineRows(e dynasty.Export, teamID, limit int) []SchoolPipelineRow {
	pipes := pipelinesForTeam(e, teamID)
	if limit > 0 && len(pipes) > limit {
		pipes = pipes[:limit]
	}
	rows := make([]SchoolPipelineRow, len(pipes))
	for i, p := range pipes {
		rows[i] = pipelineRowFromExport(p)
	}
	return rows
}

func pipelineRowFromExport(p dynasty.PipelineInfluenceExport) SchoolPipelineRow {
	tier, _ := ParsePipelineTier(p.InfluenceLevel, p.InfluenceValue)
	return SchoolPipelineRow{
		Name:      FormatPipelineName(p.Pipeline),
		Tier:      tier,
		TierClass: PipelineTierBadgeClass(tier),
		Percent:   PipelineTierPercent(tier),
	}
}

func schoolRivals(e dynasty.Export, teamID int) []SchoolRivalLink {
	seen := map[int]struct{}{}
	shortNames := BuildTeamShortNameMap(e.Teams)
	var out []SchoolRivalLink
	for _, r := range e.Rivalries {
		if r.Team1ID != nil && *r.Team1ID == teamID && r.Team2ID != nil {
			if rival := rivalLink(*r.Team2ID, r.Team2Name, seen, shortNames); rival != nil {
				out = append(out, *rival)
			}
		}
		if r.Team2ID != nil && *r.Team2ID == teamID && r.Team1ID != nil {
			if rival := rivalLink(*r.Team1ID, r.Team1Name, seen, shortNames); rival != nil {
				out = append(out, *rival)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func rivalLink(id int, name string, seen map[int]struct{}, shortNames map[int]string) *SchoolRivalLink {
	if id < 0 || name == "" {
		return nil
	}
	if !IsExportableTeamName(name) {
		return nil
	}
	if _, ok := seen[id]; ok {
		return nil
	}
	seen[id] = struct{}{}
	return &SchoolRivalLink{Name: name, ID: id, SchoolLogo: TeamLogoForTeamID(shortNames, id)}
}

func schoolMatchesQuery(s SchoolSummary, q string) bool {
	if strings.Contains(strings.ToLower(s.Name), q) {
		return true
	}
	return strings.Contains(strings.ToLower(s.Conference), q)
}

func computeGPARanks(e dynasty.Export, grades map[int]dynasty.SchoolGradesExport) map[int]int {
	type row struct {
		id  int
		gpa float64
	}
	rows := make([]row, 0, len(grades))
	for id, g := range grades {
		if gpa := programGPA(g); gpa > 0 {
			rows = append(rows, row{id: id, gpa: gpa})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].gpa != rows[j].gpa {
			return rows[i].gpa > rows[j].gpa
		}
		return rows[i].id < rows[j].id
	})
	ranks := make(map[int]int, len(rows))
	for i, r := range rows {
		ranks[r.id] = i + 1
	}
	return ranks
}

func programGPA(g dynasty.SchoolGradesExport) float64 {
	var sum float64
	var n int
	for _, f := range schoolPitchFields {
		if v, ok := gradePoints(f.get(g)); ok {
			sum += v
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

func formatGPA(gpa float64) string {
	if gpa <= 0 {
		return ""
	}
	return strconv.FormatFloat(math.Round(gpa*100)/100, 'f', 2, 64)
}

func gradePoints(raw string) (float64, bool) {
	switch FormatLetterGrade(raw) {
	case "A+":
		return 4.3, true
	case "A":
		return 4.0, true
	case "A-":
		return 3.7, true
	case "B+":
		return 3.3, true
	case "B":
		return 3.0, true
	case "B-":
		return 2.7, true
	case "C+":
		return 2.3, true
	case "C":
		return 2.0, true
	case "C-":
		return 1.7, true
	case "D+":
		return 1.3, true
	case "D":
		return 1.0, true
	case "F", "F+", "F-":
		return 0, true
	default:
		return 0, false
	}
}
