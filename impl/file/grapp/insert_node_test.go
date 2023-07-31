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

package grapp

import (
	"testing"

	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

func TestInsertNodeIntoGraph(t *testing.T) {

	// SETUP CODE
	indexSetup(t)

	testInsertNodeIntoDefaultGraph(t)
	testInsertNodeIntoNamedGraph(t)
	testInsertNodeIntoNGAfterDefaultGraph(t)
	// TEARDOWN CODE
	indexTeardown(t)

}

// Add new named graph to default graph
func testInsertNodeIntoDefaultGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testGrappName, "test1")
	fileDataset.create(ctxt)

	keys := []string{"key1"}
	values := []string{"value1"}
	var fileRe FileGrapplicationResource
	err := fileRe.InsertNode(ctxt, testGrappName, "test1", "", "", "", keys, values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	m, defaultGraph, _ := fileDataset.readDefaultGraph()
	if defaultGraph == nil || len(defaultGraph) == 0 {
		t.Errorf("Attributes not added to default graph : %v", defaultGraph)
		t.FailNow()
	}

	m2 := defaultGraph[0]
	nodeAsMap, _ := m2.(map[string]interface{})
	if nodeAsMap["@id"] == nil {
		t.Errorf("Id not added for new property  : %v", defaultGraph)
		t.FailNow()
	}

	if nodeAsMap["key1"] == nil || nodeAsMap["key1"] != "value1" {
		t.Errorf("New attributes not added to the default graph  : %v", defaultGraph)
		t.FailNow()
	}

	mtime := m[jsonld.MtimesEntryKeyName].(map[string]interface{})
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated expected size of MTime to be 3 but was  : %v", mtime)
		t.FailNow()
	}
	if len(mtime) != 3 {
		t.Errorf("Mtime node not updated correctly: %v", mtime)
		t.FailNow()
	}

}

// Add new named graph to Named graph
func testInsertNodeIntoNamedGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testGrappName, "test2")
	fileDataset.create(ctxt)
	var namedGraphName = "namedGraph1"
	err1 := fileDataset.createNamedGraph(ctxt, namedGraphName, "default")
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
		t.FailNow()
	}
	keys := []string{"key2"}
	values := []string{"value2"}
	var fileRe FileGrapplicationResource
	err := fileRe.InsertNode(ctxt, testGrappName, "test2", "", "", namedGraphName, keys, values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	m, defaultGraph, _ := fileDataset.readDefaultGraph()
	if defaultGraph == nil || len(defaultGraph) == 0 {
		t.Errorf(" default graph not found: %v", defaultGraph)
		t.FailNow()
	}

	m2 := defaultGraph[0]
	nodeAsMap, _ := m2.(map[string]interface{})

	childnodeAsMap, _ := nodeAsMap["@graph"].([]interface{})[0].(map[string]interface{})

	if childnodeAsMap["@id"] == nil {
		t.Errorf("Id not added for new property  : %v", childnodeAsMap)
		t.FailNow()
	}

	if childnodeAsMap["key2"] == nil || childnodeAsMap["key2"] != "value2" {
		t.Errorf("New attributes not added to the Named graph  : %v", childnodeAsMap)
		t.FailNow()
	}

	mtime := m[jsonld.MtimesEntryKeyName].(map[string]interface{})
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated expected size of MTime to be 3 but was  : %v", mtime)
		t.FailNow()
	}
	if len(mtime) != 4 { // 2 dataset creating, 1 for named graph node and 1 for adding properties to names graph
		t.Errorf("Mtime node not updated correctly: %v", mtime)
		t.FailNow()
	}
}

// Add new named graph to default graph
func testInsertNodeIntoNGAfterDefaultGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testGrappName, "test3")
	fileDataset.create(ctxt)

	var namedGraphName = "namedGraph2"
	err1 := fileDataset.createNamedGraph(ctxt, namedGraphName, "default")
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
		t.FailNow()
	}
	keys := []string{"key4"}
	values := []string{"value4"}
	var fileRe FileGrapplicationResource
	err := fileRe.InsertNode(ctxt, testGrappName, "test3", "", "", "", keys, values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	m, defaultGraph, _ := fileDataset.readDefaultGraph()
	if defaultGraph == nil || len(defaultGraph) == 0 {
		t.Errorf("Attributes not added to default graph : %v", defaultGraph)
		t.FailNow()
	}

	m2 := defaultGraph[1]
	nodeAsMap, _ := m2.(map[string]interface{})
	if nodeAsMap["@id"] == nil {
		t.Errorf("Id not added for new property  : %v", defaultGraph)
		t.FailNow()
	}

	if nodeAsMap["key4"] == nil || nodeAsMap["key4"] != "value4" {
		t.Errorf("New attributes not added to the default graph  : %v", defaultGraph)
		t.FailNow()
	}

	mtime := m[jsonld.MtimesEntryKeyName]
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated : %v", mtime)
		t.FailNow()
	}

	keys = []string{"key3"}
	values = []string{"value3"}

	err = fileRe.InsertNode(ctxt, testGrappName, "test3", "", "", namedGraphName, keys, values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	m, defaultGraph, _ = fileDataset.readDefaultGraph()
	if defaultGraph == nil || len(defaultGraph) == 0 {
		t.Errorf(" default graph not found: %v", defaultGraph)
		t.FailNow()
	}

	m2 = defaultGraph[0]
	nodeAsMap, _ = m2.(map[string]interface{})

	childnodeAsMap, _ := nodeAsMap["@graph"].([]interface{})[0].(map[string]interface{})

	if childnodeAsMap["@id"] == nil {
		t.Errorf("Id not added for new property  : %v", childnodeAsMap)
		t.FailNow()
	}

	if childnodeAsMap["key3"] == nil || childnodeAsMap["key3"] != "value3" {
		t.Errorf("New attributes not added to the Named graph  : %v", childnodeAsMap)
		t.FailNow()
	}

	mtime = m[jsonld.MtimesEntryKeyName]
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated : %v", mtime)
		t.FailNow()
	}
}
