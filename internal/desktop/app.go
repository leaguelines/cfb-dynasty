package desktop

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/leaguelines/cfb-dynasty/dynasty"
	"github.com/leaguelines/cfb-dynasty/internal/desktop/exporter"
	"github.com/leaguelines/cfb-dynasty/internal/desktop/prefs"
	"github.com/leaguelines/cfb-dynasty/internal/desktop/saves"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails-bound desktop application API.
type App struct {
	ctx context.Context

	mu        sync.RWMutex
	schemaDir string
	export    *dynasty.Export
	source    string
	execDir   string

	// CLI overrides applied before startup.
	flagSchemaDir string
	flagSavePath  string
}

// NewApp creates the desktop app. Optional CLI overrides may be set via WithFlags.
func NewApp() *App {
	return &App{}
}

// WithFlags sets optional CLI overrides resolved at startup.
func (a *App) WithFlags(schemaDir, savePath string) *App {
	a.flagSchemaDir = strings.TrimSpace(schemaDir)
	a.flagSavePath = strings.TrimSpace(savePath)
	return a
}

// Startup is called by Wails when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if exe, err := os.Executable(); err == nil {
		a.execDir = filepath.Dir(exe)
	}
	a.schemaDir = a.resolveSchemaDir()
	if a.flagSavePath != "" && a.schemaDir != "" {
		_ = a.OpenSave(a.flagSavePath)
	}
}

func (a *App) resolveSchemaDir() string {
	if a.flagSchemaDir != "" {
		return a.flagSchemaDir
	}
	if env := strings.TrimSpace(os.Getenv("SCHEMA_DIR")); env != "" {
		return env
	}
	if p, err := prefs.Load(); err == nil && p.SchemaDir != "" {
		return p.SchemaDir
	}
	return firstValidSchema(defaultSchemaCandidates(a.execDir)...)
}

// ConfigView is the home-screen configuration snapshot.
type ConfigView struct {
	Schema       SchemaStatus     `json:"schema"`
	HasSave      bool             `json:"hasSave"`
	SourceName   string           `json:"sourceName,omitempty"`
	Discovered   []saves.SaveFile `json:"discovered"`
	DefaultSaves string           `json:"defaultSavesDir,omitempty"`
}

// GetConfig returns current schema/save setup state.
func (a *App) GetConfig() ConfigView {
	a.mu.RLock()
	defer a.mu.RUnlock()
	discovered, _ := saves.Discover()
	return ConfigView{
		Schema:       ValidateSchemaDir(a.schemaDir),
		HasSave:      a.export != nil,
		SourceName:   a.source,
		Discovered:   discovered,
		DefaultSaves: saves.DefaultWindowsSavesDir(),
	}
}

// GetSchemaStatus validates the current schema directory.
func (a *App) GetSchemaStatus() SchemaStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return ValidateSchemaDir(a.schemaDir)
}

// SetSchemaDir sets and persists the schema directory after validation.
func (a *App) SetSchemaDir(dir string) (SchemaStatus, error) {
	status := ValidateSchemaDir(dir)
	if !status.Valid {
		return status, fmt.Errorf("%s", status.Message)
	}
	a.mu.Lock()
	a.schemaDir = status.Dir
	a.mu.Unlock()
	_ = prefs.Save(prefs.Prefs{SchemaDir: status.Dir})
	return status, nil
}

// PickSchemaDir opens a native folder dialog and sets the schema directory.
func (a *App) PickSchemaDir() (SchemaStatus, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose schema directory (containing C27_*.gz)",
	})
	if err != nil {
		return SchemaStatus{}, err
	}
	if dir == "" {
		return a.GetSchemaStatus(), nil
	}
	return a.SetSchemaDir(dir)
}

// ListDiscoveredSaves returns auto-discovered save files.
func (a *App) ListDiscoveredSaves() ([]saves.SaveFile, error) {
	return saves.Discover()
}

// PickSaveFile opens a native file dialog for a dynasty save.
func (a *App) PickSaveFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open dynasty save",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// OpenSave parses a dynasty save at path into memory.
func (a *App) OpenSave(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("save path required")
	}
	a.mu.RLock()
	schemaDir := a.schemaDir
	a.mu.RUnlock()

	status := ValidateSchemaDir(schemaDir)
	if !status.Valid {
		return fmt.Errorf("configure a schema directory before opening a save: %s", status.Message)
	}

	settings := dynasty.DefaultSettings()
	settings.AutoParse = true
	settings.SchemaDir = schemaDir

	file, err := dynasty.Open(path, &settings)
	if err != nil {
		return friendlyParseError(err, schemaDir)
	}
	exp, err := file.Export()
	if err != nil {
		return friendlyParseError(err, schemaDir)
	}
	if msg := exporter.RosterGroupingError(&exp); msg != "" {
		return fmt.Errorf("%s", msg)
	}
	exporter.Normalize(&exp)

	info, _ := os.Stat(path)
	var size int64
	if info != nil {
		size = info.Size()
	}
	_ = size

	a.mu.Lock()
	a.export = &exp
	a.source = filepath.Base(path)
	a.mu.Unlock()
	return nil
}

// CloseSave clears the in-memory export.
func (a *App) CloseSave() {
	a.mu.Lock()
	a.export = nil
	a.source = ""
	a.mu.Unlock()
}

func (a *App) requireExport() (*dynasty.Export, string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.export == nil {
		return nil, "", fmt.Errorf("no save loaded")
	}
	return a.export, a.source, nil
}

// DashboardCard is one collection card on the dashboard.
type DashboardCard struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Count int    `json:"count"`
	View  string `json:"view"`
}

// DashboardView is the dashboard payload.
type DashboardView struct {
	SourceName string          `json:"sourceName"`
	Cards      []DashboardCard `json:"cards"`
}

// GetDashboard returns collection summary cards.
func (a *App) GetDashboard() (DashboardView, error) {
	exp, source, err := a.requireExport()
	if err != nil {
		return DashboardView{}, err
	}
	cards := make([]DashboardCard, 0, len(exporter.Registry)+1)
	for _, c := range exporter.Registry {
		view := "collection:" + c.Name
		title := c.Title
		switch c.Name {
		case "schoolGrades", "pipelineInfluence", "rivalries", "recruiting":
			continue
		case "recordBook":
			view = "records"
		case "rosters":
			view = "rosters"
		case "leavingPlayers":
			view = "leaving"
		case "recruits":
			view = "recruits"
		case "teams":
			view = "schools"
			title = "Schools"
		}
		cards = append(cards, DashboardCard{
			Name:  c.Name,
			Title: title,
			Count: c.Count(*exp),
			View:  view,
		})
	}
	if len(exp.Recruits) > 0 {
		rankings := exporter.RecruitingRankings(*exp, "")
		cards = append(cards, DashboardCard{
			Name:  "recruitingRankings",
			Title: "Recruiting Rankings",
			Count: rankings.Total,
			View:  "recruiting-rankings",
		})
	}
	return DashboardView{SourceName: source, Cards: cards}, nil
}

// TableView is a generic header+rows table for the UI.
type TableView struct {
	Name   string     `json:"name"`
	Title  string     `json:"title"`
	Header []string   `json:"header"`
	Rows   [][]string `json:"rows"`
}

// GetCollection returns a flattened collection table.
func (a *App) GetCollection(name string) (TableView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return TableView{}, err
	}
	c, ok := exporter.Lookup(name)
	if !ok {
		return TableView{}, fmt.Errorf("unknown collection %q", name)
	}
	return TableView{
		Name:   c.Name,
		Title:  c.Title,
		Header: c.Header,
		Rows:   c.Rows(*exp),
	}, nil
}

// GetSchools returns the schools index.
func (a *App) GetSchools(query string) ([]exporter.SchoolSummary, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return nil, err
	}
	return exporter.SchoolsIndex(*exp, query), nil
}

// GetSchool returns school detail for a team id.
func (a *App) GetSchool(teamID int) (exporter.SchoolDetail, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return exporter.SchoolDetail{}, err
	}
	detail, ok := exporter.SchoolDetailForTeam(*exp, teamID)
	if !ok {
		return exporter.SchoolDetail{}, fmt.Errorf("school %d not found", teamID)
	}
	return detail, nil
}

// GetSchoolClass returns a team's recruiting class view.
func (a *App) GetSchoolClass(teamID int) (exporter.SchoolClassView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return exporter.SchoolClassView{}, err
	}
	view, ok := exporter.SchoolClassForTeam(*exp, teamID)
	if !ok {
		return exporter.SchoolClassView{}, fmt.Errorf("school %d not found", teamID)
	}
	return view, nil
}

// GetRosterTeams returns roster index rows.
func (a *App) GetRosterTeams() ([]exporter.RosterTeamSummary, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return nil, err
	}
	return exporter.RosterTeams(*exp), nil
}

// RosterTeamView is a per-team roster page payload.
type RosterTeamView struct {
	TeamID    int                      `json:"teamId"`
	TeamName  string                   `json:"teamName"`
	Positions []string                 `json:"positions"`
	Players   []dynasty.PlayerExport   `json:"players"`
	Header    []string                 `json:"header"`
}

// GetRoster returns a team's roster.
func (a *App) GetRoster(teamID int) (RosterTeamView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return RosterTeamView{}, err
	}
	r, ok := exporter.FindRoster(*exp, teamID)
	if !ok {
		return RosterTeamView{}, fmt.Errorf("roster for team %d not found", teamID)
	}
	return RosterTeamView{
		TeamID:    r.TeamID,
		TeamName:  r.TeamName,
		Positions: exporter.RosterPositions(r.Players),
		Players:   r.Players,
		Header:    exporter.RosterPositionHeader(),
	}, nil
}

// GetLeavingTeams returns leaving-player team summaries.
func (a *App) GetLeavingTeams() ([]exporter.LeavingTeamSummary, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return nil, err
	}
	return exporter.LeavingTeams(*exp), nil
}

// LeavingTeamView is leaving players for one team.
type LeavingTeamView struct {
	TeamID     int                            `json:"teamId"`
	TeamName   string                         `json:"teamName"`
	LeaveTypes []string                       `json:"leaveTypes"`
	Players    []dynasty.LeavingPlayerExport  `json:"players"`
	Header     []string                       `json:"header"`
	Rows       [][]string                     `json:"rows"`
}

// GetLeavingTeam returns leaving players for one team.
func (a *App) GetLeavingTeam(teamID int) (LeavingTeamView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return LeavingTeamView{}, err
	}
	players := exporter.LeavingPlayersForTeam(*exp, teamID)
	if len(players) == 0 {
		return LeavingTeamView{}, fmt.Errorf("leaving data for team %d not found", teamID)
	}
	return LeavingTeamView{
		TeamID:     teamID,
		TeamName:   exporter.LeavingTeamName(*exp, teamID),
		LeaveTypes: exporter.LeavingLeaveTypes(players),
		Players:    players,
		Header:     exporter.LeavingHeader(),
		Rows:       exporter.LeavingCollectionRows(players),
	}, nil
}

// GetRecruits returns the recruit board table rows.
func (a *App) GetRecruits() (TableView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return TableView{}, err
	}
	targets := exporter.BuildRecruitTargetByRecruitID(*exp)
	shortNames := exporter.BuildTeamShortNameMap(exp.Teams)
	return TableView{
		Name:   "recruits",
		Title:  "Recruits",
		Header: exporter.RecruitUIHeader(),
		Rows:   exporter.RecruitUIRows(exp.Recruits, shortNames, targets),
	}, nil
}

// GetRecruiting returns the recruiting pursuit table.
func (a *App) GetRecruiting() (TableView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return TableView{}, err
	}
	entries := exporter.RecruitingEntries(*exp)
	return TableView{
		Name:   "recruiting",
		Title:  "Recruiting",
		Header: exporter.RecruitingUIHeader(),
		Rows:   exporter.RecruitingUIRows(entries),
	}, nil
}

// GetRecruitingRankings returns team recruiting rankings.
func (a *App) GetRecruitingRankings(conf string) (exporter.RecruitingRankingsView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return exporter.RecruitingRankingsView{}, err
	}
	return exporter.RecruitingRankings(*exp, conf), nil
}

// GetPlayer returns player/recruit detail.
func (a *App) GetPlayer(playerID int) (map[string]any, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return nil, err
	}
	lookup, ok := exporter.FindPlayer(*exp, playerID)
	if !ok {
		return nil, fmt.Errorf("player %d not found", playerID)
	}
	isRecruit := lookup.Recruit != nil
	out := map[string]any{
		"player":    lookup.Player,
		"teamId":    lookup.TeamID,
		"teamName":  lookup.TeamName,
		"isRecruit": isRecruit,
		"recruit":   lookup.Recruit,
		"wear":      exporter.WearRows(&lookup.Player),
		"skillGroups": exporter.SkillGroupRows(&lookup.Player),
	}
	if isRecruit {
		if target, ok := exporter.FindRecruitingTarget(*exp, lookup.Recruit.ID); ok {
			out["recruiting"] = target
			shortNames := exporter.BuildTeamShortNameMap(exp.Teams)
			out["schoolInterest"] = exporter.SchoolInterestRows(target.SchoolInterest, shortNames)
			out["pitches"] = exporter.ActivePitchRows(target.ActivePitches)
		} else if len(lookup.Recruit.SchoolInterest) > 0 {
			shortNames := exporter.BuildTeamShortNameMap(exp.Teams)
			out["schoolInterest"] = exporter.SchoolInterestRows(lookup.Recruit.SchoolInterest, shortNames)
		}
	}
	return out, nil
}

// GetGame returns a game by id.
func (a *App) GetGame(gameID int) (dynasty.GameExport, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return dynasty.GameExport{}, err
	}
	for _, g := range exp.Games {
		if int(g.ID) == gameID {
			return g, nil
		}
	}
	return dynasty.GameExport{}, fmt.Errorf("game %d not found", gameID)
}

// RecordPeriodTab is one record-book period tab.
type RecordPeriodTab struct {
	Period string `json:"period"`
	Title  string `json:"title"`
	Count  int    `json:"count"`
}

// GetRecordPeriods returns available record-book periods.
func (a *App) GetRecordPeriods() ([]RecordPeriodTab, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return nil, err
	}
	tabs := make([]RecordPeriodTab, 0, len(exporter.RecordBookPeriods))
	for _, p := range exporter.RecordBookPeriods {
		tabs = append(tabs, RecordPeriodTab{
			Period: p,
			Title:  exporter.RecordBookPeriodTitle(p),
			Count:  exporter.RecordBookCount(*exp, p),
		})
	}
	return tabs, nil
}

// GetRecords returns filtered record-book rows.
func (a *App) GetRecords(period, team, statType string) (TableView, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return TableView{}, err
	}
	filter := exporter.RecordBookFilter{Period: period, Team: team, StatType: statType}
	return TableView{
		Name:   "recordBook",
		Title:  "Record Book",
		Header: exporter.RecordBookHeader,
		Rows:   exporter.RecordBookRows(*exp, filter),
	}, nil
}

func (a *App) writeExport(defaultName, filter string, write func(*bytes.Buffer) error) (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: filter, Pattern: "*"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	var buf bytes.Buffer
	if err := write(&buf); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// ExportFullJSON writes the full export as JSON via a save dialog.
func (a *App) ExportFullJSON() (string, error) {
	exp, source, err := a.requireExport()
	if err != nil {
		return "", err
	}
	base := strings.TrimSuffix(source, filepath.Ext(source))
	if base == "" {
		base = "dynasty"
	}
	return a.writeExport(base+".json", "JSON (*.json)", func(buf *bytes.Buffer) error {
		return exporter.FullJSONStream(buf, *exp)
	})
}

// ExportFullCSVZip writes all collections as a CSV zip via a save dialog.
func (a *App) ExportFullCSVZip() (string, error) {
	exp, source, err := a.requireExport()
	if err != nil {
		return "", err
	}
	base := strings.TrimSuffix(source, filepath.Ext(source))
	if base == "" {
		base = "dynasty"
	}
	return a.writeExport(base+"-csv.zip", "ZIP (*.zip)", func(buf *bytes.Buffer) error {
		return exporter.FullCSVZip(buf, *exp)
	})
}

// ExportCollectionJSON writes one collection as JSON.
func (a *App) ExportCollectionJSON(name string) (string, error) {
	exp, source, err := a.requireExport()
	if err != nil {
		return "", err
	}
	c, ok := exporter.Lookup(name)
	if !ok {
		return "", fmt.Errorf("unknown collection %q", name)
	}
	base := strings.TrimSuffix(source, filepath.Ext(source))
	if base == "" {
		base = "dynasty"
	}
	return a.writeExport(base+"-"+name+".json", "JSON (*.json)", func(buf *bytes.Buffer) error {
		return exporter.SectionJSONStream(buf, *exp, c)
	})
}

// ExportCollectionCSV writes one collection as CSV.
func (a *App) ExportCollectionCSV(name string) (string, error) {
	exp, source, err := a.requireExport()
	if err != nil {
		return "", err
	}
	c, ok := exporter.Lookup(name)
	if !ok {
		return "", fmt.Errorf("unknown collection %q", name)
	}
	base := strings.TrimSuffix(source, filepath.Ext(source))
	if base == "" {
		base = "dynasty"
	}
	return a.writeExport(base+"-"+name+".csv", "CSV (*.csv)", func(buf *bytes.Buffer) error {
		return exporter.SectionCSV(buf, *exp, c)
	})
}

// ExportSchoolClassCSV exports a school recruiting class as CSV.
func (a *App) ExportSchoolClassCSV(teamID int) (string, error) {
	exp, _, err := a.requireExport()
	if err != nil {
		return "", err
	}
	view, ok := exporter.SchoolClassForTeam(*exp, teamID)
	if !ok {
		return "", fmt.Errorf("school %d not found", teamID)
	}
	name := "school-" + fmt.Sprintf("%d", teamID) + "-class.csv"
	return a.writeExport(name, "CSV (*.csv)", func(buf *bytes.Buffer) error {
		return exporter.WriteCSV(buf, exporter.SchoolClassHeader(), exporter.SchoolClassRows(view))
	})
}
