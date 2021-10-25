package runner

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/renbou/dontstress/stresser/runner/langexecute/util"
)

const (
	STEP_INIT    = "init"
	STEP_COMPILE = "compile"
	STEP_RUN     = "run"
)

type Runner interface {
	Prepare(path string) error
	Run(r io.Reader, w io.Writer) error
}

type RunError struct {
	runStep string
	err     error
}

func (re *RunError) Step() string {
	return re.runStep
}

func (re *RunError) Unwrap() error {
	return re.err
}

func (re *RunError) Error() string {
	return fmt.Sprintf("error %s during step %s", re.err.Error(), re.runStep)
}

type StepRunner struct {
	CompileStep []string // commands to execute in the testing directory during compilation steps
	RunStep     string   // how to launch the executable
	directory   string   // temporary directory of current runner
	file        string   // base filename
}

func stepError(step string) func(error) *RunError {
	return func(err error) *RunError {
		return &RunError{step, err}
	}
}

var initError = stepError(STEP_INIT)
var compileError = stepError(STEP_COMPILE)
var runError = stepError(STEP_RUN)

func WrapError(wrapper func(error) *RunError, err error) *RunError {
	var re *RunError
	if err != nil && !errors.As(err, &re) {
		return wrapper(err)
	} else {
		return re
	}
}

func (s *StepRunner) Prepare(path string) error {
	err := util.WithTempDir(func(dir string) error {
		s.directory = dir
		s.file = filepath.Base(path)

		if err := util.CopyFile(path, s.file); err != nil {
			return initError(err)
		}
		for _, cmd := range s.CompileStep {
			if err := util.Exec(time.Second*10, cmd, s.file); err != nil {
				return compileError(err)
			}
		}
		return nil
	})
	return WrapError(initError, err)
}

func (s *StepRunner) Run(r io.Reader, w io.Writer) error {
	err := util.WithDir(func(_ string) error {
		if err := util.RWExec(r, w, time.Second*5, s.RunStep, s.file); err != nil {
			return runError(err)
		}
		return nil
	}, s.directory)
	return WrapError(runError, err)
}
