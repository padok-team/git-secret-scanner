package secret

import (
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
		return "True"
	case SecretValidityInvalid:
		return "False"
	case SecretValidityUnknown:
		return ""
	default:
		// should never be reached
		return ""
	}
}

// implements gocsv.TypeUnmarshaller to unmarshall SecretValidity the right way (instead of boolean)
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

var ErrMergeSecretsNotEqual error = errors.New("not equal secrets cannot be merge")

type Secret struct {
	Repository  string         `csv:"repository"`
	Path        string         `csv:"path"`
	Kind        SecretKind     `csv:"kind"`
	Commit      string         `csv:"commit"`
	Line        int            `csv:"line"`
	Valid       SecretValidity `csv:"valid"`
	Cleartext   string         `csv:"cleartext"`
	Fingerprint string         `csv:"fingerprint"`
}

func NewSecret(repository string, path string, kind SecretKind, commit string, line int, valid SecretValidity, cleartext string, fingerprint string) (*Secret, error) {
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

		return &Secret{
			Repository:  s.Repository,
			Path:        s.Path,
			Kind:        kind,
			Commit:      s.Commit,
			Line:        s.Line,
			Valid:       valid,
			Cleartext:   s.Cleartext,
			Fingerprint: s.Fingerprint,
		}, nil
	} else {
		return nil, ErrMergeSecretsNotEqual
	}
}
