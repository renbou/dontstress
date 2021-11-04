package util

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

func Cmd(d time.Duration, fmt string, data interface{}) (*exec.Cmd, error) {
	commandString, err := Format(fmt, data)
	if err != nil {
		return nil, err
	}
	command := strings.Split(commandString, " ")

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	return exec.CommandContext(ctx, command[0], command[1:]...), err

}

var (
	StatusOK = 0
	StatusRE = 1
	StatusML = 2
)

type ExecResult struct {
	Status int
	Code   int
	Stdout string
	Stderr string
}

func Exec(input string, d time.Duration, fmt string, data interface{}) (error, *ExecResult) {
	cmd, err := Cmd(d, fmt, data)
	if err != nil {
		return err, nil
	}

	cmd.Stdin = strings.NewReader(input)
	// 64 megabytes worth of data
	stdout := bytes.NewBuffer(make([]byte, 63_000_000))
	stderr := bytes.NewBuffer(make([]byte, 1_000_000))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			// TODO
			// Add memory limit error
			return nil, &ExecResult{
				Status: StatusRE,
				Code:   exitErr.ExitCode(),
				Stdout: stdout.String(),
				Stderr: stderr.String(),
			}
		} else {
			return err, nil
		}
	}
	return nil, &ExecResult{
		Status: StatusOK,
		Code:   0,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
}
