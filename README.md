# go-flags

A lightweight, zero-dependency type-safe command-line flag parser for Go.

The `flag` package provides a simple API for registering typed command-line options, parsing their values, defining defaults, validating required options, and formatting registered options.

It is designed to keep command-line configuration explicit while remaining easy to extend with custom parsers.

## Features

* Generic, type-safe flags using Go generics.
* Built-in support for:

  * `bool`
  * `string`
  * `int`
  * `int64`
  * `uint`
  * `uint64`
  * `float64`
  * void options without arguments
* Short and long option names.
* Optional values with configurable defaults.
* Defaults for options that are not provided.
* Required options.
* Custom value parsers.
* Customizable help formatting.
* Sentinel errors that can be checked with `errors.Is`.

The standard library's flag package is a good choice for simple command-line programs, but applications with more typed options and configuration rules may benefit from a more explicit API.

go-flags focuses on a few specific ideas:

Typed options. Each registered option has a concrete Flag[T] type, so its value can be retrieved without type assertions.
Explicit defaults. DefaultUndefined and DefaultOptional distinguish between an option that was not provided and one that was provided without a value.
Required options. Options can explicitly declare that they must be present on the command line.
Custom parsers. Any value that can be parsed from a string can be integrated through Parser[T].
Composable formatting. Option metadata is exposed through Info, allowing applications to use the default formatter or provide their own representation.
Small API surface. The package focuses on registering options, parsing arguments, storing typed values, and reporting parsing errors.

The goal is not to provide every possible command-line syntax. Instead, go-flags keeps the parser small and predictable while providing the configuration features commonly needed by typed command-line applications.

## Installation

```bash
go get github.com/Rafael24595/go-flags
```

Import the package:

```go
import "github.com/Rafael24595/go-flags/flag"
```

## Basic usage

Create a CLI, register the options you need, and parse the command line:

```go
cli := flag.NewCLI()

help := cli.Void(
    "Shows this help message",
    "-h", "--help",
)

input := cli.String(
    "Path to the local or remote HTML file containing the track script",
    "-i", "--input",
).Required()

output := cli.String(
    "Destination directory",
    "-o", "--out",
).DefaultUndefined("download")

workers := cli.Uint(
    "Number of parallel downloads",
    "-w", "--workers",
).DefaultUndefined(4)

cover := cli.Int(
    "Cover index",
    "-c", "--cover",
).DefaultUndefined(-1).DefaultOptional(0)

err := cli.Parse()
helpPresent := help.IsPresent()

if err != nil && !helpPresent {
    log.Fatal(err)
}

if helpPresent {
    cli.WriteOptions()
}

fmt.Println(input.Value())
fmt.Println(output.Value())
fmt.Println(workers.Value())
```

For example:

```bash
./app --input tracks.html --workers 8
```

## Options and aliases

An option can have one or more names:

```go
cover := cli.Int(
	"Cover index",
	"-c", "--cover",
)
```

Both names refer to the same option:

```bash
./app -c 2
```

or:

```bash
./app --cover 2
```

The first name is not treated differently from the others; all registered names are aliases for the same flag.

## Built-in types

The CLI provides constructors for the supported built-in types:

```go
cli.Void("Shows this help message", "-h", "--help",)

cli.Bool("Enable verbose output", "-v", "--verbose")

cli.String("Input file", "-i", "--input")

cli.Int("Cover index", "-c", "--cover")

cli.Int64("Maximum size", "--max-size")

cli.Uint("Number of workers", "-w", "--workers")

cli.Uint64("Maximum bytes", "--max-bytes")

cli.Float64("Scale factor", "--scale")
```

Each constructor returns a typed `*Flag[T]`, so its value can be retrieved without type assertions:

```go
workers := cli.Uint("Number of workers", "-w", "--workers")

// workers.Value() has type uint.
value := workers.Value()
```

## Void options

A void option does not consume a command-line argument.

Use `Void` for switches such as `--verbose`, `--quiet`, or `--dry-run`:

```go
verbose := cli.Void(
	"Enable verbose output",
	"-v", "--verbose",
)
```

When the option is present:

```bash
./app --verbose
```

`verbose.IsPresent()` returns `true`.

Its value is the empty string:

```go
verbose.Value() // ""
```

A void option consumes no following argument. For example:

```bash
./app --verbose --input tracks.html
```

`--verbose` does not consume `--input`; the next option is parsed normally.

## Checking whether an option was provided

Use `IsPresent` when you need to distinguish between an option that was explicitly provided and one that was not:

```go
verbose := cli.Bool(
	"Enable verbose output",
	"--verbose",
)

if verbose.IsPresent() {
	fmt.Println("verbose was provided")
}
```

This is particularly useful when an option has a default value.

## Default values

The package supports two different kinds of defaults.

### Undefined default

`DefaultUndefined` defines the value used when the option is **not provided at all**.

```go
output := cli.String(
	"Destination directory",
	"-o", "--output",
).DefaultUndefined("download")
```

Without the option:

```bash
./app
```

the value is:

```go
output.Value() // "download"
```

When the option is explicitly provided, its argument is parsed normally:

```bash
./app --output music
```

```go
output.Value() // "music"
```

### Optional default

`DefaultOptional` defines the value used when the option is provided **without an explicit value**.

```go
cover := cli.Int(
	"Cover index",
	"-c", "--cover",
).DefaultOptional(0)
```

With:

```bash
./app
```

the option is absent.

With:

```bash
./app --cover
```

the value is:

```go
cover.Value() // 0
```

With:

```bash
./app --cover 2
```

the explicit value is parsed:

```go
cover.Value() // 2
```

### Using both defaults

The two defaults can be used together:

```go
cover := cli.Int("Cover index",	"-c", "--cover").
    DefaultUndefined(-1).
	DefaultOptional(0)
```

This gives the option three possible states:

| Command line  | Value |
| ------------- | ----: |
| option absent |  `-1` |
| `--cover`     |   `0` |
| `--cover 2`   |   `2` |

This is useful when absence and explicit presence without a value have different meanings.

## Required options

Use `Required` when an option must be present on the command line:

```go
input := cli.String("Input file", "-i", "--input").Required()
```

Running the application without `--input` or `-i` returns `ErrRequiredOption`.

A required option can still define an optional default:

```go
cover := cli.Int("Cover index",	"-c", "--cover").
	Required().
	DefaultOptional(0)
```

The option must be present, but it can be provided without an explicit value:

```bash
./app --cover
```

results in:

```go
cover.Value() // 0
```

## Custom parsers

The generic `New` function allows flags to use custom parsers.

A parser is a function with this signature:

```go
type Parser[T any] func(string) (T, error)
```

For example, a custom duration parser:

```go
func DurationParser(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}
```

It can then be used to create a typed flag:

```go
timeout := flag.New(
	DurationParser,
	flag.TypeString,
	"Request timeout",
	"--timeout",
)
```

After parsing:

```go
duration := timeout.Value()
```

The returned value has type `time.Duration`.

When using a custom parser, choose an appropriate `TypeName` for the help output:

```go
flag.TypeString
flag.TypeInt
flag.TypeFloat64
```

or another `TypeName` when appropriate.

## TypeName

`TypeName` identifies the type displayed for an option in help output.

The package provides predefined names:

```go
flag.TypeUnknown
flag.TypeVoid
flag.TypeBool
flag.TypeString
flag.TypeInt
flag.TypeInt64
flag.TypeUint
flag.TypeUint64
flag.TypeFloat64
```

`TypeName` is a string type, so it can also represent custom type names:

```go
flag.TypeName("duration")
```

## Parsing errors

The package exposes sentinel errors for common parsing failures:

```go
flag.ErrUnknownOption
flag.ErrUnexpectedValue
flag.ErrMissingValue
flag.ErrRequiredOption
flag.ErrInvalidValue
```

Use `errors.Is` to inspect an error:

```go
if err := cli.Parse(); err != nil {
	switch {
	case errors.Is(err, flag.ErrUnknownOption):
		fmt.Println("unknown option")
	case errors.Is(err, flag.ErrMissingValue):
		fmt.Println("missing option value")
	case errors.Is(err, flag.ErrRequiredOption):
		fmt.Println("required option was not provided")
	case errors.Is(err, flag.ErrInvalidValue):
		fmt.Println("invalid option value")
	default:
		fmt.Println(err)
	}
}
```

Errors preserve their underlying sentinel value when additional context is included, so `errors.Is` can be used reliably.

## Formatting options

The CLI exposes three methods for formatting and writing registered options.

### FormatOptions

`FormatOptions` formats the registered options using a custom `Formatter` and
returns the resulting string:

```go
options := cli.FormatOptions(MyFormatter)
fmt.Println(options)
```

### WriteOptions

`WriteOptions` writes the formatted options to `stdout` using `DefaultFormatter`

```go
if err := cli.WriteOptions(); err != nil {
	log.Fatal(err)
}
```

### WriteOptionsWith

`WriteOptionsWith` allows both the output writer and formatter to be customized:

```go
if err := cli.WriteOptionsWith(os.Stderr, myFormatter); err != nil {
	log.Fatal(err)
}
```

### Custom formatters

The output format is controlled by the `Formatter` type:

```go
type Formatter func([]Info) string
```

The formatter receives all registered options at once, allowing it to inspect
the complete set before producing the output.

The package provides `DefaultFormatter`, which produces a human-readable
tabular representation similar to:

```text

Options:
  -h, --hwlp             Shows this message
  -i, --input    string  Path to the input file (required)
  -o, --output   string  Destination directory [default: download]
  -w, --workers  uint    Number of workers [default: 4]
  -c, --cover    int     Cover index [default: -1] [optional: 0]

```

Required options are marked with:

```text
(required)
```

Undefined defaults are displayed as:

```text
[default: value]
```

Optional defaults are displayed as:

```text
[optional: value]
```

Applications can provide their own formatter when a different representation
is required.

## Option metadata

The `Info` type contains the metadata used to describe an option:

```go
type Info struct {
	Names       []string
	Type        TypeName
	Description string

	Void     bool
	Required bool

	OptionalDefault  DefaultInfo
	UndefinedDefault DefaultInfo
}
```

`DefaultInfo` describes a configured default:

```go
type DefaultInfo struct {
	Value any
	Set   bool
}
```

The `Set` field distinguishes between a default that was not configured and a default whose value happens to be the zero value of its type.

For example:

```go
DefaultInfo{
	Value: 0,
	Set:   true,
}
```

means that `0` was explicitly configured as a default.

## Flag configuration

Flag configuration methods return the same flag, allowing method chaining:

```go
cover := cli.Int("Cover index",	"-c", "--cover").
	DefaultUndefined(-1).
	DefaultOptional(0).
	Required()
```

The main configuration methods are:

```go
Void()
DefaultUndefined(value)
DefaultOptional(value)
Required()
```

`Void` makes the option argument-less and clears optional and undefined defaults. It also clears the required state.

Configuring a default or marking a flag as required changes it from void mode.

## Design

The package separates command-line parsing from option representation and help formatting.

A `Flag[T]` contains the configuration and parsed value for a single option.

The internal flag abstraction allows the CLI to work with flags of different generic types while preserving their concrete types for callers.

`Info` is a public metadata representation used by formatters. This keeps help formatting independent from the internal flag implementation.

Parsers are generic functions:

```go
type Parser[T any] func(string) (T, error)
```

This makes adding a custom value type possible without modifying the parser itself.

## Current command-line syntax

The parser currently supports options in the following form:

```text
--name value
-n value
--flag
```

Short and long names are registered explicitly.

A value is considered an argument when the following token is not the name of another registered option. This allows negative numeric values to be used normally:

```bash
./app --cover -1
```

Positional arguments and more advanced syntaxes such as:

```text
--name=value
-nvalue
-nv
--
```

are not part of the current parser syntax for now.

## Example application

A complete small example:

```go
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/Rafael24595/go-flags/flag"
)

func main() {
	cli := flag.NewCLI()

	input := cli.String(
		"Path to the input file",
		"-i", "--input",
	).Required()

	output := cli.String(
		"Destination directory",
		"-o", "--output",
	).DefaultUndefined("download")

	workers := cli.Uint(
		"Number of parallel workers",
		"-w", "--workers",
	).DefaultUndefined(4)

	cover := cli.Int(
		"Cover index",
		"-c", "--cover",
	).
		DefaultUndefined(-1).
		DefaultOptional(0)

	verbose := cli.Void(
		"Enable verbose output",
		"-v", "--verbose",
	)

	if err := cli.Parse(); err != nil {
		switch {
		case errors.Is(err, flag.ErrRequiredOption):
			log.Fatal("required option was not provided")
		case errors.Is(err, flag.ErrUnknownOption):
			log.Fatal("unknown option")
		default:
			log.Fatal(err)
		}
	}

	fmt.Println("input:", input.Value())
	fmt.Println("output:", output.Value())
	fmt.Println("workers:", workers.Value())
	fmt.Println("cover:", cover.Value())
	fmt.Println("verbose:", verbose.IsPresent())
}
```

For example:

```bash
./app \
	--input tracks.html \
	--output downloads \
	--workers 8 \
	--cover \
	--verbose
```

produces values equivalent to:

```text
input: tracks.html
output: downloads
workers: 8
cover: 0
verbose: true
```

## License

See the repository license for details.
