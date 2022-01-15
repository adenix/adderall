//go:build tools
// +build tools

package tools

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/kisielk/errcheck"
	_ "golang.org/x/lint/golint"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
