package matcher

type Pattern struct {
	Name       string
	Binary     string
	Subcommand string
	AnyOfArgs  []string
}

type MatchResult struct {
	Matched bool
	Pattern Pattern
}

func Match(binary string, args []string, patterns []Pattern) MatchResult {
	if len(args) == 0 {
		return MatchResult{}
	}
	subcommand := args[0]

	for _, p := range patterns {
		if p.Binary == binary && p.Subcommand == subcommand {
			return MatchResult{Matched: true, Pattern: p}
		}
	}
	return MatchResult{}
}
