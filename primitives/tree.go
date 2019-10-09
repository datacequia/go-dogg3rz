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
	"reflect"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/fatih/structs"
)

const TYPE_DOGG3RZ_TREE = "dogg3rz.tree"
const D_ATTR_ENTRIES = "entries"

type dgrzTree struct {
	name    string
	entries map[string]interface{} // key = dogg3rz object name, value = hash to dogg3rz object

	parent string
}

func Dogg3rzTreeNew(treeName string) (*dgrzTree, error) {

	if ok, err := errors.PathElementValid(treeName); !ok {
		return nil, errors.InvalidArg.Wrap(err, "treeName")
	}

	return &dgrzTree{name: treeName, entries: make(map[string]interface{})}, nil

}

func (receiver *dgrzTree) PutTree(tree *dgrzTree) error {

	if ok, err := errors.NotNil(receiver); !ok {
		return errors.InvalidArg.Wrap(err, "receiver")
	}
	if ok, err := errors.NotNil(tree); !ok {
		return errors.InvalidArg.Wrap(err, "tree")
	}

	if _, ok := receiver.entries[tree.name]; ok {
		// ANOTHER ENTRY IN THIS TREE WITH SAME NAME ALREADY EXISTS
		return errors.AlreadyExists.New(tree.name)
	}

	receiver.entries[tree.name] = tree

	return nil

}

func (receiver *dgrzTree) GetEntry(name string) (interface{}, error) {

	if ok, err := errors.NotNil(receiver); !ok {
		return nil, errors.InvalidArg.Wrap(err, "receiver")
	}

	if entry, ok := receiver.entries[name]; ok {
		// ENTRY  FOUND IN TREE

		if ok, err := errors.NotNil(entry); !ok {
			// WEIRD. TREE HAD AN ENTRY WHOSE VALUE WAS nil
			panic(errors.Wrapf(err, "dgrzTree entry is nil: { tree name = %s, entry name = %s }", receiver.name, name))

		}

		return entry, nil

	}

	return nil, errors.NotFound.New(name)

}

func (receiver *dgrzTree) Dogg3rzObject() *dgrzObject {

	o := Dogg3rzObjectNew(TYPE_DOGG3RZ_TREE)

	o.Metadata[MD_ATTR_NAME] = receiver.name

	m := make(map[string]interface{})

	for k, v := range receiver.entries {
		if do, ok := v.(Dogg3rzObjectifiable); ok {
			m[k] = structs.Map(do.Dogg3rzObject())
		} else {
			log.Warningf("entry does not implement Dogg3rzObjectifiable: {name: %s, type = %s}", k, reflect.TypeOf(v).String())
		}

	}
	o.Data[D_ATTR_ENTRIES] = m

	return o
}
