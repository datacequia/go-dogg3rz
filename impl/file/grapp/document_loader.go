package grapp

import (
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/fxamacker/cbor/v2"
	"github.com/piprate/json-gold/ld"
)

const (
	// An HTTP Accept header that prefers JSONLD.
	acceptHeader = "application/ld+json, application/json;q=0.9, application/javascript;q=0.5, text/javascript;q=0.5, text/plain;q=0.2, */*;q=0.1"

	// JSON-LD link header rel
	linkHeaderRel = "http://www.w3.org/ns/json-ld#context"
)

type DocumentLoader struct {
	httpClient      *http.Client          // optional: http client to use to load documents
	grappDir        string                // base project dir
	objectsDir      string                // where to place object files
	objectFileIndex map[string]objectFile // index of processed documents cached as binary object files (cbor format)
}

type objectFile struct {
	givenLocation        string // iri or (relative) location to project of document (.jsonld) file that was processed in project
	standardizedLocation string // location standardized , for urls its the value returns from [X], for local paths it's relative path to project dir
	resolvedLocation     string // location resolved to the actual path where the JSON-LD can be dereferenced
	documentLocationHash string // hash of source path
	documentContentHash  string // hash of sourcePath content

	processedDoc interface{} // JSON-LD processor output of document
}

func NewDocumentLoader(httpClient *http.Client, grappDir string, objectsDir string) *DocumentLoader {
	rval := &DocumentLoader{httpClient: httpClient, grappDir: grappDir, objectsDir: objectsDir}

	if rval.httpClient == nil {
		rval.httpClient = http.DefaultClient
	}

	rval.objectFileIndex = make(map[string]objectFile)

	return rval
}

// Loads JSON-LD documents from local or http paths
// Implements github.com/piprate/ld/DocumentLoader interface
func (dl *DocumentLoader) LoadDocument(u string) (*ld.RemoteDocument, error) {

	var givenLocation string
	var standardizedLocation string
	var resolvedLocation string

	givenLocation = u

	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, fmt.Sprintf("error parsing URL: %s", u))
	}

	var documentBody io.ReadCloser
	var finalURL, contextURL string

	protocol := parsedURL.Scheme

	if protocol != "http" && protocol != "https" {
		// Can't use the HTTP client for those!

		var file *os.File
		var absolutePath string
		var relativePath string
		var absolutePathGrappDir string

		file, err = os.Open(u)
		if err != nil {
			return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
		}
		defer file.Close()

		// GET CANONICAL PATH
		absolutePath, err = filepath.Abs(u)
		if err != nil {
			return nil, err
		}
		absolutePathGrappDir, err = filepath.Abs(dl.grappDir)
		if err != nil {
			return nil, err
		}

		relativePath, err = filepath.Rel(absolutePathGrappDir, absolutePath)
		if err != nil {
			return nil, err
		}
		finalURL = relativePath

		standardizedLocation = relativePath
		resolvedLocation = relativePath

		documentBody = file
	} else {
		resolvedLocation = resolveIRI(givenLocation)

		req, err := http.NewRequest("GET", resolvedLocation, nil)
		if err != nil {
			return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
		}
		// We prefer application/ld+json, but fallback to application/json
		// or whatever is available
		req.Header.Add("Accept", acceptHeader)

		res, err := dl.httpClient.Do(req)
		if err != nil {
			return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed,
				fmt.Sprintf("Bad response status code: %d", res.StatusCode))
		}

		finalURL = res.Request.URL.String()
		standardizedLocation = finalURL

		//	fmt.Println("finalURL", finalURL, "resolveIRI", resolveIRI(u)) // deleteme

		//fmt.Println("finalURL", finalURL)
		contentType := res.Header.Get("Content-Type")
		linkHeader := res.Header.Get("Link")

		if len(linkHeader) > 0 && contentType != "application/ld+json" {
			header := ld.ParseLinkHeader(linkHeader)[linkHeaderRel]
			if len(header) > 1 {
				return nil, ld.NewJsonLdError(ld.MultipleContextLinkHeaders, nil)
			} else if len(header) == 1 {
				contextURL = header[0]["target"]
			}
		}

		documentBody = res.Body

	}

	// read whole document body into memory
	var buf []byte
	if buf, err = io.ReadAll(documentBody); err != nil {

		return nil, err
	}

	parsedJSON, _, err := dl.createObjectFile(givenLocation,
		standardizedLocation,
		resolvedLocation,
		buf)
	if err != nil {
		return nil, err
	}

	return &ld.RemoteDocument{DocumentURL: finalURL, Document: parsedJSON, ContextURL: contextURL}, nil

}

func (dl *DocumentLoader) createObjectFile(givenLocation string, standardizedLocation string, resolvedLocation string, data []byte) (interface{}, string, error) {

	var obj objectFile

	obj.givenLocation = givenLocation
	obj.resolvedLocation = resolvedLocation
	obj.standardizedLocation = standardizedLocation

	// compute hash on standardized path location
	iriHash := crypto.SHA1.New()
	_, err := io.Copy(iriHash, strings.NewReader(obj.standardizedLocation))
	if err != nil {
		return nil, "", err
	}
	obj.documentLocationHash = fmt.Sprintf("%x", iriHash.Sum(nil))

	// compute hash on document
	docHash := crypto.SHA1.New()
	_, err = io.Copy(docHash, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	obj.documentContentHash = fmt.Sprintf("%x", docHash.Sum(nil))

	// CREATE STAGING  FILE
	var tmp *os.File
	if tmp, err = os.CreateTemp(dl.objectsDir, "createObjectFile-*"); err != nil {
		return nil, "", err
	}
	defer tmp.Close()

	// CLEANUP TMP FILE IF IT STILL EXISTS
	cleanup := func() {
		if file.FileExists(tmp.Name()) {
			os.Remove(tmp.Name())
		}
	}
	defer cleanup()

	// PARSE DOC CONTENTS
	var jsonTree map[string]interface{}
	jsonTree, err = parseJSON(bytes.NewBuffer(data), obj.resolvedLocation)
	if err != nil {
		return nil, "", err
	}
	// RUN JSON-LD PROCESSOR WITH PARSED JSON INPUT
	var expandedDoc []interface{}
	var flattenedDoc interface{}
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.DocumentLoader = NewDocumentLoader(dl.httpClient, dl.grappDir, dl.objectsDir)

	// EXPAND DOC (I.E. EXPAND JSON-LD TERMS TO FULL IRIs)

	expandedDoc, err = proc.Expand(jsonTree, options)
	if err != nil {
		//fmt.Println("expand failed", err, iri, jsonTree)
		return nil, "", err
	}
	if len(expandedDoc) < 1 {

		return nil, "", errors.NotFound.New("No RDF statements found after JSON-LD doc expansion")

	}

	// FLATTEN THE TREE TO AN ARRAY OF N-QUADS
	flattenedDoc, err = proc.Flatten(expandedDoc, nil, options)
	if err != nil {
		return nil, "", err
	}
	obj.processedDoc = flattenedDoc

	// ENCODE OBJECT TO CBOR FORMAT
	var cborObject []byte

	cborObject, err = cbor.Marshal(obj.toMap())
	if err != nil {
		return nil, "", err
	}

	// WRITE CBOR OBJECT DATA TO STAGE FILE
	_, err = io.Copy(tmp, bytes.NewReader(cborObject))
	if err != nil {
		return nil, "", err
	}
	// CONSTRUCT OBJECT FILE  PATH NAME
	objectFilePath := path.Join(dl.objectsDir, obj.documentLocationHash)

	// MOVE STAGED FILE TO OBJECT FILE  PATH
	if err := os.Rename(tmp.Name(), objectFilePath); err != nil {
		os.Remove(tmp.Name())
		//fmt.Println("rename err", err)
		return nil, "", err
	}
	// ADD OBJECT TO DOCUMENT INDEX
	dl.objectFileIndex[obj.standardizedLocation] = obj

	return jsonTree, objectFilePath, nil

}

func parseJSON(r io.Reader, src string) (map[string]interface{}, error) {

	jsonMap := make(map[string]interface{})

	parseStats := &jsonParseStats{realReader: r}
	parseStats.newlineOffsets = make([]int64, 0)

	decoder := json.NewDecoder(parseStats)

	err := decoder.Decode(&jsonMap)
	if err != nil {
		if r, ok := err.(*json.SyntaxError); ok {
			var syntaxErrCol int64 = r.Offset
			var syntaxErrLine int64 = 1

			// compute column
			for i, nlOffset := range parseStats.newlineOffsets {

				if r.Offset <= nlOffset {
					if i > 0 {
						syntaxErrCol = r.Offset - parseStats.newlineOffsets[i-1]
					} else {
						syntaxErrCol = r.Offset
					}
					syntaxErrLine = int64(i + 1)
					break
				}

			}

			return nil, errors.Newf("%s:%d:%d: %s", src, syntaxErrLine, syntaxErrCol, err.Error())

		}
		return nil, err
	}

	return jsonMap, nil

}

// resolves Prefix IRIs in an alternate location
// where prefix IRI does not make ontology available in their designated namespace
func resolveIRI(iri string) string {

	switch strings.ToLower(iri) { // urls are CASE SENSIVE. NORMALIZE TO LOWER
	case "https://schema.org/":
		return "https://schema.org/version/latest/schemaorg-current-https.jsonld"
	case "http://schema.org/":
		return "http://schema.org/version/latest/schemaorg-current-https.jsonld"
	default:
		return iri
	}
}

func (obj objectFile) toMap() map[string]interface{} {

	m := make(map[string]interface{})

	m["documentContentHash"] = obj.documentContentHash

	m["givenLocation"] = obj.givenLocation
	m["resolvedLocation"] = obj.resolvedLocation
	m["standardizedLocation"] = obj.standardizedLocation

	m["documentLocationHash"] = obj.documentLocationHash
	m["processedDoc"] = obj.processedDoc

	return m

}
