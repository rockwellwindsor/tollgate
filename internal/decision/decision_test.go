package decision

import (
	"testing"

	"github.com/rockwellwindsor/tollgate/internal/matcher"
)

var matched = matcher.MatchResult{Matched: true}

func TestDecide(t *testing.T) {
	tests := []struct {
		name  string
		state State
		env   Env
		want  Decision
	}{
		{
			name:  "match + global off → AllowSilent",
			state: State{GlobalEnabled: false, DefaultAction: "prompt"},
			env:   Env{},
			want:  AllowSilent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decide(matched, tt.state, tt.env)
			if got != tt.want {
				t.Errorf("Decide() = %v, want %v", got, tt.want)
			}
		})
	}
}
