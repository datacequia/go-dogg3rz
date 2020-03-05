/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package primitives

import (
	"encoding/json"
	"io"

	"github.com/datacequia/go-dogg3rz/errors"
)

//const D_ATTR_NAME = "objectHeads"
const TYPE_DOGG3RZ_COMMIT = "dogg3rz.commit"

//const MD_ATTR_NAME = "name"
//const MD_ATTR_IPFS_PEER_ID = "ipfsPeerId"
const MD_ATTR_EMAIL_ADDR = "emailAddress"
const MD_ATTR_REPO_NAME = "repositoryName"
const MD_ATTR_REPO_ID = "repositoryId"

//const D_ATTR_ROOT_TREE = "rootTree"
const D_ATTR_TRIPLES = "triples"
const D_ATTR_IMPORTS = "imports"

var reservedMDAttrCommit = [...]string{MD_ATTR_REPO_NAME, MD_ATTR_EMAIL_ADDR}

var reservedDAttrCommit = [...]string{D_ATTR_TRIPLES, D_ATTR_IMPORTS}

type dgrzCommit struct {
	repositoryId string // GLOBALLY UNIQUE IMMUTABLE IDENTIFIER FOR REPOSITORY
	// THAT THIS COMMIT BELONGS TO.
	emailAddress string // EMAIL ADDRESS OF COMMITTER

	repositoryName string // REPOSITORY NAME FOR THIS COMMIT

	parents []string
}

func Dogg3rzCommitNew(repoName string, repoId string, emailAddress string, parentCommits []string) (*dgrzCommit, error) {

	c := &dgrzCommit{repositoryName: repoName, emailAddress: emailAddress,
		repositoryId: repoId, parents: parentCommits}

	return c, nil
}

// Return a dgrzCommit object fom the Reader
func Deserialize(reader io.Reader) (*dgrzCommit, error) {

	// CONVERT FROM DOGG3RZOBJECT TO COMMITOBJECT
	dgrzObj, err := Dogg3rzObjectDeserializeFromJson(reader)

	if dgrzObj.ObjectType != TYPE_DOGG3RZ_COMMIT {
		return nil, errors.UnexpectedType.Newf("expected dogg3rz type '%s', found '%s'",
			TYPE_DOGG3RZ_COMMIT, dgrzObj.ObjectType)

	}

	var (
		repoName     string
		emailAddress string
		repoId       string
		parents      []string
	)

	if len(dgrzObj.Parents) > 0 {
		parents = make([]string, len(dgrzObj.Parents))
		copy(parents, dgrzObj.Parents)

	}
	// FETCH  PEER ID FROM METADATA SECTION
	if val, ok := dgrzObj.Metadata[MD_ATTR_REPO_NAME]; !ok {
		return nil, errors.NotFound.Newf("metadata attribute value '%s' not found",
			MD_ATTR_REPO_NAME)
	} else {
		repoName = val
	}
	// FETCH  REPO ID FROM METADATA SECTION
	if val, ok := dgrzObj.Metadata[MD_ATTR_REPO_ID]; !ok {
		return nil, errors.NotFound.Newf("metadata attribute value '%s' not found",
			MD_ATTR_REPO_ID)
	} else {
		repoId = val
	}

	// FETCH COMMITTER'S EMAIL ADDRESS FROM METADATA SECTION
	if val, ok := dgrzObj.Metadata[MD_ATTR_EMAIL_ADDR]; !ok {
		return nil, errors.NotFound.Newf("metadata attribute value '%s' not found",
			MD_ATTR_EMAIL_ADDR)
	} else {
		emailAddress = val
	}

	// GET NAME OF REPO FROM ROOT TREE OBJ
	//	if val,ok := rootTree[DOGG3RZ_OBJECT_ATTR_METADATA]; !ok {
	//		return nil, errors.NotFound.Newf("data attribute")
	//	} else {

	//	}

	// CREATE COMMIT OBJECT FROM ATTRIBUTES
	// EXTRACTEED FROM DESERIALIZED DOGG3RZ Object

	commitObj, err := Dogg3rzCommitNew(repoName, repoId, emailAddress, parents)
	if err != nil {
		return nil, err
	}

	return commitObj, nil
}

func Serialize(commit *dgrzCommit, writer io.Writer) error {

	encoder := json.NewEncoder(writer)
	err := encoder.Encode(commit.ToDogg3rzObject())
	if err != nil {
		return err
	}

	return nil

}

func (receiver *dgrzCommit) ToDogg3rzObject() *dgrzObject {

	o := Dogg3rzObjectNew(TYPE_DOGG3RZ_COMMIT)

	o.Metadata[MD_ATTR_REPO_ID] = receiver.repositoryId
	o.Metadata[MD_ATTR_REPO_NAME] = receiver.repositoryName
	o.Metadata[MD_ATTR_EMAIL_ADDR] = receiver.emailAddress

	o.Parents = make([]string, len(receiver.parents))
	copy(o.Parents, receiver.parents)
	//	o.Parent = receiver.parent

	return o
}
