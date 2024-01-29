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

/*
func TestGrapplicationCachedRemoteDocument(t *testing.T) {

	var cacheDir string
	var err error

	cacheDir, err = os.MkdirTemp(os.TempDir(), "TestGrapplicationCachedRemoteDocument*")
	if err != nil {
		t.Fatal("os.MkdirTemp", err)
	}
	defer os.RemoveAll(cacheDir)
	//fmt.Println("cacheDir", cacheDir)

	var doc *os.File

	doc, err = os.Open("testfiles/schemaorg-current-https.jsonld")
	if err != nil {
		t.Fatal("os.Open", err)
	}
	var docInfo fs.FileInfo

	docInfo, err = doc.Stat()
	if err != nil {
		t.Fatal("doc.Stat", err)
	}

	dl, err := NewDocumentLoader(nil, )
	rd, err := NewGrapplicationCachedRemoteDocument(doc, cacheDir, crypto.SHA1)
	if err != nil {
		t.Fatal("NewGrapplicationCachedRemoteDocument", err)

	}
	defer rd.Close()

	// read doc
	var n, totalN int
	var buf []byte = make([]byte, 4096)

	for n, err = rd.Read(buf); err == nil; n, err = rd.Read(buf) {
		//fmt.Println("rd,Read", n, err)
		totalN += n
	}

	if err != nil && err != io.EOF {

		t.Fatal("rd.Read", err)
	}

	if err == io.EOF {
		fmt.Println("reached eof:", doc.Name())
	}
	if len(rd.cachedDocPath) < 1 {
		t.Fatal("rd.cachedDocPath not set", rd.cachedDocPath)
	}
	if !file.FileExists(rd.cachedDocPath) {
		t.Fatal("cached Doc not created at", rd.cachedDocPath)
	}

	var docInfo2 fs.FileInfo
	docInfo2, err = os.Stat(rd.cachedDocPath)
	if err != nil {
		t.Fatal("os.Stat(rd.cachedDocPath)", err)
	}

	if docInfo.Size() != docInfo2.Size() {
		t.Fatal("cached doc file size does not match original file size: original size = ", docInfo.Size(),
			"cache size = ", docInfo2.Size())

	}

	if docInfo2.Size() != int64(totalN) {
		t.Fatal("total bytes returned from Read() don't match cache file size: cache size = ", docInfo2.Size(),
			"totalN = ", totalN)

	}

	fmt.Println("n", n)

	//t.FailNow()

}
*/

func TestDocumentLoader(t *testing.T) {

	var grappDir string
	//var objectsDir string
	var err error

	grappDir, err = os.MkdirTemp(os.TempDir(), "TestGrapplicationCachedRemoteDocument*")
	if err != nil {
		t.Fatal("os.MkdirTemp", err)
	}
	fmt.Println("grappDir", grappDir)
	//defer os.RemoveAll(grappDir )
	ctxt := context.Background()
	ctxt = context.WithValue(ctxt, env.EnvDogg3rzGrapp, grappDir)

	// INIT TMPDIR AS GRAPP DIR FIRST
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
