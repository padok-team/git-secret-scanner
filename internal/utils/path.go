package utils

import (
	"os"
	"os/exec"
	"path"
)

const (
	TempDirName string = "github.padok.git-secret-scanner"
)

func TempDirPath() string {
	return path.Join(os.TempDir(), TempDirName)
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FileExistsAndNotEmpty(path string) bool {
	s, err := os.Stat(path)
	return (err == nil) && s.Size() > 0
}
