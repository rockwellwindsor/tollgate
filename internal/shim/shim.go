package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/rockwellwindsor/tollgate/internal/config"
	"github.com/rockwellwindsor/tollgate/internal/decision"
	"github.com/rockwellwindsor/tollgate/internal/matcher"
	"github.com/rockwellwindsor/tollgate/internal/patterns"
	"github.com/rockwellwindsor/tollgate/internal/prompt"
	"github.com/rockwellwindsor/tollgate/internal/resolver"
	"github.com/rockwellwindsor/tollgate/internal/state"
)

func Run(binaryName string) {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: could not determine shim path: %v\n", err)
		os.Exit(1)
	}

	real, err := resolver.ResolveRealBinary(binaryName, filepath.Dir(self))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	// Fast path: TOLLGATE=off skips everything
	if os.Getenv("TOLLGATE") == "off" {
		execOrDie(real, args)
	}

	result := matcher.Match(binaryName, args, patterns.Defaults)
	if !result.Matched {
		execOrDie(real, args)
	}

	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: config error: %v\n", err)
		os.Exit(1)
	}

	mgr := state.NewManager(state.DefaultDir())
	globalOn, err := mgr.IsGlobalOn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: state error: %v\n", err)
		os.Exit(1)
	}

	ppid := os.Getppid()
	sessionAllow, _ := mgr.IsAllowedForSession(result.Pattern.Name, ppid)
	sessionDeny, _ := mgr.IsSessionPaused(ppid)

	d := decision.Decide(result, decision.State{
		GlobalEnabled: globalOn,
		DefaultAction: cfg.DefaultAction,
		SessionAllow:  sessionAllow,
		SessionDeny:   sessionDeny,
	}, decision.Env{TollgateVar: os.Getenv("TOLLGATE")})

	switch d {
	case decision.AllowSilent:
		execOrDie(real, args)
	case decision.DenySilent:
		fmt.Fprintf(os.Stderr, "tollgate: %s %v blocked\n", binaryName, args)
		os.Exit(1)
	case decision.Prompt:
		resp, err := prompt.Ask(os.Stdout, os.Stdin, binaryName, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tollgate: prompt error: %v\n", err)
			os.Exit(1)
		}
		switch resp {
		case prompt.AllTime:
			_ = mgr.AllowForSession(result.Pattern.Name, ppid)
			execOrDie(real, args)
		case prompt.Yes:
			execOrDie(real, args)
		default:
			fmt.Fprintf(os.Stderr, "tollgate: %s %v denied\n", binaryName, args)
			os.Exit(1)
		}
	}
}

func execOrDie(real string, args []string) {
	if err := syscall.Exec(real, append([]string{real}, args...), os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "tollgate: exec failed: %v\n", err)
		os.Exit(1)
	}
}
