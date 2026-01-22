.PHONY: run
run:
	go tool air

.PHONY: deploy
deploy:
	npx wrangler r2 object put ujiprog-static/index.html --file=index.html --remote
	npx wrangler r2 object put ujiprog-static/avator.jpg --file=public/avator.jpg --remote
	npm run deploy
