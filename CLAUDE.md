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
npm run build                                      # Workers デプロイ用の Go Wasm バイナリをビルド
npm run deploy                                     # Cloudflare Workers にデプロイ
GOOS=js GOARCH=wasm go build -o ./build/app.wasm . # Go コードのビルドが通ることを確認する
make generate-articles                             # HTMLページの生成
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
- `package.json`: Node.js 依存関係とビルドスクリプト
- `Makefile`: 追加のビルドコマンド

### 開発ワークフロー
1. 開発中のホットリロードに `make run` を使用
2. 変更をコミットする前に `go test ./...` を実行
3. コードフォーマットを確保するために `go fmt ./...` を使用
4. デプロイ前に `npm run build` でビルド
5. /verify-browser で描画を確認

