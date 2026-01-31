---
title: 令和8年 VSCode から Vim(Neovim) に戻った
display_title: 令和8年\nVSCode から Vim(Neovim) に\n戻った
published_at: 2026-01-31
---

GitHub Copilot によるコード補完が流行りだした辺り(2022,23年あたり)でメインのエディタを VSCode に移行していました。
Copilot Chat や Cursor, Cline など、VSCode のエコシステムのLLMツールの流行があったので、その恩恵を受けたかったのがモチベーションでした。

https://x.com/uji_rb/status/1643625874086936583?s=20

VSCode Neovim にはとてもお世話になりました。

令和8年、また Vim(Neovim も併用、というかほぼ Neovim 使ってます)に戻ってきました。
Claude Code をはじめ、CLI ベースのコーディングエージェントが普及したのがとても大きかったです。この流れに歓喜した Vimmer は多いんじゃないでしょうか（知らんけど）
使っているエディタを問わず、LLMを使った開発ができる環境が整備されていく中で、また Vim に戻りたいなという気持ちが沸きました。

## 入力補完を切ってみている

以前 Vim を運用していたときは、[asyncomplete.vim](https://github.com/prabirshrestha/asyncomplete.vim) などを使って入力補完のサポートもこだわっていたのですが、今は標準の機能以外は利用していないです。
コードを書くことをサポートする部分の設定やプラグインの導入は基本ストップしてみています。

自分がコードを直接ゴリゴリ書くこともどんどん減っている状況で、エディタにあまり高度な機能を求めることが少なくなってきたなと思います。
一見「自分が手で書いた方が早い」と思えるようなタスクも、

- 認識してなかった影響の検知
- 疲れの軽減
- 別作業時や外出時などの時間活用

など、コーディングエージェントに任せることによるメリットは多いです。
コードの書きやすさを良くしすぎないことで、「大人しくエージェント使おう」とさせたい気持ちです。
遊びとして自分でゴリゴリコード書きたい欲求が出てくることもあるのですが、その時は標準の機能で己を鍛えます。

一方で、**コードの読みやすさ** は変わらず重要な認識で、そこに集中して斧を研いで行きたいなと考えてます。

## 使ってる Vim プラグイン

現時点で利用しているのはこちら:

- [k-takata/minpac](https://github.com/k-takata/minpac)
- [tpope/vim-abolish](https://github.com/tpope/vim-abolish)
- [tpope/vim-surround](https://github.com/tpope/vim-surround)
- [tpope/vim-rhubarb](https://github.com/tpope/vim-rhubarb)
- [ctrlpvim/ctrlp.vim](https://github.com/ctrlpvim/ctrlp.vim)
- [prabirshrestha/vim-lsp](https://github.com/prabirshrestha/vim-lsp)
- [mattn/vim-lsp-settings](https://github.com/mattn/vim-lsp-settings)
- [easymotion/vim-easymotion](https://github.com/easymotion/vim-easymotion)
- [vim-test/vim-test](https://github.com/vim-test/vim-test)
- [lambdalisue/fern.vim](https://github.com/lambdalisue/fern.vim)
- [pgr0ss/vim-github-url](https://github.com/pgr0ss/vim-github-url)

[dotfiles](https://github.com/uji/dotfiles/blob/1fe444134251e631e2d9da60e6fa36ce3f77af4f/vimrc/vim/plug-settings/minpac.vim)

こう見ると割とミニマム寄りなんですかね。スッキリしててメンテナンスしやすいです。
どのプラグインも安定しててとても気に入ってます。

Neovim のサポートのみで Vim で利用できないプラグインは、便利そうなものが多そうだなと思いつつも利用に踏み切れておらずで、Vim script 実装のもののみのラインナップになってます。
Neovim の利用で問題が出た時にすぐに Vim を同じ設定で動かせるのはメリットなのかもしれません。（基本 Neovim 使ってて問題起きたことはないですね...）

先日 Vimmer の同僚とプラグインの見せ合いをしたんですが全然違ってて面白かったです。

## 創造意欲が沸いてくる

Vim のような、ミニマムで最初から色々整備されているわけではないエディタを触っていると、「こういう設定/プラグインあったら便利そう」といった創造意欲が沸き立つ感覚があります。
プログラムを書いていたはずが気づいたら `.vimrc` を開いていて本筋の作業が進まず困ることもあったりするのですが、そういう時間にすごく心地よさを感じます。
やっぱり、自分好みの道具づくりは楽しいですね。

そして開発環境を整えると、それを使ってまた何か作りたくなる好循環が生まれるのですごいです。

## 開発環境 2026

こんな感じになってます。

```
OS: macOS (Macbook M4 Air)
エディタ(メモ・ブログ執筆にも利用): Vim / Neovim
シェル: bash
ターミナル: Ghostty (+tmux)
コーディングエージェント: Claude Code
ランチャー: Spotlight
ブラウザ: Chrome
```

これからもアップデート楽しんでいきます。
