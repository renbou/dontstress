package compile

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/renbou/dontstress/stresser/lib/fsutil"
	"github.com/renbou/dontstress/stresser/lib/safeexec"
)

type Options struct {
	// Path to executable
	Path string
	// Limit of stderr
	OutputLimit int
	// Context under which all execution takes place,
	// so that we can quit instantly if necessary
	Context context.Context
}

type LangCompiler interface {
	// Returns path to compiler binary and error
	Compile(*Options) (string, error)
}

type StepCompiler struct {
	steps     []string
	extension string
}

// Config to pass as template parameters to each step
type templateConfig struct {
	Out  string
	File string
}

func (sc *StepCompiler) Compile(options *Options) (string, error) {
	var binaryPath string

	err := fsutil.WithTempDir(func(dir string) error {
		codeFileName := fsutil.RandomFileName() + "." + sc.extension
		if err := fsutil.CopyFile(options.Path, codeFileName); err != nil {
			return err
		}

		execFileName := fsutil.RandomFileName()
		execConfig := &safeexec.Config{
			Context:     options.Context,
			OutputLimit: options.OutputLimit,
			Data: &templateConfig{
				Out:  fsutil.RandomFileName(),
				File: execFileName,
			},
		}

		// Run all compilation steps, if they fail
		// then we have a compilation error (or an internal error)
		for _, cmd := range sc.steps {
			proc, procIo, err := safeexec.Exec(execConfig, cmd)
			if err != nil {
				return err
			}

			if err = proc.Wait(); err != nil {
				ex := &exec.ExitError{}
				if errors.As(err, &ex) {
					b, err := io.ReadAll(procIo.Stderr)
					if err == nil {
						return &CompilationError{string(b)}
					}
				}
				return err
			}
		}

		binaryPath = filepath.Join(dir, execFileName)
		return nil
	})

	if err != nil {
		return "", err
	}
	return binaryPath, nil
}
