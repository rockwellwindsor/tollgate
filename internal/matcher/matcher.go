package matcher

type Pattern struct {
	Name         string
	Binary       string
	Subcommand   string
	RequiredArgs []string
	AnyOfArgs    []string
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

	// Specific patterns (with AnyOfArgs) take precedence over general ones.
	for _, p := range patterns {
		if p.Binary != binary || p.Subcommand != subcommand || len(p.AnyOfArgs) == 0 {
			continue
		}
		if !containsAny(args[1:], p.AnyOfArgs) || !containsAll(args[1:], p.RequiredArgs) {
			continue
		}
		return MatchResult{Matched: true, Pattern: p}
	}
	for _, p := range patterns {
		if p.Binary != binary || p.Subcommand != subcommand || len(p.AnyOfArgs) > 0 {
			continue
		}
		if !containsAll(args[1:], p.RequiredArgs) {
			continue
		}
		return MatchResult{Matched: true, Pattern: p}
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

func containsAll(args, required []string) bool {
	for _, r := range required {
		found := false
		for _, a := range args {
			if a == r {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
