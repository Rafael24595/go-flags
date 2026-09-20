package flag

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// CLI parses command-line arguments and stores the registered options.
type CLI struct {
	flags []flag
}

// NewCLI creates a new command-line parser.
func NewCLI() *CLI {
	return &CLI{}
}

// Bool registers a boolean option.
func (c *CLI) Bool(description string, names ...string) *Flag[bool] {
	flag := BoolFlag(description, names...)
	c.add(flag)
	return flag
}

// String registers a string option.
func (c *CLI) String(description string, names ...string) *Flag[string] {
	flag := StringFlag(description, names...)
	c.add(flag)
	return flag
}

// Int registers an integer option.
func (c *CLI) Int(description string, names ...string) *Flag[int] {
	flag := IntFlag(description, names...)
	c.add(flag)
	return flag
}

// Int64 registers a 64-bit signed integer option.
func (c *CLI) Int64(description string, names ...string) *Flag[int64] {
	flag := Int64Flag(description, names...)
	c.add(flag)
	return flag
}

// Uint registers an unsigned integer option.
func (c *CLI) Uint(description string, names ...string) *Flag[uint] {
	flag := UintFlag(description, names...)
	c.add(flag)
	return flag
}

// Uint64 registers a 64-bit unsigned integer option.
func (c *CLI) Uint64(description string, names ...string) *Flag[uint64] {
	flag := Uint64Flag(description, names...)
	c.add(flag)
	return flag
}

// Float64 registers a 64-bit floating-point option.
func (c *CLI) Float64(description string, names ...string) *Flag[float64] {
	flag := Float64Flag(description, names...)
	c.add(flag)
	return flag
}

func (c *CLI) add(f flag) {
	c.flags = append(c.flags, f)
}

// Parse parses the command-line arguments provided to the process.
//
// It returns an error if an unknown option is encountered, a required value
// is missing, an option value cannot be parsed, or a required option was not
// provided.
func (c *CLI) Parse() error {
	return c.parse(os.Args[1:])
}

func (c *CLI) parse(args []string) error {
	lookup := make(map[string]flag)
	for _, f := range c.flags {
		for _, name := range f.aliases() {
			lookup[name] = f
		}
	}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		lenArgs := len(args)

		flag, ok := lookup[arg]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownOption, arg)
		}

		flag.markPresent()

		isNextArg := false
		if lenArgs > 0 {
			_, isNextArg = lookup[args[0]]
		}

		if lenArgs == 0 || isNextArg {
			if !flag.isOptional() {
				return fmt.Errorf("%w: %s", ErrMissingValue, arg)
			}

			flag.applyOptionalDefault()
			continue
		}

		value := args[0]
		args = args[1:]

		if err := flag.parse(value); err != nil {
			return fmt.Errorf("%w for %s: %w", ErrInvalidValue, arg, err)
		}
	}

	for _, flag := range c.flags {
		if flag.isRequired() && !flag.isPresent() {
			aliases := strings.Join(flag.aliases(), ", ")
			return fmt.Errorf("%w: %s", ErrRequiredOption, aliases)
		}

		if !flag.isSet() {
			flag.applyUndefinedDefault()
		}
	}

	return nil
}

func (c *CLI) ShowHelp() error {
	return c.showHelp(os.Stdout, DefaultFormatter)
}

func (c *CLI) showHelp(
	writer io.Writer,
	formater Formatter,
) error {
	help := c.makeHelp(formater)
	_, err := writer.Write([]byte(help))
	return err
}

func (c *CLI) makeHelp(formatter Formatter) string {
	infos := make([]Info, 0, len(c.flags))
	for _, f := range c.flags {
		infos = append(infos, f.info())
	}
	return formatter(infos)
}
