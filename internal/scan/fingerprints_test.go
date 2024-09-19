package scan

import (
	"reflect"
	"testing"
)

func TestLoadIgnoredFingerprints(t *testing.T) {
	err := LoadIgnoredFingerprints("../../test/testdata/.fingerprintsignore")
	if err != nil {
		t.Fatalf(`LoadIgnoredFingerprints("test/testdata/.fingerprintsignore") = %v, want nil`, err)
	}

	want := map[string][]string{
		"test1_repo": {
			"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		"test2_repo": {
			"test2_repo:e2acefa38de1aaaaa73cd2496a4663a3c6c42509:test2_path:5",
		},
	}

	if !reflect.DeepEqual(ignoredFingerprints, want) {
		t.Fatalf(`ignoredFingerprints = %v, want %v`, ignoredFingerprints, want)
	}
}
