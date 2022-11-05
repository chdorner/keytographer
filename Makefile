KEYMAP_FILE ?= keymap.example.yaml

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done


.PHONY: install-tools
install-tools: # Install development tools
	go install github.com/cosmtrek/air@latest

.PHONY: live
live: # Start live server with automatic code reload
	air -- -w ${KEYMAP_FILE} -d