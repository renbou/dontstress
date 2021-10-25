package util

import "os"

func SwitchDir(path string) (oldpwd string, err error) {
	if oldpwd, err = os.Getwd(); err != nil {
		return "", err
	} else if err = os.Chdir(path); err != nil {
		return "", err
	}
	return oldpwd, nil
}

func WithDir(wrapped func(string) error, dir string) error {
	if old, err := SwitchDir(dir); err != nil {
		return err
	} else {
		defer SwitchDir(old)
	}

	if err := wrapped(dir); err != nil {
		return err
	}
	return nil
}

func WithTempDir(wrapped func(string) error) error {
	if dir, err := os.MkdirTemp("", "stress"); err != nil {
		return err
	} else {
		defer func() {
			if err != nil {
				os.Remove(dir)
			}
		}()
		return WithDir(wrapped, dir)
	}
}
