// Common dependent test functions that are used by other test
// in this module
package grapp

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
	"github.com/datacequia/go-dogg3rz/ipfs"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

const (
	testGrappName = "grappTest"
)

func TestInitTestNode(t *testing.T) {

	ctxt, cancelFunc, err := initTestNode("myprefix")
	if err != nil {
		t.Error(err)
	}
	defer cancelFunc()

	homeDir, ok := ctxt.Value(env.EnvDogg3rzHome).(string)
	if !ok {
		t.Errorf("value from context key '%s' is not a string", env.EnvDogg3rzHome)
		return
	}

	if !file.DirExists(homeDir) {
		t.Errorf("node init'd but dogg3rz home dir does not exist: %s", homeDir)
	}

	defer os.RemoveAll(homeDir)

	grappName, ok2 := ctxt.Value(env.EnvDogg3rzGrapp).(string)
	if !ok2 {
		t.Errorf("value from context key '%s' is not a string", env.EnvDogg3rzGrapp)
		return
	}

	if grappName != testGrappName {
		t.Errorf("expected testGrappName value '%s' in context, found '%s'",
			testGrappName, grappName)

	}

}

// Initialize a new dogg3rz node in a temp dir
// Returns non-cancellable context with value for DOGG3RZ_HOME and DOGG3RZ_GRAPP assigned
// or error initialized if failed
func initTestNode(dogg3rzHomePrefix string) (context.Context, context.CancelFunc, error) {

	if len(dogg3rzHomePrefix) < 1 {
		return nil, nil, errors.New("dogg3rzHomePrefix len must be > 0")
	}

	dogg3rzHome := filepath.Join(os.TempDir(),
		fmt.Sprintf("%s_%d", dogg3rzHomePrefix, time.Now().UnixNano()))

	ctxt, cancelFunc := context.WithCancel(context.Background())

	ctxt = context.WithValue(ctxt, env.EnvDogg3rzHome, dogg3rzHome)
	ctxt = context.WithValue(ctxt, env.EnvDogg3rzGrapp, testGrappName)

	// SPAWN AN EPHEMERAL IPFS NODE TO INTERACT W/ DOGG3RZ TEST NODE
	ipfs.SpawnEphemeral(ctxt)

	// os.Setenv("DOGG3RZ_HOME", dogg3rzHome)

	fileNodeResource := &filenode.FileNodeResource{}

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.Email = "test@dogg3rz.com"

	if err := fileNodeResource.InitNode(ctxt, dgrzConf); err != nil {
		//t.Error(err)
		return ctxt, cancelFunc, err
	}

	//t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	fileGrapplicationResource := FileGrapplicationResource{}

	if err := fileGrapplicationResource.InitGrapp(ctxt, testGrappName); err != nil {
		return ctxt, cancelFunc, err
	}

	return ctxt, cancelFunc, nil

}
