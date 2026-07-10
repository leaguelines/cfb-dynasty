package dynasty

import "sync"

type recruitingClassIndex struct {
	playerRecruits map[int]Record
	rankWeights    []int
}

var (
	recruitingClassIndexMu    sync.Mutex
	recruitingClassIndexByKey string
	recruitingClassIndexCache recruitingClassIndex
)

func (f *File) recruitingClassIndex() recruitingClassIndex {
	if f == nil {
		return recruitingClassIndex{}
	}
	key := f.path + "|" + f.settings.SchemaDir + "|" + discoverTuningPath(f.settings.SchemaDir)
	recruitingClassIndexMu.Lock()
	if key == recruitingClassIndexByKey && len(recruitingClassIndexCache.playerRecruits) > 0 {
		idx := recruitingClassIndexCache
		recruitingClassIndexMu.Unlock()
		return idx
	}
	recruitingClassIndexMu.Unlock()

	idx := buildRecruitingClassIndex(f)

	recruitingClassIndexMu.Lock()
	recruitingClassIndexByKey = key
	recruitingClassIndexCache = idx
	recruitingClassIndexMu.Unlock()
	return idx
}

func buildRecruitingClassIndex(f *File) recruitingClassIndex {
	idx := recruitingClassIndex{
		playerRecruits: make(map[int]Record),
	}
	recruitTable, ok := f.PrimaryTableByName("Recruit")
	if !ok {
		return idx
	}
	if err := recruitTable.ReadRecords(); err != nil {
		return idx
	}
	for _, record := range recruitTable.Records {
		playerRef, ok := record.Get("Player")
		if !ok || playerRef.Reference == nil {
			continue
		}
		idx.playerRecruits[int(playerRef.Reference.RowNumber)] = record
	}
	idx.rankWeights = topClassRankWeights(f)
	return idx
}

func topClassRankWeights(f *File) []int {
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
	return intArrayFromRecord(tf, table.Records[0], "TopClassesRankWeightPercentageTable")
}

func recruitingClassExportForTeam(f *File, team Record, idx recruitingClassIndex) *TeamRecruitingClassExport {
	national := topClassNationalRank(team)
	conference := topClassConferenceRank(team)
	score, commits := recruitingClassScoreForTeam(f, team, idx)

	if national == nil && conference == nil && score == 0 && commits == 0 {
		return nil
	}
	out := &TeamRecruitingClassExport{
		NationalRank:   national,
		ConferenceRank: conference,
	}
	if score > 0 {
		out.Score = score
	}
	if commits > 0 {
		out.CommitCount = commits
	}
	return out
}

func recruitingClassScoreForTeam(f *File, team Record, idx recruitingClassIndex) (score int, commits int) {
	if len(idx.playerRecruits) == 0 {
		return 0, 0
	}
	ref, ok := team.Get("CommittedPlayers")
	if !ok || ref.Reference == nil {
		return 0, 0
	}
	minWeight := topClassMinimumRankWeight(idx.rankWeights)
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		player, ok := f.RecordByReference("Player", memberRef)
		if !ok {
			continue
		}
		recruit, ok := idx.playerRecruits[player.Index]
		if !ok {
			continue
		}
		commitScore, ok := intFieldOK(recruit, "CommitScore")
		if !ok {
			continue
		}
		weight := topClassRankWeight(idx.rankWeights, recruit, minWeight)
		score += commitScore * weight / 100
		commits++
	}
	return score, commits
}

func topClassMinimumRankWeight(weights []int) int {
	if len(weights) == 0 {
		return 100
	}
	min := weights[len(weights)-1]
	for _, weight := range weights {
		if weight < min {
			min = weight
		}
	}
	if min <= 0 {
		return 3
	}
	return min
}

func topClassRankWeight(weights []int, recruit Record, fallback int) int {
	if len(weights) == 0 {
		return 100
	}
	nationalRank, ok := intFieldOK(recruit, "NationalRank")
	if !ok || nationalRank <= 0 {
		return fallback
	}
	if nationalRank > len(weights) {
		return fallback
	}
	weight := weights[nationalRank-1]
	if weight <= 0 {
		return fallback
	}
	return weight
}

func topClassNationalRank(record Record) *int {
	value, ok := intFieldOK(record, "TopClassRank")
	if !ok || value <= 0 || value > 250 {
		return nil
	}
	return &value
}

func topClassConferenceRank(record Record) *int {
	value, ok := intFieldOK(record, "TopClassConferenceRank")
	if !ok || value <= 0 || value > 31 {
		return nil
	}
	return &value
}
