package utils

import (
	"os"
	"path"
)

const (
	TempDirName string = "github.padok.git-secret-scanner"
)

func TempDirPath() string {
	return path.Join(os.TempDir(), TempDirName)
}

func FileExistsAndNotEmpty(path string) bool {
	s, err := os.Stat(path)
	return (err == nil) && s.Size() > 0
}
