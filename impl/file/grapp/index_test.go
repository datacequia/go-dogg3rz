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
	//	"fmt"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/datacequia/go-dogg3rz/errors"
	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/config"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

var dogg3rzHome string
var fileGrappIdx *fileGrapplicationIndex

// returns a cancellable cnotext that inits DOGG3RZ_HOME
// to package var dogg3rzHome
func getContext() context.Context {

	if len(dogg3rzHome) < 1 {
		panic("dogg3rzHome not set")
	}

	ctxt := context.Background()

	ctxt = context.WithValue(ctxt, "DOGG3RZ_HOME", dogg3rzHome)

	return ctxt
}

// setup temp dogg3rz env to test index
func indexSetup(t *testing.T) {

	dogg3rzHome = filepath.Join(os.TempDir(),
		fmt.Sprintf("index_test_%d", time.Now().UnixNano()))

	ctxt := getContext()
	// os.Setenv("DOGG3RZ_HOME", dogg3rzHome)

	fileNodeResource := &filenode.FileNodeResource{}

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.ActivityPubUserHandle = "@test@dogg3rz.com"

	if err := fileNodeResource.InitNode(ctxt, dgrzConf); err != nil {
		t.Error(err)
	}
	//	t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	fileGrapplicationResource := FileGrapplicationResource{}

	if err := fileGrapplicationResource.InitGrapp(ctxt, testGrappName); err != nil {
		t.Error(err)
	}

}

// removes all file resources created in tmp dir
func indexTeardown(t *testing.T) {

	os.RemoveAll(dogg3rzHome)

}

// main index test. all other tests called sequentially from here
func TestIndex(t *testing.T) {

	// SETUP CODE
	indexSetup(t)
	defer indexTeardown(t)

	//fileNodeResource = FileNo

	testNewFileGrappIndexWithNonExistentGrapp(t)

	//testNewFileGrappIndexThenCancel(t)

	testNewFileGrappIndexThenScanWithNoIndexFile(t)

	testNewFileGrappIndexThenEmptyCommit(t)

	testNewFileGrappIndexThenEmptyRollback(t)

	testFileGrappIndexUpdateAndReadBackInTx(t)

	testFileGrappIndexUpdateCommitAndReadBack(t)

	testFileGrappIndexUpdateCommitAndScanForOne(t)

	testFileGrappIndexStageThreeCommitAndStageUpdateOneReadBack(t)

	testRemoveSingleNodeResourceFromIndex(t)

	testLockIndexOnModify(t)

	testInvalidIndexEntryValidate(t)

	testNewFileGrappIndexOnNonExistentGrapp(t)

	testRemoveSingleNamedGraphResourceWithChildrenFromIndex(t)
	indexTeardown(t)

}

// returns a new index object with a context supplied cancel callback
func newFileGrappIdxWithCancelFunc(t *testing.T) (*fileGrapplicationIndex, context.Context, context.CancelFunc) {
	// TEST NEW GRAPP INDEX WITH BAD GRAPP NAME
	ctxt := getContext()

	var cancelFunc context.CancelFunc

	ctxt, cancelFunc = context.WithCancel(ctxt)

	var f *fileGrapplicationIndex
	var err error

	if f, err = newFileGrapplicationIndex(ctxt, testGrappName); err != nil {
		t.Errorf("newFileGrappIdxWithCancelFunc(): %s", err)
	}

	// make sure index lock file not there before returning
	// from previous test

	/*
		indexLockFile := f.path + file.LOCK_FILE_SUFFIX
		var i int
		for i = 0; file.FileExists(indexLockFile); i++ {
			//fmt.Printf("attempt %d: index lock file exists at %s: waiting...\n", i, indexLockFile)
			time.Sleep(time.Millisecond * 250)

		}
	*/

	return f, ctxt, cancelFunc

}

// test that failure occurs when attempt to instantiate index object from
// non-existent grapp
func testNewFileGrappIndexWithNonExistentGrapp(t *testing.T) {

	if _, err := newFileGrapplicationIndex(getContext(), testGrappName+"-noExist"); err == nil {
		t.Errorf("newFileGrapplicationIndex(): succeeded with non-existent grapp (name)")
	}

}

// tests to make sure that
// scan() will respond with success even thoough a
// new grapplication will not yet have an index file
// created because no commits have been issued yet on the index
func testNewFileGrappIndexThenScanWithNoIndexFile(t *testing.T) {

	f, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer f.close()

	var req *msgIndexResultSetRequest
	var err error

	if req, err = f.scan(func(arg1 common.StagingResource) bool {
		return false
	}); err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	var result *common.StagingResource
	//	var err error

	for result, err = req.next(); result != nil && err == nil; result, err = req.next() {
		// NOT EXPECTING ANY RESULTS
		t.Errorf("Not expecting results, found %s", result)

	}

	if err != nil {

		t.Errorf("scan() failed: %s", err)
	}

}

// tests a new instance followed by a commit with no prior pending changes
func testNewFileGrappIndexThenEmptyCommit(t *testing.T) {
	f, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer f.close()

	if _, err := os.Stat(f.path); err == nil {
		t.Errorf("expected index file at %s to NOT exist", f.path)
		return
	}

	if err := f.commit(); err != nil {
		if errors.GetType(err) != errors.EmptyCommit {
			t.Errorf("expected EmptyCommit error, found %s", err)
			return
		}
	} else {
		t.Errorf("expected Empty commit to produce error, returned nil")
		return
	}

	if _, err := os.Stat(f.path); err == nil {
		t.Errorf("expected index file at %s to NOT exist", f.path)
		return
	}

}

// tests that rollback operation is a no-op on non-existent index and
// does not throw an error
func testNewFileGrappIndexThenEmptyRollback(t *testing.T) {

	f, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer f.close()

	if _, err := os.Stat(f.path); err == nil {
		t.Errorf("expected index file at %s to NOT exist", f.path)
	}

	if err := f.rollback(); err != nil {

		t.Errorf("expected empty rollback to have no errors, found %s", err)

	}

	if _, err := os.Stat(f.path); err == nil {
		t.Errorf("index file at %s exists after empty rollback", f.path)
	}

}

// tests that index stage on empty index followed by rollback and rescan
// results in zero entries in the index
func testFileGrappIndexUpdateAndReadBackInTx(t *testing.T) {

	var requests [3]common.StagingResource

	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	// stage 3 requests
	for _, r := range requests {

		if err := index.stage(r); err != nil {
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	// READ BACK EVERYTHING
	scanReq, err := index.scan(func(sr common.StagingResource) bool { return true })

	var result *common.StagingResource
	//var err error
	var matchCnt int
	var iterateCnt int

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		for _, r := range requests {
			//t.Logf("cmp %s == %s", *result, r)
			if *result == r {
				matchCnt++
			}
		}
	}

	if err != nil {
		t.Errorf("error result from scan request: %s", err)
		return
	}

	if matchCnt != 3 {
		t.Errorf("failed retrieve uncommitted staged requests back via scan(): { matchCnt = %d, iterateCount = %d }",
			matchCnt, iterateCnt)
		return
	}

	// ROLLBACK CHANGES
	if err := index.rollback(); err != nil {
		t.Errorf("failed to rollback changes")
	}

	// SCAN AGAIN SHOULD HAVE ZERO ELEMENTS
	scanReq, err = index.scan(func(sr common.StagingResource) bool { return true })
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	iterateCnt = 0
	matchCnt = 0

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		for _, r := range requests {
			//t.Logf("cmp %s == %s", *result, r)
			if *result == r {
				matchCnt++
			}
		}
	}
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	if iterateCnt != 0 {
		t.Errorf("expected zero results after rollback found %d: { matched results = %d }",
			iterateCnt, matchCnt)
	}

}

// tests that staged and commmitted changes to index can be scanned back
func testFileGrappIndexUpdateCommitAndReadBack(t *testing.T) {

	var requests [3]common.StagingResource

	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	//t.Logf("testFileGrappIndexUpdateAndReadBack() ...")
	//defer t.Logf("testFileGrappIndexUpdateAndReadBack() !")
	for _, r := range requests {

		if err := index.stage(r); err != nil {
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	// READ BACK EVERYTHING
	scanReq, err := index.scan(func(sr common.StagingResource) bool { return true })
	if err != nil {
		t.Errorf("scan failed: %s", err)
	}

	var result *common.StagingResource
	//	var err error
	var matchCnt int
	var iterateCnt int

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		for _, r := range requests {
			//t.Logf("cmp %s == %s", *result, r)
			if *result == r {
				matchCnt++
			}
		}
	}

	if err != nil {
		t.Errorf("error result from scan request: %s", err)
		return
	}

	if matchCnt != 3 {
		t.Errorf("failed retrieve uncommitted staged requests back via scan(): { matchCnt = %d, iterateCount = %d }",
			matchCnt, iterateCnt)
		return
	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes")
		return
	}

	// CHECK THAT INDEX FILE EXIST
	if _, err := os.Stat(index.path); err != nil {
		t.Errorf("index file at %s does NOT exist after commit", index.path)
		return
	}

	t.Logf("index file exists at %s", index.path)

	//fmt.Println("doing re-scan after rollback ")

	// SCAN AGAIN SHOULD HAVE THREE ELEMENTS
	scanReq, err = index.scan(func(sr common.StagingResource) bool { return true })
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	iterateCnt = 0
	matchCnt = 0

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		for _, r := range requests {
			//t.Logf("cmp %s == %s", *result, r)
			if *result == r {
				matchCnt++
			}
		}
	}
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	if iterateCnt != 3 || matchCnt != 3 {
		t.Errorf("expected three results after commit found %d: { matched results = %d }",
			iterateCnt, matchCnt)
		return
	}

}

// tests that issuing a filtered  scan for one committed resource of multiple can be scanned
// back after committing said resources
func testFileGrappIndexUpdateCommitAndScanForOne(t *testing.T) {

	var requests [3]common.StagingResource

	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	//t.Logf("testFileGrappIndexUpdateAndReadBack() ...")
	//defer t.Logf("testFileGrappIndexUpdateAndReadBack() !")
	for _, r := range requests {

		if err := index.stage(r); err != nil {
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	searchReq := requests[2]

	// READ BACK EVERYTHING
	filterFunc := func(sr common.StagingResource) bool {
		if sr == searchReq {
			return true
		}
		return false
	}

	scanReq, err := index.scan(filterFunc)
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	var result *common.StagingResource
	//var err error
	var matchCnt int
	var iterateCnt int

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		//	for _, r := range requests {
		//t.Logf("cmp %s == %s", *result, r)
		if *result == searchReq {
			matchCnt++
		}
		//	}
	}

	if err != nil {
		t.Errorf("error result from scan request: %s", err)
		return
	}

	if matchCnt != 1 {
		t.Errorf("failed scan filtered single staged request back via scan(): { matchCnt = %d, iterateCount = %d }",
			matchCnt, iterateCnt)
		return
	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes")
		return
	}

	// CHECK THAT INDEX FILE EXIST
	if _, err := os.Stat(index.path); err != nil {
		t.Errorf("index file at %s does NOT exist after commit", index.path)
		return
	}

	t.Logf("index file exists at %s", index.path)

	//fmt.Println("doing re-scan after rollback ")

	// SCAN AGAIN SHOULD HAVE THREE ELEMENTS
	scanReq, err = index.scan(filterFunc)
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	iterateCnt = 0
	matchCnt = 0

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		//	for _, r := range requests {
		//t.Logf("cmp %s == %s", *result, r)
		if *result == searchReq {
			matchCnt++
		}
		//	}
	}
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	if matchCnt != 1 {
		t.Errorf("expected three results after commit found %d: { matched results = %d }",
			iterateCnt, matchCnt)
	}

}

func testFileGrappIndexStageThreeCommitAndStageUpdateOneReadBack(t *testing.T) {

	var requests [3]common.StagingResource

	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	// Stage 3 entries
	for _, r := range requests {

		if err := index.stage(r); err != nil {
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes")
	}

	// CHECK THAT INDEX FILE EXIST
	if _, err := os.Stat(index.path); err != nil {
		t.Errorf("index file at %s does NOT exist after commit", index.path)
		return
	}

	//t.Logf("index file exists at %s", index.path)

	// UPDATAE A SINGLE NODE RESOURCE CID  previously staged
	updatedSR := requests[2]
	updatedSR.ObjectCID = "bafyreigqrowux55fl53qzozsy3wrc3r56xeav246cnyoscjmg7lavuwe7e"

	if err := index.stage(updatedSR); err != nil {
		t.Logf("request failed: { request = %s, err = %s} ", updatedSR, err)
	}

	filterFunc := func(sr common.StagingResource) bool {
		return true
	}

	//fmt.Println("doing re-scan after rollback ")

	// SCAN AGAIN SHOULD HAVE THREE ELEMENTS
	scanReq, err := index.scan(filterFunc)
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	iterateCnt := 0
	matchCnt := 0

	var result *common.StagingResource
	//	var err error

	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		//	for _, r := range requests {
		//t.Logf("cmp %s == %s", *result, r)
		if *result == updatedSR {
			matchCnt++
		}
		//	}
	}
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	if matchCnt != 1 || iterateCnt != 3 {
		t.Errorf("expected one result, found %d: { matched results = %d }",
			iterateCnt, matchCnt)
		return
	}

	// COMMIT STAGED UPDATE
	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes")
		return
	}

	requests[2] = updatedSR // UPDATE SECOND ENTRY TO UPDATED VALUE

	// scan back everything, and ensure update there
	// READ BACK EVERYTHING

	scanReq, err = index.scan(filterFunc)
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	iterateCnt = 0
	matchCnt = 0

	// SCAN ALL STAGED REQUESTS AND TRY MATCHING UUPDATED ONE
	for result, err = scanReq.next(); result != nil && err == nil; result, err = scanReq.next() {
		iterateCnt++
		//	for _, r := range requests {
		//t.Logf("cmp %s == %s", *result, r)
		if *result == updatedSR {
			matchCnt++
		}
		//	}
	}

	if err != nil {
		t.Errorf("error result from scan request: %s", err)
		return
	}

	if matchCnt != 1 || iterateCnt != 3 {
		t.Errorf("failed retrieve uncommitted staged requests back via scan(): { matchCnt = %d, iterateCount = %d }",
			matchCnt, iterateCnt)
		return

	}

}

// tests that a removed resource from a previously committed collection
// is in fact removed
func testRemoveSingleNodeResourceFromIndex(t *testing.T) {

	var requests [3]common.StagingResource

	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	// Stage 3 entries
	for _, r := range requests {

		if err := index.stage(r); err != nil {
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes: %s", err)
		return
	}

	nodeResource := requests[2]

	// REMOVE NODE RESOURCE
	removeReq, err := index.remove(nodeResource)
	if err != nil {
		t.Errorf("remove failed: %s", err)
		return
	}

	var result *common.StagingResource
	//var err error
	for result, err = removeReq.next(); result != nil; result, err = removeReq.next() {
		t.Logf("removed %s", result)

	}

	if err != nil {
		t.Errorf("failed to remove node resource:  %s", err)
		return
	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes: %s", err)
		return
	}
	iterateCnt := 0
	resultCnt := 0
	scanReq, err := index.scan(func(common.StagingResource) bool { return true })
	if err != nil {
		t.Errorf("scan failed: %s", err)
		return
	}

	for r, e := scanReq.next(); r != nil; r, e = scanReq.next() {
		iterateCnt++
		if *r == requests[0] {
			resultCnt++
		}
		if *r == requests[1] {
			resultCnt++
		}
		err = e
	}

	if err != nil {
		t.Errorf("failed on scan after remove: %s", err)
		return
	}

	if resultCnt != 2 {
		t.Errorf("expected 2 entries after index.remove(), found %d: iterateCnt = %d", resultCnt, iterateCnt)
	}

}

// tests if the remove operation will remove a container based resource and
// its children from the index
func testRemoveSingleNamedGraphResourceWithChildrenFromIndex(t *testing.T) {

	var requests []common.StagingResource
	var entry common.StagingResource

	// OUTERMOST 'document' @context
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.ContextResource
	entry.ObjectIRI = ""
	entry.ContainerType = jsonld.DatasetResource
	entry.ContainerIRI = ""
	entry.LastModifiedNs = 1600103677854799000
	entry.ObjectCID = "bafyreigcm277jvvdmenqkudvan3mn7icvzdj2a3eygtgilkf2mypcrkgvi"

	requests = append(requests, entry)

	// NODE IN DEFAULT GRAPH
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.NodeResource
	entry.ObjectIRI = "mynode1" // DEFAULT GRAPH. NO NAME
	entry.ContainerType = jsonld.DatasetResource
	entry.ContainerIRI = ""
	entry.LastModifiedNs = 1600103677853633000
	entry.ObjectCID = "bafyreie5h75u3kv47vywgzsohnqlmxowfv4herxrj6ezdlttx47wrkpbkm"

	requests = append(requests, entry)

	// NAMED GRAPH AS CHILD WITHIN DEFAULT GRAPH
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.NamedGraphResource
	entry.ObjectIRI = "myNamedGraph"
	entry.ContainerType = jsonld.DatasetResource
	entry.ContainerIRI = ""
	entry.LastModifiedNs = 1600103677853633999
	entry.ObjectCID = "" //"bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"

	requests = append(requests, entry)

	myNamedGraph := entry

	// NODE WITHIN NAMED GRAPH
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.NodeResource
	entry.ObjectIRI = "subNode1"
	entry.ContainerType = jsonld.NamedGraphResource
	entry.ContainerIRI = "myNamedGraph"
	entry.LastModifiedNs = 1600103677853633888
	entry.ObjectCID = "bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"

	requests = append(requests, entry)

	// SECOND NODE WITHIN NAMED GRAPH
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.NodeResource
	entry.ObjectIRI = "subNode2"
	entry.ContainerType = jsonld.NamedGraphResource
	entry.ContainerIRI = "myNamedGraph"
	entry.LastModifiedNs = 1600103677853633123
	entry.ObjectCID = "bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"

	requests = append(requests, entry)

	//	requests[0], requests[1], requests[2] = getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	t.Logf("about to stage %d msgs to be later removed...", len(requests))
	// Stage  entries
	for _, r := range requests {
		t.Logf("staging %s", r)
		if err := index.stage(r); err != nil {
			//	panic("here")
			t.Errorf("request failed: { request = %s, err = %s} ", r, err)
			return
		}

	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes: %s", err)

		return
	}

	scanCnt := 0

	if scanReq, err := index.scan(func(sr common.StagingResource) bool { return true }); err != nil {
		t.Errorf("scan failed: %s", err)
		return
	} else {
		for x, _ := scanReq.next(); x != nil; x, _ = scanReq.next() {
			scanCnt++
		}
	}

	if scanCnt != len(requests) {
		t.Errorf("expected %d test index entries. found %d", len(requests), scanCnt)
		return
	}

	// REMOVE NODE RESOURCE
	removeReq, err := index.remove(myNamedGraph)
	if err != nil {
		t.Errorf("remove failed: %s", err)
		return
	}

	var result *common.StagingResource
	//	var err error

	// iterate filtered remove matches
	var removeCnt int
	for result, err = removeReq.next(); result != nil; result, err = removeReq.next() {
		t.Logf("removing %s", result)
		removeCnt++
	}

	if err != nil {
		t.Errorf("failed to remove node resource:  %s", err)
		return
	}

	if removeCnt != 3 {
		t.Errorf("expected 3 staging resouorces to be removed, %d were removed instead", removeCnt)
		return
	}

	// COMMIT CHANGES
	if err := index.commit(); err != nil {
		t.Errorf("failed to commit changes: %s", err)

	}

	// SCAN ALL AGAIN
	scanCnt = 0
	if scanReq, err := index.scan(func(sr common.StagingResource) bool { return true }); err != nil {
		t.Errorf("scan failed: %s", err)
		return
	} else {
		for x, _ := scanReq.next(); x != nil; x, _ = scanReq.next() {
			scanCnt++
		}
	}

	if scanCnt != len(requests)-removeCnt {
		t.Errorf("expected %d test index entries. found %d", len(requests)-removeCnt, scanCnt)
		return
	}

}

func testLockIndexOnModify(t *testing.T) {

	var entry common.StagingResource

	// OUTERMOST 'document' @context
	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.ContextResource
	entry.ObjectIRI = ""
	entry.ContainerType = jsonld.DatasetResource
	entry.ContainerIRI = ""
	entry.LastModifiedNs = 1600103677854799000
	entry.ObjectCID = "bafyreigcm277jvvdmenqkudvan3mn7icvzdj2a3eygtgilkf2mypcrkgvi"

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	// NOW THAT LOCK FILE EXISTS. TRY TO STAGE SOOMETHING

	if err := index.stage(entry); err != nil {

		t.Errorf("stage failed: %s", err)
		return
	}

	// 	if here lock file is in place

	index2, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index2.close()

	// NOW THAT LOCK FILE EXISTS. TRY TO STAGE SOOMETHING in SECND OBJECT
	if err := index2.stage(entry); err != nil {

		if errors.GetType(err) != errors.TryAgain {
			t.Errorf("stage failed: %s", err)
			return
		}
		// IF HERE THEN LOCK WORKED

	}

	// commit first index obbject
	if err := index.commit(); err != nil {
		t.Errorf("commit failed: %s", err)

		return
	}

	// NOW TRY TO COMMIT FROM INDEX 2
	entry.LastModifiedNs = 1600103677854799001

	if err := index2.stage(entry); err != nil {
		t.Errorf("expected stage success, got error: %s", err)
		return
	}

	if err := index2.commit(); err != nil {
		t.Errorf("commit failed: %s", err)
		return
	}

	var req *msgIndexResultSetRequest
	var err error
	if req, err = index2.scan(func(x common.StagingResource) bool { return true }); err != nil {

		t.Errorf("scan failed: %s", err)
		return
	}

	var sr *common.StagingResource
	var scanCnt int
	for sr, err = req.next(); sr != nil; sr, err = req.next() {
		scanCnt++
		if *sr != entry {
			t.Errorf("read back failed: expected %s, found %s", entry, *sr)
			return
		}
	}

	if err != nil {
		t.Errorf("scan next() failed: %s", err)
		return
	}

	if scanCnt != 1 {
		t.Errorf("scan back expected 1 result, found %d", scanCnt)
	}

}

// tests that malformed index entries or otherwise in an illegal state
// can be caught
func testInvalidIndexEntryValidate(t *testing.T) {

	entry, _, _ := getThreeEntries()

	index, _, _ := newFileGrappIdxWithCancelFunc(t)
	defer index.close()
	defer os.Remove(index.path)

	// SET BAD TYPE. SELF ASSIGNED AND NOT ALLOCATED IN primitives PACKAGE
	//var badType string = "dogg3rz.badtype"
	var badObjectType jsonld.JSONLDResourceType = math.MaxInt8

	// TEST BAD OBJECT TYPE
	holdType := entry.ObjectType
	entry.ObjectType = badObjectType
	if err := index.stage(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Type value assigned")
		return
	}

	entry.ObjectType = holdType

	// TEST BAD DATASETPATH
	holdDatasetPath := entry.DatasetPath
	entry.DatasetPath = "%badPathElement/$foo"
	if err := index.stage(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Type value assigned")
	}

	entry.DatasetPath = holdDatasetPath

	if err := index.rollback(); err != nil {
		t.Errorf("error on rollback: %s", err)
	}

}

// tests that creating a new index object on a non existent grapp fails
func testNewFileGrappIndexOnNonExistentGrapp(t *testing.T) {

	var nonExistGrapp = "not." + testGrappName

	ctxt := getContext()

	if index, err := newFileGrapplicationIndex(ctxt, nonExistGrapp); err == nil {
		index.close()
		t.Errorf("testNewFileGrappIndexOnNonExistentGrapp(): did not fail "+
			"on non-existent grapplication: %s", nonExistGrapp)
	} else {

		t.Logf("testNewFileGrappIndexOnNonExistentGrapp: %s", err)
	}

}

// getThreeEntries returns test index resources for testing
func getThreeEntries() (common.StagingResource, common.StagingResource, common.StagingResource) {

	var entry = common.StagingResource{}

	entry.DatasetPath = "data1"
	entry.ObjectType = jsonld.ContextResource
	entry.ObjectIRI = ""
	entry.ContainerType = jsonld.DatasetResource
	entry.ContainerIRI = ""
	entry.LastModifiedNs = 1600103677854799000
	entry.ObjectCID = "bafyreigcm277jvvdmenqkudvan3mn7icvzdj2a3eygtgilkf2mypcrkgvi"

	var entry2 = common.StagingResource{}

	entry2.DatasetPath = "data2"
	entry2.ObjectType = jsonld.NamedGraphResource
	entry2.ObjectIRI = ""
	entry2.ContainerType = jsonld.DatasetResource
	entry2.ContainerIRI = ""
	entry2.LastModifiedNs = 1600103677853633000
	entry2.ObjectCID = "" //"bafyreie5h75u3kv47vywgzsohnqlmxowfv4herxrj6ezdlttx47wrkpbkm"

	var entry3 = common.StagingResource{}

	entry3.DatasetPath = "data3"
	entry3.ObjectType = jsonld.NodeResource
	entry3.ObjectIRI = "http://www.doggg3rz.com/my/test#me"
	entry3.ContainerType = jsonld.NamedGraphResource
	entry3.ContainerIRI = ""
	entry3.LastModifiedNs = 1600103677853633999
	entry3.ObjectCID = "bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"

	return entry, entry2, entry3

}
