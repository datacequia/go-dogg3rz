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
	"encoding/json"
	"io"
)

const (
	TYPE_DOGG3RZ_OBJECT             = "dogg3rz.object"
	DOGG3RZ_OBJECT_ATTR_OBJECT_TYPE = "type"
	DOGG3RZ_OBJECT_ATTR_METADATA    = "metadata"
	DOGG3RZ_OBJECT_ATTR_DATA        = "data"
	DOGG3RZ_OBJECT_ATTR_PARENT      = "parent"
)

type dgrzObject struct {
	ObjectType string                 `structs:"type" json:"type"`
	Metadata   map[string]string      `structs:"metadata" json:"metadata"`
	Data       map[string]interface{} `structs:"data,omitempty" json:"data,omitempty"`
	Parents    []string               `structs:"parent,omitempty" json:"parent,omitempty"`
}

type Dogg3rzObjectifiable interface {
	ToDogg3rzObject() *dgrzObject
}

func Dogg3rzObjectDeserializeFromJson(reader io.Reader) (*dgrzObject, error) {

	decoder := json.NewDecoder(reader)

	obj := dgrzObjectNew()

	err := decoder.Decode(obj)
	if err != nil {
		return nil, err
	}

	return obj, err

}

func Dogg3rzObjectSerializeToJson(obj *dgrzObject, writer io.Writer) error {

	encoder := json.NewEncoder(writer)

	err := encoder.Encode(obj)

	return err

}

func Dogg3rzObjectNew(objectType string) *dgrzObject {

	o := dgrzObjectNew()

	o.ObjectType = objectType

	return o

}

func dgrzObjectNew() *dgrzObject {
	return &dgrzObject{ObjectType: "", Metadata: make(map[string]string),
		Data: make(map[string]interface{})}
}
