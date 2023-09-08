package fs

import (
	"errors"
	"fmt"
	"os"
)

func EnsureDir(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}

	if !info.IsDir() {
		return errors.New(fmt.Sprintf("%s exists, but not a directory", dirPath))
	}

	return nil
}
