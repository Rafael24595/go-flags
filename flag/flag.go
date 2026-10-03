package flag

import "fmt"

type flag interface {
	info() Info
	aliases() []string
	reset()
	process(value *string) (bool, error)
	validate() error
}

// Flag represents a command-line flag and its configuration.
//
// A flag can define different default values depending on whether it is
// omitted entirely or provided without an explicit value.
type Flag[T any] struct {
	parser Parser[T]

	typeName    TypeName
	names       []string
	description string

	void bool

	optionalDefault    T
	hasOptionalDefault bool

	undefinedDefault    T
	hasUndefinedDefault bool

	required bool

	present bool
	value   T
	set     bool
}

// NewFlag creates a new command-line flag using the given parser.
//
// At least one name must be provided. Multiple names can be used to define
// aliases for the same flag.
//
// For example:
//
// cover := NewFlag(IntParser, "Cover index", "-c", "--cover")
func NewFlag[T any](
	typeName TypeName,
	parser Parser[T],
	desc string,
	names ...string,
) *Flag[T] {
	if len(names) == 0 {
		panic("flag requires at least one name")
	}

	return &Flag[T]{
		parser:      parser,
		typeName:    typeName,
		names:       names,
		description: desc,
	}
}

// Void marks the flag as an option that does not accept an argument.
//
// Optional and undefined default values are cleared, and the flag is
// no longer required.
func (f *Flag[T]) Void() *Flag[T] {
	f.void = true

	var zero T

	f.undefinedDefault = zero
	f.hasUndefinedDefault = false

	f.optionalDefault = zero
	f.hasOptionalDefault = false

	f.required = false

	return f
}

// DefaultUndefined sets the value used when the flag is not provided.
//
// This default is applied when the flag is completely absent from the
// command line.
//
// For example:
//
//	cover := New(IntParser, "Cover index", "-c").
//		DefaultUndefined(-1)
//
// Running the program without -c results in a value of -1.
func (f *Flag[T]) DefaultUndefined(value T) *Flag[T] {
	f.undefinedDefault = value
	f.hasUndefinedDefault = true

	f.void = false

	return f
}

// DefaultOptional sets the value used when the flag is provided without
// an explicit value.
//
// For example:
//
//	cover := New(IntParser, "Cover index", "-c").
//		DefaultOptional(0)
//
// Running the program with -c results in a value of 0, while -c 2 results
// in a value of 2.
func (f *Flag[T]) DefaultOptional(value T) *Flag[T] {
	f.optionalDefault = value
	f.hasOptionalDefault = true

	f.void = false

	return f
}

// Required marks the flag as requiring a value when it is provided.
func (f *Flag[T]) Required() *Flag[T] {
	f.required = true
	f.void = false
	return f
}

// IsPresent returns true if the flag was provided on the command line.
func (f *Flag[T]) IsPresent() bool {
	return f.present
}

// Value returns the parsed value of the flag.
func (f *Flag[T]) Value() T {
	return f.value
}

func (f *Flag[T]) info() Info {
	names := make([]string, len(f.names))
	copy(names, f.names)

	return Info{
		Names:       names,
		Type:        f.typeName,
		Description: f.description,
		Void:        f.void,
		Required:    f.required,

		OptionalDefault: DefaultInfo{
			Value: f.optionalDefault,
			Set:   f.hasOptionalDefault,
		},

		UndefinedDefault: DefaultInfo{
			Value: f.undefinedDefault,
			Set:   f.hasUndefinedDefault,
		},
	}
}

func (f *Flag[T]) aliases() []string {
	return f.names
}

func (f *Flag[T]) reset() {
	var zero T
	f.present = false
	f.value = zero
	f.set = false
}

func (f *Flag[T]) process(value *string) (bool, error) {
	if f.present {
		return false, ErrDuplicateOption
	}

	f.present = true

	if value == nil {
		if f.void {
			return false, nil
		}

		if f.hasOptionalDefault {
			f.setValue(f.optionalDefault)
			return false, nil
		}

		return false, ErrMissingValue
	}

	if f.void {
		return false, ErrUnexpectedValue
	}

	parsed, err := f.parser(*value)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrInvalidValue, err)
	}

	f.setValue(parsed)

	return true, nil
}

func (f *Flag[T]) validate() error {
	if f.required && !f.present {
		return ErrRequiredOption
	}

	if !f.set && f.hasUndefinedDefault {
		f.setValue(f.undefinedDefault)
	}

	return nil
}

func (f *Flag[T]) setValue(value T) {
	f.value = value
	f.set = true
}

// SliceFlag represents a command-line flag that accepts multiple values
// and its configuration.
//
// The values are parsed using the provided parser and separated by a
// specified separator character. The default separator is a comma (',').
//
// A SliceFlag can define different default values depending on whether it is
// omitted entirely or provided without an explicit value.
type SliceFlag[T any] struct {
	Flag[[]T]
	separator  rune
	baseParser Parser[T]
}

// NewSliceFlag creates a new command-line flag that accepts multiple values
// using the given parser.
//
// At least one name must be provided. Multiple names can be used to define
// aliases for the same flag.
//
// For example:
//
// covers := NewSliceFlag(IntParser, "Cover indices", "-c", "--covers")
func NewSliceFlag[T any](
	typeName TypeName,
	parser Parser[T],
	desc string,
	names ...string,
) *SliceFlag[T] {
	if len(names) == 0 {
		panic("flag requires at least one name")
	}

	separator := DefaultSeparator

	return &SliceFlag[T]{
		Flag: Flag[[]T]{
			parser:      SliceParser(parser, separator),
			typeName:    typeName.Slice(),
			names:       names,
			description: desc,
		},
		separator:  separator,
		baseParser: parser,
	}
}

// Separator sets the character used to separate multiple values for the flag.
//
// For example:
//
// covers := NewSliceFlag(IntParser, "Cover indices", "-c", "--covers").
//     Separator(';')
//
// Running the program with -c 1;2;3 results in a value of []int{1, 2, 3}.
func (f *SliceFlag[T]) Separator(sep rune) *SliceFlag[T] {
	f.separator = sep
	f.parser = SliceParser(f.baseParser, sep)
	return f
}

func (f *SliceFlag[T]) info() Info {
	info := f.Flag.info()
	info.Description = f.description()
	return info
}

func (f *SliceFlag[T]) description() string {
	return fmt.Sprintf("%s (separator: %q)", f.Flag.description, f.separator)
}

func (f *SliceFlag[T]) process(value *string) (bool, error) {
	f.present = true

	if value == nil {
		return false, ErrMissingValue
	}

	parsed, err := f.parser(*value)
	if err != nil {
		return false, err
	}

	f.setValue(
		append(f.value, parsed...),
	)

	return true, nil
}
