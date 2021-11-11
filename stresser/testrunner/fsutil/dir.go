package fsutil

import "os"

func SwitchDir(path string) (oldpwd string, err error) {
	if oldpwd, err = os.Getwd(); err != nil {
		return "", err
	} else if err = os.Chdir(path); err != nil {
		return "", err
	}
	return oldpwd, nil
}

func WithDir(wrapped func() error, dir string) error {
	old, err := SwitchDir(dir)
	if err != nil {
		return err
	}

	defer SwitchDir(old)
	return wrapped()
}

func WithTempDir(wrapped func() error) error {
	dir, err := os.MkdirTemp("", "stress")
	if err != nil {
		return err
	}

	defer os.Remove(dir)
	return WithDir(wrapped, dir)
}
