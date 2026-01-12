# Astroベストプラクティス

## パフォーマンス最適化

### 1. 画像最適化
```astro
---
import { getImage } from 'astro:assets';
import heroImage from '../assets/hero.jpg';

// 異なるフォーマットで最適化
const optimizedImage = await getImage({
  src: heroImage,
  format: 'webp', // 'avif', 'png', 'jpg'も可
  width: 1200,
  height: 600,
});
---

<!-- 最適化された画像を使用 -->
<img {...optimizedImage} alt="ヒーロー画像" loading="eager" />

<!-- レスポンシブ画像の場合 -->
<picture>
  <source srcset={optimizedImage.src} type="image/webp" />
  <img src={heroImage.src} alt="フォールバック" />
</picture>
```

### 2. コード分割
```astro
---
// 必要なときだけコンポーネントをロード
const InteractiveComponent = Astro.resolve('./InteractiveComponent.astro');
---

<!-- 必要なときだけハイドレーション -->
<InteractiveComponent client:load />
<!-- 重要でないコンポーネントはclient:idle -->
<InteractiveComponent client:idle />
```

### 3. バンドル最適化
```javascript
// astro.config.mjs
export default defineConfig({
  vite: {
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['react', 'react-dom'],
            ui: ['@headlessui/react'],
          },
        },
      },
    },
  },
});
```

## アクセシビリティガイドライン

### 1. セマンティックHTML構造
```astro
---
// 適切な見出し階層を使用
---

<header>
  <h1>サイトタイトル</h1>
</header>

<main>
  <section aria-labelledby="section1-title">
    <h2 id="section1-title">セクションタイトル</h2>
    <h3>サブセクションタイトル</h3>
  </section>
</main>

<aside aria-labelledby="sidebar-title">
  <h2 id="sidebar-title">サイドバー</h2>
</aside>

<footer>
  <p>&copy; 2024 サイト名</p>
</footer>
```

### 2. フォームアクセシビリティ
```astro
---
// 適切なフォームラベリング
---

<form>
  <div class="form-group">
    <label for="email">メールアドレス</label>
    <input 
      type="email" 
      id="email" 
      name="email" 
      required 
      aria-describedby="email-help"
      aria-invalid="false"
    />
    <div id="email-help" class="help-text">
      あなたのメールアドレスを他の誰とも共有しません。
    </div>
  </div>
  
  <button type="submit">送信</button>
</form>
```

### 3. ナビゲーションアクセシビリティ
```astro
---
// アクセシブルなナビゲーション
---

<nav role="navigation" aria-label="メインナビゲーション">
  <ul>
    <li><a href="/" aria-current="page">ホーム</a></li>
    <li><a href="/about/">概要</a></li>
    <li><a href="/blog/">ブログ</a></li>
  </ul>
</nav>

<!-- キーボードユーザー向けスキップリンク -->
<a href="#main-content" class="skip-link">メインコンテンツへスキップ</a>
```

## セキュリティベストプラクティス

### 1. コンテンツセキュリティポリシー
```astro
---
// src/layouts/BaseLayout.astro
const csp = `
  default-src 'self';
  script-src 'self' 'unsafe-inline';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  font-src 'self';
`;
---

<head>
  <meta http-equiv="Content-Security-Policy" content={csp} />
</head>
```

### 2. 入力検証
```astro
---
// ユーザー入力を検証
interface Props {
  searchQuery?: string;
}

const { searchQuery } = Astro.props;

// 検索クエリをサニタイズ
const sanitizedQuery = searchQuery 
  ? searchQuery.replace(/[<>]/g, '').trim()
  : '';
---

<!-- サニタイズした入力を使用 -->
{sanitizedQuery && (
  <p>検索結果: {sanitizedQuery}</p>
)}
```

### 3. 環境変数
```javascript
// .env.example
PUBLIC_API_URL=https://api.example.com
SECRET_API_KEY=your-secret-key

// astro.config.mjs
export default defineConfig({
  vite: {
    define: {
      'process.env.PUBLIC_API_URL': process.env.PUBLIC_API_URL,
    },
  },
});
```

## SEOベストプラクティス

### 1. メタタグ構造
```astro
---
// src/components/SEOHead.astro
interface Props {
  title: string;
  description: string;
  ogImage?: string;
  keywords?: string[];
  canonical?: string;
}

const { title, description, ogImage, keywords, canonical } = Astro.props;
const siteUrl = Astro.site.href;
---

<!-- 必須メタタグ -->
<title>{title}</title>
<meta name="description" content={description} />
{keywords && <meta name="keywords" content={keywords.join(', ')} />}
{canonical && <link rel="canonical" href={`${siteUrl}${canonical}`} />}

<!-- Open Graphタグ -->
<meta property="og:title" content={title} />
<meta property="og:description" content={description} />
<meta property="og:type" content="website" />
<meta property="og:url" content={siteUrl} />
{ogImage && <meta property="og:image" content={ogImage} />}

<!-- Twitter Cardタグ -->
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content={title} />
<meta name="twitter:description" content={description} />
{ogImage && <meta name="twitter:image" content={ogImage} />}
```

### 2. 構造化データ
```astro
---
// 記事用のJSON-LD
<script type="application/ld+json">
{{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": title,
  "description": description,
  "author": {
    "@type": "Person",
    "name": author
  },
  "datePublished": publishDate.toISOString(),
  "dateModified": modifiedDate.toISOString(),
  "image": ogImage,
  "publisher": {
    "@type": "Organization",
    "name": "ujiprog.com",
    "url": siteUrl
  }
}}
</script>
```

### 3. サイトマップ生成
```javascript
// scripts/generate-sitemap.js
import { writeFileSync } from 'fs';
import { glob } from 'astro/loaders';

export async function generateSitemap(pages) {
  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${pages.map(page => `
  <url>
    <loc>${page.url}</loc>
    <lastmod>${page.lastModified}</lastmod>
    <changefreq>${page.changeFreq}</changefreq>
    <priority>${page.priority}</priority>
  </url>
`).join('')}
</urlset>`;

  writeFileSync('./public/sitemap.xml', sitemap);
}
```

## コード品質

### 1. TypeScriptベストプラクティス
```typescript
// 厳密な型付けを使用
interface BlogPost {
  id: string;
  title: string;
  content: string;
  publishDate: Date;
  tags: readonly string[];
  author: {
    name: string;
    email: string;
  };
}

// ユーティリティ型を使用
type PartialBlogPost = Partial<BlogPost>;
type BlogPostPreview = Pick<BlogPost, 'id' | 'title' | 'publishDate'>;

// 再利用可能コンポーネントにジェネリクスを使用
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => Astro.Component;
}
```

### 2. コンポーネント構成
```astro
---
// 小さなコンポーネントを組み合わせる
import Card from './Card.astro';
import Badge from './Badge.astro';
import Button from './Button.astro';

interface Props {
  post: BlogPost;
}

const { post } = Astro.props;
---

<Card>
  <h2>{post.title}</h2>
  <div class="meta">
    {post.tags.map(tag => <Badge>{tag}</Badge>)}
  </div>
  <p>{post.excerpt}</p>
  <Button href={`/blog/${post.id}`}>もっと読む</Button>
</Card>
```

### 3. エラーハンドリング
```astro
---
// グレースフルなエラーハンドリング
interface Props {
  postId: string;
}

const { postId } = Astro.props;

let post = null;
let error = null;

try {
  post = await getPost(postId);
} catch (err) {
  error = err;
  console.error('投稿の読み込み失敗:', err);
}
---

{error ? (
  <div class="error">
    <h2>投稿が見つかりません</h2>
    <p>申し訳ありません、お探しの投稿が見つかりませんでした。</p>
    <a href="/blog/">ブログに戻る</a>
  </div>
) : post ? (
  <article>{post.content}</article>
) : (
  <div class="loading">読み込み中...</div>
)}
```

## CSSベストプラクティス

### 1. CSSアーキテクチャ
```astro
---
// テーマ用にCSSカスタムプロパティを使用
---

<style define:vars={{ primaryColor: '#3b82f6', textColor: '#111827' }}>
  :root {
    --color-primary: var(--primaryColor);
    --color-text: var(--textColor);
  }
  
  .button {
    background-color: var(--color-primary);
    color: var(--color-text);
  }
</style>
```

### 2. レスポンシブデザイン
```astro
---
// モバイルファーストのレスポンシブデザイン
---

<style>
  .container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  @media (min-width: 768px) {
    .container {
      padding: 2rem;
    }
  }
  
  @media (min-width: 1024px) {
    .container {
      padding: 3rem;
    }
  }
</style>
```

### 3. パフォーマンス最適化CSS
```astro
---
// 効率的なセレクタを使用
---

<style>
  /* 良い: クラスベースセレクタ */
  .card {
    padding: 1rem;
  }
  
  .card-title {
    font-size: 1.5rem;
  }
  
  /* 避ける: ユニバーサルセレクタ */
  /* * { margin: 0; } */
  
  /* 避ける: 過剰修飾セレクタ */
  /* div.card-container ul.card-list li.card-item { } */
</style>
```

## テストガイドライン

### 1. コンポーネントテスト
```astro
---
// src/components/__tests__/Button.astro
import Button from '../Button.astro';

interface Props {
  variant?: 'primary' | 'secondary';
  size?: 'small' | 'medium' | 'large';
}

const { variant = 'primary', size = 'medium' } = Astro.props;
---

<div 
  class={`button button--${variant} button--${size}`}
  role="button"
  tabindex="0"
>
  <slot />
</div>

<style>
  .button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  
  .button--primary {
    background-color: #3b82f6;
    color: white;
  }
  
  .button--secondary {
    background-color: #6b7280;
    color: white;
  }
  
  .button--small {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }
  
  .button--large {
    padding: 0.75rem 1.5rem;
    font-size: 1.125rem;
  }
</style>
```

### 2. 統合テスト
```typescript
// src/test/api.test.ts
import { test, expect } from '@playwright/test';

test('ブログAPIが正しいデータを返す', async ({ request }) => {
  const response = await request.get('/api/posts');
  expect(response.ok()).toBeTruthy();
  
  const posts = await response.json();
  expect(posts).toHaveLength(10);
  expect(posts[0]).toHaveProperty('title');
  expect(posts[0]).toHaveProperty('content');
});
```

## Gitワークフロー

### 1. コミットメッセージ規約
```bash
# 機能コミット
git commit -m "feat(blog): 投稿一覧コンポーネントを追加"

# バグ修正
git commit -m "fix(nav): モバイルメニュートグルの問題を解決"

# ドキュメント
git commit -m "docs(readme): インストール手順を更新"

# リファクタリング
git commit -m "refactor(components): 共通ボタンスタイルを抽出"
```

### 2. ブランチ戦略
```bash
# 機能ブランチ
git checkout -b feature/blog-post-listing
git checkout -b fix/mobile-navigation

# リリースブランチ
git checkout -b release/v1.2.0

# ホットフィックスブランチ
git checkout -b hotfix/critical-security-patch
```

### 3. コードレビューチェックリスト
- [ ] コードがプロジェクトスタイルガイドに従っている
- [ ] TypeScript型が適切に定義されている
- [ ] コンポーネントがアクセシブルである
- [ ] パフォーマンスへの影響が検討されている
- [ ] 適切な場所にテストが含まれている
- [ ] ドキュメントが更新されている
- [ ] 機密データがコミットされていない

## デプロイベストプラクティス

### 1. ビルド最適化
```javascript
// astro.config.mjs
export default defineConfig({
  build: {
    format: 'directory', // ホスティングのため
  },
  output: 'static',
  compress: true,
});
```

### 2. 環境設定
```javascript
// astro.config.mjs
export default defineConfig({
  site: 'https://ujiprog.com',
  output: 'static',
  trailingSlash: 'always',
  
  // 環境固有設定
  experimental: {
    env: {
      schema: {
        PUBLIC_API_URL: String,
        SECRET_KEY: String,
      },
    },
  },
});
```

### 3. パフォーマンス監視
```astro
---
// パフォーマンス監視を追加
---

<script>
  // Core Web Vitals監視
  import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';

  function sendToAnalytics(metric) {
    // 分析サービスに送信
    console.log(metric);
  }

  getCLS(sendToAnalytics);
  getFID(sendToAnalytics);
  getFCP(sendToAnalytics);
  getLCP(sendToAnalytics);
  getTTFB(sendToAnalytics);
</script>
```