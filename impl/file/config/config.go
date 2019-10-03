package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"

	//	"github.com/datacequia/go-dogg3rz/impl/config"
	//	"github.com/datacequia/go-dogg3rz/impl/config"

	//"github.com/datacequia/go-dogg3rz/impl/config/file"
	"github.com/xeipuuv/gojsonschema"
)

const configFileName = "config"

type FileConfigResource struct {
}

func (configResource *FileConfigResource) GetConfig() (*resourceconfig.Dogg3rzConfig, error) {

	err := validateConfig(configPath())
	if err != nil {
		return nil, err
	}

	dgrzCfg := &resourceconfig.Dogg3rzConfig{}

	byteValue, err := ioutil.ReadFile(configPath())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &dgrzCfg)
	if err != nil {
		return nil, err
	}

	return dgrzCfg, nil

}

// Returns Path to dogg3rz configuration file
func configPath() string {

	// gojsonschema reequires the schema 'file://' prepeended to path
	return path.Join("file://", file.DotDirPath(), configFileName)

}

func validateConfig(path string) error {

	schemaLoader := gojsonschema.NewStringLoader(resourceconfig.CONFIG_JSON_SCHEMA)
	documentLoader := gojsonschema.NewReferenceLoader(path)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {

		return err

	}

	if !result.Valid() {

		log.Printf("configuration is not valid: %s:\n", path)

		s := make([]string, len(result.Errors()))

		for i, desc := range result.Errors() {
			s[i] = desc.Description()
		}

		return dgrzerr.ConfigError.Newf("schema validation failed: %s: %s", path, strings.Join(s, ": "))

	}

	return nil

}
