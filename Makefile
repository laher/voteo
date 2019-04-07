
.PHONY: db
db: ## db
	docker run -p 8000:8000 -v ${HOME}/voteo-data:/data amazon/dynamodb-local -jar DynamoDBLocal.jar -dbPath /data

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