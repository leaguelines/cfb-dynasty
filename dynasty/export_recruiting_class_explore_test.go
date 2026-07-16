package dynasty

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// classScoreVariant names alternative scoring formulas explored against the save.
type classScoreVariant struct {
	name string
	fn   func(f *File, team Record, idx recruitingClassIndex, teamID int) int
}

// teamClassRow holds per-team scoring comparison data for exploration tests.
type teamClassRow struct {
	name            string
	nationalRank    int
	scores          map[string]int
	hsCommits       int
	transferCommits int
}

func TestExploreRecruitingClassRanking(t *testing.T) {
	skipIfShortIntegration(t)
	f := openSeasonSave(t)
	idx := buildRecruitingClassIndex(f)
	starTable := topClassStarTable(f)

	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		t.Fatal("Team table missing")
	}
	if err := teamTable.ReadRecords(); err != nil {
		t.Fatal(err)
	}

	variants := []classScoreVariant{
		{name: "weighted_commit", fn: func(f *File, team Record, idx recruitingClassIndex, teamID int) int {
			score, _ := recruitingClassScoreForTeam(f, team, idx)
			return score
		}},
		{name: "raw_commit_sum", fn: classScoreRawCommitSum},
		{name: "hs_only_weighted", fn: classScoreHSOnlyWeighted},
		{name: "influence_share", fn: classScoreInfluenceShare},
		{name: "weighted_plus_stars", fn: func(f *File, team Record, idx recruitingClassIndex, teamID int) int {
			return classScoreWeightedPlusStars(f, team, idx, teamID, starTable)
		}},
		{name: "signed_100pct", fn: classScoreSigned100PctWeighted},
		{name: "star_table_sum", fn: func(f *File, team Record, idx recruitingClassIndex, _ int) int {
			return classScoreStarTableSum(f, team, idx, starTable)
		}},
		{name: "weighted_star_table", fn: func(f *File, team Record, idx recruitingClassIndex, _ int) int {
			return classScoreWeightedStarTable(f, team, idx, starTable)
		}},
		{name: "top100_count", fn: classScoreTop100Count},
		{name: "hs_raw_sum", fn: classScoreHSRawSum},
	}

	type teamRow = teamClassRow
	var teams []teamRow
	for _, record := range teamTable.Records {
		if !recordIsActive(record, teamTable) {
			continue
		}
		name := stringField(record, "LongName")
		if !isOfficialTeamName(name) {
			continue
		}
		rank := intField(record, "TopClassRank")
		if rank <= 0 {
			continue
		}
		teamID := intField(record, "TeamIndex")
		hs, xfer := countCommitClasses(f, record, idx)
		row := teamRow{name: name, nationalRank: rank, hsCommits: hs, transferCommits: xfer, scores: make(map[string]int)}
		for _, v := range variants {
			row.scores[v.name] = v.fn(f, record, idx, teamID)
		}
		teams = append(teams, row)
	}

	// How well does each variant's score order match stored national rank?
	for _, v := range variants {
		mismatch := rankOrderMismatch(teams, v.name)
		t.Logf("variant %s: %d rank inversions vs score sort (of %d ranked teams)", v.name, mismatch, len(teams))
	}
	top15Match := func(variant string) int {
		byScore := append([]teamClassRow(nil), teams...)
		sort.Slice(byScore, func(i, j int) bool { return byScore[i].scores[variant] > byScore[j].scores[variant] })
		match := 0
		for i := 0; i < 15 && i < len(byScore); i++ {
			if byScore[i].nationalRank == i+1 {
				match++
			}
		}
		return match
	}
	for _, name := range []string{"weighted_commit", "raw_commit_sum", "hs_raw_sum", "star_table_sum", "weighted_star_table", "five_star_count"} {
		t.Logf("top-15 exact rank match (%s): %d/15", name, top15Match(name))
	}

	// Focus on top-15 stored ranks where score order diverges.
	sort.Slice(teams, func(i, j int) bool { return teams[i].nationalRank < teams[j].nationalRank })
	t.Log("\nTop 15 by stored national rank:")
	for i := 0; i < 15 && i < len(teams); i++ {
		team := teams[i]
		t.Logf("  #%d %-14s weighted=%4d wStars=%3d stars=%3d top100=%2d 5*=%d",
			team.nationalRank, team.name,
			team.scores["weighted_commit"], team.scores["weighted_star_table"],
			team.scores["star_table_sum"], team.scores["top100_count"], team.scores["five_star_count"])
	}

	// Show score-sorted top 15 for comparison.
	byWeighted := append([]teamRow(nil), teams...)
	sort.Slice(byWeighted, func(i, j int) bool {
		return byWeighted[i].scores["weighted_commit"] > byWeighted[j].scores["weighted_commit"]
	})
	t.Log("\nTop 15 by weighted_commit score:")
	for i := 0; i < 15 && i < len(byWeighted); i++ {
		team := byWeighted[i]
		t.Logf("  score=%4d stored=#%-2d %-14s hs/xfer=%d/%d",
			team.scores["weighted_commit"], team.nationalRank, team.name,
			team.hsCommits, team.transferCommits)
	}

	// Deep dive: Oregon (#4) vs Alabama (#5) vs Miami (#6)
	for _, target := range []string{"Oregon", "Alabama", "Miami", "Indiana", "Texas"} {
		t.Logf("\n--- %s commit breakdown ---", target)
		logTeamCommitBreakdown(t, f, teamTable, idx, target)
	}
}

func rankOrderMismatch(teams []teamClassRow, variant string) int {
	sorted := append([]teamClassRow(nil), teams...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].scores[variant] > sorted[j].scores[variant]
	})
	inversions := 0
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			// if higher stored rank (worse) has higher score than lower stored rank (better)
			if sorted[i].nationalRank > sorted[j].nationalRank &&
				sorted[i].scores[variant] > sorted[j].scores[variant] {
				inversions++
			}
		}
	}
	return inversions
}

func classScoreRawCommitSum(f *File, team Record, idx recruitingClassIndex, _ int) int {
	score := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		if cs, ok := intFieldOK(recruit, "CommitScore"); ok {
			score += cs
		}
	}
	return score
}

func classScoreHSOnlyWeighted(f *File, team Record, idx recruitingClassIndex, _ int) int {
	minWeight := topClassMinimumRankWeight(idx.rankWeights)
	score := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		if !isHighSchoolRecruit(recruit) {
			continue
		}
		commitScore, ok := intFieldOK(recruit, "CommitScore")
		if !ok {
			continue
		}
		weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
		score += commitScore * weight / 100
	}
	return score
}

func classScoreInfluenceShare(f *File, team Record, idx recruitingClassIndex, teamID int) int {
	minWeight := topClassMinimumRankWeight(idx.rankWeights)
	score := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		commitScore, ok := intFieldOK(recruit, "CommitScore")
		if !ok {
			continue
		}
		share := schoolInfluenceShare(f, recruit, teamID)
		if share <= 0 {
			share = 100
		}
		weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
		score += commitScore * share / 100 * weight / 100
	}
	return score
}

func classScoreWeightedPlusStars(f *File, team Record, idx recruitingClassIndex, teamID int, starTable []int) int {
	base := classScoreInfluenceShare(f, team, idx, teamID)
	if len(starTable) == 0 {
		return base
	}
	for _, recruit := range committedRecruits(f, team, idx) {
		stars := recruitStarLevel(recruit)
		if stars <= 0 || stars > len(starTable) {
			continue
		}
		base += starTable[stars-1]
	}
	return base
}

func classScoreSigned100PctWeighted(f *File, team Record, idx recruitingClassIndex, _ int) int {
	minWeight := topClassMinimumRankWeight(idx.rankWeights)
	score := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		commitScore, ok := intFieldOK(recruit, "CommitScore")
		if !ok {
			continue
		}
		share := 100
		if stringField(recruit, "RecruitStage") != "Signed" {
			share = 100 // committed players should all be signed; keep 100
		}
		weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
		score += commitScore * share / 100 * weight / 100
	}
	return score
}

func classScoreStarTableSum(f *File, team Record, idx recruitingClassIndex, starTable []int) int {
	if len(starTable) == 0 {
		return 0
	}
	total := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		stars := recruitProspectStars(f, recruit)
		if stars <= 0 || stars > len(starTable) {
			continue
		}
		total += starTable[stars-1]
	}
	return total
}

func classScoreWeightedStarTable(f *File, team Record, idx recruitingClassIndex, starTable []int) int {
	if len(starTable) == 0 {
		return 0
	}
	minWeight := topClassMinimumRankWeight(idx.rankWeights)
	total := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		stars := recruitProspectStars(f, recruit)
		if stars <= 0 || stars > len(starTable) {
			continue
		}
		weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
		total += starTable[stars-1] * weight / 100
	}
	return total
}

func classScoreTop100Count(f *File, team Record, idx recruitingClassIndex, _ int) int {
	count := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		rank, ok := intFieldOK(recruit, "NationalRank")
		if ok && rank > 0 && rank <= 100 {
			count++
		}
	}
	return count
}

func recruitProspectStars(f *File, recruit Record) int {
	playerRef, ok := recruit.Get("Player")
	if !ok || playerRef.Reference == nil {
		return 0
	}
	player, ok := f.RecordByReference("Player", playerRef.Reference)
	if !ok {
		return 0
	}
	star := strings.ToUpper(stringField(player, "ProspectStarRating"))
	switch {
	case strings.Contains(star, "FIVE"), star == "5":
		return 5
	case strings.Contains(star, "FOUR"), star == "4":
		return 4
	case strings.Contains(star, "THREE"), star == "3":
		return 3
	case strings.Contains(star, "TWO"), star == "2":
		return 2
	case strings.Contains(star, "ONE"), star == "1":
		return 1
	}
	return 0
}

func classScoreHSRawSum(f *File, team Record, idx recruitingClassIndex, _ int) int {
	score := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		if !isHighSchoolRecruit(recruit) {
			continue
		}
		if cs, ok := intFieldOK(recruit, "CommitScore"); ok {
			score += cs
		}
	}
	return score
}

func classScoreFiveStarCount(f *File, team Record, idx recruitingClassIndex, _ int) int {
	count := 0
	for _, recruit := range committedRecruits(f, team, idx) {
		if recruitProspectStars(f, recruit) >= 5 {
			count++
		}
	}
	return count
}

func committedRecruits(f *File, team Record, idx recruitingClassIndex) []Record {
	ref, ok := team.Get("CommittedPlayers")
	if !ok || ref.Reference == nil {
		return nil
	}
	var recruits []Record
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		player, ok := f.RecordByReference("Player", memberRef)
		if !ok {
			continue
		}
		recruit, ok := idx.playerRecruits[player.Index]
		if !ok {
			continue
		}
		recruits = append(recruits, recruit)
	}
	return recruits
}

func isHighSchoolRecruit(recruit Record) bool {
	class := strings.ToLower(stringField(recruit, "Class"))
	return strings.Contains(class, "highschool") || class == "0"
}

func schoolInfluenceShare(f *File, recruit Record, teamID int) int {
	ref, ok := recruit.Get("TopSchoolsList")
	if !ok || ref.Reference == nil {
		return 100
	}
	total := 0
	teamInfl := 0
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("ProspectTargetSchool", memberRef)
		if !ok {
			continue
		}
		infl, ok := intFieldOK(row, "TeamInfluence")
		if !ok {
			continue
		}
		total += infl
		if intField(row, "TeamId") == teamID {
			teamInfl = infl
		}
	}
	if total <= 0 {
		return 100
	}
	return teamInfl * 100 / total
}

func recruitStarLevel(recruit Record) int {
	// ProductionGrade often encodes star tier in recruiting saves.
	if grade, ok := intFieldOK(recruit, "ProductionGrade"); ok && grade > 0 && grade <= 5 {
		return grade
	}
	return 0
}

func countCommitClasses(f *File, team Record, idx recruitingClassIndex) (hs, xfer int) {
	for _, recruit := range committedRecruits(f, team, idx) {
		if isHighSchoolRecruit(recruit) {
			hs++
		} else {
			xfer++
		}
	}
	return hs, xfer
}

func topClassStarTable(f *File) []int {
	path := f.settings.TuningPath
	if path == "" {
		path = discoverTuningPath(f.settings.SchemaDir)
	}
	if path == "" {
		return nil
	}
	settings := DefaultSettings()
	settings.SchemaDir = f.settings.SchemaDir
	settings.TuningPath = path
	settings.AutoParse = true
	tf, err := Open(path, &settings)
	if err != nil {
		return nil
	}
	table, ok := tf.PrimaryTableByName("RecruitingTunables")
	if !ok {
		return nil
	}
	if err := table.ReadRecords(); err != nil || len(table.Records) == 0 {
		return nil
	}
	return intArrayFromRecord(tf, table.Records[0], "TopClassesStarRatingTable")
}

func logTeamCommitBreakdown(t *testing.T, f *File, teamTable *Table, idx recruitingClassIndex, school string) {
	t.Helper()
	for _, record := range teamTable.Records {
		if stringField(record, "LongName") != school {
			continue
		}
		teamID := intField(record, "TeamIndex")
		minWeight := topClassMinimumRankWeight(idx.rankWeights)
		type line struct {
			rank   int
			score  int
			commit int
			share  int
			class  string
			stage  string
			stars  int
			weighted int
		}
		var lines []line
		for _, recruit := range committedRecruits(f, record, idx) {
			rank, _ := intFieldOK(recruit, "NationalRank")
			commit, _ := intFieldOK(recruit, "CommitScore")
			share := schoolInfluenceShare(f, recruit, teamID)
			weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
			lines = append(lines, line{
				rank: rank, commit: commit, share: share,
				class: stringField(recruit, "Class"),
				stage: stringField(recruit, "RecruitStage"),
				stars: recruitStarLevel(recruit),
				weighted: commit * share / 100 * weight / 100,
			})
		}
		sort.Slice(lines, func(i, j int) bool { return lines[i].rank < lines[j].rank })
		for _, l := range lines {
			t.Logf("  rank=%3d commit=%4d share=%3d%% weight=%d weighted=%d class=%s stage=%s stars=%d",
				l.rank, l.commit, l.share, topClassRankWeight(idx.rankWeights, Record{Fields: map[string]FieldValue{"NationalRank": {Int: int64(l.rank)}}}, minWeight),
				l.weighted, l.class, l.stage, l.stars)
		}
		return
	}
	t.Log("team not found:", school)
}

// silence unused import if build tags change
var _ = fmt.Sprint
