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

package common

import (
	"context"
	"reflect"
	"testing"

	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

type testStagingResourceCollector struct {
	locations []StagingResourceLocation
	resources []interface{}
	ctxt      context.Context
}

func (o *testStagingResourceCollector) CollectStart(ctxt context.Context, resource interface{}, location StagingResourceLocation) error {

	o.locations = append(o.locations, location)
	o.resources = append(o.resources, resource)

	return nil

}

func (o *testStagingResourceCollector) CollectEnd(ctxt context.Context, resource interface{}, location StagingResourceLocation) {
	// NOTHIING TO CLEANUP FOR THIS IMPL.

}

func TestAssertValids(t *testing.T) {

	var sr StagingResource = StagingResource{}

	var err error

	sr.ObjectType = 0
	if err = sr.AssertValid(); err == nil {
		t.Errorf("Expected error on bad %T.ObjctType", sr)
	}
	sr.ObjectType = jsonld.NodeResource
	sr.ContainerType = 0
	if err = sr.AssertValid(); err == nil {
		t.Errorf("Expected error on bad %T.ContainerType", sr)
	}

	// test with values that cannot be staged
	sr.ObjectType = jsonld.ContextResource
	sr.ContainerType = jsonld.ContextResource

	if err = sr.AssertValid(); err == nil {

		t.Errorf("Expected error on unstagable object/container type combination: %v", err)
	}

	// test with values that Can be   staged
	// but wants CID but it's not initialized
	sr.ObjectType = jsonld.NodeResource
	sr.ContainerType = jsonld.DatasetResource
	sr.ObjectCID = ""
	if err = sr.AssertValid(); err == nil {
		t.Errorf("Expected error on stageable resource that should have it's CID populated")
	}

	// test with values that Can be   staged
	// but does not want  CID but it's  initialized
	sr.ObjectType = jsonld.NamedGraphResource
	sr.ContainerType = jsonld.DatasetResource
	sr.ObjectCID = "QmQHYbaGrvwAW1u4YWohKbyp4Uj23pG1yZaxPJX13zT6HZ"
	if err = sr.AssertValid(); err == nil {
		t.Errorf("Expected error on stageable resource that should have it's CID populated")
	}

}

func TestFindStageableResources(t *testing.T) {

	//	var sr = StagingResourceLocation{}

	var c = testStagingResourceCollector{}

	var doc = map[string]interface{}{
		//		"a": "b",
		"@context": map[string]interface{}{

			"abc": "http://www.abc.com/#",
			"def": "http://www.def.com/#",
		},
		"@graph": []interface{}{

			map[string]interface{}{
				"@id": "node1",
			},
			map[string]interface{}{
				"@id": "namedGraph1",
				"@graph": []interface{}{
					map[string]interface{}{
						"@context": map[string]interface{}{},
						"@id":      "node2",
					},
				},
			},
			map[string]interface{}{
				"@id": "node3",
			},
		},
	}

	const datasetPath = "a/b/c"

	if err := FindStageableResources(context.Background(), datasetPath, doc, &c); err != nil {

		t.Errorf("FindStageableResources failed on good doc: %s", err)
	}

	var ctxtResCnt int
	var ngResCnt int
	var nodeResCnt int

	for i, _ := range c.locations {

		var loc = c.locations[i]
		//var res = c.resources[i]

		switch loc.ObjectType {

		case jsonld.ContextResource:

			ctxtResCnt++

			if loc.ContainerType != jsonld.DatasetResource {
				t.Errorf("expected %s, found %s", jsonld.DatasetResource, loc.ContainerType)
			}
			if loc.ContainerIRI != "" {
				t.Errorf("expected blank IRI for DatasetResource container, found '%s'",
					loc.ContainerIRI)
			}
			if loc.ObjectIRI != "" {
				t.Errorf("expected blank IRI for ContextResource object, found '%s'",
					loc.ContainerIRI)
			}
			if loc.DatasetPath != datasetPath {
				t.Errorf("expected '%s' for .DatasetPath, found %s", datasetPath, loc.DatasetPath)
			}

		case jsonld.NamedGraphResource:
			ngResCnt++

			if loc.ObjectIRI != "namedGraph1" {
				t.Errorf("expected '%s' for .ObjectIRI where .ObjectType=jsonld.NamedGraphResource, found %s", "namedGraph1", loc.ObjectIRI)

			}
			if loc.ContainerType != jsonld.DatasetResource {
				t.Errorf("expected .ContainerType to be jsonld.DatasetResource, found %s", loc.ContainerType)
			}
			if loc.ContainerIRI != "" {
				t.Errorf("expected blank IRI for DatasetResource container, found '%s'",
					loc.ContainerIRI)
			}
			if loc.DatasetPath != datasetPath {
				t.Errorf("expected '%s' for .DatasetPath, found %s", datasetPath, loc.DatasetPath)
			}

		case jsonld.NodeResource:
			nodeResCnt++

			switch loc.ObjectIRI {
			case "node1":
				if loc.ContainerType != jsonld.DatasetResource {
					t.Errorf("expected %s, found %s", jsonld.DatasetResource, loc.ContainerType)
				}
				if loc.ContainerIRI != "" {
					t.Errorf("expected blank IRI for DatasetResource container, found '%s'",
						loc.ContainerIRI)
				}

			case "node2":
				if loc.ContainerType != jsonld.NamedGraphResource {
					t.Errorf("expected %s, found %s", jsonld.DatasetResource, loc.ContainerType)
				}
				if loc.ContainerIRI != "namedGraph1" {
					t.Errorf("expected '%s' IRI for DatasetResource container, found '%s'",
						"namedGraph1", loc.ContainerIRI)
				}

			case "node3":
				if loc.ContainerType != jsonld.DatasetResource {
					t.Errorf("expected %s, found %s", jsonld.DatasetResource, loc.ContainerType)
				}
				if loc.ContainerIRI != "" {
					t.Errorf("expected blank IRI for DatasetResource container, found '%s'",
						loc.ContainerIRI)
				}
			default:
				t.Errorf("expected 'node1','node2', or 'node3' for .ObjectIRI  found %s", loc.ObjectIRI)

			}

			if loc.DatasetPath != datasetPath {
				t.Errorf("expected '%s' for .DatasetPath, found %s", datasetPath, loc.DatasetPath)
			}

		}

	}

	if nodeResCnt != 3 {
		t.Errorf("expected total nodeResCnt to be '3', found '%d'", nodeResCnt)
	}
	if ngResCnt != 1 {
		t.Errorf("expected total nodeResCnt to be '1', found '%d'", nodeResCnt)
	}
	if ctxtResCnt != 1 {
		t.Errorf("expected total nodeResCnt to be '1', found '%d'", ctxtResCnt)
	}

}

func TestHasFuncs(t *testing.T) {

	m := map[string]interface{}{
		"key1": "value1",
	}

	if v, ok := hasEntryKey("key1", m); ok {
		if v != "value1" {
			t.Errorf("expected value '%s', found %s", "value1", v)
		}
	} else {
		t.Errorf("expected return value true, got %v", ok)
	}

	//hasEntryWithMapValue(entryKey string, obj map[string]interface{}) (map[string]interface{}, bool)
	m2 := map[string]interface{}{
		"key1": map[string]interface{}{
			"key2": "val2",
		},
	}

	if mVal, ok := hasEntryWithMapValue("key1", m2); !ok {

		t.Errorf("expected return value true, got %v", ok)
	} else {
		if val, ok := mVal["key2"]; ok {
			if val != "val2" {
				t.Errorf("expected 'val2', got '%s'", val)
			}
		} else {
			t.Errorf("expected map value with key 'key2', not found")
		}
	}

	if _, ok := hasEntryWithMapValue("key1", m); ok {
		t.Errorf("expected return value false, got %v", ok)
	}

	//
	m2a := map[string]interface{}{
		"key1": []interface{}{
			"val1", "val2",
		},
	}

	if ar, ok := hasEntryWithArrayValue("key1", m2a); !ok {

		t.Errorf("expected return value true, got %v", ok)
	} else {

		if len(ar) != 2 {
			t.Errorf("Expected array with 2 elements, found %d", len(ar))
		}

		if ar[0] != "val1" {

			t.Errorf("expected 'val2', got '%s'", ar[1])

		}

		if ar[1] != "val2" {

			t.Errorf("expected 'val2', got '%s'", ar[1])

		}

	}

}

func TestIsNodeObject(t *testing.T) {

	//isNodeObject(obj map[string]interface{}, existsOutsideContextObject bool, topMostMapInJSONLDDoc bool) bool

	node1 := make(map[string]interface{})

	// check that if caller says existsOutsideContextObject is false, it returns
	// false
	if isNodeObject(node1, false, false) {
		t.Errorf("expected return false, got true")
	}

	// ADD @value,@list,@set keywords and ensure it's not a node

	// BELOW SETS UP NODE SO THATA IT LOOKS LIKE A TOP LEVEL NODE
	// WHICH HAS BOTH @graph and @context
	//  this isolates
	node1["@graph"] = "ddd"
	node1["@context"] = "dkdk"

	for _, kw := range []string{"@value", "@list", "@set"} {

		//fmt.Println(kw)
		node1[kw] = "val"

		if isNodeObject(node1, true, true) {
			t.Errorf("expected return false, got true when passing invalid keywords in a node object")
		}

		delete(node1, kw)

	}

}

func TestIsGraphObject(t *testing.T) {

	graph1 := make(map[string]interface{})

	graph1["@graph"] = "dummy"

	var existsOutsideContextObject = true
	var topMostMapInJSONLDDoc = false
	if ok, _ := isGraphObject(graph1, existsOutsideContextObject, topMostMapInJSONLDDoc); !ok {
		t.Errorf("expected object with only @graph keyword that exists outside context object and no topmost map to be a graph object")

	}

	existsOutsideContextObject = false
	if ok, _ := isGraphObject(graph1, existsOutsideContextObject, topMostMapInJSONLDDoc); ok {
		t.Errorf("expected object with only @graph keyword that exists inside  context object to be not a graph object")

	}

	existsOutsideContextObject = true
	topMostMapInJSONLDDoc = true
	if ok, _ := isGraphObject(graph1, existsOutsideContextObject, topMostMapInJSONLDDoc); ok {
		t.Errorf("expected object with only @graph keyword that is topmost map in JSON document to be not a graph object")

	}

}

func TestIsContextValue(t *testing.T) {

	// isContextValue(v interface{}) bool

	//var cv interface{}

	// test nil as a valid context value
	if !isContextValue(nil) {
		t.Errorf("Expected type %T to be valid context value type, got false", nil)
	}

	// test string
	if !isContextValue("http://www.x.com/#") {
		t.Errorf("Expected type %T to be valid context value type, got false", nil)
	}

	// test map
	if !isContextValue(make(map[string]interface{})) {
		t.Errorf("Expected type %T to be valid context value type, got false", nil)
	}

	// test array
	//	x := []string{"http://www.x.com/#"}
	if !isContextValue([1]string{"http://www.x.com/#"}) {
		t.Errorf("Expected type []string to be valid context value type, got false")
	}

	// test slice (of string )
	if !isContextValue([]string{"http://www.x.com/#"}) {
		t.Errorf("Expected type []string to be valid context value type, got false")
	}

	// test nil
	if !isContextValue(nil) {
		t.Errorf("Expected type nil to be valid context value type, got false")
	}

	//  TYPES NOT ACCEPTEED FOR A CONTEXT VALUE.
	// NOT EXCLUSIVE. IF IT'S NOT ABOVE TYPES IT SHOULD FAIL

	// test int (should be bad type)
	if isContextValue(123) {
		t.Errorf("Expected type map to be valid context value type, got false")
	}
	// test float (should be bad type)
	if isContextValue(123.444) {
		t.Errorf("Expected type map to be valid context value type, got false")
	}

	// test chan (shoould be bad type)
	if isContextValue(make(chan int)) {
		t.Errorf("Expected type map to be valid context value type, got false")
	}

}

func TestGenerateSHA256Key_TypeCast(t *testing.T) {

	// TEST THAT jsonld.JSONLDResourceType size equal to sizeof byte
	// because it is casted to byte in this function

	var x jsonld.JSONLDResourceType
	var y byte

	if reflect.TypeOf(x).Size() > reflect.TypeOf(y).Size() {
		t.Errorf("Generate256Key may fail (overflow error) when jsonld.ResourceType is cast to byte")
	}

}

func TestGenerateSHA256Key(t *testing.T) {

	var a JSONLDDocumentLocation
	var b JSONLDDocumentLocation

	a.ObjectType = jsonld.ContextResource
	a.ObjectIRI = ""
	a.ContainerType = jsonld.DatasetResource
	a.ContainerIRI = ""

	b = a

	if key1, err := a.GenerateSHA256Key(); err != nil {
		t.Errorf("JSONLDDocumentLocation.Generate256Key() = %s", err)
	} else {
		if key2, err := b.GenerateSHA256Key(); err != nil {
			t.Errorf("JSONLDDocumentLocation.Generate256Key() = %s", err)
		} else {
			if key1 != key2 {
				t.Errorf("JSONLDDocumentation sha256 keys for same value are not equal: { a = %v, b = %v, a(key) = %s, b(key) = %s }",
					a, b, key1, key2)
			}
		}
	}

	a.ObjectType = jsonld.NodeResource
	a.ObjectIRI = "http://www.dummy.com/#myobj"
	a.ContainerType = jsonld.NamedGraphResource
	a.ContainerIRI = "http://www.dummy.com/#mynamedgraph"

	b = a

	if key1, err := a.GenerateSHA256Key(); err != nil {
		t.Errorf("JSONLDDocumentLocation.Generate256Key() = %s", err)
	} else {
		if key2, err := b.GenerateSHA256Key(); err != nil {
			t.Errorf("JSONLDDocumentLocation.Generate256Key() = %s", err)
		} else {
			if key1 != key2 {
				t.Errorf("JSONLDDocumentation sha256 keys for same value are not equal: { a = %v, b = %v, a(key) = %s, b(key) = %s }",
					a, b, key1, key2)
			}
		}
	}

}
