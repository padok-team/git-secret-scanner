package scan

import (
	"reflect"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
)

func TestLoadBaseline(t *testing.T) {
	err := LoadBaseline("../../test/testdata/baseline.csv")
	if err != nil {
		t.Fatalf(`LoadBaseline("test/testdata/baseline.csv") = %v, want nil`, err)
	}

	want := map[string]secret.SecretSet{
		"test1_repo": {
			"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &secret.Secret{
				Repository:  "test1_repo",
				Path:        "test1_path",
				Kind:        secret.SecretKindGeneric,
				Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
				Line:        10,
				Valid:       secret.SecretValidityInvalid,
				Cleartext:   "test1_cleartext",
				Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
			},
			"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:5": &secret.Secret{
				Repository:  "test1_repo",
				Path:        "test1_path",
				Kind:        secret.SecretKindAMQP,
				Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
				Line:        5,
				Valid:       secret.SecretValidityUnknown,
				Cleartext:   "test1_cleartext",
				Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:5",
			},
		},
		"test2_repo": {
			"hardcoded_fingerprint": &secret.Secret{
				Repository:  "test2_repo",
				Path:        "test2_path",
				Kind:        secret.SecretKindGeneric,
				Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
				Line:        42,
				Valid:       secret.SecretValidityUnknown,
				Cleartext:   "",
				Fingerprint: "hardcoded_fingerprint",
			},
		},
	}

	if !reflect.DeepEqual(baseline, want) {
		t.Fatalf(`baseline = %v, want %v`, baseline, want)
	}
}
