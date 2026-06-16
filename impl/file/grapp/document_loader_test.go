package grapp

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/piprate/json-gold/ld"
)

func TestDocumentLoader(t *testing.T) {

	var grappDir string
	var err error

	// Create temp grapp project dir
	grappDir, err = os.MkdirTemp(os.TempDir(), "TestGrapplicationCachedRemoteDocument*")
	if err != nil {
		t.Fatal("os.MkdirTemp", err)
	}

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

	cwd, _ := os.Getwd()

	docPath := path.Join(cwd, "testfiles", "schemaorg-current-https.jsonld")
	doc, err = dl.LoadDocument(docPath)
	if err != nil {
		t.Fatal("dl.LoadDocument", err)

	}
	fmt.Println("DocumentURL", doc.DocumentURL)
	fmt.Println("ContextURL", doc.ContextURL)

	if doc.Document == nil {
		t.Fatal("doc.Document is nil", doc.Document)
	}
	//fmt.Println(gd)

	//t.FailNow()

}

func TestToProjectFileSubPath(t *testing.T) {

	// TODO: TEST WHEN ON WINDOWS

	type input struct {
		projectPath   string
		filePath      string
		expectedError error
		returnValue   string
	}

	//grappDir = ""
	//filePath = ""
	//var err error

	test1 := input{}
	test1.projectPath = "/some/random/non-existent/project/path-2423rt34asdfadf234rq23fad"
	test1.filePath = "blah"
	test1.expectedError = &fs.PathError{}
	test1.returnValue = ""

	// MAKE PROJECT DIR PATH A FILE (NOT A DIR)
	test2 := input{}

	if f, err := os.CreateTemp(os.TempDir(), "test2-*"); err != nil {

		t.Fatal("test2: os.CreateTemp", err)

	} else {
		test2.projectPath = f.Name() // assign project path a file path
		defer os.Remove(f.Name())
		defer f.Close()

	}

	test2.filePath = "my.jsonld"
	test2.expectedError = errors.InvalidValue.New("")
	test2.returnValue = ""

	// test3
	// TEST 3
	test3 := input{}
	if f, err := os.MkdirTemp(os.TempDir(), "test3-*"); err != nil {

		t.Fatal("test3: os.CreateTemp", err)

	} else {
		test3.projectPath = f
		defer os.Remove(f)
		//defer f.Close()

	}
	test3.filePath = filepath.Join(test3.projectPath, "my.jsonld")
	test3.expectedError = nil
	test3.returnValue = "my.jsonld"
	if err := file.Touch(test3.filePath); err != nil {
		t.Fatal("test3: file.Touch", err)
	}

	// test 4
	// project file path is relative
	test4 := input{}
	if f, err := os.MkdirTemp(os.TempDir(), "test4-*"); err != nil {

		t.Fatal("test4: os.CreateTemp", err)

	} else {
		test4.projectPath = f
		defer os.Remove(f)
		//defer f.Close()

	}
	test4.filePath = "my.jsonld"
	test4.expectedError = nil
	test4.returnValue = "my.jsonld"

	if err := file.Touch(filepath.Join(test4.projectPath, test4.filePath)); err != nil {
		t.Fatal("test4: os.Mkdir", err)
	}

	// test4_1
	// test 4
	// project file path is relative
	test4_1 := input{}
	if f, err := os.MkdirTemp(os.TempDir(), "test4_1-*"); err != nil {

		t.Fatal("test4_1: os.CreateTemp", err)

	} else {
		test4_1.projectPath = f
		defer os.Remove(f)
		//defer f.Close()

	}
	test4_1.filePath = "sub//path//my.jsonld" // create extra empy path separators. should be cleaned
	test4_1.expectedError = nil
	test4_1.returnValue = "sub/path/my.jsonld"

	if err := os.MkdirAll(filepath.Join(test4_1.projectPath, "sub", "path"), 0755); err != nil {
		t.Fatal("test4_1: os.MkdirAll", err)

	}

	if err := file.Touch(filepath.Join(test4_1.projectPath, test4_1.filePath)); err != nil {
		t.Fatal("test4_1: os.Mkdir", err)
	}

	// test 5 -
	// project file is relative but does not exist
	test5 := input{}
	if f, err := os.MkdirTemp(os.TempDir(), "test5-*"); err != nil {

		t.Fatal("test5: os.CreateTemp", err)

	} else {
		test5.projectPath = f
		defer os.Remove(f)
		//defer f.Close()

	}
	test5.filePath = "test5.jsonld"
	test5.expectedError = &fs.PathError{}
	test5.returnValue = "test5.jsonld"

	// test6
	// project file path is relative but is a directory
	test6 := input{}
	if f, err := os.MkdirTemp(os.TempDir(), "test6-*"); err != nil {

		t.Fatal("test6: os.CreateTemp", err)

	} else {
		test6.projectPath = f
		defer os.Remove(f)
	}
	test6.filePath = "test6.jsonld"
	test6.expectedError = errors.InvalidValue.New("")
	test6.returnValue = "test6.jsonld"
	if err := os.Mkdir(filepath.Join(test6.projectPath, test6.filePath), 0755); err != nil {
		t.Fatal("test6: os.Mkdir", err)
	}

	inputList := []input{
		test1, //
		test2, //
		test3,
		test4,
		test4_1,
		test5,
		test6,
	}

	sameErrorType := func(a, b error) (string, bool) {

		if reflect.TypeOf(a) == reflect.TypeOf(b) {

			if t, ok := a.(errors.BadDogg3rz); ok {
				if t2, ok2 := b.(errors.BadDogg3rz); ok2 {

					if t.ErrorType != t2.ErrorType {
						return fmt.Sprintf("Expected BadDoggrz error Type %s, found %s: %s", a, b, b.Error()), false
					}
				}
			}
			return "", true
		}

		return fmt.Sprintf("Expected error type %T, found %T", a, b), false

	}

	for _, i := range inputList {
		var pp string
		var err error
		//i.preCond()
		if pp, err = toProjectFileSubPath(i.projectPath, i.filePath); err != nil {
			if msg, ok := sameErrorType(i.expectedError, err); !ok {
				t.Fatalf("%s: { projectPath=%s, filePath=%s}", msg, i.projectPath, i.filePath)
			}
		} else {
			if i.returnValue != pp {
				t.Fatalf("expected return value '%s' when returning error type %T", i.returnValue, i.expectedError)
			}
		}
	}
	t.FailNow()
}

func TestUrlParse(t *testing.T) {

	if parsedURL, err := url.Parse("/how/now/test"); err == nil {
		if parsedURL.Host != "" {
			t.Fatalf("expected empty Host, found '%s'", parsedURL.Host)

		}
		if parsedURL.Scheme != "" {
			t.Fatalf("expected empty Schema, found '%s'", parsedURL.Scheme)

		}
	}

	if parsedURL, err := url.Parse("\\how\\now\\test"); err == nil {
		if parsedURL.Host != "" {
			t.Fatalf("expected empty Host, found '%s'", parsedURL.Host)

		}
		if parsedURL.Scheme != "" {
			t.Fatalf("expected empty Schema, found '%s'", parsedURL.Scheme)

		}
		if parsedURL.Path != "/how/now/test" {
			t.Fatalf("expected Path '%s', found '%s'", "/how/now/test", parsedURL.Path)

		}
	}
}

func TestNormalizePrefixIRI(t *testing.T) {

	var err error
	var normUrl string
	var givenUrl string
	var expectedUrl string

	// test lcasing and duplicate slashes and preserve fragment
	givenUrl = "http://Example.COM//#"
	expectedUrl = "http://example.com/#"
	if normUrl, err = normalizePrefixIRI(givenUrl); err != nil {
		t.FailNow()
	}
	if normUrl != expectedUrl {
		t.Errorf("expected %s, found %s", expectedUrl, normUrl)

	}

	// test rm'ing dot segments
	givenUrl = "http://Example.COM/./foo/"
	expectedUrl = "http://example.com/foo/"
	if normUrl, err = normalizePrefixIRI(givenUrl); err != nil {
		t.FailNow()
	}
	if normUrl != expectedUrl {
		t.Errorf("expected %s, found %s", expectedUrl, normUrl)

	}

	givenUrl = "HTTp://Example.COM//foo/"
	expectedUrl = "http://example.com/foo/"
	if normUrl, err = normalizePrefixIRI(givenUrl); err != nil {
		t.FailNow()
	}
	if normUrl != expectedUrl {
		t.Errorf("expected %s, found %s", expectedUrl, normUrl)

	}

}
