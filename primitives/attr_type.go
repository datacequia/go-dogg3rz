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
	"fmt"

	"github.com/datacequia/go-dogg3rz/errors"
)

// INTERNAL DOGG3RZ ATTR TYPE REPRESENTED
// AS A NUMBER FOR EACH PRIMITIVE TYPE
type Dogg3rzObjectType uint32

var dogg3rzObjectTypeMap = map[Dogg3rzObjectType]string{

	TYPE_DOGG3RZ_MEDIA:   "dogg3rz.media",
	TYPE_DOGG3RZ_OBJECT:  "dogg3rz.object",
	TYPE_DOGG3RZ_SCHEMA:  "dogg3rz.schema",
	TYPE_DOGG3RZ_SERVICE: "dogg3rz.service",
	TYPE_DOGG3RZ_TREE:    "dogg3rz.tree",
	TYPE_DOGG3RZ_TRIPLE:  "dogg3rz.triple",
}

var dogg3rzObjectTypeSMap map[string]Dogg3rzObjectType

var dogg3rzObjectTypes []Dogg3rzObjectType

var dogg3rzObjectTypesS []string

func init() {
	dogg3rzObjectTypeSMap = make(map[string]Dogg3rzObjectType)

	dogg3rzObjectTypes = make([]Dogg3rzObjectType, len(dogg3rzObjectTypeMap))
	dogg3rzObjectTypesS = make([]string, len(dogg3rzObjectTypeMap))

	var i = 0
	for k, v := range dogg3rzObjectTypeMap {
		dogg3rzObjectTypeSMap[v] = k
		dogg3rzObjectTypes[i] = k
		dogg3rzObjectTypesS[i] = v
		i++
	}
}

func (t Dogg3rzObjectType) String() string {

	if value, ok := dogg3rzObjectTypeMap[t]; ok {
		return value
	}

	panic(fmt.Sprintf("Dogg3rzObjectType.String(): unmapped type encountered: %d", t))

}

func (t Dogg3rzObjectType) Valid() bool {

	_, ok := dogg3rzObjectTypeMap[t]

	return ok
}

func Dogg3rzObjectTypes() []Dogg3rzObjectType {

	return dogg3rzObjectTypes

}

func Dogg3rzObjectTypeFromString(s string) (Dogg3rzObjectType, error) {
	//fmt.Println("dogg3rzObjectTypeSMap", dogg3rzObjectTypeSMap)
	if t, ok := dogg3rzObjectTypeSMap[s]; ok {
		return t, nil
	} else {
		return t, errors.NotFound.Newf(
			"Dogg3rzObjectTypeFromString(): '%s': expected one of %v",
			s,
			dogg3rzObjectTypes)

	}
}

/*
type Dogg3rzObjectType struct {
	value dogg3rzObjectType
}

// CREATE NEW INSTANCE
func NewDogg3rzObjectType(ot dogg3rzObjectType) Dogg3rzObjectType {
	return Dogg3rzObjectType{value: ot}
}

// CREATE NEW INSTANCE USING STRING REPRESENTATION
// AS INPUT
func NewDogg3rzObjectTypeS(s string) Dogg3rzObjectType {

	d := Dogg3rzObjectType{}

	d.value = mapStringToType(s)

	return d

}

// CONVERT TO NUMERIC VALUE
func (t Dogg3rzObjectType) ToUint32() uint32 {

	_ = mapTypeToString(t.value) // WILL PANIC IF TYPE NOT INITIALIZED

	return uint32(t.value)
}

// CONVERT INSTANCE TO STRING REPRESENTATION
func (t Dogg3rzObjectType) String() string {

	return mapTypeToString(t.value)

}

func (t Dogg3rzObjectType) Value() dogg3rzObjectType {
	return t.value

}

// COMPARE INSTANCE WITH  ANOTHER INTERNAL Type
// NOTE VALUES FOR THESE TYPEES ARE EXPORTED FROM THIS
// PACKAGE BUT NOT CREATBLE/ASSIGNABLE
func (l Dogg3rzObjectType) Equal(r dogg3rzObjectType) bool {

	return l.value == r

}

func mapTypeToString(t dogg3rzObjectType) string {

	var s string

	switch t {
	case TYPE_DOGG3RZ_MEDIA:
		s = "dogg3rz.media"
	case TYPE_DOGG3RZ_OBJECT:
		s = "dogg3rz.object"
	case TYPE_DOGG3RZ_SCHEMA:
		s = "dogg3rz.schema"
	case TYPE_DOGG3RZ_SERVICE:
		s = "dogg3rz.service"
	case TYPE_DOGG3RZ_TREE:
		s = "dogg3rz.tree"
	case TYPE_DOGG3RZ_TRIPLE:
		s = "dogg3rz.triple"
	default:
		if t == 0 {
			log.Panicf("mapTypeToString(): ininitialized  Dogg3rzObjectType: { value = %d }", t)
		}
		log.Panicf("mapTypeToString(): unmapped Dogg3rzObjectType: { value = %d }", t)
	}

	return s

}

func mapStringToType(s string) dogg3rzObjectType {

	var d dogg3rzObjectType

	switch s {

	case "dogg3rz.media":
		d = TYPE_DOGG3RZ_MEDIA
	case "dogg3rz.object":
		d = TYPE_DOGG3RZ_OBJECT
	case "dogg3rz.schema":
		d = TYPE_DOGG3RZ_SCHEMA
	case "dogg3rz.service":
		d = TYPE_DOGG3RZ_SERVICE
	case "dogg3rz.tree":
		d = TYPE_DOGG3RZ_TREE
	case "dogg3rz.triple":
		d = TYPE_DOGG3RZ_TRIPLE

	default:

		log.Panicf("mapStringToType(): unmapped Dogg3rzObjectType: '%s'", s)

	}

	return d
}
*/
