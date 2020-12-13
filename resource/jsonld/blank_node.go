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

	"github.com/google/uuid"
)

const BlankNodeIDPrefix = "_"

type BlankNodeID struct {
	id uuid.UUID
}

func NewBlankNodeID() BlankNodeID {

	blankNode := BlankNodeID{}

	blankNode.id = uuid.New()

	return blankNode
}

func (b BlankNodeID) String() string {

	return fmt.Sprintf("%s:%s", BlankNodeIDPrefix, b.id.String())

}
