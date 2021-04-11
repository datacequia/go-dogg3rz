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

	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

func TestCreateNamedGraph(t *testing.T) {

	// SETUP CODE
	indexSetup(t)

	testCreateNamedGraphInDefaultGraph(t)
	testCreateNamedGraphInAnotherNamedGraph(t)
	testDuplicateCreateNamedGraphInDefaultGraph(t)
	testDuplicateCreateNamedGraphInAnotherNamedGraph(t)
	// TEARDOWN CODE
	indexTeardown(t)

}

// Add new named graph to default graph
func testCreateNamedGraphInDefaultGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test1")
	fileDataset.create(ctxt)
	mOriginal, _, _ := fileDataset.readDefaultGraph()
	var childGraph = "test1"
	err := fileDataset.createNamedGraph(ctxt, childGraph, "default")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	m, defaultGraph, _ := fileDataset.readDefaultGraph()
	createdGraph, err1 := common.GetGraph(defaultGraph, childGraph)
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
		t.FailNow()
	}
	if createdGraph == nil {
		t.Errorf("Newly created graph not found: %v", childGraph)
		t.FailNow()
	}
	mtime := m[jsonld.MtimesEntryKeyName]
	if mtime == nil {
		t.Errorf("Newly created graph Mtime not updated : %v", mtime)
		t.FailNow()
	}
}

// Add new named graph to another named graph
func testCreateNamedGraphInAnotherNamedGraph(t *testing.T) {
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test2")
	fileDataset.create(ctxt)
	var parentGraph = "testParent2"
	var childGraph = "childGraphTest2"
	err1 := fileDataset.createNamedGraph(ctxt, parentGraph, "default")
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
		t.FailNow()
	}

	err := fileDataset.createNamedGraph(ctxt, childGraph, parentGraph)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
	_, defaultGraph, _ := fileDataset.readDefaultGraph()
	createdGraph, err1 := common.GetGraph(defaultGraph, childGraph)
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)

		t.FailNow()
	}
	if createdGraph == nil {
		t.Errorf("Newly created graph not found: %v", childGraph)
		t.FailNow()
	}
}

// Add duplicate named graph to default graph
func testDuplicateCreateNamedGraphInDefaultGraph(t *testing.T) {

	// TEST NEW REPO INDEX WITH BAD REPO NAME
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test2")
	fileDataset.create(ctxt)

	var childGraph = "childGraphTest3"
	err := fileDataset.createNamedGraph(ctxt, childGraph, "default")
	if err != nil {
		t.Errorf("Unexpected error happened cannot add named graph: %v, %v", childGraph, err)
		t.FailNow()
	}
	// Try adding graph with same id again
	err = fileDataset.createNamedGraph(ctxt, childGraph, "default")
	if err == nil {
		t.Errorf("Expected to throw error on adding duplicate graph, no error was thrown for dupicate graph: %v", childGraph)
		t.FailNow()
	}

}

// Add duplicate named graph to default graph
func testDuplicateCreateNamedGraphInAnotherNamedGraph(t *testing.T) {

	// TEST NEW REPO INDEX WITH BAD REPO NAME
	ctxt := getContext()
	fileDataset, _ := newFileDataset(ctxt, testRepoName, "test4")
	fileDataset.create(ctxt)

	var parentGraph = "testParent5"
	var childGraph = "childGraphTest5"
	fileDataset.createNamedGraph(ctxt, parentGraph, "default")
	err := fileDataset.createNamedGraph(ctxt, childGraph, parentGraph)
	if err != nil {
		t.Errorf("Unexpected error happened cannot add named graph: %v, %v", childGraph, err)
		t.FailNow()
	}
	// Try adding graph with same id again
	err = fileDataset.createNamedGraph(ctxt, childGraph, parentGraph)
	if err == nil {
		t.Errorf("Expected to throw error on adding duplicate graph, no error was thrown for dupicate graph: %v", childGraph)
		t.FailNow()
	}

}
