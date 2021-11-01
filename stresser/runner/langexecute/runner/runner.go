package runner

import (
	"path/filepath"
	"time"

	"github.com/renbou/dontstress/stresser/runner/langexecute/util"
)

// User-time error with info about when and what
type UserError struct {
	cause   string
	message string
}

func UserErrorWithCause(cause string) func(string) *UserError {
	return func(message string) *UserError {
		return &UserError{cause, message}
	}
}

var (
	NewCompilationError = UserErrorWithCause("ce")
	NewRuntimeError     = UserErrorWithCause("re")
	NewMemoryLimitError = UserErrorWithCause("ml")
)

func (ue *UserError) Error() string {
	return ue.message
}

func (ue *UserError) Cause() string {
	return ue.cause
}

type Runner interface {
	Prepare(path string) error
	Run(input string) (error, string)
}

type StepRunner struct {
	CompileStep []string // commands to execute in the testing directory during compilation steps
	RunStep     string   // how to launch the executable
	directory   string   // temporary directory of current runner
	file        string   // base filename
	config      util.KV  // config for templates
}

func (s *StepRunner) Prepare(path string) error {
	err, _ := util.WithTempDir(func(dir string) (error, interface{}) {
		s.directory = dir
		s.file = filepath.Base(path)
		s.config = util.KV{
			"File": s.file,
			"Out":  util.RandomFileName(),
		}

		if err := util.CopyFile(path, s.file); err != nil {
			return err, nil
		}
		for _, cmd := range s.CompileStep {
			if err, result := util.Exec("", time.Second*10, cmd, s.config); err != nil {
				return err, nil
			} else if result.Status != util.StatusOK {
				return NewCompilationError(result.Stderr), nil
			}
		}
		return nil, nil
	})
	return err
}

func (s *StepRunner) Run(input string) (error, string) {
	err, data := util.WithDir(func(_ string) (error, interface{}) {
		if err, result := util.Exec(input, time.Second*5, s.RunStep, s.config); err != nil {
			return err, ""
		} else {
			switch result.Status {
			case util.StatusML:
				return NewMemoryLimitError(result.Stderr), result.Stdout
			case util.StatusRE:
				return NewRuntimeError(result.Stderr), result.Stdout
			default:
				return nil, result.Stdout
			}
		}
	}, s.directory)
	return err, data.(string)
}
