.PHONY: help up down clean dev- db-migrate db-apply db-revert dev-user sass js
.DEFAULT_GOAL := help

ts := $(shell date +%Y%m%d-%H%M%S)


up: dev-setup db-apply db-seed  ## Create local docker env

dev-setup:
	docker-compose up -d nginx

down:  ## Tear down local docker env
	docker-compose down --remove-orphans

clean:  ## Remove all seed db data
	sudo rm -rf db/data

db-migrate: check-MSG ## Create a new DB migration
	$(eval $@_MSG := $(shell echo ${MSG} | sed -E 's/[^a-zA-Z0-9]+/_/g' | sed -E 's/^_+|_+$$//g' | tr A-Z a-z))
	docker-compose run -e MSG="$($@_MSG)" --rm migrator ./migrate.sh migrate
	sudo chown -R $(shell whoami):$(shell whoami) $(shell pwd)/db/migrations/

db-apply:  ## Apply DB migrations
	docker-compose run --rm migrator ./migrate.sh apply

db-revert:  ## Revert latest DB migration
	docker-compose run --rm migrator ./migrate.sh revert

db-seed:
	docker-compose run --rm migrator ./migrate.sh seed

dev-player:  ## Create a local player
	curl -XPOST -H 'Content-Type: application/json' -d '{"username": "player-$(ts)", "email": "player-$(ts)@test.com", "pw": "testpass"}' 'http://localhost:8080/api/v1/players' | jq

test:  ## Run tests in local dev env
	@docker-compose exec app go test -count=1 -timeout 5s ./...

js:  ## Build client js
	@cd client && /usr/bin/npx webpack

js-watch:  ## Build client js (watching for changes)
	@cd client && NODE_ENV=dev /usr/bin/npx webpack --watch

sass:
	sass --no-source-map client/sass:client/assets/css

sass-watch:
	sass --no-source-map --watch client/sass/entry:client/assets/css

check-%:
	@if [ -z '${${*}}' ]; then echo 'Environment variable $* not set' && exit 1; fi

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'