package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/renbou/dontstress/stresser/testrunner/fsutil"
	"github.com/renbou/dontstress/stresser/testrunner/safeexec"
)

// User-time error with info about when and what
type UserError struct {
	cause   string
	message string
	code    int
}

func UserErrorWithCause(cause string) func(int, string) *UserError {
	return func(code int, message string) *UserError {
		return &UserError{cause, message, code}
	}
}

var (
	NewCompilationError  = UserErrorWithCause("ce")
	NewRuntimeError      = UserErrorWithCause("re")
	NewPresentationError = UserErrorWithCause("pe")
)

func (ue *UserError) Error() string {
	return ue.message
}

func (ue *UserError) Cause() string {
	return ue.cause
}

func (ue *UserError) Code() int {
	return ue.code
}

type Options struct {
	// Path to executable
	Path string
	// Limit of stdout and stderr
	OutputLimit int
	// Context under which all execution takes place,
	// so that we can quit instantly if necessary
	Context context.Context
	// Channel to send communication channels with launched process
	// MUST send data here first (before Result)
	IO chan<- *safeexec.ExecIO
	// Channel to communicate end of run
	Result chan<- error
}

type LangRunner interface {
	// Returns the output
	Run(*Options)
}

type StepRunner struct {
	compileSteps []string
	execStep     string
}

// Config to pass as template parameters to each step
type templateConfig struct {
	Out  string
	File string
}

// If err is an ExitError, then builds a new error using the given builder
func handleExitError(err error, stderr io.Reader, errBuilder func(int, string) *UserError) error {
	ex := &exec.ExitError{}
	if errors.As(err, &ex) {
		b, err := io.ReadAll(stderr)
		if err != nil {
			return err
		} else {
			return errBuilder(ex.ExitCode(), string(b))
		}
	}
	return err
}

func (s *StepRunner) Run(options *Options) {
	// Make sure we send io only once
	var sentIO = false
	err := fsutil.WithTempDir(func() error {
		// Copy the requested file into our temp directory
		filename := fsutil.RandomFileName()
		if err := fsutil.CopyFile(options.Path, filename); err != nil {
			return err
		}

		execConfig := &safeexec.Config{
			Context:     options.Context,
			OutputLimit: options.OutputLimit,
			Data: &templateConfig{
				Out:  fsutil.RandomFileName(),
				File: filename,
			},
		}

		// Run all compilation steps, if they fail
		// then we have a compilation error (or an internal error)
		for _, cmd := range s.compileSteps {
			proc, procIo, err := safeexec.Exec(execConfig, cmd)
			if err != nil {
				return err
			}

			if err = proc.Wait(); err != nil {
				return handleExitError(err, procIo.Stderr, NewCompilationError)
			}
		}

		proc, procIo, err := safeexec.Exec(execConfig, s.execStep)
		if err != nil {
			return err
		}

		// Finally...
		options.IO <- procIo
		sentIO = true

		// Wait until execution finishes, kill children and then handle errors
		err = proc.Wait()
		proc.KillChildren()
		if err != nil {
			return handleExitError(err, procIo.Stderr, NewRuntimeError)
		}
		if proc.StdoutThresholdExceeded() {
			// Written more than we were allowed to!
			return NewPresentationError(0,
				fmt.Sprintf("written more bytes to stdout than the limit (%d)", options.OutputLimit))
		}

		return nil
	})
	if !sentIO {
		// send nil io if we have errored out before reaching execution
		options.IO <- nil
	}
	options.Result <- err
}

// Gets runner for lang. If runner is not defined returns false as second value
func GetRunner(lang string) (LangRunner, bool) {
	runner, ok := runners[lang]
	return runner, ok
}
