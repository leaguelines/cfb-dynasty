package exporter

import (
	"sort"
	"strconv"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// Recruit sort keys for the recruits browser.
const (
	RecruitSortNationalRank = "nationalRank"
	RecruitSortPositionRank = "positionRank"
	RecruitSortStateRank    = "stateRank"
	RecruitSortOverall      = "overall"
	RecruitSortStarRating   = "starRating"
	RecruitSortName         = "name"
	RecruitSortClass        = "class"
	RecruitSortCommitScore  = "commitScore"
	RecruitSortStage        = "recruitStage"
	RecruitSortNIL          = "nil"
	RecruitSortHomeTown     = "homeTown"
)

// RecruitUIColumn describes one recruits table column.
type RecruitUIColumn struct {
	Label   string
	SortKey string
}

// RecruitUIColumns is the prospect list table layout.
var RecruitUIColumns = []RecruitUIColumn{
	{Label: "Rank", SortKey: RecruitSortNationalRank},
	{Label: "School"},
	{Label: "Pos"},
	{Label: "Name", SortKey: RecruitSortName},
	{Label: "Rating", SortKey: RecruitSortStarRating},
	{Label: "Gem"},
	{Label: "NIL", SortKey: RecruitSortNIL},
	{Label: "Stage", SortKey: RecruitSortStage},
	{Label: "Hometown", SortKey: RecruitSortHomeTown},
	{Label: "Pipe"},
	{Label: "Class", SortKey: RecruitSortClass},
}

// RecruitPositions returns distinct recruit positions in tab order.
func RecruitPositions(e dynasty.Export) []string {
	seen := make(map[string]struct{})
	for i := range e.Recruits {
		if pos := recruitPosition(&e.Recruits[i]); pos != "" {
			seen[pos] = struct{}{}
		}
	}
	return sortPositionKeys(seen)
}

// RecruitFilters narrows the prospect list.
type RecruitFilters struct {
	Position string
	State    string
	Stars    string
}

// RecruitHomeStates returns distinct recruit home states sorted alphabetically.
func RecruitHomeStates(e dynasty.Export) []string {
	seen := make(map[string]struct{})
	for i := range e.Recruits {
		if r := e.Recruits[i].Player; r != nil && r.HomeState != "" {
			seen[r.HomeState] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// ValidRecruitState reports whether any recruit has the given home state.
func ValidRecruitState(e dynasty.Export, state string) bool {
	for i := range e.Recruits {
		if r := e.Recruits[i].Player; r != nil && r.HomeState == state {
			return true
		}
	}
	return false
}

// FilterRecruits returns recruits matching the given filters (empty values mean all).
func FilterRecruits(e dynasty.Export, f RecruitFilters) []dynasty.RecruitExport {
	starFilter := 0
	if f.Stars != "" {
		if n, err := strconv.Atoi(f.Stars); err == nil && n >= 1 && n <= 5 {
			starFilter = n
		}
	}
	out := make([]dynasty.RecruitExport, 0)
	for i := range e.Recruits {
		r := e.Recruits[i]
		if f.Position != "" && recruitPosition(&r) != f.Position {
			continue
		}
		if f.State != "" {
			if r.Player == nil || r.Player.HomeState != f.State {
				continue
			}
		}
		if starFilter > 0 {
			if r.Player == nil || ParseStarRating(r.Player.StarRating) != starFilter {
				continue
			}
		}
		out = append(out, r)
	}
	return out
}

// RecruitsForPosition returns recruits at the given position.
func RecruitsForPosition(e dynasty.Export, position string) []dynasty.RecruitExport {
	out := make([]dynasty.RecruitExport, 0)
	for _, r := range e.Recruits {
		if recruitPosition(&r) == position {
			out = append(out, r)
		}
	}
	return out
}

// ValidRecruitPosition reports whether position has recruits.
func ValidRecruitPosition(e dynasty.Export, position string) bool {
	for _, r := range e.Recruits {
		if recruitPosition(&r) == position {
			return true
		}
	}
	return false
}

// RecruitUIHeader returns flat header labels for CSV-style tables.
func RecruitUIHeader() []string {
	h := make([]string, len(RecruitUIColumns))
	for i, c := range RecruitUIColumns {
		h[i] = c.Label
	}
	return h
}

// RecruitListRows builds prospect list rows for the recruits browser.
func RecruitListRows(recruits []dynasty.RecruitExport, teamShortNames map[int]string, targets map[int]dynasty.RecruitingTargetExport) []RecruitListRow {
	rows := make([]RecruitListRow, 0, len(recruits))
	for i := range recruits {
		target := targets[recruits[i].ID]
		rows = append(rows, recruitListRowFromExport(&recruits[i], teamShortNames, target))
	}
	return rows
}

// RecruitUIRows builds flat table rows for downloads and legacy tables.
func RecruitUIRows(recruits []dynasty.RecruitExport, teamShortNames map[int]string, targets map[int]dynasty.RecruitingTargetExport) [][]string {
	rows := make([][]string, 0, len(recruits))
	for i := range recruits {
		row := recruitListRowFromExport(&recruits[i], teamShortNames, targets[recruits[i].ID])
		rows = append(rows, []string{
			row.NationalRank,
			row.SchoolName,
			row.Position,
			row.DisplayName,
			strings.TrimSpace(strings.ReplaceAll(row.StarsText, "☆", "")),
			row.QualityLabel,
			row.NIL,
			row.Stage,
			row.Hometown,
			row.HomePipelineLabel,
			row.Class,
		})
	}
	return rows
}

func recruitListRowFromExport(r *dynasty.RecruitExport, teamShortNames map[int]string, target dynasty.RecruitingTargetExport) RecruitListRow {
	p := r.Player
	row := RecruitListRow{
		NationalRank: fmtIntPtr(r.NationalRank),
		Stage:        FormatRecruitStage(r.RecruitStage),
		Class:        FormatRecruitClass(r.Class),
	}
	if bar, label, ok := RecruitCommitGauge(r, &target); ok {
		row.HasCommitGauge = true
		row.CommitGaugeBar = bar
		row.CommitGaugeLabel = label
	}
	if name, teamID, ok := RecruitCommittedSchool(r, &target); ok {
		row.SchoolName = name
		row.SchoolHref = SchoolClassHref(teamID)
		if teamShortNames != nil {
			if shortName := teamShortNames[teamID]; shortName != "" {
				row.SchoolLogo = TeamLogoPath(shortName)
			}
		}
	}
	if p != nil {
		row.ID = p.ID
		row.Href = "/player/" + fmtInt(p.ID)
		row.Position = p.Position
		row.IsAth = p.IsAth
		row.DisplayName = FormatRecruitDisplayName(p.FirstName, p.LastName)
		row.StarRating, row.StarsText = StarRatingDisplay(p.StarRating)
		row.QualityLabel, row.QualityClass = FormatQualityModifier(r.QualityModifier)
		row.NIL = FormatNILDisplay(p)
		row.HomeState = p.HomeState
		row.Hometown = FormatHometown(p.HomeTown, p.HomeState)
		row.HomePipeline = p.HomePipeline
		row.HomePipelineLabel = FormatPipelineName(p.HomePipeline)
	}
	return row
}

// SortRecruits sorts recruits in place for the given sort key.
func SortRecruits(recruits []dynasty.RecruitExport, key string, desc bool) {
	if key == "" {
		key = RecruitSortNationalRank
	}
	sort.SliceStable(recruits, func(i, j int) bool {
		less := compareRecruits(recruits[i], recruits[j], key)
		if desc {
			return !less
		}
		return less
	})
}

// RecruitSortDesc reports the default descending flag for a sort key.
func RecruitSortDesc(key string) bool {
	switch key {
	case RecruitSortOverall, RecruitSortStarRating, RecruitSortCommitScore, RecruitSortNIL, RecruitSortStage:
		return true
	default:
		return false
	}
}

func recruitPosition(r *dynasty.RecruitExport) string {
	if r == nil || r.Player == nil {
		return ""
	}
	return r.Player.Position
}

func compareRecruits(a, b dynasty.RecruitExport, key string) bool {
	switch key {
	case RecruitSortNationalRank:
		return rankOrLast(a.NationalRank) < rankOrLast(b.NationalRank)
	case RecruitSortPositionRank:
		return rankOrLast(a.PositionRank) < rankOrLast(b.PositionRank)
	case RecruitSortStateRank:
		return rankOrLast(a.StateRank) < rankOrLast(b.StateRank)
	case RecruitSortOverall:
		return playerOverall(a.Player) < playerOverall(b.Player)
	case RecruitSortStarRating:
		return starRatingLess(a.Player, b.Player)
	case RecruitSortName:
		an := recruitName(a)
		bn := recruitName(b)
		if an != bn {
			return an < bn
		}
		return recruitNameLast(a) < recruitNameLast(b)
	case RecruitSortClass:
		if a.Class != b.Class {
			return a.Class < b.Class
		}
		return rankOrLast(a.NationalRank) < rankOrLast(b.NationalRank)
	case RecruitSortCommitScore:
		return intOrLast(a.CommitScore) < intOrLast(b.CommitScore)
	case RecruitSortStage:
		sa, sb := RecruitStageRank(a.RecruitStage), RecruitStageRank(b.RecruitStage)
		if sa != sb {
			return sa < sb
		}
		return rankOrLast(a.NationalRank) < rankOrLast(b.NationalRank)
	case RecruitSortNIL:
		return nilValue(a.Player) < nilValue(b.Player)
	case RecruitSortHomeTown:
		at, bt := recruitHomeTown(a), recruitHomeTown(b)
		if at != bt {
			return at < bt
		}
		return rankOrLast(a.NationalRank) < rankOrLast(b.NationalRank)
	default:
		return rankOrLast(a.NationalRank) < rankOrLast(b.NationalRank)
	}
}

func recruitHomeTown(r dynasty.RecruitExport) string {
	if r.Player == nil {
		return ""
	}
	return strings.ToLower(FormatHometown(r.Player.HomeTown, r.Player.HomeState))
}

func nilValue(p *dynasty.PlayerExport) int {
	if p == nil {
		return 0
	}
	if p.NILBaseValue != nil {
		return *p.NILBaseValue
	}
	if p.NILCompensation != nil {
		return *p.NILCompensation
	}
	return 0
}

func recruitName(r dynasty.RecruitExport) string {
	if r.Player == nil {
		return ""
	}
	return strings.ToLower(r.Player.FirstName)
}

func recruitNameLast(r dynasty.RecruitExport) string {
	if r.Player == nil {
		return ""
	}
	return strings.ToLower(r.Player.LastName)
}

func starRatingLess(a, b *dynasty.PlayerExport) bool {
	return starRatingRank(a) < starRatingRank(b)
}

func starRatingRank(p *dynasty.PlayerExport) int {
	if p == nil {
		return 0
	}
	return ParseStarRating(p.StarRating)
}

func intOrLast(v *int) int {
	if v == nil || *v <= 0 {
		return 1<<31 - 1
	}
	return *v
}
