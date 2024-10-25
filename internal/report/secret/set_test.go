package secret

import (
	"reflect"
	"testing"
)

var testSlice []*Secret = []*Secret{&test1, &test2, &test3}

var test4 Secret = Secret{
	Repository:  "test1_repo",
	Path:        "test1_path",
	Kind:        SecretKindAWS,
	Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
	Line:        32,
	Valid:       SecretValidityUnknown,
	Scanners:    SecretScannersGitleaks,
	Cleartext:   "test1_cleartext",
	Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32",
}

var test5 Secret = Secret{
	Repository:  "test5_repo",
	Path:        "test5_path",
	Kind:        SecretKindAbstract,
	Commit:      "00000fa38de1bb02673cd2496a4663a3c6c42508",
	Line:        28,
	Valid:       SecretValidityValid,
	Scanners:    SecretScannersTrufflehog,
	Cleartext:   "test5_cleartext",
	Fingerprint: "test5_repo:00000fa38de1bb02673cd2496a4663a3c6c42508:test5_path:28",
}

var test6 Secret = Secret{
	Repository:  "test6_repo",
	Path:        "test6_path",
	Kind:        SecretKindAWS,
	Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
	Line:        32,
	Valid:       SecretValidityUnknown,
	Scanners:    SecretScannersAll,
	Cleartext:   "test1_cleartext",
	Fingerprint: "test6_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test6_path:32",
}

func TestNewSet(t *testing.T) {
	var test interface{} = NewSet()
	want := make(SecretSet)

	if _, ok := test.(SecretSet); ok == false {
		t.Fatalf(`NewSet().(type) is not %q`, "SecretSet")
	}

	if !reflect.DeepEqual(test, want) {
		t.Fatalf(`NewSet() = %v, want %v`, test, want)
	}
}

func TestNewSetFromSlice(t *testing.T) {
	test := NewSetFromSlice(testSlice)
	want := SecretSet{
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		"hardcoded_fingerprint": &Secret{
			Repository:  "test2_repo",
			Path:        "test2_path",
			Kind:        SecretKindGeneric,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
			Line:        42,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "",
			Fingerprint: "hardcoded_fingerprint",
		},
	}

	if !reflect.DeepEqual(test, want) {
		t.Fatalf(`NewSetFromSlice([]{test1, test2, test3}) = %v, want %v`, test, want)
	}
}

func TestSecretAdd(t *testing.T) {
	set := NewSetFromSlice(testSlice)
	want := SecretSet{
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		"hardcoded_fingerprint": &Secret{
			Repository:  "test2_repo",
			Path:        "test2_path",
			Kind:        SecretKindGeneric,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
			Line:        42,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "",
			Fingerprint: "hardcoded_fingerprint",
		},
	}

	set.Add(&test1)
	want["test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10"] = &Secret{
		Repository:  "test1_repo",
		Path:        "test1_path",
		Kind:        SecretKindAMQP,
		Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
		Line:        10,
		Valid:       SecretValidityInvalid,
		Scanners:    SecretScannersAll,
		Cleartext:   "test1_cleartext",
		Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Add(test4) -> set = %v, want %v`, set, want)
	}

	set.Add(&test4)
	want["test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32"] = &Secret{
		Repository:  "test1_repo",
		Path:        "test1_path",
		Kind:        SecretKindAWS,
		Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
		Line:        32,
		Valid:       SecretValidityUnknown,
		Scanners:    SecretScannersGitleaks,
		Cleartext:   "test1_cleartext",
		Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32",
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Add(test4) -> set = %v, want %v`, set, want)
	}

	set.Add(&test5)
	want["test5_repo:00000fa38de1bb02673cd2496a4663a3c6c42508:test5_path:28"] = &Secret{
		Repository:  "test5_repo",
		Path:        "test5_path",
		Kind:        SecretKindAbstract,
		Commit:      "00000fa38de1bb02673cd2496a4663a3c6c42508",
		Line:        28,
		Valid:       SecretValidityValid,
		Scanners:    SecretScannersTrufflehog,
		Cleartext:   "test5_cleartext",
		Fingerprint: "test5_repo:00000fa38de1bb02673cd2496a4663a3c6c42508:test5_path:28",
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Add(test5) -> set = %v, want %v`, set, want)
	}

	set.Add(&test6)
	want["test6_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test6_path:32"] = &Secret{
		Repository:  "test6_repo",
		Path:        "test6_path",
		Kind:        SecretKindAWS,
		Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
		Line:        32,
		Valid:       SecretValidityUnknown,
		Scanners:    SecretScannersAll,
		Cleartext:   "test1_cleartext",
		Fingerprint: "test6_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test6_path:32",
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Add(test6) -> set = %v, want %v`, set, want)
	}
}

func TestSecretHas(t *testing.T) {
	set := NewSetFromSlice(testSlice)

	test := set.Has(&test1)
	want := true

	if test != want {
		t.Fatalf(`set.Has(test1) = %t, want %t`, test, want)
	}

	test = set.Has(&test2)
	want = true

	if test != want {
		t.Fatalf(`set.Has(test2) = %t, want %t`, test, want)
	}

	test = set.Has(&test3)
	want = true

	if test != want {
		t.Fatalf(`set.Has(test3) = %t, want %t`, test, want)
	}

	test = set.Has(&test4)
	want = false

	if test != want {
		t.Fatalf(`set.Has(test4) = %t, want %t`, test, want)
	}

	test = set.Has(&test5)
	want = false

	if test != want {
		t.Fatalf(`set.Has(test5) = %t, want %t`, test, want)
	}

	test = set.Has(&test6)
	want = false

	if test != want {
		t.Fatalf(`set.Has(test6) = %t, want %t`, test, want)
	}

	set.Add(&test6)

	test = set.Has(&test6)
	want = true

	if test != want {
		t.Fatalf(`set.Has(test6) = %t, want %t`, test, want)
	}
}

func TestSecretRemove(t *testing.T) {
	set := NewSetFromSlice(testSlice)
	set.Add(&test4)
	set.Add(&test5)
	set.Add(&test6)

	want := SecretSet{
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		"hardcoded_fingerprint": &Secret{
			Repository:  "test2_repo",
			Path:        "test2_path",
			Kind:        SecretKindGeneric,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
			Line:        42,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "",
			Fingerprint: "hardcoded_fingerprint",
		},
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAWS,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        32,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersGitleaks,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32",
		},
		"test5_repo:00000fa38de1bb02673cd2496a4663a3c6c42508:test5_path:28": &Secret{
			Repository:  "test5_repo",
			Path:        "test5_path",
			Kind:        SecretKindAbstract,
			Commit:      "00000fa38de1bb02673cd2496a4663a3c6c42508",
			Line:        28,
			Valid:       SecretValidityValid,
			Scanners:    SecretScannersTrufflehog,
			Cleartext:   "test5_cleartext",
			Fingerprint: "test5_repo:00000fa38de1bb02673cd2496a4663a3c6c42508:test5_path:28",
		},
		"test6_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test6_path:32": &Secret{
			Repository:  "test6_repo",
			Path:        "test6_path",
			Kind:        SecretKindAWS,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        32,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test6_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test6_path:32",
		},
	}

	set.Remove(&test2)
	delete(want, test2.Fingerprint)

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Remove(test2) -> set = %v, want %v`, set, want)
	}

	set.Remove(&test3)
	delete(want, test3.Fingerprint)

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Remove(test3) -> set = %v, want %v`, set, want)
	}

	set.Remove(&test6)
	delete(want, test6.Fingerprint)

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set.Remove(test6) -> set = %v, want %v`, set, want)
	}
}

func TestSecretLength(t *testing.T) {
	set := NewSet()

	test := set.Length()
	want := 0

	if test != want {
		t.Fatalf(`set.Length() = %d, want %d`, test, want)
	}

	set = NewSetFromSlice(testSlice)

	test = set.Length()
	want = 2

	if test != want {
		t.Fatalf(`set.Length() = %d, want %d`, test, want)
	}

	set.Add(&test6)

	test = set.Length()
	want = 3

	if test != want {
		t.Fatalf(`set.Length() = %d, want %d`, test, want)
	}

	set.Add(&test4)

	test = set.Length()
	want = 4

	if test != want {
		t.Fatalf(`set.Length() = %d, want %d`, test, want)
	}
}

func TestSecretClone(t *testing.T) {
	set := NewSetFromSlice(testSlice)

	test := set.Clone()
	want := set

	if !reflect.DeepEqual(test, want) {
		t.Fatalf(`set.Clone() = %v, want %v`, test, want)
	}
}

func TestSecretUnion(t *testing.T) {
	set1 := NewSetFromSlice(testSlice)

	set := set1.Union(NewSet())
	want := set1

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set1.Union(NewSet()) = %v, want %v`, set, want)
	}

	set2 := NewSet()
	set2.Add(&test2)
	set2.Add(&test3)
	set2.Add(&test4)

	set = set1.Union(set2)
	want = SecretSet{
		"hardcoded_fingerprint": &Secret{
			Repository:  "test2_repo",
			Path:        "test2_path",
			Kind:        SecretKindGeneric,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
			Line:        42,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "",
			Fingerprint: "hardcoded_fingerprint",
		},
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAWS,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        32,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersGitleaks,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32",
		},
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set1.Union(set2) = %v, want %v`, set, want)
	}
}

func TestSecretDiff(t *testing.T) {
	set1 := NewSetFromSlice(testSlice)

	set := set1.Diff(NewSet())
	want := set1

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set1.Diff(NewSet()) = %v, want %v`, set, want)
	}

	set2 := NewSet()
	set2.Add(&test2)
	set2.Add(&test3)
	set2.Add(&test4)

	set = set1.Diff(set2)
	want = SecretSet{}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set1.Diff(set2) = %v, want %v`, set, want)
	}

	set = set2.Diff(set1)
	want = SecretSet{
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAWS,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        32,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersGitleaks,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:32",
		},
	}

	if !reflect.DeepEqual(set, want) {
		t.Fatalf(`set2.Diff(set1) = %v, want %v`, set, want)
	}
}

func TestSecretDropFingerprints(t *testing.T) {
	set := NewSetFromSlice(testSlice)
	fps := []string{"hardcoded_fingerprint"}

	test := set.DropFingerprints(fps)
	want := SecretSet{
		"test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10": &Secret{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
	}

	if !reflect.DeepEqual(test, want) {
		t.Fatalf(`set.DropFingerprints(fps) = %v, want %v`, test, want)
	}
}

func TestSecretToSlice(t *testing.T) {
	set := NewSetFromSlice(testSlice)

	test := set.ToSlice()
	want := []*Secret{
		{
			Repository:  "test1_repo",
			Path:        "test1_path",
			Kind:        SecretKindAMQP,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
			Line:        10,
			Valid:       SecretValidityInvalid,
			Scanners:    SecretScannersAll,
			Cleartext:   "test1_cleartext",
			Fingerprint: "test1_repo:e2acefa38de1bb02673cd2496a4663a3c6c42508:test1_path:10",
		},
		{
			Repository:  "test2_repo",
			Path:        "test2_path",
			Kind:        SecretKindGeneric,
			Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
			Line:        42,
			Valid:       SecretValidityUnknown,
			Scanners:    SecretScannersAll,
			Cleartext:   "",
			Fingerprint: "hardcoded_fingerprint",
		},
	}

	if !reflect.DeepEqual(test, want) {
		t.Fatalf(`set.ToSlice() = %v, want %v`, test, want)
	}
}
