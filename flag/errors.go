package flag

import "errors"

var (
	// ErrUnknownOption indicates that an unrecognized option was provided.
	ErrUnknownOption = errors.New("unknown option")
	// ErrUnexpectedValue indicates that a void option was given a value.
	ErrUnexpectedValue = errors.New("option does not accept a value")
	// ErrMissingValue indicates that an option requires a value.
	ErrMissingValue = errors.New("option requires a value")
	// ErrRequiredOption indicates that a required option was not provided.
	ErrRequiredOption = errors.New("required option was not provided")
	// ErrInvalidValue indicates that an option value could not be parsed.
	ErrInvalidValue = errors.New("invalid option value")
)
