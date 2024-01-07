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
	var dgrzConf config.Dogg3rzConfig
	var defaultConfS string

	dgrzConf.User.ActivityPubUserHandle = "@test@datacequia.com"
	
	defaultConfS, err = config.GenerateDefault(dgrzConf)
	if err != nil {
		t.FailNow()
	}
	//fmt.Println(defaultConfS)

	defaultConfigPath := path.Join(os.TempDir(), "dogg3rz-default-config.json")
	err = ioutil.WriteFile(defaultConfigPath,
		[]byte(defaultConfS), os.FileMode(0777))
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
