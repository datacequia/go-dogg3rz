package grapp

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/piprate/json-gold/ld"
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

func (grapp *FileGrapplicationResource) Validate(ctxt context.Context) error {

	//fmt.Println("file vaalidate ", grappDir)
	grappDir, err := file.GrapplicationDirPath(ctxt)
	if err != nil {
		return err
	}

	return validate(ctxt, grappDir)
}

func validate(ctxt context.Context, grappDir string) error {

	projectFiles, err := listJsonLdFiles(grappDir)
	if err != nil {
		return err
	}

	if len(projectFiles) < 1 {
		return errors.NotFound.Newf("%s: no JSON-LD files found.", grappDir)
	}

	// process JSON-LD files aganst JSON-LD processor for well-formedness
	for _, jsonLdFile := range projectFiles {

		var err error
		var parsedJson map[string]interface{}

		if parsedJson, err = parseJSON(jsonLdFile); err != nil {
			return err
		}
		var expandedJsonLd interface{}

		if p, err = process(parsedJson); err != nil {

			return err
		}

		//fmt.Printf("expanded type %T", p)

		if len(p) < 1 {
			return errors.NotFound.Newf("%s: no RDF statement were found after processing", jsonLdFile)
		}

		//fmt.Println("after process")
		ld.PrintDocument(jsonLdFile, expandedJsonLd)

		/*
			termNsMap := make(map[string]string)
			if termNsMap, err = extractNamespacesFromContext(parsedJson); err != nil {
				return err
			}

			for k, v := range termNsMap {
				fmt.Printf("term = %s, ns = %s\n", k, v)

			}
		*/
	}

	return nil

}

func parseJSON(path string) (map[string]interface{}, error) {

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	jsonMap := make(map[string]interface{})

	parseStats := &jsonParseStats{realReader: fp}
	parseStats.newlineOffsets = make([]int64, 0)

	decoder := json.NewDecoder(parseStats)

	err = decoder.Decode(&jsonMap)
	if err != nil {
		//fmt.Println("InputOffset is ", decoder.InputOffset())
		//fmt.Printf("decoder.Decode returned type %T\n", err)
		if r, ok := err.(*json.SyntaxError); ok {
			var syntaxErrCol int64 = r.Offset
			var syntaxErrLine int64 = 1

			// compute column
			for i, nlOffset := range parseStats.newlineOffsets {

				if r.Offset <= nlOffset {
					if i > 0 {
						//fmt.Println("syntaxErrCol", i, r.Offset, parseStats.newlineOffsets[i-1])
						syntaxErrCol = r.Offset - parseStats.newlineOffsets[i-1]
					} else {
						syntaxErrCol = r.Offset
						//fmt.Println("syntaxErrCol(i<=0)", syntaxErrCol)
					}
					syntaxErrLine = int64(i + 1)
					break
				}

			}

			return nil, errors.Newf("%s:%d:%d: %s", path, syntaxErrLine, syntaxErrCol, err.Error())

		}
		//fmt.Println("other decode err", err)
		return nil, err
	}

	return jsonMap, nil

}

func extractNamespacesFromContext(jsonMap map[string]interface{}) (map[string]string, error) {

	m := make(map[string]string)

	return m, nil
}

func process(jsonMap map[string]interface{}) ([]interface{}, error) {

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"
	//var fileContent []byte
	var err error
	/*expanded*/

	flattened, err := proc.Flatten(jsonMap, nil, options)

	if err != nil {
		return nil, err
	}

	var flattenedList []interface{}

	if t, ok := flattened.([]interface{}); !ok {
		return nil, errors.UnexpectedType.Newf("expected type %T returned from proc.Flatten(), got %T",
			flattenedList, flattened)
	} else {
		flattenedList = t
	}

	//ld.PrintDocument("process", expanded)

	return flattenedList, nil

}

func listJsonLdFiles(grappDir string) ([]string, error) {

	// LIST ALL JSONLD FILES
	files, err := os.ReadDir(grappDir)

	if err != nil {
		return nil, err
	}

	var jsonLDFiles []string

	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(strings.ToLower(file.Name()), ".jsonld") {
			jsonLDFiles = append(jsonLDFiles, filepath.Join(grappDir, file.Name()))
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
