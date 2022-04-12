.DEFAULT_GOAL := all

.PHONY: all
all: lint test currency_converter

.PHONY: test
test: generate
	go test -race -timeout 10s -count 1 ./...

.PHONY: generate
generate:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	go generate ./...

.PHONY: lint
lint:
	golangci-lint --timeout=5m run

.PHONY: currency_converter
currency_converter:
	go build -o ./bin/currency_converter ./cmd
