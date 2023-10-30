package fs

import (
	"fmt"
	"os"
	"time"
)

func Touch(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Create an empty file if it doesn't exist
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		return err
	}

	// Update the access and modification times to the current time
	currentTime := time.Now()
	return os.Chtimes(filename, currentTime, currentTime)
}

func EnsureDir(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s exists, but is not a directory", dirPath)
	}

	return nil
}

// IsStale checks if the srcFilename is updated after the derivedFilename was made.
// Returns true if the source file is updated after the derived was made.
func IsStale(srcFilename, derivedFilename string) (bool, error) {
	srcInfo, err := os.Stat(srcFilename)
	if err != nil {
		return false, err
	}

	derivedInfo, err := os.Stat(derivedFilename)
	if err != nil {
		return false, err
	}

	return srcInfo.ModTime().After(derivedInfo.ModTime()), nil
}

// GetDurationSinceModified returns the duration since the file was last modified.
func GetDurationSinceModified(filename string) (time.Duration, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	modificationTime := fileInfo.ModTime()
	currentTime := time.Now()
	durationSinceModified := currentTime.Sub(modificationTime)

	return durationSinceModified, nil
}
