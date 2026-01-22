.PHONY: run
run:
	go tool air

.PHONY: fetch-zenn
fetch-zenn:
	@curl -s "https://zenn.dev/api/articles?username=uji" | \
	jq '{zenn: [.articles[] | {title: .title, url: "https://zenn.dev/uji/articles/\(.slug)", published_at: .published_at}]}' > public/articles.json
	@echo "Updated public/articles.json"

.PHONY: deploy
deploy:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --remote
	npx wrangler r2 object put ujiprog-static/articles.json --file=public/articles.json --remote
	npm run deploy
