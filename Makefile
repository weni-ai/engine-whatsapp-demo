.PHONY: build
build:
	docker-compose -f ./docker/docker-compose.yml build

.PHONY: up
up:
	docker-compose -f ./docker/docker-compose.yml up -d

.PHONY: down
down:
	docker-compose -f ./docker/docker-compose.yml down

.PHONY: run
run:
	go run ./cmd