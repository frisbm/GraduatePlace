.PHONY: help
# help:
#    Print this help message
help:
	@grep -o '^\#.*' Makefile | cut -d" " -f2-

.PHONY: install
# install:
#    Install required dependencies for the project
install:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest; \
	go install golang.org/x/tools/cmd/goimports@latest; \
	go mod tidy; \
	  go mod download; \
	  go mod verify

.PHONY: start
# start:
#    Start up all the services with docker compose
start:
	docker compose up

.PHONY: stop
# stop:
#    Stop all the services with docker compose
stop:
	docker compose stop

.PHONY: generate
# generate:
#    Generate sqlc
generate:
	sqlc generate
	$(MAKE) fmt

.PHONY: reset
# reset:
#    Reset data containers and their volumes for clean instance
reset:
	docker compose rm s3 postgresql -f || true
	rm -rf ./tmp/s3/

.PHONY: fmt
# fmt:
#    Format go code
fmt:
	goimports -local github.com/frisbm -w ./
