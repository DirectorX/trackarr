.DEFAULT_GOAL := build
CMD           := trackarr
GOARCH        := $(shell go env GOARCH)
GOOS          := $(shell go env GOOS)
TARGET        := ${GOOS}_${GOARCH}
PLATFORMS     := darwin linux windows
ARCHITECTURES := amd64
DIST_PATH     := dist
BUILD_PATH    := ${DIST_PATH}/${TARGET}
DESTDIR       := /usr/local/bin
GO_FILES      := $(shell find . -path ./vendor -prune -or -type f -name '*.go' -print)
GO_PACKAGES   := $(shell go list -mod vendor ./...)
GIT_COMMIT    := $(shell git rev-parse --short HEAD)
GIT_BRANCH    := $(shell git symbolic-ref --short HEAD)
VERSION_PATH  := VERSION
VERSION       := $(shell cat ${VERSION_PATH})
TIMESTAMP     := $(shell date +%s)
WEB_DIR       := web
RICE_FILE     := rice-box.go
WEB_RICE_FILE := ${WEB_DIR}/${RICE_FILE}
WEB_FILES     := $(shell find ${WEB_DIR} -type f -and -not -name ${RICE_FILE} -print)

.PHONY: all
all: test lint build

.PHONY: test-all
test-all: test lint

.PHONY: test
test:
	@echo "*** go test ***"
	go test -cover -mod vendor -v -race ${GO_PACKAGES}

.PHONY: lint
lint:
	@echo "*** golangci-lint ***"
	golangci-lint run

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: vendor_update
vendor_update:
	go get -u ./...
	${MAKE} vendor

.PHONY: build
build: rice ${BUILD_PATH}/${CMD}

.PHONY: rice
rice: ${WEB_RICE_FILE}

.PHONY: build-all
build-all:
	@$(foreach GOOS,$(PLATFORMS), $(foreach GOARCH,$(ARCHITECTURES), ${MAKE} build GOOS=${GOOS} GOARCH=${GOARCH};))

${BUILD_PATH}/${CMD}: ${GO_FILES} go.sum
	@echo "Building for ${TARGET}..." && \
	mkdir -p ${BUILD_PATH} && \
	CGO_ENABLED=1 go build \
		-mod vendor \
		-trimpath \
		-ldflags "-s -w -X main.buildVersion=${VERSION} -X main.buildGitCommit=${GIT_COMMIT} -X main.buildTimestamp=${TIMESTAMP}" \
		-o ${BUILD_PATH}/${CMD} \
		.

${WEB_RICE_FILE}: ${WEB_FILES} go.sum
	@echo "Generating rice..." && \
	cd ${WEB_DIR} && rice embed-go

.PHONY: install
install: build
	install -m 0755 ${BUILD_PATH}/${CMD} ${DESTDIR}/${CMD}

.PHONY: clean
clean:
	rm -f ${WEB_RICE_FILE}
	rm -rf ${DIST_PATH}
