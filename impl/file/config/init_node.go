package config

import (
	"fmt"
	"os"
)

const DOT_DIR_NAME = ".dogg3rz"

func InitNode() error {

	// make dot dgrzInitRepo
	var userHomeDir string
	var err error

	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.Mkdir(fmt.Sprintf("%s%c%s", userHomeDir, os.PathSeparator, DOT_DIR_NAME), 0700)
	if err != nil {
		return err
	}

	return nil

}
