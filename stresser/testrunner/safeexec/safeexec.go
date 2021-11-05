package safeexec

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/nanmu42/limitio"
	"github.com/renbou/dontstress/stresser/testrunner/template"
)

// Custom cmd struct with given method to check if stdout size was exceeded
type Cmd struct {
	*exec.Cmd
	stdoutLimit   int
	limitedStdout *limitio.Writer
}

func (c *Cmd) StdoutThresholdExceeded() bool {
	return c.limitedStdout.Written() > c.stdoutLimit
}

// KillChildren tries killing process' children, should be called after exit.
// You should also set Setpgid and Setsid in SysProcAttr during creation
func (c *Cmd) KillChildren() {
	// First try killing using groupid
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGKILL)
	}
	// Then try killing using session id
	syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
}

func createCmd(ctx context.Context, fmt string, data interface{}) (*exec.Cmd, error) {
	commandString, err := template.Format(fmt, data)
	if err != nil {
		return nil, err
	}
	command := strings.Split(commandString, " ")

	return exec.CommandContext(ctx, command[0], command[1:]...), err

}

type ExecIO struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

type Config struct {
	Context     context.Context
	OutputLimit int
	Data        interface{}
}

// Launches a program with the context defined by the config, immediately returns
// the io stuff to communicate with the launched program. Does not wait for exit
// (caller should though)
func Exec(config *Config, command string) (*Cmd, *ExecIO, error) {
	cmd, err := createCmd(config.Context, command, config.Data)
	if err != nil {
		return nil, nil, err
	}

	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutPipeW := io.Pipe()
	stderrR, stderrPipeW := io.Pipe()

	// Low level limit on all the writers. This will allow writes over the limit so that the program
	// doesn't crash due to EPIPE and we can safely check later that the limit was exceeded
	// Also make these buffered so that the process can start writing straight away
	limited := func(w io.Writer) *limitio.Writer { return limitio.NewWriter(w, config.OutputLimit, true) }
	stdoutW := limited(bufio.NewWriter(stdoutPipeW))
	stderrW := limited(bufio.NewWriter(stderrPipeW))

	// Assign all io to the created pipes
	cmd.Stdin = stdinR
	cmd.Stdout = stdoutW
	cmd.Stderr = stderrW

	// No env for u, bad boi!
	// (prevents AWS credentials leakage, etc... lol)
	cmd.Env = []string{}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Setsid:  true,
		// Pdeathsig: syscall.SIGKILL, // should be available for linux, so remove comment in prod
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return &Cmd{cmd, config.OutputLimit, stdoutW}, &ExecIO{
		Stdin:  stdinW,
		Stdout: stdoutR,
		Stderr: stderrR,
	}, err
}
