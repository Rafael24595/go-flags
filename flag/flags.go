package flag

// BoolFlag creates a flag that parses boolean values.
//
// It is equivalent to calling New with BoolParser.
func BoolFlag(description string, names ...string) *Flag[bool] {
	return New(BoolParser, "bool", description, names...)
}

// StringFlag creates a flag that parses string values.
//
// It is equivalent to calling New with StringParser.
func StringFlag(desc string, names ...string) *Flag[string] {
	return New(StringParser, "string", desc, names...)
}

// IntFlag creates a flag that parses integer values.
//
// It is equivalent to calling New with IntParser.
func IntFlag(desc string, names ...string) *Flag[int] {
	return New(IntParser, "int", desc, names...)
}

// Int64Flag creates a flag that parses 64-bit integer values.
//
// It is equivalent to calling New with Int64Parser.
func Int64Flag(description string, names ...string) *Flag[int64] {
	return New(Int64Parser, "int64", description, names...)
}

// UintFlag creates a flag that parses unsigned integer values.
//
// It is equivalent to calling New with UintParser.
func UintFlag(desc string, names ...string) *Flag[uint] {
	return New(UintParser, "uint", desc, names...)
}

// Uint64Flag creates a flag that parses 64-bit unsigned integer values.
//
// It is equivalent to calling New with Uint64Parser.
func Uint64Flag(description string, names ...string) *Flag[uint64] {
	return New(Uint64Parser, "uint64", description, names...)
}

// Float64Flag creates a flag that parses 64-bit floating-point values.
//
// It is equivalent to calling New with Float64Parser.
func Float64Flag(description string, names ...string) *Flag[float64] {
	return New(Float64Parser, "float64", description, names...)
}
