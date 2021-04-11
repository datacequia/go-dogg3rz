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
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/ipfs"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

type stagingResourceLocationFile struct {
	rescom.StagingResourceLocation

	fileDataset *fileDataset
}

type FileRepositoryResourceStager struct {
	repoName string

	index *fileRepositoryIndex

	//stagedResources []rescom.StagingResource

	//startList []stagingResourceLocationFile

	// CONTAINS MAP OF ALL IN-MEMORY PARSED DATASETS FOR repoName
	jsonLDParsedDocMap map[string]map[string]interface{} // key = REPO_NAME:DATASET_PATH

	// SPECIFIES CURRENT PARSED JSON-LD DOC BEING PROCESSED
	currentJSONLDParsedDoc map[string]interface{}
	// SPECIFIES THE STAGING RESOURCEES BEING SEARCHED FOR STAGING
	// AT THE MOMENT THE SEARCH IS PROCESSING
	currentStagingResourceLocation rescom.StagingResourceLocation
	// INDICATES IF PARSER IS CURRENTLY ITERATING THAT RESOURCE AND/OR RESOURCES
	// WITHIN IF IT'S A CONTAINER RESOURCE
	inCurrentStagingResourceLocation bool

	// Counter of number of resources staged since last commit
	count int
}

type stageCmdCollector struct {
	stager *FileRepositoryResourceStager
}

type removeCmdCollector struct {
	stager *FileRepositoryResourceStager
}

// Create new file repository stager instance
func NewFileRepositoryResourceStager(ctxt context.Context, repoName string) (*FileRepositoryResourceStager, error) {

	var err error
	stager := &FileRepositoryResourceStager{}

	stager.repoName = repoName

	if stager.index, err = newFileRepositoryIndex(ctxt, repoName); err != nil {
		return nil, err
	}

	stager.jsonLDParsedDocMap = make(map[string]map[string]interface{})

	return stager, nil

}

// Add new resource (location) to staging index
func (stager *FileRepositoryResourceStager) Add(ctxt context.Context, srl rescom.StagingResourceLocation) error {

	var err error

	oldCount := stager.count

	if err = stager.stageResource(ctxt, srl); err != nil {
		//fsr.index.rollback()

		return err
	}

	if stager.count == oldCount {
		// COUNT DID NOT CHANGE
		return errors.NotFound.Newf("staging resource not found: %s", srl)
	}

	return err

}

// CollectStart implements the StagingResourceCollector interface CollectStart method
func (collector stageCmdCollector) CollectStart(ctxt context.Context, resource interface{}, location rescom.StagingResourceLocation) error {

	if location == collector.stager.currentStagingResourceLocation {
		// flag that current location is the current staging resource location
		collector.stager.inCurrentStagingResourceLocation = true
	}

	if !collector.stager.inCurrentStagingResourceLocation {
		// NOT ITERATING WITHIN CURRENT TAARGETE STAGING RESOURCE
		// RETURN
		return nil
	}

	var err error

	var entry rescom.StagingResource

	entry.StagingResourceLocation = location

	var i interface{}
	var mtimesMap map[string]interface{}
	var ok bool

	// GET MTIMES MAP FROM PARSED JSON-LD DOC
	if i, ok = collector.stager.currentJSONLDParsedDoc[jsonld.MtimesEntryKeyName]; !ok {
		return errors.NotFound.Newf("can't find parsed mtimes map within parsed "+
			"JSON-LD document for dataset '%s', repository '%s'",
			location.DatasetPath,
			collector.stager.repoName)
	}

	if mtimesMap, ok = i.(map[string]interface{}); !ok {
		return errors.UnexpectedType.Newf("expected parsed JSON-LD document "+
			"entry '%s' value type to be"+
			" type %T, found type %T: %s",
			jsonld.MtimesEntryKeyName,
			mtimesMap, i, location.String())
	}

	// GET LAST MODIFIED TIME FOR LOCATION
	var shaKey string

	if shaKey, err = location.GenerateSHA256Key(); err != nil {
		return err
	}

	var lastModifiedNs interface{}

	if lastModifiedNs, ok = mtimesMap[shaKey]; !ok {
		return errors.NotFound.Newf("can't find mtimes value in parsed JSON-LD document "+
			"for stageable resource: { shaKey: %s, location: %s }",
			shaKey, location.String(),
		)
	}

	var lastModifiedAsFloat64 float64

	// NOTE: encoding/json deserializes JSON Number values as float64
	if lastModifiedAsFloat64, ok = lastModifiedNs.(float64); !ok {
		return errors.UnexpectedType.Newf("expected 'last modified time' attribute  "+
			"for object to be type %T, found type %T: { location: %s }",
			entry.LastModifiedNs, lastModifiedNs, location.String())
	}

	if lastModifiedAsFloat64 > math.MaxInt64 {
		return errors.OutOfRange.Newf("expected 'last modified time' attribute "+
			"for object to be within int64 range ( %v - %v ): found value of %v",
			math.MinInt64, math.MaxInt64, lastModifiedAsFloat64)
	}

	entry.LastModifiedNs = int64(lastModifiedAsFloat64)

	// STAGE OBJECT INTO IPFS AS AN IPLD OBJECT AND GET CID
	if location.WantsCID() {

		// SUBMIT TO IPFS ONLY IF THE RESOURCE LOCATION
		// INDICATES THAT IT WANTS TO HAVE A CID
		if entry.ObjectCID, err = ipfs.DagPut(resource); err != nil {

			return err
		}

		fmt.Println("CID", entry.ObjectCID)
	}

	if collector.stager.index == nil {
		panic("fileStageResource.index is uninitialized!")
	}

	err = collector.stager.index.stage(entry)
	if err == nil {
		collector.stager.count++

	}
	return err

}

// // CollectStart implements the StagingResourceCollector interface CollectEnd method
func (collector stageCmdCollector) CollectEnd(ctxt context.Context, resource interface{}, location rescom.StagingResourceLocation) {

	if location == collector.stager.currentStagingResourceLocation {
		collector.stager.inCurrentStagingResourceLocation = false
	}

}

// String renders FileRepositoryResource instance as a string
func (stager *FileRepositoryResourceStager) String() string {

	return fmt.Sprintf("fileStageResource = { } ")

}

/////////////////////////////////////
// FileResourceStager METHODS
/////////////////////////////////////

// Removes a resource (and its children) from the staging index
func (stager *FileRepositoryResourceStager) Remove(ctxt context.Context, sr rescom.StagingResourceLocation) error {

	return nil
}

// Commits changes to the staging index
func (stager *FileRepositoryResourceStager) Commit(ctxt context.Context) error {

	err := stager.index.commit()
	if err == nil {
		stager.count = 0
	}
	return err
}

// Rollbacks discards any changes to the staging index since last call
// to Commit
func (stager *FileRepositoryResourceStager) Rollback(ctxt context.Context) error {

	err := stager.index.rollback()
	if err == nil {
		stager.count = 0
	}
	return err
}

// Close closes all resources associated with the staging index file
// Note: this must be called after the object is not longer needed
func (stager *FileRepositoryResourceStager) Close(ctxt context.Context) error {
	stager.index.close()
	return nil
}

// Repository returns the repository name for this staging index
func (stager *FileRepositoryResourceStager) Repository() string {

	return stager.repoName

}

//////////////////////////
// UNEXPORTED METHODS
/////////////////////////

// loadUnderlyingDatasetToCache loads dataset identified within srl
// into memory and place into stager's map of cached datasets
// Returns the dataset as a map if successful, otherwise returns an error
func (stager *FileRepositoryResourceStager) loadUnderlyingDatasetToCache(ctxt context.Context, srl rescom.StagingResourceLocation) (map[string]interface{}, error) {

	var err error

	var fds *fileDataset

	// CREATE NEW (FILE) DATASET OBJECT AND ASSERT IT'S VALID (PATH ETC.)
	if fds, err = newFileDataset(ctxt, stager.repoName, srl.DatasetPath); err != nil {
		return nil, err
	}

	var datasetExists bool

	// ENSURE REPO /  DATASET JSON-LD DOCUMENT FILE EXIST
	if datasetExists, err = fds.assertState(ctxt, true); !datasetExists {
		return nil, err
	}

	// IS CURENT RESOURCE LOCATION ASSIGNED A VALID STATE?
	if err = srl.AssertValid(); err != nil {
		// NOT STAGEABLE OR ILLEGAL STATE
		return nil, err
	}

	// IS RESOURCE STAGEABLE?
	if !srl.CanStage() {
		return nil, errors.InvalidValue.Newf("resource not stageable: %s",
			srl.String())
	}

	// PARSE JSON-LD DATASET, IF NOT ALREADY NOT ALREADY PARSED
	repoDatasetKey := makeJSONLDParsedDocMapKey(stager.repoName, srl.DatasetPath)

	// POPULATE MAP OF PARSED JSON-LD DOCUMENTS IDENTIFIED BY KEY
	//	if _, found := stager.jsonLDParsedDocMap[repoDatasetKey]; !found {
	var parsedJSONDoc map[string]interface{}

	if parsedJSONDoc, err = parseJSONFile(fds.operatingSystemPath); err != nil {
		return nil, err
	}

	//fmt.Println("parseJSONDoc", parsedJSONDoc, fds.operatingSystemPath)
	//fmt.Println("repoDatasetKey", repoDatasetKey)
	// POPULATE PARSED JSON-LD DOC MAP
	stager.jsonLDParsedDocMap[repoDatasetKey] = parsedJSONDoc

	//	}

	return parsedJSONDoc, nil

}

// makeJSONLDParsedDocMapKey produces a a unique key using
// repoName and datasetName as input
func makeJSONLDParsedDocMapKey(repoName string, datasetName string) string {

	return fmt.Sprintf("%s:%s", repoName, datasetName)
}

// parseJSONFile parses json-ld document specified by path.
// Returns a parsed map representation of the parsed json-ld document
func parseJSONFile(path string) (map[string]interface{}, error) {

	var fp *os.File
	var err error
	var m map[string]interface{}

	if fp, err = os.Open(path); err != nil {
		return nil, err
	}
	defer fp.Close()

	m = make(map[string]interface{})

	decoder := json.NewDecoder(fp)

	if err = decoder.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil

}

// stageResource stages a resoource location into a file based index
func (stager *FileRepositoryResourceStager) stageResource(ctxt context.Context, srl rescom.StagingResourceLocation) error {

	var jsonLDDoc map[string]interface{}
	var ok bool

	var err error

	if err = srl.AssertValid(); err != nil {
		return err
	}

	mapKey := makeJSONLDParsedDocMapKey(stager.repoName, srl.DatasetPath)

	// GET PARSED DATASET FROM MAP
	if jsonLDDoc, ok = stager.jsonLDParsedDocMap[mapKey]; !ok {
		// DATASET NOT CACHED. LOAD IT.

		jsonLDDoc, err = stager.loadUnderlyingDatasetToCache(ctxt, srl)
		if err != nil {
			return err
		}

	}

	// SET CONTEXT BEFORE FINDING STAGEABLE REESOURCS WITHIN DOC
	stager.currentJSONLDParsedDoc = jsonLDDoc
	stager.currentStagingResourceLocation = srl
	stager.inCurrentStagingResourceLocation = false

	// CREATE JSON-LD DOC COLLECTOR FROM THIS STAGER OBJECT
	// TO ITERATE THE DOCUMENT
	collector := stageCmdCollector{stager}

	// USE COLLECTOR TO LOCATE RESOURCE AND IT'S CHILDREN (IF ANY) AND STAGE
	if err = rescom.FindStageableResources(ctxt, srl.DatasetPath, jsonLDDoc, collector); err != nil {
		return err
	}

	return nil

}
