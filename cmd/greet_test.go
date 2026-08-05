package cmd

import (
	"bytes"
	"testing"
)

func TestGreetCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default name",
			args:     []string{},
			expected: "Hello, World!\n",
		},
		{
			name:     "custom name",
			args:     []string{"--name", "Alice"},
			expected: "Hello, Alice!\n",
		},
		{
			name:     "short flag",
			args:     []string{"-n", "Bob"},
			expected: "Hello, Bob!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(append([]string{"greet"}, tt.args...))

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			got := buf.String()
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
