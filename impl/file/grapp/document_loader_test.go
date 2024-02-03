package grapp

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/piprate/json-gold/ld"
)

func TestDocumentLoader(t *testing.T) {

	var grappDir string
	//var objectsDir string
	var err error

	// Create temp grapp project dir
	grappDir, err = os.MkdirTemp(os.TempDir(), "TestGrapplicationCachedRemoteDocument*")
	if err != nil {
		t.Fatal("os.MkdirTemp", err)
	}
	fmt.Println("grappDir", grappDir)
	defer os.RemoveAll(grappDir)
	// set context/env var to point to tmp project dir
	ctxt := context.Background()
	ctxt = context.WithValue(ctxt, env.EnvDogg3rzGrapp, grappDir)

	// INIT TMPDIR AS GRAPP DIR
	if err := initGrappDir(ctxt, grappDir); err != nil {
		t.Fatal(err)
	}

	if val, ok := ctxt.Value(env.EnvDogg3rzGrapp).(string); !ok {
		t.Fatal("can't retrieve env.EnvDogg3rzGrapp after setting it")
	} else {
		if val != grappDir {
			t.Fatalf("context returned '%s', expected '%s'", val, grappDir)
		}
	}
	var gd string
	var od string

	gd, err = file.GrapplicationDirPath(ctxt)
	if err != nil {
		t.Fatal("file.GrapplicationsDirPath", err)
	}

	od, err = file.GrapplicationObjectsDirPath(ctxt)
	if err != nil {
		t.Fatal("file.GrapplicationObjectsDirPath", err)
	}

	dl := NewDocumentLoader(nil, gd, od)

	var doc *ld.RemoteDocument

	doc, err = dl.LoadDocument("testfiles/schemaorg-current-https.jsonld")
	if err != nil {
		t.Fatal("dl.LoadDocument", err)

	}
	fmt.Println("DocumentURL", doc.DocumentURL)
	fmt.Println("ContextURL", doc.ContextURL)

	if doc.Document == nil {
		t.Fatal("doc.Document is nil", doc.Document)
	}

	//t.FailNow()

}
