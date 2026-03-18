# MCP Practice Project

このリポジトリは、Model Context Protocol (MCP) を Go 言語で実装し、動作を確認するための練習用プロジェクトです。
MCP のコアロジック、サーバーの実装、およびそれを利用するクライアントの実装が含まれています。

## プロジェクト構成

- **mcp-core/**: MCP のコアライブラリ。
  - `core.go`, `dto.go`: JSON-RPC 2.0 に基づくリクエスト/レスポンス構造体や、MCP プロトコル（Initialize, ListTools 等）の定義。
  - `stdio_client.go`, `stdio_server.go`: 標準入出力（stdio）を介した MCP 通信のヘルパー。
- **mcp-date/**: サンプルの MCP サーバー実装。
  - 都市ごとの現在時刻の取得、タイムゾーンのリスト表示、特定の都市の時刻を尋ねるプロンプトの提供などを行います。
- **client/**: `mcp-date` サーバーと通信するサンプルクライアント。
  - サーバーをサブプロセスとして起動し、stdio 経由で MCP メソッド（Initialize, ListTools, CallTool 等）を呼び出します。

## 主な機能

### MCP Core
- JSON-RPC 2.0 準拠のメッセージング
- MCP プロトコルの各種パラメータ（Capabilities, Tools, Prompts, Resources）の定義
- StdIO を用いたメッセージの送受信処理

### Date Server (`mcp-date`)
- **Tools**: `get_city_time`, `list_timezones`
- **Prompts**: `city_time_prompt`
- **Resources**: 各都市の現在の時刻情報を提供

## 使い方

### 1. サーバーのビルド
```bash
cd mcp-date
go build -o date-server main.go
```

### 2. クライアントの実行
クライアントはデフォルトで `./mcp-date/date-server` を呼び出すように構成されています。

```bash
cd client
go run main.go
```

特定のツールを直接呼び出す場合：
```bash
go run main.go -tool get_city_time -args '{"city": "Tokyo"}'
```

## テスト
```bash
go test ./...
```
