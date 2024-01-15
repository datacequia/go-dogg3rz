//go:build nothing

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

package common

import (
	"testing"
)

func TestGrapplicationPathNew(t *testing.T) {

	if _, err := GrapplicationPathNew(""); err == nil {
		t.Error("empty path allowed ")
	}

	if rp, _ := GrapplicationPathNew("test/"); rp != nil && !rp.lastCharPathElement {
		t.Error("specified path separator as last char in path but not detected")

	}

	if _, err := GrapplicationPathNew("/"); err == nil {
		t.Errorf("no error on zero path elements: err = %s", err)

	}

	if _, err := GrapplicationPathNew(".dkdkd"); err == nil {
		t.Errorf("no error on path element starting with '.'")

	}

	if _, err := GrapplicationPathNew("_dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '_'")

	}

	if _, err := GrapplicationPathNew("-dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '-'")

	}

	if _, err := GrapplicationPathNew("dkdkd/.dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '.'")

	}

	if rp, _ := GrapplicationPathNew("first//////second/"); rp != nil {

		if rp.Size() != 2 {
			t.Errorf("path with empty elements should return 2. returned %d", rp.Size())
		}

	}

	if rp, _ := GrapplicationPathNew("first//"); rp != nil {

		if rp.Size() != 1 {
			t.Errorf("path with empty elements should return 1. returned %d", rp.Size())
		}
	}

	if rp, err := GrapplicationPathNew("first//second"); err != nil {
		if rp.ToString() != "first/second" {
			t.Errorf("expected 'first/second' found '%s'", rp.ToString())
		}
	}

	if rp, err := GrapplicationPathNew("first//second/"); err == nil {
		if rp.ToString() != "first/second/" {
			t.Errorf("expected 'first/second/' found '%s'", rp.ToString())
		}
	}

}
