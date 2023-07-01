package ipfs

import (
	"io/fs"
	"os"
	"testing"
)

func TestNewIpfsContainer(t *testing.T) {

	o := newIpfsContainer()

	tmpPath := o.tempDirLocation()

	if _, err := os.Stat(tmpPath); err != nil {
		t.Errorf("not a valid temp dir: %s", tmpPath)

	}

	if _, err := o.createTempDir(tmpPath, "test_*"); err != nil {
		t.Errorf("failed to create temp dir in %s", tmpPath)
	}

}

func TestGetStdoutStderrRedirectFiles(t *testing.T) {

	o := newIpfsContainer()

	if sout, serr, err := o.getStdoutStderrRedirectFiles("test_*"); err != nil {
		if sout == nil {
			t.Errorf("returned uninitialized stdout")
		}
		if serr == nil {
			t.Errorf("returned uninitialized stderr")
		}

		t.Errorf("failed to get redirect files: %s", err)

	}

	o.tempDirLocation = func() string { return "bad dir " }

	if _, _, err := o.getStdoutStderrRedirectFiles("test_*"); err == nil {
		t.Errorf("expect failed on bad dir location")
	} else {
		if _, ok := err.(*fs.PathError); !ok {
			t.Errorf("expected fs.PathError, found %T", err)
		}
	}

}

func TestPull(t *testing.T) {

	o := newIpfsContainer()

	o.tempDirLocation = func() string { return "bad dir   xx" }

	if err := o.pull(); err == nil {
		t.Errorf("Expected fail with bad dir, returned success")

	}

	//o = newIpfsContainer()

}
