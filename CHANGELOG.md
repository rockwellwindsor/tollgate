# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-05-10

Initial release.

### Added

- `shim-git` and `shim-gh` binaries intercept watched commands and prompt for confirmation
- Pattern matching on `git push`, `git push --force/-f`, `gh pr create`, `gh pr merge`, `gh repo create`, `gh repo delete`, `gh release create`
- Three independent toggle layers: `TOLLGATE=off` env var, session pause, global on/off
- `tollgate status` shows global state, session pause, active session-allow patterns, and audit log entry count
- `tollgate pause` / `tollgate resume` disable and re-enable prompts for the current shell session
- `tollgate on` / `tollgate off` enable and disable tollgate across all sessions
- `tollgate clear-logs` deletes the audit log with `--dry-run` and `--yes` flag support
- `tollgate install` writes shim binaries and default config to `~/.tollgate/bin/`
- Append-only JSONL audit log at `~/.tollgate/audit.log`
- Config file at `~/.tollgate/config.json` with `default_action` field (`prompt`, `allow`, `deny`)
- `TOLLGATE_HOME` environment variable overrides the default state and config directory
- Session-allow (`a` at prompt) skips future prompts for a pattern within the same shell session
- Stale session state pruned automatically when the parent shell process is no longer alive
