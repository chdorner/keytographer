KEYMAP_FILE ?= keymap.example.yaml

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: build
build: # Build keytographer
	go build ./...

.PHONY: test
test: # Test keytographer
	go test -cover ./...

.PHONY: install
install: build # Install keytographer to /usr/local/bin
	cp keytographer /usr/local/bin/


.PHONY: install-tools
install-tools: # Install development tools
	brew install goreleaser/tap/goreleaser anchore/syft/syft
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: lint
lint: # Lint source code
	golangci-lint run

.PHONY: run-live
run-live: # Start live server with automatic code reload
	air -- -d live -c ${KEYMAP_FILE}

.PHONY: run-validate
run-validate: # Validates a keymap
	go run main.go -d validate -c ${KEYMAP_FILE}
