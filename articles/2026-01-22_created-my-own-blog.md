---
title: 自分のブログを syumai/workers(Go) で作った話
published_at: 2026-01-22
---

こんにちは、ujiです。

## 構成

- syumai/workers
  - 簡単なハンドラー

goldmark: 現在のGoにおける標準的なMarkdownパーサー。CommonMark準拠で拡張性が高いです。

html/template: Go標準パッケージ。XSS対策（エスケープ処理）が強力で、安全にHTMLを生成できます。

Astro を使ってスタンダードにやるつもりが気づけば Gopher スタックになってました。

## syumai/workers

https://github.com/syumai/workers

Claude Code で簡単にできた。
https://github.com/vercel-labs/agent-browser や 公式の frontend-design Skill
個人では Pro の契約で、レートリミットが来たら他のことをするスタイル。
Web の Claude Code でコツコツ修正しつつ
