package flag

import (
	"errors"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewWithoutNames(t *testing.T) {
	assert.Panic(t, func() {
		NewFlag(TypeInt, IntParser, "test")
	})
}

func TestFlagCustomParser(t *testing.T) {
	expected := errors.New("parse error")

	parser := func(string) (int, error) {
		return 0, expected
	}

	flag := NewFlag(TypeUnknown, parser, "test", "-t")

	value := "value"
	consumed, err := flag.process(&value)

	assert.False(t, consumed)
	assert.False(t, flag.set)
	assert.ErrorIs(t, expected, err)
}

func TestFlagVoid(t *testing.T) {
	flag := NewFlag(TypeVoid, VoidParser, "test", "-v").
		Void()

	assert.True(t, flag.void)
	assert.False(t, flag.hasOptionalDefault)
	assert.False(t, flag.required)
	assert.False(t, flag.set)
	assert.False(t, flag.present)
}

func TestFlagVoidResetsConfiguration(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "int", "-c").
		DefaultOptional(0).
		DefaultUndefined(-1).
		Required().
		Void()

	assert.True(t, flag.void)
	assert.False(t, flag.hasOptionalDefault)
	assert.False(t, flag.required)
	assert.False(t, flag.hasOptionalDefault)
	assert.False(t, flag.hasUndefinedDefault)
}

func TestFlagDefaultUndefined(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t").
		DefaultUndefined(-1)

	assert.True(t, flag.hasUndefinedDefault)

	err := flag.validate()

	assert.Nil(t, err)
	assert.True(t, flag.set)
	assert.Equal(t, -1, flag.Value())
}

func TestFlagDefaultUndefinedDisablesVoid(t *testing.T) {
	flag := NewFlag(TypeVoid, VoidParser, "void", "-v").
		Void().
		DefaultUndefined("")

	assert.False(t, flag.void)
}

func TestFlagDefaultOptional(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t").
		DefaultOptional(10)

	assert.True(t, flag.hasOptionalDefault)

	consumed, err := flag.process(nil)

	assert.False(t, consumed)
	assert.Nil(t, err)
	assert.True(t, flag.set)
	assert.Equal(t, 10, flag.Value())
}

func TestFlagDefaultOptionalDisablesVoid(t *testing.T) {
	flag := NewFlag(TypeVoid, VoidParser, "void", "-v").
		Void().
		DefaultOptional("")

	assert.False(t, flag.void)
	assert.True(t, flag.hasOptionalDefault)
}

func TestFlagRequired(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t").
		Required()

	assert.True(t, flag.required)
}

func TestFlagRequiredDisablesVoid(t *testing.T) {
	flag := NewFlag(TypeVoid, VoidParser, "void", "-v").
		Void().
		Required()

	assert.False(t, flag.void)
	assert.True(t, flag.required)
}

func TestFlagParse(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t")

	value := "42"
	consumed, err := flag.process(&value)

	assert.True(t, consumed)
	assert.Nil(t, err)

	assert.True(t, flag.set)
	assert.Equal(t, 42, flag.Value())
}

func TestFlagParseError(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t")

	value := "invalid"
	consumed, err := flag.process(&value)

	assert.False(t, consumed)
	assert.NotNil(t, err)
	assert.False(t, flag.set)
}

func TestFlagInfo(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "Cover index", "-c", "--cover").
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
	flag := NewFlag(TypeInt, IntParser, "test", "-t", "--test")

	got := flag.aliases()

	assert.Size(t, 2, got)
	assert.Equal(t, "-t", got[0])
	assert.Equal(t, "--test", got[1])
}

func TestFlagDuplicateOption(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t")

	val1 := "10"
	consumed, err := flag.process(&val1)
	assert.True(t, consumed)
	assert.Nil(t, err)

	val2 := "20"
	consumed, err = flag.process(&val2)
	assert.False(t, consumed)
	assert.ErrorIs(t, ErrDuplicateOption, err)
}

func TestFlagMissingValue(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t")

	consumed, err := flag.process(nil)
	assert.False(t, consumed)
	assert.ErrorIs(t, ErrMissingValue, err)
}

func TestFlagVoidUnexpectedValue(t *testing.T) {
	flag := VoidFlag("verbose", "-v")

	val := "unexpected"
	consumed, err := flag.process(&val)
	assert.False(t, consumed)
	assert.ErrorIs(t, ErrUnexpectedValue, err)
}

func TestFlagRequiredMissing(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t").Required()

	err := flag.validate()
	assert.ErrorIs(t, ErrRequiredOption, err)

	val := "1"
	_, _ = flag.process(&val)
	assert.Nil(t, flag.validate())
}

func TestFlagReset(t *testing.T) {
	flag := NewFlag(TypeInt, IntParser, "test", "-t")

	val := "42"
	_, _ = flag.process(&val)
	assert.Equal(t, 42, flag.Value())

	flag.reset()
	assert.False(t, flag.set)
	assert.False(t, flag.present)
	assert.Equal(t, 0, flag.Value())
}
