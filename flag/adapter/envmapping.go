package adapter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strings"
)

// ErrCannotReadDotEnvFile indicates an error occurred while attempting to read or open a .env file.
var ErrCannotReadDotEnvFile = errors.New("cannot read .env file")

const (
	// EnvMappingName is the identifier for the environment variable mapping adapter.
	EnvMappingName        = "env-mapping"
	// EnvMappingDescription briefly explains the purpose of the environment mapping adapter.
	EnvMappingDescription = "Environment variables mapping"
)

// LookupEnv defines a function signature for retrieving environment variable values by key.
// It returns the value and a boolean indicating whether the key was found.
type LookupEnv func(key string) (string, bool)

// NewEnvMapping creates an Adapter that maps environment variables to CLI flags using os.LookupEnv.
//
// Arguments already provided on the command line take precedence over environment variables.
func NewEnvMapping(mapping map[string]string) Adapter {
	return NewEnvMappingWith(mapping, os.LookupEnv)
}

// NewEnvMappingWith creates an Adapter that maps environment variables to CLI flags
// using a custom LookupEnv function.
func NewEnvMappingWith(
	mapping map[string]string,
	lookup LookupEnv,
) Adapter {
	return Adapter{
		Name:        EnvMappingName,
		Description: EnvMappingDescription,
		Transform:   envMappingTransform(mapping, lookup),
	}
}

func envMappingTransform(
	mapping map[string]string,
	lookup LookupEnv,
) Transform {
	return func(args []string) []string {
		provided := make(map[string]bool)
		for _, arg := range args {
			provided[arg] = true
		}

		result := append([]string{}, args...)
		
		keys := envMapKeys(mapping)
		for _, envVar := range keys {
			flagName := mapping[envVar]

			if provided[flagName] {
				continue
			}

			val, exists := lookup(envVar)
			if !exists {
				continue
			}

			result = append(result, flagName, val)
		}

		return result
	}
}

func envMapKeys(mapping map[string]string) []string {
	keys := make([]string, 0, len(mapping))
	for envVar := range mapping {
		keys = append(keys, envVar)
	}

	slices.Sort(keys)
	return keys
}

// LookupEnvWithDot loads environment key-value pairs from a .env file located at path
// within the current working directory, returning a LookupEnv function.
//
// The returned LookupEnv checks os.LookupEnv first, falling back to the parsed .env values.
// If the file cannot be read, it returns a wrapped ErrCannotReadDotEnvFile along with a valid fallback lookup function.
func LookupEnvWithDot(path string) (LookupEnv, error) {
	return lookupEnvWith(os.DirFS("."), path, os.LookupEnv)
}

func lookupEnvWith(
	sys fs.FS,
	path string,
	lookup LookupEnv,
) (LookupEnv, error) {
	dotEnv, err := readDotEnvWith(sys, path)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrCannotReadDotEnvFile, err)
	}

	return func(key string) (string, bool) {
		if val, exists := lookup(key); exists {
			return val, true
		}

		val, exists := dotEnv[key]
		return val, exists
	}, err
}

func readDotEnvWith(sys fs.FS, path string) (map[string]string, error) {
	file, err := sys.Open(path)
	if err != nil {
		return make(map[string]string), err
	}

	defer func() {
		err = errors.Join(err, file.Close())
	}()

	return parseDotEnvReader(file)
}

func parseDotEnvReader(reader io.Reader) (envs map[string]string, err error) {
	envs = make(map[string]string)

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if key, value, ok := parseEnvLine(line); ok {
			envs[key] = value
		}
	}

	return envs, scanner.Err()
}

func parseEnvLine(line string) (string, string, bool) {
	key, val, found := strings.Cut(line, "=")
	if !found {
		return "", "", false
	}

	key = strings.TrimSpace(key)
	val = strings.TrimSpace(val)

	if isQuoteWrapped(val) {
		val = val[1 : len(val)-1]
	}

	return key, val, key != ""
}

func isQuoteWrapped(value string) bool {
	if len(value) < 2 {
		return false
	}

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return true
	}

	return strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")
}
