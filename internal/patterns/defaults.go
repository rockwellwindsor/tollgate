package patterns

import "github.com/rockwellwindsor/tollgate/internal/matcher"

var Defaults = []matcher.Pattern{
	{Name: "git-push", Binary: "git", Subcommand: "push"},
	{Name: "git-push-force", Binary: "git", Subcommand: "push", AnyOfArgs: []string{"--force", "-f"}},
	{Name: "gh-pr-create", Binary: "gh", Subcommand: "pr", RequiredArgs: []string{"create"}},
	{Name: "gh-pr-merge", Binary: "gh", Subcommand: "pr", RequiredArgs: []string{"merge"}},
	{Name: "gh-repo-create", Binary: "gh", Subcommand: "repo", RequiredArgs: []string{"create"}},
	{Name: "gh-repo-delete", Binary: "gh", Subcommand: "repo", RequiredArgs: []string{"delete"}},
	{Name: "gh-release-create", Binary: "gh", Subcommand: "release", RequiredArgs: []string{"create"}},
}
