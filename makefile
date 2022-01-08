PKG_NAME := "adderall"
PKG_HOST := "go.adenix.dev"
PKG := "${PKG_HOST}/${PKG_NAME}"
PKG_LIST := $(shell go list ${PKG}/... 		\
	| grep -v /vendor | grep -v /vendor/ 	\
	| grep -v /example | grep -v /example/)

.PHONY: list
list:
	@for i in ${PKG_LIST}; do echo $$i; done

.PHONY: fmt
fmt:
	@go fmt ${PKG_LIST}

.PHONY: vet
vet:
	@go vet ${PKG_LIST}

.PHONY: golint
golint:
	@go run golang.org/x/lint/golint ${PKG_LIST}

.PHONY: staticcheck
staticcheck:
	@go run honnef.co/go/tools/cmd/staticcheck ${PKG_LIST}

.PHONY: errcheck
errcheck:
	@go run github.com/kisielk/errcheck ${PKG_LIST}

.PHONY: lint
lint: fmt vet golint staticcheck errcheck

.PHONY: clean
clean:
	go clean -i ./...