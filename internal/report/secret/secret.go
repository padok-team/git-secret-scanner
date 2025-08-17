package secret

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SecretValidity int

const (
	SecretValidityUnknown SecretValidity = iota
	SecretValidityValid
	SecretValidityInvalid
)

func (sv SecretValidity) String() string {
	switch sv {
	case SecretValidityValid:
		return "true"
	case SecretValidityInvalid:
		return "false"
	case SecretValidityUnknown:
		return ""
	default:
		// should never be reached
		return ""
	}
}

// implements json.Marshaler
func (sv SecretValidity) MarshalJSON() ([]byte, error) {
	return json.Marshal(sv.String())
}

// implements json.Unmarshaler
func (sv *SecretValidity) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case SecretValidityValid.String():
		*sv = SecretValidityValid
		return nil
	case SecretValidityInvalid.String():
		*sv = SecretValidityInvalid
		return nil
	case SecretValidityUnknown.String():
		*sv = SecretValidityUnknown
		return nil
	default:
		return fmt.Errorf(
			"validity must be one of \"%s\", \"%s\" or \"%s\"",
			SecretValidityValid.String(),
			SecretValidityInvalid.String(),
			SecretValidityUnknown.String(),
		)
	}
}

// implements gocsv.TypeUnmarshaller
func (sv *SecretValidity) UnmarshalCSV(s string) error {
	switch s {
	case SecretValidityValid.String():
		*sv = SecretValidityValid
		return nil
	case SecretValidityInvalid.String():
		*sv = SecretValidityInvalid
		return nil
	case SecretValidityUnknown.String():
		*sv = SecretValidityUnknown
		return nil
	default:
		return fmt.Errorf(
			"validity must be one of \"%s\", \"%s\" or \"%s\"",
			SecretValidityValid.String(),
			SecretValidityInvalid.String(),
			SecretValidityUnknown.String(),
		)
	}
}

type SecretScanners int

const (
	SecretScannersAll SecretScanners = iota
	SecretScannersGitleaks
	SecretScannersTrufflehog
)

func (ss SecretScanners) String() string {
	switch ss {
	case SecretScannersAll:
		return "all"
	case SecretScannersGitleaks:
		return "gitleaks"
	case SecretScannersTrufflehog:
		return "trufflehog"
	default:
		// should never be reached
		return ""
	}
}

// implements json.Marshaler
func (ss SecretScanners) MarshalJSON() ([]byte, error) {
	return json.Marshal(ss.String())
}

// implements json.Unmarshaler
func (ss *SecretScanners) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case SecretScannersAll.String():
		*ss = SecretScannersAll
		return nil
	case SecretScannersGitleaks.String():
		*ss = SecretScannersGitleaks
		return nil
	case SecretScannersTrufflehog.String():
		*ss = SecretScannersTrufflehog
		return nil
	default:
		return fmt.Errorf(
			"scanners must be one of \"%s\", \"%s\" or \"%s\"",
			SecretScannersAll.String(),
			SecretScannersGitleaks.String(),
			SecretScannersTrufflehog.String(),
		)
	}
}

// implements gocsv.TypeUnmarshaller
func (ss *SecretScanners) UnmarshalCSV(s string) error {
	switch s {
	case SecretScannersAll.String():
		*ss = SecretScannersAll
		return nil
	case SecretScannersGitleaks.String():
		*ss = SecretScannersGitleaks
		return nil
	case SecretScannersTrufflehog.String():
		*ss = SecretScannersTrufflehog
		return nil
	default:
		return fmt.Errorf(
			"scanners must be one of \"%s\", \"%s\" or \"%s\"",
			SecretScannersAll.String(),
			SecretScannersGitleaks.String(),
			SecretScannersTrufflehog.String(),
		)
	}
}

var ErrMergeSecretsNotEqual error = errors.New("not equal secrets cannot be merge")

type Secret struct {
	Repository  string         `csv:"repository" json:"repository"`
	Path        string         `csv:"path" json:"path"`
	Kind        SecretKind     `csv:"kind" json:"kind"`
	Commit      string         `csv:"commit" json:"commit"`
	Line        int            `csv:"line" json:"line"`
	Valid       SecretValidity `csv:"valid" json:"valid"`
	Scanners    SecretScanners `csv:"scanners" json:"scanners"`
	Cleartext   string         `csv:"cleartext" json:"cleartext"`
	Fingerprint string         `csv:"fingerprint" json:"fingerprint"`
}

func NewSecret(repository string, path string, kind SecretKind, commit string, line int, valid SecretValidity, scanners SecretScanners, cleartext string, fingerprint string) (*Secret, error) {
	if fingerprint == "" {
		fingerprint = fmt.Sprintf("%s:%s:%s:%d", repository, commit, path, line)
	}

	return &Secret{
		Repository:  repository,
		Path:        path,
		Kind:        kind,
		Commit:      commit,
		Line:        line,
		Valid:       valid,
		Scanners:    scanners,
		Cleartext:   cleartext,
		Fingerprint: fingerprint,
	}, nil
}

func (s Secret) IsEqual(other *Secret) bool {
	return s.Fingerprint == other.Fingerprint
}

func (s Secret) Merge(other *Secret) (*Secret, error) {
	if s.IsEqual(other) {
		kind := s.Kind
		if kind == SecretKindGeneric {
			kind = other.Kind
		}

		valid := s.Valid
		if valid == SecretValidityUnknown {
			valid = other.Valid
		}

		scanners := s.Scanners
		if scanners != SecretScannersAll && scanners != other.Scanners {
			scanners = SecretScannersAll
		}

		return &Secret{
			Repository:  s.Repository,
			Path:        s.Path,
			Kind:        kind,
			Commit:      s.Commit,
			Line:        s.Line,
			Valid:       valid,
			Scanners:    scanners,
			Cleartext:   s.Cleartext,
			Fingerprint: s.Fingerprint,
		}, nil
	} else {
		return nil, ErrMergeSecretsNotEqual
	}
}
