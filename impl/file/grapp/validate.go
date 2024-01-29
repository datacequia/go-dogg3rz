package grapp

/*
 * Copyright (c) 2019-2024 Datacequia LLC. All rights reserved.
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
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

type jsonParseStats struct {
	realReader     io.Reader
	bytesRead      int64
	newlineOffsets []int64
}

type triple struct {
	subject     string
	predicate   string
	object      string
	object_type interface{} // IRI for XSD TYPE OR "@"
}

func (grapp *FileGrapplicationResource) Validate(ctxt context.Context, vw io.Writer) error {

	var objectsDir string
	var grappDir string
	var err error

	objectsDir, err = file.GrapplicationObjectsDirPath(ctxt)
	if err != nil {
		return err
	}

	grappDir, err = file.GrapplicationDirPath(ctxt)
	if err != nil {
		return err
	}

	err = validateGrappProjectFiles(ctxt, grappDir, objectsDir, vw)
	if err != nil {
		verbose(vw, "Validation completed with errors: ", err)
	} else {
		verbose(vw, "Validation completed successfully!")
	}
	return err

}

func validateGrappProjectFiles(ctxt context.Context, grappDir string, objectsDir string, vw io.Writer) error {

	verbose(vw, "Listing .jsonld files in project directory at %s...", grappDir)
	projectFiles, err := listJsonLdFiles(grappDir, vw)
	if err != nil {
		return err
	}

	if len(projectFiles) < 1 {
		return errors.NotFound.Newf("%s: no JSON-LD files found.", grappDir)
	}

	// process JSON-LD files against JSON-LD processor for well-formedness
	for _, jsonLdFile := range projectFiles {

		loader := NewDocumentLoader(nil, grappDir, objectsDir)

		if _, err := loader.LoadDocument(jsonLdFile); err != nil {
			return err
		}

	}

	return nil

}

func extractNamespacesFromContext(jsonMap map[string]interface{}) (map[string]string, error) {

	m := make(map[string]string)

	return m, nil
}

/*
func process(doc interface{}) ([]interface{}, error) {

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	grappDocumentLoader := NewDocumentLoader(nil, "")

	//options.Format = "application/n-quads"
	options.DocumentLoader = grappDocumentLoader

	var err error

	expanded, err := proc.Expand(doc, options)

	if err != nil {
		return nil, err
	}

	var flattened interface{}

	flattened, err = proc.Flatten(expanded, nil, options)

	var flattenedList []interface{}

	if t, ok := flattened.([]interface{}); !ok {
		return nil, errors.UnexpectedType.Newf("expected type %T returned from proc.Normalize(), got %T",
			flattenedList, flattened)
	} else {
		flattenedList = t
	}

	ld.PrintDocument("Normalize", flattened)

	return flattenedList, nil

}
*/

func listJsonLdFiles(grappDir string, vw io.Writer) ([]string, error) {

	// LIST ALL JSONLD FILES
	files, err := os.ReadDir(grappDir)

	if err != nil {
		return nil, err
	}

	var jsonLDFiles []string

	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(strings.ToLower(file.Name()), ".jsonld") {
			newFile := filepath.Join(grappDir, file.Name())
			jsonLDFiles = append(jsonLDFiles, newFile)

			verbose(vw, newFile)
		}
	}

	return jsonLDFiles, nil

}

func (s *jsonParseStats) Read(p []byte) (int, error) {

	bytesRead, err := s.realReader.Read(p)
	//fmt.Println("bytesRead", bytesRead, err)
	if err == nil {

		for i, b := range p[:bytesRead] {
			if b == '\n' {
				s.newlineOffsets = append(s.newlineOffsets, s.bytesRead+int64(i+1))
			}
		}
		s.bytesRead += int64(bytesRead)
	}

	return bytesRead, err

}

func verbose(w io.Writer, msg string, args ...interface{}) (int, error) {
	if w == nil {
		return 0, nil
	}

	return fmt.Fprintf(w, msg+"\n", args...)

}
