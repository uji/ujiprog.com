# Astroコードパターン

## コンポーネントパターン

### 1. ページコンポーネントパターン
```astro
---
// src/pages/blog/[slug].astro
import { getCollection } from 'astro:content';
import BlogLayout from '../../layouts/BlogLayout.astro';
import type { CollectionEntry } from 'astro:content';

export async function getStaticPaths() {
  const posts = await getCollection('blog');
  return posts.map((post) => ({
    params: { slug: post.slug },
    props: post,
  }));
}

interface Props {
  post: CollectionEntry<'blog'>;
}

const { post } = Astro.props;
const { Content, headings } = await post.render();
---

<BlogLayout title={post.data.title}>
  <article>
    <header>
      <h1>{post.data.title}</h1>
      <time datetime={post.data.pubDate}>
        {post.data.pubDate.toLocaleDateString()}
      </time>
    </header>
    <Content />
  </article>
</BlogLayout>
```

### 2. データ取得コンポーネント
```astro
---
// src/components/PostList.astro
import { getCollection } from 'astro:content';
import type { CollectionEntry } from 'astro:content';

interface Props {
  limit?: number;
  category?: string;
}

const { limit = 10, category } = Astro.props;

const posts = await getCollection('blog', ({ data }) => {
  if (category) {
    return data.category === category;
  }
  return true;
});

const filteredPosts = posts
  .sort((a, b) => b.data.pubDate.valueOf() - a.data.pubDate.valueOf())
  .slice(0, limit);
---

<section class="post-list">
  <h2>最新の投稿</h2>
  {filteredPosts.map((post) => (
    <article class="post-item">
      <h3>
        <a href={`/blog/${post.slug}/`}>{post.data.title}</a>
      </h3>
      <time datetime={post.data.pubDate}>
        {post.data.pubDate.toLocaleDateString()}
      </time>
      {post.data.description && (
        <p>{post.data.description}</p>
      )}
    </article>
  ))}
</section>

<style>
  .post-list {
    display: grid;
    gap: 2rem;
  }
  
  .post-item {
    padding: 1.5rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
  }
  
  .post-item h3 {
    margin: 0 0 0.5rem 0;
  }
  
  .post-item time {
    color: #6b7280;
    font-size: 0.875rem;
  }
</style>
```

### 3. インタラクティブコンポーネントパターン
```astro
---
// src/components/Counter.astro
interface Props {
  initial?: number;
  max?: number;
}

const { initial = 0, max = 100 } = Astro.props;
---

<div id="counter" data-initial={initial} data-max={max}>
  <button id="decrement">-</button>
  <span id="count">{initial}</span>
  <button id="increment">+</button>
</div>

<script>
  class Counter {
    constructor(element) {
      this.element = element;
      this.countEl = element.querySelector('#count');
      this.decrementBtn = element.querySelector('#decrement');
      this.incrementBtn = element.querySelector('#increment');
      
      this.count = parseInt(element.dataset.initial);
      this.max = parseInt(element.dataset.max);
      
      this.init();
    }
    
    init() {
      this.decrementBtn.addEventListener('click', () => this.decrement());
      this.incrementBtn.addEventListener('click', () => this.increment());
      this.updateDisplay();
    }
    
    decrement() {
      if (this.count > 0) {
        this.count--;
        this.updateDisplay();
      }
    }
    
    increment() {
      if (this.count < this.max) {
        this.count++;
        this.updateDisplay();
      }
    }
    
    updateDisplay() {
      this.countEl.textContent = this.count;
      this.decrementBtn.disabled = this.count === 0;
      this.incrementBtn.disabled = this.count === this.max;
    }
  }
  
  // カウンターを初期化
  document.getElementById('counter').forEach(el => new Counter(el));
</script>
```

## レイアウトパターン

### 1. 基本レイアウトパターン
```astro
---
// src/layouts/BaseLayout.astro
import '../styles/global.css';

interface Props {
  title: string;
  description?: string;
  ogImage?: string;
}

const { title, description, ogImage } = Astro.props;
---

<!doctype html>
<html lang="ja">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <link rel="icon" type="image/ico" href="/favicon.ico" />
    <meta name="generator" content={Astro.generator} />
    
    <title>{title}</title>
    {description && <meta name="description" content={description} />}
    
    <!-- Open Graph -->
    <meta property="og:title" content={title} />
    {description && <meta property="og:description" content={description} />}
    <meta property="og:type" content="website" />
    {ogImage && <meta property="og:image" content={ogImage} />}
    
    <!-- Twitter Card -->
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content={title} />
    {description && <meta name="twitter:description" content={description} />}
    {ogImage && <meta name="twitter:image" content={ogImage} />}
  </head>
  <body>
    <slot />
  </body>
</html>

<style is:global>
  body {
    font-family: system-ui, -apple-system, sans-serif;
    line-height: 1.6;
  }
</style>
```

### 2. ブログレイアウトパターン
```astro
---
// src/layouts/BlogLayout.astro
import BaseLayout from './BaseLayout.astro';
import Header from '../components/Header.astro';
import Footer from '../components/Footer.astro';

interface Props {
  title: string;
  description?: string;
  publishDate?: Date;
  author?: string;
  tags?: string[];
}

const { title, description, publishDate, author, tags } = Astro.props;
---

<BaseLayout title={title} description={description}>
  <Header />
  <main class="blog-main">
    <article class="blog-article">
      <header class="blog-header">
        <h1>{title}</h1>
        {publishDate && (
          <time datetime={publishDate.toISOString()}>
            {publishDate.toLocaleDateString('ja-JP')}
          </time>
        )}
        {author && <p class="author">投稿者: {author}</p>}
        {tags && (
          <div class="tags">
            {tags.map((tag) => (
              <span class="tag">#{tag}</span>
            ))}
          </div>
        )}
      </header>
      <div class="blog-content">
        <slot />
      </div>
    </article>
  </main>
  <Footer />
</BaseLayout>

<style>
  .blog-main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }
  
  .blog-header {
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .blog-header h1 {
    margin: 0 0 0.5rem 0;
    font-size: 2.5rem;
  }
  
  .blog-header time {
    color: #6b7280;
    font-size: 0.875rem;
  }
  
  .tags {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }
  
  .tag {
    padding: 0.25rem 0.5rem;
    background-color: #f3f4f6;
    border-radius: 0.25rem;
    font-size: 0.875rem;
  }
</style>
```

## データパターン

### 1. コンテンツコレクションパターン
```typescript
// src/content/config.ts
import { defineCollection, z } from 'astro:content';

const blogCollection = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    description: z.string().optional(),
    pubDate: z.coerce.date(),
    author: z.string().optional(),
    tags: z.array(z.string()).default([]),
    category: z.string().optional(),
    draft: z.boolean().default(false),
    featured: z.boolean().default(false),
  }),
});

export const collections = {
  blog: blogCollection,
};
```

### 2. APIデータ取得パターン
```astro
---
// src/components/ExternalData.astro
interface Props {
  endpoint: string;
  limit?: number;
}

const { endpoint, limit = 10 } = Astro.props;

let data = [];
let error = null;

try {
  const response = await fetch(endpoint);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  data = await response.json();
  data = data.slice(0, limit);
} catch (err) {
  error = err.message;
  console.error('データ取得失敗:', err);
}
---

<div class="data-container">
  {error ? (
    <div class="error">
      <p>データ読み込みエラー: {error}</p>
    </div>
  ) : (
    <div class="data-list">
      {data.map((item, index) => (
        <div class="data-item" key={index}>
          <h3>{item.title || item.name}</h3>
          <p>{item.description || item.summary}</p>
        </div>
      ))}
    </div>
  )}
</div>

<style>
  .data-container {
    padding: 1rem;
  }
  
  .error {
    color: #dc2626;
    background-color: #fef2f2;
    padding: 1rem;
    border-radius: 0.5rem;
  }
  
  .data-item {
    padding: 1rem;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .data-item:last-child {
    border-bottom: none;
  }
</style>
```

## フォームパターン

### 1. お問い合わせフォームパターン
```astro
---
// src/components/ContactForm.astro
---

<form id="contact-form" class="contact-form">
  <div class="form-group">
    <label for="name">名前</label>
    <input type="text" id="name" name="name" required />
  </div>
  
  <div class="form-group">
    <label for="email">メールアドレス</label>
    <input type="email" id="email" name="email" required />
  </div>
  
  <div class="form-group">
    <label for="subject">件名</label>
    <input type="text" id="subject" name="subject" required />
  </div>
  
  <div class="form-group">
    <label for="message">メッセージ</label>
    <textarea id="message" name="message" rows="5" required></textarea>
  </div>
  
  <button type="submit">メッセージを送信</button>
  
  <div id="form-status" class="form-status"></div>
</form>

<style>
  .contact-form {
    max-width: 600px;
    margin: 0 auto;
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }
  
  .form-group input,
  .form-group textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 0.375rem;
    font-size: 1rem;
  }
  
  .form-group input:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  
  button {
    background-color: #3b82f6;
    color: white;
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 0.375rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button:hover {
    background-color: #2563eb;
  }
  
  .form-status {
    margin-top: 1rem;
    padding: 0.75rem;
    border-radius: 0.375rem;
  }
  
  .form-status.success {
    background-color: #dcfce7;
    color: #166534;
  }
  
  .form-status.error {
    background-color: #fef2f2;
    color: #dc2626;
  }
</style>

<script>
  document.getElementById('contact-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const form = e.target;
    const statusEl = document.getElementById('form-status');
    const formData = new FormData(form);
    
    try {
      const response = await fetch('/api/contact', {
        method: 'POST',
        body: formData,
      });
      
      if (response.ok) {
        statusEl.textContent = 'メッセージが正常に送信されました！';
        statusEl.className = 'form-status success';
        form.reset();
      } else {
        throw new Error('メッセージの送信に失敗しました');
      }
    } catch (error) {
      statusEl.textContent = 'メッセージ送信エラー。もう一度お試しください。';
      statusEl.className = 'form-status error';
    }
  });
</script>
```

## ナビゲーションパターン

### 1. ヘッダーナビゲーションパターン
```astro
---
// src/components/Header.astro
---

<header class="site-header">
  <nav class="nav-container">
    <div class="nav-brand">
      <a href="/">ujiprog.com</a>
    </div>
    
    <ul class="nav-menu">
      <li><a href="/">ホーム</a></li>
      <li><a href="/about/">概要</a></li>
      <li><a href="/blog/">ブログ</a></li>
      <li><a href="/contact/">お問い合わせ</a></li>
    </ul>
    
    <button id="nav-toggle" class="nav-toggle" aria-label="ナビゲーションを切り替え">
      <span></span>
      <span></span>
      <span></span>
    </button>
  </nav>
</header>

<style>
  .site-header {
    background-color: white;
    border-bottom: 1px solid #e5e7eb;
    position: sticky;
    top: 0;
    z-index: 100;
  }
  
  .nav-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .nav-brand a {
    font-size: 1.5rem;
    font-weight: bold;
    text-decoration: none;
    color: #111827;
  }
  
  .nav-menu {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 2rem;
  }
  
  .nav-menu a {
    text-decoration: none;
    color: #4b5563;
    font-weight: 500;
    transition: color 0.2s;
  }
  
  .nav-menu a:hover {
    color: #111827;
  }
  
  .nav-toggle {
    display: none;
    flex-direction: column;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
  }
  
  .nav-toggle span {
    width: 25px;
    height: 3px;
    background-color: #4b5563;
    margin: 3px 0;
    transition: 0.3s;
  }
  
  @media (max-width: 768px) {
    .nav-menu {
      display: none;
      position: absolute;
      top: 100%;
      left: 0;
      right: 0;
      background-color: white;
      flex-direction: column;
      padding: 1rem;
      border-bottom: 1px solid #e5e7eb;
      gap: 1rem;
    }
    
    .nav-menu.active {
      display: flex;
    }
    
    .nav-toggle {
      display: flex;
    }
  }
</style>

<script>
  const navToggle = document.getElementById('nav-toggle');
  const navMenu = document.querySelector('.nav-menu');
  
  navToggle.addEventListener('click', () => {
    navMenu.classList.toggle('active');
  });
  
  // 外側クリックでメニューを閉じる
  document.addEventListener('click', (e) => {
    if (!navToggle.contains(e.target) && !navMenu.contains(e.target)) {
      navMenu.classList.remove('active');
    }
  });
</script>
```

## ユーティリティパターン

### 1. 日付フォーマットパターン
```astro
---
// src/components/FormattedDate.astro
interface Props {
  date: Date;
  format?: 'short' | 'long' | 'iso';
  locale?: string;
}

const { date, format = 'short', locale = 'ja-JP' } = Astro.props;

function formatDate(date: Date, format: string, locale: string): string {
  switch (format) {
    case 'short':
      return date.toLocaleDateString(locale);
    case 'long':
      return date.toLocaleDateString(locale, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      });
    case 'iso':
      return date.toISOString();
    default:
      return date.toLocaleDateString(locale);
  }
}

const formattedDate = formatDate(date, format, locale);
---

<time datetime={date.toISOString()}>{formattedDate}</time>
```

### 2. SEOパターン
```astro
---
// src/components/SEO.astro
interface Props {
  title: string;
  description?: string;
  ogImage?: string;
  ogType?: string;
  canonical?: string;
  keywords?: string[];
  author?: string;
}

const { 
  title, 
  description, 
  ogImage, 
  ogType = 'website',
  canonical,
  keywords,
  author 
} = Astro.props;

const siteUrl = Astro.site.href;
const fullCanonical = canonical ? `${siteUrl}${canonical}` : siteUrl;
---

<!-- Metaタグ -->
{description && <meta name="description" content={description} />}
{keywords && <meta name="keywords" content={keywords.join(', ')} />}
{author && <meta name="author" content={author} />}
<link rel="canonical" href={fullCanonical} />

<!-- Open Graph -->
<meta property="og:title" content={title} />
{description && <meta property="og:description" content={description} />}
<meta property="og:type" content={ogType} />
<meta property="og:url" content={fullCanonical} />
{ogImage && <meta property="og:image" content={ogImage} />}

<!-- Twitter Card -->
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content={title} />
{description && <meta name="twitter:description" content={description} />}
{ogImage && <meta name="twitter:image" content={ogImage} />}

<!-- JSON-LD構造化データ -->
<script type="application/ld+json">
{{
  "@context": "https://schema.org",
  "@type": ogType === 'article' ? "Article" : "WebPage",
  "name": title,
  "description": description,
  "url": fullCanonical,
  "image": ogImage,
  "author": author ? {
    "@type": "Person",
    "name": author
  } : undefined,
  "publisher": {
    "@type": "Organization",
    "name": "ujiprog.com",
    "url": siteUrl
  }
}}
</script>
```