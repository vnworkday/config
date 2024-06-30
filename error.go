package config

import "github.com/pkg/errors"

var (
	// ErrInvalidTarget is returned when the target is not a struct pointer.
	ErrInvalidTarget = errors.New("target must be a struct pointer")

	// ErrUnsupportedField is returned when a field type is not supported.
	ErrUnsupportedField = errors.New("unsupported field type")
)
