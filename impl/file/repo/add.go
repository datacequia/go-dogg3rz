package repo

import (
	"context"

	"fmt"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

func add(ctxt context.Context, repoName string, path string) error {

	if !file.RepositoryExist(ctxt, repoName) {
		return errors.NotFound.Newf("repository '%s' does not exist",
			repoName)
	}

	repoFsPath := file.RepositoryDgrzDirPath(ctxt, repoName)

	fmt.Println(repoFsPath)

	return nil
}
