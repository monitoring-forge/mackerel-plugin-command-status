# mackerel-plugin-command-status

Mackerel カスタムプラグイン。指定したコマンドを実行し、その**終了ステータス**と**実行時間**をメトリクスとして収集します。

 cron で定期的にバッチ処理やスクリプトを実行し、その成否・所要時間を Mackerel で監視・可視化したい場合に有用です。

## メトリクス

実行ごとに以下の 2 つのメトリクスが出力されます。

| メトリクス名 | 概要 | 値の例 |
|---|---|---|
| `command-status.time-taken.<name>` | コマンドの実行時間（秒） | `0.008186` |
| `command-status.exit-code.<name>` | コマンドの終了ステータスコード | `0` |

### 終了ステータスコードの意味

| コード | 意味 |
|---|---|
| `0` | 正常終了（OK） |
| `1〜126` | コマンド自体の終了ステータス（アプリが返す値） |
| `127` | コマンドの起動に失敗した場合（存在しないコマンドなど） |
| `137` | タイムアウトにより `SIGKILL` された場合（`9 + 128`） |

> **注意**: コマンドがシグナルで終了した場合、ExitCode は `-1` の値になるため、本プラグインは `127` に正規化して出力します。

## 使い方

### 実行形式

```
mackerel-plugin-command-status [OPTIONS] -- command [args...]
```

`--` 以降に実行したいコマンドとその引数を指定します。

### オプション

| オプション | 略称 | 説明 | 既定値 |
|---|---|---|---|
| `--timeout=` | — | コマンドのタイムアウト時間 | `30s` |
| `--name=` | `-n` | メトリクス名（ユニークな名前を付けましょう） | *必須* |
| `--quiet` | `-q` | サブコマンドのエラー出力を抑制する | `false` |
| `--version` | `-v` | バージョン情報を表示する | — |
| `--help` | `-h` | ヘルプを表示する | — |

### 実行例

#### 正常終了する場合

```bash
$ mackerel-plugin-command-status --name update-cache -- /path/to/fetch-cache
command-status.time-taken.update-cache  0.008186  1606958816
command-status.exit-code.update-cache   0         1606958816
```

#### エラー終了する場合（ログ出力あり）

```bash
$ mackerel-plugin-command-status -n test -- false
2025/04/15 16:17:07 Command false exit with err: exit status 1
command-status.time-taken.test  0.001459  1744701427
command-status.exit-code.test   1         1744701427
```

#### エラー終了する場合（`--quiet` でログ抑制）

```bash
$ mackerel-plugin-command-status --quiet -n test -- false
command-status.time-taken.test  0.001316  1744701374
command-status.exit-code.test   1         1744701374
```

#### タイムアウトした場合

```bash
$ mackerel-plugin-command-status -n sleep --timeout 3s -- sleep 30
2020/12/03 10:29:11 Command sleep timeout. killed
command-status.time-taken.sleep  3.016507  1614308642
command-status.exit-code.sleep   137       1614308642
```

タイムアウト時はプロセスが `SIGKILL` されるため、終了コードは `137` になります。

## Mackerel への設定

### mackerel-agent.conf

```ini
[plugin.metrics.update-cache]
command = "/path/to/mackerel-plugin-command-status --name update-cache --timeout 10s -- /path/to/cmd-fetch-cache"
```

### Monitor（監視ルール）の設定

Mackerel の Web UI で以下のような Monitor を設定します。

- **Condition**: `warning: > 0` 、`critical: > 0` （終了コードが 0 以外なら警告/異常）
- **Max Check Attempts**: 目的の SLO に応じて適切に設定（例: `2` で warning、`3` で critical など）

実行時間に監視を追加したい場合は、`command-status.time-taken.<name>` を対象に閾値を設定します。

## Install

### Release ページからダウンロード

[GitHub Releases](https://github.com/monitoring-forge/mackerel-plugin-command-status/releases) から該当プラットフォームのバイナリをダウンロードし、`mackerel-agent` が実行できる場所に配置してください。

### mkr コマンドでインストール

```bash
mkr plugin install monitoring-forge/mackerel-plugin-command-status
```

## License

MIT
