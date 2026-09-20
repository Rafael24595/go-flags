package flag

import (
	"fmt"
	"strings"
)

// Formatter formats command-line option metadata into a help message.
//
// The formatter receives all registered options at once, allowing it to
// calculate column widths or otherwise use information from the complete
// option set before producing the output.
type Formatter func([]Info) string

// DefaultFormatter formats command-line options as a human-readable table.
//
// Required options are marked with "(required)". Optional defaults are shown
// with "[optional: value]", while defaults used when an option is not provided
// are shown with "[default: value]".
//
// The formatter aligns option names and value types into columns based on the
// widest value in the complete option set.
func DefaultFormatter(flags []Info) string {
	if len(flags) == 0 {
		return ""
	}

	nameWidth, typeWidth := maxWidths(flags)

	var sb strings.Builder

	sb.WriteString("\nOptions:\n\n")

	for _, flag := range flags {
		names := strings.Join(flag.Names, ", ")

		sb.WriteString("  ")
		fmt.Fprintf(&sb, "%-*s", nameWidth, names)
		sb.WriteString("  ")

		if flag.Type != "" {
			fmt.Fprintf(&sb, "%-*s", typeWidth, flag.Type)
			sb.WriteString("  ")
		}

		sb.WriteString(flag.Description)

		if flag.Required {
			sb.WriteString(" (required)")
		}

		if flag.UndefinedDefault.Set {
			fmt.Fprintf(&sb, " [default: %v]", flag.UndefinedDefault.Value)
		}

		if flag.OptionalDefault.Set {
			fmt.Fprintf(&sb, " [optional: %v]", flag.OptionalDefault.Value)
		}

		sb.WriteByte('\n')
	}

	sb.WriteByte('\n')

	return sb.String()
}

func maxWidths(flags []Info) (int, int) {
	nameWidth := 0
	typeWidth := 0

	for _, flag := range flags {
		names := strings.Join(flag.Names, ", ")

		if len(names) > nameWidth {
			nameWidth = len(names)
		}

		if len(flag.Type) > typeWidth {
			typeWidth = len(flag.Type)
		}
	}

	return nameWidth, typeWidth
}
