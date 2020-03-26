package repo

//	"os"

//	"bytes"
//	"crypto/sha1"
//	"encoding/binary"

//	"github.com/datacequia/go-dogg3rz/errors"
//	"github.com/datacequia/go-dogg3rz/impl/file"
//	rescom "github.com/datacequia/go-dogg3rz/resource/common"
//	"github.com/datacequia/go-dogg3rz/util"
import (
	"io"
	"os"
	"path/filepath"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/google/uuid"

	"github.com/datacequia/go-dogg3rz/primitives"
)

type fileCreateSchema struct {
	repoName       string
	schemaSubpath  common.RepositoryPath
	fileSystemPath string
}

//

func (cs *fileCreateSchema) createSchema(repoName string, schemaSubpath string, schemaReader io.Reader) error {

	rp, err := common.RepositoryPathNew(schemaSubpath)
	if err != nil {
		return err
	}

	if rp.EndsWithPathSeparator() {
		return errors.InvalidValue.Newf("path to schema object cannot end "+
			"with path separater, found %s",
			schemaSubpath)
	}

	if schemaPath, err := cs.createRepositoryResourcePath(
		rp, repoName, primitives.TYPE_DOGG3RZ_SCHEMA,
		schemaReader); err != nil {
		return err
	} else {

		cs.fileSystemPath = schemaPath

	}

	cs.repoName = repoName
	cs.schemaSubpath = *rp

	return nil

}

// CREATES THE RESOURCE PATH IN THE DESIGNATED REPOSITORY OF A SPECIFIC
// RESOURCE TYPE
func (cs *fileCreateSchema) createRepositoryResourcePath(
	resPath *rescom.RepositoryPath,
	repoName string,
	resType primitives.Dogg3rzObjectType,
	bodyReader io.Reader) (string, error) {

	if !file.RepositoryExist(repoName) {
		return "", errors.NotFound.Newf("repository '%s' does not exist. please create it first", repoName)
	}
	// REPO DOES EXIST. CREATE EACH PATH ELEMENT IF NECESSARY
	curPath := filepath.Join(file.RepositoriesDirPath(), repoName)
	curResType := primitives.TYPE_DOGG3RZ_TREE
	mkDirCount := 0
	success := false

	removePathOnFail := func(path string) {
		if !success {
			os.Remove(path)
		}
	}

	for pathElementIndex, path := range resPath.PathElements() {

		curPath = filepath.Join(curPath, path)

		lastPathElement := (pathElementIndex == (resPath.Size() - 1))

		if lastPathElement {
			// LAST ELEMENT. MAKE CUR RESOURCE TYPE
			// THE DESIRED RESOURCE TYPE
			curResType = resType

		}

		// EVAL CURRENT REPO PATH TO ENSURE IT'S A DIRECTORY AND
		// A DOGG3RZ TREE OBJECT
		if _, err := os.Stat(curPath); err != nil {

			if os.IsNotExist(err) {

				// CURRENT PATH DOES NOT EXIST. CREATE IT
				if err := os.Mkdir(curPath, os.FileMode(0700)); err != nil {
					return "", err
				} else {
					defer removePathOnFail(curPath)
				}

				mkDirCount++

				// CREAT EXCLUSIVE LOCK ON THIS DIR WHILE
				// THIS FUNCTION IS RUNNING. THIS WILL DISALLOW
				// ANY READERS/WRITERS THAT TRAVERSE THROUGH THIS
				// PATH WHILE THIS METHOD RUNS AND CREATES RESOURCES
				if mkDirCount == 1 {

					dirLockFilePath := filepath.Join(curPath, file.DirLockFileName)

					if lockFile, err := os.OpenFile(dirLockFilePath,
						os.O_RDWR|os.O_CREATE|os.O_EXCL,
						os.FileMode(0600)); err != nil {
						if os.IsExist(err) {
							return "", errors.TryAgain.Wrapf(err,
								"resource %s is temporarily unavailable. try operation again later...",
								curPath)
						}
						return "", err

					} else {
						// LOCK FILE CREATED! RELEASE FILE HANDLE
						lockFile.Close()

						defer os.Remove(dirLockFilePath)
					}

				}

				// CREATE ID ATTR BEFORE TYPE
				// SO IF TYPE ATTR FILE FAILS TO BE CREATED
				//  IT WON'T BE
				// RECOGNIZED AS A RESOURCE DIR WHEN
				// SUBSEQUENTLY EVALUATED
				if attrPath, err := file.PutResourceAttributeS(curPath,
					primitives.DOGG3RZ_OBJECT_ATTR_ID,
					uuid.New().String()); err != nil {
					return "", err
				} else {
					defer removePathOnFail(attrPath)
				}

				//	fmt.Println("1")

				if attrPath, err := file.PutResourceAttributeS(curPath,
					primitives.DOGG3RZ_OBJECT_ATTR_TYPE,
					curResType.String()); err != nil {
					return "", err
				} else {
					defer removePathOnFail(attrPath)
				}
				//	fmt.Println("2")
			} else {
				// SOME OTHER (SYSTEM?) ERROR OCCURRED. RETURN IT
				return "", err
			}
		} else {
			// PATH EXISTS
			if lastPathElement {
				return "", errors.AlreadyExists.Newf("resource path already exists: %s",
					curPath)
			}

			if _, err := os.Stat(filepath.Join(curPath, file.DirLockFileName)); err == nil {
				// RESOURCE IS LOCKED.

				return "", errors.TryAgain.Newf(
					"resource %s is temporarily unavailable. try operation again later...",
					curPath)

			}

			// IS IT A DOGG3RZ TREE OBJECT?

			dgrzType, err := file.GetResourceAttributeS(curPath,
				primitives.DOGG3RZ_OBJECT_ATTR_TYPE)
			if err != nil {
				return "", err
			}

			//	fmt.Println("3")
			// READ CONTENT. IS IT A DOGGERZ TREE OBJECT?
			if dgrzType != primitives.TYPE_DOGG3RZ_TREE.String() {

				if pathElementIndex == (resPath.Size() - 1) {
					// THIS IS LAST PATH ELEMENT AND A RESOURCE
					// ALREADY EXISTS HERE
					return "", errors.AlreadyExists.Newf(
						"%s: %s = %s",
						resPath.ToString(),
						primitives.DOGG3RZ_OBJECT_ATTR_TYPE,
						dgrzType)

				} else {

					return "", errors.InvalidPathElement.Newf(
						"encountered invalid base path '%s' during creation of "+
							"repository resource '%s': want %s '%s', found %s '%s'",
						curPath,
						resPath.ToString(),
						primitives.DOGG3RZ_OBJECT_ATTR_TYPE,
						primitives.TYPE_DOGG3RZ_TREE,
						primitives.DOGG3RZ_OBJECT_ATTR_TYPE,
						dgrzType)

				}

				// IS A TREE OBJECT. GTG...
			}

		}
	}
	//	fmt.Println("body1")
	if _, err := file.PutResourceAttribute(
		curPath,
		primitives.DOGG3RZ_OBJECT_ATTR_BODY,
		bodyReader); err != nil {

		return "", err
	}

	success = true

	return curPath, nil

}
