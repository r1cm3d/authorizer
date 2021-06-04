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

build:
	@echo "\nBuilding application"
	@go build -o application cmd/main.go

assemble: clean
	@echo "\nCreating Docker container"
	@docker build --tag authorizer .
	@printf "\nContainer created. Run %bmake install%b to install authorizer\n" "\e[1;33m" "\e[0m"

install:
	@echo "\nCreating executable"
	@$(shell scripts/link_script.sh)
	@echo "\nGiven execution permission to executable file"
	@chmod +x authorizer
	@printf "\nAll setup. %bauthorizer%b is ready\n %b./authorizer < EVENTS_FILE%b to use\n" "\e[1;32m" "\e[0m" "\e[1;32m" "\e[0m"

clean:
	@echo "\nRemove authorizer executable"
	-@rm authorizer 2>/dev/null || echo "\nExecutable file authorizer not found to remove"
	@echo "\nRemove authorizer image"
	-@docker rmi authorizer 2>/dev/null || echo "\nDocker image authorizer not found to remove"