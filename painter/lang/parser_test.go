package lang

import (
	"strings"
	"testing"
)

func TestParseCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected int
	}{
		{
			name:     "valid white command",
			input:    "white",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid green command",
			input:    "green",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid update command",
			input:    "update",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid bgrect command",
			input:    "bgrect 0.1 0.2 0.3 0.4",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid figure command",
			input:    "figure 0.5 0.5",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid move command",
			input:    "move 0.2 0.1",
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "valid reset command",
			input:    "reset",
			wantErr:  false,
			expected: 1,
		},
		{
			name:    "invalid command",
			input:   "paint",
			wantErr: true,
		},
		{
			name:    "missing arguments for bgrect",
			input:   "bgrect 0.1 0.2",
			wantErr: true,
		},
		{
			name:    "invalid float argument",
			input:   "figure one two",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	parser := &Parser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			ops, err := parser.Parse(r)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
			if !tt.wantErr && len(ops) != tt.expected {
				t.Errorf("expected %d operation(s), got %d", tt.expected, len(ops))
			}
		})
	}
}
