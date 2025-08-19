package utils

import "testing"

func TestIsVersionValid(t *testing.T) {
	test, err := IsVersionValid("1.1.0", "1.1.0")
	want := true

	if test != want || err != nil {
		t.Fatalf(`IsVersionValid("1.1.0", "1.1.0") = %v, %v, want true, nil`, test, err)
	}

	test, err = IsVersionValid("1.1.0", "1.2.0")
	want = false

	if test != want || err != nil {
		t.Fatalf(`IsVersionValid("1.1.0", "1.2.0") = %v, %v, want false, nil`, test, err)
	}

	test, err = IsVersionValid("1.2.0", "1.1.0")
	want = true

	if test != want || err != nil {
		t.Fatalf(`IsVersionValid("1.2.0", "1.1.0") = %v, %v, want true, nil`, test, err)
	}
}
