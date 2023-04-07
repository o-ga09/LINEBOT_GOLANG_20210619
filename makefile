.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:	#	build docker image to deploy
	docker build -t taiti09/cloud-run-test:${DOCKER_TAG} \
			--target deploy ./

build-gcp: #	build docker image to cloud deploy
	docker build -t gcr.io/todo-app-20221107/cloud-run-test:${DOCKER_TAG} \
			--target deploy ./

build-local:
	docker compose build --no-cache

up:
	docker compose up

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

help:
	@grep -E '[^a-zA-Z_-]+:,*?## .*$$' $(MAKEFILE_LIST) | \
			awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'