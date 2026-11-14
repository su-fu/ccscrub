# ccscrub

[English](README.md) | [日本語](README.ja.md) | [中文](README.zh.md)

> **Unofficial.** Not affiliated with or endorsed by Anthropic.

`ccscrub` scans Claude Code's session transcripts (`~/.claude/projects/**/*.jsonl`) for secrets such as API keys and tokens, and masks them in place. Every line stays valid JSON, so you can still resume the session.

> **This is a learning project.** I built it to learn Go (streaming I/O, `regexp`, safe file rewriting). For real use, consider more complete tools such as [agentleaks](https://github.com/Arthur031221/agentleaks) or [claude-code-scrubify](https://github.com/benjamin-awd/claude-code-scrubify). They cover more AI tools and detect more kinds of secrets.

## Why this exists

Hooks that redact secrets (`PreToolUse` / `PostToolUse`) only control what is sent to the model. The transcript that Claude Code writes to disk can still contain the raw value. `ccscrub` deals with that "after it was written" part.

## Install

Requires Go 1.27.4 or later.

```sh
go install github.com/su-fu/ccscrub/cmd/ccscrub@latest
```

## Usage

The recommended flow is: **scan → preview → fix**.

```sh
ccscrub scan             # find secrets (read-only)
ccscrub fix --dry-run    # show what would be replaced (read-only)
ccscrub fix              # replace them in place
```

### Example

```sh
$ ccscrub scan
C:\Users\you\.claude\projects\my-project\0a1b2c3d.jsonl:42:118: github-token ghp_****AAAA
secrets found: 1, files: 1
```

Each line is `file:line:column: rule-id preview`. The preview shows only the first and last 4 characters, so `ccscrub` never prints a full secret. The summary line goes to stderr.

After `fix`, the secret in the transcript becomes:

```json
{"text":"[REDACTED:github-token]"}
```

### Flags

| Flag | Commands | Description |
|---|---|---|
| `--dry-run` | `fix` | Show what would change without modifying any files |
| `--min-age` | `scan`, `fix` | Skip transcripts modified within this duration, as they may be in use (default `30m`) |
| `--lang` | `scan`, `fix` | Message language: `en`, `ja`, `zh` (default `en`). You can also set `CCSCRUB_LANG` |

### Exit codes

| Code | Meaning |
|---|---|
| `0` | No secrets found (`scan`, `fix --dry-run`), or `fix` finished |
| `1` | Secrets found (`scan`, `fix --dry-run`) |
| `2` | Error (bad flag, `~/.claude/projects` not found, etc.) |

## What it detects

| Rule ID | Example prefix |
|---|---|
| `github-token` | `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`, `github_pat_` |
| `aws-access-key-id` | `AKIA`, `ASIA` |
| `anthropic-api-key` | `sk-ant-api03-` |
| `anthropic-admin-api-key` | `sk-ant-admin01-` |
| `openai-api-key` | `sk-proj-`, `sk-svcacct-`, `sk-admin-`, older `sk-` keys |
| `slack-token` | `xoxb-`, `xoxp-`, `xoxe-` |
| `google-api-key` | `AIza` |
| `private-key` | `-----BEGIN ... PRIVATE KEY-----` |

Detection is based on regular expressions. It can miss secrets in other formats.

## Safety

- **Sessions in use are skipped.** Files modified within `--min-age` (30 minutes by default) are not touched. `fix --min-age=0` prints a warning.
- **Atomic writes.** `fix` writes to a temporary file in the same folder, then renames it over the original. If anything fails, the original file is left as it was.
- **Valid JSON.** If replacing a secret would break a valid JSON line, that line is left unchanged.
- **Timestamps are kept.** `fix` restores the original modification time, so the order of sessions in `/resume` does not change.
- **No backups.** A backup would keep the secret on disk. Run `fix --dry-run` first instead.

## Limitations

- **Rotate the key.** If a secret was written to a transcript, it may already have been copied elsewhere. The real fix is to revoke the key and issue a new one. `ccscrub` only lowers the risk of it leaking from these files in the future.
- **Resumed sessions.** After `fix`, the model in a resumed session does not know the removed value. Its earlier messages may no longer match the history, and it may get confused about what happened.
- Only Claude Code's `~/.claude/projects` is scanned.

## License

[MIT](LICENSE)
