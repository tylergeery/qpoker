.PHONY: help 
.DEFAULT_GOAL := help

prod_image_http=tylergeery/qcards-http
prod_image_ws=tylergeery/qcards-ws
prod_image_migrator=tylergeery/qcards-migrator
prod_image_nginx=tylergeery/qcards-nginx

build-prod-images: ## Build latest docker images
	docker build -t $(prod_image_http) --target http_prod -f ./go/Dockerfile ..
	docker build -t $(prod_image_ws) --target ws_prod -f ./go/Dockerfile ..
	docker build -t $(prod_image_migrator) --target migrator -f ./go/Dockerfile ..
	docker build -t $(prod_image_nginx) --target prod -f ./nginx/Dockerfile ..

push-prod-images: ## Push latest docker images
	docker push $(prod_image_http)
	docker push $(prod_image_ws)
	docker push $(prod_image_migrator)
	docker push $(prod_image_nginx)

prod-images: build-prod-images push-prod-images ## Build and push new prod images

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'