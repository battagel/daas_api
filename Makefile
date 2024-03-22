GO := go
MOCKERY := mockery
HELM := helm
K3D := k3d
K3D_CONF := k3d-conf.yaml
KUBECTL := kubectl
DOCKER := docker
MIGRATE := migrate

NAME := daas_api
VER := latest
CMD_DIR := $(CURDIR)/cmd
BIN_DIR := $(CURDIR)/bin
MIGRATIONS_DIR := $(CURDIR)/pkg/sqlite/migrations
HELM_DIR := $(CURDIR)/helm
MAIN_LOCATION := $(CMD_DIR)/$(NAME)/main.go

## help: Print this message
.PHONY: help
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'

## build: Create the binary
.PHONY: build
build: vendor
	@$(GO) build -o $(BIN_DIR)/$(NAME) -mod=vendor $(MAIN_LOCATION)


## linux: Create a binary for linux
.PHONY: deploy
deploy:
## Copy over config.env, database.go and clone the repo. Launch make run

## run: Run the binary
.PHONY: run
run: build
	@$(BIN_DIR)/$(NAME)

## vendor: Download the vendored dependencies
.PHONY: vendor
vendor:
	@$(GO) mod tidy
	@$(GO) mod vendor

## test: Run the tests
.PHONY: test
test:
	@$(GO) test -v ./... --cover

## mock: Generate the mocks for testing
.PHONY: mock
mock:
	@$(MOCKERY) --dir ./internal -r --all --config .mockery.yaml

## db: delete then recreate the sqlite db
.PHONY: db
db:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database sqlite3://database.db down
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database sqlite3://database.db up
