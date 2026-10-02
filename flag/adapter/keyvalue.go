package adapter

import "strings"

const (
	KeyValueName        = "KeyValue"
	KeyValueDescription = "Supports 'flag=value' assignment syntax for configured prefixes"
)

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
