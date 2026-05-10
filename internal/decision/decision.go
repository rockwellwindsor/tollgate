package decision

import "github.com/rockwellwindsor/tollgate/internal/matcher"

type Decision int

const (
	Prompt      Decision = iota
	AllowSilent Decision = iota
	DenySilent  Decision = iota
)

type State struct {
	GlobalEnabled bool
	DefaultAction string // "prompt", "allow", "deny"
	SessionAllow  bool
	SessionDeny   bool
}

type Env struct {
	TollgateVar string // value of TOLLGATE env var
}

func Decide(match matcher.MatchResult, state State, env Env) Decision {
	return applyPrecedence(state, env)
}

// applyPrecedence evaluates the decision chain in priority order:
// env var > global toggle > session rules > default_action > prompt
func applyPrecedence(state State, env Env) Decision {
	// TOLLGATE=off bypasses everything — single-invocation kill switch
	if env.TollgateVar == "off" {
		return AllowSilent
	}
	// Global off — tollgate disabled across all sessions
	if !state.GlobalEnabled {
		return AllowSilent
	}
	// Session deny takes precedence over session allow
	if state.SessionDeny {
		return DenySilent
	}
	// Session allow — user said "yes for the rest of this session"
	if state.SessionAllow {
		return AllowSilent
	}
	// default_action configured as "deny" — block without prompting
	if state.DefaultAction == "deny" {
		return DenySilent
	}
	// default_action configured as "allow" — soft-off, log but never prompt
	if state.DefaultAction == "allow" {
		return AllowSilent
	}
	return Prompt
}
