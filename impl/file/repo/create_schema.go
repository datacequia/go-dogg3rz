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

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"

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
		return errors.InvalidArg.Newf("path to schema object cannot end "+
			"with path separater, found %s",
			schemaSubpath)
	}

	if schemaPath, err := file.CreateRepositoryResourcePath(
		rp, repoName, primitives.TYPE_DOGG3RZ_SCHEMA, schemaReader); err != nil {
		return err
	} else {

		cs.fileSystemPath = schemaPath

	}

	cs.repoName = repoName
	cs.schemaSubpath = *rp

	return nil

}
