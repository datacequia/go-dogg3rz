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
	"crypto/sha256"
	"encoding/hex"
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

type symbolTable struct {
	subjectIndex map[string]*symbolEntry // key is Node/Subject IRI (i.e. @id value), value is symbolEntry where it occurs
	// in its resolution
	typeIndex map[string][]*symbolEntry // key is Object IRI that indicates an rdf:type(i.e. @type value), value is symbolEntry where it occurs
}

type symbolEntry struct {
	jsonLDNodeObject map[string]interface{} // each value represents a flattened representation
	// of a JSON-LD Node Object (i.e. a node in graph)
	srcObjectFile ObjectFile // src object file info of jsonld node object
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

	projectFileLoaders := make([]*DocumentLoader, 0)

	// process JSON-LD files against JSON-LD processor for well-formedness
	for _, jsonLdFile := range projectFiles {

		loader := NewDocumentLoader(nil, grappDir, objectsDir)

		if _, err := loader.LoadDocument(jsonLdFile); err != nil {
			return err
		}

		projectFileLoaders = append(projectFileLoaders, loader)
	}
	//fmt.Println("#prjFileLoaders", len(projectFileLoaders))
	var symTbl *symbolTable

	symTbl, err = createSymbolTable(projectFileLoaders)
	if err != nil {
		return err
	}

	fmt.Println("symbol table size", len(symTbl.subjectIndex))

	resolveSymbolTable(symTbl, objectsDir)

	return nil

}

func extractNamespacesFromContext(jsonMap map[string]interface{}) (map[string]string, error) {

	m := make(map[string]string)

	return m, nil
}

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

func createSymbolTable(projectFiles []*DocumentLoader) (*symbolTable, error) {

	st := symbolTable{}

	// init symbol entries map
	st.subjectIndex = make(map[string]*symbolEntry)
	st.typeIndex = make(map[string][]*symbolEntry)

	// iterate over project json-ld documents (and their prefix dependencies)
	for _, processedSrcFile := range projectFiles {

		curObjectFileIndex := processedSrcFile.objectFileIndex
		for normalizedPrefixIRI := range curObjectFileIndex {
			objectFile := curObjectFileIndex[normalizedPrefixIRI]

			if jsonLDObjectNodeList, ok := objectFile.processedDoc.([]interface{}); ok {

				// iterate JSON-LD object nodes in flattened json-ld graph structure
				for _, jsonLDNodeObject := range jsonLDObjectNodeList {

					if jsonLDNodeObjectMap, ok := jsonLDNodeObject.(map[string]interface{}); ok {

						//fmt.Println("5")

						if nodeIdStr, err := createGloballyUniqueBlankNodesIds(jsonLDNodeObjectMap, normalizedPrefixIRI); err == nil {
							if _, ok := st.subjectIndex[nodeIdStr]; ok {
								//fmt.Println("already have symbol ", nodeIdStr, "in table: prefix iri = ", normalizedPrefixIRI)
								continue
							}

							if !strings.HasPrefix(nodeIdStr, "_:") {
								if normalizedNodeIdStr, err := normalizePrefixIRI(nodeIdStr); err != nil {
									return nil, err
								} else {
									nodeIdStr = normalizedNodeIdStr
								}
							}

							sEntry := &symbolEntry{}
							st.subjectIndex[nodeIdStr] = sEntry
							sEntry.jsonLDNodeObject = jsonLDNodeObjectMap
							sEntry.srcObjectFile = objectFile

							// check for @type property and if found
							// add to type index

							//fmt.Printf("node types for %s: %d: %v\n", nodeIdStr, len(getNodeTypes(jsonLDNodeObjectMap)), jsonLDNodeObjectMap)

							for _, nodeType := range getNodeTypes(jsonLDNodeObjectMap) {
								sEntry2 := &symbolEntry{}

								/*
									if normalizedNodeType, err := normalizePrefixIRI(nodeType); err != nil {
										return nil, err
									} else {
										nodeType = normalizedNodeType
									}
								*/

								if _, ok := st.typeIndex[nodeType]; !ok {
									st.typeIndex[nodeType] = make([]*symbolEntry, 0)

								}
								st.typeIndex[nodeType] = append(st.typeIndex[nodeType], sEntry2)

								sEntry2.jsonLDNodeObject = jsonLDNodeObjectMap
								sEntry2.srcObjectFile = objectFile

							}
							//fmt.Println("jsonLDNodeObject:", objectFile.givenLocation, "@id:", nodeIdStr)

						} else {
							return nil, err

						}

					} else {
						return nil, errors.UnexpectedType.Newf("Expecting map[string]interface{} element while iterating flattened json-ld list structure, found type %T", jsonLDNodeObject)

					}

				}

			} else {
				return nil, errors.UnexpectedType.Newf("Expecting []interface{} as outtermost element type while iterating flattened json-ld structure, found type %T", objectFile.processedDoc)

			}
		}
	}

	return &st, nil

}

func getNodeTypes(jsondLDNodeObject map[string]interface{}) []string {

	typeList := make([]string, 0)

	if nodeType, ok := jsondLDNodeObject["@type"]; ok {
		//fmt.Printf("type @type: %T\n", nodeType)
		if nodeTypeStr, ok := nodeType.(string); ok {
			// add @type IRI to type index
			typeList = append(typeList, nodeTypeStr)

		} else {
			if nodeTypeList, ok := nodeType.([]interface{}); ok {

				for _, nodeType := range nodeTypeList {
					if nodeTypeStr, ok := nodeType.(string); ok {
						typeList = append(typeList, nodeTypeStr)

					}

				}

			}
		}
	}

	return typeList

}

func getNodeId(jsondLDNodeObject map[string]interface{}) (string, error) {

	var nodeId interface{}
	var nodeIdStr string
	var ok bool

	if nodeId, ok = jsondLDNodeObject["@id"]; !ok {
		return "", errors.NotFound.Newf("expected '@id' property: %v", jsondLDNodeObject)
	}

	if nodeIdStr, ok = nodeId.(string); ok {
		// add @type IRI to type index
		return nodeIdStr, nil
	} else {
		return "", errors.UnexpectedType.Newf("expected node @id property value to be type string, found type %T: %v", nodeId, jsondLDNodeObject)
	}

}

func createGloballyUniqueBlankNodeTerm(originalBlankNodeId string, prefixIRI string) string {
	hasher := sha256.New()

	hasher.Write([]byte(prefixIRI + originalBlankNodeId))

	return hex.EncodeToString(hasher.Sum(nil))

}

// Search jsonLDNodeObjectMap for blank node identifiers (i.e. _:<some ephemeral id>)
// and if found, replace with sha256 sum of ephemeral id + it's prefixIRI (i.e. global id)
// return @id property value of jsonLDNodeObjectMap as  string
func createGloballyUniqueBlankNodesIds(jsonLDNodeObjectMap map[string]interface{}, prefixIRI string) (string, error) {

	// check to see if "@id" property is a generated blank node id,
	// if so replace with a global identifier based on on the
	// sha256 hash of its prefix iri and blank node id

	var nodeIdStr string

	if nodeId, ok := jsonLDNodeObjectMap["@id"]; !ok {
		return "", errors.NotFound.Newf("@id property not found for JSON-LD object node: { prefix iri = '%s'}", prefixIRI)

	} else {
		if s, ok := nodeId.(string); !ok {
			return "", errors.UnexpectedValue.Newf("expected JSON-LD node @id value to be a string found type %T", nodeId)

		} else {
			nodeIdStr = s
		}

	}

	if strings.HasPrefix(nodeIdStr, "_:") {
		// this is a blank node. swap identifier with globally unique one
		nodeIdStr = "_:" + createGloballyUniqueBlankNodeTerm(nodeIdStr[2:], prefixIRI)
		jsonLDNodeObjectMap["@id"] = nodeIdStr
	}

	var f func(interface{})

	f = func(object1 interface{}) {

		//fmt.Printf("object1 type: %T\n", object1)
		switch v := object1.(type) {
		case map[string]interface{}:
			//fmt.Println("createGloballyUniqueBlankNodesIds: found map[string]interface{}")
			// node object is a map. check for @id value

			for p1, o1 := range v {

				if idValue, ok := o1.(string); ok {
					if strings.HasPrefix(idValue, "_:") {
						globallyUniqueBlankNodeId := "_:" + createGloballyUniqueBlankNodeTerm(idValue[2:], prefixIRI)
						v[p1] = globallyUniqueBlankNodeId
					}

				} else {
					f(o1)
				}

			}

		case []interface{}:
			//	fmt.Println("createGloballyUniqueBlankNodesIds: found []interface{}")
			// value is a list. iterate list recursively
			for _, i := range v {
				//fmt.Println("found list:", object1)
				f(i)
			}
		default:
			//fmt.Printf("createGloballyUniqueBlankNodesIds: default type: %T\n", v)
		}

	}

	// iterate predicate values of node object
	for _, object := range jsonLDNodeObjectMap {
		//fmt.Printf("f(): %s %s %T\n", prefixIRI, p, object)
		f(object)

	}

	return nodeIdStr, nil

}

// Iterate symbol references in table to ensure they have been defined
func resolveSymbolTable(symTbl *symbolTable, objectsDir string) error {

	//for _, symbolEntry := range symTbl.subjectIndex {

	//findSymbolReferences(symbolEntry.jsonLDNodeObject, symTbl)
	resolveTypes(symTbl, objectsDir)
	//}
	return nil
}

func resolveTypes(symTbl *symbolTable, objectsDir string) error {
	// if an id has a type but type prefix iri is not imported then fail

	//fmt.Println("type index size ", len(symTbl.typeIndex))

	for nodeType, symbolEntries := range symTbl.typeIndex {
		//fmt.Println("iterate ", nodeType)
		/*()
		lastPathIndexStart := strings.LastIndex(nodeType, "#")
		if lastPathIndexStart < 0 {
			lastPathIndexStart = strings.LastIndex(nodeType, "/")
		}
		*/

		if _, ok := symTbl.subjectIndex[nodeType]; !ok {
			// type NOT found as subject (i.e. not defined)
			for _, se := range symbolEntries {

				if nodeId, err := getNodeId(se.jsonLDNodeObject); err == nil {
					fmt.Printf("type symbol '%s' not found: referenced from %s: { '@id': '%s'} \n", nodeType, se.srcObjectFile.givenLocation, nodeId)
				} else {
					return err // error getting @id attribute
				}
			}
		}

	}

	return nil

}
func findSymbolReferences(obj interface{}, symTbl *symbolTable) {

	switch v := obj.(type) {

	case map[string]interface{}:

		for p, o := range v {

			switch p {
			case "@id":
				if oStr, ok := o.(string); ok {

					// is iri reference in symtbl?
					if _, ok := symTbl.subjectIndex[oStr]; !ok {
						fmt.Println("cannot resolve symbol ", oStr, "referenced in ", v)
					}
				}

			case "@type":

			default:
				findSymbolReferences(o, symTbl)
			}

		}
	case []interface{}:
		for _, i := range v {
			findSymbolReferences(i, symTbl)
		}
	}
}
