package flag

// VoidFlag creates a flag that does not require an argument.
//
// When the flag is present, its value is the empty string.
//
// It is equivalent to calling New with VoidParser and marking the flag
// as void.
func VoidFlag(description string, names ...string) *Flag[string] {
	return New(TypeVoid, VoidParser, description, names...).Void()
}

// BoolFlag creates a flag that parses boolean values.
//
// It is equivalent to calling New with BoolParser.
func BoolFlag(description string, names ...string) *Flag[bool] {
	return New(TypeBool, BoolParser, description, names...)
}

// StringFlag creates a flag that parses string values.
//
// It is equivalent to calling New with StringParser.
func StringFlag(desc string, names ...string) *Flag[string] {
	return New(TypeString, StringParser, desc, names...)
}

// IntFlag creates a flag that parses integer values.
//
// It is equivalent to calling New with IntParser.
func IntFlag(desc string, names ...string) *Flag[int] {
	return New(TypeInt, IntParser, desc, names...)
}

// Int64Flag creates a flag that parses 64-bit integer values.
//
// It is equivalent to calling New with Int64Parser.
func Int64Flag(description string, names ...string) *Flag[int64] {
	return New(TypeInt64, Int64Parser, description, names...)
}

// UintFlag creates a flag that parses unsigned integer values.
//
// It is equivalent to calling New with UintParser.
func UintFlag(desc string, names ...string) *Flag[uint] {
	return New(TypeUint, UintParser, desc, names...)
}

// Uint64Flag creates a flag that parses 64-bit unsigned integer values.
//
// It is equivalent to calling New with Uint64Parser.
func Uint64Flag(description string, names ...string) *Flag[uint64] {
	return New(TypeUint64, Uint64Parser, description, names...)
}

// Float64Flag creates a flag that parses 64-bit floating-point values.
//
// It is equivalent to calling New with Float64Parser.
func Float64Flag(description string, names ...string) *Flag[float64] {
	return New(TypeFloat64, Float64Parser, description, names...)
}
