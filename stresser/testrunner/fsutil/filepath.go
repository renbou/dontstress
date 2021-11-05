package fsutil

import (
	"path/filepath"
	"strings"
)

func TrimExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
