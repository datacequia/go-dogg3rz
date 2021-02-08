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

// package jsonld provides general use constants/structs/funcs that are
// used when interacting with JSON-LD documents
package jsonld

import (
	"fmt"

	"github.com/datacequia/go-dogg3rz/errors"
)

type JSONLDResourceType int8

const (
	DatasetResource    JSONLDResourceType = 1
	ContextResource    JSONLDResourceType = 2
	NamedGraphResource JSONLDResourceType = 3
	NodeResource       JSONLDResourceType = 4
	// NOTE: IF ANOTHER RESOURCE CONSTANT IS ADDED, PLEASE ADD TO
	// validValues ARRAY OR WILL FAIL ASSERTION METHOD BELOW
)

// MtimesKeyName is used as an outermost entry key in a JSON-LD DOCUMENT
// which assumes a map as a value that stores mtimes as entries where
// the key is a SHA-256 output of a JSONLDDocumentLocation and value is
// the epoc time of the last time the entry was modified
const MtimesEntryKeyName = "__dgrz_mtimes__"

var validValues []JSONLDResourceType = []JSONLDResourceType{DatasetResource, ContextResource, NamedGraphResource, NodeResource}

func (x JSONLDResourceType) AssertValid() error {

	for _, v := range validValues {
		if x == v {
			return nil
		}
	}

	return errors.UnexpectedValue.Newf("Unexpected JSONLDResourceType value. Found %d, want one of %v", x, validValues)
}
func (x JSONLDResourceType) String() string {

	switch x {
	case DatasetResource:
		return "Dataset"
	case ContextResource:
		return "Context"
	case NamedGraphResource:
		return "NamedGraph"
	case NodeResource:
		return "Node"
	default:
		return fmt.Sprintf("<bad JSONLDResourceType %d>", x)

	}

}

// isContiner returns true of the resourcetype can contain other resources
func (x JSONLDResourceType) IsContainer() bool {

	switch x {
	case DatasetResource, NamedGraphResource, NodeResource:
		return true
	}
	return false

}
