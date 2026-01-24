.PHONY: run
run:
	go tool air

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm .

.PHONY: dev
dev:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --local
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --local
	npx wrangler r2 object put ujiprog-static/favicon.ico --file=public/favicon.ico --local
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --local
	@for f in .generated/articles/*; do \
		if [ -f "$$f" ]; then \
			echo "Uploading: $$f"; \
			npx wrangler r2 object put "ujiprog-static/articles/$$(basename $$f)" --file="$$f" --local; \
		fi; \
	done
	npx wrangler dev

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
		-og-template=templates/blog-ogp-tmpl.png \
		-ascii-font=fonts/DMSans-Bold.ttf \
		-japanese-font=fonts/NotoSansJP-Bold.ttf \
		-font-size=56 \
		-articles-json=public/articles.json
	@echo "Article generation complete"

.PHONY: deploy
deploy:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --remote
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --remote
	@for file in .generated/articles/*.html .generated/articles/*.png; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file"); \
			echo "Uploading articles/$$filename..."; \
			npx wrangler r2 object put "ujiprog-static/articles/$$filename" --file="$$file" --remote; \
		fi \
	done
	$(MAKE) build
	npx wrangler deploy
