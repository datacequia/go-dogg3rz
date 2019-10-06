package node

import (
	"os"

	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/impl/file/config"
)

type FileNodeResource struct {
}

func (node *FileNodeResource) InitNode() error {

	file.DotDirPath()

	createDirList := []string{file.DotDirPath(), file.DataDirPath(), file.RepositoriesDirPath()}

	for _, d := range createDirList {
		// CREATE DIR SO THAT ONLY USER CAN R/W
		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			return err
		}
	}

	err := config.SetConfigDefault()
	if err != nil {
		return err
	}

	return nil

}
