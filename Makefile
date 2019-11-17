.DEFAULT_GOAL  := build
CMD            := trackarr
GOARCH         := $(shell go env GOARCH)
GOOS           := $(shell go env GOOS)
TARGET         := ${GOOS}_${GOARCH}
PLATFORMS      := darwin linux windows
ARCHITECTURES  := amd64
DIST_PATH      := dist
BUILD_PATH     := ${DIST_PATH}/${CMD}_${TARGET}
DESTDIR        := /usr/local/bin
GO_FILES       := $(shell find . -path ./vendor -prune -or -type f -name '*.go' -print)
GO_PACKAGES    := $(shell go list -mod vendor ./...)
GIT_COMMIT     := $(shell git rev-parse --short HEAD)
GIT_BRANCH     := $(shell git symbolic-ref --short HEAD)
TIMESTAMP      := $(shell date +%s)
WEB_DIR        := web
RICE_FILE      := rice-box.go
WEB_RICE_FILE  := ${WEB_DIR}/${RICE_FILE}
WEB_UI_PATH    := ${WEB_DIR}/trackarr-ui
WEB_UI_MODULES := ${WEB_UI_PATH}/node_modules
WEB_UI_DIST    := ${WEB_UI_PATH}/${DIST_PATH}
WEB_GO_FILES   := $(shell find ${WEB_DIR} -type f -and -name '*.go' -and -not -name ${RICE_FILE} -print)
WEB_UI_FILES   := $(shell find ${WEB_UI_PATH}/src ${WEB_UI_PATH}/public -type f -print)

# Deps
.PHONY: check_rice
check_rice:
	@command -v rice >/dev/null || (echo "rice is required."; exit 1)
.PHONY: check_golangci
check_golangci:
	@command -v golangci-lint >/dev/null || (echo "golangci-lint is required."; exit 1)
.PHONY: check_goreleaser
check_goreleaser:
	@command -v goreleaser >/dev/null || (echo "goreleaser is required."; exit 1)
.PHONY: check_yarn
check_yarn:
	@command -v yarn >/dev/null || (echo "yarn is required."; exit 1)

.PHONY: all ## Run tests, linting and build
all: test lint build

.PHONY: test-all ## Run tests and linting
test-all: test lint

.PHONY: test
test: ## Run tests
	@echo "*** go test ***"
	go test -cover -mod vendor -v -race ${GO_PACKAGES}

.PHONY: lint
lint: check_golangci ## Run linting
	@echo "*** golangci-lint ***"
	golangci-lint run

.PHONY: vendor
vendor: ## Vendor files and tidy go.mod
	go mod vendor
	go mod tidy

.PHONY: vendor_update
vendor_update: ## Update vendor dependencies
	go get -u ./...
	${MAKE} vendor

.PHONY: build
build: fetch rice web ${BUILD_PATH}/${CMD} ## Build application

.PHONY: rice
rice: check_rice ${WEB_RICE_FILE} ## Generate embedded web files

.PHONY: web
web: check_yarn ${WEB_UI_MODULES} ${WEB_UI_DIST} ## Package web frontend

# Binary
${BUILD_PATH}/${CMD}: ${GO_FILES} go.sum
	@echo "Building for ${TARGET}..." && \
	mkdir -p ${BUILD_PATH} && \
	CGO_ENABLED=0 go build \
		-mod vendor \
		-trimpath \
		-ldflags "-s -w -X main.buildVersion=0.0.0-dev -X main.buildGitCommit=${GIT_COMMIT} -X main.buildTimestamp=${TIMESTAMP}" \
		-o ${BUILD_PATH}/${CMD} \
		.

# Go generated embed files
${WEB_RICE_FILE}: ${WEB_GO_FILES} ${WEB_UI_DIST} go.sum
	@echo "Generating rice..." && \
	cd ${WEB_DIR} && rice embed-go

# Web files build
${WEB_UI_DIST}: SKIP_WEB=false
${WEB_UI_DIST}: ${WEB_UI_FILES} ${WEB_UI_PATH}/vue.config.js ${WEB_UI_PATH}/yarn.lock
	@[ "${SKIP_WEB}" = "true" ] || \
	(echo "Building Web UI..." && \
	cd ${WEB_UI_PATH} && yarn build)

# Web node modules
${WEB_UI_MODULES}: ${WEB_UI_PATH}/package.json ${WEB_UI_PATH}/yarn.lock
	@echo "Fetching node_modules..." && \
	cd ${WEB_UI_PATH} && yarn install

.PHONY: install
install: build ## Install binary
	install -m 0755 ${BUILD_PATH}/${CMD} ${DESTDIR}/${CMD}

.PHONY: clean
clean: ## Cleanup
	rm -f ${WEB_RICE_FILE}
	rm -rf ${DIST_PATH}
	rm -rf ${WEB_UI_DIST}

.PHONY: fetch
fetch: ## Fetch vendor files
	go mod vendor

.PHONY: release
release: check_goreleaser fetch rice ## Generate a release, but don't publish
	goreleaser --skip-validate --skip-publish --rm-dist

.PHONY: publish
publish: check_goreleaser fetch rice ## Generate a release, and publish
	goreleaser --rm-dist

.PHONY: snapshot
snapshot: check_goreleaser fetch rice ## Generate a snapshot release
	goreleaser --snapshot --skip-validate --skip-publish --rm-dist

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
