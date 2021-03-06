package helper

import (
	"os"
)

// PathExists ...
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkdirAll ...
func MkdirAll(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

// IfIsDir ...
func IfIsDir(path string) bool {
	f, _ := os.Stat(path)
	return f.IsDir()
}
