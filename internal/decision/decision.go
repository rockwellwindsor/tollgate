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
	if !state.GlobalEnabled {
		return AllowSilent
	}
	if state.SessionAllow {
		return AllowSilent
	}
	return Prompt
}
