package primitives

import (
	"testing"
)

func TestNew(t *testing.T) {

	// SEND NIL FOR subject
	var (
		subject   string
		predicate string
		object    string
	)

	subject = ""
	predicate = "isA"
	object = "dog"
	triple, err := Dogg3rzTripleNew(subject, predicate, object)
	if triple != nil { // EXPECTED NIL BECAUSE ZERO LEN SUBJECT PASSED
		t.Errorf("did not fail on zero length subject: err=%v", err)
	}

	subject = "kai"
	predicate = ""
	triple, err = Dogg3rzTripleNew(subject, predicate, object)
	if triple != nil { // EXPECTED NIL BECAUSE ZERO LEN PREDICATE PASSED
		t.Errorf("did not fail on zero length predicate: err=%v", err)
	}

	subject = "kai"
	predicate = "isA"
	object = ""
	triple, err = Dogg3rzTripleNew(subject, predicate, object)
	if triple != nil { // EXPECTED NIL BECAUSE ZERO LEN object PASSED
		t.Errorf("did not fail on zero length object: err=%v", err)
	}

	subject = "fido"
	predicate = "isA"
	object = "realdog"
	x := []string{subject, predicate, object}
	y := []string{"subject", "predicate", "object"}
	for i, s := range x {
		hold := x[i]
		x[i] = s + "1"
		triple, err = Dogg3rzTripleNew(x[0], x[1], x[2])
		if triple != nil {
			t.Errorf("did not fail on %s parameter that contained non alpha character: %v", y[i], err)

		}
		x[i] = hold

	}

}
