package matcher_test

import (
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/matcher"
	"github.com/rockwellwindsor/tollgate/internal/patterns"
)

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
		{
			name:        "git status not matched",
			binary:      "git",
			args:        []string{"status"},
			wantMatched: false,
		},
		{
			name:        "git push --force matched as separate pattern",
			binary:      "git",
			args:        []string{"push", "--force", "origin", "main"},
			wantMatched: true,
			wantPattern: "git-push-force",
		},
		{
			name:        "gh pr create matched",
			binary:      "gh",
			args:        []string{"pr", "create", "--title", "my pr"},
			wantMatched: true,
			wantPattern: "gh-pr-create",
		},
		{
			name:        "gh pr list not matched",
			binary:      "gh",
			args:        []string{"pr", "list"},
			wantMatched: false,
		},
		{
			name:        "args in different order still matched",
			binary:      "git",
			args:        []string{"push", "origin", "--force", "main"},
			wantMatched: true,
			wantPattern: "git-push-force",
		},
		{
			name:        "gh release create matched",
			binary:      "gh",
			args:        []string{"release", "create", "v1.0.0"},
			wantMatched: true,
			wantPattern: "gh-release-create",
		},
		{
			name:        "gh repo delete matched",
			binary:      "gh",
			args:        []string{"repo", "delete", "my-repo"},
			wantMatched: true,
			wantPattern: "gh-repo-delete",
		},
		{
			name:        "gh repo create matched",
			binary:      "gh",
			args:        []string{"repo", "create", "my-repo"},
			wantMatched: true,
			wantPattern: "gh-repo-create",
		},
		{
			name:        "gh pr merge matched",
			binary:      "gh",
			args:        []string{"pr", "merge", "42"},
			wantMatched: true,
			wantPattern: "gh-pr-merge",
		},
		{
			name:        "short flag -f matches git-push-force",
			binary:      "git",
			args:        []string{"push", "-f", "origin", "main"},
			wantMatched: true,
			wantPattern: "git-push-force",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.Match(tt.binary, tt.args, patterns.Defaults)
			if got.Matched != tt.wantMatched {
				t.Errorf("Matched = %v, want %v", got.Matched, tt.wantMatched)
			}
			if tt.wantMatched && got.Pattern.Name != tt.wantPattern {
				t.Errorf("Pattern.Name = %q, want %q", got.Pattern.Name, tt.wantPattern)
			}
		})
	}
}
