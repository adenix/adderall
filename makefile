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

GOCOVERCOBERTURA = ${GOBIN}/gocover-cobertura
MOCKGEN = ${GOBIN}/mockgen
ERRCHECK = ${GOBIN}/errcheck
GOLINT = ${GOBIN}/golint
GOTESTSUM = ${GOBIN}/gotestsum
STATICCHECK = ${GOBIN}/staticcheck

${GOCOVERCOBERTURA}: tools/go.mod
	cd tools && go install github.com/boumenot/gocover-cobertura

${MOCKGEN}: tools/go.mod
	cd tools && go install github.com/golang/mock/mockgen

${ERRCHECK}: tools/go.mod
	cd tools && go install github.com/kisielk/errcheck

${GOLINT}: tools/go.mod
	cd tools && go install golang.org/x/lint/golint

${GOTESTSUM}: tools/go.mod
	cd tools && go install gotest.tools/gotestsum

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
tidy:
	$(foreach dir,$(MODULES),(cd $(dir) && go mod tidy) &&) true

.PHONY: generate
generate: ${MOCKGEN} clean-mock
	@go generate

.PHONY: clean-test
clean-test:
	@rm -rf test

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
clean: clean-test clean-cover clean-bench clean-mock
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
test: clean-test clean-mock generate
	@go test -cover ${PKG_LIST}

.PHONY: test-report
test-report: ${GOTESTSUM} clean-test clean-mock generate
	@mkdir -p test
	@${GOTESTSUM} --junitfile test/report.xml --format testname

.PHONY: cover
cover: clean-cover clean-mock generate
	@mkdir -p cover
	@echo 'mode: count' > cover/coverage.out
	@echo ${PKG_LIST} | xargs -n1 -I{} sh -c 'go test -covermode=count -coverprofile=cover/coverage.tmp {} && tail -n +2 cover/coverage.tmp >> cover/coverage.out' && rm cover/coverage.tmp
	@go tool cover -func=cover/coverage.out

.PHONY: cover-report
cover-report: ${GOCOVERCOBERTURA} cover
	@${GOCOVERCOBERTURA} < cover/coverage.out > cover/report.xml

.PHONY: cover-html
cover-html: cover
	@go tool cover -html=cover/coverage.out

.PHONY: race
race:
	@go test -race -short ${PKG_LIST}

.PHONY: msan
msan:
	@go test -msan -short ${PKG_LIST}