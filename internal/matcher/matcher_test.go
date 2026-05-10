package matcher

import "testing"

var defaultPatterns = []Pattern{
	{Name: "git-push", Binary: "git", Subcommand: "push"},
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name        string
		binary      string
		args        []string
		wantMatched bool
		wantPattern string
	}{
		{
			name:        "git push matched",
			binary:      "git",
			args:        []string{"push", "origin", "main"},
			wantMatched: true,
			wantPattern: "git-push",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match(tt.binary, tt.args, defaultPatterns)
			if got.Matched != tt.wantMatched {
				t.Errorf("Matched = %v, want %v", got.Matched, tt.wantMatched)
			}
			if tt.wantMatched && got.Pattern.Name != tt.wantPattern {
				t.Errorf("Pattern.Name = %q, want %q", got.Pattern.Name, tt.wantPattern)
			}
		})
	}
}
