package dynasty

// ExportOptions controls which sections are included in an export payload.
// When no section flags are set, all sections are included.
type ExportOptions struct {
	// Sections lists explicitly requested export sections. When Sections is zero
	// (no flags set), every section is exported.
	Sections ExportSections

	// OmitGameStats drops team and player stat lines from game exports.
	OmitGameStats bool
}

// ExportSections identifies top-level export sections selected on the CLI.
type ExportSections struct {
	Games    bool
	Recruits bool
	Season   bool
	Teams    bool
	Rosters  bool
	Recruiting bool
	SeasonStats bool
	Coaches bool
	LeavingPlayers bool
	Injuries bool
	DepthCharts bool
	History bool
	SchoolGrades bool
	Pipelines bool
	Rivalries bool
	PositionChanges bool
	Draft bool
	Bowls bool
}

// DefaultExportOptions exports every section including game stats.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{}
}

// IncludeGames reports whether games should be exported.
func (o ExportOptions) IncludeGames() bool {
	return o.Sections.all() || o.Sections.Games
}

// IncludeRecruits reports whether recruits should be exported.
func (o ExportOptions) IncludeRecruits() bool {
	return o.Sections.all() || o.Sections.Recruits
}

// IncludeSeason reports whether season metadata should be exported.
func (o ExportOptions) IncludeSeason() bool {
	return o.Sections.all() || o.Sections.Season
}

// IncludeTeams reports whether team records should be exported.
func (o ExportOptions) IncludeTeams() bool {
	return o.Sections.all() || o.Sections.Teams
}

// IncludeRosters reports whether team rosters should be exported.
func (o ExportOptions) IncludeRosters() bool {
	return o.Sections.all() || o.Sections.Rosters
}

// IncludeRecruiting reports whether recruiting pursuit data should be exported.
func (o ExportOptions) IncludeRecruiting() bool {
	return o.Sections.all() || o.Sections.Recruiting
}

// IncludeSeasonStats reports whether season stat totals should be exported.
func (o ExportOptions) IncludeSeasonStats() bool {
	return o.Sections.all() || o.Sections.SeasonStats
}

// IncludeCoaches reports whether coaching staff should be exported.
func (o ExportOptions) IncludeCoaches() bool {
	return o.Sections.all() || o.Sections.Coaches
}

// IncludeLeavingPlayers reports whether offseason exit players should be exported.
func (o ExportOptions) IncludeLeavingPlayers() bool {
	return o.Sections.all() || o.Sections.LeavingPlayers
}

// IncludeInjuries reports whether active injuries should be exported.
func (o ExportOptions) IncludeInjuries() bool {
	return o.Sections.all() || o.Sections.Injuries
}

// IncludeDepthCharts reports whether team depth charts should be exported.
func (o ExportOptions) IncludeDepthCharts() bool {
	return o.Sections.all() || o.Sections.DepthCharts
}

// IncludeHistory reports whether awards and league history should be exported.
func (o ExportOptions) IncludeHistory() bool {
	return o.Sections.all() || o.Sections.History
}

// IncludeSchoolGrades reports whether school recruiting grades should be exported.
func (o ExportOptions) IncludeSchoolGrades() bool {
	return o.Sections.all() || o.Sections.SchoolGrades
}

// IncludePipelines reports whether school pipeline influence should be exported.
func (o ExportOptions) IncludePipelines() bool {
	return o.Sections.all() || o.Sections.Pipelines
}

// IncludeRivalries reports whether rivalry rows should be exported.
func (o ExportOptions) IncludeRivalries() bool {
	return o.Sections.all() || o.Sections.Rivalries
}

// IncludePositionChanges reports whether player position change history should be exported.
func (o ExportOptions) IncludePositionChanges() bool {
	return o.Sections.all() || o.Sections.PositionChanges
}

// IncludeDraft reports whether draft pick rows should be exported.
func (o ExportOptions) IncludeDraft() bool {
	return o.Sections.all() || o.Sections.Draft
}

// IncludeBowls reports whether bowl game metadata should be exported.
func (o ExportOptions) IncludeBowls() bool {
	return o.Sections.all() || o.Sections.Bowls
}

// IncludeGameStats reports whether per-game stat lines should be attached.
func (o ExportOptions) IncludeGameStats() bool {
	return !o.OmitGameStats
}

func (s ExportSections) all() bool {
	return !s.Games && !s.Recruits && !s.Season && !s.Teams && !s.Rosters &&
		!s.Recruiting && !s.SeasonStats && !s.Coaches && !s.LeavingPlayers &&
		!s.Injuries && !s.DepthCharts && !s.History && !s.SchoolGrades && !s.Pipelines &&
		!s.Rivalries && !s.PositionChanges && !s.Draft && !s.Bowls
}
