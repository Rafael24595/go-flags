package flag

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Rafael24595/go-flags/flag/adapter"
)

// CLI parses command-line arguments and stores the registered options.
type CLI struct {
	flags    []flag
	adapters []adapter.Adapter
}

// NewCLI creates a new command-line parser.
func NewCLI() *CLI {
	return &CLI{}
}

// Use registers one or more command-line adapters.
func (c *CLI) Use(adapters ...adapter.Adapter) {
	c.adapters = append(c.adapters, adapters...)
}

// Add registers a command-line flag.
func (c *CLI) Add(f flag) {
	c.flags = append(c.flags, f)
}

// Void registers a void option.
func (c *CLI) Void(description string, names ...string) *Flag[string] {
	flag := VoidFlag(description, names...)
	c.Add(flag)
	return flag
}

// Bool registers a boolean option.
func (c *CLI) Bool(description string, names ...string) *Flag[bool] {
	flag := BoolFlag(description, names...)
	c.Add(flag)
	return flag
}

// String registers a string option.
func (c *CLI) String(description string, names ...string) *Flag[string] {
	flag := StringFlag(description, names...)
	c.Add(flag)
	return flag
}

// Int registers an integer option.
func (c *CLI) Int(description string, names ...string) *Flag[int] {
	flag := IntFlag(description, names...)
	c.Add(flag)
	return flag
}

// Int64 registers a 64-bit signed integer option.
func (c *CLI) Int64(description string, names ...string) *Flag[int64] {
	flag := Int64Flag(description, names...)
	c.Add(flag)
	return flag
}

// Uint registers an unsigned integer option.
func (c *CLI) Uint(description string, names ...string) *Flag[uint] {
	flag := UintFlag(description, names...)
	c.Add(flag)
	return flag
}

// Uint64 registers a 64-bit unsigned integer option.
func (c *CLI) Uint64(description string, names ...string) *Flag[uint64] {
	flag := Uint64Flag(description, names...)
	c.Add(flag)
	return flag
}

// Float64 registers a 64-bit floating-point option.
func (c *CLI) Float64(description string, names ...string) *Flag[float64] {
	flag := Float64Flag(description, names...)
	c.Add(flag)
	return flag
}

// SliceBool registers a boolean slice option.
func (c *CLI) SliceBool(description string, names ...string) *SliceFlag[bool] {
	flag := SliceBoolFlag(description, names...)
	c.Add(flag)
	return flag
}

// SliceString registers a string slice option.
func (c *CLI) SliceString(description string, names ...string) *SliceFlag[string] {
	flag := SliceStringFlag(description, names...)
	c.Add(flag)
	return flag
}

// SliceInt registers an integer slice option.
func (c *CLI) SliceInt(description string, names ...string) *SliceFlag[int] {
	flag := SliceIntFlag(description, names...)
	c.Add(flag)
	return flag
}

// SliceInt64 registers a 64-bit signed integer slice option.
func (c *CLI) SliceInt64(description string, names ...string) *SliceFlag[int64] {
	flag := SliceInt64Flag(description, names...)
	c.Add(flag)
	return flag
}

// SliceUint registers an unsigned integer slice option.
func (c *CLI) SliceUint(description string, names ...string) *SliceFlag[uint] {
	flag := SliceUintFlag(description, names...)
	c.Add(flag)
	return flag
}

// SliceUint64 registers a 64-bit unsigned integer slice option.
func (c *CLI) SliceUint64(description string, names ...string) *SliceFlag[uint64] {
	flag := SliceUint64Flag(description, names...)
	c.Add(flag)
	return flag
}

// SliceFloat64 registers a 64-bit floating-point slice option.
func (c *CLI) SliceFloat64(description string, names ...string) *SliceFlag[float64] {
	flag := SliceFloat64Flag(description, names...)
	c.Add(flag)
	return flag
}

// Parse parses the command-line arguments provided to the process.
//
// It returns an error if an unknown option is encountered, a required value
// is missing, an option value cannot be parsed, or a required option was not
// provided.
func (c *CLI) Parse() error {
	return c.ParseWith(os.Args[1:])
}

// ParseWith parses the given command-line arguments.
//
// It returns an error if an unknown option is encountered, a required value
// is missing, an option value cannot be parsed, or a required option was not
// provided.
func (c *CLI) ParseWith(args []string) error {
	for _, adapter := range c.adapters {
		args = adapter.Transform(args)
	}

	lookup := make(map[string]flag)
	for _, f := range c.flags {
		f.reset()
		for _, name := range f.aliases() {
			lookup[name] = f
		}
	}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		flag, ok := lookup[arg]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownOption, arg)
		}

		var value *string
		if len(args) > 0 {
			nextArg := args[0]

			_, isNextFlag := lookup[nextArg]
			if !isNextFlag {
				value = &nextArg
			}
		}

		consumed, err := flag.process(value)
		if err != nil {
			return fmt.Errorf("%w: %s", err, arg)
		}

		if consumed {
			args = args[1:]
		}
	}

	for _, flag := range c.flags {
		if err := flag.validate(); err != nil {
			aliases := strings.Join(flag.aliases(), ", ")
			return fmt.Errorf("%w: %s", err, aliases)
		}
	}

	return nil
}

// WriteOptions writes the registered options and adapters to stdout using the default formatter.
func (c *CLI) WriteOptions() error {
	return c.WriteOptionsWith(os.Stdout, DefaultFormatter)
}

// WriteOptionsWith writes the registered options and adapters to the given writer using the
// given formatter.
func (c *CLI) WriteOptionsWith(
	writer io.Writer,
	formater Formatter,
) error {
	help := c.FormatOptions(formater)
	_, err := writer.Write([]byte(help))
	return err
}

// FormatOptions formats the registered options and adapters using the given formatter.
func (c *CLI) FormatOptions(formatter Formatter) string {
	flags := make([]Info, 0, len(c.flags))
	for _, f := range c.flags {
		flags = append(flags, f.info())
	}

	adapters := make([]string, 0, len(c.adapters))
	for _, a := range c.adapters {
		adapters = append(adapters, a.Description)
	}

	return formatter(CLIInfo{
		Flags:    flags,
		Adapters: adapters,
	})
}
