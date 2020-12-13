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
	"github.com/datacequia/go-dogg3rz/errors"
)

const TYPE_DOGG3RZ_TRIPLE Dogg3rzObjectType = 1 << 6

const ATTR_SUBJECT = "subject"
const ATTR_PREDICATE = "predicate"
const ATTR_OBJECT = "object"

var reservedAttributesTriple = [...]string{ATTR_SUBJECT, ATTR_PREDICATE, ATTR_OBJECT}

// See https://en.wikipedia.org/wiki/Semantic_triple
type dgrzTriple struct {
	subject   string // needs to be an ipfs cid that points todogg3rz primitive
	predicate string // some
	object    string // needs to be an ipfs cid that points to a dogg3rz primitive

	parent string
}

func Dogg3rzTripleNew(subject string, predicate string, object string) (*dgrzTriple, error) {

	// TODO: validate args 'subject' and 'object' to check if they are valid
	// CIDs using IPFS API
	if ok, err := errors.StrlenGtZero(subject); !ok {
		return nil, errors.InvalidValue.Wrap(err, "subject")
	}

	if ok, err := errors.StrlenGtZero(predicate); !ok {
		return nil, errors.InvalidValue.Wrap(err, "predicate")
	}

	if ok, err := errors.StrlenGtZero(object); !ok {
		return nil, errors.InvalidValue.Wrap(err, "object")
	}

	if ok, err := errors.StrAlpha(subject); !ok {
		return nil, errors.InvalidValue.Wrap(err, "subject")
	}
	if ok, err := errors.StrAlpha(predicate); !ok {
		return nil, errors.InvalidValue.Wrap(err, "predicate")
	}

	if ok, err := errors.StrAlpha(object); !ok {
		return nil, errors.InvalidValue.Wrap(err, "object")
	}

	triple := &dgrzTriple{subject: subject, predicate: predicate, object: object}

	return triple, nil

}
