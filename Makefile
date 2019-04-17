
.PHONY: db
db: ## db
	 docker run -p 5432:5432 -e POSTGRES_DB=voteo -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -v ${HOME}/voteo-data:/var/lib/postgresql/data postgres:10-alpine

.PHONY: cleandb
cleandb: ## clean out local db directory
	rm -rf ${HOME}/voteo-data

.PHONY: test
test: ## test
	go test -v

.PHONY: run
run: ## run
	go build
	./voteo

.PHONY: help
help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
