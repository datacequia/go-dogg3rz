// Common dependent test functions that are used by other test
// in this module
package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/impl/file"
	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

const (
	testRepoName = "repoTest"
)

func TestInitTestNode(t *testing.T) {

	ctxt, err := initTestNode("myprefix")
	if err != nil {
		t.Error(err)
	}

	homeDir, ok := ctxt.Value(env.EnvDogg3rzHome).(string)
	if !ok {
		t.Errorf("value from context key '%s' is not a string", env.EnvDogg3rzHome)
		return
	}

	if !file.DirExists(homeDir) {
		t.Errorf("node init'd but dogg3rz home dir does not exist: %s", homeDir)
	}

	defer os.RemoveAll(homeDir)

	repoName, ok2 := ctxt.Value(env.EnvDogg3rzRepo).(string)
	if !ok2 {
		t.Errorf("value from context key '%s' is not a string", env.EnvDogg3rzRepo)
		return
	}

	if repoName != testRepoName {
		t.Errorf("expected testRepoName value '%s' in context, found '%s'",
			testRepoName, repoName)

	}

}

// Initialize a new dogg3rz node in a temp dir
// Returns non-cancellable context with value for DOGG3RZ_HOME and DOGG3RZ_REPO assigned
// or error initialized if failed
func initTestNode(dogg3rzHomePrefix string) (context.Context, error) {

	if len(dogg3rzHomePrefix) < 1 {
		return nil, errors.New("dogg3rzHomePrefix len must be > 0")
	}

	//testRepoName := "test"

	dogg3rzHome := filepath.Join(os.TempDir(),
		fmt.Sprintf("%s_%d", dogg3rzHomePrefix, time.Now().UnixNano()))

	ctxt := context.Background()

	ctxt = context.WithValue(ctxt, env.EnvDogg3rzHome, dogg3rzHome)
	ctxt = context.WithValue(ctxt, env.EnvDogg3rzRepo, testRepoName)

	// os.Setenv("DOGG3RZ_HOME", dogg3rzHome)

	fileNodeResource := &filenode.FileNodeResource{}

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.Email = "test@dogg3rz.com"

	if err := fileNodeResource.InitNode(ctxt, dgrzConf); err != nil {
		//t.Error(err)
		return nil, err
	}

	//t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	fileRepositoryResource := FileRepositoryResource{}

	if err := fileRepositoryResource.InitRepo(ctxt, testRepoName); err != nil {
		return nil, err
	}

	return ctxt, nil

}
