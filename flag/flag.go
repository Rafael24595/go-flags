package flag

type flag interface {
	info() Info
	aliases() []string
	isSet() bool
	isPresent() bool
	isOptional() bool
	isRequired() bool
	markPresent()
	parse(string) error
	applyOptionalDefault()
	applyUndefinedDefault()
}

// DefaultInfo describes a default value associated with an option.
type DefaultInfo struct {
	// Value is the default value.
	Value any
	// Set indicates whether a default value is defined.
	Set bool
}

// Info contains the metadata used to format a command-line option.
type Info struct {
	// Names is a list of names and aliases for the option.
	Names []string
	// Type is the name of the value type for the option.
	Type TypeName
	// Description is a human-readable description of the option.
	Description string

	// Required indicates whether the option is required.
	Required bool

	// OptionalDefault is the default value used when the option is provided without
	OptionalDefault  DefaultInfo
	// UndefinedDefault is the default value used when the option is not provided.
	UndefinedDefault DefaultInfo
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

	optionalDefault    T
	hasOptionalDefault bool

	undefinedDefault    T
	hasUndefinedDefault bool

	required bool

	present bool
	value   T
	set     bool
}

// New creates a new command-line flag using the given parser.
//
// At least one name must be provided. Multiple names can be used to define
// aliases for the same flag.
//
// For example:
//
// cover := New(IntParser, "Cover index", "-c", "--cover")
func New[T any](
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
	return f
}

// Required marks the flag as requiring a value when it is provided.
func (f *Flag[T]) Required() *Flag[T] {
	f.required = true
	return f
}

// IsPresent returns true if the flag was provided on the command line.
func (f *Flag[T]) IsPresent() bool {
	return f.isPresent()
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

func (f *Flag[T]) isSet() bool {
	return f.set
}

func (f *Flag[T]) isPresent() bool {
	return f.present
}

func (f *Flag[T]) isOptional() bool {
	return f.hasOptionalDefault
}

func (f *Flag[T]) isRequired() bool {
	return f.required
}

func (f *Flag[T]) markPresent() {
	f.present = true
}

func (f *Flag[T]) parse(value string) error {
	parsed, err := f.parser(value)
	if err != nil {
		return err
	}

	f.value = parsed
	f.set = true

	return nil
}

func (f *Flag[T]) applyOptionalDefault() {
	if f.hasOptionalDefault {
		f.value = f.optionalDefault
		f.set = true
	}
}

func (f *Flag[T]) applyUndefinedDefault() {
	if f.hasUndefinedDefault {
		f.value = f.undefinedDefault
		f.set = true
	}
}
