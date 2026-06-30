package dynasty

import (
	"encoding/json"
	"time"
)

// Export is the stable JSON-oriented view of dynasty data for downstream apps.
type Export struct {
	Source     ExportSource     `json:"source"`
	Season     *SeasonExport    `json:"season,omitempty"`
	Teams      []TeamExport     `json:"teams,omitempty"`
	Games      []GameExport     `json:"games,omitempty"`
	ExportedAt time.Time        `json:"exportedAt"`
	Parser     ExportParserInfo `json:"parser"`
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
	ID              uint32 `json:"id,omitempty"`
	ShortName       string `json:"shortName,omitempty"`
	LongName        string `json:"longName,omitempty"`
	Conference      string `json:"conference,omitempty"`
	OverallWins     int    `json:"overallWins,omitempty"`
	OverallLosses   int    `json:"overallLosses,omitempty"`
	ConferenceWins  int    `json:"conferenceWins,omitempty"`
	ConferenceLosses int   `json:"conferenceLosses,omitempty"`
}

// GameExport is a normalized schedule/score row.
type GameExport struct {
	ID         uint32 `json:"id,omitempty"`
	SeasonYear int    `json:"seasonYear,omitempty"`
	Week       int    `json:"week,omitempty"`
	WeekType   string `json:"weekType,omitempty"`
	Status     string `json:"status,omitempty"`
	HomeTeam   string `json:"homeTeam,omitempty"`
	AwayTeam   string `json:"awayTeam,omitempty"`
	HomeScore  *int   `json:"homeScore,omitempty"`
	AwayScore  *int   `json:"awayScore,omitempty"`
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
