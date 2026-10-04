package adapter

import (
	"testing"
	"testing/fstest"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestEnvMappingTransform(t *testing.T) {
	t.Parallel()

	mockEnv := map[string]string{
		"APP_PORT":     "9000",
		"DATABASE_URL": "postgres://localhost:5432/db",
	}

	mockLookup := func(key string) (string, bool) {
		val, ok := mockEnv[key]
		return val, ok
	}

	mapping := map[string]string{
		"APP_PORT":     "--port",
		"DATABASE_URL": "--db-url",
		"LOG_LEVEL":    "--log-level",
	}

	envAdapter := NewEnvMappingWith(mapping, mockLookup)

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "injects environment variables when flags are missing",
			args:     []string{"--config", "app.json"},
			expected: []string{"--config", "app.json", "--port", "9000", "--db-url", "postgres://localhost:5432/db"},
		},
		{
			name:     "gives precedence to command-line flags over environment variables",
			args:     []string{"--port", "8080"},
			expected: []string{"--port", "8080", "--db-url", "postgres://localhost:5432/db"},
		},
		{
			name:     "does not inject anything when all flags are provided",
			args:     []string{"--port", "8080", "--db-url", "custom-url"},
			expected: []string{"--port", "8080", "--db-url", "custom-url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := envAdapter.Transform(tt.args)

			assert.Size(t, len(tt.expected), result)

			for i := range result {
				assert.Equal(t, tt.expected[i], result[i])
			}
		})
	}
}

func TestLookupEnvWith(t *testing.T) {
	t.Parallel()

	var data string

	data += "# Application Settings\n"
	data += "APP_PORT=9000\n"
	data += "DATABASE_URL=\"postgres://localhost:5432/db\"\n"
	data += "GREETING='hello world'\n\n"

	data += "# Invalid or empty lines\n"
	data += "INVALID_LINE\n"
	data += "EMPTY_VAL=\n"

	mockFS := fstest.MapFS{
		".env": &fstest.MapFile{
			Data: []byte(data),
		},
	}

	mockOSLookup := func(key string) (string, bool) {
		if key == "APP_PORT" {
			return "8080", true
		}
		return "", false
	}

	lookup, err := lookupEnvWith(mockFS, ".env", mockOSLookup)

	assert.Nil(t, err)

	tests := []struct {
		name       string
		key        string
		expected   string
		wantExists bool
	}{
		{
			name:       "gives precedence to OS environment variable over .env file",
			key:        "APP_PORT",
			expected:   "8080",
			wantExists: true,
		},
		{
			name:       "falls back to .env file when OS variable is absent",
			key:        "DATABASE_URL",
			expected:   "postgres://localhost:5432/db",
			wantExists: true,
		},
		{
			name:       "strips single quotes from values in .env file",
			key:        "GREETING",
			expected:   "hello world",
			wantExists: true,
		},
		{
			name:       "parses empty values correctly",
			key:        "EMPTY_VAL",
			expected:   "",
			wantExists: true,
		},
		{
			name:       "returns false for non-existent keys in both OS and .env",
			key:        "MISSING_KEY",
			expected:   "",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotVal, gotExists := lookup(tt.key)

			assert.Equal(t, tt.wantExists, gotExists)
			assert.Equal(t, tt.expected, gotVal)
		})
	}
}

func TestLookupEnvWith_FileNotFound(t *testing.T) {
	t.Parallel()

	emptyFS := fstest.MapFS{}

	mockOSLookup := func(key string) (string, bool) {
		if key == "PORT" {
			return "3000", true
		}
		return "", false
	}

	lookup, err := lookupEnvWith(emptyFS, ".env", mockOSLookup)

	assert.ErrorIs(t, ErrCannotReadDotEnvFile, err)

	val, exists := lookup("PORT")

	assert.True(t, exists)
	assert.Equal(t, "3000", val)

	_, exists = lookup("OTHER_KEY")

	assert.False(t, exists)
}
