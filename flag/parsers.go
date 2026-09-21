package flag

import "strconv"

// Parser converts a command-line argument from a string into a value of type T.
type Parser[T any] func(string) (T, error)

// VoidParser returns the input string unchanged.
//
// It can be used by options that do not require an argument.
func VoidParser(value string) (string, error) {
	return StringParser(value)
}

// BoolParser parses a boolean value.
//
// Accepted values are those supported by strconv.ParseBool.
func BoolParser(value string) (bool, error) {
	return strconv.ParseBool(value)
}

// StringParser returns the input string unchanged.
func StringParser(value string) (string, error) {
	return value, nil
}

// IntParser parses a base-10 signed integer.
func IntParser(value string) (int, error) {
	return strconv.Atoi(value)
}

// Int64Parser parses a base-10 signed 64-bit integer.
func Int64Parser(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// UintParser parses a base-10 unsigned integer.
func UintParser(value string) (uint, error) {
	n, err := strconv.ParseUint(value, 10, strconv.IntSize)
	return uint(n), err
}

// Uint64Parser parses a base-10 unsigned 64-bit integer.
func Uint64Parser(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

// Float64Parser parses a base-10 floating-point number.
func Float64Parser(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
