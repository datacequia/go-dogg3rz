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
	"github.com/datacequia/go-dogg3rz/errors"
)

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
		return nil, errors.InvalidArg.Wrap(err, "subject")
	}

	if ok, err := errors.StrlenGtZero(predicate); !ok {
		return nil, errors.InvalidArg.Wrap(err, "predicate")
	}

	if ok, err := errors.StrlenGtZero(object); !ok {
		return nil, errors.InvalidArg.Wrap(err, "object")
	}

	if ok, err := errors.StrAlpha(subject); !ok {
		return nil, errors.InvalidArg.Wrap(err, "subject")
	}
	if ok, err := errors.StrAlpha(predicate); !ok {
		return nil, errors.InvalidArg.Wrap(err, "predicate")
	}

	if ok, err := errors.StrAlpha(object); !ok {
		return nil, errors.InvalidArg.Wrap(err, "object")
	}

	triple := &dgrzTriple{subject: subject, predicate: predicate, object: object}

	return triple, nil

}
