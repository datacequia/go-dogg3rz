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

func newFileDataset(repoName string, datasetPath string, ctxt context.Context) (*fileDataset, error) {

	var err error

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
func (ds *fileDataset) assertState(state bool, ctxt context.Context) (bool, error) {

	if state {
		// WANT ASSERT DATASET DOES EXIST

		if !file.RepositoryExist(ds.repoName, ctxt) {
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

func (fds *fileDataset) appendNodeToDefaultGraph(newNode map[string]interface{}, ctxt context.Context) error {

	// ASSERT THAT DATASET EXISTS
	if state, err := fds.assertState(true, ctxt); !state {
		return err
	}

	// READ JSON-LD DOCUMENT INTO MEMORY
	var doc *os.File
	var err error

	if doc, err = os.Open(fds.operatingSystemPath); err != nil {
		return err
	}
	defer doc.Close()

	// TODO: convert this func to stream changes to json-ld document
	// instead of renderinig doc in memory
	callback := func() (io.Reader, error) {

		buf := &bytes.Buffer{}

		m := make(map[string]interface{})

		decoder := json.NewDecoder(doc)

		// DESERIALIZE JSON-LD DOC INTO MEMORY
		if err1 := decoder.Decode(&m); err1 != nil {
			return nil, err1
		}

		// GET DEFAULT GRAPH OBJECT FROM DOC
		var defaultGraphRaw interface{}
		var defaultGraph []interface{}
		var success bool

		// ENSURE DEFAULT GRAPH EXISTS
		if defaultGraphRaw, success = m["@graph"]; !success {
			return nil, errors.NotFound.New(
				"default graph not found as outermost " +
					"attribute '@graph' in JSON-LD document")

		}

		if defaultGraph, success = defaultGraphRaw.([]interface{}); !success {

			return nil, errors.UnexpectedType.Newf(
				"expected '@graph' value to be type %T , found type %T",
				defaultGraph, m["@graph"])

		}

		// CHECK IF NEW NODE HAS ID
		var newNodeIDValue string

		if newNodeID, ok := newNode["@id"]; ok {
			// NEW NODE HAS ID. EXTRACT AS TYPE string
			if newNodeIDStr, ok := newNodeID.(string); ok {
				newNodeIDValue = newNodeIDStr

			} else {
				return nil, errors.UnexpectedType.Newf(
					"expected @id value of node to append to be type %T, found type %T",
					newNodeIDValue, newNodeID)
			}
		}

		for _, node := range defaultGraph {

			var nodeAsMap map[string]interface{}
			var isMap bool
			if nodeAsMap, isMap = node.(map[string]interface{}); !isMap {
				continue
			}

			if curID, ok := nodeAsMap["@id"]; ok {
				if curIDValue, ok := curID.(string); ok {
					// TODO: need to expand id's to IRI before comparison
					if curIDValue == newNodeIDValue {
						return nil, errors.AlreadyExists.Newf(
							"node with @id value of '%s' already exists in the default graph",
							newNodeIDValue)

					}
				}
			}
		}

		// UPDATE IN MEMORY JSON-LD DOC
		m["@graph"] = append(defaultGraph, newNode)

		// UPDATE MTIME FOR THIS RESOURCE (NODE)
		var loc common.JSONLDDocumentLocation

		loc.ContainerType = jsonld.DatasetResource // SINCE IT'S DEFAULT GRAPH, CONTAINER IS DATASET
		loc.ContainerIRI = ""
		loc.ObjectType = jsonld.NodeResource
		loc.ObjectIRI = newNodeIDValue

		if err = common.UpdateResourceMtimeToNow(m, loc); err != nil {
			return nil, err
		}

		//  SERIALIZE UPDATED  JSON-LD DOC
		encoder := json.NewEncoder(buf)

		if err1 := encoder.Encode(&m); err1 != nil {
			return nil, err1
		}

		return buf, nil
	}

	if _, err = file.WriteToFileAtomic(callback, fds.operatingSystemPath); err != nil {
		return err
	}

	return nil
}
