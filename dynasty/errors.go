package dynasty

import "errors"

var (
	// ErrNotImplemented is returned when a feature depends on format details
	// that have not been reverse-engineered yet.
	ErrNotImplemented = errors.New("cfb-dynasty: not implemented")

	// ErrInvalidSave is returned when a file does not appear to be a dynasty save.
	ErrInvalidSave = errors.New("cfb-dynasty: invalid save file")

	// ErrNotLoaded is returned when tables or records are accessed before Parse.
	ErrNotLoaded = errors.New("cfb-dynasty: file not parsed")

	// ErrSchemaRequired is returned when schema data is missing for parsing.
	ErrSchemaRequired = errors.New("cfb-dynasty: schema required")
)

// IsNotImplemented reports whether err indicates unimplemented functionality.
func IsNotImplemented(err error) bool {
	return errors.Is(err, ErrNotImplemented)
}
