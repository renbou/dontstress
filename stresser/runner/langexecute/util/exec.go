package util

import (
	"context"
	"errors"
	"io"
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

func Exec(d time.Duration, fmt string, data interface{}) error {
	if cmd, err := Cmd(d, fmt, data); err != nil {
		return err
	} else if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}
	return nil
}

func RWExec(r io.Reader, w io.Writer, d time.Duration, fmt string, data interface{}) error {
	cmd, err := Cmd(d, fmt, data)
	if err != nil {
		return err
	}
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
