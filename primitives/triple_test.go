/*
 * Copyright (c) 2019-2020 Datacequia LLC. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
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
