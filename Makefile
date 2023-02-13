.PHONY: help
help: # Print this help message
	@grep '^[a-z]' Makefile

.PHONY: install
install: # Install required dependencies for the project
	go mod tidy; \
	  go mod download; \
	  go mod verify
	cd frontend; \
	  npm ci; \
	  npm run build

.PHONY: start
start: # Run server and database with docker compose with hot reloading using air
	docker compose up

.PHONY: stop
stop: # Stop running docker compose server and database
	docker compose stop

.PHONY: generate
generate: # Generate the ent types required for this project
	cd ent; \
	  go run -mod=mod entgo.io/ent/cmd/ent generate --feature privacy ./schema
	$(MAKE) fmt

.PHONY: fmt
fmt: # Format go code
	goimports -local github.com/MatthewFrisby -w ./
