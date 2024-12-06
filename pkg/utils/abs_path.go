package utils

import (
	"os"
	"path/filepath"
)

var absPath = ""

func InitAbsolutePath() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	rootPath, err := findProjectRoot(wd)
	if err != nil {
		return err
	}
	absPath = rootPath
	return nil
}

func GetAbsPath() string {
	return absPath
}

func findProjectRoot(startDir string) (string, error) {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", os.ErrNotExist
		}
		dir = parentDir
	}
}
