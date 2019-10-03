package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/datacequia/go-dogg3rz/impl/config"
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

	testFilesGood := []string{"file://testfiles/good-config.json"}

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

	err = validateConfig("file://" + defaultConfigPath)
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
