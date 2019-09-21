package primitives

import (
	"encoding/json"
	"io"
  "github.com/datacequia/go-dogg3rz/errors"
	"github.com/fatih/structs"
	"reflect"
)

//const D_ATTR_NAME = "objectHeads"
const TYPE_DOGG3RZ_COMMIT = "dogg3rz.commit"

//const MD_ATTR_NAME = "name"
const MD_ATTR_IPFS_PEER_ID = "ipfsPeerId"
const MD_ATTR_EMAIL_ADDR = "emailAddress"

const D_ATTR_ROOT_TREE = "rootTree"
const D_ATTR_TRIPLES = "triples"
const D_ATTR_IMPORTS = "imports"

var reservedMDAttrCommit = [...]string{MD_ATTR_IPFS_PEER_ID, MD_ATTR_EMAIL_ADDR}

var reservedDAttrCommit = [...]string{D_ATTR_ROOT_TREE, D_ATTR_TRIPLES, D_ATTR_IMPORTS}

type dgrzCommit struct {
	IpfsPeerId   string
	emailAddress string
	rootTree     *dgrzTree

	parent string
}

func Dogg3rzCommitNew(repoName string, ipfsPeerId string, emailAddress string) (*dgrzCommit, error) {

	rootTree, err := Dogg3rzTreeNew(repoName)
	if err != nil {
		return nil, err
	}

	c := &dgrzCommit{IpfsPeerId: ipfsPeerId, emailAddress: emailAddress, rootTree: rootTree}

	return c, nil
}

// Return a dgrzCommit object fom the Reader
func Deserialize(reader io.Reader) (*dgrzCommit, error) {

	// CONVERT FROM DOGG3RZOBJECT TO COMMITOBJECT
	dgrzObj, err := Dogg3rzObjectDeserializeFromJson(reader)

	if dgrzObj.ObjectType != TYPE_DOGG3RZ_COMMIT {
		return nil, errors.UnexpectedType.Newf("expected dogg3rz type '%s', found '%s'",
			 TYPE_DOGG3RZ_COMMIT,dgrzObj.ObjectType)

	}

	var (
		ipfsPeerId string
		emailAddr string
		rootTree map[string]interface{}
	)

	// FETCH  PEER ID FROM METADATA SECTION
	if val,ok := dgrzObj.Metadata[MD_ATTR_IPFS_PEER_ID]; !ok {
		return nil, errors.NotFound.Newf("metadata attribute value '%s' not found",
			MD_ATTR_IPFS_PEER_ID)
	} else {
		ipfsPeerId = val
	}

	// FETCH COMMITTER'S EMAIL ADDRESS FROM METADATA SECTION
	if val,ok := dgrzObj.Metadata[MD_ATTR_EMAIL_ADDR]; !ok {
		return nil, errors.NotFound.Newf("metadata attribute value '%s' not found",
			MD_ATTR_EMAIL_ADDR)
	} else {
		emailAddr = val
	}

	// FETCH ROOT TREE OBJECT FROM DATA SECTION
	if val,ok := dgrzObj.Data[D_ATTR_ROOT_TREE]; !ok {
		return nil, errors.NotFound.Newf("data attribute value '%s' not found")
	} else {
		if val2,ok := val.(map[string]interface{}); !ok {
			return nil,errors.UnexpectedType.Newf("expected type %v for data attribute %s, found '%s'",
			reflect.TypeOf(rootTree),D_ATTR_ROOT_TREE, reflect.TypeOf(val))
		} else {
			rootTree = val2
		}

	}

	// GET NAME OF REPO FROM ROOT TREE OBJ
//	if val,ok := rootTree[DOGG3RZ_OBJECT_ATTR_METADATA]; !ok {
//		return nil, errors.NotFound.Newf("data attribute")
//	} else {

//	}

	// CREATE COMMIT OBJECT FROM ATTRIBUTES
	// EXTRACTEED FROM DESERIALIZED DOGG3RZ Object

	commitObj := Dogg3rzCommitNew(repoName, ipfsPeerId, emailAddress)


	// GET METADATA ATTR_OBJECT
	commitObj.IpfsPeerId, ok := dgrzObj.Metadata[D_ATTR_IPFS_PEER_ID]
	if ! ok {
		return errors.NotFound.Newf("")
	}
	rootTree,ok := dgrzObj.Data[D_ATTR_ROOT_TREE]
	if ! ok {

	}


	return nil, nil
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

	o.Metadata[MD_ATTR_IPFS_PEER_ID] = receiver.IpfsPeerId
	o.Metadata[MD_ATTR_EMAIL_ADDR] = receiver.emailAddress

	o.Data[D_ATTR_ROOT_TREE] = structs.Map(receiver.rootTree.Dogg3rzObject())

	return o
}
