package dynasty

import (
	"encoding/json"
	"time"
)

// Export is the stable JSON-oriented view of dynasty data for downstream apps.
type Export struct {
	Source              ExportSource               `json:"source"`
	Season              *SeasonExport              `json:"season,omitempty"`
	Teams               []TeamExport               `json:"teams,omitempty"`
	Rosters             []RosterExport             `json:"rosters,omitempty"`
	Recruiting          []RecruitingTargetExport   `json:"recruiting,omitempty"`
	SeasonPlayerStats   []PlayerSeasonStatsExport  `json:"seasonPlayerStats,omitempty"`
	SeasonTeamStats     []TeamSeasonStatsExport    `json:"seasonTeamStats,omitempty"`
	Coaches             []CoachExport              `json:"coaches,omitempty"`
	LeavingPlayers      []LeavingPlayerExport      `json:"leavingPlayers,omitempty"`
	Injuries            []InjuryExport             `json:"injuries,omitempty"`
	DepthCharts         []DepthChartExport         `json:"depthCharts,omitempty"`
	PlayerAwards        []PlayerAwardExport        `json:"playerAwards,omitempty"`
	LeagueAwards        []LeagueHistoryAwardExport `json:"leagueAwards,omitempty"`
	ConferenceChampions []ConferenceChampionExport `json:"conferenceChampions,omitempty"`
	RecordBook          []RecordBookEntry          `json:"recordBook,omitempty"`
	Games               []GameExport               `json:"games,omitempty"`
	Recruits            []RecruitExport            `json:"recruits,omitempty"`
	SchoolGrades        []SchoolGradesExport       `json:"schoolGrades,omitempty"`
	PipelineInfluence   []PipelineInfluenceExport  `json:"pipelineInfluence,omitempty"`
	Rivalries           []RivalryExport            `json:"rivalries,omitempty"`
	PositionChanges     []PositionChangeExport     `json:"positionChanges,omitempty"`
	DraftPicks          []DraftPickExport          `json:"draftPicks,omitempty"`
	BowlGames           []BowlGameExport           `json:"bowlGames,omitempty"`
	ExportedAt          time.Time                  `json:"exportedAt"`
	Parser              ExportParserInfo           `json:"parser"`
}

// ExportSource identifies the save file that produced the export.
type ExportSource struct {
	Path          string         `json:"path"`
	Size          int64          `json:"size"`
	SHA256        string         `json:"sha256,omitempty"`
	GameYear      int            `json:"gameYear"`
	Format        Format         `json:"format"`
	Compressed    bool           `json:"compressed"`
	SchemaVersion *SchemaVersion `json:"schemaVersion,omitempty"`
}

// ExportParserInfo describes the library build used for export.
type ExportParserInfo struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// SeasonExport holds calendar and standings context.
type SeasonExport struct {
	Year     int    `json:"year,omitempty"`
	Week     int    `json:"week,omitempty"`
	WeekType string `json:"weekType,omitempty"`
	Phase    string `json:"phase,omitempty"`
	Periods  *SeasonPeriodsExport `json:"periods,omitempty"`
}

// SeasonPeriodsExport reports which offseason/recruiting phases are active.
type SeasonPeriodsExport struct {
	IsRecruitingPeriodActive           bool `json:"isRecruitingPeriodActive,omitempty"`
	IsSigningPeriodActive              bool `json:"isSigningPeriodActive,omitempty"`
	IsVisitingPeriodActive             bool `json:"isVisitingPeriodActive,omitempty"`
	IsPitchingPeriodActive             bool `json:"isPitchingPeriodActive,omitempty"`
	IsScholarshipPeriodActive          bool `json:"isScholarshipPeriodActive,omitempty"`
	IsScoutingPeriodActive             bool `json:"isScoutingPeriodActive,omitempty"`
	IsTransferPortalNewlyAvailable     bool `json:"isTransferPortalNewlyAvailable,omitempty"`
	IsTransferSignPeriodActive         bool `json:"isTransferSignPeriodActive,omitempty"`
	IsDraftPeriodActive                bool `json:"isDraftPeriodActive,omitempty"`
	IsDraftScoutingActive              bool `json:"isDraftScoutingActive,omitempty"`
	IsGoalsPeriodActive                bool `json:"isGoalsPeriodActive,omitempty"`
	IsCarouselPeriodActive             bool `json:"isCarouselPeriodActive,omitempty"`
	IsStaffHiringPeriodActive          bool `json:"isStaffHiringPeriodActive,omitempty"`
	IsWeeklyAwardPeriodActive          bool `json:"isWeeklyAwardPeriodActive,omitempty"`
	IsAnnualAwardPeriodActive          bool `json:"isAnnualAwardPeriodActive,omitempty"`
	RegularSeasonLastWeekScheduled     *int `json:"regularSeasonLastWeekScheduled,omitempty"`
	PostSeasonNumWeeks                 *int `json:"postSeasonNumWeeks,omitempty"`
}

// TeamExport is a normalized team record.
type TeamExport struct {
	ID               int    `json:"id"`
	ShortName        string `json:"shortName,omitempty"`
	LongName         string `json:"longName,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`
	Conference       string `json:"conference,omitempty"`
	OverallWins      *int   `json:"overallWins,omitempty"`
	OverallLosses    *int   `json:"overallLosses,omitempty"`
	ConferenceWins   *int   `json:"conferenceWins,omitempty"`
	ConferenceLosses *int   `json:"conferenceLosses,omitempty"`
	CoachesPollRank  *int   `json:"coachesPollRank,omitempty"`
	MediaPollRank    *int   `json:"mediaPollRank,omitempty"`
	CFPPollRank      *int   `json:"cfpPollRank,omitempty"`
	OffensiveRank    *int   `json:"offensiveRank,omitempty"`
	DefensiveRank    *int   `json:"defensiveRank,omitempty"`
	PrestigeRank     *int   `json:"prestigeRank,omitempty"`
	TeamPrestige     *int   `json:"teamPrestige,omitempty"` // half-star units: 10 = 5★, 9 = 4.5★, …
	OffensiveScheme  string `json:"offensiveScheme,omitempty"`
	DefensiveScheme  string `json:"defensiveScheme,omitempty"`
	Philosophy       string `json:"philosophy,omitempty"`
	PlayoffStatus    string `json:"playoffStatus,omitempty"`
	PlayoffRoundReached string `json:"playoffRoundReached,omitempty"`
	CoachesPollPoints *int  `json:"coachesPollPoints,omitempty"`
	MediaPollPoints   *int  `json:"mediaPollPoints,omitempty"`
	CFPPollPoints     *int  `json:"cfpPollPoints,omitempty"`
	StaffIDs           *TeamStaffIDsExport `json:"staffIds,omitempty"`
	CoachesPollFirstPlaceVotes *int `json:"coachesPollFirstPlaceVotes,omitempty"`
	CoachesPollLastWeekRank    *int `json:"coachesPollLastWeekRank,omitempty"`
	MediaPollLastWeekRank      *int `json:"mediaPollLastWeekRank,omitempty"`
	CFPPollLastWeekRank        *int `json:"cfpPollLastWeekRank,omitempty"`
	ConferenceStanding         *int `json:"conferenceStanding,omitempty"`
	DivisionStanding           *int `json:"divisionStanding,omitempty"`
	SeasonWinLossStreak        *int `json:"seasonWinLossStreak,omitempty"`
	ProgramPointBudget         *int `json:"programPointBudget,omitempty"`
	RemainingProgramPoints     *int `json:"remainingProgramPoints,omitempty"`
	RecruitProgramPointsSpent  *int `json:"recruitProgramPointsSpent,omitempty"`
	ScoutingPoints             *int `json:"scoutingPoints,omitempty"`
	RecruitingBoard            *TeamRecruitingBoardExport `json:"recruitingBoard,omitempty"`
}

// TeamRecruitingBoardExport is a team's recruiting hour budget state.
type TeamRecruitingBoardExport struct {
	HoursAssigned  *int `json:"hoursAssigned,omitempty"`
	HoursProcessed *int `json:"hoursProcessed,omitempty"`
	HoursTotal     *int `json:"hoursTotal,omitempty"`
}

// TeamStaffIDsExport links a team to its primary coaching staff row indices.
type TeamStaffIDsExport struct {
	HeadCoachID               *int `json:"headCoachId,omitempty"`
	OffensiveCoordinatorID    *int `json:"offensiveCoordinatorId,omitempty"`
	DefensiveCoordinatorID    *int `json:"defensiveCoordinatorId,omitempty"`
}

// RosterExport is a team's active roster with player attributes.
type RosterExport struct {
	TeamID   int            `json:"teamId"`
	TeamName string         `json:"teamName,omitempty"`
	Players  []PlayerExport `json:"players,omitempty"`
}

// GameExport is a normalized schedule/score row.
type GameExport struct {
	ID              uint32                  `json:"id,omitempty"`
	SeasonYear      int                     `json:"seasonYear,omitempty"`
	Week            int                     `json:"week,omitempty"`
	WeekType        string                  `json:"weekType,omitempty"`
	Status          string                  `json:"status,omitempty"`
	HomeTeam        string                  `json:"homeTeam,omitempty"`
	AwayTeam        string                  `json:"awayTeam,omitempty"`
	HomeScore       *int                    `json:"homeScore,omitempty"`
	AwayScore       *int                    `json:"awayScore,omitempty"`
	HomeQuarterScores []int                 `json:"homeQuarterScores,omitempty"`
	AwayQuarterScores []int                 `json:"awayQuarterScores,omitempty"`
	HomeScoreOT       *int                    `json:"homeScoreOT,omitempty"`
	AwayScoreOT       *int                    `json:"awayScoreOT,omitempty"`
	Attendance        *int                    `json:"attendance,omitempty"`
	Temperature       *int                    `json:"temperature,omitempty"`
	Weather           string                  `json:"weather,omitempty"`
	Wind              string                  `json:"wind,omitempty"`
	WindSpeed         *int                    `json:"windSpeed,omitempty"`
	IsSimmed          bool                    `json:"isSimmed,omitempty"`
	IsOvertime        bool                    `json:"isOvertime,omitempty"`
	BroadcastNetwork  string                  `json:"broadcastNetwork,omitempty"`
	StadiumName       string                  `json:"stadiumName,omitempty"`
	SeasonGameNum     *int                    `json:"seasonGameNum,omitempty"`
	DayOfWeek         string                  `json:"dayOfWeek,omitempty"`
	GameDateMonth     *int                    `json:"gameDateMonth,omitempty"`
	GameDateDay       *int                    `json:"gameDateDay,omitempty"`
	Precipitation     string                  `json:"precipitation,omitempty"`
	CloudCover        string                  `json:"cloudCover,omitempty"`
	PrecipitationChance *int                  `json:"precipitationChance,omitempty"`
	IsGameOfTheWeek   bool                    `json:"isGameOfTheWeek,omitempty"`
	IsKickoffGame     bool                    `json:"isKickoffGame,omitempty"`
	IsChallengeGame   bool                    `json:"isChallengeGame,omitempty"`
	IsRematch         bool                    `json:"isRematch,omitempty"`
	BowlGameID        *int                    `json:"bowlGameId,omitempty"`
	BowlGameName      string                  `json:"bowlGameName,omitempty"`
	ScoringPlays      []ScoringPlayExport     `json:"scoringPlays,omitempty"`
	KickingStats      []GameKickingStatsExport `json:"kickingStats,omitempty"`
	OLineStats        []GameOLineStatsExport  `json:"olineStats,omitempty"`
	HomeTeamStats   *TeamStatsExport        `json:"homeTeamStats,omitempty"`
	AwayTeamStats   *TeamStatsExport        `json:"awayTeamStats,omitempty"`
	PlayerGameStats []PlayerGameStatsExport `json:"playerGameStats,omitempty"`
}

// ScoringPlayExport is one scoring event from a game's scoring summary.
type ScoringPlayExport struct {
	Quarter           int    `json:"quarter,omitempty"`
	TimeStampSec      *int   `json:"timeStampSec,omitempty"`
	HomeScore         *int   `json:"homeScore,omitempty"`
	AwayScore         *int   `json:"awayScore,omitempty"`
	HomePreviousScore *int   `json:"homePreviousScore,omitempty"`
	AwayPreviousScore *int   `json:"awayPreviousScore,omitempty"`
	Conversion        string `json:"conversion,omitempty"`
	HomeScorerCount   *int   `json:"homeScorerCount,omitempty"`
	AwayScorerCount   *int   `json:"awayScorerCount,omitempty"`
}

// GameKickingStatsExport is kicking/punting stats for one team in a game.
type GameKickingStatsExport struct {
	FieldGoalsMade       *int `json:"fieldGoalsMade,omitempty"`
	FieldGoalsAttempted  *int `json:"fieldGoalsAttempted,omitempty"`
	FieldGoalLongest     *int `json:"fieldGoalLongest,omitempty"`
	ExtraPointsMade      *int `json:"extraPointsMade,omitempty"`
	ExtraPointsAttempted *int `json:"extraPointsAttempted,omitempty"`
	PuntAttempts         *int `json:"puntAttempts,omitempty"`
	PuntYards            *int `json:"puntYards,omitempty"`
	PuntLongest          *int `json:"puntLongest,omitempty"`
	PuntsInside20        *int `json:"puntsInside20,omitempty"`
}

// GameOLineStatsExport is offensive line stats for one team in a game.
type GameOLineStatsExport struct {
	Pancakes     *int `json:"pancakes,omitempty"`
	SacksAllowed *int `json:"sacksAllowed,omitempty"`
}

// TeamStatsExport is a team stat line (game or season).
type TeamStatsExport struct {
	Wins             *int `json:"wins,omitempty"`
	Losses           *int `json:"losses,omitempty"`
	TotalYards       *int `json:"totalYards,omitempty"`
	PassYards        *int `json:"passYards,omitempty"`
	RushYards        *int `json:"rushYards,omitempty"`
	PassAttempts     *int `json:"passAttempts,omitempty"`
	PassCompletions  *int `json:"passCompletions,omitempty"`
	PassTDs          *int `json:"passTDs,omitempty"`
	PassInts         *int `json:"passInts,omitempty"`
	RushAttempts     *int `json:"rushAttempts,omitempty"`
	RushTDs          *int `json:"rushTDs,omitempty"`
	FirstDowns       *int `json:"firstDowns,omitempty"`
	Turnovers        *int `json:"turnovers,omitempty"`
	Sacks            *int `json:"sacks,omitempty"`
	Ties             *int `json:"ties,omitempty"`
	Takeaways        *int `json:"takeaways,omitempty"`
	Giveaways        *int `json:"giveaways,omitempty"`
	DefPassYards     *int `json:"defPassYards,omitempty"`
	DefRushYards     *int `json:"defRushYards,omitempty"`
	KickReturnYards  *int `json:"kickReturnYards,omitempty"`
	PuntReturnYards  *int `json:"puntReturnYards,omitempty"`
	PuntYards        *int `json:"puntYards,omitempty"`
	SpecialTeamYards *int `json:"specialTeamYards,omitempty"`
}

// PlayerSeasonStatsExport is a player's accumulated season stat line.
type PlayerSeasonStatsExport struct {
	PlayerID     int                         `json:"playerId"`
	FirstName    string                      `json:"firstName,omitempty"`
	LastName     string                      `json:"lastName,omitempty"`
	TeamID       *int                        `json:"teamId,omitempty"`
	SeasonSlot   *int                        `json:"seasonSlot,omitempty"`
	SeasonYear   *int                        `json:"seasonYear,omitempty"`
	GamesPlayed  *int                        `json:"gamesPlayed,omitempty"`
	GamesStarted *int                        `json:"gamesStarted,omitempty"`
	Offense      *SeasonOffensiveStatsExport `json:"offense,omitempty"`
	Defense      *SeasonDefensiveStatsExport `json:"defense,omitempty"`
	SpecialTeams *SpecialTeamsStatsExport    `json:"specialTeams,omitempty"`
}

// SeasonOffensiveStatsExport is a season offensive stat line.
type SeasonOffensiveStatsExport struct {
	GameRating      *int `json:"gameRating,omitempty"`
	PassYards       *int `json:"passYards,omitempty"`
	PassTDs         *int `json:"passTDs,omitempty"`
	PassAttempts    *int `json:"passAttempts,omitempty"`
	PassCompletions *int `json:"passCompletions,omitempty"`
	PassInts        *int `json:"passInts,omitempty"`
	RushYards       *int `json:"rushYards,omitempty"`
	RushTDs         *int `json:"rushTDs,omitempty"`
	RushAttempts    *int `json:"rushAttempts,omitempty"`
	RecYards        *int `json:"recYards,omitempty"`
	RecTDs          *int `json:"recTDs,omitempty"`
	Receptions      *int `json:"receptions,omitempty"`
	FirstDowns      *int `json:"firstDowns,omitempty"`
}

// SeasonDefensiveStatsExport is a season defensive stat line.
type SeasonDefensiveStatsExport struct {
	GameRating       *int `json:"gameRating,omitempty"`
	Tackles          *int `json:"tackles,omitempty"`
	AssistTackles    *int `json:"assistTackles,omitempty"`
	Sacks            *int `json:"sacks,omitempty"`
	TacklesForLoss   *int `json:"tacklesForLoss,omitempty"`
	Ints             *int `json:"ints,omitempty"`
	ForcedFumbles    *int `json:"forcedFumbles,omitempty"`
	FumbleRecoveries *int `json:"fumbleRecoveries,omitempty"`
	PassDeflections  *int `json:"passDeflections,omitempty"`
}

// TeamSeasonStatsExport is a team's accumulated season stat line.
type TeamSeasonStatsExport struct {
	TeamID   int              `json:"teamId"`
	TeamName string           `json:"teamName,omitempty"`
	Stats    *TeamStatsExport `json:"stats,omitempty"`
}

// PlayerGameStatsExport holds one player's stat line for one game.
type PlayerGameStatsExport struct {
	PlayerID     *int                      `json:"playerId,omitempty"`
	Player       *PlayerExport             `json:"player,omitempty"`
	Offense      *OffensiveGameStatsExport `json:"offense,omitempty"`
	Defense      *DefensiveGameStatsExport `json:"defense,omitempty"`
	SpecialTeams *SpecialTeamsStatsExport  `json:"specialTeams,omitempty"`
}

// SpecialTeamsStatsExport is a kick/punt return stat line (game or season).
type SpecialTeamsStatsExport struct {
	KickReturns       *int `json:"kickReturns,omitempty"`
	KickReturnYards   *int `json:"kickReturnYards,omitempty"`
	KickReturnLongest *int `json:"kickReturnLongest,omitempty"`
	KickReturnTDs     *int `json:"kickReturnTDs,omitempty"`
	PuntReturns       *int `json:"puntReturns,omitempty"`
	PuntReturnYards   *int `json:"puntReturnYards,omitempty"`
	PuntReturnLongest *int `json:"puntReturnLongest,omitempty"`
	PuntReturnTDs     *int `json:"puntReturnTDs,omitempty"`
}

// OffensiveGameStatsExport is a single-game offensive stat line.
type OffensiveGameStatsExport struct {
	GameRating           *int `json:"gameRating,omitempty"`
	GamesStarted         *int `json:"gamesStarted,omitempty"`
	DownsPlayed          *int `json:"downsPlayed,omitempty"`
	PassYards            *int `json:"passYards,omitempty"`
	PassTDs              *int `json:"passTDs,omitempty"`
	PassAttempts         *int `json:"passAttempts,omitempty"`
	PassCompletions      *int `json:"passCompletions,omitempty"`
	PassInts             *int `json:"passInts,omitempty"`
	PassLongest          *int `json:"passLongest,omitempty"`
	PassSacked           *int `json:"passSacked,omitempty"`
	RushYards            *int `json:"rushYards,omitempty"`
	RushTDs              *int `json:"rushTDs,omitempty"`
	RushAttempts         *int `json:"rushAttempts,omitempty"`
	RushLongest          *int `json:"rushLongest,omitempty"`
	RushBrokenTackles    *int `json:"rushBrokenTackles,omitempty"`
	RushFumbles          *int `json:"rushFumbles,omitempty"`
	Rush20YardRuns       *int `json:"rush20YardRuns,omitempty"`
	RushYardsAfterFirstHit *int `json:"rushYardsAfterFirstHit,omitempty"`
	RecYards             *int `json:"recYards,omitempty"`
	RecTDs               *int `json:"recTDs,omitempty"`
	Receptions           *int `json:"receptions,omitempty"`
	RecLongest           *int `json:"recLongest,omitempty"`
	RecDrops             *int `json:"recDrops,omitempty"`
	RecYardsAfterCatch   *int `json:"recYardsAfterCatch,omitempty"`
}

// DefensiveGameStatsExport is a single-game defensive stat line.
type DefensiveGameStatsExport struct {
	GameRating            *int `json:"gameRating,omitempty"`
	GamesStarted          *int `json:"gamesStarted,omitempty"`
	DownsPlayed           *int `json:"downsPlayed,omitempty"`
	Tackles               *int `json:"tackles,omitempty"`
	AssistTackles         *int `json:"assistTackles,omitempty"`
	Sacks                 *int `json:"sacks,omitempty"`
	TacklesForLoss        *int `json:"tacklesForLoss,omitempty"`
	Ints                  *int `json:"ints,omitempty"`
	IntReturnYards        *int `json:"intReturnYards,omitempty"`
	IntReturnLongest      *int `json:"intReturnLongest,omitempty"`
	IntTDs                *int `json:"intTDs,omitempty"`
	ForcedFumbles         *int `json:"forcedFumbles,omitempty"`
	FumbleRecover         *int `json:"fumbleRecoveries,omitempty"`
	FumbleRecoverYards    *int `json:"fumbleRecoverYards,omitempty"`
	FumbleTDs             *int `json:"fumbleTDs,omitempty"`
	PassDeflections       *int `json:"passDeflections,omitempty"`
	BigHits               *int `json:"bigHits,omitempty"`
	CatchesAllowed        *int `json:"catchesAllowed,omitempty"`
	Safeties              *int `json:"safeties,omitempty"`
	Blocks                *int `json:"blocks,omitempty"`
}

// RecruitExport is a recruiting board prospect with nested player attributes.
type RecruitExport struct {
	ID                     int           `json:"id"`
	Class                  string        `json:"class,omitempty"`
	NationalRank           *int          `json:"nationalRank,omitempty"`
	StateRank              *int          `json:"stateRank,omitempty"`
	PositionRank           *int          `json:"positionRank,omitempty"`
	CommitScore            *int          `json:"commitScore,omitempty"`
	RecruitStage           string        `json:"recruitStage,omitempty"`
	ProductionGrade        *int          `json:"productionGrade,omitempty"`
	QualityModifier        string        `json:"qualityModifier,omitempty"`
	TotalScholarshipOffers *int          `json:"totalScholarshipOffers,omitempty"`
	AlternatePosition1     string        `json:"alternatePosition1,omitempty"`
	AlternatePosition2     string                          `json:"alternatePosition2,omitempty"`
	RecruitStageAdvance      string                          `json:"recruitStageAdvance,omitempty"`
	SchoolInterest           []RecruitingSchoolInterestExport `json:"schoolInterest,omitempty"`
	Player                 *PlayerExport `json:"player,omitempty"`
}

// RecruitingTargetExport is per-recruit pursuit state and school interest.
type RecruitingTargetExport struct {
	RecruitID                      int                              `json:"recruitId"`
	TopSchool                      *RecruitingSchoolInterestExport  `json:"topSchool,omitempty"`
	SchoolInterest                 []RecruitingSchoolInterestExport `json:"schoolInterest,omitempty"`
	ActivePitches                  []RecruitingPitchExport          `json:"activePitches,omitempty"`
	UnlockedIntelBitfield          *int                             `json:"unlockedIntelBitfield,omitempty"`
	ScholarshipStatus              string                          `json:"scholarshipStatus,omitempty"`
	CurrentNILOffer                *int                            `json:"currentNilOffer,omitempty"`
	NILExpectation                 *int                            `json:"nilExpectation,omitempty"`
	OriginalNILExpectation         *int                            `json:"originalNilExpectation,omitempty"`
	CurrentScholarshipBonus        *int                            `json:"currentScholarshipBonus,omitempty"`
	ProspectInfluenceTotal         *int                            `json:"prospectInfluenceTotal,omitempty"`
	ProspectInfluenceDelta         *int                            `json:"prospectInfluenceDelta,omitempty"`
	ProspectInfluenceTotalLastWeek *int                            `json:"prospectInfluenceTotalLastWeek,omitempty"`
	ProspectHoursSpentCurrent      *int                            `json:"prospectHoursSpentCurrent,omitempty"`
	CommittedWeekNumber            *int                            `json:"committedWeekNumber,omitempty"`
	SwayPitch                      string                          `json:"swayPitch,omitempty"`
	ScheduledVisit                 *RecruitingVisitExport          `json:"scheduledVisit,omitempty"`
	ContactFriendsAndFamily        bool                            `json:"contactFriendsAndFamily,omitempty"`
	ContactHighSchoolCoaches       bool                            `json:"contactHighSchoolCoaches,omitempty"`
	SearchSocialMedia              bool                            `json:"searchSocialMedia,omitempty"`
	SendTheHouse                   bool                            `json:"sendTheHouse,omitempty"`
	VisitRecruitsSchool            bool                            `json:"visitRecruitsSchool,omitempty"`
	IsFavorite                     bool                            `json:"isFavorite,omitempty"`
}

// RecruitingSchoolInterestExport is one school's interest in a prospect.
type RecruitingSchoolInterestExport struct {
	TeamID    int    `json:"teamId"`
	TeamName  string `json:"teamName,omitempty"`
	Influence int    `json:"influence,omitempty"`
}

// RecruitingPitchExport is an active recruiting pitch on a prospect.
type RecruitingPitchExport struct {
	Pitch     string `json:"pitch,omitempty"`
	Intensity string `json:"intensity,omitempty"`
}

// SchoolGradesExport is a team's recruiting pitch grades from MySchoolTrackingTable.
type SchoolGradesExport struct {
	TeamID                      int    `json:"teamId"`
	TeamName                    string `json:"teamName,omitempty"`
	AcademicPrestige            string `json:"academicPrestige,omitempty"`
	AthleticFacilities          string `json:"athleticFacilities,omitempty"`
	BrandExposure               string `json:"brandExposure,omitempty"`
	CampusLifestyle             string `json:"campusLifestyle,omitempty"`
	ChampionshipContender         string `json:"championshipContender,omitempty"`
	CoachPrestige               string `json:"coachPrestige,omitempty"`
	CoachStability              string `json:"coachStability,omitempty"`
	ConferencePrestige          string `json:"conferencePrestige,omitempty"`
	ProgramTradition            string `json:"programTradition,omitempty"`
	StadiumAtmosphere           string `json:"stadiumAtmosphere,omitempty"`
	ProPotentialQB              string `json:"proPotentialQB,omitempty"`
	ProPotentialRB              string `json:"proPotentialRB,omitempty"`
	ProPotentialWR              string `json:"proPotentialWR,omitempty"`
	ProPotentialTE              string `json:"proPotentialTE,omitempty"`
	ProPotentialOL              string `json:"proPotentialOL,omitempty"`
	ProPotentialDL              string `json:"proPotentialDL,omitempty"`
	ProPotentialLB              string `json:"proPotentialLB,omitempty"`
	ProPotentialDB              string `json:"proPotentialDB,omitempty"`
	ProPotentialK               string `json:"proPotentialK,omitempty"`
	ProPotentialP               string `json:"proPotentialP,omitempty"`
	AthleticFacilitiesScore     *int   `json:"athleticFacilitiesScore,omitempty"`
	CampusLifestyleScore        *int   `json:"campusLifestyleScore,omitempty"`
	BrandExposureNationalTVPlayed *int `json:"brandExposureNationalTVPlayed,omitempty"`
	BrandExposureNationalTVWins   *int `json:"brandExposureNationalTVWins,omitempty"`
	BrandExposureStreamingPlayed  *int `json:"brandExposureStreamingPlayed,omitempty"`
	BrandExposureStreamingWins    *int `json:"brandExposureStreamingWins,omitempty"`
	BrandExposureGOTWPlayed       *int `json:"brandExposureGamesOfTheWeekPlayed,omitempty"`
	BrandExposureGOTWWins         *int `json:"brandExposureGamesOfTheWeekWins,omitempty"`
	ChampionshipContenderCurrentRank *int `json:"championshipContenderCurrentRank,omitempty"`
}

// PipelineInfluenceExport is one pipeline region's influence for a school.
type PipelineInfluenceExport struct {
	TeamID         int    `json:"teamId"`
	TeamName       string `json:"teamName,omitempty"`
	Pipeline       string `json:"pipeline,omitempty"`
	InfluenceLevel string `json:"influenceLevel,omitempty"`
	InfluenceValue *int   `json:"influenceValue,omitempty"`
}

// RivalryExport is a head-to-head rivalry between two teams.
type RivalryExport struct {
	ID              int    `json:"id"`
	Name            string `json:"name,omitempty"`
	Team1ID         *int   `json:"team1Id,omitempty"`
	Team1Name       string `json:"team1Name,omitempty"`
	Team2ID         *int   `json:"team2Id,omitempty"`
	Team2Name       string `json:"team2Name,omitempty"`
	Team1Wins       *int   `json:"team1Wins,omitempty"`
	Team2Wins       *int   `json:"team2Wins,omitempty"`
	StreakLength    *int   `json:"streakLength,omitempty"`
	StreakTeamID    *int   `json:"streakTeamId,omitempty"`
	Team1LastScore  *int   `json:"team1LastScore,omitempty"`
	Team2LastScore  *int   `json:"team2LastScore,omitempty"`
}

// PositionChangeExport is a player position or archetype change event.
type PositionChangeExport struct {
	ID           int    `json:"id"`
	PlayerID     *int   `json:"playerId,omitempty"`
	OldPosition  string `json:"oldPosition,omitempty"`
	NewPosition  string `json:"newPosition,omitempty"`
	OldArchetype string `json:"oldArchetype,omitempty"`
	NewArchetype string `json:"newArchetype,omitempty"`
	SeasonYear   *int   `json:"seasonYear,omitempty"`
	SeasonWeek   *int   `json:"seasonWeek,omitempty"`
	SeasonStage  string `json:"seasonStage,omitempty"`
	OldTeamID    *int   `json:"oldTeamId,omitempty"`
	NewTeamID    *int   `json:"newTeamId,omitempty"`
}

// DraftPickExport is a draft selection slot assignment.
type DraftPickExport struct {
	ID            int    `json:"id"`
	YearOffset    *int   `json:"yearOffset,omitempty"`
	Round         *int   `json:"round,omitempty"`
	PickNumber    *int   `json:"pickNumber,omitempty"`
	PositionGroup string `json:"positionGroup,omitempty"`
	TeamID        *int   `json:"teamId,omitempty"`
	TeamName      string `json:"teamName,omitempty"`
}

// BowlGameExport is bowl game metadata from the save.
type BowlGameExport struct {
	ID           int    `json:"id"`
	Name         string `json:"name,omitempty"`
	IsPlayoffBowl bool  `json:"isPlayoffBowl,omitempty"`
	PlayoffBracketSlot *int `json:"playoffBracketSlot,omitempty"`
}

// RecruitingVisitExport is a scheduled recruiting visit.
type RecruitingVisitExport struct {
	Week     int    `json:"week,omitempty"`
	WeekType string `json:"weekType,omitempty"`
	Activity string `json:"activity,omitempty"`
}

// CoachExport is a coaching staff member with contract and career context.
// Staff contract fields come from the Coach table. PlayerBaseContract / ContractYearSummary
// (Madden/pro salary tables) are intentionally not exported for CFB dynasty use.
type CoachExport struct {
	ID                     int                     `json:"id"`
	FirstName              string                  `json:"firstName,omitempty"`
	LastName               string                  `json:"lastName,omitempty"`
	TeamID                 *int                    `json:"teamId,omitempty"`
	TeamName               string                  `json:"teamName,omitempty"`
	Position               string                  `json:"position,omitempty"`
	Age                    *int                    `json:"age,omitempty"`
	Level                  *int                    `json:"level,omitempty"`
	IsUserControlled       bool                    `json:"isUserControlled,omitempty"`
	ContractSalary         *int                    `json:"contractSalary,omitempty"`
	ContractYearsRemaining *int                    `json:"contractYearsRemaining,omitempty"`
	ContractLength         *int                    `json:"contractLength,omitempty"`
	ContractStatus         string                  `json:"contractStatus,omitempty"`
	SeasonsWithTeam        *int                    `json:"seasonsWithTeam,omitempty"`
	OffensiveScheme        string                  `json:"offensiveScheme,omitempty"`
	DefensiveScheme        string                  `json:"defensiveScheme,omitempty"`
	TeamPhilosophy         string                  `json:"teamPhilosophy,omitempty"`
	HomeTown               string                  `json:"homeTown,omitempty"`
	HomeState              string                  `json:"homeState,omitempty"`
	Career                 *CoachCareerStatsExport `json:"career,omitempty"`
	JobSecurityStatus      string                  `json:"jobSecurityStatus,omitempty"`
	JobSecurityPercent     *int                    `json:"jobSecurityPercent,omitempty"`
	CoachPrestige          string                  `json:"coachPrestige,omitempty"`
	CoachPrestigeScore     *int                    `json:"coachPrestigeScore,omitempty"`
	DominantArchetype      string                  `json:"dominantArchetype,omitempty"`
	SpecialtyType          string                  `json:"specialtyType,omitempty"`
	PositionRatings        map[string]int          `json:"positionRatings,omitempty"`
	SeasonGoal             string                  `json:"seasonGoal,omitempty"`
	CurrentWinStreak       *int                    `json:"currentWinStreak,omitempty"`
	OffTendencyRunPass     *int                    `json:"offTendencyRunPass,omitempty"`
	DefTendencyRunPass     *int                    `json:"defTendencyRunPass,omitempty"`
	OffTendencyAggressive  *int                    `json:"offTendencyAggressive,omitempty"`
	DefTendencyAggressive  *int                    `json:"defTendencyAggressive,omitempty"`
	CoachDemeanor          string                  `json:"coachDemeanor,omitempty"`
	EarnedContractPointsThisYear *int              `json:"earnedContractPointsThisYear,omitempty"`
}

// CoachCareerStatsExport is a coach's career record summary.
type CoachCareerStatsExport struct {
	Wins                  *int `json:"wins,omitempty"`
	Losses                *int `json:"losses,omitempty"`
	Ties                  *int `json:"ties,omitempty"`
	WinsAtCurrentSchool   *int `json:"winsAtCurrentSchool,omitempty"`
	LossesAtCurrentSchool *int `json:"lossesAtCurrentSchool,omitempty"`
	BowlWins              *int `json:"bowlWins,omitempty"`
	BowlLosses            *int `json:"bowlLosses,omitempty"`
	NCWins                *int `json:"nationalChampionshipWins,omitempty"`
	NCLosses              *int `json:"nationalChampionshipLosses,omitempty"`
	PlayoffWins           *int `json:"playoffWins,omitempty"`
	PlayoffLosses         *int `json:"playoffLosses,omitempty"`
	ConfChampWins         *int `json:"conferenceChampionshipWins,omitempty"`
	ConfChampLosses       *int `json:"conferenceChampionshipLosses,omitempty"`
	Top25Wins             *int `json:"top25Wins,omitempty"`
	RivalWins             *int `json:"rivalWins,omitempty"`
}

// LeavingPlayerExport is a player entering the offseason exit pipeline.
type LeavingPlayerExport struct {
	ID                 int    `json:"id"`
	PlayerID           int    `json:"playerId,omitempty"`
	FirstName          string `json:"firstName,omitempty"`
	LastName           string `json:"lastName,omitempty"`
	Position           string `json:"position,omitempty"`
	TeamID             *int   `json:"teamId,omitempty"`
	TeamName           string `json:"teamName,omitempty"`
	LeaveType          string `json:"leaveType,omitempty"`
	LeaveStatus        string `json:"leaveStatus,omitempty"`
	DraftClassPosition string `json:"draftClassPosition,omitempty"`
	ProjectRound       *int   `json:"projectRound,omitempty"`
	PersuadeAttempts   *int   `json:"persuadeAttempts,omitempty"`
}

// InjuryExport is an active player injury.
type InjuryExport struct {
	ID        int    `json:"id"`
	PlayerID  int    `json:"playerId,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	TeamID    *int   `json:"teamId,omitempty"`
	TeamName  string `json:"teamName,omitempty"`
	Type      string `json:"type,omitempty"`
	Severity  string `json:"severity,omitempty"`
	MinWeeks  *int   `json:"minWeeks,omitempty"`
	MaxWeeks  *int   `json:"maxWeeks,omitempty"`
}

// DepthChartExport is a team's depth chart slots.
type DepthChartExport struct {
	TeamID   int                    `json:"teamId"`
	TeamName string                 `json:"teamName,omitempty"`
	Slots    []DepthChartSlotExport `json:"slots,omitempty"`
}

// DepthChartSlotExport is one position slot on a depth chart.
type DepthChartSlotExport struct {
	Position     string `json:"position,omitempty"`
	Depth        int    `json:"depth,omitempty"`
	LockedDepth  *int   `json:"lockedDepth,omitempty"`
	PlayerID     int    `json:"playerId,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	UserEditable bool   `json:"userEditable,omitempty"`
}

// PlayerAwardExport is a player honor or award line.
type PlayerAwardExport struct {
	ID          int    `json:"id"`
	PlayerID    int    `json:"playerId,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	TeamID      *int   `json:"teamId,omitempty"`
	TeamName    string `json:"teamName,omitempty"`
	Position    string `json:"position,omitempty"`
	AwardType   string `json:"awardType,omitempty"`
	Period      string `json:"period,omitempty"`
	PeriodIndex *int   `json:"periodIndex,omitempty"`
	AwardScore  *int   `json:"awardScore,omitempty"`
}

// LeagueHistoryAwardExport is a league-wide historical award entry.
type LeagueHistoryAwardExport struct {
	ID              int    `json:"id"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Position        string `json:"position,omitempty"`
	AwardType       string `json:"awardType,omitempty"`
	TeamDisplayName string `json:"teamDisplayName,omitempty"`
	TeamName        string `json:"teamName,omitempty"`
}

// ConferenceChampionExport is a conference championship game result.
type ConferenceChampionExport struct {
	ID                int    `json:"id"`
	ConferenceName    string `json:"conferenceName,omitempty"`
	WinningTeamName   string `json:"winningTeamName,omitempty"`
	LosingTeamName    string `json:"losingTeamName,omitempty"`
	WinningScore      *int   `json:"winningScore,omitempty"`
	LosingScore       *int   `json:"losingScore,omitempty"`
	WinningCoachFirst string `json:"winningCoachFirstName,omitempty"`
	WinningCoachLast  string `json:"winningCoachLastName,omitempty"`
	WinningTeamRank   *int   `json:"winningTeamRank,omitempty"`
	LosingTeamRank    *int   `json:"losingTeamRank,omitempty"`
}

// RecordBookEntry is one stat-record-book line: a record holder within a given
// scope (league / conference / team) and period (career / season / game).
//
// The game keeps a record book per scope entity: the whole FBS (league), each
// conference, and each team. League boards store a ranked top-N per stat
// category; conference and team boards store a single record holder (rank 1).
type RecordBookEntry struct {
	Scope        string `json:"scope"`               // league | conference | team
	ScopeName    string `json:"scopeName,omitempty"` // conference/team name; empty for league
	Period       string `json:"period"`              // career | season | game
	StatType     string `json:"statType"`
	Rank         int    `json:"rank"`
	StatValue    int    `json:"statValue"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Position     string `json:"position,omitempty"`
	TeamName     string `json:"teamName,omitempty"` // record holder's team
	CalendarYear *int   `json:"calendarYear,omitempty"`
}

// PlayerExport is a normalized player attribute bundle.
type PlayerExport struct {
	ID                 int            `json:"id"`
	FirstName          string         `json:"firstName,omitempty"`
	LastName           string         `json:"lastName,omitempty"`
	Position           string         `json:"position,omitempty"`
	IsAth              bool           `json:"isAth,omitempty"`
	IsImpactPlayer     bool           `json:"isImpactPlayer,omitempty"`
	IsCaptain          bool           `json:"isCaptain,omitempty"`
	Archetype          string         `json:"archetype,omitempty"`
	ArchetypeLabel     string         `json:"archetypeLabel,omitempty"`
	SchoolYear         string         `json:"schoolYear,omitempty"`
	Age                *int           `json:"age,omitempty"`
	Height             *int           `json:"height,omitempty"`
	Weight             *int           `json:"weight,omitempty"`
	DevTrait           string         `json:"devTrait,omitempty"`
	Overall            *int           `json:"overall,omitempty"`
	StarRating         string         `json:"starRating,omitempty"`
	HomeTown           string         `json:"homeTown,omitempty"`
	HomeState          string         `json:"homeState,omitempty"`
	Jersey             *int           `json:"jersey,omitempty"`
	TeamIndex          *int           `json:"teamIndex,omitempty"`
	Ratings            map[string]int              `json:"ratings,omitempty"`
	SkillGroupCaps            []int              `json:"skillGroupCaps,omitempty"`            // greyed/capped upgrade slots per bucket
	SkillGroupUnlockedSlots   []int              `json:"skillGroupUnlockedSlots,omitempty"`   // unlocked upgrade slots saved per bucket (0..20)
	SkillGroupLabels          []string           `json:"skillGroupLabels,omitempty"`
	SkillGroupAttributeCounts []int              `json:"skillGroupAttributeCounts,omitempty"` // attribute definitions per bucket from tuning
	SkillGroups               []SkillGroupExport `json:"skillGroups,omitempty"`
	SkillGroupCapTotal        int                `json:"skillGroupCapTotal,omitempty"`        // total greyed/capped upgrade slots
	SkillGroupUnlockedTotal   int                `json:"skillGroupUnlockedTotal,omitempty"`     // total unlocked upgrade slots (recruit ceiling proxy)
	RedshirtStatus     string                      `json:"redshirtStatus,omitempty"`
	IsNIL              bool                        `json:"isNIL,omitempty"`
	NILBaseValue       *int                        `json:"nilBaseValue,omitempty"`
	NILCompensation    *int                        `json:"nilCompensation,omitempty"`
	HomePipeline       string                      `json:"homePipeline,omitempty"`
	Scheme             string                      `json:"scheme,omitempty"`
	Motivations        []string                    `json:"motivations,omitempty"`
	RecruitingDealbreaker string                   `json:"recruitingDealbreaker,omitempty"`
	IdealRecruitingPitch  string                   `json:"idealRecruitingPitch,omitempty"`
	Handedness         string                      `json:"handedness,omitempty"`
	QBStyle            string                      `json:"qbStyle,omitempty"`
	InjuryStatus       string                      `json:"injuryStatus,omitempty"`
	WasPreviouslyInjured bool                      `json:"wasPreviouslyInjured,omitempty"`
	ExperiencePoints   *int                        `json:"experiencePoints,omitempty"`
	LegacyScore        *int                        `json:"legacyScore,omitempty"`
	PrevTeamIndex      *int                        `json:"prevTeamIndex,omitempty"`
	PhysicalAbilities  []PlayerPhysicalAbilityExport `json:"physicalAbilities,omitempty"`
	MentalAbilities    []PlayerMentalAbilityExport   `json:"mentalAbilities,omitempty"`
	CareerStats        *PlayerCareerStatsExport    `json:"careerStats,omitempty"`
	ArchetypeTraits    []string                    `json:"archetypeTraits,omitempty"`
	Personality        string                      `json:"personality,omitempty"`
	PracticePlan       string                      `json:"practicePlan,omitempty"`
	Fatigue            *int                        `json:"fatigue,omitempty"`
	TransferChance     *int                        `json:"transferChance,omitempty"`
	WearAndTear        map[string]int              `json:"wearAndTear,omitempty"`
	DraftPick          *int                        `json:"draftPick,omitempty"`
	DraftRound         *int                        `json:"draftRound,omitempty"`
	ConsecutiveYearsWithTeam *int                  `json:"consecutiveYearsWithTeam,omitempty"`
}

// SkillGroupExport pairs per-bucket cap state with its tuning label when available.
type SkillGroupExport struct {
	Slot           int    `json:"slot,omitempty"` // save slot 1..6 (SkillGroupCap1..6)
	Label          string `json:"label,omitempty"`
	CappedSlots    int    `json:"cappedSlots,omitempty"`    // greyed-out upgrade slots in the UI
	UnlockedSlots  int    `json:"unlockedSlots,omitempty"`  // upgrade slots still available (saved value)
	AttributeCount int    `json:"attributeCount,omitempty"` // number of attributes in the bucket from tuning
}

// PlayerPhysicalAbilityExport is one equipped physical ability slot on a player.
type PlayerPhysicalAbilityExport struct {
	Slot int    `json:"slot,omitempty"`
	Name string `json:"name,omitempty"`
	Tier string `json:"tier,omitempty"`
}

// PlayerMentalAbilityExport is one equipped mental ability on a player.
type PlayerMentalAbilityExport struct {
	Name string `json:"name,omitempty"`
	Rank *int   `json:"rank,omitempty"`
}

// PlayerCareerStatsExport is lifetime playing-time and production from career stat rows.
type PlayerCareerStatsExport struct {
	GamesPlayed  *int                        `json:"gamesPlayed,omitempty"`
	GamesStarted *int                        `json:"gamesStarted,omitempty"`
	DownsPlayed  *int                        `json:"downsPlayed,omitempty"`
	GameRating   *int                        `json:"gameRating,omitempty"`
	Offense      *SeasonOffensiveStatsExport `json:"offense,omitempty"`
	Defense      *SeasonDefensiveStatsExport `json:"defense,omitempty"`
	SpecialTeams *SpecialTeamsStatsExport    `json:"specialTeams,omitempty"`
}

// MarshalJSON writes export data with stable field ordering defaults.
func (e Export) MarshalJSON() ([]byte, error) {
	type alias Export
	if e.ExportedAt.IsZero() {
		e.ExportedAt = time.Now().UTC()
	}
	if e.Parser.Module == "" {
		e.Parser.Module = "github.com/leaguelines/cfb-dynasty/dynasty"
	}
	if e.Parser.Version == "" {
		e.Parser.Version = "pre-alpha"
	}
	return json.Marshal(alias(e))
}

// ToJSON returns indented JSON for the export payload.
func (e Export) ToJSON() ([]byte, error) {
	type alias Export
	if e.ExportedAt.IsZero() {
		e.ExportedAt = time.Now().UTC()
	}
	if e.Parser.Module == "" {
		e.Parser.Module = "github.com/leaguelines/cfb-dynasty/dynasty"
	}
	if e.Parser.Version == "" {
		e.Parser.Version = "pre-alpha"
	}
	return json.MarshalIndent(alias(e), "", "  ")
}
