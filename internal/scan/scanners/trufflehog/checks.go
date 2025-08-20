package trufflehog

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/utils"
)

const MinVersion string = "v3.90.3"

var trufflehogCommand string = "trufflehog"

func SetCommandPath(path string) {
	trufflehogCommand = path
}

func Version() (string, error) {
	cmd := exec.Command(trufflehogCommand, "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	split := strings.Split(string(output), " ")
	version := fmt.Sprintf("v%s", strings.TrimPrefix(strings.TrimSuffix(split[len(split)-1], "\n"), "v"))

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
