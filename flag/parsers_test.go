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

func TestSliceStringParser(t *testing.T) {
	t.Parallel()

	t.Run("basic comma-separated strings", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(StringParser, DefaultSeparator)
		got, err := parser("foo, bar, baz")

		assert.Nil(t, err)
		assert.DeepEqual(t, []string{"foo", "bar", "baz"}, got)
	})

	t.Run("quotes preserving internal commas with leading space", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(StringParser, DefaultSeparator)
		got, err := parser(`foo, "bar, with comma", baz`)

		assert.Nil(t, err)
		assert.DeepEqual(t, []string{"foo", "bar, with comma", "baz"}, got)
	})

	t.Run("quotes preserving internal commas without leading space", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(StringParser, DefaultSeparator)
		got, err := parser(`foo,"bar,with comma",baz`)

		assert.Nil(t, err)
		assert.DeepEqual(t, []string{"foo", "bar,with comma", "baz"}, got)
	})

	t.Run("custom separator", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(StringParser, ';')
		got, err := parser("a; b; c")

		assert.Nil(t, err)
		assert.DeepEqual(t, []string{"a", "b", "c"}, got)
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(StringParser, DefaultSeparator)
		got, err := parser("   ")

		assert.Nil(t, err)
		assert.Size(t, 0, got)
	})
}

func TestSliceIntParser(t *testing.T) {
	t.Parallel()

	t.Run("valid integers with surrounding spaces", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(IntParser, DefaultSeparator)
		got, err := parser("1,  2, 3 , 4")

		assert.Nil(t, err)
		assert.DeepEqual(t, []int{1, 2, 3, 4}, got)
	})

	t.Run("invalid integer error propagation", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(IntParser, DefaultSeparator)
		_, err := parser("1, invalid, 3")

		assert.NotNil(t, err)
	})

	t.Run("custom pipe separator", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(IntParser, '|')
		got, err := parser("10|20|30")

		assert.Nil(t, err)
		assert.DeepEqual(t, []int{10, 20, 30}, got)
	})
}

func TestSliceBoolParser(t *testing.T) {
	t.Parallel()

	t.Run("valid booleans", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(BoolParser, DefaultSeparator)
		got, err := parser("true, false, 1, 0")

		assert.Nil(t, err)
		assert.DeepEqual(t, []bool{true, false, true, false}, got)
	})

	t.Run("invalid boolean error propagation", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(BoolParser, DefaultSeparator)
		_, err := parser("true, not_a_bool")

		assert.NotNil(t, err)
	})
}

func TestSliceFloat64Parser(t *testing.T) {
	t.Parallel()

	t.Run("valid floats", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(Float64Parser, DefaultSeparator)
		got, err := parser("1.1, 2.5, 3.1415")

		assert.Nil(t, err)
		assert.DeepEqual(t, []float64{1.1, 2.5, 3.1415}, got)
	})
}

func TestSliceNumericTypes(t *testing.T) {
	t.Parallel()

	t.Run("int64 parser", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(Int64Parser, DefaultSeparator)
		got, err := parser("100, 200")

		assert.Nil(t, err)
		assert.DeepEqual(t, []int64{100, 200}, got)
	})

	t.Run("uint parser", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(UintParser, DefaultSeparator)
		got, err := parser("10, 20")

		assert.Nil(t, err)
		assert.DeepEqual(t, []uint{10, 20}, got)
	})

	t.Run("uint64 parser", func(t *testing.T) {
		t.Parallel()

		parser := SliceParser(Uint64Parser, DefaultSeparator)
		got, err := parser("1000, 2000")

		assert.Nil(t, err)
		assert.DeepEqual(t, []uint64{1000, 2000}, got)
	})
}
