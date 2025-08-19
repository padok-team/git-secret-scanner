package trufflehog

import (
	"os/exec"
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/utils"
)

const MinVersion string = "3.90.3"

func CommandExists() bool {
	return utils.CommandExists("trufflehog")
}

func Version() (string, error) {
	cmd := exec.Command("trufflehog", "--version")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	split := strings.Split(string(output), " ")
	version := strings.TrimSuffix(split[len(split)-1], "\n")

	return version, nil
}

func IsVersionValid() (bool, string, error) {
	version, err := Version()
	if err != nil {
		return false, "", err
	}

	ok, err := utils.IsVersionValid(version, MinVersion)
	return ok, version, err
}
