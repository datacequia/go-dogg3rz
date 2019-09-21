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
