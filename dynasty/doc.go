// Package dynasty reads EA Sports College Football dynasty save files and
// exposes parsed league data for export.
//
// The save format is under active reverse engineering. Most parsing methods
// return [ErrNotImplemented] until the CFB 27 container and table layout are
// mapped. [Open] and [File.Inspect] work today for container-level analysis.
package dynasty
