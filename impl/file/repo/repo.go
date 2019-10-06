package repo

import (
	"os"
	"path"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

type FileRepositoryResource struct {
}

func (repo *FileRepositoryResource) InitRepo(name string) error {

	repoDir := path.Join(file.RepositoriesDirPath(), name)

	err := os.Mkdir(repoDir, os.FileMode(0700))
	if err != nil {
		if os.IsNotExist(err) {
			// BASE REPO DIR DOES NOT EXIST
			return dgrzerr.NotFound.Wrapf(err, file.RepositoriesDirPath())
		}
		if os.IsExist(err) {
			return dgrzerr.AlreadyExists.Wrapf(err, name)
		}
		return err
	}

	return nil

}
