package exporter

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// RecruitingRankingRow is one team on the recruiting rankings board.
type RecruitingRankingRow struct {
	Rank        int
	TeamID      int
	TeamName    string
	Conference  string
	SchoolLogo  string
	ClassHref   string
	Commits     int
	FiveStar    int
	FourStar    int
	ThreeStar   int
	AvgOverall  string
	TotalPoints string
}

// RecruitingRankingsView is the team recruiting rankings page.
type RecruitingRankingsView struct {
	ClassYear      int
	ClassLabel     string
	Conference     string
	ConferenceLogo string
	Conferences    []ConferenceFilter
	Rows           []RecruitingRankingRow
	Total          int
}

// ConferenceFilter is one conference pill on the rankings page.
type ConferenceFilter struct {
	Name string
	Href string
	Logo string
}

// RecruitingClassYear estimates the HS recruiting class year from the dynasty season.
func RecruitingClassYear(e dynasty.Export) int {
	if e.Season != nil && e.Season.Year > 0 {
		return e.Season.Year + 1
	}
	return 0
}

// RecruitingRankings builds team recruiting rankings from game-exported class data.
func RecruitingRankings(e dynasty.Export, conference string) RecruitingRankingsView {
	targets := BuildRecruitTargetByRecruitID(e)
	shortNames := BuildTeamShortNameMap(e.Teams)
	recruitStats := recruitCommitStatsByTeam(e, targets)

	confFilter := strings.TrimSpace(conference)
	type sortRow struct {
		teamID      int
		teamName    string
		conference  string
		class       *dynasty.TeamRecruitingClassExport
		stats       recruitCommitStats
		displayRank int
		sortRank    int
		sortScore   int
	}
	const unranked = 1_000_000

	rows := make([]sortRow, 0, len(e.Teams))
	for _, team := range e.Teams {
		name := teamExportName(team)
		if !IsExportableTeamName(name) {
			continue
		}
		conf := strings.TrimSpace(team.Conference)
		if confFilter != "" && !strings.EqualFold(conf, confFilter) {
			continue
		}

		stats := recruitStats[team.ID]
		rc := team.RecruitingClass
		if rc == nil && stats.commits == 0 {
			continue
		}

		sortRank := unranked
		displayRank := 0
		if confFilter != "" {
			if rc != nil && rc.ConferenceRank != nil {
				sortRank = *rc.ConferenceRank
				displayRank = *rc.ConferenceRank
			}
		} else if rc != nil && rc.NationalRank != nil {
			sortRank = *rc.NationalRank
			displayRank = *rc.NationalRank
		}

		sortScore := 0
		if rc != nil {
			sortScore = rc.Score
		}

		rows = append(rows, sortRow{
			teamID:      team.ID,
			teamName:    name,
			conference:  conf,
			class:       rc,
			stats:       stats,
			displayRank: displayRank,
			sortRank:    sortRank,
			sortScore:   sortScore,
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].sortRank != rows[j].sortRank {
			if rows[i].sortRank == unranked {
				return false
			}
			if rows[j].sortRank == unranked {
				return true
			}
			return rows[i].sortRank < rows[j].sortRank
		}
		if rows[i].sortScore != rows[j].sortScore {
			return rows[i].sortScore > rows[j].sortScore
		}
		if rows[i].stats.avgOverall != rows[j].stats.avgOverall {
			return rows[i].stats.avgOverall > rows[j].stats.avgOverall
		}
		return rows[i].stats.commits > rows[j].stats.commits
	})

	out := make([]RecruitingRankingRow, len(rows))
	for i, row := range rows {
		rank := row.displayRank
		if rank <= 0 {
			rank = i + 1
		}

		commits := row.stats.commits
		if row.class != nil && row.class.CommitCount > 0 {
			commits = row.class.CommitCount
		}

		score := 0
		if row.class != nil {
			score = row.class.Score
		}

		logo := ""
		if sn := shortNames[row.teamID]; sn != "" {
			logo = TeamLogoPath(sn)
		}

		out[i] = RecruitingRankingRow{
			Rank:        rank,
			TeamID:      row.teamID,
			TeamName:    row.teamName,
			Conference:  row.conference,
			SchoolLogo:  logo,
			ClassHref:   SchoolClassHref(row.teamID),
			Commits:     commits,
			FiveStar:    row.stats.five,
			FourStar:    row.stats.four,
			ThreeStar:   row.stats.three,
			AvgOverall:  formatRecruitingAvg(row.stats.avgOverall),
			TotalPoints: formatRecruitingScore(score),
		}
	}

	classYear := RecruitingClassYear(e)
	label := ""
	if classYear > 0 {
		label = fmtInt(classYear)
	}

	return RecruitingRankingsView{
		ClassYear:      classYear,
		ClassLabel:     label,
		Conference:     confFilter,
		ConferenceLogo: ConferenceLogoPath(confFilter),
		Conferences:    recruitingRankingConferenceFilters(e.Teams),
		Rows:           out,
		Total:          len(out),
	}
}

type recruitCommitStats struct {
	commits    int
	five       int
	four       int
	three      int
	avgOverall int // tenths: 65.4 -> 654
}

func recruitCommitStatsByTeam(e dynasty.Export, targets map[int]dynasty.RecruitingTargetExport) map[int]recruitCommitStats {
	teamMeta := teamMetaByID(e.Teams)
	out := make(map[int]recruitCommitStats)

	for i := range e.Recruits {
		r := &e.Recruits[i]
		target := targets[r.ID]
		_, teamID, ok := RecruitCommittedSchool(r, &target)
		if !ok || teamID <= 0 {
			continue
		}
		if _, ok := teamMeta[teamID]; !ok {
			continue
		}

		stats := out[teamID]
		stats.commits++

		stars := 0
		if r.Player != nil {
			stars = ParseStarRating(r.Player.StarRating)
		}
		switch stars {
		case 5:
			stats.five++
		case 4:
			stats.four++
		case 3:
			stats.three++
		}

		if overall := playerOverall(r.Player); overall > 0 {
			stats.avgOverall = (stats.avgOverall*(stats.commits-1) + overall*10) / stats.commits
		}

		out[teamID] = stats
	}
	return out
}

type teamMeta struct {
	name       string
	conference string
}

func teamMetaByID(teams []dynasty.TeamExport) map[int]teamMeta {
	m := make(map[int]teamMeta, len(teams))
	for _, t := range teams {
		name := teamExportName(t)
		if !IsExportableTeamName(name) {
			continue
		}
		m[t.ID] = teamMeta{name: name, conference: strings.TrimSpace(t.Conference)}
	}
	return m
}

func recruitingRankingConferenceFilters(teams []dynasty.TeamExport) []ConferenceFilter {
	names := recruitingRankingConferences(teams)
	out := make([]ConferenceFilter, len(names))
	for i, name := range names {
		out[i] = ConferenceFilter{
			Name: name,
			Href: "/recruiting/rankings?conf=" + url.QueryEscape(name),
			Logo: ConferenceLogoPath(name),
		}
	}
	return out
}

func recruitingRankingConferences(teams []dynasty.TeamExport) []string {
	seen := make(map[string]struct{})
	for _, t := range teams {
		name := teamExportName(t)
		if !IsExportableTeamName(name) {
			continue
		}
		conf := strings.TrimSpace(t.Conference)
		if conf == "" {
			continue
		}
		seen[conf] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for conf := range seen {
		out = append(out, conf)
	}
	sort.Strings(out)
	return out
}

func formatRecruitingAvg(avgTenths int) string {
	if avgTenths <= 0 {
		return "—"
	}
	return fmt.Sprintf("%.1f", float64(avgTenths)/10)
}

func formatRecruitingScore(score int) string {
	if score <= 0 {
		return "—"
	}
	return fmtInt(score)
}
