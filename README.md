# tollgate

> A configurable checkpoint for sensitive operations in your dev environment. Stop. Confirm. Proceed.

[![CI](https://github.com/rockwellwindsor/tollgate/actions/workflows/ci.yml/badge.svg)](https://github.com/rockwellwindsor/tollgate/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rockwellwindsor/tollgate)](https://goreportcard.com/report/github.com/rockwellwindsor/tollgate)

`tollgate` is a small CLI tool that intercepts sensitive operations against external services and asks you to confirm them. It works by placing shim binaries named `git` and `gh` on your `$PATH` ahead of the real ones. When a watched command runs (`git push`, `gh pr create`, `gh repo delete`), the shim pauses, prompts you in the terminal, and only proceeds if you say yes.

It exists because agents now write and run code in your environment, and the operations that touch the outside world deserve a moment of human attention. `tollgate` is that moment.

```
$ git push origin main
tollgate: allow git [push origin main]? [y/n/a] y
[ ... real git push proceeds ... ]
```

---

## Why this exists

Most agent-driven workflows lack human-in-the-loop checkpoints for destructive or external operations. Pushing to a remote, creating a pull request, deleting a repository, deploying to production: these are operations where a one-token mistake by an agent has real consequences.

`tollgate` gives you a configurable set of "ask first" gates. It's not a sandbox; it's a tap on the shoulder. The friction is the feature.

The design doesn't try to distinguish "agent" from "user." Both go through the same prompt. If you don't want a prompt for your own work, you have an environment variable kill switch. The model stays honest.

---

## How it works

When installed, `tollgate` puts two shim binaries (`git` and `gh`) into `~/.tollgate/bin/`. You add that directory to the front of your `$PATH`. From then on, any process that runs `git` or `gh` hits the shim first.

The shim's decision sequence:

1. Check the `TOLLGATE` environment variable. If set to `off`, exec the real binary immediately and exit.
2. Match the invocation against a configurable watchlist (`git push`, `gh pr create`, etc.). If unwatched, exec the real binary and exit.
3. Consult current state: global on/off, session pause, prior session-allow rules.
4. If a prompt is needed, write a question to `/dev/tty` and read the response.
5. On allow, exec the real binary. On deny, exit non-zero.
6. Log the decision to `~/.tollgate/audit.log`.

Unwatched invocations (`git status`, `gh pr list`, etc.) pass through with negligible overhead. The shim uses `syscall.Exec` to replace itself with the real binary rather than forking and waiting on it.

---

## Install

### From source

```bash
git clone https://github.com/rockwellwindsor/tollgate
cd tollgate
make build
./bin/tollgate install
```

### Finishing the install

`tollgate install` writes the shim binaries into `~/.tollgate/bin/` and prints the one-line addition for your shell profile:

```bash
export PATH="$HOME/.tollgate/bin:$PATH"
```

After updating your shell config, reload it (`source ~/.zshrc`) or open a new terminal.

Verify the install:

```bash
which git
# /Users/you/.tollgate/bin/git    <- good, this is the shim

tollgate status
# tollgate: ENABLED
# audit: 0 entries
```

---

## Usage

### Commands

```
tollgate status              show current state, active rules, and audit log count
tollgate pause               disable prompts for this shell session
tollgate resume              re-enable prompts for this shell session
tollgate on                  enable globally
tollgate off                 disable globally (across all sessions)
tollgate clear-logs          delete the audit log (with confirmation)
tollgate install             write/update shim binaries and config
```

### The prompt

When a watched command is intercepted:

```
tollgate: allow git [push origin main]? [y/n/a]
```

| Key | Meaning |
|-----|---------|
| `y` | Allow this once |
| `n` | Deny this once |
| `a` | Allow this pattern for the rest of this shell session |

### The kill switch

Bypass `tollgate` for a single invocation:

```bash
TOLLGATE=off git push
```

Or for a whole session:

```bash
tollgate pause
# ...do stuff freely...
tollgate resume
```

Or persistently across all sessions:

```bash
tollgate off    # stays off until you run `tollgate on`
```

The three layers (env var, session pause, global toggle) are independent. Use whichever fits the situation.

---

## Configuration

The config file lives at `~/.tollgate/config.json` (override the directory with `TOLLGATE_HOME`).

### Default watchlist

Out of the box, `tollgate` watches:

- `git push` (any push to a remote)
- `git push --force` / `-f` (separate pattern from plain push)
- `gh pr create`
- `gh pr merge`
- `gh repo create`
- `gh repo delete`
- `gh release create`

Everything else passes through silently.

### `default_action`

The `default_action` field controls what happens when no session rule applies:

| Value | Behavior |
|-------|----------|
| `prompt` | Ask every time (default) |
| `allow` | Log but never prompt, "soft off" |
| `deny` | Block all watched ops without prompting |

```json
{
  "default_action": "prompt"
}
```

---

## Toggling

Three layers of disable, in increasing scope:

| Method | Scope | When to use |
|--------|-------|-------------|
| `TOLLGATE=off cmd` | One invocation | One-off bypass during manual work |
| `tollgate pause` / `tollgate resume` | One shell session | A chunk of manual work in one terminal |
| `tollgate on` / `tollgate off` | All sessions, persistent | Stepping away from agent work entirely |

No time-bounded resume. Time-bounded toggles behave like a cache, and caches in dev tools tend to mislead you about current state. If `tollgate` is on, it's on. If you turned it off, it's off until you turn it back on.

---

## Audit log

Every watched invocation writes a line to `~/.tollgate/audit.log` as JSONL:

```json
{"binary":"git","args":["push","origin","main"],"pattern":"git-push","decision":"allowed-once"}
```

Pass-through (unwatched) invocations are not logged. Only watched ones.

To clear the log:

```bash
tollgate clear-logs              # interactive confirmation
tollgate clear-logs --dry-run    # show what would be deleted, delete nothing
tollgate clear-logs --yes        # skip confirmation
```

`clear-logs` only ever touches the audit log. It will not remove your config, your shims, or any session state.

---

## What I learned building this

- **Errors as values is verbose, but honest.** Every place a thing can fail is visible in the code. After a while, I stopped wishing for `try`/`catch`.
- **Table-driven tests make TDD feel natural.** A fast compile loop and a standard library `testing` package with no setup ceremony made test-first feel like the path of least resistance, not a discipline.
- **`syscall.Exec` is a magic word.** The unwatched-path performance trick (replacing the shim process with the real binary rather than forking and waiting) is the kind of thing you can only do when the language exposes the OS directly. Worth knowing about.
- **Single-binary distribution is the underrated feature.** No runtime, no `node_modules`, no virtualenv. The install story is `make build` and that's it.

A longer writeup of the design decisions and what changed about my thinking is at [windsordevelopmentstudio.io/post/tollgate-confirmation-layer-for-ai-agents](https://windsordevelopmentstudio.io/post/tollgate-confirmation-layer-for-ai-agents).

---

## Roadmap

### v1.x

- Per-directory enable/disable (a `.tollgate-off` file walks up from cwd)
- Optional desktop notifications alongside terminal prompts
- Windows support (PATHEXT handling, signal semantics, cmd.exe quirks)
- Configurable log rotation policy

### v2

- Vercel support with patterns for `vercel deploy --prod`, `vercel env`, `vercel domains`
- Shims for `npm publish`, `cargo publish`, and similar package-publish operations

Open issues track all of the above. Contributions and design suggestions welcome.

---

## Development

Built with Go 1.22+. Tests live alongside the code they test. End-to-end tests live in `tests/e2e/` and spawn real shim binaries against a synthetic `$PATH`.

```bash
make test          # run all tests
make lint          # run golangci-lint
make build         # build all binaries to ./bin/
make test-coverage # generate HTML coverage report
```

---

## License

MIT. See [LICENSE](./LICENSE).

---

*Built by [Rockwell Windsor Rice](https://windsordevelopmentstudio.io/)*
