package grapp

import (
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/iricache"
	"github.com/fxamacker/cbor/v2"
	"github.com/piprate/json-gold/ld"

	"github.com/PuerkitoBio/purell"
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
	objectFileIndex map[string]ObjectFile // index of processed documents cached as binary object files (cbor format)
	// the key(string) is a normalized prefix IRI and value is the object code of JSON-LD document
	// representing the graph referenced by prefix IRI  plus metadata
	stagedDocumentIndex map[string]bool // index of documents that are currently in process of loading
}

type ObjectFile struct {
	givenLocation        string // iri or (relative) location to project of document (.jsonld) file that was processed in project
	normalizedLocation   string // location standardized , for urls its the value returns from [X], for local paths it's relative path to project dir
	resolvedLocation     string // location resolved to the actual path where the JSON-LD can be dereferenced
	documentLocationHash string // hash of source path
	documentContentHash  string // hash of sourcePath content

	processedDoc interface{} // JSON-LD processor output of document
}

type documentLocation struct {
	given      string // the prefix iri provided in the json-ld doc context
	normalized string // a standardized version of document
	resolved   string // the actual location where document was located
	contextURL string // only set if document location was fetched via http. Set within
	// http response 'Link' header that provides more (metadata) info
	// about document
}

func NewDocumentLoader(httpClient *http.Client, grappDir string, objectsDir string) *DocumentLoader {
	rval := &DocumentLoader{httpClient: httpClient, grappDir: grappDir, objectsDir: objectsDir}

	if rval.httpClient == nil {
		rval.httpClient = http.DefaultClient
	}

	rval.objectFileIndex = make(map[string]ObjectFile)
	rval.stagedDocumentIndex = make(map[string]bool)

	return rval
}

// Loads JSON-LD documents from local or http paths
// Implements github.com/piprate/ld/DocumentLoader interface
func (dl *DocumentLoader) LoadDocument(u string) (*ld.RemoteDocument, error) {

	//var documentBody io.ReadCloser
	var finalURL string
	var docLoc *documentLocation
	var docContent []byte
	var err error
	//fmt.Println("laod document called", u)
	if isStaging, InIndex := dl.stagedDocumentIndex[u]; InIndex && isStaging {
		return nil, errors.InProgress.Newf("already loading: %s", u)
	}

	dl.stagedDocumentIndex[u] = true
	defer func() {
		// remote staging status upon exit
		dl.stagedDocumentIndex[u] = false
	}()

	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, fmt.Sprintf("error parsing URL: %s", u))
	}

	protocol := parsedURL.Scheme

	// default to checking in embedded json-ld doc cache first...
	if docContent, docLoc, err = dl.loadDocumentFromEmbeddedCache(u); err != nil {

		// Error loading doc from embedded cache
		if errors.GetType(err) != errors.NotFound {
			return nil, err
		}

		// Not in cache... See if it's an http(s) url
		if protocol == "http" || protocol == "https" {
			if docContent, docLoc, err = dl.loadDocumentFromHttp(u); err != nil {
				return nil, err
			}
		} else {
			// Assume it's an operating filesystem path (within project folder) by deduction...
			if docContent, docLoc, err = dl.loadDocumentFromOperatingSystemPath(u); err != nil {
				return nil, err
			}
		}

	}

	if len(docContent) < 1 {
		return nil, errors.OutOfRange.Newf("expected JSON-LD document content length to be > 0: %s", u)
	}

	if docLoc == nil {
		panic("docLoc is nil")
	}

	// Create an binary object file version of document
	parsedJSON, _, err := dl.createObjectFile(docLoc, docContent)
	if err != nil && errors.GetType(err) != errors.AlreadyExists {
		return nil, err
	}

	return &ld.RemoteDocument{DocumentURL: finalURL, Document: parsedJSON, ContextURL: docLoc.contextURL}, nil

}

// Load a JSON-LD document from an Operating System file path
func (dl *DocumentLoader) loadDocumentFromOperatingSystemPath(loc string) ([]byte, *documentLocation, error) {

	// Not an http(s) protocol. Try to open as a local file..
	var file *os.File
	var absolutePathGrappDir string
	var err error
	var docLoc = documentLocation{}

	docLoc.given = loc

	file, err = os.Open(loc)
	if err != nil {
		// not a local file
		return nil, nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
	}
	defer file.Close()

	// Get absolute path of grapplication dir
	// so that local document path can be converted to a relative path
	// to grapp dir.
	absolutePathGrappDir, err = filepath.Abs(dl.grappDir)
	if err != nil {
		return nil, nil, err
	}

	// STANDARDIZED LOCATION FOR A PROJECT FILE IS THE RELATIVE WRT TO PROJECT DIR
	if docLoc.normalized, err = toProjectFileSubPath(absolutePathGrappDir, loc); err != nil {
		return nil, nil, err
	}

	docLoc.resolved = filepath.Join(absolutePathGrappDir, docLoc.normalized)

	var content []byte
	if content, err = io.ReadAll(file); err != nil {
		return nil, nil, err
	}

	return content, &docLoc, nil

}

// Load JSON-LD document from embedded file within this program binary
// Returns content and its documentLocation's on success
// Returns github.com/datacequia/go-dogg3rz/errors.NotFound if not found in cache
// otherwise returns some other error
func (dl *DocumentLoader) loadDocumentFromEmbeddedCache(loc string) ([]byte, *documentLocation, error) {

	var content []byte
	var docLoc = documentLocation{}
	var err error

	docLoc.given = loc

	// SINCE NOT CALLING THIS URL AND USING EMBEDDED DOC (IN THIS BINARY)
	// NORMALIZE/STANDARDIZE URL USING LIBRARY CALL BELOW
	if docLoc.normalized, err = normalizePrefixIRI(docLoc.given); err != nil {
		return nil, nil, err
	}

	if content, docLoc.resolved, err = iricache.Get(docLoc.normalized); err != nil {
		return nil, nil, err
	}

	return content, &docLoc, nil

}

// Load JSON-LD document from http(s) URL
// Returns content, content location contexts on success
// returns error object on failure
func (dl *DocumentLoader) loadDocumentFromHttp(location string) ([]byte, *documentLocation, error) {

	var content []byte
	var docLoc = documentLocation{}
	var err error

	docLoc.given = location

	// Document is an http(s) URL
	docLoc.resolved = resolveIRI(location)

	req, err := http.NewRequest("GET", docLoc.resolved, nil)
	if err != nil {
		return nil, nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
	}
	// We prefer application/ld+json, but fallback to application/json
	// or whatever is available
	req.Header.Add("Accept", acceptHeader)

	res, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, nil, ld.NewJsonLdError(ld.LoadingDocumentFailed, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, ld.NewJsonLdError(ld.LoadingDocumentFailed,
			fmt.Sprintf("Bad response status code: %d", res.StatusCode))
	}

	docLoc.normalized = res.Request.URL.String()

	contentType := res.Header.Get("Content-Type")
	linkHeader := res.Header.Get("Link") // See https://www.w3.org/wiki/LinkHeader

	if len(linkHeader) > 0 && contentType != "application/ld+json" {
		header := ld.ParseLinkHeader(linkHeader)[linkHeaderRel]
		if len(header) > 1 {
			return nil, nil, ld.NewJsonLdError(ld.MultipleContextLinkHeaders, nil)
		} else if len(header) == 1 {
			docLoc.contextURL = header[0]["target"]
		}
	}

	if content, err = io.ReadAll(res.Body); err != nil {
		return nil, nil, err
	}

	return content, &docLoc, nil

}

// Callback invoked by Document Loader when a JSON-LD Term Definition
// is encountered. See https://www.w3.org/ns/json-ld#TermDefinition
func (dl *DocumentLoader) OnTermDefinition(termDef ld.TermDefinition) {

	if termDef.Value().IsPrefix {

		// Document references a dependent ontology. Load it
		// See https://www.w3.org/TR/json-ld11/#dfn-prefix
		_, err := dl.LoadDocument(termDef.Value().IdKeyword) // use @id keyword
		if err != nil {

			if dgrzErr, ok := err.(errors.BadDogg3rz); ok {
				if dgrzErr.ErrorType == errors.AlreadyExists {
					//fmt.Println("OnTermDefinition: already exists:", termDef.Value().IdKeyword)

					return
				} else if dgrzErr.ErrorType == errors.InProgress {
					//fmt.Println("OnTermDefinition: in progress: ", termDef.Value().IdKeyword)
					return
				}

			}
			fmt.Println("failed to load prefix namespace ", termDef.Value().IdKeyword, err)
		}
	}
}

// Create object file from JSON-LD Document body
func (dl *DocumentLoader) createObjectFile(docLoc *documentLocation, data []byte) (interface{}, string, error) {

	var obj ObjectFile

	obj.givenLocation = docLoc.given
	obj.resolvedLocation = docLoc.resolved
	obj.normalizedLocation = docLoc.normalized

	// compute hash on standardized path location
	iriHash := crypto.SHA1.New()
	_, err := io.Copy(iriHash, strings.NewReader(obj.normalizedLocation))
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

	// CONSTRUCT OBJECT FILE  PATH NAME
	objectFilePath := path.Join(dl.objectsDir, obj.documentLocationHash)

	// CHECK TO SEE IF OBJECT FILE ALREADY EXISTS IN CACHE
	// BEFORE TRYING TO CREATE A NEW ONE
	if file.FileExists(objectFilePath) {

		return nil, "", errors.AlreadyExists.Newf("%s at %s",
			obj.normalizedLocation, objectFilePath)

	}

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

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// Check if document has a suffix. If it's not '.json(ld)' don't attempt o parse
	suffix := strings.ToLower(path.Ext(docLoc.resolved))
	switch suffix {
	case ".json":
	case ".jsonld":
	case "":
		break
	default:
		return nil, "", errors.NotImplemented.Newf("Cannot handle documents with extension '%s'", suffix)
	}

	// PARSE DOC CONTENTS

	var jsonObject interface{}
	jsonObject, err = parseJSONToMap(bytes.NewBuffer(data), obj.resolvedLocation)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		jsonObject, err = parseJSONToArray(bytes.NewBuffer(data), obj.resolvedLocation)
	}
	if err != nil {
		return nil, "", err
	}

	// RUN JSON-LD PROCESSOR WITH PARSED JSON INPUT
	var expandedDoc []interface{}
	var flattenedDoc interface{}
	options.DocumentLoader = NewDocumentLoader(dl.httpClient, dl.grappDir, dl.objectsDir)
	options.Listener = dl // LISTEN FOR PROCESSING EVENTS

	// EXPAND DOC (I.E. EXPAND JSON-LD TERMS TO FULL IRIs)
	expandedDoc, err = proc.Expand(jsonObject, options)
	if err != nil {
		//fmt.Println("expand failed", err, iri, jsonTree)
		return nil, "", err
	}

	if len(expandedDoc) < 1 {
		return nil, "", errors.NotFound.Newf("No RDF statements found after JSON-LD doc expansion: %s", docLoc.resolved)
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

	// MOVE STAGED FILE TO OBJECT FILE  PATH
	if err := os.Rename(tmp.Name(), objectFilePath); err != nil {
		os.Remove(tmp.Name())
		//fmt.Println("rename err", err)
		return nil, "", err
	}

	// ADD OBJECT TO DOCUMENT INDEX
	dl.objectFileIndex[obj.normalizedLocation] = obj

	return jsonObject, objectFilePath, nil

}

func parseJSONToMap(r io.Reader, src string) (map[string]interface{}, error) {

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

func parseJSONToArray(r io.Reader, src string) ([]interface{}, error) {

	jsonMap := make([]interface{}, 0)

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
// where prefix IRI does not make ontology available in their designated prefix iri
func resolveIRI(iri string) string {

	switch strings.ToLower(iri) { // urls are CASE SENSIVE. NORMALIZE TO LOWER
	case "https://schema.org/":
		return "https://schema.org/version/latest/schemaorg-current-https.jsonld"
	case "http://schema.org/":
		return "http://schema.org/version/latest/schemaorg-current-http.jsonld"
	default:
		return iri
	}
}

func (obj ObjectFile) toMap() map[string]interface{} {

	m := make(map[string]interface{})

	m["documentContentHash"] = obj.documentContentHash

	m["givenLocation"] = obj.givenLocation
	m["resolvedLocation"] = obj.resolvedLocation
	m["standardizedLocation"] = obj.normalizedLocation

	m["documentLocationHash"] = obj.documentLocationHash
	m["processedDoc"] = obj.processedDoc

	return m

}

func toProjectFileSubPath(grappDir string, filePath string) (string, error) {

	//var absLocation string
	var err error

	if fileInfo, err := os.Stat(grappDir); err != nil {
		//fmt.Println("test1")
		return "", err
	} else {
		//fmt.Println("1.2")
		// grapp dir path exists. is it a directory?
		if !fileInfo.IsDir() {
			//fmt.Println("test2", grappDir)
			return "", errors.InvalidValue.Newf("%s: not a directory", grappDir)
		}
	}

	// IS LOCATION ABSOLUTE?
	if !filepath.IsAbs(filePath) {
		//fmt.Println("2.1")

		// IT IS A RELATIVE PATH. JOIN WITH GRAPP PROJECT DIR
		// TO CANONICALIZE PATH
		var absFilePath string

		if absFilePath, err = filepath.Abs(filepath.Join(grappDir, filePath)); err == nil {
			var fileInfo fs.FileInfo
			if fileInfo, err = os.Stat(absFilePath); err == nil {
				// PATH EXISTS. IS IT A FILE?
				if !fileInfo.IsDir() {
					// IS A FILE. RETURN REL
					var projectSubPath string
					if projectSubPath, err = filepath.Rel(grappDir, absFilePath); err == nil {
						//fmt.Println("test4", absFilePath)
						return projectSubPath, nil
					} else {
						// test x
					}
				} else {
					// relative file path is a dir!
					//fmt.Println("test6")
					err = errors.InvalidValue.Newf("%s: not a file", absFilePath)
				}
			} else {
				//fmt.Println("test5")
			}
		} else {
			// abs fails on relpath joined to project dir
		}
	} else {
		// FILEPATH IS ABSOLUTE. IS IT A PROJECT (SUB) PATH?
		var projectSubPath string
		if projectSubPath, err = filepath.Rel(grappDir, filePath); err == nil {
			//fmt.Println("test3")
			return projectSubPath, nil
		}

	}

	if err == nil {
		panic("toProjectFileSubPath: err is not initialized")
	}

	return "", err

}

func normalizePrefixIRI(u string) (string, error) {

	var parsedURL string
	var err error
	var flags purell.NormalizationFlags

	u = strings.TrimSpace(u) // make sure no leading/training spaces

	flags |= purell.FlagLowercaseScheme           // HTTP://host -> http://host, applied by default in Go1.1
	flags |= purell.FlagLowercaseHost             // http://HOST -> http://host
	flags |= purell.FlagUppercaseEscapes          // http://host/t%ef -> http://host/t%EF
	flags |= purell.FlagDecodeUnnecessaryEscapes  // http://host/t%41 -> http://host/tA
	flags |= purell.FlagEncodeNecessaryEscapes    // http://host/!"#$ -> http://host/%21%22#$
	flags |= purell.FlagRemoveDefaultPort         // http://host:80 -> http://host
	flags |= purell.FlagRemoveEmptyQuerySeparator // http://host/path? -> http://host/path

	// Usually safe normalizations
	//FlagRemoveTrailingSlash // http://host/path/ -> http://host/path
	//FlagAddTrailingSlash    // http://host/path -> http://host/path/ (should choose only one of these add/remove trailing slash flags)
	flags |= purell.FlagRemoveDotSegments // http://host/path/./a/b/../c -> http://host/path/a/c

	// Unsafe normalizations
	//FlagRemoveDirectoryIndex   // http://host/path/index.html -> http://host/path/
	//FlagRemoveFragment         // http://host/path#fragment -> http://host/path
	//FlagForceHTTP              // https://host -> http://host
	flags |= purell.FlagRemoveDuplicateSlashes // http://host/path//a///b -> http://host/path/a/b
	//FlagRemoveWWW              // http://www.host/ -> http://host/
	//FlagAddWWW                 // http://host/ -> http://www.host/ (should choose only one of these add/remove WWW flags)
	flags |= purell.FlagSortQuery // http://host/path?c=3&b=2&a=1&b=1 -> http://host/path?a=1&b=1&b=2&c=3

	// Normalizations not in the wikipedia article, required to cover tests cases
	// submitted by jehiah
	flags |= purell.FlagDecodeDWORDHost           // http://1113982867 -> http://66.102.7.147
	flags |= purell.FlagDecodeOctalHost           // http://0102.0146.07.0223 -> http://66.102.7.147
	flags |= purell.FlagDecodeHexHost             // http://0x42660793 -> http://66.102.7.147
	flags |= purell.FlagRemoveUnnecessaryHostDots // http://.host../path -> http://host/path
	flags |= purell.FlagRemoveEmptyPortSeparator  // http://host:/path -> http://host/path

	parsedURL, err = purell.NormalizeURLString(u, flags)

	// add back trailing fragment char (important to preserve for iri prefixes) if originally there
	lastChar := u[len(u)-1:]
	if lastChar == "#" {
		parsedURL = parsedURL + lastChar
	}

	return parsedURL, err

}
