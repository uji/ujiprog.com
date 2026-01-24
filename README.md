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

## OG画像生成のためのフォント設定

記事の OG 画像を生成するには、以下のフォントファイルが必要です：

- `fonts/DMSans_36pt-Regular.ttf` - ASCII/英数字用
 - [Google Fonts](https://fonts.google.com/specimen/DM+Sans) からダウンロード
- `fonts/NotoSansJP-Regular.ttf` - 日本語用
 - [Google Fonts](https://fonts.google.com/noto/specimen/Noto+Sans+JP) からダウンロード
