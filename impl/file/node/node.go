package node

import (
	"os"
	"path"

	"github.com/datacequia/go-dogg3rz/impl/file"
)

type FileNodeResource struct {
}

func (node *FileNodeResource) InitNode() error {

	file.DotDirPath()

	dataDir := path.Join(file.DotDirPath(), "data")
	repoDir := path.Join(dataDir, "repositories")

	createDirList := []string{file.DotDirPath(), dataDir, repoDir}

	for _, d := range createDirList {
		// CREATE DIR SO THAT ONLY USER CAN R/W
		err := os.Mkdir(d, os.FileMode(0700))
		if err != nil {
			return err
		}
	}

	return nil

}
