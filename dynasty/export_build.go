package dynasty

import "fmt"

// buildExport assembles the normalized export payload from parsed tables.
func (f *File) buildExport(opts ExportOptions) (Export, error) {
	if !f.loaded {
		return Export{}, ErrNotLoaded
	}
	if f.schema == nil {
		return Export{}, ErrSchemaRequired
	}

	export := Export{
		Source: ExportSource{
			Path:          f.path,
			Size:          f.info.Size,
			SHA256:        f.info.SHA256,
			GameYear:      f.info.GameYear,
			Format:        f.info.Format,
			Compressed:    f.info.Compressed,
			SchemaVersion: schemaVersionCopy(f.schema.Version),
		},
	}

	if opts.IncludeTeams() {
		teams, err := f.buildTeamExports()
		if err != nil {
			return Export{}, err
		}
		export.Teams = teams
	}

	if opts.IncludeRosters() {
		rosters, err := f.buildRosterExports()
		if err != nil {
			return Export{}, err
		}
		export.Rosters = rosters
	}

	if opts.IncludeGames() {
		games, err := f.buildGameExports(opts)
		if err != nil {
			return Export{}, err
		}
		export.Games = games
	}

	if opts.IncludeRecruits() {
		recruits, err := f.buildRecruitExports()
		if err != nil {
			return Export{}, err
		}
		export.Recruits = recruits
	}

	if opts.IncludeRecruiting() {
		recruiting, err := f.buildRecruitingExports()
		if err != nil {
			return Export{}, err
		}
		export.Recruiting = recruiting
	}

	if opts.IncludeSeasonStats() {
		playerStats, err := f.buildPlayerSeasonStatsExports()
		if err != nil {
			return Export{}, err
		}
		export.SeasonPlayerStats = playerStats

		teamStats, err := f.buildTeamSeasonStatsExports()
		if err != nil {
			return Export{}, err
		}
		export.SeasonTeamStats = teamStats
	}

	if opts.IncludeCoaches() {
		coaches, err := f.buildCoachExports()
		if err != nil {
			return Export{}, err
		}
		export.Coaches = coaches
	}

	if opts.IncludeLeavingPlayers() {
		leaving, err := f.buildLeavingPlayerExports()
		if err != nil {
			return Export{}, err
		}
		export.LeavingPlayers = leaving
	}

	if opts.IncludeInjuries() {
		injuries, err := f.buildInjuryExports()
		if err != nil {
			return Export{}, err
		}
		export.Injuries = injuries
	}

	if opts.IncludeDepthCharts() {
		depthCharts, err := f.buildDepthChartExports()
		if err != nil {
			return Export{}, err
		}
		export.DepthCharts = depthCharts
	}

	if opts.IncludeHistory() {
		playerAwards, err := f.buildPlayerAwardExports()
		if err != nil {
			return Export{}, err
		}
		export.PlayerAwards = playerAwards

		leagueAwards, err := f.buildLeagueHistoryAwardExports()
		if err != nil {
			return Export{}, err
		}
		export.LeagueAwards = leagueAwards

		champions, err := f.buildConferenceChampionExports()
		if err != nil {
			return Export{}, err
		}
		export.ConferenceChampions = champions

		recordBook, err := f.buildRecordBookExports()
		if err != nil {
			return Export{}, err
		}
		export.RecordBook = recordBook
	}

	if opts.IncludeSeason() {
		export.Season = f.buildSeasonExport()
	}

	if opts.IncludeSchoolGrades() {
		grades, err := f.buildSchoolGradesExports()
		if err != nil {
			return Export{}, err
		}
		export.SchoolGrades = grades
	}

	if opts.IncludePipelines() {
		pipelines, err := f.buildPipelineInfluenceExports()
		if err != nil {
			return Export{}, err
		}
		export.PipelineInfluence = pipelines
	}

	if opts.IncludeRivalries() {
		rivalries, err := f.buildRivalryExports()
		if err != nil {
			return Export{}, err
		}
		export.Rivalries = rivalries
	}

	if opts.IncludePositionChanges() {
		changes, err := f.buildPositionChangeExports()
		if err != nil {
			return Export{}, err
		}
		export.PositionChanges = changes
	}

	if opts.IncludeDraft() {
		picks, err := f.buildDraftPickExports()
		if err != nil {
			return Export{}, err
		}
		export.DraftPicks = picks
	}

	if opts.IncludeBowls() {
		bowls, err := f.buildBowlGameExports()
		if err != nil {
			return Export{}, err
		}
		export.BowlGames = bowls
	}

	return export, nil
}

func (f *File) buildGameExports(opts ExportOptions) ([]GameExport, error) {
	seasonGame, ok := f.PrimaryTableByName("SeasonGame")
	if !ok {
		return nil, fmt.Errorf("cfb-dynasty: export: SeasonGame table not found")
	}
	if err := seasonGame.ReadRecords(); err != nil {
		return nil, err
	}

	teamTable, ok := f.PrimaryTableByName("Team")
	if !ok {
		return nil, fmt.Errorf("cfb-dynasty: export: Team table not found")
	}
	if err := teamTable.ReadRecords(); err != nil {
		return nil, err
	}

	teamNames := buildTeamNameIndex(teamTable)

	var statsIndex gameStatsIndex
	var detailIndex gameDetailIndex
	var err error
	if opts.IncludeGameStats() {
		statsIndex, err = f.buildGameStatsIndex(len(seasonGame.Records))
		if err != nil {
			return nil, err
		}
	}
	detailIndex = f.buildGameDetailIndex(len(seasonGame.Records))

	games := make([]GameExport, 0, seasonGame.ActiveRecordCount())
	for _, record := range seasonGame.Records {
		if isEmptyReference(record, "HomeTeam") && isEmptyReference(record, "AwayTeam") {
			continue
		}
		game := GameExport{
			ID:         uint32(record.Index),
			SeasonYear: intField(record, "SeasonYear"),
			Week:       intField(record, "SeasonWeek"),
			WeekType:   stringField(record, "SeasonWeekType"),
			Status:     stringField(record, "GameStatus"),
			HomeTeam:   resolveTeamName(teamNames, record, "HomeTeam", f),
			AwayTeam:   resolveTeamName(teamNames, record, "AwayTeam", f),
		}
		if score, ok := intFieldOK(record, "HomeScore"); ok {
			game.HomeScore = &score
		}
		if score, ok := intFieldOK(record, "AwayScore"); ok {
			game.AwayScore = &score
		}
		applyGameDetails(&game, record, f, detailIndex)
		if opts.IncludeGameStats() {
			attachTeamGameStats(f, &game, record)
			if playerStats := buildPlayerGameStatsExports(record.Index, statsIndex); len(playerStats) > 0 {
				game.PlayerGameStats = playerStats
			}
		}
		games = append(games, game)
	}
	return games, nil
}

func (f *File) buildSeasonExport() *SeasonExport {
	seasonInfo, ok := f.GetTableByName("SeasonInfo")
	if !ok {
		return nil
	}
	if err := seasonInfo.ReadRecords(); err != nil || len(seasonInfo.Records) == 0 {
		return nil
	}
	row := seasonInfo.Records[0]
	return &SeasonExport{
		Year:     intField(row, "CurrentSeasonYear"),
		Week:     intField(row, "CurrentWeek"),
		WeekType: stringField(row, "CurrentWeekType"),
		Phase:    stringField(row, "CurrentStage"),
		Periods:  buildSeasonPeriodsExport(row),
	}
}

type teamKey struct {
	tableID   uint32
	rowNumber uint32
}

func buildTeamNameIndex(teamTable *Table) map[teamKey]string {
	names := make(map[teamKey]string, teamTable.ActiveRecordCount())
	tableID := teamTable.Header.TableID
	for _, record := range teamTable.Records {
		key := teamKey{tableID: tableID, rowNumber: uint32(record.Index)}
		if name := bestTeamName(record); name != "" {
			names[key] = name
		}
	}
	return names
}

func bestTeamName(record Record) string {
	for _, field := range []string{"LongName", "DisplayName", "ShortName"} {
		if name := stringField(record, field); isOfficialTeamName(name) {
			return name
		}
	}
	return ""
}

func resolveTeamName(index map[teamKey]string, record Record, field string, f *File) string {
	value, ok := record.Get(field)
	if !ok || value.Reference == nil {
		return ""
	}
	ref := value.Reference
	if name, ok := index[teamKey{tableID: ref.TableID, rowNumber: ref.RowNumber}]; ok {
		return name
	}
	if row, ok := f.RecordByReference("Team", ref); ok {
		if name := bestTeamName(row); name != "" {
			return name
		}
	}
	return referenceID(ref)
}
