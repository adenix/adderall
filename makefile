PKG_NAME := "adderall"
PKG_HOST := "go.adenix.dev"
PKG := "${PKG_HOST}/${PKG_NAME}"
PKG_LIST := $(shell go list ${PKG}/... 		\
	| grep -v /vendor | grep -v /vendor/ 	\
	| grep -v /example | grep -v /example/  \
	| grep -v /mock/)

.PHONY: list
list:
	@for i in ${PKG_LIST}; do echo $$i; done

.PHONY: generate
generate:
	@go generate

.PHONY: clean-test
clean-test:
	@rm -rf test

.PHONY: clean-cover
clean-cover:
	@rm -rf cover

.PHONY: clean-mock
clean-mock:
	@rm -rf internal/mock

.PHONY: clean
clean: clean-test clean-cover clean-mock

.PHONY: fmt
fmt:
	@go fmt ${PKG_LIST}

.PHONY: vet
vet: generate
	@go vet ${PKG_LIST}

.PHONY: golint
golint:
	@go run golang.org/x/lint/golint ${PKG_LIST}

.PHONY: staticcheck
staticcheck: generate
	@go run honnef.co/go/tools/cmd/staticcheck ${PKG_LIST}

.PHONY: errcheck
errcheck: generate
	@go run github.com/kisielk/errcheck ${PKG_LIST}

.PHONY: lint
lint: fmt vet golint staticcheck errcheck

.PHONY: test
test: clean-test clean-mock generate
	@go test -cover ${PKG_LIST}

.PHONY: test-report
test-report: clean-test clean-mock generate
	@mkdir -p test
	@go run gotest.tools/gotestsum --junitfile test/report.xml --format testname

.PHONY: cover
cover: clean-cover clean-mock generate
	@mkdir -p cover
	@echo 'mode: count' > cover/coverage.out
	@echo ${PKG_LIST} | xargs -n1 -I{} sh -c 'go test -covermode=count -coverprofile=cover/coverage.tmp {} && tail -n +2 cover/coverage.tmp >> cover/coverage.out' && rm cover/coverage.tmp
	@go tool cover -func=cover/coverage.out

.PHONY: cover-report
cover-report: cover ## Generate global code coverage report in cobertura
	@go run github.com/boumenot/gocover-cobertura < cover/coverage.out > cover/report.xml

.PHONY: cover-html
cover-html: cover ## Generate global code coverage report in HTML
	@go tool cover -html=cover/coverage.out

.PHONY: race
race:
	@go test -race -short ${PKG_LIST}

.PHONY: msan
msan:
	@go test -msan -short ${PKG_LIST}