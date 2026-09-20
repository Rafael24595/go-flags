package flag

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewCLI(t *testing.T) {
	cli := NewCLI()

	assert.NotNil(t, cli)
}

func TestCLIRegistersFlags(t *testing.T) {
	cli := NewCLI()

	boolFlag := cli.Bool("boolean", "-b")
	stringFlag := cli.String("string", "-s")
	intFlag := cli.Int("integer", "-i")
	int64Flag := cli.Int64("integer", "--int64")
	uintFlag := cli.Uint("unsigned", "-u")
	uint64Flag := cli.Uint64("unsigned", "--uint64")
	floatFlag := cli.Float64("float", "-f")

	assert.NotNil(t, boolFlag)
	assert.NotNil(t, stringFlag)
	assert.NotNil(t, intFlag)
	assert.NotNil(t, int64Flag)
	assert.NotNil(t, uintFlag)
	assert.NotNil(t, uint64Flag)
	assert.NotNil(t, floatFlag)

	assert.Size(t, 7, cli.flags)
}

func TestCLIParseValue(t *testing.T) {
	cli := NewCLI()

	port := cli.Int("port", "-p")

	err := cli.parse([]string{"-p", "8080"})

	assert.Nil(t, err)
	assert.True(t, port.IsPresent())
	assert.True(t, port.isSet())
	assert.Equal(t, 8080, port.Value())
}

func TestCLIParseAlias(t *testing.T) {
	cli := NewCLI()

	port := cli.Int("port", "-p", "--port")

	err := cli.parse([]string{"--port", "8080"})

	assert.Nil(t, err)
	assert.True(t, port.IsPresent())
	assert.Equal(t, 8080, port.Value())
}

func TestCLIParseDefaultUndefined(t *testing.T) {
	cli := NewCLI()

	port := cli.Int("port", "-p").
		DefaultUndefined(8080)

	err := cli.parse(nil)

	assert.Nil(t, err)
	assert.False(t, port.IsPresent())
	assert.True(t, port.isSet())
	assert.Equal(t, 8080, port.Value())
}

func TestCLIParseDefaultOptional(t *testing.T) {
	cli := NewCLI()

	verbose := cli.Bool("verbose", "-v").
		DefaultOptional(true)

	err := cli.parse([]string{"-v"})

	assert.Nil(t, err)
	assert.True(t, verbose.IsPresent())
	assert.True(t, verbose.isSet())
	assert.True(t, verbose.Value())
}

func TestCLIParseOptionalBeforeAnotherFlag(t *testing.T) {
	cli := NewCLI()

	verbose := cli.Bool("verbose", "-v").
		DefaultOptional(true)

	port := cli.Int("port", "-p")

	err := cli.parse([]string{"-v", "-p", "8080"})

	assert.Nil(t, err)
	assert.True(t, verbose.IsPresent())
	assert.True(t, verbose.Value())
	assert.Equal(t, 8080, port.Value())
}

func TestCLIParseExplicitValueOverridesOptionalDefault(t *testing.T) {
	cli := NewCLI()

	verbose := cli.Bool("verbose", "-v").
		DefaultOptional(true)

	err := cli.parse([]string{"-v", "false"})

	assert.Nil(t, err)
	assert.True(t, verbose.IsPresent())
	assert.False(t, verbose.Value())
}

func TestCLIParseRequired(t *testing.T) {
	cli := NewCLI()

	input := cli.String("input", "-i").
		Required()

	err := cli.parse([]string{"-i", "file.go"})

	assert.Nil(t, err)
	assert.True(t, input.IsPresent())
	assert.Equal(t, "file.go", input.Value())
}

func TestCLIParseRequiredMissing(t *testing.T) {
	cli := NewCLI()

	input := cli.String("input", "-i").
		Required()

	err := cli.parse(nil)

	assert.ErrorIs(t, ErrRequiredOption, err)
	assert.False(t, input.IsPresent())
}

func TestCLIParseRequiredWithOptionalDefault(t *testing.T) {
	cli := NewCLI()

	input := cli.String("input", "-i").
		Required().
		DefaultOptional("stdin")

	err := cli.parse([]string{"-i"})

	assert.Nil(t, err)
	assert.True(t, input.IsPresent())
	assert.Equal(t, "stdin", input.Value())
}

func TestCLIParseMissingValue(t *testing.T) {
	cli := NewCLI()

	cli.String("input", "-i")

	err := cli.parse([]string{"-i"})

	assert.ErrorIs(t, ErrMissingValue, err)
}

func TestCLIParseMissingValueBeforeAnotherFlag(t *testing.T) {
	cli := NewCLI()

	cli.String("input", "-i")
	cli.Bool("verbose", "-v").
		DefaultOptional(true)

	err := cli.parse([]string{"-i", "-v"})

	assert.ErrorIs(t, ErrMissingValue, err)
}

func TestCLIParseInvalidValue(t *testing.T) {
	cli := NewCLI()

	value := cli.Int("value", "-v")

	err := cli.parse([]string{"-v", "invalid"})

	assert.ErrorIs(t, ErrInvalidValue, err)
	assert.False(t, value.isSet())
}

func TestCLIParseUnknownOption(t *testing.T) {
	cli := NewCLI()

	err := cli.parse([]string{"--unknown"})

	assert.ErrorIs(t, ErrUnknownOption, err)
}

func TestCLIParseNegativeValue(t *testing.T) {
	cli := NewCLI()

	value := cli.Int("value", "-v")

	err := cli.parse([]string{"-v", "-10"})

	assert.Nil(t, err)
	assert.Equal(t, -10, value.Value())
}

func TestCLIParseMultipleFlags(t *testing.T) {
	cli := NewCLI()

	input := cli.String("input", "-i")
	port := cli.Int("port", "-p")
	verbose := cli.Bool("verbose", "-v").
		DefaultOptional(true)

	err := cli.parse([]string{
		"-i", "input.mp3",
		"-p", "8080",
		"-v",
	})

	assert.Nil(t, err)

	assert.Equal(t, "input.mp3", input.Value())
	assert.Equal(t, 8080, port.Value())
	assert.Equal(t, true, verbose.Value())

	assert.True(t, input.IsPresent())
	assert.True(t, port.IsPresent())
	assert.True(t, verbose.IsPresent())
}
