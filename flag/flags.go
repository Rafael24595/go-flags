package flag

// VoidFlag creates a flag that does not require an argument.
//
// When the flag is present, its value is the empty string.
//
// It is equivalent to calling New with VoidParser and marking the flag
// as void.
func VoidFlag(description string, names ...string) *Flag[string] {
	return NewFlag(TypeVoid, VoidParser, description, names...).Void()
}

// BoolFlag creates a flag that parses boolean values.
//
// It is equivalent to calling New with BoolParser.
func BoolFlag(description string, names ...string) *Flag[bool] {
	return NewFlag(TypeBool, BoolParser, description, names...)
}

// StringFlag creates a flag that parses string values.
//
// It is equivalent to calling New with StringParser.
func StringFlag(desc string, names ...string) *Flag[string] {
	return NewFlag(TypeString, StringParser, desc, names...)
}

// IntFlag creates a flag that parses integer values.
//
// It is equivalent to calling New with IntParser.
func IntFlag(desc string, names ...string) *Flag[int] {
	return NewFlag(TypeInt, IntParser, desc, names...)
}

// Int64Flag creates a flag that parses 64-bit integer values.
//
// It is equivalent to calling New with Int64Parser.
func Int64Flag(description string, names ...string) *Flag[int64] {
	return NewFlag(TypeInt64, Int64Parser, description, names...)
}

// UintFlag creates a flag that parses unsigned integer values.
//
// It is equivalent to calling New with UintParser.
func UintFlag(desc string, names ...string) *Flag[uint] {
	return NewFlag(TypeUint, UintParser, desc, names...)
}

// Uint64Flag creates a flag that parses 64-bit unsigned integer values.
//
// It is equivalent to calling New with Uint64Parser.
func Uint64Flag(description string, names ...string) *Flag[uint64] {
	return NewFlag(TypeUint64, Uint64Parser, description, names...)
}

// Float64Flag creates a flag that parses 64-bit floating-point values.
//
// It is equivalent to calling New with Float64Parser.
func Float64Flag(description string, names ...string) *Flag[float64] {
	return NewFlag(TypeFloat64, Float64Parser, description, names...)
}

// SliceBoolFlag creates a flag that parses a list of boolean values.
func SliceBoolFlag(description string, names ...string) *SliceFlag[bool] {
	return NewSliceFlag(TypeBool, BoolParser, description, names...)
}

// SliceStringFlag creates a flag that parses a list of string values.
func SliceStringFlag(description string, names ...string) *SliceFlag[string] {
	return NewSliceFlag(TypeString, StringParser, description, names...)
}

// SliceIntFlag creates a flag that parses a list of integer values.
func SliceIntFlag(description string, names ...string) *SliceFlag[int] {
	return NewSliceFlag(TypeInt, IntParser, description, names...)
}

// SliceInt64Flag creates a flag that parses a list of 64-bit integer values.
func SliceInt64Flag(description string, names ...string) *SliceFlag[int64] {
	return NewSliceFlag(TypeInt64, Int64Parser, description, names...)
}

// SliceUintFlag creates a flag that parses a list of unsigned integer values.
func SliceUintFlag(description string, names ...string) *SliceFlag[uint] {
	return NewSliceFlag(TypeUint, UintParser, description, names...)
}

// SliceUint64Flag creates a flag that parses a list of 64-bit unsigned integer values.
func SliceUint64Flag(description string, names ...string) *SliceFlag[uint64] {
	return NewSliceFlag(TypeUint64, Uint64Parser, description, names...)
}

// SliceFloat64Flag creates a flag that parses a list of 64-bit floating-point values.
func SliceFloat64Flag(description string, names ...string) *SliceFlag[float64] {
	return NewSliceFlag(TypeFloat64, Float64Parser, description, names...)
}
