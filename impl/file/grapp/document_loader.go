package grapp

import (
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"hash"
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
	httpClient *http.Client // optional: http client to use to load documents
	grappDir   string       // base project dir
	objectsDir string       // where to place object files
	//cachedDocumentIndex map[string]
}

type CachedDocument struct {
	//ld.RemoteDocument

	cacheDir   string        // baseDir to cached document files
	tmpFile    *os.File      // Open file where cache writes are sent
	realReader io.ReadCloser // actual reader that is proxied
	hash       hash.Hash     // hash object that computes hash of cached doc

}

func NewDocumentLoader(httpClient *http.Client, grappDir string, objectsDir string) *DocumentLoader {
	rval := &DocumentLoader{httpClient: httpClient, grappDir: grappDir, objectsDir: objectsDir}

	if rval.httpClient == nil {
		rval.httpClient = http.DefaultClient
	}

	return rval
}

// Loads JSON-LD documents from local or http paths
// Implements github.com/piprate/ld/DocumentLoader interface
func (dl *DocumentLoader) LoadDocument(u string) (*ld.RemoteDocument, error) {

	//fmt.Println("loading document...", u)
	/*
		f := func() {
			fmt.Println("exiting loading document ", u)
		}
		defer f()
	*/
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, fmt.Sprintf("error parsing URL: %s", u))
	}

	var documentBody io.ReadCloser
	var finalURL, contextURL string
	//var loadedDocument *CachedDocument

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

		documentBody = file
	} else {

		req, err := http.NewRequest("GET", resolveIRI(u), nil)
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

	parsedJSON, _, err := dl.createObjectFile(finalURL, buf)
	if err != nil {
		return nil, err
	}
	//fmt.Println("after createObjectFile returns ", objectFilePath)
	return &ld.RemoteDocument{DocumentURL: finalURL, Document: parsedJSON, ContextURL: contextURL}, nil

}

func (dl *DocumentLoader) createObjectFile(iri string, data []byte) (interface{}, string, error) {

	//fmt.Println("1.", iri)

	// compute hash on iri
	iriHash := crypto.SHA1.New()
	_, err := io.Copy(iriHash, strings.NewReader(iri))
	if err != nil {
		return nil, "", err
	}
	iriHashStr := fmt.Sprintf("%x", iriHash.Sum(nil))
	//fmt.Println("2.", iri)
	// compute hash on document
	docHash := crypto.SHA1.New()
	_, err = io.Copy(docHash, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	//	docHashStr := fmt.Sprintf("%x", hash.Sum(nil))

	//fmt.Println("3.", iri)
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

	//fmt.Println("tmp created at ", tmp.Name())
	//fmt.Println("4.", iri)
	// PARSE DOC CONTENTS
	var jsonTree map[string]interface{}
	jsonTree, err = parseJSON(bytes.NewBuffer(data), iri)
	if err != nil {
		return nil, "", err
	}
	//fmt.Println("5.", iri, len(jsonTree))
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

	//fmt.Println("6.", iri)
	// FLATTEN THE TREE TO AN ARRAY OF N-QUADS
	flattenedDoc, err = proc.Flatten(expandedDoc, nil, options)
	if err != nil {
		return nil, "", err
	}

	//fmt.Println("7.", iri)
	// ENCODE OBJECT TO CBOR FORMAT
	var cborObject []byte

	cborObject, err = cbor.Marshal(flattenedDoc)
	if err != nil {
		return nil, "", err
	}
	//fmt.Println("8.", iri)
	// WRITE CBOR OBJECT DATA TO STAGE FILE
	_, err = io.Copy(tmp, bytes.NewReader(cborObject))
	if err != nil {
		return nil, "", err
	}
	//fmt.Println("9.", iri)
	// CONSTRUCT OBJECT FILE  PATH NAME
	objectFilePath := path.Join(dl.objectsDir, iriHashStr)

	//tmp.Close()

	// MOVE STAGED FILE TO OBJECT FILE  PATH
	if err := os.Rename(tmp.Name(), objectFilePath); err != nil {
		os.Remove(tmp.Name())
		//fmt.Println("rename err", err)
		return nil, "", err
	}
	//fmt.Println("10.", iri)
	return jsonTree, objectFilePath, nil

}

func DocumentFromReader(documentBody io.Reader, src string) (interface{}, error) {

	return parseJSON(documentBody, src)
}

func parseJSON(r io.Reader, src string) (map[string]interface{}, error) {

	jsonMap := make(map[string]interface{})

	parseStats := &jsonParseStats{realReader: r}
	parseStats.newlineOffsets = make([]int64, 0)

	decoder := json.NewDecoder(parseStats)

	err := decoder.Decode(&jsonMap)
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

			return nil, errors.Newf("%s:%d:%d: %s", src, syntaxErrLine, syntaxErrCol, err.Error())

		}
		//fmt.Println("other decode err", err)
		return nil, err
	}

	return jsonMap, nil

}

// resolves Prefix IRIs in an alternate location
// where prefix IRI does not make ontology available in their designated namespace
func resolveIRI(iri string) string {

	switch iri {
	case "https://schema.org/":
		return "https://schema.org/version/latest/schemaorg-current-https.jsonld"
	case "http://schema.org/":
		return "http://schema.org/version/latest/schemaorg-current-https.jsonld"
	default:
		return iri
	}
}

func NewCachedDocument(r io.ReadCloser, cacheDir string, docHashType crypto.Hash) (*CachedDocument, error) {

	var err error

	cachedDoc := &CachedDocument{}
	if !file.DirExists(cacheDir) {

		return nil, errors.NotFound.Newf("%s: not found or is not a directory", cacheDir)
	}
	cachedDoc.cacheDir = cacheDir

	// CREATE CACHE TEMP FILE
	if cachedDoc.tmpFile, err = os.CreateTemp(cacheDir, "NewGrapplicationCachedRemoteDocument-*"); err != nil {
		return nil, err
	}
	// assign reader to be proxied
	cachedDoc.realReader = r
	// create hash object using caller supplied crypto hash type
	cachedDoc.hash = docHashType.New()

	return cachedDoc, nil

}

func (d *CachedDocument) Read(p []byte) (int, error) {
	// pass read request to real reader
	bytesRead, err := d.realReader.Read(p)

	if err != nil {
		if err == io.EOF {
			//Println("is eof")
			if bytesRead > 0 {
				// add final bytes ...
				if _, err = d.cacheBytesReadAndAddToHash(p[:bytesRead]); err != nil {
					return bytesRead, err
				}
			}

			// flush cache file writes
			if err = d.tmpFile.Sync(); err != nil {
				return bytesRead, err
			}

			// MOVE CACHED doc TO FINAL DESTINATION
			// USING HASH OF DOC AS PREFIX FILENAME
			docPath := path.Join(d.cacheDir, fmt.Sprintf("%x.jsonld", d.hash.Sum(nil)))

			if err = os.Rename(d.tmpFile.Name(), docPath); err != nil {
				return bytesRead, err
			}
			//d.cachedDocPath = docPath

			//fmt.Println("moved cache Doc to ", docPath)
			// RETURN  EOF
			return bytesRead, io.EOF

		} else {
			// some i/o error
			return bytesRead, err
		}
	}

	if _, err = d.cacheBytesReadAndAddToHash(p[:bytesRead]); err != nil {
		return bytesRead, err
	}

	return bytesRead, nil

}

func (d *CachedDocument) Close() error {

	err := d.realReader.Close()
	if err != nil {
		d.tmpFile.Close()
		// return first error
		return err
	}

	return d.tmpFile.Close()

}

func (d *CachedDocument) cacheBytesReadAndAddToHash(p []byte) (int, error) {

	// WRITER OBJECTS WHICH COMPUTE DIGEST OF BYTES READ
	// AND WRITE TO TEMP CACHE FILE
	mr := io.MultiWriter(d.hash, d.tmpFile)

	return mr.Write(p)

}
