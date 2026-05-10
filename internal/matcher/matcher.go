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

	// Check specific patterns (with AnyOfArgs) before general ones.
	for _, specific := range []bool{true, false} {
		for _, p := range patterns {
			if p.Binary != binary || p.Subcommand != subcommand {
				continue
			}
			if specific != (len(p.AnyOfArgs) > 0) {
				continue
			}
			if len(p.AnyOfArgs) > 0 && !containsAny(args[1:], p.AnyOfArgs) {
				continue
			}
			return MatchResult{Matched: true, Pattern: p}
		}
	}
	return MatchResult{}
}

func containsAny(args, targets []string) bool {
	for _, t := range targets {
		for _, a := range args {
			if a == t {
				return true
			}
		}
	}
	return false
}
