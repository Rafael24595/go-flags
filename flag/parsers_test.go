package flag

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestVoidParser(t *testing.T) {
	tests := []string{
		"",
		"true",
		"false",
		"anything",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			value, err := VoidParser(input)

			assert.Nil(t, err)
			assert.Equal(t, input, value)
		})
	}
}

func TestBoolParser(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: "true", want: true},
		{input: "false", want: false},
	}

	for _, tt := range tests {
		value, err := BoolParser(tt.input)

		assert.Nil(t, err)
		assert.Equal(t, tt.want, value)
	}
}

func TestBoolParserError(t *testing.T) {
	_, err := BoolParser("invalid")

	assert.NotNil(t, err)
}

func TestStringParser(t *testing.T) {
	value, err := StringParser("hello")

	assert.Nil(t, err)
	assert.Equal(t, "hello", value)
}

func TestIntParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "positive",
			input: "42",
			want:  42,
		},
		{
			name:  "negative",
			input: "-42",
			want:  -42,
		},
		{
			name:  "zero",
			input: "0",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := IntParser(tt.input)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, value)
		})
	}
}

func TestIntParserError(t *testing.T) {
	_, err := IntParser("invalid")

	assert.NotNil(t, err)
}

func TestUintParser(t *testing.T) {
	value, err := UintParser("42")

	assert.Nil(t, err)
	assert.Equal(t, 42, value)
}

func TestUintParserError(t *testing.T) {
	_, err := UintParser("-1")

	assert.NotNil(t, err)
}

func TestInt64Parser(t *testing.T) {
	value, err := Int64Parser("9223372036854775807")

	assert.Nil(t, err)
	assert.Equal(t, 9223372036854775807, value)
}

func TestInt64ParserError(t *testing.T) {
	_, err := Int64Parser("9223372036854775808")

	assert.NotNil(t, err)
}

func TestUint64Parser(t *testing.T) {
	value, err := Uint64Parser("18446744073709551615")

	assert.Nil(t, err)
	assert.Equal(t, 18446744073709551615, value)
}

func TestUint64ParserError(t *testing.T) {
	_, err := Uint64Parser("18446744073709551616")

	assert.NotNil(t, err)
}

func TestFloat64Parser(t *testing.T) {
	value, err := Float64Parser("3.14159")

	assert.Nil(t, err)
	assert.Equal(t, 3.14159, value)
}

func TestFloat64ParserError(t *testing.T) {
	_, err := Float64Parser("invalid")

	assert.NotNil(t, err)
}
