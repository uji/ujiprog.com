# ujiprog.com

- Uses [`syumai/workers`](https://github.com/syumai/workers) package to run an HTTP server.

## Requirements

- Node.js
- Go 1.24.0 or later

## Commands

ビルドスクリプトは Makefile に統一しています。package.json には npm スクリプトを定義しません。

```bash
make run               # Air を使用してホットリロードで開発サーバーを起動
make build             # Go Wasm バイナリをビルド
make dev               # ローカル R2 にアセットを配置して wrangler dev を起動
make generate-articles # 記事 HTML ページを生成
make deploy            # R2 にアセットをアップロードし、Cloudflare Workers にデプロイ
```
