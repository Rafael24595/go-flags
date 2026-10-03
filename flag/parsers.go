package flag

import (
	"encoding/csv"
	"strconv"
	"strings"
)

const (
	// DefaultSeparator is the default character used to separate multiple values
	DefaultSeparator = ','
)

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

func SliceParser[T any](parser Parser[T], separator rune) Parser[[]T] {
	return func(value string) ([]T, error) {
		value = strings.TrimSpace(value)
		if value == "" {
			return []T{}, nil
		}

		reader := csv.NewReader(
			strings.NewReader(value),
		)

		reader.Comma = separator
		reader.LazyQuotes = true
		reader.TrimLeadingSpace = true

		records, err := reader.Read()
		if err != nil {
			return nil, err
		}

		result := make([]T, 0, len(records))

		for _, part := range records {
			cleaned := strings.TrimSpace(part)

			parsed, err := parser(cleaned)
			if err != nil {
				return nil, err
			}

			result = append(result, parsed)
		}

		return result, nil
	}
}
