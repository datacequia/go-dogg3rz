package grapp

import (
	"context"
	"os"
	"testing"
)

func TestNextIPFSAPIPort(t *testing.T) {

	var baseDir string

	if tmpDir, err := os.MkdirTemp("", "TestNextIPFSAPIPort"); err != nil {
		t.Fatal("os.MkdirTemp", err)
	} else {
		baseDir = tmpDir
	}
	defer os.RemoveAll(baseDir)

	ctxt := context.Background()

	const maxReservedPort = 1024

	if _, err := nextIPFSAPIPort(ctxt, maxReservedPort, baseDir); err == nil {

		t.Fatal("didn't fail with non-user level port", maxReservedPort)
	}

	badDir := "/this/is a / bad dir / path"
	if _, err := nextIPFSAPIPort(ctxt, 1025, badDir); err == nil {
		t.Fatal("didn't failed with bad dir", badDir)
	}

	if port, err := nextIPFSAPIPort(ctxt, 1025, baseDir); err != nil {
		t.Fatal("should've succeeded with good base port and dir: got ", err)

		if port != 1025 {
			t.Fatal("expected basePort on new tmpDir, got ", port)
		}
	}

	for i := 1; i <= 10; i++ {
		if port, err := nextIPFSAPIPort(ctxt, 1025, baseDir); err != nil {
			t.Fatal("should've succeeded with good base port and dir: got ", err)

			if port != 1025+i {
				t.Fatalf("expected port %d on existing tmpDir, got %d", 1025+i, port)
			}
		}
	}

	//t.Fail()
}
