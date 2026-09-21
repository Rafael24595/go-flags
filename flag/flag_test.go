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

func TestFlagVoid(t *testing.T) {
	flag := New(TypeVoid, VoidParser, "test", "-v").
		Void()

	assert.True(t, flag.isVoid())
	assert.False(t, flag.isOptional())
	assert.False(t, flag.isRequired())
	assert.False(t, flag.isSet())
	assert.False(t, flag.isPresent())
}

func TestFlagVoidResetsConfiguration(t *testing.T) {
	flag := New(TypeInt, IntParser, "int", "-c").
		DefaultOptional(0).
		DefaultUndefined(-1).
		Required().
		Void()

	assert.True(t, flag.isVoid())
	assert.False(t, flag.isOptional())
	assert.False(t, flag.isRequired())
	assert.False(t, flag.hasOptionalDefault)
	assert.False(t, flag.hasUndefinedDefault)
}

func TestFlagDefaultUndefined(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		DefaultUndefined(-1)

	flag.applyUndefinedDefault()

	assert.True(t, flag.isSet())
	assert.Equal(t, -1, flag.Value())
}

func TestFlagDefaultUndefinedDisablesVoid(t *testing.T) {
	flag := New(TypeVoid, VoidParser, "void", "-v").
		Void().
		DefaultUndefined("")

	assert.False(t, flag.isVoid())
}

func TestFlagDefaultOptional(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		DefaultOptional(10)

	assert.True(t, flag.isOptional())

	flag.applyOptionalDefault()

	assert.True(t, flag.isSet())
	assert.Equal(t, 10, flag.Value())
}

func TestFlagDefaultOptionalDisablesVoid(t *testing.T) {
	flag := New(TypeVoid, VoidParser, "void", "-v").
		Void().
		DefaultOptional("")

	assert.False(t, flag.isVoid())
	assert.True(t, flag.isOptional())
}

func TestFlagRequired(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t").
		Required()

	assert.True(t, flag.isRequired())
}

func TestFlagRequiredDisablesVoid(t *testing.T) {
	flag := New(TypeVoid, VoidParser, "void", "-v").
		Void().
		Required()

	assert.False(t, flag.isVoid())
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

func TestFlagInfo(t *testing.T) {
	flag := New(TypeInt, IntParser, "Cover index", "-c", "--cover").
		DefaultOptional(0).
		DefaultUndefined(-1).
		Required()

	info := flag.info()

	assert.DeepEqual(t, []string{"-c", "--cover"}, info.Names)
	assert.Equal(t, TypeInt, info.Type)
	assert.Equal(t, "Cover index", info.Description)
	assert.False(t, info.Void)
	assert.True(t, info.Required)

	assert.True(t, info.OptionalDefault.Set)
	assert.Equal(t, 0, info.OptionalDefault.Value)

	assert.True(t, info.UndefinedDefault.Set)
	assert.Equal(t, -1, info.UndefinedDefault.Value)
}

func TestVoidFlagInfo(t *testing.T) {
	flag := VoidFlag("Enable verbose output", "-v", "--verbose")

	info := flag.info()

	assert.DeepEqual(t, []string{"-v", "--verbose"}, info.Names)
	assert.Equal(t, TypeVoid, info.Type)
	assert.Equal(t, "Enable verbose output", info.Description)
	assert.True(t, info.Void)
	assert.False(t, info.Required)

	assert.False(t, info.OptionalDefault.Set)
	assert.False(t, info.UndefinedDefault.Set)
}

func TestFlagAliases(t *testing.T) {
	flag := New(TypeInt, IntParser, "test", "-t", "--test")

	got := flag.aliases()

	assert.Size(t, 2, got)
	assert.Equal(t, "-t", got[0])
	assert.Equal(t, "--test", got[1])
}
