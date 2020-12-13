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

func TestRepositoryPathNew(t *testing.T) {

	if _, err := RepositoryPathNew(""); err == nil {
		t.Error("empty path allowed ")
	}

	if rp, _ := RepositoryPathNew("test/"); rp != nil && !rp.lastCharPathElement {
		t.Error("specified path separator as last char in path but not detected")

	}

	if _, err := RepositoryPathNew("/"); err == nil {
		t.Errorf("no error on zero path elements: err = %s", err)

	}

	if _, err := RepositoryPathNew(".dkdkd"); err == nil {
		t.Errorf("no error on path element starting with '.'")

	}

	if _, err := RepositoryPathNew("_dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '_'")

	}

	if _, err := RepositoryPathNew("-dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '-'")

	}

	if _, err := RepositoryPathNew("dkdkd/.dkdkd"); err == nil {
		t.Errorf("no error with path element starting with '.'")

	}

	if rp, _ := RepositoryPathNew("first//////second/"); rp != nil {

		if rp.Size() != 2 {
			t.Errorf("path with empty elements should return 2. returned %d", rp.Size())
		}

	}

	if rp, _ := RepositoryPathNew("first//"); rp != nil {

		if rp.Size() != 1 {
			t.Errorf("path with empty elements should return 1. returned %d", rp.Size())
		}
	}

	if rp, err := RepositoryPathNew("first//second"); err != nil {
		if rp.ToString() != "first/second" {
			t.Errorf("expected 'first/second' found '%s'", rp.ToString())
		}
	}

	if rp, err := RepositoryPathNew("first//second/"); err == nil {
		if rp.ToString() != "first/second/" {
			t.Errorf("expected 'first/second/' found '%s'", rp.ToString())
		}
	}

}
