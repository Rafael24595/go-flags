package flag

// CLIInfo contains the metadata used to format a command-line interface.
type CLIInfo struct {
	// Flags is a list of all registered command-line options.
	Flags []Info
	// Adapters is a list of all registered command-line adapters.
	Adapters []string
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

	// Void indicates whether the option is a void flag that does not require an argument.
	Void bool
	// Required indicates whether the option is required.
	Required bool

	// UndefinedDefault is the default value used when the option is not provided.
	UndefinedDefault DefaultInfo
	// OptionalDefault is the default value used when the option is provided without
	OptionalDefault DefaultInfo
}
