package exporter

import "github.com/leaguelines/cfb-dynasty/dynasty"

var seasonCollection = Collection{
	Name:  "season",
	Title: "Season",
	Count: func(e dynasty.Export) int {
		if e.Season == nil {
			return 0
		}
		return 1
	},
	JSON: func(e dynasty.Export) any { return e.Season },
	Header: []string{
		"year", "week", "weekType", "phase",
		"isRecruitingPeriodActive", "isSigningPeriodActive", "isVisitingPeriodActive", "isPitchingPeriodActive",
		"isScholarshipPeriodActive", "isScoutingPeriodActive", "isTransferPortalNewlyAvailable", "isTransferSignPeriodActive",
		"isDraftPeriodActive", "isDraftScoutingActive", "isGoalsPeriodActive", "isCarouselPeriodActive",
		"isStaffHiringPeriodActive", "isWeeklyAwardPeriodActive", "isAnnualAwardPeriodActive",
	},
	Rows: func(e dynasty.Export) [][]string {
		if e.Season == nil {
			return nil
		}
		s := e.Season
		row := []string{fmtInt(s.Year), fmtInt(s.Week), s.WeekType, s.Phase}
		if s.Periods == nil {
			return [][]string{append(row, blanks(15)...)}
		}
		p := s.Periods
		row = append(row,
			fmtBool(p.IsRecruitingPeriodActive), fmtBool(p.IsSigningPeriodActive), fmtBool(p.IsVisitingPeriodActive), fmtBool(p.IsPitchingPeriodActive),
			fmtBool(p.IsScholarshipPeriodActive), fmtBool(p.IsScoutingPeriodActive), fmtBool(p.IsTransferPortalNewlyAvailable), fmtBool(p.IsTransferSignPeriodActive),
			fmtBool(p.IsDraftPeriodActive), fmtBool(p.IsDraftScoutingActive), fmtBool(p.IsGoalsPeriodActive), fmtBool(p.IsCarouselPeriodActive),
			fmtBool(p.IsStaffHiringPeriodActive), fmtBool(p.IsWeeklyAwardPeriodActive), fmtBool(p.IsAnnualAwardPeriodActive),
		)
		return [][]string{row}
	},
}

var teamsCollection = Collection{
	Name:   "teams",
	Title:  "Teams",
	Count:  func(e dynasty.Export) int { return len(e.Teams) },
	JSON:   func(e dynasty.Export) any { return e.Teams },
	Header: []string{"id", "shortName", "longName", "displayName", "conference", "overallWins", "overallLosses", "conferenceWins", "conferenceLosses", "coachesPollRank", "mediaPollRank", "cfpPollRank", "offensiveRank", "defensiveRank", "offensiveRating", "defensiveRating", "overallRating", "prestigeRank"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Teams))
		for _, t := range e.Teams {
			rows = append(rows, []string{
				fmtInt(t.ID), t.ShortName, t.LongName, t.DisplayName, t.Conference,
				fmtIntPtr(t.OverallWins), fmtIntPtr(t.OverallLosses), fmtIntPtr(t.ConferenceWins), fmtIntPtr(t.ConferenceLosses),
				fmtIntPtr(t.CoachesPollRank), fmtIntPtr(t.MediaPollRank), fmtIntPtr(t.CFPPollRank),
				fmtIntPtr(t.OffensiveRank), fmtIntPtr(t.DefensiveRank),
				fmtIntPtr(t.OffensiveRating), fmtIntPtr(t.DefensiveRating), fmtIntPtr(t.OverallRating),
				fmtIntPtr(t.PrestigeRank),
			})
		}
		return rows
	},
}

var rostersCollection = Collection{
	Name:  "rosters",
	Title: "Rosters",
	Count: func(e dynasty.Export) int { return len(e.Rosters) },
	JSON:  func(e dynasty.Export) any { return e.Rosters },
	Header: []string{
		"firstName", "lastName", "position", "isAth", "overall", "starRating", "class",
		"team", "archetype", "age", "height", "weight", "jersey",
		"homeState", "homeTown", "skillGroupCapTotal", "skillGroupUnlockedTotal", "devTrait",
		"redshirtStatus", "isNIL", "nilCompensation", "fatigue", "injuryStatus",
		"isImpactPlayer", "isCaptain", "wearAndTear", "skillGroups",
		"teamId", "playerId", "teamIndex",
	},
	Rows: func(e dynasty.Export) [][]string {
		var rows [][]string
		for _, r := range e.Rosters {
			if !IsExportableRosterTeam(&r) {
				continue
			}
			for i := range r.Players {
				p := &r.Players[i]
				if !IsExportablePlayer(p) {
					continue
				}
				v := playerVals(p)
				rows = append(rows, []string{
					v.FirstName, v.LastName, v.Position, v.IsAth, v.Overall, v.StarRating, v.SchoolYear,
					r.TeamName, v.ArchetypeLabel, v.Age, v.Height, v.Weight, v.Jersey,
					v.HomeState, v.HomeTown, v.SkillGroupCapTotal, fmtInt(p.SkillGroupUnlockedTotal), p.DevTrait,
					p.RedshirtStatus, fmtBool(p.IsNIL), fmtIntPtr(p.NILCompensation), fmtIntPtr(p.Fatigue), p.InjuryStatus,
					fmtBool(p.IsImpactPlayer), fmtBool(p.IsCaptain), fmtMapInt(p.WearAndTear), fmtSkillGroupsSummary(p.SkillGroups),
					fmtInt(r.TeamID), v.ID, v.TeamIndex,
				})
			}
		}
		return rows
	},
}

var gamesCollection = Collection{
	Name:       "games",
	Title:      "Games",
	Count:      func(e dynasty.Export) int { return len(e.Games) },
	JSON:       func(e dynasty.Export) any { return e.Games },
	LinkPrefix: "/game/",
	LinkColumn: 0,
	Header:     []string{"id", "seasonYear", "week", "weekType", "status", "homeTeam", "awayTeam", "homeScore", "awayScore"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Games))
		for _, g := range e.Games {
			rows = append(rows, []string{
				fmtUint32(g.ID), fmtInt(g.SeasonYear), fmtInt(g.Week), g.WeekType, g.Status,
				g.HomeTeam, g.AwayTeam, fmtIntPtr(g.HomeScore), fmtIntPtr(g.AwayScore),
			})
		}
		return rows
	},
}

var recruitsCollection = Collection{
	Name:  "recruits",
	Title: "Recruits",
	Count: func(e dynasty.Export) int { return len(e.Recruits) },
	JSON:  func(e dynasty.Export) any { return e.Recruits },
	Header: []string{
		"firstName", "lastName", "position", "isAth", "overall", "starRating",
		"nationalRank", "positionRank", "stateRank", "class", "recruitStage", "recruitStageAdvance",
		"archetype", "height", "weight", "homeState", "homeTown",
		"commitScore", "totalScholarshipOffers", "productionGrade", "qualityModifier",
		"skillGroupCapTotal", "skillGroupUnlockedTotal", "schoolInterest",
		"altPosition1", "altPosition2", "recruitId", "playerId",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Recruits))
		for i := range e.Recruits {
			r := e.Recruits[i]
			p := playerVals(r.Player)
			unlocked := ""
			if r.Player != nil {
				unlocked = fmtInt(r.Player.SkillGroupUnlockedTotal)
			}
			rows = append(rows, []string{
				p.FirstName, p.LastName, p.Position, p.IsAth, p.Overall, p.StarRating,
				fmtIntPtr(r.NationalRank), fmtIntPtr(r.PositionRank), fmtIntPtr(r.StateRank), r.Class, r.RecruitStage, r.RecruitStageAdvance,
				p.ArchetypeLabel, p.Height, p.Weight, p.HomeState, p.HomeTown,
				fmtIntPtr(r.CommitScore), fmtIntPtr(r.TotalScholarshipOffers), fmtIntPtr(r.ProductionGrade), r.QualityModifier,
				p.SkillGroupCapTotal, unlocked, fmtSchoolInterestList(r.SchoolInterest),
				r.AlternatePosition1, r.AlternatePosition2, fmtInt(r.ID), p.ID,
			})
		}
		return rows
	},
}

var recruitingCollection = Collection{
	Name:  "recruiting",
	Title: "Recruiting",
	Count: func(e dynasty.Export) int { return len(e.Recruiting) },
	JSON:  func(e dynasty.Export) any { return e.Recruiting },
	Header: []string{
		"recruitId", "topSchoolTeamId", "topSchoolTeamName", "topSchoolInfluence",
		"schoolInterest", "activePitches", "scheduledVisit",
		"scholarshipStatus", "currentNilOffer", "nilExpectation", "originalNilExpectation",
		"currentScholarshipBonus", "prospectInfluenceTotal", "prospectInfluenceDelta",
		"prospectInfluenceTotalLastWeek", "prospectHoursSpentCurrent", "committedWeekNumber",
		"swayPitch", "contactFriendsAndFamily", "contactHighSchoolCoaches", "searchSocialMedia",
		"sendTheHouse", "visitRecruitsSchool", "isFavorite",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Recruiting))
		for _, r := range e.Recruiting {
			var teamID, teamName, influence string
			if r.TopSchool != nil {
				teamID = fmtInt(r.TopSchool.TeamID)
				teamName = r.TopSchool.TeamName
				influence = fmtInt(r.TopSchool.Influence)
			}
			rows = append(rows, []string{
				fmtInt(r.RecruitID), teamID, teamName, influence,
				fmtSchoolInterestList(r.SchoolInterest), fmtActivePitches(r.ActivePitches), fmtScheduledVisit(r.ScheduledVisit),
				r.ScholarshipStatus, fmtIntPtr(r.CurrentNILOffer), fmtIntPtr(r.NILExpectation), fmtIntPtr(r.OriginalNILExpectation),
				fmtIntPtr(r.CurrentScholarshipBonus), fmtIntPtr(r.ProspectInfluenceTotal), fmtIntPtr(r.ProspectInfluenceDelta),
				fmtIntPtr(r.ProspectInfluenceTotalLastWeek), fmtIntPtr(r.ProspectHoursSpentCurrent), fmtIntPtr(r.CommittedWeekNumber),
				r.SwayPitch, fmtBool(r.ContactFriendsAndFamily), fmtBool(r.ContactHighSchoolCoaches), fmtBool(r.SearchSocialMedia),
				fmtBool(r.SendTheHouse), fmtBool(r.VisitRecruitsSchool), fmtBool(r.IsFavorite),
			})
		}
		return rows
	},
}

var seasonPlayerStatsCollection = Collection{
	Name:       "seasonPlayerStats",
	Title:      "Season Player Stats",
	LinkPrefix: "/player/",
	LinkColumn: 0,
	Count:      func(e dynasty.Export) int { return len(e.SeasonPlayerStats) },
	JSON:       func(e dynasty.Export) any { return e.SeasonPlayerStats },
	Header:     []string{"playerId", "firstName", "lastName", "teamId", "seasonYear", "gamesPlayed", "gamesStarted", "passYards", "passTDs", "passInts", "rushYards", "rushTDs", "recYards", "recTDs", "receptions", "tackles", "sacks", "ints", "forcedFumbles", "kickReturns", "kickReturnYards", "puntReturns", "puntReturnYards"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.SeasonPlayerStats))
		for _, s := range e.SeasonPlayerStats {
			var passYds, passTDs, passInts, rushYds, rushTDs, recYds, recTDs, rec string
			if s.Offense != nil {
				o := s.Offense
				passYds, passTDs, passInts = fmtIntPtr(o.PassYards), fmtIntPtr(o.PassTDs), fmtIntPtr(o.PassInts)
				rushYds, rushTDs = fmtIntPtr(o.RushYards), fmtIntPtr(o.RushTDs)
				recYds, recTDs, rec = fmtIntPtr(o.RecYards), fmtIntPtr(o.RecTDs), fmtIntPtr(o.Receptions)
			}
			var tackles, sacks, ints, ff string
			if s.Defense != nil {
				d := s.Defense
				tackles, sacks, ints, ff = fmtIntPtr(d.Tackles), fmtIntPtr(d.Sacks), fmtIntPtr(d.Ints), fmtIntPtr(d.ForcedFumbles)
			}
			var kr, krYds, pr, prYds string
			if s.SpecialTeams != nil {
				st := s.SpecialTeams
				kr, krYds = fmtIntPtr(st.KickReturns), fmtIntPtr(st.KickReturnYards)
				pr, prYds = fmtIntPtr(st.PuntReturns), fmtIntPtr(st.PuntReturnYards)
			}
			rows = append(rows, []string{
				fmtInt(s.PlayerID), s.FirstName, s.LastName, fmtIntPtr(s.TeamID), fmtIntPtr(s.SeasonYear),
				fmtIntPtr(s.GamesPlayed), fmtIntPtr(s.GamesStarted),
				passYds, passTDs, passInts, rushYds, rushTDs, recYds, recTDs, rec,
				tackles, sacks, ints, ff, kr, krYds, pr, prYds,
			})
		}
		return rows
	},
}

var seasonTeamStatsCollection = Collection{
	Name:   "seasonTeamStats",
	Title:  "Season Team Stats",
	Count:  func(e dynasty.Export) int { return len(e.SeasonTeamStats) },
	JSON:   func(e dynasty.Export) any { return e.SeasonTeamStats },
	Header: []string{"teamId", "teamName", "wins", "losses", "totalYards", "passYards", "rushYards", "passTDs", "rushTDs", "firstDowns", "turnovers", "sacks", "defPassYards", "defRushYards", "kickReturnYards", "puntReturnYards", "puntYards", "specialTeamYards"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.SeasonTeamStats))
		for _, t := range e.SeasonTeamStats {
			base := []string{fmtInt(t.TeamID), t.TeamName}
			if t.Stats != nil {
				s := t.Stats
				base = append(base,
					fmtIntPtr(s.Wins), fmtIntPtr(s.Losses), fmtIntPtr(s.TotalYards), fmtIntPtr(s.PassYards), fmtIntPtr(s.RushYards),
					fmtIntPtr(s.PassTDs), fmtIntPtr(s.RushTDs), fmtIntPtr(s.FirstDowns), fmtIntPtr(s.Turnovers), fmtIntPtr(s.Sacks),
					fmtIntPtr(s.DefPassYards), fmtIntPtr(s.DefRushYards),
					fmtIntPtr(s.KickReturnYards), fmtIntPtr(s.PuntReturnYards), fmtIntPtr(s.PuntYards), fmtIntPtr(s.SpecialTeamYards),
				)
			} else {
				base = append(base, blanks(16)...)
			}
			rows = append(rows, base)
		}
		return rows
	},
}

var coachesCollection = Collection{
	Name:  "coaches",
	Title: "Coaches",
	Count: func(e dynasty.Export) int { return len(e.Coaches) },
	JSON:  func(e dynasty.Export) any { return e.Coaches },
	Header: []string{
		"id", "firstName", "lastName", "teamId", "teamName", "position", "age", "level", "isUserControlled",
		"contractSalary", "contractYearsRemaining", "contractLength", "contractStatus", "seasonsWithTeam",
		"offensiveScheme", "defensiveScheme", "teamPhilosophy", "homeTown", "homeState",
		"jobSecurityStatus", "jobSecurityPercent", "coachPrestige", "coachPrestigeScore",
		"dominantArchetype", "specialtyType", "positionRatings", "seasonGoal",
		"currentWinStreak", "offTendencyRunPass", "defTendencyRunPass", "offTendencyAggressive", "defTendencyAggressive",
		"coachDemeanor", "earnedContractPointsThisYear",
		"career.wins", "career.losses", "career.ties", "career.winsAtCurrentSchool", "career.lossesAtCurrentSchool",
		"career.bowlWins", "career.bowlLosses", "career.playoffWins", "career.playoffLosses", "career.nationalChampionshipWins",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Coaches))
		for _, c := range e.Coaches {
			base := []string{
				fmtInt(c.ID), c.FirstName, c.LastName, fmtIntPtr(c.TeamID), c.TeamName, c.Position,
				fmtIntPtr(c.Age), fmtIntPtr(c.Level), fmtBool(c.IsUserControlled),
				fmtIntPtr(c.ContractSalary), fmtIntPtr(c.ContractYearsRemaining), fmtIntPtr(c.ContractLength), c.ContractStatus, fmtIntPtr(c.SeasonsWithTeam),
				c.OffensiveScheme, c.DefensiveScheme, c.TeamPhilosophy, c.HomeTown, c.HomeState,
				c.JobSecurityStatus, fmtIntPtr(c.JobSecurityPercent), c.CoachPrestige, fmtIntPtr(c.CoachPrestigeScore),
				c.DominantArchetype, c.SpecialtyType, fmtPositionRatings(c.PositionRatings), c.SeasonGoal,
				fmtIntPtr(c.CurrentWinStreak), fmtIntPtr(c.OffTendencyRunPass), fmtIntPtr(c.DefTendencyRunPass),
				fmtIntPtr(c.OffTendencyAggressive), fmtIntPtr(c.DefTendencyAggressive),
				c.CoachDemeanor, fmtIntPtr(c.EarnedContractPointsThisYear),
			}
			if c.Career != nil {
				cr := c.Career
				base = append(base,
					fmtIntPtr(cr.Wins), fmtIntPtr(cr.Losses), fmtIntPtr(cr.Ties),
					fmtIntPtr(cr.WinsAtCurrentSchool), fmtIntPtr(cr.LossesAtCurrentSchool),
					fmtIntPtr(cr.BowlWins), fmtIntPtr(cr.BowlLosses), fmtIntPtr(cr.PlayoffWins), fmtIntPtr(cr.PlayoffLosses), fmtIntPtr(cr.NCWins),
				)
			} else {
				base = append(base, blanks(10)...)
			}
			rows = append(rows, base)
		}
		return rows
	},
}

var leavingPlayersCollection = Collection{
	Name:   "leavingPlayers",
	Title:  "Leaving Players",
	Count:  func(e dynasty.Export) int { return len(e.LeavingPlayers) },
	JSON:   func(e dynasty.Export) any { return e.LeavingPlayers },
	Header: []string{"id", "playerId", "firstName", "lastName", "position", "teamId", "teamName", "leaveType", "leaveStatus", "draftClassPosition", "projectRound", "persuadeAttempts"},
	Rows: func(e dynasty.Export) [][]string {
		return LeavingCollectionRows(e.LeavingPlayers)
	},
}

var injuriesCollection = Collection{
	Name:   "injuries",
	Title:  "Injuries",
	Count:  func(e dynasty.Export) int { return len(e.Injuries) },
	JSON:   func(e dynasty.Export) any { return e.Injuries },
	Header: []string{"id", "playerId", "firstName", "lastName", "teamId", "teamName", "type", "severity", "minWeeks", "maxWeeks"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Injuries))
		for _, in := range e.Injuries {
			rows = append(rows, []string{
				fmtInt(in.ID), fmtInt(in.PlayerID), in.FirstName, in.LastName,
				fmtIntPtr(in.TeamID), in.TeamName, in.Type, in.Severity,
				fmtIntPtr(in.MinWeeks), fmtIntPtr(in.MaxWeeks),
			})
		}
		return rows
	},
}

var depthChartsCollection = Collection{
	Name:   "depthCharts",
	Title:  "Depth Charts",
	Count:  func(e dynasty.Export) int { return len(e.DepthCharts) },
	JSON:   func(e dynasty.Export) any { return e.DepthCharts },
	Header: []string{"teamId", "teamName", "position", "depth", "lockedDepth", "playerId", "firstName", "lastName", "userEditable"},
	Rows: func(e dynasty.Export) [][]string {
		var rows [][]string
		for _, dc := range e.DepthCharts {
			for _, s := range dc.Slots {
				rows = append(rows, []string{
					fmtInt(dc.TeamID), dc.TeamName, s.Position, fmtInt(s.Depth), fmtIntPtr(s.LockedDepth),
					fmtInt(s.PlayerID), s.FirstName, s.LastName, fmtBool(s.UserEditable),
				})
			}
		}
		return rows
	},
}

var playerAwardsCollection = Collection{
	Name:   "playerAwards",
	Title:  "Player Awards",
	Count:  func(e dynasty.Export) int { return len(e.PlayerAwards) },
	JSON:   func(e dynasty.Export) any { return e.PlayerAwards },
	Header: []string{"id", "playerId", "firstName", "lastName", "teamId", "teamName", "position", "awardType", "period", "periodIndex", "awardScore"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.PlayerAwards))
		for _, a := range e.PlayerAwards {
			rows = append(rows, []string{
				fmtInt(a.ID), fmtInt(a.PlayerID), a.FirstName, a.LastName, fmtIntPtr(a.TeamID), a.TeamName,
				a.Position, a.AwardType, a.Period, fmtIntPtr(a.PeriodIndex), fmtIntPtr(a.AwardScore),
			})
		}
		return rows
	},
}

var leagueAwardsCollection = Collection{
	Name:   "leagueAwards",
	Title:  "League Awards",
	Count:  func(e dynasty.Export) int { return len(e.LeagueAwards) },
	JSON:   func(e dynasty.Export) any { return e.LeagueAwards },
	Header: []string{"id", "firstName", "lastName", "position", "awardType", "teamDisplayName", "teamName"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.LeagueAwards))
		for _, a := range e.LeagueAwards {
			rows = append(rows, []string{
				fmtInt(a.ID), a.FirstName, a.LastName, a.Position, a.AwardType, a.TeamDisplayName, a.TeamName,
			})
		}
		return rows
	},
}

var conferenceChampionsCollection = Collection{
	Name:   "conferenceChampions",
	Title:  "Conference Champions",
	Count:  func(e dynasty.Export) int { return len(e.ConferenceChampions) },
	JSON:   func(e dynasty.Export) any { return e.ConferenceChampions },
	Header: []string{"id", "conferenceName", "winningTeamName", "losingTeamName", "winningScore", "losingScore", "winningCoachFirstName", "winningCoachLastName", "winningTeamRank", "losingTeamRank"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.ConferenceChampions))
		for _, c := range e.ConferenceChampions {
			rows = append(rows, []string{
				fmtInt(c.ID), c.ConferenceName, c.WinningTeamName, c.LosingTeamName,
				fmtIntPtr(c.WinningScore), fmtIntPtr(c.LosingScore),
				c.WinningCoachFirst, c.WinningCoachLast, fmtIntPtr(c.WinningTeamRank), fmtIntPtr(c.LosingTeamRank),
			})
		}
		return rows
	},
}

var recordBookCollection = Collection{
	Name:   "recordBook",
	Title:  "Record Book",
	Count:  func(e dynasty.Export) int { return len(e.RecordBook) },
	JSON:   func(e dynasty.Export) any { return e.RecordBook },
	Header: []string{"scope", "scopeName", "period", "statType", "rank", "statValue", "firstName", "lastName", "position", "teamName", "calendarYear"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.RecordBook))
		for _, r := range e.RecordBook {
			rows = append(rows, []string{
				r.Scope, r.ScopeName, r.Period, r.StatType, fmtInt(r.Rank), fmtInt(r.StatValue),
				r.FirstName, r.LastName, r.Position, r.TeamName, fmtIntPtr(r.CalendarYear),
			})
		}
		return rows
	},
}

var schoolGradesCollection = Collection{
	Name:  "schoolGrades",
	Title: "School Grades",
	Count: func(e dynasty.Export) int { return len(e.SchoolGrades) },
	JSON:  func(e dynasty.Export) any { return e.SchoolGrades },
	Header: []string{
		"teamId", "teamName", "academicPrestige", "athleticFacilities", "brandExposure", "campusLifestyle",
		"championshipContender", "coachPrestige", "coachStability", "conferencePrestige", "programTradition",
		"stadiumAtmosphere", "proPotentialQB", "proPotentialRB", "proPotentialWR", "proPotentialTE",
		"proPotentialOL", "proPotentialDL", "proPotentialLB", "proPotentialDB", "proPotentialK", "proPotentialP",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.SchoolGrades))
		for _, g := range e.SchoolGrades {
			rows = append(rows, []string{
				fmtInt(g.TeamID), g.TeamName,
				g.AcademicPrestige, g.AthleticFacilities, g.BrandExposure, g.CampusLifestyle,
				g.ChampionshipContender, g.CoachPrestige, g.CoachStability, g.ConferencePrestige, g.ProgramTradition,
				g.StadiumAtmosphere, g.ProPotentialQB, g.ProPotentialRB, g.ProPotentialWR, g.ProPotentialTE,
				g.ProPotentialOL, g.ProPotentialDL, g.ProPotentialLB, g.ProPotentialDB, g.ProPotentialK, g.ProPotentialP,
			})
		}
		return rows
	},
}

var pipelineInfluenceCollection = Collection{
	Name:   "pipelineInfluence",
	Title:  "Pipeline Influence",
	Count:  func(e dynasty.Export) int { return len(e.PipelineInfluence) },
	JSON:   func(e dynasty.Export) any { return e.PipelineInfluence },
	Header: []string{"teamId", "teamName", "pipeline", "influenceLevel", "influenceValue"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.PipelineInfluence))
		for _, p := range e.PipelineInfluence {
			rows = append(rows, []string{
				fmtInt(p.TeamID), p.TeamName, p.Pipeline, p.InfluenceLevel, fmtIntPtr(p.InfluenceValue),
			})
		}
		return rows
	},
}

var rivalriesCollection = Collection{
	Name:  "rivalries",
	Title: "Rivalries",
	Count: func(e dynasty.Export) int { return len(e.Rivalries) },
	JSON:  func(e dynasty.Export) any { return e.Rivalries },
	Header: []string{
		"id", "name", "team1Id", "team1Name", "team2Id", "team2Name",
		"team1Wins", "team2Wins", "streakLength", "streakTeamId", "team1LastScore", "team2LastScore",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.Rivalries))
		for _, r := range e.Rivalries {
			rows = append(rows, []string{
				fmtInt(r.ID), r.Name, fmtIntPtr(r.Team1ID), r.Team1Name, fmtIntPtr(r.Team2ID), r.Team2Name,
				fmtIntPtr(r.Team1Wins), fmtIntPtr(r.Team2Wins), fmtIntPtr(r.StreakLength), fmtIntPtr(r.StreakTeamID),
				fmtIntPtr(r.Team1LastScore), fmtIntPtr(r.Team2LastScore),
			})
		}
		return rows
	},
}

var positionChangesCollection = Collection{
	Name:  "positionChanges",
	Title: "Position Changes",
	Count: func(e dynasty.Export) int { return len(e.PositionChanges) },
	JSON:  func(e dynasty.Export) any { return e.PositionChanges },
	Header: []string{
		"id", "playerId", "oldPosition", "newPosition", "oldArchetype", "newArchetype",
		"seasonYear", "seasonWeek", "seasonStage", "oldTeamId", "newTeamId",
	},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.PositionChanges))
		for _, c := range e.PositionChanges {
			rows = append(rows, []string{
				fmtInt(c.ID), fmtIntPtr(c.PlayerID), c.OldPosition, c.NewPosition, c.OldArchetype, c.NewArchetype,
				fmtIntPtr(c.SeasonYear), fmtIntPtr(c.SeasonWeek), c.SeasonStage, fmtIntPtr(c.OldTeamID), fmtIntPtr(c.NewTeamID),
			})
		}
		return rows
	},
}

var draftPicksCollection = Collection{
	Name:   "draftPicks",
	Title:  "Draft Picks",
	Count:  func(e dynasty.Export) int { return len(e.DraftPicks) },
	JSON:   func(e dynasty.Export) any { return e.DraftPicks },
	Header: []string{"id", "yearOffset", "round", "pickNumber", "positionGroup", "teamId", "teamName"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.DraftPicks))
		for _, d := range e.DraftPicks {
			rows = append(rows, []string{
				fmtInt(d.ID), fmtIntPtr(d.YearOffset), fmtIntPtr(d.Round), fmtIntPtr(d.PickNumber),
				d.PositionGroup, fmtIntPtr(d.TeamID), d.TeamName,
			})
		}
		return rows
	},
}

var bowlGamesCollection = Collection{
	Name:   "bowlGames",
	Title:  "Bowl Games",
	Count:  func(e dynasty.Export) int { return len(e.BowlGames) },
	JSON:   func(e dynasty.Export) any { return e.BowlGames },
	Header: []string{"id", "name", "isPlayoffBowl", "playoffBracketSlot"},
	Rows: func(e dynasty.Export) [][]string {
		rows := make([][]string, 0, len(e.BowlGames))
		for _, b := range e.BowlGames {
			rows = append(rows, []string{
				fmtInt(b.ID), b.Name, fmtBool(b.IsPlayoffBowl), fmtIntPtr(b.PlayoffBracketSlot),
			})
		}
		return rows
	},
}
