package repo

import (
	//	"fmt"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"time"

	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	"github.com/datacequia/go-dogg3rz/primitives"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

var dogg3rzHome string
var fileRepoIdx *fileRepositoryIndex

const (
	testRepoName = "index_test"
)

func indexSetup(t *testing.T) {

	dogg3rzHome = filepath.Join(os.TempDir(),
		fmt.Sprintf("index_test_%d", time.Now().UnixNano()))

	os.Setenv("DOGG3RZ_HOME", dogg3rzHome)

	fileNodeResource := &filenode.FileNodeResource{}

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.Email = "test@dogg3rz.com"

	if err := fileNodeResource.InitNode(dgrzConf); err != nil {
		t.Error(err)
	}
	t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	fileRepositoryResource := FileRepositoryResource{}

	if err := fileRepositoryResource.InitRepo(testRepoName); err != nil {
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
	indexTeardown(t)

}

// ADD 3 NEW ENTRIES TO THE INDEX AND READ THEM
// BACK FROM INDEX AND COMPARE . SHOULD BE EXACTLY SAME
func testNewFileRepoIndex(t *testing.T) {

	// TEST NEW REPO INDEX WITH BAD REPO NAME

	if _, err := newFileRepositoryIndex(testRepoName + "-badName"); err == nil {
		t.Errorf("newFileRepositoryIndex(): succeeded with non-existent repo (name)")
	}

	if fri, err := newFileRepositoryIndex(testRepoName); err != nil {
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

			switch e.Uuid {
			case entry.Uuid:
				if e != entry {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry, e)
				}

			case entry2.Uuid:
				if e != entry2 {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry2, e)
				}

			case entry3.Uuid:

				if e != entry3 {
					t.Errorf("testFileRepoIndexUpdate(): "+
						"single updated entry retrieved != entry updated: { update entry = %s, retrieve entry = %s }",
						entry3, e)
				}

			}
		}

	}

}

func getThreeEntries() (indexEntry, indexEntry, indexEntry) {

	var entry = indexEntry{}

	entry.Type = primitives.TYPE_DOGG3RZ_SCHEMA.String()

	entry.Uuid = "cc424aad-0ff8-4d9d-b5d2-bf0c17db124c"
	entry.MtimeNs = time.Now().Unix()
	entry.FileSize = int64(1024)
	entry.Multihash = "bafyreigcm277jvvdmenqkudvan3mn7icvzdj2a3eygtgilkf2mypcrkgvi"
	entry.Subpath = "my/object"

	var entry2 = indexEntry{}
	entry2.Type = primitives.TYPE_DOGG3RZ_SCHEMA.String()

	entry2.Uuid = "5222cb2e-831f-4be8-8e88-8dc9c7fb8447"
	entry2.MtimeNs = time.Now().Unix()
	entry2.FileSize = int64(2048)
	entry2.Multihash = "bafyreie5h75u3kv47vywgzsohnqlmxowfv4herxrj6ezdlttx47wrkpbkm"
	entry2.Subpath = "my/object2"

	var entry3 = indexEntry{}
	entry3.Type = primitives.TYPE_DOGG3RZ_SCHEMA.String()
	entry3.Uuid = "d821e954-cab8-4829-9b3d-59581f5627a2"
	entry3.MtimeNs = time.Now().Unix()
	entry3.FileSize = int64(4096)
	entry3.Multihash = "bafyreiffn3ktxl4xdhtha4bvqx5ezganq5mwk4lbp3yhjyw7phle4kgc4m"
	entry3.Subpath = "my/object3"

	return entry, entry2, entry3

}

func testUpdateExistingAndReadBack(t *testing.T) {

	// ENTRY ALREADY EXISTS BUT WILL CHANGE SOME VALUES
	entry, entry2, entry3 := getThreeEntries()

	entry2.FileSize = 999
	entry2.MtimeNs = time.Now().Unix()
	entry2.Multihash = "bafyreihldlxvtuzd6wix4kyb2ijngpxwlrh5esifrf6aqrx2ijip6qo6fy"
	entry2.Subpath = "my/updated/object2"

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
			if e.Uuid == entry2.Uuid {
				if e != entry2 {
					t.Errorf("testUpdateExistingAndReadBack(): updated index entry "+
						"changed after read: {expected: %v, found: %v} ", entry2, e)
				}
			} else if e.Uuid == entry.Uuid {
				if e != entry {
					t.Errorf("testUpdateExistingAndReadBack(): non-updated index entry "+
						"changed after read. found %v", e)
				}
			} else if e.Uuid == entry3.Uuid {
				if e != entry3 {
					t.Errorf("testUpdateExistingAndReadBack(): non-updated index entry "+
						"changed after read. found %v", e)
				}
			}
		}
	}

}

func testInvalidIndexEntryValidate(t *testing.T) {

	entry, _, _ := getThreeEntries()

	// SET BAD TYPE. SELF ASSIGNED AND NOT ALLOCATED IN primitives PACKAGE
	var badType string = "dogg3rz.badtype"

	holdType := entry.Type
	entry.Type = badType

	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Type value assigned")
	}
	entry.Type = holdType

	holdFileSize := entry.FileSize
	entry.FileSize = -1

	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.FileSize value assigned")
	}

	entry.FileSize = holdFileSize

	holdMultihash := entry.Multihash
	entry.Multihash = "242323423asdfadfasdf"

	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Multihash value assigned")
	}

	entry.Multihash = holdMultihash

	holdSubpath := entry.Subpath

	entry.Subpath = "." + entry.Subpath
	if err := fileRepoIdx.update(entry); err == nil {
		t.Errorf("testInvalidIndexEntryValidate(): update did not fail on bad " +
			"indexEntry.Subpath value assigned")
	}

	entry.Subpath = holdSubpath

}

func testNewFileRepoIndexOnNonExistentRepo(t *testing.T) {

	var nonExistRepo = "not." + testRepoName

	if _, err := newFileRepositoryIndex(nonExistRepo); err == nil {
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
