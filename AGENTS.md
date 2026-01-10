# AGENTS.md

このファイルは、このAstroプロジェクトリポジトリで作業するエージェントコーディングエージェント向けのガイドラインとコマンドを含んでいます。

## プロジェクト概要

これはTypeScriptを厳格な型チェックで使用するAstro 5.16.8プロジェクトです。最小限の構成を持つ基本的なスターターテンプレートです。

## ビルドコマンド

### 開発
- `npm run dev` - `localhost:4321`でローカル開発サーバーを起動
- `npm run build` - 本番用サイトを`./dist/`にビルド
- `npm run preview` - ローカルでビルドをプレビュー

### Astro CLI
- `npm run astro ...` - Astro CLIコマンドを実行
- `npm run astro add` - インテグレーションを追加
- `npm run astro check` - 型チェック
- `npm run astro -- --help` - Astro CLIのヘルプを表示

### テスト
このプロジェクトでは現在テストフレームワークが設定されていません。テストを追加する場合は、以下のパターンに従ってください：
- 単体テスト: `*.test.ts`または`*.spec.ts`
- 統合テスト: `test/**/*.test.ts`
- コンポーネントテスト: `src/components/**/*.test.astro`

## コードスタイルガイドライン

### ファイル構造
- `src/pages/` - ルートページ（ファイルベースルーティング）
- `src/components/` - 再利用可能なAstroコンポーネント
- `src/layouts/` - ページレイアウトコンポーネント
- `src/assets/` - 静的アセット（画像など）
- `public/` - 直接提供される静的ファイル

### Astroコンポーネント
- コンポーネントには`.astro`拡張子を使用
- フロントマタースクリプトブロックは`---`区切りを使用
- インポート文はフロントマターの先頭に配置
- コンポーネントのpropsはフロントマターで`interface Props`を介して定義
- 型安全性のためにTypeScriptインターフェースを使用

### TypeScript設定
- 厳格モードが有効（`extends: "astro/tsconfigs/strict"`）
- `dist/`を除くすべてのファイルが対象
- 必要に応じてJSファイルで`@ts-check`コメントを使用

### インポートスタイル
- `../`と`./`表記で相対インポートを使用
- Astroコンポーネントは完全な拡張子でインポート: `import Component from './Component.astro'`
- アセットインポート: `import logo from '../assets/logo.svg'`

### 命名規則
- コンポーネント: PascalCase（例: `Welcome.astro`, `Layout.astro`）
- ファイルとディレクトリ: アセットはkebab-case、コンポーネントはPascalCase
- 変数: camelCase
- CSSクラス: kebab-case
- ID: kebab-case

### HTML/JSXスタイル
- セマンティックなHTML5要素を使用
- 空要素は自己終了タグを使用: `<img />`, `<br />`, `<hr />`
- 属性は二重引用符を使用
- Astro属性: 動的値には`{variable}`構文を使用
- ファーストビュー画像には`fetchpriority="high"`を使用

### CSSガイドライン
- コンポーネントではスコープ付き`<style>`ブロックを使用
- モバイルファーストのレスポンシブデザインで`@media`クエリを使用
- 必要に応じてテーマ用にCSSカスタムプロパティを使用
- 最新のCSS機能を使用（Grid, Flexboxなど）
- 絶対に必要でない限り`!important`は避ける

### エラーハンドリング
- 非同期操作にはtry-catchブロックを使用
- コンポーネントのフロントマターでpropsを検証
- コンパイル時の型チェックにTypeScriptインターフェースを使用
- 欠落しているアセットを適切に処理

### パフォーマンス
- 適切な形式で画像を最適化（WebP, AVIF）
- 重要な画像には`fetchpriority`属性を使用
- 重要でない画像は遅延読み込み
- 未使用のインポートをツリーシェーキングしてバンドルサイズを最小化

### アクセシビリティ
- 画像には常に`alt`属性を含める
- 適切な見出し階層を使用（h1, h2, h3など）
- 十分なカラーコントラストを確保
- 必要に応じてARIAラベルを追加
- セマンティックなHTML要素を使用

### Git規約
- 可能な場合は慣例的なコミットメッセージを使用
- 関連する変更をまとめてコミット
- `dist/`ディレクトリをコミットしない
- コミットは焦点を絞り、原子単位に保つ

## 開発ワークフロー

1. `npm run dev`で開発サーバーを起動
2. `src/`内のコンポーネントを変更
3. ブラウザはファイル変更時に自動リロード
4. `npm run build`で本番ビルドをテスト
5. `npm run preview`で本番ビルドをローカルでテスト

## 新機能の追加

- 新しいページ: `src/pages/`に追加（ファイルベースルーティング）
- 新しいコンポーネント: `src/components/`に追加
- 新しいレイアウト: `src/layouts/`に追加
- 新しいアセット: `src/assets/`または`public/`に追加
- 新しいインテグレーション: `npm run astro add [integration]`を使用

## 型安全性

- すべてのコンポーネントはTypeScriptを使用すべき
- コンポーネントのprops用にインターフェースを定義
- JavaScriptファイルでは`@ts-check`を使用
- `npm run astro check`で型を検証

## 一般的なパターン

### コンポーネントProps
```typescript
---
interface Props {
  title: string;
  description?: string;
}

const { title, description = "" } = Astro.props;
---
```

### アセットインポート
```typescript
---
import logo from '../assets/logo.svg';
---
<img src={logo.src} alt="Logo" />
```

### 条件付きレンダリング
```astro
{showContent && (
  <div>Content here</div>
)}
```

## ツールと拡張機能

- 推奨: Astro VS Code拡張機能
- 推奨: TypeScriptとJavaScript言語機能
- オプション: フォーマット用Prettier（必要に応じてプロジェクトに追加）
- オプション: リント用ESLint（必要に応じてプロジェクトに追加）