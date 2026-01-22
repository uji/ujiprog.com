.PHONY: run
run:
	go tool air

.PHONY: fetch-articles
fetch-articles:
	@echo "Fetching articles from Zenn and note..."
	@( \
		curl -s "https://zenn.dev/api/articles?username=uji" | \
		jq '[.articles[] | {title: .title, url: "https://zenn.dev/uji/articles/\(.slug)", published_at: .published_at, platform: "zenn"}]'; \
		curl -s "https://note.com/api/v2/creators/ujiii/contents?kind=note&page=1" | \
		jq '[.data.contents[] | {title: .name, url: .noteUrl, published_at: .publishAt, platform: "note"}]' \
	) | jq -s 'add | sort_by(.published_at) | reverse | {articles: .}' > public/articles.json
	@echo "Updated public/articles.json"

.PHONY: deploy
deploy:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --remote
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --remote
	npm run deploy
