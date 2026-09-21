package flag

import (
	"errors"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewWithoutNames(t *testing.T) {
	assert.Panic(t, func() {
		New(TypeInt, IntParser, "test")
	})
}

func TestFlagCustomParser(t *testing.T) {
	expected := errors.New("parse error")

	parser := func(string) (int, error) {
		return 0, expected
	}

	flag := New(TypeUnknown, parser, "test", "-t")

	err := flag.parse("value")

	assert.False(t, flag.isSet())
	assert.ErrorIs(t, expected, err)
}

func TestFlagDefaultOptional(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		DefaultOptional(10)

	assert.True(t, flag.isOptional())

	flag.applyOptionalDefault()

	assert.True(t, flag.isSet())
	assert.Equal(t, 10, flag.Value())
}

func TestFlagDefaultUndefined(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		DefaultUndefined(-1)

	flag.applyUndefinedDefault()

	assert.True(t, flag.isSet())
	assert.Equal(t, -1, flag.Value())
}

func TestFlagRequired(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		Required()

	assert.True(t, flag.isRequired())
}

func TestFlagParse(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t")

	assert.Nil(t, flag.parse("42"))

	assert.True(t, flag.isSet())
	assert.Equal(t, 42, flag.Value())
}

func TestFlagParseError(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t")

	assert.NotNil(t, flag.parse("invalid"))
	assert.False(t, flag.isSet())
}

func TestFlagAliases(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t", "--test")

	got := flag.aliases()

	assert.Size(t, 2, got)
	assert.Equal(t, "-t", got[0])
	assert.Equal(t, "--test", got[1])
}
