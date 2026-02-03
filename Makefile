.PHONY: run
run:
	@# 両方のAir監視を並列起動
	go tool air -c .air.toml & go tool air -c .air-articles.toml & wait

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm .

# 初回セットアップ（全アセットアップロード）
.PHONY: dev-init
dev-init: generate-articles
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --local
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --local
	npx wrangler r2 object put ujiprog-static/favicon.ico --file=public/favicon.ico --local
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --local
	npx wrangler r2 object put ujiprog-static/style.css --file=public/style.css --local
	npx wrangler r2 object put ujiprog-static/article.css --file=public/article.css --local
	npx wrangler r2 object put ujiprog-static/main.js --file=public/main.js --local
	npx wrangler r2 object put ujiprog-static/article.js --file=public/article.js --local
	@# Upload article HTML files
	@for f in .generated/articles/*.html; do \
		if [ -f "$$f" ]; then \
			echo "Uploading: $$f"; \
			npx wrangler r2 object put "ujiprog-static/articles/$$(basename $$f)" --file="$$f" --local; \
		fi; \
	done
	@# Upload OG metadata and assets for dynamic OG image generation
	npx wrangler r2 object put ujiprog-static/og-meta.json --file=.generated/og-meta.json --local
	npx wrangler r2 object put ujiprog-static/fonts/DMSans-Bold.ttf --file=fonts/DMSans/DMSans-Bold.ttf --local
	npx wrangler r2 object put ujiprog-static/fonts/NotoSansJP-Bold.ttf --file=fonts/NotoSansJP/NotoSansJP-Bold.ttf --local
	npx wrangler r2 object put ujiprog-static/templates/blog-ogp-tmpl.png --file=templates/blog-ogp-tmpl.png --local

# 開発サーバー起動のみ（高速）
.PHONY: dev
dev:
	npx wrangler dev

# Git差分のMarkdownのみ処理
.PHONY: dev-sync-articles
dev-sync-articles:
	@changed=$$(git diff --name-only articles/ 2>/dev/null); \
	untracked=$$(git ls-files --others --exclude-standard articles/ 2>/dev/null); \
	all_changed="$$changed $$untracked"; \
	for md in $$all_changed; do \
		if [ -f "$$md" ] && [ "$${md##*.}" = "md" ]; then \
			slug=$$(basename "$$md" .md); \
			echo "Processing: $$md"; \
			go run ./cmd/generate \
				-single="$$md" \
				-output=.generated/articles \
				-template=templates/article.html \
				-articles-json=public/articles.json \
				-og-meta=.generated/og-meta.json; \
			npx wrangler r2 object put "ujiprog-static/articles/$$slug.html" \
				--file=".generated/articles/$$slug.html" --local; \
		fi; \
	done; \
	npx wrangler r2 object put ujiprog-static/og-meta.json --file=.generated/og-meta.json --local

.PHONY: fetch-articles
fetch-articles:
	@echo "Fetching articles from Zenn, note, and SpeakerDeck..."
	@( \
		curl -s "https://zenn.dev/api/articles?username=uji" | \
		jq '[.articles[] | {title: .title, url: "https://zenn.dev/uji/articles/\(.slug)", published_at: .published_at, platform: "zenn"}]'; \
		curl -s "https://note.com/api/v2/creators/ujiii/contents?kind=note&page=1" | \
		jq '[.data.contents[] | {title: .name, url: .noteUrl, published_at: .publishAt, platform: "note"}]'; \
		python3 -c 'import json,urllib.request,xml.etree.ElementTree as ET; \
		r=urllib.request.urlopen("https://speakerdeck.com/uji.atom"); \
		root=ET.fromstring(r.read().decode("utf-8")); \
		ns={"a":"http://www.w3.org/2005/Atom"}; \
		print(json.dumps([{"title":e.find("a:title",ns).text,"url":e.find("a:link[@rel=\"alternate\"]",ns).get("href"),"published_at":e.find("a:published",ns).text,"platform":"speakerdeck"} for e in root.findall("a:entry",ns)],ensure_ascii=False))' \
	) | jq -s 'add | sort_by(.published_at) | reverse | {articles: .}' > public/articles.json
	@echo "Updated public/articles.json"

.PHONY: generate-articles
generate-articles:
	@echo "Generating articles from markdown..."
	mkdir -p .generated/articles
	go run ./cmd/generate \
		-articles=articles \
		-output=.generated/articles \
		-template=templates/article.html \
		-articles-json=public/articles.json \
		-og-meta=.generated/og-meta.json
	@echo "Article generation complete"

.PHONY: deploy
deploy: generate-articles
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --remote
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --remote
	npx wrangler r2 object put ujiprog-static/style.css --file=public/style.css --remote
	npx wrangler r2 object put ujiprog-static/article.css --file=public/article.css --remote
	npx wrangler r2 object put ujiprog-static/main.js --file=public/main.js --remote
	npx wrangler r2 object put ujiprog-static/article.js --file=public/article.js --remote
	@# Upload article HTML files (OG images are generated dynamically)
	@for file in .generated/articles/*.html; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file"); \
			echo "Uploading articles/$$filename..."; \
			npx wrangler r2 object put "ujiprog-static/articles/$$filename" --file="$$file" --remote; \
		fi \
	done
	@# Upload OG metadata and assets for dynamic OG image generation
	npx wrangler r2 object put ujiprog-static/og-meta.json --file=.generated/og-meta.json --remote
	npx wrangler r2 object put ujiprog-static/fonts/DMSans-Bold.ttf --file=fonts/DMSans/DMSans-Bold.ttf --remote
	npx wrangler r2 object put ujiprog-static/fonts/NotoSansJP-Bold.ttf --file=fonts/NotoSansJP/NotoSansJP-Bold.ttf --remote
	npx wrangler r2 object put ujiprog-static/templates/blog-ogp-tmpl.png --file=templates/blog-ogp-tmpl.png --remote
	$(MAKE) build
	npx wrangler deploy
