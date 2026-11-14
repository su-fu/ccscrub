# ccscrub

[English](README.md) | [日本語](README.ja.md) | [中文](README.zh.md)

> **非公式ツールです。** Anthropicとは無関係であり、公認・提携もしていません。

`ccscrub` は、Claude Codeのセッションのトランスクリプト（`~/.claude/projects/**/*.jsonl`）から、APIキーやトークンなどのシークレットを探し、その場でマスクするCLIです。すべての行が正しいJSONのまま保たれるので、マスクした後もセッションを再開できます。

> **これは学習プロジェクトです。** Go（ストリーミングI/O、`regexp`、安全なファイルの書き換え）を学ぶために作りました。実際に使うなら、[agentleaks](https://github.com/Arthur031221/agentleaks) や [claude-code-scrubify](https://github.com/benjamin-awd/claude-code-scrubify) など、より充実したツールを検討してください。対応しているAIツールも、検出できるシークレットの種類も多いです。

## 作った理由

シークレットを伏せるフック（`PreToolUse` / `PostToolUse`）が制御できるのは、モデルに送られる内容だけです。Claude Codeがディスクに書き込むトランスクリプトには、元の値が残ることがあります。`ccscrub` は、その「書き込まれた後」の部分を扱います。

## インストール

Go 1.27.4 以降が必要です。

```sh
go install github.com/su-fu/ccscrub/cmd/ccscrub@latest
```

## 使い方

おすすめの流れは **scan → 確認 → fix** です。

```sh
ccscrub scan             # シークレットを探す（読み取りのみ）
ccscrub fix --dry-run    # 置き換わる場所を表示する（読み取りのみ）
ccscrub fix              # その場で置き換える
```

### 例

```sh
$ ccscrub scan --lang=ja
C:\Users\you\.claude\projects\my-project\0a1b2c3d.jsonl:42:118: github-token ghp_****AAAA
1 件のシークレットが 1 個のファイルで見つかりました
```

各行は `ファイル:行:列: ルールID プレビュー` の形式です。プレビューは先頭と末尾の4文字だけを表示するので、`ccscrub` がシークレット全体を出力することはありません。最後の集計行は標準エラー出力に出ます。

`fix` の後、トランスクリプトの中のシークレットは次のようになります。

```json
{"text":"[REDACTED:github-token]"}
```

### フラグ

| フラグ | コマンド | 説明 |
|---|---|---|
| `--dry-run` | `fix` | ファイルを変更せずに、変更される内容を表示する |
| `--min-age` | `scan`、`fix` | この時間以内に更新されたトランスクリプトは、使用中の可能性があるのでスキップする（既定値 `30m`） |
| `--lang` | `scan`、`fix` | メッセージの言語：`en`、`ja`、`zh`（既定値 `en`）。環境変数 `CCSCRUB_LANG` でも指定できる |

### 終了コード

| コード | 意味 |
|---|---|
| `0` | シークレットが見つからなかった（`scan`、`fix --dry-run`）、または `fix` が完了した |
| `1` | シークレットが見つかった（`scan`、`fix --dry-run`） |
| `2` | エラー（フラグの誤り、`~/.claude/projects` が見つからない、など） |

## 検出できるもの

| ルールID | 先頭の例 |
|---|---|
| `github-token` | `ghp_`、`gho_`、`ghu_`、`ghs_`、`ghr_`、`github_pat_` |
| `aws-access-key-id` | `AKIA`、`ASIA` |
| `anthropic-api-key` | `sk-ant-api03-` |
| `anthropic-admin-api-key` | `sk-ant-admin01-` |
| `openai-api-key` | `sk-proj-`、`sk-svcacct-`、`sk-admin-`、古い形式の `sk-` |
| `slack-token` | `xoxb-`、`xoxp-`、`xoxe-` |
| `google-api-key` | `AIza` |
| `private-key` | `-----BEGIN ... PRIVATE KEY-----` |

検出は正規表現によるものです。これ以外の形式のシークレットは見逃すことがあります。

## 安全のための仕組み

- **使用中のセッションはスキップします。** `--min-age`（既定値30分）以内に更新されたファイルには触りません。`fix --min-age=0` では警告を出します。
- **アトミックな書き込み。** `fix` は同じフォルダの一時ファイルに書き込み、それを元のファイルに rename します。途中で失敗しても、元のファイルはそのまま残ります。
- **JSONを壊しません。** 置き換えると正しいJSONの行が壊れる場合は、その行を変更しません。
- **更新日時を保ちます。** `fix` は元の更新日時に戻すので、`/resume` の一覧の並び順は変わりません。
- **バックアップは作りません。** バックアップにはシークレットが残ってしまうためです。代わりに、先に `fix --dry-run` を実行してください。

## 制限事項

- **キーを無効化して、再発行してください。** トランスクリプトに書き込まれたシークレットは、すでに別の場所にコピーされている可能性があります。根本的な対処は、キーを無効化して新しいキーを発行することです。`ccscrub` は、これらのファイルから今後漏れるリスクを下げるだけです。
- **再開したセッション。** `fix` の後でセッションを再開すると、モデルは消された値を知りません。過去の発言と履歴がかみ合わなくなり、何が起きたかについてモデルが混乱することがあります。
- 調べるのは Claude Code の `~/.claude/projects` だけです。

## ライセンス

[MIT](LICENSE)
