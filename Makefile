#!/usr/bin/make
# Makefile readme (ru): <http://linux.yaroslavl.ru/docs/prog/gnu_make_3-79_russian_manual.html>
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>

SHELL = /bin/sh

DOCKER_BIN = $(shell command -v docker 2> /dev/null)
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)"
APP_DIR = $(notdir $(CURDIR))
APP_NAME:=$(shell echo "${APP_DIR}" | tr '[:upper:]' '[:lower:]')

COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
TAG = $(shell git describe --tags 2>/dev/null | cut -d- -f1)
BUILDVERSION = $(shell dd if=/dev/urandom bs=1 count=4 2>/dev/null | hexdump -e '1/1 "%u"')

ifeq (${TAG},)
    TAG:=devel
endif

LDFLAGS = '-X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH} -X main.build=${BUILDVERSION}'

.PHONY : help build fmt lint gotest test cover shell image clean
.DEFAULT_GOAL : help
.SILENT : test shell

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build app binary file
	mkdir -p ./build
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app sh -c "go mod download && go build -a -ldflags=$(LDFLAGS) -o /build/app /src/src"

fmt: ## Run source code formatter tools
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app sh -c 'GO111MODULE=off go get golang.org/x/tools/cmd/goimports && $$GOPATH/bin/goimports -d -w .'
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app gofmt -s -w -d .

lint: ## Run app linters
	$(DOCKER_BIN) run --rm -t -v "$(shell pwd):/app" -w /app golangci/golangci-lint:v1.24-alpine golangci-lint run -v

gotest: ## Run app tests
	$(DC_BIN) run $(DC_RUN_ARGS) app go test -v -race ./...

test: lint gotest ## Run app tests and linters
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'All tests passed!';

cover: ## Run app tests with coverage report
	$(DC_BIN) run $(DC_RUN_ARGS) app sh -c 'go test -race -covermode=atomic -coverprofile /tmp/cp.out ./... && go tool cover -html=/tmp/cp.out -o ./coverage.html'
	-sensible-browser ./coverage.html && sleep 2 && rm -f ./coverage.html

shell: ## Start shell into container with golang
	$(DC_BIN) run $(DC_RUN_ARGS) app sh

image: ## Build docker image with app
	$(DOCKER_BIN) build -f ./Dockerfile -t $(APP_NAME):local .
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Now you can use image like `docker run --rm $(APP_NAME):local ...`';

clean: ## Make clean
	$(DC_BIN) down -v -t 1
	$(DOCKER_BIN) rmi $(APP_NAME):local -f
