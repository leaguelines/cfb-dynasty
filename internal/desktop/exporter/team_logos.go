package exporter

import "github.com/leaguelines/cfb-dynasty/dynasty"

// teamLogoShortNames lists dynasty shortNames with a vendored PNG in web/static/teams/.
// To add a logo, drop SHORT.png there and add the shortName below.
var teamLogoShortNames = map[string]struct{}{
	"AFA": {}, "AKRN": {}, "APP": {}, "ARK": {}, "ARMY": {}, "ARST": {}, "ASU": {}, "AUB": {},
	"BALL": {}, "BAMA": {}, "BAY": {}, "BC": {}, "BGSU": {}, "BOI": {}, "BUFF": {}, "BYU": {},
	"CAL": {}, "CCAR": {}, "CIN": {}, "CLEM": {}, "CLT": {}, "CMU": {}, "CONN": {}, "CSU": {},
	"CU": {}, "CUSE": {}, "DEL": {}, "DUKE": {}, "ECU": {}, "EMU": {}, "FAU": {}, "FIU": {},
	"FRES": {}, "FSU": {}, "GASO": {}, "GAST": {}, "GT": {}, "HAW": {}, "HOU": {}, "ILL": {},
	"IOWA": {}, "ISU": {}, "IU": {}, "JMU": {}, "JXST": {}, "KENN": {}, "KENT": {}, "KSU": {},
	"KU": {}, "LOU": {}, "LSU": {}, "LTU": {}, "LU": {}, "M-OH": {}, "MASS": {}, "MEM": {},
	"MIA": {}, "MICH": {}, "MINN": {}, "MISS": {}, "MIZZ": {}, "MOST": {}, "MRSH": {}, "MSST": {},
	"MSU": {}, "MTSU": {}, "NAVY": {}, "NCST": {}, "ND": {}, "NDSU": {}, "NEB": {}, "NEV": {},
	"NIU": {}, "NMSU": {}, "NW": {}, "ODU": {}, "OHIO": {}, "OKLA": {}, "OKST": {}, "ORE": {},
	"ORST": {}, "OSU": {}, "PITT": {}, "PSU": {}, "PUR": {}, "RICE": {}, "RU": {}, "SAC": {},
	"SCAR": {}, "SDSU": {}, "SHSU": {}, "SJSU": {}, "SMU": {}, "STAN": {}, "TAMU": {}, "TCU": {},
	"TEM": {}, "TENN": {}, "TERPS": {}, "TEX": {}, "TLSA": {}, "TOL": {}, "TROY": {}, "TTU": {},
	"TUL": {}, "TXST": {}, "UAB": {}, "UCF": {}, "UCLA": {}, "UF": {}, "UGA": {}, "UK": {},
	"UL": {}, "ULM": {}, "UNC": {}, "UNLV": {}, "UNM": {}, "UNT": {}, "USA": {}, "USC": {},
	"USF": {}, "USM": {}, "USU": {}, "UTAH": {}, "UTEP": {}, "UTSA": {}, "UVA": {}, "VAND": {},
	"VT": {}, "WAKE": {}, "WASH": {}, "WISC": {}, "WKU": {}, "WMU": {}, "WSU": {}, "WVU": {},
	"WYO": {}, "ZONA": {},
}

// TeamLogoPath returns the static URL for a vendored team logo, or "" if unavailable.
func TeamLogoPath(shortName string) string {
	if shortName == "" {
		return ""
	}
	if _, ok := teamLogoShortNames[shortName]; !ok {
		return ""
	}
	return "/static/teams/" + shortName + ".png"
}

// TeamLogoForTeam returns the logo path for a team export row.
func TeamLogoForTeam(t dynasty.TeamExport) string {
	return TeamLogoPath(t.ShortName)
}

// TeamLogoForTeamID returns the logo path for a team ID using a shortName map.
func TeamLogoForTeamID(shortNames map[int]string, teamID int) string {
	if shortNames == nil || teamID <= 0 {
		return ""
	}
	return TeamLogoPath(shortNames[teamID])
}

// TeamShortNameByID returns the dynasty shortName for a team ID (e.g. "BAMA").
func TeamShortNameByID(teams []dynasty.TeamExport, teamID int) string {
	for _, t := range teams {
		if t.ID == teamID {
			return t.ShortName
		}
	}
	return ""
}

// BuildTeamShortNameMap indexes teams by ID for recruit/school display helpers.
func BuildTeamShortNameMap(teams []dynasty.TeamExport) map[int]string {
	if len(teams) == 0 {
		return nil
	}
	m := make(map[int]string, len(teams))
	for _, t := range teams {
		m[t.ID] = t.ShortName
	}
	return m
}
