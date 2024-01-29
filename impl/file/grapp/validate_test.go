package grapp

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

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

func TestValidate(t *testing.T) {

	var grappDir string

	if tmpDir, err := os.MkdirTemp("", "TestValidate"); err != nil {
		t.Fatal("os.MkdirTemp", err)
	} else {
		grappDir = tmpDir
	}
	//defer os.RemoveAll(grappDir)

	// SET NEW TMP DIR AS GRAPP DIR VIA CONTEXT VAR OVERRIDE (FROM PWD)
	ctxt := context.Background()
	ctxt = context.WithValue(ctxt, env.EnvDogg3rzGrapp, grappDir)

	// INIT TMPDIR AS GRAPP DIR FIRST
	if err := initGrappDir(ctxt, grappDir); err != nil {
		t.Fatal(err)
	}

	// get object dir
	var od string
	var err error
	od, err = file.GrapplicationObjectsDirPath(ctxt)
	if err != nil {
		t.Fatal("file.GrapplicationObjectsDirPath", err)
	}

	// stage good json
	if jsonLdFilePath, err := stageFile("good.jsonld", grappDir); err != nil {
		t.Fatal(err)
	} else {
		if err := validateGrappProjectFiles(ctxt, grappDir, od, os.Stdout); err != nil {
			//fmt.Println("failed here 111")
			t.Fatal(err)

		}
		os.Remove(jsonLdFilePath)

	}

	// stage malformed json file

	if jsonLdFilePath, err := stageFile("malformed-unclosed-quote.jsonld", grappDir); err != nil {
		t.Fatal(err)
	} else {

		if err := validateGrappProjectFiles(ctxt, grappDir, od, os.Stdout); err == nil {
			t.Fatal("expected error on malformed file", jsonLdFilePath, err)

		}
		os.Remove(jsonLdFilePath)

	}

	// stage no-rdf-stmts json file

	if jsonLdFilePath, err := stageFile("no-rdf-stmts.jsonld", grappDir); err != nil {
		t.Fatal(err)
	} else {

		if err := validateGrappProjectFiles(ctxt, grappDir, od, os.Stdout); err == nil {
			t.Fatal("expected error on jsonld file with no rdf statements produced after expansion", jsonLdFilePath, err)

		}
		os.Remove(jsonLdFilePath)

	}
	//t.FailNow()

}

func stageFile(filename string, dir string) (string, error) {

	src := filepath.Join("testfiles", "validate", filename)
	dst := filepath.Join(dir, filename)

	var src_fp, dst_fp *os.File
	var err error

	src_fp, err = os.Open(src)
	if err != nil {
		return "", err
	}
	defer src_fp.Close()

	dst_fp, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	defer dst_fp.Close()

	fmt.Println("copy", src, "to", dst, "...")
	_, err = io.Copy(dst_fp, src_fp)
	if err != nil {
		return "", err
	}

	return dst, err

}
