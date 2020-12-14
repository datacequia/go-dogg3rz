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
	//	"fmt"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/config"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

var dogg3rzHome string
var fileRepoIdx *fileRepositoryIndex

const (
	testRepoName = "index_test"
)

func getContext() context.Context {
	ctxt := context.Background()

	ctxt = context.WithValue(ctxt, "DOGG3RZ_HOME", dogg3rzHome)

	return ctxt
}

func indexSetup(t *testing.T) {

	dogg3rzHome = filepath.Join(os.TempDir(),
		fmt.Sprintf("index_test_%d", time.Now().UnixNano()))

	ctxt := getContext()
	// os.Setenv("DOGG3RZ_HOME", dogg3rzHome)

	fileNodeResource := &filenode.FileNodeResource{}

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.Email = "test@dogg3rz.com"

	if err := fileNodeResource.InitNode(ctxt, dgrzConf); err != nil {
		t.Error(err)
	}
	t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	fileRepositoryResource := FileRepositoryResource{}

	if err := fileRepositoryResource.InitRepo(ctxt, testRepoName); err != nil {
		t.Error(err)
	}

}

func indexTeardown(t *testing.T) {

	os.RemoveAll(dogg3rzHome)

}
func TestIndex(t *testing.T) {

	// SETUP CODE
	indexSetup(t)

	//fileNodeResource = FileNo
	testNewFileRepoIndex(t)

	testFileRepoIndexUpdateAndReadBack(t)

	testUpdateExistingAndReadBack(t)

	testInvalidIndexEntryValidate(t)

	testNewFileRepoIndexOnNonExistentRepo(t)
	testReadIndexFileFailsOnUpdate(t)

	// TEARDOWN CODE
	//	indexTeardown(t)

}

// ADD 3 NEW ENTRIES TO THE INDEX AND READ THEM
// BACK FROM INDEX AND COMPARE . SHOULD BE EXACTLY SAME
func testNewFileRepoIndex(t *testing.T) {

	// TEST NEW REPO INDEX WITH BAD REPO NAME
	ctxt := getContext()

	if _, err := newFileRepositoryIndex(ctxt, testRepoName+"-badName"); err == nil {
		t.Errorf("newFileRepositoryIndex(): succeeded with non-existent repo (name)")
	}

	if fri, err := newFileRepositoryIndex(ctxt, testRepoName); err != nil {
		t.Errorf("newFileRepositoryIndex(): failed  with existing repo (name): %v", err)
	} else {
		fileRepoIdx = fri
	}

	if fileRepoIdx == nil {
		t.Errorf("newFileRepositoryIndex(): returnd nil object on success")
	}
}

func testFileRepoIndexUpdateAndReadBack(t *testing.T) {

	entry, entry2, entry3 := getThreeEntries()

	if err := fileRepoIdx.update(entry); err != nil {
		t.Errorf("testFileRepoIndexUpdate(): fileRepositoryIndex.update() failed: %s", err)
	}
	if err := fileRepoIdx.update(entry2); err != nil {
		t.Errorf("testFileRepoIndexUpdate(): fileRepositoryIndex.update() failed: %s", err)
	}
	if err := fileRepoIdx.update(entry3); err != nil {
		t.Errorf("testFileRepoIndexUpdate(): fileRepositoryIndex.update() failed: %s", err)
	}

	//fileRepoIdx.
	if indexEntries, err := fileRepoIdx.readIndexFile(); err != nil {
		t.Errorf("testFileRepoIndexUpdate(): readIndexFile() failed after update(): %v", err)
	} else {

		if len(indexEntries) != 3 {
			t.Errorf("testFileRepoIndexUpdate(): expected 3 entries, found %d", len(indexEntries))
		}

		// COMPARE THE INDEX ENTRIES RETRIEEVE FROM THE Index
		// WITH THE 3 USED TO UPDATE INDEX. COMPARE THOSE entries
		// WITH THE SAME UUID

		for _, e := range indexEntries {

			switch e.ObjectType {
			case jsonld.ContextResource:
				if e != entry {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry, e)
				}

			case jsonld.NamedGraphResource:
				if e != entry2 {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry2, e)
				}

			case jsonld.NodeResource:

				if e != entry3 {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry3, e)
				}

			}
		}

	}

}

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

	entry2.DatasetPath = "data/two"
	entry2.ObjectType = jsonld.NamedGraphResource
	entry2.ObjectIRI = ""
	entry2.ContainerType = jsonld.DatasetResource
	entry2.ContainerIRI = ""
	entry2.LastModifiedNs = 1600103677853633000
	entry2.ObjectCID = "" //"bafyreie5h75u3kv47vywgzsohnqlmxowfv4herxrj6ezdlttx47wrkpbkm"

	var entry3 = common.StagingResource{}

	entry3.DatasetPath = "data/is/three"
	entry3.ObjectType = jsonld.NodeResource
	entry3.ObjectIRI = "http://www.doggg3rz.com/my/test#me"
	entry3.ContainerType = jsonld.NamedGraphResource
	entry3.ContainerIRI = ""
	entry3.LastModifiedNs = 1600103677853633999
	entry3.ObjectCID = "bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"

	return entry, entry2, entry3

}

func testUpdateExistingAndReadBack(t *testing.T) {

	// ENTRY ALREADY EXISTS BUT WILL CHANGE SOME VALUES
	entry, entry2, entry3 := getThreeEntries()

	entry2.DatasetPath = "data/two"
	entry2.ObjectType = jsonld.NamedGraphResource
	entry2.ObjectIRI = ""
	entry2.ContainerType = jsonld.DatasetResource
	entry2.ContainerIRI = ""
	entry2.LastModifiedNs = time.Now().UnixNano()
	entry2.ObjectCID = "" //"bafyreidwx2fvfdiaox32v2mnn6sxu3j4qoxeqcuenhtgrv5qv6litfnmoe"

	if err := fileRepoIdx.update(entry2); err != nil {
		t.Errorf("testUpdateExistingAndReadBack(): %s", err)
	}

	if indexEntries, err := fileRepoIdx.readIndexFile(); err != nil {
		t.Errorf("testUpdateExistingAndReadBack(): %s", err)
	} else {
		if len(indexEntries) != 3 {
			t.Errorf("testUpdateExistingAndReadBack(): expected 3 entries after updating "+
				"existing entry. found %d", len(indexEntries))
		}

		for _, e := range indexEntries {
			if e.ObjectType == entry2.ObjectType {
				if e != entry2 {
					t.Errorf("testUpdateExistingAndReadBack(): updated index entry "+
						"changed after read: {expected: %v, found: %v } ", entry2, e)
				}
			} else if e.ObjectType == entry.ObjectType {
				if e != entry {
					t.Errorf("testUpdateExistingAndReadBack(): non-updated index entry "+
						"changed after read. { expected %v, found %v }", entry, e)
				}
			} else if e.ObjectType == entry3.ObjectType {
				if e != entry3 {
					t.Errorf("testUpdateExistingAndReadBack(): non-updated index entry "+
						"changed after read. { expected %v, found %v }", entry3, e)
				}
			}
		}
	}

}

func testInvalidIndexEntryValidate(t *testing.T) {

	entry, _, _ := getThreeEntries()

	// SET BAD TYPE. SELF ASSIGNED AND NOT ALLOCATED IN primitives PACKAGE
	//var badType string = "dogg3rz.badtype"
	var badObjectType jsonld.JSONLDResourceType = math.MaxInt8

	// TEST BAD OBJECT TYPE
	holdType := entry.ObjectType
	entry.ObjectType = badObjectType

	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Type value assigned")
	}
	entry.ObjectType = holdType

	// TEST BAD DATASETPATH
	holdDatasetPath := entry.DatasetPath
	entry.DatasetPath = "%badPathElement/$foo"
	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Type value assigned")
	}
	entry.DatasetPath = holdDatasetPath

}

func testNewFileRepoIndexOnNonExistentRepo(t *testing.T) {

	var nonExistRepo = "not." + testRepoName

	ctxt := getContext()

	if _, err := newFileRepositoryIndex(ctxt, nonExistRepo); err == nil {
		t.Errorf("testNewFileRepoIndexOnNonExistentRepo(): did not fail "+
			"on non-existent repository: %s", nonExistRepo)
	}

}

func testReadIndexFileFailsOnUpdate(t *testing.T) {

	// MOVE REPO DIR AND THEN CALL UPDATE
	moveDir := fileRepoIdx.repoDir + ".move"
	if err := os.Rename(fileRepoIdx.repoDir, moveDir); err != nil {
		t.Fail()
	}

	e, _, _ := getThreeEntries()

	if err := fileRepoIdx.update(e); err == nil {
		t.Errorf("testReadIndexFileFailsOnUpdate(): did not fail "+
			"when repo dir moved to %s", moveDir)
	}
	// MOVE REPO DIR BACK

	if err := os.Rename(moveDir, fileRepoIdx.repoDir); err != nil {
		t.Fail()
	}

}
