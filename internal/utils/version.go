package utils

import (
	"fmt"

	"golang.org/x/mod/semver"
)

func IsVersionValid(version string, minVersion string) (bool, error) {
	if !semver.IsValid(fmt.Sprintf("v%s", version)) {
		return false, fmt.Errorf("invalid version: %q", version)
	}

	if semver.Compare(fmt.Sprintf("v%s", version), fmt.Sprintf("v%s", minVersion)) < 0 {
		return false, nil
	}

	return true, nil
}
