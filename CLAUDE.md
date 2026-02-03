# CLAUDE.md

このファイルには、このリポジトリで作業するエージェント型コーディングエージェント向けのガイドラインとコマンドが含まれています。

## プロジェクト概要

これは `syumai/workers` フレームワークを使用した Cloudflare Workers デプロイ用の Go ベースの Web アプリケーションです。このアプリケーションは HTTP ハンドラーを提供し、workers プラットフォーム用の WebAssembly としてビルドできます。

## 開発コマンド

### アプリケーションの実行
```bash
make run           # Air を使用してホットリロードで開発サーバーを起動
```

### ビルドとデプロイ
```bash
make build              # Go Wasm バイナリをビルド
make generate-articles  # HTMLページの生成
make deploy             # R2 にアセットをアップロードし、Cloudflare Workers にデプロイ
```

### ローカル開発 (Wrangler)
```bash
make dev-init           # 初回セットアップ（全アセットをローカルR2にアップロード）
make dev                # wrangler dev を起動（高速）
```

## コードスタイルガイドライン

### HTTP ハンドラー
- ルート登録に `http.HandleFunc` を使用
- 適切なヘッダーを設定 (Content-Type, CSP, セキュリティヘッダー)
- 異なる HTTP メソッドを適切に処理
- 適切な HTTP ステータスコードを返す

### セキュリティヘッダー
アプリケーションはデフォルトでこれらのセキュリティヘッダーを設定:
- `Content-Security-Policy`: `default-src 'self'; style-src 'unsafe-inline'; img-src 'self' blob:;`
- `Cross-Origin-Opener-Policy`: `same-origin`
- `Strict-Transport-Security`: `max-age=31536000; includeSubDomains; preload`

### HTML/テンプレートガイドライン
- セマンティックな HTML5 要素を使用
- 適切な DOCTYPE と lang 属性を含める
- スタイリングに Tailwind CSS クラスを使用 (既存コード通り)
- flexbox/grid ユーティリティでレスポンシブデザインを確保

### ファイル構造
- `main.go`: エントリーポイントと HTTP ハンドラー
- `go.mod`: Go モジュール定義
- `package.json`: Node.js 依存関係 (wrangler のみ、スクリプトは定義しない)
- `Makefile`: ビルド・開発・デプロイコマンド

### 重要: ビルドスクリプトの統一
- **ビルドスクリプトは Makefile に統一する**
- package.json には npm スクリプトを定義しない (依存関係の管理のみ)
- wrangler.jsonc の `build.command` も `make build` を使用する

### 開発ワークフロー

#### ローカル開発
```bash
make dev-init      # 初回セットアップ（全アセットをローカルR2にアップロード）
make run           # Air によるホットリロード開発（並列監視）
                   # - Go/HTML/CSS/JS 変更 → ビルド → wrangler dev 再起動
                   # - Markdown 変更 → 差分ファイルのみ処理 → R2 アップロード
```

#### 外部記事の取得（Zenn/note/SpeakerDeck）
```bash
make fetch-articles    # APIから最新記事一覧を取得 → public/articles.json
```

#### 本番デプロイ
```bash
make fetch-articles      # 外部記事を最新化（任意）
make generate-articles   # ローカル記事をビルド
make deploy              # R2 アップロード + Workers デプロイ
```

#### コマンド早見表
| コマンド | 用途 |
|----------|------|
| `make dev-init` | 初回セットアップ（全アセットアップロード） |
| `make run` | ホットリロード開発（Go/HTML + Markdown 並列監視） |
| `make dev` | wrangler dev のみ起動 |
| `make dev-sync-articles` | Git差分のMarkdownのみ処理・アップロード |
| `make fetch-articles` | 外部記事一覧を更新 |
| `make deploy` | 本番環境へデプロイ |

#### 品質チェック
```bash
make build         # ビルドエラーの確認（wasmビルド、go build ./... は使わない）
go test ./...      # テスト実行
go fmt ./...       # フォーマット
/verify-browser    # ブラウザで描画確認
```

