---
name: astro-development
description: Guide Astro development
---

# Astro開発スキル

## クイックリファレンス

### 基本コマンド
```bash
# 開発
npm run dev          # localhost:4321で開発サーバーを起動
npm run build        # 本番用ビルドを./dist/に生成
npm run preview      # 本番ビルドをローカルでプレビュー

# Astro CLI
npm run astro add [integration]    # インテグレーションを追加
npm run astro check                # 型チェック
npm run astro -- --help            # Astro CLIヘルプ
```

### ファイル構造
```
src/
├── pages/          # ルートページ（ファイルベースルーティング）
├── components/     # 再利用可能なAstroコンポーネント
├── layouts/        # ページレイアウトコンポーネント
├── assets/         # 静的アセット（画像など）
└── styles/         # グローバルCSSファイル

public/             # 直接提供される静的ファイル
```

### よく使うインポートパターン
```typescript
// Astroコンポーネント
import Component from './Component.astro';
import Layout from '../layouts/Layout.astro';

// アセット
import logo from '../assets/logo.svg';

// スタイル
import '../styles/global.css';
```

## 開発ワークフロー

### 1. 開発開始
```bash
npm run dev
```
- `localhost:4321`で開く
- ファイル変更時に自動リロード

### 2. コンポーネント作成
- ページ: `src/pages/[route].astro`
- コンポーネント: `src/components/[Component].astro`
- レイアウト: `src/layouts/[Layout].astro`

### 3. 型チェック
```bash
npm run astro check
```
- TypeScriptの型を検証
- コンパイルエラーをチェック

### 4. ビルド検証
```bash
npm run build
npm run preview
```
- 本番ビルドをテスト
- デプロイ前にローカルで検証

## コードパターン

### 基本的なコンポーネント構造
```astro
---
// フロントマター - サーバーサイドコード
import Layout from '../layouts/Layout.astro';

interface Props {
  title: string;
  description?: string;
}

const { title, description = "" } = Astro.props;
---

<!-- HTMLテンプレート -->
<Layout>
  <h1>{title}</h1>
  <p>{description}</p>
</Layout>

<style>
  /* コンポーネントスコープCSS */
  h1 {
    color: #333;
  }
</style>
```

### Propsインターフェースパターン
```typescript
---
interface Props {
  // 必須
  title: string;
  
  // オプション（デフォルト値あり）
  description?: string;
  showContent?: boolean;
  
  // 複雑な型
  items: Array<{ id: string; name: string }>;
}

const { 
  title, 
  description = "デフォルト説明", 
  showContent = true,
  items 
} = Astro.props;
---
```

### アセットインポートパターン
```astro
---
import logo from '../assets/logo.svg';
import backgroundImage from '../assets/hero-bg.jpg';
---

<!-- 最適化された画像 -->
<img src={logo.src} alt="Logo" width="115" height="48" />
<img 
  src={backgroundImage.src} 
  alt="ヒーロー背景" 
  fetchpriority="high"
/>
```

### 条件付きレンダリング
```astro
---
interface Props {
  showContent?: boolean;
  user?: { name: string } | null;
}

const { showContent = true, user } = Astro.props;
---

{showContent && (
  <div>
    <p>コンテンツが表示されています</p>
    {user && <p>ようこそ、{user.name}さん！</p>}
  </div>
)}
```

### レイアウトパターン
```astro
---
// src/layouts/Layout.astro
interface Props {
  title: string;
}

const { title } = Astro.props;
---

<!doctype html>
<html lang="ja">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <title>{title}</title>
  </head>
  <body>
    <slot />
  </body>
</html>
```

## よく使うタスク

### 新しいページを追加
1. `src/pages/[route].astro`を作成
2. レイアウトコンポーネントを追加
3. `localhost:4321/[route]`でテスト

### 新しいコンポーネントを追加
1. `src/components/[Component].astro`を作成
2. Propsインターフェースを定義
3. ページ/レイアウトでインポート

### インテグレーションを追加
```bash
npm run astro add [integration-name]
```

### スタイルを追加
- コンポーネントスコープ: コンポーネント内の`<style>`ブロック
- グローバル: レイアウトでCSSをインポート
- CSSモジュール: `import styles from './styles.module.css'`

## エラーハンドリング

### 必須Propsの検証
```typescript
---
interface Props {
  title: string;
}

// 必須propsを検証
if (!Astro.props.title) {
  throw new Error("Titleプロップは必須です");
}
---
```

### アセット読み込み
```astro
---
import { getImage } from 'astro:assets';
import imageSrc from '../assets/image.jpg';

const image = await getImage({ src: imageSrc, format: 'webp' });
---

<img {...image} alt="説明" />
```

## 詳細ガイド

包括的なパターンとベストプラクティスについては、以下を参照してください：

### 📋 コードパターン＆テンプレート
**ファイル**: `.claude/skills/astro/PATTERNS.md`
- 詳細なコンポーネントパターン（ページ、レイアウト、インタラクティブコンポーネント）
- データ取得パターン（コンテンツコレクション、API連携）
- フォームとナビゲーションパターン
- ユーティリティコンポーネント（日付フォーマット、SEO）

### 🎯 ベストプラクティスガイドライン
**ファイル**: `.claude/skills/astro/BEST_PRACTICES.md`
- パフォーマンス最適化テクニック
- アクセシビリティ（a11y）ガイドライン
- セキュリティベストプラクティス
- SEO最適化戦略
- コード品質基準
- テスト手法
- Gitワークフロー規約
- デプロイ最適化

### 🚀 クイックスタート参照
- **新しいコンポーネント**: PATTERNS.mdの「コンポーネントパターン」から基本パターンをコピー
- **パフォーマンスチェック**: BEST_PRACTICES.mdの「パフォーマンス最適化」を確認
- **アクセシビリティ監査**: BEST_PRACTICES.mdの「アクセシビリティガイドライン」に従う
- **SEO設定**: PATTERNS.mdの「SEOパターン」＋BEST_PRACTICES.mdの「SEOベストプラクティス」を実装

## プロジェクト固有の注意事項

- **言語**: 日本語（lang="ja"）
- **TypeScript**: 厳格モード有効
- **スタイリング**: グローバルCSSは`src/styles/global.css`
- **ビルドターゲット**: モダンブラウザ
- **デプロイ**: 静的サイト生成

## テスト

現在テストフレームワークは設定されていません。テストを追加する場合：
- 単体テスト: `*.test.ts`または`*.spec.ts`
- 統合テスト: `test/**/*.test.ts`
- コンポーネントテスト: `src/components/**/*.test.astro`