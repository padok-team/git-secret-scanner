package utils

import (
	"fmt"

	"golang.org/x/mod/semver"
)

func IsVersionValid(version string, minVersion string) (bool, error) {
	if !semver.IsValid(version) {
		return false, fmt.Errorf("invalid version: %q", version)
	}

	if semver.Compare(version, minVersion) < 0 {
		return false, nil
	}

	return true, nil
}
