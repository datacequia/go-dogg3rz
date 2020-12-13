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

package jsonld

import (
	"testing"
)

func TestJSONLDResourceType(t *testing.T) {

	for _, x := range validValues {

		if err := x.AssertValid(); err != nil {
			t.Errorf("Valid JSONLDResourceType failed assertion: %v", x)
		}
	}

	var badVal JSONLDResourceType = 99

	if err := badVal.AssertValid(); err == nil {
		t.Errorf("Invalid JSONLDResourceType failed assertion: %v", badVal)
	}

}
