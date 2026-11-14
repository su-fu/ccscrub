# ccscrub

[English](README.md) | [日本語](README.ja.md) | [中文](README.zh.md)

> **非官方工具。** 与 Anthropic 无关，也未获得其认可或合作。

`ccscrub` 是一个命令行工具，用于在 Claude Code 的会话记录（`~/.claude/projects/**/*.jsonl`）中查找 API 密钥、令牌等机密信息，并在原文件中将其脱敏。每一行都保持为有效的 JSON，因此脱敏后仍然可以恢复会话。

> **这是一个学习项目。** 我开发它是为了学习 Go（流式 I/O、`regexp`、安全地改写文件）。如需实际使用，请考虑功能更完善的工具，例如 [agentleaks](https://github.com/Arthur031221/agentleaks) 或 [claude-code-scrubify](https://github.com/benjamin-awd/claude-code-scrubify)。它们支持更多 AI 工具，也能检测更多种类的机密信息。

## 为什么开发这个工具

用于隐藏机密信息的钩子（`PreToolUse` / `PostToolUse`）只能控制发送给模型的内容。Claude Code 写入磁盘的会话记录中仍可能保留原始值。`ccscrub` 处理的正是"写入之后"的这一部分。

## 安装

需要 Go 1.27.4 或更高版本。

```sh
go install github.com/su-fu/ccscrub/cmd/ccscrub@latest
```

## 使用方法

推荐的流程是：**scan → 预览 → fix**。

```sh
ccscrub scan             # 查找机密信息（只读）
ccscrub fix --dry-run    # 显示将被替换的位置（只读）
ccscrub fix              # 在原文件中替换
```

### 示例

```sh
$ ccscrub scan --lang=zh
C:\Users\you\.claude\projects\my-project\0a1b2c3d.jsonl:42:118: github-token ghp_****AAAA
发现 1 处机密信息，涉及 1 个文件
```

每一行的格式为 `文件:行:列: 规则ID 预览`。预览只显示开头和结尾各 4 个字符，因此 `ccscrub` 不会输出完整的机密信息。最后的汇总行输出到标准错误。

执行 `fix` 后，会话记录中的机密信息会变成：

```json
{"text":"[REDACTED:github-token]"}
```

### 参数

| 参数 | 命令 | 说明 |
|---|---|---|
| `--dry-run` | `fix` | 不修改任何文件，只显示将要更改的内容 |
| `--min-age` | `scan`、`fix` | 跳过在此时间内修改过的会话记录，因为它们可能正在使用（默认 `30m`） |
| `--lang` | `scan`、`fix` | 消息语言：`en`、`ja`、`zh`（默认 `en`）。也可以通过环境变量 `CCSCRUB_LANG` 设置 |

### 退出码

| 代码 | 含义 |
|---|---|
| `0` | 未发现机密信息（`scan`、`fix --dry-run`），或 `fix` 已完成 |
| `1` | 发现了机密信息（`scan`、`fix --dry-run`） |
| `2` | 错误（参数错误、找不到 `~/.claude/projects` 等） |

## 可检测的内容

| 规则 ID | 前缀示例 |
|---|---|
| `github-token` | `ghp_`、`gho_`、`ghu_`、`ghs_`、`ghr_`、`github_pat_` |
| `aws-access-key-id` | `AKIA`、`ASIA` |
| `anthropic-api-key` | `sk-ant-api03-` |
| `anthropic-admin-api-key` | `sk-ant-admin01-` |
| `openai-api-key` | `sk-proj-`、`sk-svcacct-`、`sk-admin-`、旧格式的 `sk-` |
| `slack-token` | `xoxb-`、`xoxp-`、`xoxe-` |
| `google-api-key` | `AIza` |
| `private-key` | `-----BEGIN ... PRIVATE KEY-----` |

检测基于正则表达式，可能会漏掉其他格式的机密信息。

## 安全机制

- **跳过正在使用的会话。** 不会修改在 `--min-age`（默认 30 分钟）内更新过的文件。使用 `fix --min-age=0` 时会显示警告。
- **原子写入。** `fix` 先写入同一文件夹中的临时文件，再将其重命名为原文件。如果中途失败，原文件保持不变。
- **保持 JSON 有效。** 如果替换会破坏一个有效的 JSON 行，则不修改该行。
- **保留时间戳。** `fix` 会恢复原来的修改时间，因此 `/resume` 列表中会话的顺序不会改变。
- **不创建备份。** 备份会把机密信息留在磁盘上。请先运行 `fix --dry-run` 进行确认。

## 限制

- **请吊销并重新签发密钥。** 写入会话记录的机密信息可能已经被复制到其他地方。根本的解决办法是吊销该密钥并签发新的密钥。`ccscrub` 只能降低今后从这些文件中泄露的风险。
- **恢复的会话。** 执行 `fix` 后恢复会话时，模型并不知道被删除的值。它之前的发言可能与历史记录对不上，模型可能会对发生的事情感到困惑。
- 只扫描 Claude Code 的 `~/.claude/projects`。

## 许可证

[MIT](LICENSE)
