# Include variables from .envrc files
-include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo "\t##IMPORTANT##: please run 'echo .envrc >> .gitignore' at very first time"
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/main: run the cmd/main binary file
.PHONY: run/main
run/main:
	go run ./cmd

## run/debug: debug app use dlv
.PHONY: run/debug
run/debug:
	dlv debug ./cmd --headless --listen :4040

## run/test: runs go test with default values
.PHONY: run/test
run/test:
	go test -timeout 300s -v -count=1 -race ./...

## run/update: runs go get -u && go mod tidy
.PHONY: run/update
run/update:
	go get -u ./...
	go mod tidy

## db/psql: connection to the database using psql
.PHONY: db/psql
db/psql:
	@psql ${PG_DSN}

## db/generate use sqlc generated models and queries
.PHONY: db/generate
db/generate:
	@echo 'sqlc generate in internal/sqlc fold'
	@cd internal/sqlc && sqlc generate && cd ../..

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migrate files for ${name}'
	@migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@migrate -path ./migrations -database ${PG_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	#staticcheck ./...  # go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

#current_time = $(shell date --iso-8601=seconds)
current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/main: audit
	@echo 'Building cmd/...'
	go build -ldflags=${linker_flags} -o=./bin/cmd ./cmd
	#go tool dist list
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/cmd ./cmd

## build/dlv-debug: build the application with dlv gcflags
.PHONY: build/dlv-debug
build/dlv-debug: 
	@echo "Building for delve debug..."
	@go build \
	-ldflags ${linker_flags} \
	-ldflags=-compressdwarf=false \
	-gcflags=all=-d=checkptr \
	-gcflags="all=-N -l" \
	-o ./bin/debug ./cmd
