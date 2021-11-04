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

func WithDir(wrapped func(string) (error, interface{}), dir string) (error, interface{}) {
	if old, err := SwitchDir(dir); err != nil {
		return err, nil
	} else {
		defer SwitchDir(old)
	}
	return wrapped(dir)
}

func WithTempDir(wrapped func(string) (error, interface{})) (error, interface{}) {
	if dir, err := os.MkdirTemp("", "stress"); err != nil {
		return err, nil
	} else {
		defer func() {
			if err != nil {
				os.Remove(dir)
			}
		}()
		return WithDir(wrapped, dir)
	}
}
