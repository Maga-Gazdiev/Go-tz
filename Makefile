.PHONY: help run migrate test vet lint format check up down

help:
	@echo "run      - start the API locally"
	@echo "migrate  - apply Goose migrations"
	@echo "test     - run tests"
	@echo "vet      - run go vet"
	@echo "lint     - run golangci-lint"
	@echo "format   - format Go source files"
	@echo "check    - run test, vet and lint"
	@echo "up       - start Docker Compose"
	@echo "down     - stop Docker Compose"

run:
	go run ./cmd/api

migrate:
	go run ./cmd/migrate

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

format:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

check: test vet lint

up:
	docker compose up --build

down:
	docker compose down
