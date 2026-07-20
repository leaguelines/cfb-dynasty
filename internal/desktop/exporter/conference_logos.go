package exporter

import "strings"

// conferenceLogoFiles lists dynasty conference names with a vendored PNG in
// web/static/conferences/. Filenames use spaces replaced by hyphens (e.g. "Big 12" → Big-12.png).
var conferenceLogoFiles = map[string]struct{}{
	"ACC":         {},
	"American":    {},
	"Big 12":      {},
	"Big Ten":     {},
	"CUSA":        {},
	"Independent": {},
	"MAC":         {},
	"MWC":         {},
	"Pac-12":      {},
	"SEC":         {},
	"Sun Belt":    {},
}

// ConferenceLogoPath returns the static URL for a vendored conference logo, or "" if unavailable.
func ConferenceLogoPath(conference string) string {
	conf := strings.TrimSpace(conference)
	if conf == "" {
		return ""
	}
	if _, ok := conferenceLogoFiles[conf]; !ok {
		return ""
	}
	return "/static/conferences/" + conferenceLogoFilename(conf)
}

func conferenceLogoFilename(conference string) string {
	return strings.ReplaceAll(conference, " ", "-") + ".png"
}
