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
	"strings"
	"testing"
)

func TestDogg3rzObjectNew(t *testing.T) {

	obj := Dogg3rzObjectNew("test")

	if obj.ObjectType != "test" {
		t.Errorf("expected object type '%s', got '%s'", "test", obj.ObjectType)
	}

	if obj.Metadata == nil {
		t.Errorf("expected initialized Metadata attribute. got nil { obj.Metadata = %v}", obj.Metadata)
	}

	if obj.Data == nil {
		t.Errorf("expected initialized Metadata attribute. got nil { obj.Metadata = %v}", obj.Data)
	}

}

func TestDogg3rzObjectSerializeDeserialize(t *testing.T) {

	obj := Dogg3rzObjectNew(TYPE_DOGG3RZ_OBJECT)

	obj.Metadata["attr1"] = "attr_value_1"
	obj.Metadata["attr2"] = "attr_value_2"

	obj.Data["data_attr1"] = 5
	obj.Data["data_attr2"] = "me"

	var s1 strings.Builder

	// SERIALIZE OBJECT TO STRING BUFFER
	err := Dogg3rzObjectSerializeToJson(obj, &s1)
	if err != nil {
		t.Errorf("failed to serialize dogg3rz object to json: { error = %v}", err)
	}

	// DESERIALIZE FROM STRING BUFFER CREATED ABOVE

	reader := strings.NewReader(s1.String())
	obj2, err := Dogg3rzObjectDeserializeFromJson(reader)
	if err != nil {
		t.Errorf("failed to deserialize dogg3rz object from json: { error %v}", err)
	}

	if obj2.ObjectType != obj.ObjectType {
		t.Errorf("expected dogg3rz object type '%s', found '%s'", obj.ObjectType, obj2.ObjectType)
	}

	if obj2.Metadata["attr1"] != "attr_value_1" {
		t.Errorf("expected '%s', got '%s'", "attr_value_1", obj2.Metadata["attr1"])
	}
	if obj2.Metadata["attr2"] != "attr_value_2" {
		t.Errorf("expected '%s', got '%s'", "attr_value_2", obj2.Metadata["attr2"])
	}
	// TEST RECALL OF DATA ATTR
	if val, ok := obj2.Data["data_attr1"].(float64); ok {
		if val != 5 {
			t.Errorf("expected '%d', got '%s'", 5, obj2.Data["data_attr1"])
		}
	} else {
		t.Errorf("expected int type value, got %v", reflect.TypeOf(obj2.Data["data_attr1"]))
	}

	if val, ok := obj2.Data["data_attr2"].(string); ok {
		if val != "me" {
			t.Errorf("expected '%s', got '%s'", "me", obj2.Data["data_attr2"])
		}
	} else {
		t.Errorf("expected string type value, got %v", reflect.TypeOf(obj2.Data["data_attr2"]))
	}

}
