package flag

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestDefaultFormatter(t *testing.T) {
	flags := []Info{
		{
			Names:       []string{"-b", "--bool"},
			Type:        TypeBool,
			Description: "Enable verbose output",
		},
		{
			Names:       []string{"-s", "--string"},
			Type:        TypeString,
			Description: "Input file",
			Required:    true,
			UndefinedDefault: DefaultInfo{
				Value: "input.txt",
				Set:   true,
			},
		},
		{
			Names:       []string{"-i", "--int"},
			Type:        TypeInt,
			Description: "Number of items",
			UndefinedDefault: DefaultInfo{
				Value: 10,
				Set:   true,
			},
		},
		{
			Names:       []string{"--int64"},
			Type:        TypeInt64,
			Description: "64-bit integer",
			UndefinedDefault: DefaultInfo{
				Value: int64(100),
				Set:   true,
			},
		},
		{
			Names:       []string{"--uint"},
			Type:        TypeUint,
			Description: "Unsigned integer",
			UndefinedDefault: DefaultInfo{
				Value: uint(20),
				Set:   true,
			},
		},
		{
			Names:       []string{"--uint64"},
			Type:        TypeUint64,
			Description: "64-bit unsigned integer",
			UndefinedDefault: DefaultInfo{
				Value: uint64(30),
				Set:   true,
			},
		},
		{
			Names:       []string{"--float"},
			Type:        TypeFloat64,
			Description: "Floating-point value",
			UndefinedDefault: DefaultInfo{
				Value: 1.5,
				Set:   true,
			},
		},
	}

	got := DefaultFormatter(flags)

	want := "\nOptions:\n\n"
	want += "  -b, --bool    bool     Enable verbose output\n"
	want += "  -s, --string  string   Input file (required) [default: input.txt]\n"
	want += "  -i, --int     int      Number of items [default: 10]\n"
	want += "  --int64       int64    64-bit integer [default: 100]\n"
	want += "  --uint        uint     Unsigned integer [default: 20]\n"
	want += "  --uint64      uint64   64-bit unsigned integer [default: 30]\n"
	want += "  --float       float64  Floating-point value [default: 1.5]\n\n"

	assert.Equal(t, want, got)
}

func TestDefaultFormatterRequired(t *testing.T) {
	flags := []Info{
		{
			Names:       []string{"-i", "--input"},
			Type:        TypeString,
			Description: "Input file",
			Required:    true,
		},
	}

	got := DefaultFormatter(flags)

	want := "\nOptions:\n\n"
	want += "  -i, --input  string  Input file (required)\n\n"

	assert.Equal(t, want, got)
}

func TestDefaultFormatterOptionalAndUndefinedDefaults(t *testing.T) {
	flags := []Info{
		{
			Names:       []string{"-c", "--cover"},
			Type:        TypeInt,
			Description: "Cover index",

			OptionalDefault: DefaultInfo{
				Value: 0,
				Set:   true,
			},

			UndefinedDefault: DefaultInfo{
				Value: -1,
				Set:   true,
			},
		},
	}

	got := DefaultFormatter(flags)

	want := "\nOptions:\n\n"
	want += "  -c, --cover  int  Cover index [default: -1] [optional: 0]\n\n"

	assert.Equal(t, want, got)
}

func TestDefaultFormatterEmpty(t *testing.T) {
	got := DefaultFormatter(nil)

	assert.Empty(t, got)
}

func TestDefaultFormatterWithoutType(t *testing.T) {
	flags := []Info{
		{
			Names:       []string{"-v", "--verbose"},
			Description: "Enable verbose output",
		},
	}

	got := DefaultFormatter(flags)

	want := "\nOptions:\n\n"
	want += "  -v, --verbose  Enable verbose output\n\n"

	assert.Equal(t, want, got)
}
