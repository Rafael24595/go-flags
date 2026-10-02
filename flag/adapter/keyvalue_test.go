package adapter

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewKeyValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		prefixes []string
		input    []string
		want     []string
	}{
		{
			name:  "default prefixes with standard double dash flag",
			input: []string{"--output=/tmp/dir", "--workers=4", "-v"},
			want:  []string{"--output", "/tmp/dir", "--workers", "4", "-v"},
		},
		{
			name:     "default prefixes with standard single dash flag",
			prefixes: nil,
			input:    []string{"-o=/tmp/dir", "-w=4", "-v"},
			want:     []string{"-o", "/tmp/dir", "-w", "4", "-v"},
		},
		{
			name:     "custom prefix (slash for Windows style)",
			prefixes: []string{"/"},
			input:    []string{"/out=C:\\logs", "/verbose"},
			want:     []string{"/out", "C:\\logs", "/verbose"},
		},
		{
			name:     "custom empty prefix (accepts any key=value)",
			prefixes: []string{""},
			input:    []string{"dir=/tmp", "workers=4", "standalone"},
			want:     []string{"dir", "/tmp", "workers", "4", "standalone"},
		},
		{
			name:     "multiple equal signs outside quotes splits only at first equal",
			prefixes: nil,
			input:    []string{"--config=key=value=extra"},
			want:     []string{"--config", "key=value=extra"},
		},
		{
			name:     "equal sign inside double quotes is preserved",
			prefixes: nil,
			input:    []string{"--filter=\"a=b\"", "--query=\"x=1\""},
			want:     []string{"--filter", "\"a=b\"", "--query", "\"x=1\""},
		},
		{
			name:     "equal sign inside single quotes is preserved",
			prefixes: nil,
			input:    []string{"--filter='a=b'"},
			want:     []string{"--filter", "'a=b'"},
		},
		{
			name:     "value with equals outside and inside quotes",
			prefixes: nil,
			input:    []string{"--env=FOO=\"bar=baz\""},
			want:     []string{"--env", "FOO=\"bar=baz\""},
		},
		{
			name:     "arguments without equals are untouched",
			prefixes: nil,
			input:    []string{"--verbose", "-w", "4", "positional"},
			want:     []string{"--verbose", "-w", "4", "positional"},
		},
		{
			name:     "unmatched prefix is not transformed",
			prefixes: []string{"--"},
			input:    []string{"-o=/tmp/dir"},
			want:     []string{"-o=/tmp/dir"},
		},
		{
			name:     "empty slice input",
			prefixes: nil,
			input:    []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kvAdapter := NewKeyValue(tt.prefixes...)
			got := kvAdapter.Transform(tt.input)

			assert.DeepEqual(t, tt.want, got)
		})
	}
}

func TestNewKeyValueMetadata(t *testing.T) {
	t.Parallel()

	kvAdapter := NewKeyValue()

	assert.Equal(t, KeyValueName, kvAdapter.Name)
	assert.Equal(t, KeyValueDescription, kvAdapter.Description)
	assert.NotNil(t, kvAdapter.Transform)
}
