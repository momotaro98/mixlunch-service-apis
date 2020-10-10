
# -----------------------------------------------------------------
#    ENV VARIABLE
# -----------------------------------------------------------------

GOCMD       := go
GOBUILD     := $(GOCMD) build
GOCLEAN     := $(GOCMD) clean
GOFMT       := $(GOCMD) fmt
GOTEST      := $(GOCMD) test
GOVET       := $(GOCMD) vet
#GOLINT      := golint -set_exit_status TODO: Apply lint in the future
GOGET       := $(GOCMD) get
GOIMPORTS   := goimports -w
GOTOOL      := $(GOCMD) tool
GOGENERATE  := $(GOCMD) generate
GOMOD       := $(GOCMD) mod

BINARY_DIR  := build
BINARY_NAME := mixlunch-service-api


TEST_INTEGRATION_PATH := ./integration-test

# -----------------------------------------------------------------
#    Main targets
# -----------------------------------------------------------------

.PHONY: all
all: generate update-proto tidy fmt test lint build ## Update generated code, Format, Run tests, and Build

.PHONY: clean
clean: ## Remove binaries
	@$(GOCLEAN)
	@rm -rf $(BINARY_DIR)

.PHONY: build
build: ## Build app
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v

.PHONY: fmt
fmt: ## Format golang codes
	@$(GOIMPORTS) `find . -type f -not -name '*_mock.go' -not -name '*.pb.go' -name '*.go' -not -path '*/testmock/*'`
	@$(GOFMT) `go list ./...`
	@strictimportsort -w -exclude "*_mock.go,*.pb.go" -exclude-dir "testmock" -local "github.com/momotaro98/mixlunch-service-api" .

.PHONY: tidy
tidy: ## Tidy up dependencies
	@$(GOMOD) tidy

.PHONY: lint
lint: ## Run linter
#@$(GOLINT) ./...
	@$(GOVET) `go list ./...`
	@strictimportsort -exclude "*_mock.go,*.pb.go" -exclude-dir "testmock" -local "github.com/momotaro98/mixlunch-service-api" .

.PHONY: test
test: ## Run all tests
	$(GOTEST) -v ./...

.PHONY: test-coverage-output
test-coverage-output: ## Run all the tests and out coverage file
	$(GOTEST) -cover ./... -coverprofile=c.out
	$(GOTOOL) cover -html=c.out

.PHONY: generate
generate: ## Apply generated code
	$(GOGENERATE) -x ./...

.PHONY: update-all-packages
update-all-package: ## Update dependency packages
	$(GOGET) -u -v -t ./...

.PHONY: get-dependencies
get-dependencies: ## Get dependency packages
	$(GOGET) -v -t -d ./...

.PHONY: update-proto
update-proto: ## Update source code for protocol buffer file
	protoc -I pb/ pb/mixlunch.proto --go_out=plugins=grpc:pb

.PHONY: help
help: ## Help command
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: hot-reload
hot-reload:
	@air

# -----------------------------------------------------------------
#    Setup targets
# -----------------------------------------------------------------

.PHONY: setup
setup: install-tools
	@true

.PHONY: install-tools
install-tools:
	$(GOMOD) download
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install % # copied from https://marcofranssen.nl/manage-go-tools-via-go-modules/

# -----------------------------------------------------------------
#    local targets
# -----------------------------------------------------------------

.PHONY: docker-run
docker-run: ## Run docker containers with apps and middlewares
	docker-compose down; docker-compose build; docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

.PHONY: docker-min-run
docker-mid-run: ## Run docker containers with only middlewares
	docker-compose down; docker-compose -f docker-compose.mid.yml up

.PHONY: docker-stop
docker-stop: ## Stop and remove docker containers
	docker-compose down

.PHONY: integration-test
integration-test: ## Run integration test, which requires running docekr containers
	@newman run -e $(TEST_INTEGRATION_PATH)/Local.postman_environment.json $(TEST_INTEGRATION_PATH)/mixlunch-service-api.postman_collection.json
	@bash ./integration-test/test_grpc.sh

.PHONY: integration-test-local
integration-test-local: ## Run integration test with running docker containers
	docker-compose down; docker-compose build; docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
	sleep 15 # Bad way but need it to wait for preparation of the containers
	@newman run -e $(TEST_INTEGRATION_PATH)/Local.postman_environment.json $(TEST_INTEGRATION_PATH)/mixlunch-service-api.postman_collection.json
	@bash ./integration-test/test_grpc.sh
	docker-compose down
