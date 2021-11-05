package fsutil

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"strconv"
)

func ValidFile(path string) (bool, error) {
	if s, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return s.Mode().IsRegular(), nil
	} else {
		return false, err
	}
}

func CopyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if _, err := io.Copy(toFile, fromFile); err != nil {
		return err
	}
	return nil
}

func RandomFileName() string {
	return strconv.Itoa(rand.Int())
}
