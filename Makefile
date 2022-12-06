COMPOSE_FILE ?= deploy/local/docker-compose.yml
PROJECT_NAME := image-share

.PHONY: build
build:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) build --progress=plain

.PHONY: up
up:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up

.PHONY: down
down:
	docker-compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down

.PHONY: local
local: build up

YELLOW := "\e[1;33m"
GREEN := "\e[1;32m"
NC := "\e[0m"

INFO := @bash -c '\
	printf $(GREEN); \
	echo "=> $$1"; \
	printf $(NC)' VALUE