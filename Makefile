VERSION := 0.0.3
# Name of this service/application
SERVICE_NAME := source-ranges-controller

# Docker image name for this project
IMAGE_NAME := jeffersongirao/$(SERVICE_NAME)
REPOSITORY := quay.io/$(IMAGE_NAME)

# Shell to use for running scripts
SHELL := $(shell which bash)

# Commit hash from git
COMMIT=$(shell git rev-parse HEAD)

# Branch from git
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Get the unix user id for the user running make
UID := $(shell id -u)

DEV_DIR := ./docker/development
APP_DIR := ./docker/app
WORKDIR := /go/src/github.com/jeffersongirao/source-ranges-controller

# Get docker path or an empty string
DOCKER := $(shell command -v docker)

UNIT_TEST_CMD := GOCACHE=off go test `go list ./... | grep -v /vendor/` -v

.PHONY: deps-development
# Test if the dependencies we need to run this Makefile are installed
deps-development:
ifndef DOCKER
	@echo "Docker is not available. Please install docker"
	@exit 1
endif

# Run the development environment in non-daemonized mode (foreground)
.PHONY: docker-build
docker-build: deps-development
	docker build \
	-t $(REPOSITORY)-dev:latest \
	-t $(REPOSITORY)-dev:$(COMMIT) \
	-f $(DEV_DIR)/Dockerfile \
	.

# Build source-ranges-controller executable file
.PHONY: build
build: docker-build
	docker run -ti --rm -v $(PWD):$(WORKDIR) -u $(UID):$(UID) --name $(SERVICE_NAME) $(REPOSITORY)-dev ./scripts/build.sh

# Test stuff in dev
.PHONY: unit-test
unit-test: docker-build
	docker run -ti --rm -v $(PWD):$(WORKDIR) -u $(UID):$(UID) --name $(SERVICE_NAME) $(REPOSITORY)-dev /bin/sh -c '$(UNIT_TEST_CMD)'

.PHONY: test
test: unit-test

.PHONY: get-deps
get-deps: docker-build
	docker run -ti --rm -v $(PWD):$(WORKDIR) -u $(UID):$(UID) --name $(SERVICE_NAME) $(REPOSITORY)-dev /bin/sh -c '$(GET_DEPS_CMD)'

.PHONY: update-deps
update-deps: docker-build
	docker run -ti --rm -v $(PWD):$(WORKDIR) -u $(UID):$(UID) --name $(SERVICE_NAME) $(REPOSITORY)-dev /bin/sh -c '$(UPDATE_DEPS_CMD)'

# Build the production image based on the public one
.PHONY: image
image: deps-development
	docker build \
	-t $(SERVICE_NAME) \
	-t $(REPOSITORY):latest \
	-t $(REPOSITORY):$(COMMIT) \
	-t $(REPOSITORY):$(BRANCH) \
	-f $(APP_DIR)/Dockerfile \
	.

.PHONY: tag
tag:
	git tag $(VERSION)

.PHONY: publish
publish:
	@COMMIT_VERSION="$$(git rev-list -n 1 $(VERSION))"; \
	docker tag $(REPOSITORY):"$$COMMIT_VERSION" $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):$(VERSION)
	docker push $(REPOSITORY):latest

.PHONY: release
release: tag image publish
