package gitleaks

import (
	"os/exec"
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/utils"
)

const MinVersion string = "8.28.0"

func CommandExists() bool {
	return utils.CommandExists("gitleaks")
}

func Version() (string, error) {
	cmd := exec.Command("gitleaks", "version")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func IsVersionValid() (bool, string, error) {
	version, err := Version()
	if err != nil {
		return false, "", err
	}

	ok, err := utils.IsVersionValid(version, MinVersion)
	return ok, version, err
}
