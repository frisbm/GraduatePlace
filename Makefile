.PHONY: help
help: # Print this help message
	@grep '^[a-z]' Makefile

.PHONY: install
install: # Install required dependencies for the project
	go mod tidy; \
	  go mod download; \
	  go mod verify

.PHONY: start
start: # Run server and database with docker compose with hot reloading using air
	docker compose up

.PHONY: stop
stop: # Stop running docker compose server and database
	docker compose stop

.PHONY: generate
generate: # Generate sqlc
	sqlc generate
	$(MAKE) fmt

.PHONY: reset
reset: # Reset data containers and their volumes for clean instance
	docker compose rm s3 postgresql -f || true
	rm -rf ./tmp/s3/

.PHONY: fmt
fmt: # Format go code
	goimports -local github.com/MatthewFrisby -w ./
