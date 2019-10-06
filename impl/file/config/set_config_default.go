package config

import (
	"io/ioutil"
	"os"

	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"
)

func SetConfigDefault() error {

	err := ioutil.WriteFile(configPath(), []byte(resourceconfig.CONFIG_JSON_DEFAULT), os.FileMode(0660))

	return err

}
