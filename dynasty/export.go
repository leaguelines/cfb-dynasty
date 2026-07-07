package dynasty

import (
	"encoding/json"
	"time"
)

// Export is the stable JSON-oriented view of dynasty data for downstream apps.
type Export struct {
	Source     ExportSource      `json:"source"`
	Season     *SeasonExport     `json:"season,omitempty"`
	Teams      []TeamExport      `json:"teams,omitempty"`
	Rosters    []RosterExport    `json:"rosters,omitempty"`
	Recruiting []RecruitingTargetExport `json:"recruiting,omitempty"`
	SeasonPlayerStats []PlayerSeasonStatsExport `json:"seasonPlayerStats,omitempty"`
	SeasonTeamStats   []TeamSeasonStatsExport   `json:"seasonTeamStats,omitempty"`
	Coaches           []CoachExport             `json:"coaches,omitempty"`
	LeavingPlayers    []LeavingPlayerExport     `json:"leavingPlayers,omitempty"`
	Injuries          []InjuryExport            `json:"injuries,omitempty"`
	DepthCharts       []DepthChartExport        `json:"depthCharts,omitempty"`
	PlayerAwards      []PlayerAwardExport       `json:"playerAwards,omitempty"`
	LeagueAwards      []LeagueHistoryAwardExport `json:"leagueAwards,omitempty"`
	ConferenceChampions []ConferenceChampionExport `json:"conferenceChampions,omitempty"`
	StatRecords       []StatRecordExport        `json:"statRecords,omitempty"`
	Games      []GameExport      `json:"games,omitempty"`
	Recruits   []RecruitExport   `json:"recruits,omitempty"`
	ExportedAt time.Time         `json:"exportedAt"`
	Parser     ExportParserInfo  `json:"parser"`
}

// ExportSource identifies the save file that produced the export.
type ExportSource struct {
	Path           string         `json:"path"`
	Size           int64          `json:"size"`
	SHA256         string         `json:"sha256,omitempty"`
	GameYear       int            `json:"gameYear"`
	Format         Format         `json:"format"`
	Compressed     bool           `json:"compressed"`
	SchemaVersion  *SchemaVersion `json:"schemaVersion,omitempty"`
}

// ExportParserInfo describes the library build used for export.
type ExportParserInfo struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// SeasonExport holds calendar and standings context.
type SeasonExport struct {
	Year      int    `json:"year,omitempty"`
	Week      int    `json:"week,omitempty"`
	WeekType  string `json:"weekType,omitempty"`
	Phase     string `json:"phase,omitempty"`
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
	HomeTeamStats   *TeamStatsExport        `json:"homeTeamStats,omitempty"`
	AwayTeamStats   *TeamStatsExport        `json:"awayTeamStats,omitempty"`
	PlayerGameStats []PlayerGameStatsExport `json:"playerGameStats,omitempty"`
}

// TeamStatsExport is a team stat line (game or season).
type TeamStatsExport struct {
	Wins            *int `json:"wins,omitempty"`
	Losses          *int `json:"losses,omitempty"`
	TotalYards      *int `json:"totalYards,omitempty"`
	PassYards       *int `json:"passYards,omitempty"`
	RushYards       *int `json:"rushYards,omitempty"`
	PassAttempts    *int `json:"passAttempts,omitempty"`
	PassCompletions *int `json:"passCompletions,omitempty"`
	PassTDs         *int `json:"passTDs,omitempty"`
	PassInts        *int `json:"passInts,omitempty"`
	RushAttempts    *int `json:"rushAttempts,omitempty"`
	RushTDs         *int `json:"rushTDs,omitempty"`
	FirstDowns      *int `json:"firstDowns,omitempty"`
	Turnovers       *int `json:"turnovers,omitempty"`
	Sacks           *int `json:"sacks,omitempty"`
	Ties            *int `json:"ties,omitempty"`
	Takeaways       *int `json:"takeaways,omitempty"`
	Giveaways       *int `json:"giveaways,omitempty"`
	DefPassYards    *int `json:"defPassYards,omitempty"`
	DefRushYards    *int `json:"defRushYards,omitempty"`
}

// PlayerSeasonStatsExport is a player's accumulated season stat line.
type PlayerSeasonStatsExport struct {
	PlayerID     int                         `json:"playerId"`
	FirstName    string                      `json:"firstName,omitempty"`
	LastName     string                      `json:"lastName,omitempty"`
	TeamID       *int                        `json:"teamId,omitempty"`
	SeasonYear   *int                        `json:"seasonYear,omitempty"`
	GamesPlayed  *int                        `json:"gamesPlayed,omitempty"`
	GamesStarted *int                        `json:"gamesStarted,omitempty"`
	Offense      *SeasonOffensiveStatsExport `json:"offense,omitempty"`
	Defense      *SeasonDefensiveStatsExport `json:"defense,omitempty"`
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
	GameRating      *int `json:"gameRating,omitempty"`
	Tackles         *int `json:"tackles,omitempty"`
	AssistTackles   *int `json:"assistTackles,omitempty"`
	Sacks           *int `json:"sacks,omitempty"`
	TacklesForLoss  *int `json:"tacklesForLoss,omitempty"`
	Ints            *int `json:"ints,omitempty"`
	ForcedFumbles   *int `json:"forcedFumbles,omitempty"`
	FumbleRecoveries *int `json:"fumbleRecoveries,omitempty"`
	PassDeflections *int `json:"passDeflections,omitempty"`
}

// TeamSeasonStatsExport is a team's accumulated season stat line.
type TeamSeasonStatsExport struct {
	TeamID   int              `json:"teamId"`
	TeamName string           `json:"teamName,omitempty"`
	Stats    *TeamStatsExport `json:"stats,omitempty"`
}

// PlayerGameStatsExport holds per-player lines for one game.
type PlayerGameStatsExport struct {
	Slot      int                       `json:"slot,omitempty"`
	PlayerID  *int                      `json:"playerId,omitempty"`
	Player    *PlayerExport             `json:"player,omitempty"`
	Offense   *OffensiveGameStatsExport `json:"offense,omitempty"`
	Defense   *DefensiveGameStatsExport `json:"defense,omitempty"`
}

// OffensiveGameStatsExport is a single-game offensive stat line.
type OffensiveGameStatsExport struct {
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
}

// DefensiveGameStatsExport is a single-game defensive stat line.
type DefensiveGameStatsExport struct {
	GameRating    *int `json:"gameRating,omitempty"`
	Tackles       *int `json:"tackles,omitempty"`
	AssistTackles *int `json:"assistTackles,omitempty"`
	Sacks         *int `json:"sacks,omitempty"`
	Ints          *int `json:"ints,omitempty"`
	ForcedFumbles *int `json:"forcedFumbles,omitempty"`
	FumbleRecover *int `json:"fumbleRecoveries,omitempty"`
	PassDeflections *int `json:"passDeflections,omitempty"`
}

// RecruitExport is a recruiting board prospect with nested player attributes.
type RecruitExport struct {
	ID                     int          `json:"id"`
	Class                  string       `json:"class,omitempty"`
	NationalRank           *int         `json:"nationalRank,omitempty"`
	StateRank              *int         `json:"stateRank,omitempty"`
	PositionRank           *int         `json:"positionRank,omitempty"`
	CommitScore            *int         `json:"commitScore,omitempty"`
	RecruitStage           string       `json:"recruitStage,omitempty"`
	ProductionGrade        *int         `json:"productionGrade,omitempty"`
	QualityModifier        string       `json:"qualityModifier,omitempty"`
	TotalScholarshipOffers *int         `json:"totalScholarshipOffers,omitempty"`
	AlternatePosition1     string       `json:"alternatePosition1,omitempty"`
	AlternatePosition2     string       `json:"alternatePosition2,omitempty"`
	Player                 *PlayerExport `json:"player,omitempty"`
}

// RecruitingTargetExport is per-recruit pursuit state and school interest.
type RecruitingTargetExport struct {
	RecruitID                      int                            `json:"recruitId"`
	TopSchool                      *RecruitingSchoolInterestExport  `json:"topSchool,omitempty"`
	ScholarshipStatus              string                         `json:"scholarshipStatus,omitempty"`
	CurrentNILOffer                *int                           `json:"currentNilOffer,omitempty"`
	NILExpectation                 *int                           `json:"nilExpectation,omitempty"`
	OriginalNILExpectation         *int                           `json:"originalNilExpectation,omitempty"`
	CurrentScholarshipBonus        *int                           `json:"currentScholarshipBonus,omitempty"`
	ProspectInfluenceTotal         *int                           `json:"prospectInfluenceTotal,omitempty"`
	ProspectInfluenceDelta         *int                           `json:"prospectInfluenceDelta,omitempty"`
	ProspectInfluenceTotalLastWeek *int                           `json:"prospectInfluenceTotalLastWeek,omitempty"`
	ProspectHoursSpentCurrent      *int                           `json:"prospectHoursSpentCurrent,omitempty"`
	CommittedWeekNumber            *int                           `json:"committedWeekNumber,omitempty"`
	SwayPitch                      string                         `json:"swayPitch,omitempty"`
	ScheduledVisit                 *RecruitingVisitExport         `json:"scheduledVisit,omitempty"`
	ContactFriendsAndFamily        bool                           `json:"contactFriendsAndFamily,omitempty"`
	ContactHighSchoolCoaches       bool                           `json:"contactHighSchoolCoaches,omitempty"`
	SearchSocialMedia              bool                           `json:"searchSocialMedia,omitempty"`
	SendTheHouse                   bool                           `json:"sendTheHouse,omitempty"`
	VisitRecruitsSchool             bool                           `json:"visitRecruitsSchool,omitempty"`
	IsFavorite                     bool                           `json:"isFavorite,omitempty"`
}

// RecruitingSchoolInterestExport is one school's interest in a prospect.
type RecruitingSchoolInterestExport struct {
	TeamID     int    `json:"teamId"`
	TeamName   string `json:"teamName,omitempty"`
	Influence  int    `json:"influence,omitempty"`
}

// RecruitingVisitExport is a scheduled recruiting visit.
type RecruitingVisitExport struct {
	Week     int    `json:"week,omitempty"`
	WeekType string `json:"weekType,omitempty"`
	Activity string `json:"activity,omitempty"`
}

// CoachExport is a coaching staff member with contract and career context.
type CoachExport struct {
	ID                     int                    `json:"id"`
	FirstName              string                 `json:"firstName,omitempty"`
	LastName               string                 `json:"lastName,omitempty"`
	TeamID                 *int                   `json:"teamId,omitempty"`
	TeamName               string                 `json:"teamName,omitempty"`
	Position               string                 `json:"position,omitempty"`
	Age                    *int                   `json:"age,omitempty"`
	Level                  *int                   `json:"level,omitempty"`
	IsUserControlled       bool                   `json:"isUserControlled,omitempty"`
	ContractSalary         *int                   `json:"contractSalary,omitempty"`
	ContractYearsRemaining *int                   `json:"contractYearsRemaining,omitempty"`
	ContractLength         *int                   `json:"contractLength,omitempty"`
	ContractStatus         string                 `json:"contractStatus,omitempty"`
	SeasonsWithTeam        *int                   `json:"seasonsWithTeam,omitempty"`
	OffensiveScheme        string                 `json:"offensiveScheme,omitempty"`
	DefensiveScheme        string                 `json:"defensiveScheme,omitempty"`
	TeamPhilosophy         string                 `json:"teamPhilosophy,omitempty"`
	HomeTown               string                 `json:"homeTown,omitempty"`
	HomeState              string                 `json:"homeState,omitempty"`
	Career                 *CoachCareerStatsExport `json:"career,omitempty"`
}

// CoachCareerStatsExport is a coach's career record summary.
type CoachCareerStatsExport struct {
	Wins                *int `json:"wins,omitempty"`
	Losses              *int `json:"losses,omitempty"`
	Ties                *int `json:"ties,omitempty"`
	WinsAtCurrentSchool *int `json:"winsAtCurrentSchool,omitempty"`
	LossesAtCurrentSchool *int `json:"lossesAtCurrentSchool,omitempty"`
	BowlWins            *int `json:"bowlWins,omitempty"`
	BowlLosses          *int `json:"bowlLosses,omitempty"`
	NCWins              *int `json:"nationalChampionshipWins,omitempty"`
	NCLosses            *int `json:"nationalChampionshipLosses,omitempty"`
	PlayoffWins         *int `json:"playoffWins,omitempty"`
	PlayoffLosses       *int `json:"playoffLosses,omitempty"`
	ConfChampWins       *int `json:"conferenceChampionshipWins,omitempty"`
	ConfChampLosses     *int `json:"conferenceChampionshipLosses,omitempty"`
	Top25Wins           *int `json:"top25Wins,omitempty"`
	RivalWins           *int `json:"rivalWins,omitempty"`
}

// LeavingPlayerExport is a player entering the offseason exit pipeline.
type LeavingPlayerExport struct {
	ID                   int    `json:"id"`
	PlayerID             int    `json:"playerId,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	Position             string `json:"position,omitempty"`
	TeamID               *int   `json:"teamId,omitempty"`
	TeamName             string `json:"teamName,omitempty"`
	LeaveType            string `json:"leaveType,omitempty"`
	LeaveStatus          string `json:"leaveStatus,omitempty"`
	DraftClassPosition   string `json:"draftClassPosition,omitempty"`
	ProjectRound         *int   `json:"projectRound,omitempty"`
	PersuadeAttempts     *int   `json:"persuadeAttempts,omitempty"`
}

// InjuryExport is an active player injury.
type InjuryExport struct {
	ID         int    `json:"id"`
	PlayerID   int    `json:"playerId,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	TeamID     *int   `json:"teamId,omitempty"`
	TeamName   string `json:"teamName,omitempty"`
	Type       string `json:"type,omitempty"`
	Severity   string `json:"severity,omitempty"`
	MinWeeks   *int   `json:"minWeeks,omitempty"`
	MaxWeeks   *int   `json:"maxWeeks,omitempty"`
}

// DepthChartExport is a team's depth chart slots.
type DepthChartExport struct {
	TeamID   int                  `json:"teamId"`
	TeamName string               `json:"teamName,omitempty"`
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
	ID         int    `json:"id"`
	PlayerID   int    `json:"playerId,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	TeamID     *int   `json:"teamId,omitempty"`
	TeamName   string `json:"teamName,omitempty"`
	Position   string `json:"position,omitempty"`
	AwardType  string `json:"awardType,omitempty"`
	Period     string `json:"period,omitempty"`
	PeriodIndex *int  `json:"periodIndex,omitempty"`
	AwardScore *int   `json:"awardScore,omitempty"`
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
	ID                  int    `json:"id"`
	ConferenceName      string `json:"conferenceName,omitempty"`
	WinningTeamName     string `json:"winningTeamName,omitempty"`
	LosingTeamName      string `json:"losingTeamName,omitempty"`
	WinningScore        *int   `json:"winningScore,omitempty"`
	LosingScore         *int   `json:"losingScore,omitempty"`
	WinningCoachFirst   string `json:"winningCoachFirstName,omitempty"`
	WinningCoachLast    string `json:"winningCoachLastName,omitempty"`
	WinningTeamRank     *int   `json:"winningTeamRank,omitempty"`
	LosingTeamRank      *int   `json:"losingTeamRank,omitempty"`
}

// StatRecordExport is a career or season statistical record holder entry.
type StatRecordExport struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Position     string `json:"position,omitempty"`
	TeamName     string `json:"teamName,omitempty"`
	StatType     string `json:"statType,omitempty"`
	StatValue    *int   `json:"statValue,omitempty"`
	CalendarYear *int   `json:"calendarYear,omitempty"`
}

// PlayerExport is a normalized player attribute bundle.
type PlayerExport struct {
	ID          int            `json:"id"`
	FirstName   string         `json:"firstName,omitempty"`
	LastName    string         `json:"lastName,omitempty"`
	Position    string         `json:"position,omitempty"`
	Archetype      string         `json:"archetype,omitempty"`
	ArchetypeLabel string         `json:"archetypeLabel,omitempty"`
	SchoolYear     string         `json:"schoolYear,omitempty"`
	Age         *int           `json:"age,omitempty"`
	Height      *int           `json:"height,omitempty"`
	Weight      *int           `json:"weight,omitempty"`
	Overall     *int           `json:"overall,omitempty"`
	StarRating  string         `json:"starRating,omitempty"`
	HomeTown    string         `json:"homeTown,omitempty"`
	HomeState   string         `json:"homeState,omitempty"`
	Jersey      *int           `json:"jersey,omitempty"`
	TeamIndex   *int           `json:"teamIndex,omitempty"`
	Ratings     map[string]int `json:"ratings,omitempty"`
	SkillGroupCaps     []int `json:"skillGroupCaps,omitempty"`
	SkillGroupCapTotal int   `json:"skillGroupCapTotal,omitempty"`
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
