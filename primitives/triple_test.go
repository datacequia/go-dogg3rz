/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
