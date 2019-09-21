package errors

import (
	"testing"
)

func TestNotNil(t *testing.T) {

	notNil, err := NotNil(nil)
	if !notNil {
		t.Errorf("expected false on nil value: err =  %v", err)
	}

	var x *int

	notNil, err = NotNil(x)
	if notNil {
		t.Errorf("failed to detect null pointer variable: notNil = %v, x=%v", notNil, x)
	}

	// A NON POINTER VALUE (EVEN IF ZERRO) SHOULD NOT BE NIL
	var xx int = 0
	notNil, err = NotNil(xx)
	if notNil == false {
		t.Errorf("failed detect a non pointer type as not nil: notNil=%v,x=%v", notNil, xx)
	}

	var s string

	notNil, err = NotNil(s)
	if !notNil {
		t.Errorf("failed detect a non pointer type as not nil: notNil=%v,s='%v'", notNil, s)
	}

}

func TestStrLenGtZero(t *testing.T) {

	// SLICE OF NON STRINGS GT ZERO WITH ALPHA,NUMERIC,AND SPECIAL CHARS PERMUTATIONS
	slist := []string{" ", "abc", "123", "@#$@", "s1@#$^"}

	for i, s := range slist {
		if ok, err := StrlenGtZero(s); !ok {
			t.Errorf("failed to detect empthy space string gt zero length: i=%d: %v", i, err)
		}
	}

	s := ""
	if ok, err := StrlenGtZero(s); ok {
		t.Errorf("failed to detect zero length string error=%v", err)
	}

}

func TestStrAlpha(t *testing.T) {

	slist := []string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcABC"}

	for i, s := range slist {
		if ok, err := StrAlpha(s); !ok {
			t.Errorf("failed to detect alpha string '%s': i=%d: error=%v", s, i, err)
		}
	}

}

func TestPathElementValid(t *testing.T) {

	slist := []string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789", "abc123", "ABC123", "abc 12-ab_CD"}

	for i, s := range slist {
		if ok, err := PathElementValid(s); !ok {
			t.Errorf("failed to detect valid path element '%s': i=%d, error=%v", s, i, err)

		}
	}

	slist = []string{"abc()", "#sfsd", "23:*", "dfd\\dds"}
	for i, s := range slist {
		if ok, err := PathElementValid(s); ok {
			t.Errorf("failed to detect invalid path element '%s': i=%d, error=%v", s, i, err)

		}
	}

}
