PKG_NAME := "adderall"
PKG_HOST := "go.adenix.dev"
PKG := "${PKG_HOST}/${PKG_NAME}"
PKG_LIST := $(shell go list ${PKG}/... 		\
	| grep -v /vendor | grep -v /vendor/ 	\
	| grep -v /example | grep -v /example/  \
	| grep -v /mock/)

MODULES = . ./tools


### Tools Binaries ###

export GOBIN ?= $(shell pwd)/bin

MOCKGEN = ${GOBIN}/mockgen
ERRCHECK = ${GOBIN}/errcheck
GOLINT = ${GOBIN}/golint
STATICCHECK = ${GOBIN}/staticcheck

${MOCKGEN}: tools/go.mod
	cd tools && go install github.com/golang/mock/mockgen

${ERRCHECK}: tools/go.mod
	cd tools && go install github.com/kisielk/errcheck

${GOLINT}: tools/go.mod
	cd tools && go install golang.org/x/lint/golint

${STATICCHECK}: tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck


### Commands ###

.PHONY: list
list:
	@for i in ${PKG_LIST}; do echo $$i; done

.PHONY: install
install:
	$(foreach dir,$(MODULES),( cd $(dir) && go mod download) && ) true

.PHONY: tidy
tidy: generate
	$(foreach dir,$(MODULES),(cd $(dir) && go mod tidy) &&) true

.PHONY: generate
generate: ${MOCKGEN} clean-mock
	@go generate

.PHONY: clean-cover
clean-cover:
	@rm -rf cover

.PHONY: clean-bench
clean-bench:
	@rm -rf bench

.PHONY: clean-mock
clean-mock:
	@rm -rf internal/mock

.PHONY: clean
clean: clean-cover clean-bench clean-mock
	@rm -rf bin/

.PHONY: fmt
fmt:
	@go fmt ${PKG_LIST}

.PHONY: vet
vet: generate
	@go vet ${PKG_LIST}

.PHONY: golint
golint: ${GOLINT}
	@${GOLINT} ${PKG_LIST}

.PHONY: staticcheck
staticcheck: ${STATICCHECK} generate
	@${STATICCHECK} ${PKG_LIST}

.PHONY: errcheck
errcheck: ${ERRCHECK} generate
	@${ERRCHECK} ${PKG_LIST}

.PHONY: lint
lint: fmt vet golint staticcheck errcheck

.PHONY: test
test: clean-mock generate
	@go test -cover ${PKG_LIST}

.PHONY: cover
cover: clean-cover clean-mock generate
	@mkdir -p cover
	@go test -race -covermode=atomic ${PKG_LIST}

.PHONY: cover-html
cover-html: cover
	@go tool cover -html=cover/coverage.txt

.PHONY: race
race:
	@go test -race -short ${PKG_LIST}

.PHONY: msan
msan:
	@go test -msan -short ${PKG_LIST}