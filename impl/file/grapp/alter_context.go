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
	"bytes"
	"encoding/json"

	"context"
	"io"
	"os"
	"time"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

type addNamespaceNode struct {
	targetLoc   common.StagingResourceLocation
	contextTerm string
	contextIRI  string

	addFlag bool // set to true if added

}

func addNamespaceDataset(ctxt context.Context, grappName string, datasetPath string, term string, iri string) error {

	var fds *fileDataset
	var err error

	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	/* ASSERT THAT DATASET EXISTS */
	if dsExists, noExistErr := fds.assertState(ctxt, true); !dsExists {
		return noExistErr
	}

	// READ JSON-LD DOCUMENT INTO MEMORY
	var doc *os.File
	if doc, err = os.Open(fds.operatingSystemPath); err != nil {
		return err
	}
	defer doc.Close()

	// TODO: convert this func to stream changes to json-ld document
	// instead of renderinig doc in memory
	callback := func() (io.Reader, error) {

		var err error

		m := make(map[string]interface{})

		decoder := json.NewDecoder(doc)

		if err1 := decoder.Decode(&m); err1 != nil {
			return nil, err1
		}

		// add term
		var contextValue map[string]interface{}
		var success bool
		// CAST VALUE FOR @CONTEXT TO A MAP
		if contextValue, success = m["@context"].(map[string]interface{}); !success {

			return nil, errors.UnexpectedType.Newf(
				"expected '@context' value to be type %T , found type %T",
				contextValue, m["@context"])

		}

		// CHECK TO SEE IF TERM ALREADY EXISTS
		if val, ok := contextValue[term]; ok {
			return nil, errors.AlreadyExists.Newf(
				"dataset @context already has term '%s': { value = '%s'}",
				term, val)
		}

		// ADD TERM TO CONTEXT OBJECT
		contextValue[term] = iri

		// UPDATE MTIME FOR CONTEXT OBJECT
		loc := common.JSONLDDocumentLocation{}

		loc.ObjectType = jsonld.ContextResource
		loc.ObjectIRI = "" // context has no @id
		loc.ContainerType = jsonld.DatasetResource
		loc.ContainerIRI = "" // singular Dataset doesn't need IRI

		var key string

		if key, err = loc.GenerateSHA256Key(); err != nil {
			return nil, err
		}

		if v, ok := m[jsonld.MtimesEntryKeyName]; !ok {
			return nil, errors.NotFound.Newf(
				"can't find entry key '%s' in JSON-LD document at %s",
				jsonld.MtimesEntryKeyName,
				fds.operatingSystemPath)
		} else {
			if mtimes, ok := v.(map[string]interface{}); !ok {
				return nil, errors.UnexpectedType.Newf(
					"expected type %T for JSON-LD document entry value for key '%s'"+
						" in JSON-LD document at %s",
					mtimes, jsonld.MtimesEntryKeyName, fds.operatingSystemPath)
			} else {
				// SUCCESS. ASSIGN CURRENT TIMESTAMP
				mtimes[key] = time.Now().UnixNano()
			}
		}

		// SERIALIZE TO JSON DOCUMENT
		buf := &bytes.Buffer{}

		encoder := json.NewEncoder(buf)

		if err1 := encoder.Encode(&m); err1 != nil {
			return nil, err1
		}
		// RETURN BYTE BUFFER OBJECT AS io.Reader
		return buf, nil
	}

	if _, err = file.WriteToFileAtomic(callback, fds.operatingSystemPath); err != nil {
		return err
	}

	return nil
}

// TODO: need to implement this
func (o *addNamespaceNode) execute(ctxt context.Context, grappName string, datasetPath string, nodeID string, term string, iri string) error {

	var fds *fileDataset
	var err error

	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	/* ASSERT THAT DATASET EXISTS */
	if dsExists, noExistErr := fds.assertState(ctxt, true); !dsExists {
		return noExistErr
	}

	// OPEN JSON-LD DOC FILE FOR READING
	var doc *os.File
	if doc, err = os.Open(fds.operatingSystemPath); err != nil {
		return err
	}
	defer doc.Close()

	callback := func() (io.Reader, error) {

		// DESERIALIZE JSON-LD DOC
		dataset := make(map[string]interface{})
		decoder := json.NewDecoder(doc)

		if err = decoder.Decode(&dataset); err != nil {
			return nil, err
		}

		o.targetLoc.DatasetPath = fds.datasetPath.ToString()
		o.targetLoc.ObjectType = jsonld.ContextResource
		o.targetLoc.ObjectIRI = "" // CONTEXT HAS NO IRI
		o.targetLoc.ContainerType = jsonld.NodeResource
		o.targetLoc.ContainerIRI = nodeID

		o.contextTerm = term
		o.contextIRI = iri

		if err = common.FindStageableResources(ctxt, datasetPath, dataset, o); err != nil {
			return nil, err
		}

		// CHECK addNamespaceNode.addFlag TO VERIFY THAT NODE
		// WAS FOUND AND MODIFIED. IF NOT RETURN ERROR
		if !o.addFlag {
			return nil, errors.NotFound.Newf(
				"unable to find node resource identified by '%s', dataset '%s'",
				nodeID, datasetPath)
		}

		// SERIALIZE BACK TO JSON DOCUMENT WITH UPDATES
		buf := &bytes.Buffer{}

		encoder := json.NewEncoder(buf)

		if err = encoder.Encode(&dataset); err != nil {
			return nil, err
		}

		return buf, nil

	}

	// REWRITE JSON-LD DOC WITH UPDATES (NEW NS ADDED TO NODE )
	if _, err = file.WriteToFileAtomic(callback, fds.operatingSystemPath); err != nil {
		return err
	}

	return nil

}

func (o *addNamespaceNode) CollectStart(ctxt context.Context, resource interface{}, location common.StagingResourceLocation) error {

	//	var err error

	if location.ObjectType != o.targetLoc.ContainerType {
		return nil
	}
	if location.ObjectIRI != o.targetLoc.ContainerIRI {
		return nil
	}

	// FOUND CONTEXT'S NODE CONTAINER
	var nodeObject map[string]interface{}
	var contextObj map[string]interface{}
	var ok bool

	// CREATE CONTEXT OBJECT IF IT DOESN'T EXIST
	// IN THE NODE CONTAINER

	if nodeObject, ok = resource.(map[string]interface{}); !ok {
		// NODE RESOURCE IS NOT A  MAP
		return errors.UnexpectedType.Newf(
			"expected to find map object for node identified by '%s', found type %T",
			location.ObjectIRI, resource)

	}

	if c, ok := nodeObject["@context"]; ok {
		if contextObj, ok = c.(map[string]interface{}); !ok {
			return errors.UnexpectedType.Newf(
				"expected to find type %T for context resource within node resource identified by '%s', found type %T instead",
				contextObj, location.ObjectIRI, c)

		}
	} else {
		// NODE OBJECT DOESN'T HAVE CONTEXT OBJECT. CREATE IT
		contextObj = make(map[string]interface{})

		nodeObject["@context"] = contextObj

	}

	// CHECK TO SEE IF TERM ALREADY EXISTS
	if _, ok := contextObj[o.contextTerm]; ok {

		return errors.AlreadyExists.Newf(
			"the @context resource within the node resource identified by '%s' already has the term '%s'",
			location.ObjectIRI, o.contextTerm)
	}

	// ADD TERM/IRI ENTRY  TO CONTEXT
	contextObj[o.contextTerm] = o.contextIRI
	// set addFlag to true
	o.addFlag = true

	return nil
}

func (o *addNamespaceNode) CollectEnd(ctxt context.Context, resource interface{}, location common.StagingResourceLocation) {

}
