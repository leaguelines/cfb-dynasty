package dynasty

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var schemaFilePattern = regexp.MustCompile(`^(?i)(?:[mc](\d+)_)?(\d+)_(\d+)\.gz$`)

type schemaCandidate struct {
	Path     string
	GameYear int
	Major    int
	Minor    int
}

func pickSchemaFile(dir string, version SchemaVersion) (string, error) {
	candidates, err := listSchemaCandidates(dir)
	if err != nil {
		return "", err
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("%w: no .gz schema files in %q", ErrSchemaNotFound, dir)
	}

	best := selectSchemaCandidate(candidates, version.Major, version.Minor)
	if best == nil {
		return "", fmt.Errorf("%w: no matching schema in %q for major=%d minor=%d", ErrSchemaNotFound, dir, version.Major, version.Minor)
	}
	return best.Path, nil
}

func listSchemaCandidates(dir string) ([]schemaCandidate, error) {
	var candidates []schemaCandidate
	var walk func(string) error
	walk = func(current string) error {
		files, err := os.ReadDir(current)
		if err != nil {
			return err
		}
		for _, entry := range files {
			path := filepath.Join(current, entry.Name())
			if entry.IsDir() {
				if err := walk(path); err != nil {
					return err
				}
				continue
			}
			if !strings.EqualFold(filepath.Ext(entry.Name()), ".gz") {
				continue
			}
			candidate, ok := parseSchemaFilename(path, entry.Name())
			if !ok {
				continue
			}
			candidates = append(candidates, candidate)
		}
		return nil
	}
	if err := walk(dir); err != nil {
		return nil, fmt.Errorf("cfb-dynasty: read schema dir %q: %w", dir, err)
	}
	return candidates, nil
}

func parseSchemaFilename(path, name string) (schemaCandidate, bool) {
	matches := schemaFilePattern.FindStringSubmatch(name)
	if matches == nil {
		return schemaCandidate{}, false
	}
	gameYear := 0
	if matches[1] != "" {
		gameYear = atoi(matches[1])
	}
	return schemaCandidate{
		Path:     path,
		GameYear: gameYear,
		Major:    atoi(matches[2]),
		Minor:    atoi(matches[3]),
	}, true
}

func selectSchemaCandidate(candidates []schemaCandidate, major, minor int) *schemaCandidate {
	if len(candidates) == 1 {
		return &candidates[0]
	}

	exactMajor := filterCandidates(candidates, func(c schemaCandidate) bool {
		return c.Major == major
	})
	if len(exactMajor) > 0 {
		return pickClosestMinor(exactMajor, minor)
	}

	closestMajor := closestInt(uniqueMajors(candidates), major)
	majorMatches := filterCandidates(candidates, func(c schemaCandidate) bool {
		return c.Major == closestMajor
	})
	if len(majorMatches) == 0 {
		return nil
	}
	return pickClosestMinor(majorMatches, minor)
}

func pickClosestMinor(candidates []schemaCandidate, minor int) *schemaCandidate {
	if len(candidates) == 0 {
		return nil
	}
	best := candidates[0]
	bestDelta := abs(best.Minor - minor)
	for i := 1; i < len(candidates); i++ {
		delta := abs(candidates[i].Minor - minor)
		if delta < bestDelta {
			best = candidates[i]
			bestDelta = delta
		}
	}
	return &best
}

func uniqueMajors(candidates []schemaCandidate) []int {
	seen := make(map[int]struct{}, len(candidates))
	var majors []int
	for _, c := range candidates {
		if _, ok := seen[c.Major]; ok {
			continue
		}
		seen[c.Major] = struct{}{}
		majors = append(majors, c.Major)
	}
	return majors
}

func closestInt(values []int, goal int) int {
	if len(values) == 0 {
		return goal
	}
	best := values[0]
	for _, v := range values[1:] {
		if abs(v-goal) < abs(best-goal) {
			best = v
		}
	}
	return best
}

func filterCandidates(candidates []schemaCandidate, keep func(schemaCandidate) bool) []schemaCandidate {
	var out []schemaCandidate
	for _, c := range candidates {
		if keep(c) {
			out = append(out, c)
		}
	}
	return out
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
