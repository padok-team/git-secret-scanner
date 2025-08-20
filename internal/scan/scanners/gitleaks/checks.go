package gitleaks

import (
	"os/exec"
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/utils"
)

const MinVersion string = "8.28.0"

var gitleaksCommand string = "gitleaks"

func SetCommandPath(path string) {
	gitleaksCommand = path
}

func Version() (string, error) {
	cmd := exec.Command(gitleaksCommand, "version")

	output, err := cmd.CombinedOutput()
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
