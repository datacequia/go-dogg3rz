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
	"fmt"
	"path/filepath"

	"github.com/datacequia/go-dogg3rz/errors"

	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/impl/file/config"
	"github.com/datacequia/go-dogg3rz/primitives"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
)

type fileCreateSnapshot struct {
	repoName    string
	snapshotMap map[uuid.UUID]snapshotResource

	fileRepoIdx  *fileRepositoryIndex
	indexEntries []rescom.StagingResource

	createTreePathElementContext []string
}

type snapshotIndexEntry struct {
	entry    *rescom.StagingResource
	repoPath *rescom.RepositoryPath
}

type snapshotResource struct {
	staging      stagingResource
	snapshotHead snapshotHeadResource
}

type workingTreeResource struct {
}

type stagingResource struct {
}

type snapshotHeadResource struct {
}

func (cs *fileCreateSnapshot) createSnapshot(ctxt context.Context, repoName string) error {

	if !file.RepositoryExist(repoName, ctxt) {
		return errors.NotFound.Newf("repository '%s' does not exist", repoName)
	}

	// REPO VALID. ASSIGN REPO NAME TO COMMAND STRUCT
	cs.repoName = repoName

	var ssIndexEntries *[]snapshotIndexEntry
	if i, err := cs.getIndexEntries(ctxt); err != nil {
		return err
	} else {
		ssIndexEntries = i
	}

	if rootTree, err := cs.createTree(nil, 0, ssIndexEntries, cs.repoName); err != nil {
		return err
	} else {

		buf := &bytes.Buffer{}

		var dgrzSnapshotObject *map[string]interface{}

		if p, err := createSnapshotObject(cs.repoName, rootTree, ctxt); err != nil {
			return err
		} else {
			dgrzSnapshotObject = p
		}

		e := json.NewEncoder(buf)
		if err := e.Encode(dgrzSnapshotObject); err != nil {
			return err
		}

		//	fmt.Printf("rootTree: %v", (*dgrzSnapshotObject))
		sh := shell.NewShell("localhost:5001")
		//	fmt.Println(s)
		if cid, err := sh.DagPut(buf, "json", "cbor"); err != nil {
			//	fmt.Printf("dagput err: %s: %s", err, string(resBytes))
			return err
		} else {

			if err := file.WriteCommitHashToCurrentBranchHeadFile(cs.repoName, cid, ctxt); err != nil {
				return err
			}

			fmt.Printf("%s\n", cid)
		}
	}

	return nil
}

func (cs *fileCreateSnapshot) getIndexEntries(ctxt context.Context) (*[]snapshotIndexEntry, error) {

	var fileRepoIdx *fileRepositoryIndex
	if i, err := newFileRepositoryIndex(cs.repoName, ctxt); err != nil {
		return nil, err
	} else {
		fileRepoIdx = i
	}

	var indexEntries []rescom.StagingResource
	if ie, err := fileRepoIdx.readIndexFile(); err != nil {
		return nil, err
	} else {
		indexEntries = ie
	}

	ssIndexEntries := make([]snapshotIndexEntry, len(indexEntries))

	// commented out temporaarily until fix
	/*
		for i := 0; i < len(indexEntries); i++ {
			var err error
			ssIndexEntries[i].entry = &indexEntries[i]
			ssIndexEntries[i].repoPath, err = rescom.RepositoryPathNew(indexEntries[i].Subpath)
			if err != nil {
				return nil, err
			}
		}
	*/

	return &ssIndexEntries, nil

}

func (cs *fileCreateSnapshot) createTree(parent *map[string]interface{},
	level int, pathList *[]snapshotIndexEntry, attrName string) (*map[string]interface{}, error) {

	const typeAttrName = "." + primitives.DOGG3RZ_OBJECT_ATTR_TYPE

	cs.createTreePathElementContext = append(cs.createTreePathElementContext, attrName)

	popPathElementFunc := func() {
		// POP TREE PATH ELEMENT
		lengthSlice := len(cs.createTreePathElementContext)
		if lengthSlice > 0 {
			cs.createTreePathElementContext = cs.createTreePathElementContext[:lengthSlice-1]
		}
	}

	defer popPathElementFunc()

	if parent == nil {
		// CREATE PARENT IF NOT PASSED IN
		m := make(map[string]interface{})
		parent = &m
		//	fmt.Printf("parent is: %v %v %s\n", parent, *parent, typeAttrName)
		(*parent)[typeAttrName] = primitives.TYPE_DOGG3RZ_TREE.String()
		//fmt.Printf("parent is (after): %v %v %s\n", parent, *parent, typeAttrName)
	}

	var createdTree *map[string]interface{}

	if m, ok := getMapValueFromKey(parent, attrName); !ok {
		// PARENT DOES NOT HAVE NAME
		// CREATE IT
		theMap := make(map[string]interface{})
		createdTree = &theMap
		(*createdTree)[typeAttrName] = primitives.TYPE_DOGG3RZ_TREE.String()
		(*parent)[attrName] = *createdTree

	} else {
		createdTree = m

		if attrValue, ok := getStringValueFromKey(createdTree, typeAttrName); ok && attrValue != primitives.TYPE_DOGG3RZ_TREE.String() {
			return nil, errors.AlreadyExists.Newf("path %s already exists and is not "+
				"dogg3rz type '%s': found '%s'",
				filepath.Join(cs.createTreePathElementContext...),
				primitives.TYPE_DOGG3RZ_TREE.String(), attrValue)
		}

	}

	/* commented out temporarily untix fix
	for _, entry := range *pathList {
		entryPathElements := entry.repoPath.PathElements()
		if level < entry.repoPath.Size() {
			if level == entry.repoPath.Size()-1 {
				// THIS IS A LEAF ELEMENT (A NON TREE DOGG3RZ OBJECT)

				if _, ok := (*createdTree)[entryPathElements[level]]; ok {
					return nil, errors.AlreadyExists.Newf("path element %s already exists",
						filepath.Join(filepath.Join(cs.createTreePathElementContext...),
							entryPathElements[level]))
				}
				(*createdTree)[entryPathElements[level]] = map[string]string{"/": entry.entry.Multihash}
			} else {
				// level < entry.RepoPath.Size()
				var entryListOfOne []snapshotIndexEntry = []snapshotIndexEntry{entry}

				if _, err := cs.createTree(createdTree,
					level+1, &entryListOfOne, entryPathElements[level]); err != nil {
					return nil, err
				}
			}
		}
	}
	*/
	return parent, nil
}

func getMapValueFromKey(m *map[string]interface{}, key string) (*map[string]interface{}, bool) {

	if value, ok := (*m)[key]; !ok {
		return nil, false
	} else {
		if v, ok := value.(map[string]interface{}); ok {
			return &v, true
		}
		panic(fmt.Sprintf("expected type map[string]interface{} found type %T",
			value))

	}

}

func getStringValueFromKey(m *map[string]interface{}, key string) (string, bool) {

	if value, ok := (*m)[key]; !ok {
		return "", false
	} else {
		if v, ok := value.(string); ok {
			return v, true
		} else {
			panic(fmt.Sprintf("expected type string found type %T", value))

		}

	}
}

func createSnapshotObject(repoName string, rootTree *map[string]interface{}, ctxt context.Context) (*map[string]interface{}, error) {

	dogg3rzObject := make(map[string]interface{})

	var cr config.FileConfigResource
	var err error

	c, err := cr.GetConfig(ctxt)
	if err != nil {
		return &dogg3rzObject, err
	}

	meta := make(map[string]string)

	meta[primitives.META_ATTR_REPO_NAME] = repoName
	//meta[primitives.META_ATTR_REPO_ID] = ""
	meta[primitives.META_ATTR_EMAIL_ADDR] = c.User.Email

	body := make(map[string]interface{})

	body[primitives.BODY_ATTR_ROOT_TREE] = *rootTree

	dogg3rzObject[primitives.DOGG3RZ_OBJECT_ATTR_TYPE] = primitives.TYPE_DOGG3RZ_SNAPSHOT
	dogg3rzObject[primitives.DOGG3RZ_OBJECT_ATTR_METADATA] = meta
	dogg3rzObject[primitives.DOGG3RZ_OBJECT_ATTR_BODY] = body

	//	m[primitives.DOGG3RZ_OBJECT_ATTR_METADATA]
	return &dogg3rzObject, nil
}
