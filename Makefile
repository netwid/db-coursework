.PHONY: build
build:
	go build

.PHONY: deps
deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrate
migrate:
	goose -dir ./migrations postgres "postgres://dev:dev@localhost:5432/postgres?sslmode=disable" up
