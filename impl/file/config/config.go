/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"

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
	return path.Join(file.DotDirPath(), configFileName)

}

func validateConfig(path string) error {

	schemaLoader := gojsonschema.NewStringLoader(resourceconfig.CONFIG_JSON_SCHEMA)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	documentLoader := gojsonschema.NewReferenceLoader("file://" + absPath)

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
