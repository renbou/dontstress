package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/renbou/dontstress/stresser/lib/compile"
)

var (
	errInternal    = errors.New("internal error")
	errCompilation = errors.New("compilation error")
)

type compilationResult struct {
	err error
	// Error message if err is not nil
	message string
	// Path to compiled executable
	path string
}

const (
	// 128 kb limit on compilation logs
	CompilationLogLimit = 128_000
	// 12 seconds for compilation should be more than enough
	Timeout = time.Second * 12
)

func compileCode(path, lang string) *compilationResult {
	compiler, ok := compile.GetCompiler(lang)
	if !ok {
		return &compilationResult{errCompilation, fmt.Sprintf("%s is not supported yet", lang), ""}
	}

	ctx, _ := context.WithTimeout(context.Background(), Timeout)
	execPath, err := compiler.Compile(&compile.Options{
		Path:        path,
		OutputLimit: CompilationLogLimit,
		Context:     ctx,
	})

	if err != nil {
		// oh boy, here comes the error handling
		ce := &compile.CompilationError{}
		if errors.As(err, &ce) {
			return &compilationResult{errCompilation, ce.Error(), ""}
		} else {
			return &compilationResult{errInternal, err.Error(), ""}
		}
	}
	return &compilationResult{nil, "", execPath}
}

type CompilationRequest struct {
	fileId string
}

func LambdaHandler(cr *CompilationRequest) error {

}

func main() {

}
