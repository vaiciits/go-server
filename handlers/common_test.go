package handlers

import "testing"

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "standard authorization header",
			lines: []string{"Authorization: Bearer abc123"},
			want:  "abc123",
		},
		{
			name:  "case-insensitive header and scheme",
			lines: []string{"aUtHoRiZaTiOn: bEaReR abc123"},
			want:  "abc123",
		},
		{
			name:  "non-bearer authorization",
			lines: []string{"Authorization: Basic abc123"},
		},
		{
			name:  "bearer value without token",
			lines: []string{"Authorization: Bearer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractToken(tt.lines); got != tt.want {
				t.Errorf("ExtractToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
