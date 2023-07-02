PWD = $(shell pwd)

all: fmt lint test 

install:
	@brew install librdkafka
	@go install github.com/golangci/golangci-lint@latest
	
fmt:
	@echo "==> Formatting source code..."
	@go fmt ./...

lint:
	@echo "==> Running lint check..."
	@golangci-lint --config docker/.golangci.yml run ./...
	@go vet ./...

test:
	@echo "==> Running tests...$(PWD)"
	@docker exec gokit-localstack aws --no-sign-request --endpoint-url=http://localstack:4566 s3 mb s3://test
	@go clean -testcache
	@go test -tags=unit_test -vet=off `go list ./... | grep -v /lib/test` -p 1 -timeout 30s --cover -coverprofile cover.out.tmp
	@cat cover.out.tmp | grep -v ".mock.go" > cover.out && echo "Total Test Coverage: " && go tool cover -func cover.out | tail -1 | tail -c 6 && rm cover.out cover.out.tmp

test-up:
	@docker compose -f docker/compose.ci.yaml up -d
	@docker ps

test-down:
	@docker compose -f docker/compose.ci.yaml down -v --rmi local

.PHONY: fmt lint test test-up test-down