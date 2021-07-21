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

package repo

import (
	"testing"

	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

func TestInsertNodeIntoGraph(t *testing.T) {

	// SETUP CODE
	indexSetup(t)

	testInsertNodeIntoDefaultGraph(t)
	testInsertNodeIntoNamedGraph(t)

	// TEARDOWN CODE
	indexTeardown(t)

}

// Add new named graph to default graph
func testInsertNodeIntoDefaultGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test1")
	fileDataset.create(ctxt)

	keys := []string{"key1"}
	values := []string{"value1"}
	var fileRe FileRepositoryResource
	err := fileRe.InsertNode(ctxt, testRepoName, "test1", "", "", "", keys, values)

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

	mtime := m[jsonld.MtimesEntryKeyName]
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated : %v", mtime)
		t.FailNow()
	}
}

// Add new named graph to Named graph
func testInsertNodeIntoNamedGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test2")
	fileDataset.create(ctxt)
	var parentGraph = "testParent2"
	err1 := fileDataset.createNamedGraph(ctxt, parentGraph, "default")
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
		t.FailNow()
	}
	keys := []string{"key2"}
	values := []string{"value2"}
	var fileRe FileRepositoryResource
	err := fileRe.InsertNode(ctxt, testRepoName, "test2", "", "", parentGraph, keys, values)

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

	childnodeAsMap, _ := nodeAsMap["@graph"].([]interface{})[0].(map[string]interface{})

	if childnodeAsMap["@id"] == nil {
		t.Errorf("Id not added for new property  : %v", childnodeAsMap)
		t.FailNow()
	}

	if childnodeAsMap["key2"] == nil || childnodeAsMap["key2"] != "value2" {
		t.Errorf("New attributes not added to the Named graph  : %v", childnodeAsMap)
		t.FailNow()
	}

	mtime := m[jsonld.MtimesEntryKeyName]
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated : %v", mtime)
		t.FailNow()
	}
}
