package adapter

import "strings"

const (
	// KeyValueName is the identifier for the key-value argument adapter.
	KeyValueName        = "key-value"
	// KeyValueDescription briefly explains the purpose of the key-value adapter.
	KeyValueDescription = "Supports 'flag=value' assignment syntax for configured prefixes"
)

// NewKeyValue creates an Adapter that expands key-value style arguments (e.g., "--flag=value")
// into distinct flag and value tokens (e.g., "--flag", "value").
//
// Optional prefixes restrict which arguments are split. If no prefixes are provided,
// defaults to matching any argument containing an unquoted '=' sign.
func NewKeyValue(prefixes ...string) Adapter {
	if len(prefixes) == 0 {
		prefixes = []string{""}
	}

	return Adapter{
		Name:        KeyValueName,
		Description: KeyValueDescription,
		Transform:   keyValueTransform(prefixes...),
	}
}

func keyValueTransform(prefixes ...string) Transform {
	return func(args []string) []string {
		result := make([]string, 0, len(args))

		for _, arg := range args {
			if !hasPrefix(prefixes, arg) || !strings.Contains(arg, "=") {
				result = append(result, arg)
				continue
			}

			key, value, found := splitFirstUnquotedEqual(arg)
			if found {
				result = append(result, key, value)
			} else {
				result = append(result, arg)
			}
		}

		return result
	}
}

func hasPrefix(prefixes []string, arg string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}

func splitFirstUnquotedEqual(s string) (key, value string, found bool) {
	inQuotes := false
	var quoteChar rune

	for i, r := range s {
		switch r {
		case '\'', '"':
			if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else if quoteChar == r {
				inQuotes = false
			}
		case '=':
			if !inQuotes {
				return s[:i], s[i+1:], true
			}
		}
	}

	return "", "", false
}
