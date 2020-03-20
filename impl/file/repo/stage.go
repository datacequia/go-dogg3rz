package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/primitives"
	"github.com/datacequia/go-dogg3rz/resource/common"

	shell "github.com/ipfs/go-ipfs-api"
)

type fileStageResource struct {
	repoName       string
	schemaSubpathS string
	repositoryPath *common.RepositoryPath
	filesystemPath string // filesystem path to resource

	attrType         string
	attrId           string
	attrMeta         map[string]interface{}
	attrBodyCid      string
	attrBodyFileInfo os.FileInfo
	resourceContent  string
	resourceCid      string
}

var stageableResources = []primitives.Dogg3rzObjectType{
	primitives.TYPE_DOGG3RZ_SCHEMA,
	primitives.TYPE_DOGG3RZ_MEDIA,
	primitives.TYPE_DOGG3RZ_OBJECT,
	primitives.TYPE_DOGG3RZ_TRIPLE,
	primitives.TYPE_DOGG3RZ_SERVICE,
}

//
func (s *fileStageResource) stageResource(repoName string, schemaSubpath string) error {

	if !file.RepositoryExist(repoName) {
		return errors.NotFound.Newf("repository '%s' does not exist. please create it first.", repoName)
	}

	s.repoName = repoName
	s.schemaSubpathS = schemaSubpath

	rp, err := common.RepositoryPathNew(schemaSubpath)
	if err != nil {
		return err
	} else {
		s.repositoryPath = rp
	}

	if s.repositoryPath.EndsWithPathSeparator() {
		return errors.InvalidValue.Newf("path to resource object to be staged cannot end "+
			"with path separater, found %s",
			schemaSubpath)
	}

	s.filesystemPath = filepath.Join(file.RepositoriesDirPath(), repoName, rp.ToString())

	if err := s.loadResourceAttributes(s.filesystemPath); err != nil {
		return err
	}

	// CREATE DOGG3RZ OBJECT
	resource := make(map[string]interface{})
	//fmt.Println("create resource ")
	resource[primitives.DOGG3RZ_OBJECT_ATTR_TYPE] = s.attrType
	resource[primitives.DOGG3RZ_OBJECT_ATTR_ID] = s.attrId
	resource[primitives.DOGG3RZ_OBJECT_ATTR_METADATA] = s.attrMeta
	//	fmt.Println("putted resource ")
	mapBody := make(map[string]interface{})
	mapBody["/"] = s.attrBodyCid // POINT TO BODY CONTENT USING IPFS CID

	resource[primitives.DOGG3RZ_OBJECT_ATTR_BODY] = mapBody

	if resBytes, err := json.Marshal(resource); err != nil {
		return err
	} else {

		sh := shell.NewShell("localhost:5001")
		//	fmt.Println(s)
		if cid, err := sh.DagPut(resBytes, "json", "cbor"); err != nil {
			//	fmt.Printf("dagput err: %s: %s", err, string(resBytes))
			return err
		} else {
			//	fmt.Println("resourceCid", cid)
			s.resourceCid = cid

		}
	}

	var entry indexEntry

	entry.Type = s.attrType
	entry.Uuid = s.attrId
	entry.FileSize = s.attrBodyFileInfo.ModTime().UnixNano()
	entry.MtimeNs = s.attrBodyFileInfo.Size()
	entry.Multihash = s.resourceCid
	entry.Subpath = rp.ToString()
	//fmt.Println(entry.String())
	if fileRepoIdx, err := newFileRepositoryIndex(repoName); err != nil {
		return err
	} else {
		if err := fileRepoIdx.update(entry); err != nil {
			return err
		}
	}

	// PRINT THE MULTIHASH OF THE RESOURCE THAT WAS STAGED (FOR NOW)
	fmt.Println(s.resourceCid)

	return nil

}

func (s *fileStageResource) String() string {

	return fmt.Sprintf("fileStageResource = { "+
		"repoName: %v, "+
		"schemaSubpathS: %v, "+
		"repositoryPath: %v, "+
		"filesystemPath: %v, "+
		"attrType: %v, "+
		"attrId: %v, "+
		"attrMeta: %v, "+
		"attrBodyCid: %v, "+
		"resourceCid: %v }",
		s.repoName, s.schemaSubpathS, s.repositoryPath,
		s.filesystemPath, s.attrType, s.attrId,
		s.attrMeta, s.attrBodyCid, s.resourceCid)

}

func isRepositoryStageable(resType primitives.Dogg3rzObjectType) bool {

	for _, rt := range stageableResources {

		if resType == rt {
			return true
		}
	}
	return false

}

func (s *fileStageResource) loadResourceAttributes(resPath string) error {

	//	var resType string

	if attrValue, err := file.GetResourceAttributeS(resPath,
		primitives.DOGG3RZ_OBJECT_ATTR_TYPE); err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound.Wrapf(err, "repository resource attribute not found")
		}
		return err
	} else {
		s.attrType = attrValue
	}

	if dt, err := primitives.Dogg3rzObjectTypeFromString(s.attrType); err != nil {

		if !isRepositoryStageable(dt) {
			return errors.UnexpectedType.Newf(
				"expected repository resource type to be one of %v, found '%s'",
				stageableResources, s.attrType)
		}
	}

	if attrValue, err := file.GetResourceAttributeS(resPath,
		primitives.DOGG3RZ_OBJECT_ATTR_ID); err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound.Wrapf(err, "repository resource attribute not found")
		}
		return err
	} else {
		s.attrId = attrValue
	}

	if attrValue, err := file.GetResourceAttributeS(resPath,
		primitives.DOGG3RZ_OBJECT_ATTR_METADATA); err != nil {
		if !os.IsNotExist(err) {
			// METADATA IS OPTIONAL. ASSIGN EMPTY JSON ObjectType
			// ONLY RETURN ERROR IF IT'S NOT A 'NOTEXIST' ERROR

			return err
		}
		// MAKE EMPTY OBJECT
		s.attrMeta = make(map[string]interface{})
	} else {

		if err := json.Unmarshal([]byte(attrValue), &s.attrMeta); err != nil {
			return errors.UnexpectedValue.Wrapf(err, "expected JSON map when loading resource %s attribute",
				primitives.DOGG3RZ_OBJECT_ATTR_METADATA)

		}

	}
	//fmt.Println("call GetResourceAttributeCB begin...")
	if err := file.GetResourceAttributeCB(resPath,
		primitives.DOGG3RZ_OBJECT_ATTR_BODY,
		s.loadBody); err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound.Wrapf(err, "repository resource attribute '%s' not found", primitives.DOGG3RZ_OBJECT_ATTR_BODY)
		}
		return err
	}
	//fmt.Println("call GetResourceAttributeCB end")

	return nil
}

func (s *fileStageResource) expectJSONBody() bool {

	if t, err := primitives.Dogg3rzObjectTypeFromString(s.attrType); err == nil {

		switch t {
		case primitives.TYPE_DOGG3RZ_SCHEMA:
			return true
		case primitives.TYPE_DOGG3RZ_OBJECT:
			return true
		case primitives.TYPE_DOGG3RZ_MEDIA:
			return false
		case primitives.TYPE_DOGG3RZ_TRIPLE:
			return true
		case primitives.TYPE_DOGG3RZ_SERVICE:
			return true
		default:
			// SHOULD NOT HAVE GOTTEN THIS FAR. ONLY ABOVE RESOURCE TYPES
			// CAN BE STAGED
			panic(fmt.Sprintf("expectJSONBody(): "+
				"unexpected Dogg3rzObjectType value encountered: '%s': { Dogg3zObjectType = %s }",
				s.attrType, t))
		}

	} else {

		panic(fmt.Sprintf("expectJSONBody(): "+
			"unexpected Dogg3rzObjectType value encountered: '%s': %s",
			s.attrType, err))
	}

}

func (s *fileStageResource) loadBody(bodyReader io.Reader, fileInfo os.FileInfo) error {

	var cid string
	var err error
	//fmt.Println("loadBody start")
	sh := shell.NewShell("localhost:5001")

	if s.expectJSONBody() {
		//fmt.Println("loadBody: before dag put")
		cid, err = sh.DagPut(bodyReader, "json", "cbor")
		//fmt.Println("loadBody: after dag put")
	} else {
		cid, err = sh.Add(bodyReader)
	}

	if err == nil {
		s.attrBodyCid = cid
		s.attrBodyFileInfo = fileInfo
	}

	return err
}
