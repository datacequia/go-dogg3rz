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
	"os"
	"path"
	"testing"

	"github.com/datacequia/go-dogg3rz/resource/config"
	//"github.com/datacequia/go-dogg3rz/impl/config"
)

func TestValidateConfig(t *testing.T) {

	err := validateConfig("/some/bad/path/safasdfsd")
	if err == nil {
		t.Error("config validated with a bad path")
	}

	testFilesBad := []string{"file://testfiles/malformed-config.json"}

	for _, tf := range testFilesBad {
		err = validateConfig(tf)
		if err == nil {
			t.Errorf("config validation succeeded, expected fail: %s", tf)
		}
	}

	testFilesGood := []string{"testfiles/good-config.json"}

	for _, tf := range testFilesGood {
		err = validateConfig(tf)
		if err != nil {
			t.Errorf("config validation failed, expected success: %s: %s", tf, err)
		}
	}

	// TEST THE DEFAULT CONFIG
	defaultConfigPath := path.Join(os.TempDir(), "dogg3rz-default-config.json")
	err = ioutil.WriteFile(defaultConfigPath,
		[]byte(config.CONFIG_JSON_DEFAULT), os.FileMode(0777))
	if err != nil {
		// FAILED TO WRITE DEFAULT
		t.Fail()
	}
	defer os.Remove(defaultConfigPath)

	// validateConfig arg requires URL protocol schema (i.e. file://)
	err = validateConfig(defaultConfigPath)
	if err != nil {
		t.Errorf("default config (config.CONFIG_JSON_DEFAULT) is malformed or config.CONFIG_JSON_SCHEMA is malformed: %s", err)

	}

	// NOW UNMARSHAL GOOD CONFIG AND MAKE SURE ALL ATTR VALUES GEET ASSIGNED
	byteValue, err := ioutil.ReadFile("testfiles/good-config.json")
	if err != nil {
		t.Fail()
	}

	dgrzCfg := &config.Dogg3rzConfig{}

	err = json.Unmarshal(byteValue, &dgrzCfg)
	if err != nil {
		t.Errorf("failed to unmarshal config: %s", err)
	}

	// CHCK THAT GOOD CONFIG HAS APIENDPINT ASSIGNED

	if dgrzCfg.IPFS.ApiEndpoint != "http://localhost:5001/" {
		t.Errorf("Expected %s, found %s", "http://localhost:5001", dgrzCfg.IPFS.ApiEndpoint)
	}

}
