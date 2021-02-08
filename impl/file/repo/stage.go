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
	"io"
	"math"
	"os"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/ipfs"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

type stagingResourceLocationFile struct {
	rescom.StagingResourceLocation

	fileDataset *fileDataset
}

type fileStageResource struct {
	repoName string

	stagedResources []rescom.StagingResource

	startList []stagingResourceLocationFile

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
}

type FileResourceStager struct {
	running  bool // true if stager go-routine is running
	repoName string
	request  chan rescom.StagingResource
	response chan fileResourceStagerResponse
	outcome  chan fileResourceStagerOutcome
}

type fileResourceStagerResponse struct {
	request rescom.StagingResource
	err     error
}

type fileResourceStagerOutcome struct {
	err error
}

func (fsr *fileStageResource) stageResources(ctxt context.Context, repoName string, startList []rescom.StagingResourceLocation) ([]rescom.StagingResource, error) {

	var err error

	// INIT ARRAY THAT WILL HOLD STAGED resources
	fsr.stagedResources = make([]rescom.StagingResource, 0)
	fsr.startList = make([]stagingResourceLocationFile, len(startList))
	fsr.jsonLDParsedDocMap = make(map[string]map[string]interface{})

	fsr.repoName = repoName

	// VALIDATE START LIST
	for i, srl := range startList {
		// copy over to internal start list
		fsr.startList[i].StagingResourceLocation = srl

		// CREATE NEW (FILE) DATASET OBJECT AND ASSERT IT'S VALID (PATH ETC.)
		if fsr.startList[i].fileDataset, err = newFileDataset(ctxt, repoName, srl.DatasetPath); err != nil {
			return fsr.stagedResources, err
		}

		var datasetExists bool

		// ENSURE REPO /  DATASET JSON-LD DOCUMENT FILE EXIST
		if datasetExists, err = fsr.startList[i].fileDataset.assertState(ctxt, true); !datasetExists {
			return fsr.stagedResources, err
		}

		// IS CURENT RESOURCE STAGEABLE?
		if err = fsr.startList[i].StagingResourceLocation.AssertValid(); err != nil {
			// NOT STAGEABLE OR ILLEGAL STATE
			return fsr.stagedResources, err
		}

		// PARSE JSON-LD DATASET, IF NOT ALREADY NOT ALREADY PARSED
		repoDatasetKey := makeJSONLDParsedDocMapKey(repoName, fsr.startList[i].fileDataset)

		// POPULATE MAP OF PARSED JSON-LD DOCUMENTS IDENTIFIED BY KEY
		if _, found := fsr.jsonLDParsedDocMap[repoDatasetKey]; !found {
			var parsedJSONDoc map[string]interface{}

			if parsedJSONDoc, err = parseJSONFile(fsr.startList[i].fileDataset.operatingSystemPath); err != nil {
				return fsr.stagedResources, err
			}

			// POPULATE PARSED JSON-LD DOC MAP
			fsr.jsonLDParsedDocMap[repoDatasetKey] = parsedJSONDoc

		}

	}

	// ITERATE STAGED RESOURCED LOCATIONS IN startList AND STAGE
	// EACH LOCATION IN THE LIST
	for _, srlf := range fsr.startList {
		if err = fsr.stageResource(ctxt, srlf); err != nil {
			return fsr.stagedResources, err
		}
	}

	return fsr.stagedResources, nil

}

// makeJSONLDParsedDocMapKey produces a a unique key using
// repoName and fileDataset as input
func makeJSONLDParsedDocMapKey(repoName string, ds *fileDataset) string {

	return fmt.Sprintf("%s:%s", repoName, ds.datasetPath.ToString())
}

func parseJSONFile(path string) (map[string]interface{}, error) {

	var fp *os.File
	var err error
	var m map[string]interface{}

	if fp, err = os.Open(path); err != nil {
		return nil, err
	}

	m = make(map[string]interface{})

	decoder := json.NewDecoder(fp)

	if err = decoder.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil

}

func (fsr *fileStageResource) stageResource(ctxt context.Context, srlf stagingResourceLocationFile) error {

	//	var err error

	var jsonLDDoc map[string]interface{}
	var ok bool
	var datasetPath = srlf.fileDataset.datasetPath.ToString()
	var err error

	mapKey := makeJSONLDParsedDocMapKey(fsr.repoName, srlf.fileDataset)

	// GET PARSED DATASET FROM MAP
	if jsonLDDoc, ok = fsr.jsonLDParsedDocMap[mapKey]; !ok {
		// SOMETHING HAPPENED UP CALL CHAIN WHERE THE DATASET WAS
		// NOT RETRIEVED
		return errors.NotFound.Newf("expected to find parsed dataset '%s', repository '%s'",
			datasetPath, fsr.repoName)

	}
	// SET CONOTEEXT BEFORE FINDING STAGEABLE REESOURCS WITHIN DOC
	fsr.currentJSONLDParsedDoc = jsonLDDoc
	fsr.currentStagingResourceLocation = srlf.StagingResourceLocation
	fsr.inCurrentStagingResourceLocation = false
	if err = rescom.FindStageableResources(ctxt, datasetPath, jsonLDDoc, fsr); err != nil {

		return err
	}

	return nil

}

func (fsr *fileStageResource) CollectStart(ctxt context.Context, resource interface{}, location rescom.StagingResourceLocation) error {

	if location == fsr.currentStagingResourceLocation {
		fsr.inCurrentStagingResourceLocation = true
	}

	if !fsr.inCurrentStagingResourceLocation {
		// NOT ITERATING WITHIN CURRENT TAARGETE STAGING RESOURCE
		// RETURN
		//fmt.Println("!fsr.inCurrentStagingResourceLocation")
		return nil
	}

	//fmt.Println("in fileStageResource.CollectStart", location)

	var err error
	var index *fileRepositoryIndex

	if index, err = newFileRepositoryIndex(ctxt, fsr.repoName); err != nil {
		return err
	}

	var entry rescom.StagingResource

	entry.StagingResourceLocation = location

	var i interface{}
	var mtimesMap map[string]interface{}
	var ok bool

	// GET MTIMES MAP FROM PARSED JSON-LD DOC
	if i, ok = fsr.currentJSONLDParsedDoc[jsonld.MtimesEntryKeyName]; !ok {
		return errors.NotFound.Newf("can't find parsed mtimes map within parsed "+
			"JSON-LD document for dataset '%s', repository '%s'",
			location.DatasetPath,
			fsr.repoName)
	} else {

		if mtimesMap, ok = i.(map[string]interface{}); !ok {
			return errors.UnexpectedType.Newf("expected parsed JSON-LD document "+
				"entry '%s' value type to be"+
				" type %T, found type %T: %s",
				jsonld.MtimesEntryKeyName,
				mtimesMap, i, location.String())
		}

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

	// STAGE OBJECT INTO IPFS AND GET CID
	if location.WantsCID() {
		//fmt.Println("wants cid ", location)

		// SUBMIT TO IPFS ONLY IF THE RESOURCE LOCATION
		// INDICATES THAT IT WANTS TO HAVE A CID
		if entry.ObjectCID, err = ipfs.DagPut(resource); err != nil {

			return err
		}
		fmt.Println("CID", entry.ObjectCID)
	}

	//return index.update(entry)
	//	err  = index.stage(entry)

	return index.stage(entry) // TODO  RUN AN INDEX UPDATE TX HERE

	//return err
}

func (fsr *fileStageResource) CollectEnd(ctxt context.Context, resource interface{}, location rescom.StagingResourceLocation) {

	if location == fsr.currentStagingResourceLocation {
		fsr.inCurrentStagingResourceLocation = false
	}

	return

}

func (fsr *fileStageResource) String() string {

	return fmt.Sprintf("fileStageResource = { } ")

}

/////////////////////////////////////
// FileResourceStager METHODS
/////////////////////////////////////

// CREATE NEW INSTANCE. THIS SHOOULD BE THE ONLY WAY TO
// INTERACT WITH THIS STRUCT (I.E. CALL IT'S PUBLIC METHODS
func NewFileResourceStager(ctxt context.Context, repoName string) *FileResourceStager {

	stager := &FileResourceStager{}

	// INIT MEMBERS
	stager.repoName = repoName
	stager.request = make(chan rescom.StagingResource)
	stager.response = make(chan fileResourceStagerResponse)
	stager.outcome = make(chan fileResourceStagerOutcome)

	return stager

}

// INITIALIZE MEMBERS AND START GOLANG FUNCTION TO
// OBTAIN LOCK ON INDEX FILE AND WAIT FOR ITEMS
// TO STAGE FROM REQUEST CHANNEL
func (stager *FileResourceStager) runStager(ctxt context.Context) {

	stager_request := stager.request
	stagerRunFunc := func() (io.Reader, error) {

		// TODO: Read index file into memory

		for keepRunning := true; keepRunning; {

			select {

			case stagingResource, ok := <-stager_request:

				keepRunning = ok // UPDATE LOOP CONDITIONAL BASED ON CHANNEL READ

				if !ok {
					// NO MORE REQUESTS
					continue
				}
				// UPDATE MEMCACHED STAGING INDEX
				fmt.Println("runStager: need toupdate mem cached staging index here ", stagingResource)

				// WRITE RESPONSE TO RESPONSE CHANNEL

			case <-ctxt.Done():
				// OPERATION CANCELLED. EXIT GO-ROUTINE
				return nil, errors.Cancelled.New("stagerRunFunc cancelled")
			}

		}
		// TODO: COMMIT INDEX CACHE

		return nil, nil

	}

	repoName := stager.repoName

	stager_outcome := stager.outcome
	// START GO-ROUTINE TO STAGE NEW STAGING RESOURCES
	go func() {

		_, err := file.WriteToFileAtomic(stagerRunFunc,
			file.IndexFilePath(ctxt, repoName))

		outcome := fileResourceStagerOutcome{err: err}

		select {
		case stager_outcome <- outcome:
			// OUTCOME WRITTEN TO OUTCOME CHANNEL
		case <-ctxt.Done():
		}

	}()

	stager.running = true

}

func (stager *FileResourceStager) Add(ctxt context.Context, res rescom.StagingResource) error {

	if !stager.running {
		stager.runStager(ctxt)
	}

	// SUBMIT STAGING RESOURCE TO REQUEST CHANNEL
	select {
	case stager.request <- res:
		// SENT RESOURCE FOR STAGING
	case <-ctxt.Done():
		// OPERATION CANCELLED
		return errors.Newf("%T operation cancelled", stager)

	}

	// WAIT FOR RESPONSE TO REQUEST ON RESPONSE CHANNEL
	select {
	case response := <-stager.response:
		if res != response.request {
			// TODO: WRITE ERROR BACK THAT DESCRIBES REQUEST/RESPONSE MISMATCH!
			// THIS SHOULD NOT HAPPEN AS CHANNEL ONLY TAKES A SINGLE ENTRY
			// BEFORE BLOCKING

			return errors.UnexpectedValue.Newf("expected response to echo request value '%s', found '%s'",
				res, response.request)
		}

		if response.err != nil {
			return response.err
		}

	case <-ctxt.Done():
		// OPERATION CANCELLED
		return errors.Newf("%T operation cancelled", stager)

	}

	return nil
}

func (stager *FileResourceStager) Commit(ctxt context.Context) error {

	if !stager.running {
		// NOTHING TO COMMIT
		return nil
	}

	// CLOSE REQUEST CHANNEL WHICH WILL PROMPT
	// GO-ROUTINE TO FLUSH TO DISK
	close(stager.request)

	// WAIT FOR FINAL OUTCOME OF GO-ROUTINE
	select {
	case outcome := <-stager.outcome:
		return outcome.err
	case <-ctxt.Done():
		// OPERATION CANCELLED
		return errors.Newf("%T operation cancelled", stager)

	}

	return nil
}

func (stager *FileResourceStager) Repository() string {

	return stager.repoName

}
