# Load environment variables from .env
include .env
export

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	cd test && go test -v -race -buildvcs

## run: run the cmd/api application
.PHONY: run
run:
	go run cmd/api/main.go

## build: build the cmd/api application
.PHONY: build
build:
	@go build -o ./app.bin ./cmd/api/main.go

# ==================================================================================== #
# STAGING & PRODUCTION
# ==================================================================================== #

## stage: stage api
.PHONY: stage
stage:
	@git checkout master
	@git fetch origin
	@git merge origin/master
	@docker compose --file compose.yaml --file compose.stage.yaml up --build backend --no-deps --detach backend

# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

DB_DSN="mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_SCHEMA}"

## mig/new name=$1: create a new database migration
.PHONY: mig/new
mig/new:
ifdef table
	@migrate create -dir db/migrations -ext sql ${table}
else
	@echo "must define \`table\` argument"
endif

## mig/up: apply all up database migrations
.PHONY: mig/up
mig/up:
	@migrate -path=db/migrations -database=${DB_DSN} up

## mig/down: apply all down database migrations
.PHONY: mig/down
mig/down:
	@migrate -path=db/migrations -database=${DB_DSN} down

## mig/goto version=$1: migrate to a specific version number
.PHONY: mig/goto
mig/goto:
	@migrate -path=db/migrations -database=${DB_DSN} goto ${version}

## mig/force version=$1: force database migration
.PHONY: mig/force
mig/force:
	@migrate -path=db/migrations -database=${DB_DSN} force ${version}
.PHONY: mig/version

## mig/version: print the current in-use migration version
mig/version:
	@migrate -path=db/migrations -database=${DB_DSN} version

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
