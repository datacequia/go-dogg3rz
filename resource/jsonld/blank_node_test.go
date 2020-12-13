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
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestNotBlankNodeID(t *testing.T) {

	blankNode := NewBlankNodeID()

	//fmt.Println("blankNode", blankNode)

	blankNodeStr := fmt.Sprintf("%s", blankNode)

	if blankNodeStr[0:1] != BlankNodeIDPrefix {
		t.Errorf("generated blank node ID %s does not start with blank node prefix '%s'",
			blankNodeStr, BlankNodeIDPrefix)
	}

	if blankNodeStr[1:2] != ":" {
		t.Errorf("generated blank node ID %s contain colon separator following prefix '%s'",
			blankNodeStr, BlankNodeIDPrefix)
	}

	if _, err := uuid.Parse(blankNodeStr[2:]); err != nil {
		t.Errorf("generated blank node ID  identifier  %s does not contain a UUID  as a suffix: %s",
			blankNodeStr, err)
	}

}
