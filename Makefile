.PHONY: run
run:
	go tool air

.PHONY: deploy
deploy:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npm run deploy
