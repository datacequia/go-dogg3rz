//go:build nothing

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
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

//const jsonLdDocumentName = ".document.jsonld"

func (ds *fileDataset) create(ctxt context.Context) error {

	if state, err := ds.assertState(ctxt, false); !state {
		return err
	}

	// CREATE PARNENT DATASET DIR IF IT DOESN'T ALREADY EXIST
	if err := os.MkdirAll(ds.parentDirPath, 0700); err != nil {
		return err
	}

	// CREATE JSON-LD DOCUMENT WHICH HOLDS DATASET

	if err := createBlankJSONLdDoc(ds.operatingSystemPath); err != nil {
		return err
	}

	return nil
}

func createBlankJSONLdDoc(jsonLdDocPath string) error {

	doc := make(map[string]interface{})

	doc["@context"] = make(map[string]interface{})
	doc["@graph"] = make([]interface{}, 0)
	// CREATE MTIMES ENTRY TO HOLD LAST MODIFIED TIMES OF JSON-LD DOC RESOURCES
	doc[jsonld.MtimesEntryKeyName] = make(map[string]interface{})

	// ADD MTIME ENTRY FOR THE DATASET ITSELF
	loc := common.JSONLDDocumentLocation{
		ObjectType:    jsonld.DatasetResource,
		ObjectIRI:     "",
		ContainerType: jsonld.DatasetResource,
		ContainerIRI:  "",
	}

	locCtxt := common.JSONLDDocumentLocation{
		ObjectType:    jsonld.ContextResource,
		ObjectIRI:     "",
		ContainerType: jsonld.DatasetResource,
		ContainerIRI:  "",
	}

	now := time.Now().UnixNano()

	if err := common.UpdateResourceMtime(doc, loc, now); err != nil {
		return err
	}
	if err := common.UpdateResourceMtime(doc, locCtxt, now); err != nil {
		return err
	}

	// ADD MTIME ENTRY FOR DEFAULT OUTERMOST CONTEXT OBJECT
	// CREATED BY DEFAULT

	b := &strings.Builder{}

	encoder := json.NewEncoder(b)

	if err := encoder.Encode(doc); err != nil {
		return err
	}
	//func() (io.Reader, error) { return strings.NewReader(s), nil }
	if _, err := file.WriteToFileAtomic(
		func() (io.Reader, error) { return strings.NewReader(b.String()), nil },
		jsonLdDocPath); err != nil {

		return err
	}

	return nil

}
