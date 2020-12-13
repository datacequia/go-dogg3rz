/*
 * Copyright (c) 2019-2020 Datacequia LLC. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

package config

import (
	"context"
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

func (configResource *FileConfigResource) GetConfig(ctxt context.Context) (*resourceconfig.Dogg3rzConfig, error) {

	err := validateConfig(configPath(ctxt))
	if err != nil {
		return nil, err
	}

	dgrzCfg := &resourceconfig.Dogg3rzConfig{}

	byteValue, err := ioutil.ReadFile(configPath(ctxt))
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
func configPath(ctxt context.Context) string {

	// gojsonschema reequires the schema 'file://' prepeended to path
	return path.Join(file.DotDirPath(ctxt), configFileName)

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
