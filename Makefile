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

.PHONY: run-dev
run-dev:
	gnome-terminal -- bash -c "ngrok http 8000";
	gnome-terminal --tab -- bash -c "source env.sh && go run ./cmd"

.PHONY: test
	go test -p=1 test ./...