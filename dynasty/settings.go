package dynasty

// Settings configures how a dynasty save is opened and parsed.
type Settings struct {
	// AutoParse calls Parse during Open when true.
	AutoParse bool

	// GameYearOverride skips automatic year detection (default target: CFB 27).
	GameYearOverride *int

	// SchemaDir is the directory containing extracted schema XML/JSON files.
	// Required once table parsing is implemented.
	SchemaDir string

	// SchemaOverride pins major/minor schema versions instead of reading them
	// from the save header.
	SchemaOverride *SchemaVersion

	// TuningPath is the path to dynasty-tuning-binary.FTC from the game install.
	// When empty, Open discovers cfb27-db-data/*/dynasty-tuning-binary.FTC under SchemaDir.
	TuningPath string
}

// DefaultSettings returns settings suitable for read-only inspection.
func DefaultSettings() Settings {
	year := DefaultGameYear
	return Settings{
		GameYearOverride: &year,
	}
}

// GameYear returns the configured or detected game year.
func (s Settings) GameYear() int {
	if s.GameYearOverride != nil {
		return *s.GameYearOverride
	}
	return DefaultGameYear
}
