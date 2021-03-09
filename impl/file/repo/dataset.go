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
	"bytes"
	"context"
	"encoding/json"

	"io"
	"os"
	"path/filepath"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

const defaultGraphID = "default"

type fileDataset struct {
	// NAME OF REPOSITORY
	repoName string
	// THE RELATIVE & STANDARDIZED PATH TO THE DATASET RESOURCE IN DOGG3RZ FORM
	datasetPath *common.RepositoryPath

	// CANONICAL PATH TO THE PARENT DIR OF THE JSON-LD DOC FILE
	// THAT HOLDS THE DATASET
	parentDirPath string
	// CANONICAL PATH TO THE JSON-LD DOCUMENT FILE
	operatingSystemPath string
}

func newFileDataset(ctxt context.Context, repoName string, datasetPath string) (*fileDataset, error) {

	var err error

	if len(repoName) < 1 {
		return nil, errors.InvalidValue.New("empty repoName")
	}

	if len(datasetPath) < 1 {
		return nil, errors.InvalidValue.New("emtpy datasetPath")

	}

	fds := &fileDataset{}

	fds.repoName = repoName

	if fds.datasetPath, err = common.RepositoryPathNew(datasetPath); err != nil {
		return nil, err
	}

	fds.parentDirPath = filepath.Join(file.RepositoriesDirPath(ctxt),
		fds.repoName, fds.datasetPath.ToOperatingSystemPath())

	fds.operatingSystemPath = filepath.Join(fds.parentDirPath, file.JSONLDDocumentName)

	return fds, nil

}

/* assets whether a dataset exists (state=true) or does not exisst (state=false)
   fileDataset members are initialized in standardized form
   upon function return
*/
func (ds *fileDataset) assertState(ctxt context.Context, state bool) (bool, error) {

	if state {
		// WANT ASSERT DATASET DOES EXIST

		if !file.RepositoryExist(ctxt, ds.repoName) {
			return false, errors.NotFound.Newf("repository '%s' does not exist",
				ds.repoName)
		}

		if !file.FileExists(ds.operatingSystemPath) {
			return false, errors.NotFound.New(ds.operatingSystemPath)
		}

	} else {

		// WANT TO ASSERT DATASET DOES  NOT EXIST

		if file.FileExists(ds.operatingSystemPath) {
			return false, errors.AlreadyExists.New(ds.operatingSystemPath)
		}

	}

	return true, nil

}

func (ds *fileDataset) createNamedGraph(ctxt context.Context, graphID string, parentGraphID string) error {

	if len(graphID) == 0 {
		return errors.InvalidValue.New("graph id cannot be nil or blank")
	}
	if state, err := ds.assertState(ctxt, true); !state {
		return err
	}

	var defaultGraph []interface{}

	var err error

	if _, defaultGraph, err = ds.getDefaultGraph(); err != nil {
		return err
	}

	//	var child interface{}

	if child, _ := getGraph(defaultGraph, graphID); child != nil {
		// means that names graph with this id already exists, return error
		return errors.InvalidValue.New("Named graph with id: " + graphID + " already exists.")
	}

	var nodeAsMap map[string]interface{}
	nodeAsMap = make(map[string]interface{})
	nodeAsMap["@id"] = graphID
	nodeAsMap["@graph"] = make([]interface{}, 0)

	return ds.appendNodeToGraph(ctxt, nodeAsMap, parentGraphID)
}

func (ds *fileDataset) appendNodeToDefaultGraph(ctxt context.Context, newNode map[string]interface{}) error {
	return ds.appendNodeToGraph(ctxt, newNode, defaultGraphID)
}

func (ds *fileDataset) appendNodeToGraph(ctxt context.Context, newNode map[string]interface{}, parentGraphID string) error {

	// ASSERT THAT DATASET EXISTS
	if state, err := ds.assertState(ctxt, true); !state {
		return err
	}

	// READ JSON-LD DOCUMENT INTO MEMORY
	var m map[string]interface{}

	var defaultGraph []interface{}

	var err error
	if m, defaultGraph, err = ds.getDefaultGraph(); err != nil {
		return err
	}

	newNodeIDValue, err1 := getNodeID(newNode)
	if len(newNodeIDValue) == 0 || err != nil {
		return err1
	}

	m["@graph"], err = addNodeToGraph(defaultGraph, newNode, parentGraphID)

	if err != nil {
		return err
	}
	//updateMIME(m)

	ds.writeNodeToFile(m)
	// if default graph is nil and new node id value id

	// if

	return nil
}

func (ds *fileDataset) getDefaultGraph() (map[string]interface{}, []interface{}, error) {

	// READ JSON-LD DOCUMENT INTO MEMORY
	var doc *os.File
	var err error

	if doc, err = os.Open(ds.operatingSystemPath); err != nil {
		return nil, nil, err
	}
	defer doc.Close()

	m := make(map[string]interface{})

	decoder := json.NewDecoder(doc)

	// DESERIALIZE JSON-LD DOC INTO MEMORY
	if err1 := decoder.Decode(&m); err1 != nil {
		return nil, nil, err1
	}

	// GET DEFAULT GRAPH OBJECT FROM DOC
	var defaultGraphRaw interface{}
	var success bool

	// ENSURE DEFAULT GRAPH EXISTS
	if defaultGraphRaw, success = m["@graph"]; !success {
		return nil, nil, errors.NotFound.New(
			"default graph not found as outermost " +
				"attribute '@graph' in JSON-LD document")

	}
	var defaultGraph []interface{}
	if defaultGraphRaw == nil || defaultGraphRaw == "" {
		defaultGraph = make([]interface{}, 0)
	} else {

		if defaultGraph, success = defaultGraphRaw.([]interface{}); !success {
			return nil, nil, errors.InvalidValue.New(
				"Connot convert default graph to list of interface")

		}
	}

	return m, defaultGraph, nil
}

func getNodeID(newNode map[string]interface{}) (string, error) {
	// CHECK IF NEW NODE HAS ID
	var newNodeIDValue string

	if newNodeID, ok := newNode["@id"]; ok {
		// NEW NODE HAS ID. EXTRACT AS TYPE string
		if newNodeIDStr, ok := newNodeID.(string); ok {
			newNodeIDValue = newNodeIDStr

		} else {
			return "", errors.UnexpectedType.Newf(
				"expected @id value of node to append to be type %T, found type %T",
				newNodeIDValue, newNodeID)
		}
	} else {
		newNodeIDValue = ""
	}
	return newNodeIDValue, nil
}

func addNodeToGraph(defaultGraph []interface{}, newNode map[string]interface{}, parentGraphID string) (interface{}, error) {
	//find the parent nodeID

	if parentGraphID == "default" {
		return append(defaultGraph, newNode), nil
	}
	var parentGraphRaw interface{}
	var err error
	if parentGraphRaw, err = getGraph(defaultGraph, parentGraphID); err != nil {
		return nil, err
	}
	var parentGraph []interface{}
	var success bool

	if parentGraph, success = parentGraphRaw.([]interface{}); !success {
		return nil, errors.InvalidValue.New("Failure to convert graph object to list of graphs")
	}
	if nodeid, err1 := getNodeID(newNode); err1 != nil {
		parentMap := parentGraphRaw.(map[string]interface{})
		updateMTIME(parentMap, nodeid)
	}

	return append(parentGraph, newNode), nil

}

func updateMTIME(m map[string]interface{}, newNodeIDValue string) error {

	// UPDATE MTIME FOR THIS RESOURCE (NODE)
	var loc common.JSONLDDocumentLocation

	loc.ContainerType = jsonld.DatasetResource // SINCE IT'S DEFAULT GRAPH, CONTAINER IS DATASET
	loc.ContainerIRI = ""
	loc.ObjectType = jsonld.NodeResource
	loc.ObjectIRI = newNodeIDValue

	if err := common.UpdateResourceMtimeToNow(m, loc); err != nil {
		return err
	}
	return nil
}

func (ds *fileDataset) writeNodeToFile(graph map[string]interface{}) error {
	//  SERIALIZE UPDATED  JSON-LD DOC
	callback := func() (io.Reader, error) {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)

		if err1 := encoder.Encode(&graph); err1 != nil {
			return nil, err1
		}

		return buf, nil
	}

	if _, err1 := file.WriteToFileAtomic(callback, ds.operatingSystemPath); err1 != nil {
		return err1
	}
	return nil
}

func getGraph(graph []interface{}, graphID string) (interface{}, error) {

	if len(graph) == 0 && graphID == defaultGraphID {
		return make([]interface{}, 0), nil
	}
	for _, node := range graph {

		var nodeAsMap map[string]interface{}
		var isMap bool
		if nodeAsMap, isMap = node.(map[string]interface{}); isMap {
			// if graphID is blank then check that @id is not there

			curID, ok := nodeAsMap["@id"]

			if (graphID == defaultGraphID && !ok) || (ok && getIDValue(curID) == graphID) {
				childGraphRaw, _ := nodeAsMap["@graph"]
				return childGraphRaw, nil
			}

			childGraphRaw, success1 := nodeAsMap["@graph"]
			if !success1 {
				return nil, errors.NotFound.New(
					"Graph node found in Graph: " + graphID)

			}

			childGraph, success2 := childGraphRaw.([]interface{})
			if !success2 {
				return nil, errors.NotFound.New(
					"Graph: " + graphID + " graph node cannot be converted to list")

			}
			if len(childGraph) != 0 {
				return getGraph(childGraph, graphID)
			}
		}

	}
	return nil, errors.NotFound.New(
		"Graph: " + graphID + " not found in the map")

}

func getIDValue(id interface{}) string {
	if value, ok := id.(string); ok {
		return value
	}
	return ""

}
