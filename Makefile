.PHONY: all clean test

all: assemble

lint:
	@echo "\nApplying golint\n"
	@golint ./...

integration-test:
	@echo "\nRunning integration tests\n"
	@go test -cover -run Integration ./...

unit-test:
	@echo "\nRunning unit tests\n"
	@go test -cover -short ./...

test: unit-test integration-test
	@echo "\nRunning tests\n"

assemble:
	@echo "\nBuilding application"
	@go build -o authorizer cmd/main.go

clean:
	@echo "\nRemove authorizer executable"
	-@rm authorizer 2>/dev/null || echo "\nExecutable file authorizer not found to remove"