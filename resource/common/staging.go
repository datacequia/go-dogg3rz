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

package common

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"reflect"
	"time"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

// StagingResourceLocation is a struct whose methods collectively describe
// the location of a stageable JSON-LD resource object within a JSON-LD document
//
// Note that some of the resource objects
// , as defined in github.com/datacequia/go-dogg3rz/resource/jsonld/JSONLDResourceType,
// do not have unique IRIs themselves such as Context and Dataset objects. but can be
// uniquely identified relative to their container which either contain be uniquely
// identified by  an IRI (e.g. Node) or is singularly occuring (Dataset)

type JSONLDDocumentLocation struct {
	ObjectType    jsonld.JSONLDResourceType // TYPE OF JSON-LD OBJECT TO STAGE
	ObjectIRI     string                    // IRI (@id) of JSON-LD OBJECT IF TYPE IS Node or NamedGraph
	ContainerType jsonld.JSONLDResourceType // THE JSON-LD CONTAINER OBJECT TYPE OF ObjectType
	ContainerIRI  string                    // THE IRI OF CONTAINER OBJECT

}

type StagingResourceLocation struct {
	JSONLDDocumentLocation
	DatasetPath string // REPOSITORY PATH TO DATASET
}

// StagingResource is a struct whose members describe the location within
// the JSON-LD document the resource object resided at the time it was modified
// (LastModifiedTimeNS), and the state of the object as identified by it's
// IPFS/IPLD CID (ObjectCID)
type StagingResource struct {
	StagingResourceLocation

	LastModifiedNs int64  // THE LAST TIME JSON-LD OBJECT WAS MODIFIED IN NANOSECONDS SINCE EPOCH
	ObjectCID      string // THE IPFS/IPLD CONTENT IDENTFIER  OF THE STAGED OBJECT
	//
}

// StagingResourceCollector is an interface that is called whenever a stageable
// resource is encountered by FindStageablResources() while traversing
// a JSON-LD document. CollectStart() is called when a given stageable resoource
// located at location is first encountered and CollectEnd() is called if  CollectStart()
// returns without error  AND other child/nested  stageable resources have been encountered
// within location (if any)
type StagingResourceCollector interface {
	CollectStart(ctxt context.Context, resource interface{}, location StagingResourceLocation) error
	CollectEnd(ctxt context.Context, resource interface{}, location StagingResourceLocation)
}

// AssertValid checks that the state of StagingResourceLocation  is in a valid state
// Not all permutations are an acceptable state
func (sr StagingResourceLocation) AssertValid() error {

	var err error

	unHandledErr := errors.UnhandledValue.Newf("type %T, value %v", sr.ObjectType,
		sr.ObjectType)

	invalidObjectTypeContainerTypeErr := errors.InvalidValue.Newf("%T.ObjectType=%s, %T.ContainerType=%s",
		sr.ObjectType, sr.ObjectType, sr.ContainerType, sr.ContainerType)

	if err = sr.ObjectType.AssertValid(); err != nil {
		return err
	}

	if err = sr.ContainerType.AssertValid(); err != nil {
		return err
	}

	switch sr.ObjectType {

	case jsonld.ContextResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:
		case jsonld.NodeResource:
		case jsonld.NamedGraphResource:
		case jsonld.DatasetResource:
		default:
			err = unHandledErr
		}
	case jsonld.NodeResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.NodeResource:
		case jsonld.NamedGraphResource:
		case jsonld.DatasetResource:
		default:
			err = unHandledErr
		}

	case jsonld.NamedGraphResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.NodeResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.NamedGraphResource:
		case jsonld.DatasetResource:

		default:
			err = unHandledErr
		}

	case jsonld.DatasetResource:
		switch sr.ContainerType {
		case jsonld.ContextResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.NodeResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.NamedGraphResource:
			err = invalidObjectTypeContainerTypeErr
		case jsonld.DatasetResource:
		default:
			err = unHandledErr
		}

		return err

	}

	if !sr.CanStage() {
		err = errors.InvalidValue.Newf("%T object cannot be staged where %T.ObjectType="+
			"%s and %T.ContainerType=%s", sr, sr, sr.ObjectType, sr, sr.ContainerType)
	} else {

	}

	return err
}

func (sr StagingResourceLocation) String() string {

	return fmt.Sprintf("StagingResourceLocation { DatasetPath: %s, "+
		"ObjectType: %s, ObjectIRI: %s, ContainerType: %s, ContainerIRI: %s }",
		sr.DatasetPath, sr.ObjectType, sr.ObjectIRI,
		sr.ContainerType, sr.ContainerIRI)
}

// AssertValid checks that the state of StagingResource  is in a valid state
// Not all permutations are an acceptable state
func (sr StagingResource) AssertValid() error {

	if err := sr.StagingResourceLocation.AssertValid(); err != nil {
		return err
	}

	if sr.WantsCID() {
		if len(sr.ObjectCID) < 1 {
			return errors.UnexpectedValue.Newf(
				"expected  populated %T.ObjectCID value when  %T.ObjectType=%s and %T.ContainerType=%s",
				sr, sr, sr.ObjectType, sr, sr.ContainerType)
		}
	} else {

		if len(sr.ObjectCID) > 0 {
			return errors.UnexpectedValue.Newf(
				"expected  unpopulated %T.ObjectCID value when  %T.ObjectType=%s and %T.ContainerType=%s, found '%s'",
				sr, sr, sr.ObjectType, sr, sr.ContainerType, sr.ObjectCID)
		}

	}

	return nil

}

// CanStage function tells whether the StagiongResource object in
func (sr StagingResourceLocation) CanStage() bool {

	var answer bool

	switch sr.ObjectType {

	case jsonld.ContextResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:
			answer = true
		case jsonld.DatasetResource:
			answer = true
		}

	case jsonld.NodeResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:
			answer = true
		case jsonld.DatasetResource:
			answer = true
		}

	case jsonld.NamedGraphResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:
			answer = true
		case jsonld.DatasetResource:
			answer = true
		}

	case jsonld.DatasetResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:

		case jsonld.DatasetResource:
			answer = true
		}

	}
	return answer
}

// WantsCID determines whether or not the resource object value should
// be staged to IPFS with an IPLD object that represents it's contents
// and resulting CID value be assigned
// to attribute value StagingResource.ObjectCID based on the state attribute values
//  StagingResource.ObjectType and StagingResource.ContainerType
func (sr StagingResourceLocation) WantsCID() bool {

	var answer bool

	switch sr.ObjectType {

	case jsonld.ContextResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:

		case jsonld.DatasetResource:
			answer = true
		}

	case jsonld.NodeResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:
			answer = true
		case jsonld.DatasetResource:
			answer = true
		}

	case jsonld.NamedGraphResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:

		case jsonld.DatasetResource:
		}

	case jsonld.DatasetResource:

		switch sr.ContainerType {
		case jsonld.ContextResource:

		case jsonld.NodeResource:

		case jsonld.NamedGraphResource:

		case jsonld.DatasetResource:
		}

	}

	return answer
}

func (jdl JSONLDDocumentLocation) GenerateSHA256Key() (string, error) {

	var err error

	h := sha256.New()

	bb := &bytes.Buffer{}

	if err = bb.WriteByte(byte(jdl.ObjectType)); err != nil {
		return "", err
	}

	if _, err = bb.WriteString(jdl.ObjectIRI); err != nil {
		return "", err
	}

	if err = bb.WriteByte(byte(jdl.ContainerType)); err != nil {
		return "", err
	}
	if _, err = bb.WriteString(jdl.ContainerIRI); err != nil {
		return "", err
	}

	if _, err = h.Write(bb.Bytes()); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}

// FindStageableResources searches JSON-LD document doc for stageable
// resources within the document. When found, it calls srCollector.Collect()
func FindStageableResources(ctxt context.Context, datasetPath string, doc map[string]interface{}, srCollector StagingResourceCollector) error {

	var err error

	var loc StagingResourceLocation

	loc.DatasetPath = datasetPath
	loc.ContainerIRI = ""
	loc.ContainerType = jsonld.DatasetResource
	loc.ObjectIRI = ""
	loc.ObjectType = jsonld.DatasetResource

	if err = findInJSONDocumentObject(ctxt, doc, loc, 0, srCollector); err != nil {
		return err
	}

	return nil

}

func UpdateResourceMtimeToNow(doc map[string]interface{}, loc JSONLDDocumentLocation) error {
	return UpdateResourceMtime(doc, loc, time.Now().UnixNano())

}
func UpdateResourceMtime(doc map[string]interface{}, loc JSONLDDocumentLocation, lastModifiedNs int64) error {

	var m interface{}
	var mtimesMap map[string]interface{}
	var ok bool
	var err error

	// GET MTIMES MAP ENTRY FROM JSON-LD DOC
	if m, ok = doc[jsonld.MtimesEntryKeyName]; !ok {
		return errors.NotFound.Newf("JSON-LD document entry '%s'",
			jsonld.MtimesEntryKeyName)
	} else {
		if mtimesMap, ok = m.(map[string]interface{}); !ok {
			return errors.UnexpectedType.Newf(
				"expected to find %T type value for JSON-LD document entry '%s', found %T",
				mtimesMap, jsonld.MtimesEntryKeyName, m)

		}
	}

	var key string

	if key, err = loc.GenerateSHA256Key(); err != nil {
		return err
	}

	mtimesMap[key] = lastModifiedNs

	return nil

}

// findInJSONDocumentObject seaarches recursively for map objects starting with
// obj which is identified as objectType by the caller and contaianed within
// a container of containerType as determined by the caller
func findInJSONDocumentObject(ctxt context.Context, data interface{}, loc StagingResourceLocation, objectLevel int, srCollector StagingResourceCollector) error {

	//fmt.Println("level", objectLevel, "obj", data)

	// EVAL obj map entries for JSON-LD keywords that
	// map to jsonld.JSONLDResourceType 's
	var obj map[string]interface{}
	var err error
	var ok bool

	// SEND STAGEABLE DATA TO COLLECTOR  IF DATA STAGING LOCATION
	// INDICATES IT CAN BE STAGED
	//fmt.Println("findInJSONDocumentObject", loc)
	if loc.CanStage() {
		if err = srCollector.CollectStart(ctxt, data, loc); err != nil {
			return err
		}
		// NO ERRORS DURING COLLECT STAART.
		// CALL CollectEnd() BEFORE EXIT
		defer srCollector.CollectEnd(ctxt, data, loc)
	}

	switch loc.ObjectType {
	// what object type am i?
	case jsonld.DatasetResource:
		// CHECK FOR (OPTIONAL) CONTEXT

		if ok, obj = isObject(data); !ok {
			return errors.UnexpectedType.Newf("expected DatasetResource to be type map[string]interface{},  found type %T: { loc = %s }",
				data, loc)
		}

		if contextValue, ok := hasEntryKey("@context", obj); ok {

			var contextLoc StagingResourceLocation

			contextLoc.ContainerIRI = loc.ObjectIRI
			contextLoc.ContainerType = loc.ObjectType
			contextLoc.ObjectIRI = ""
			contextLoc.ObjectType = jsonld.ContextResource
			contextLoc.DatasetPath = loc.DatasetPath

			if err = findInJSONDocumentObject(ctxt, contextValue, contextLoc,
				objectLevel+1, srCollector); err != nil {
				return err
			}
		}

		if err = handleGraphProperty(ctxt, obj, loc, objectLevel, srCollector); err != nil {
			return err
		}

	case jsonld.ContextResource, jsonld.NodeResource:

	case jsonld.NamedGraphResource:

		if ok, obj = isObject(data); !ok {
			return errors.UnexpectedType.Newf("expected DatasetResource to be type map[string]interface{},  found type %T: { loc = %s }",
				data, loc)
		}

		if err = handleGraphProperty(ctxt, obj, loc, objectLevel, srCollector); err != nil {
			return err
		}

	} // end switch

	return nil
}

func hasEntryKey(entryKey string, obj map[string]interface{}) (interface{}, bool) {

	if v, ok := obj[entryKey]; ok {
		return v, true
	}

	return nil, false

}

func hasEntryWithMapValue(entryKey string, obj map[string]interface{}) (map[string]interface{}, bool) {

	if v, ok := hasEntryKey(entryKey, obj); ok {

		if cMap, ok := v.(map[string]interface{}); ok {

			return cMap, true
		}
	}

	return nil, false
}

func hasEntryWithArrayValue(entryKey string, obj map[string]interface{}) ([]interface{}, bool) {

	var v interface{}
	var ok bool
	if v, ok = hasEntryKey(entryKey, obj); ok {
		if cSlice, ok := v.([]interface{}); ok {
			return cSlice, true
		}
	}

	return nil, false
}

// isNodeObject determines whethe  obj is a node object as defined
// here https://www.w3.org/TR/json-ld11/#node-objects where caller must
// tell function of obj occurs as the top-most map (topMostMapInJSONLDDoc)
// in the JSON-LD document
func isNodeObject(obj map[string]interface{}, existsOutsideContextObject bool, topMostMapInJSONLDDoc bool) bool {

	// DOES OBJECT HAVE ANY OTHER ENTRIES OTHER THAN @graph or @context

	var otherEntryCount int
	var hasGraphKeyword bool
	var hasContextKeyword bool
	var hasValueListSetKeywords bool

	for k, v := range obj {
		switch k {
		case "@graph":
			hasGraphKeyword = true
		case "@context":
			if !isContextValue(v) {
				return false
			}

			hasContextKeyword = true
		case "@value", "@list", "@set":
			hasValueListSetKeywords = true
		default:
			otherEntryCount++
		}
	}

	var isGraphObj bool

	isGraphObj, _ = isGraphObject(obj, existsOutsideContextObject, topMostMapInJSONLDDoc)

	//  A map is a node object if it exists outside of a JSON-LD context and
	if existsOutsideContextObject &&
		// it is not the top-most map in the JSON-LD document consisting of no other entries than @graph and @context,
		(!(topMostMapInJSONLDDoc && hasGraphKeyword && hasContextKeyword && otherEntryCount == 0)) &&
		// t does not contain the @value, @list, or @set keywords, and
		(!hasValueListSetKeywords) &&
		// t is not a graph object.
		!isGraphObj {
		return true
	}

	return false

}

// isGraphObject evaluates obj against rules that define a graph object
// as defined here https://www.w3.org/TR/json-ld11/#graph-objects
// where caller must
// tell function of obj occurs as the top-most map (topMostMapInJSONLDDoc)
// in the JSON-LD document and if it exists outside of a context object
func isGraphObject(obj map[string]interface{}, existsOutsideContextObject bool, topMostMapInJSONLDDoc bool) (bool, string) {

	var hasGraphKeyword bool
	var hasOtherEntries bool
	var idValue string

	for k, v := range obj {
		switch k {
		case "@graph":
			// TODO: NEEED TO RESOLVEE FOR ALIASES TO @graph keyword as well
			hasGraphKeyword = true

		case "@index":

		case "@id":
			if s, ok := v.(string); ok {
				idValue = s
			}

		case "@context":

			if !isContextValue(v) {
				return false, ""
			}

		default:

			hasOtherEntries = true
		}
	}

	return existsOutsideContextObject && hasGraphKeyword && (!topMostMapInJSONLDDoc) &&
		(!hasOtherEntries), idValue

}

// isContextValue determines if value assigned to @context keyword in a node object
// conforms to context definition defined here https://www.w3.org/TR/json-ld11/#context-definitions
func isContextValue(v interface{}) bool {

	switch reflect.ValueOf(v).Kind() {
	case reflect.Invalid:
	// Nil
	case reflect.String:
	case reflect.Map:
	// context value is map
	case reflect.Array, reflect.Slice:

		// TODO: EVAL THE ELEMENT TYPE IN ARRAY AND SEE IF IT IS
		//       EITHER INVALID,STRING OR MAP

	default:
		// NONE OF THE ABOVE. INVALID
		return false

	}

	return true

}

// isObject determines if data is a map[string]interface{} object
// returns true with the data cast as a map[string]{interface} if true
func isObject(data interface{}) (bool, map[string]interface{}) {

	if obj, ok := data.(map[string]interface{}); ok {
		return ok, obj
	}
	return false, nil

}

func handleGraphProperty(ctxt context.Context, obj map[string]interface{}, loc StagingResourceLocation, objectLevel int, srCollector StagingResourceCollector) error {

	var array []interface{}
	var ok bool
	var err error
	//	fmt.Println("handleGraphProperty", "objectLevel", objectLevel)

	// CHECK FOR default @graph
	if array, ok = hasEntryWithArrayValue("@graph", obj); !ok {
		var v interface{}

		v, _ = obj["@graph"]

		return errors.NotFound.Newf("can't find @graph property value is not array:"+
			" { loc = %s, @graph value type = %T }",
			loc, v)

	}

	// ITERATE OBJECTS IN DEFAULT GRAPH
	for _, item := range array {
		if m, ok := item.(map[string]interface{}); ok {
			//	fmt.Println("iterate graph entries", "entry", m)
			if ok, _ := isGraphObject(m, loc.ContainerType != jsonld.ContextResource,
				objectLevel+1 == 0); ok {

				var graphLoc StagingResourceLocation

				graphLoc.ContainerIRI = loc.ObjectIRI
				graphLoc.ContainerType = loc.ObjectType
				graphLoc.ObjectIRI = ""
				graphLoc.ObjectType = jsonld.NamedGraphResource
				graphLoc.DatasetPath = loc.DatasetPath

				if graphLoc.ObjectIRI, err = getIdPropertyValue(m, graphLoc); err != nil {
					return err
				}

				if err = findInJSONDocumentObject(ctxt, m, graphLoc,
					objectLevel+1, srCollector); err != nil {
					return err
				}

			} else if isNodeObject(m, loc.ContainerType != jsonld.ContextResource,
				objectLevel+1 == 0) {

				var nodeLoc StagingResourceLocation

				nodeLoc.ContainerIRI = loc.ObjectIRI
				nodeLoc.ContainerType = loc.ObjectType
				nodeLoc.ObjectIRI = ""
				nodeLoc.ObjectType = jsonld.NodeResource
				nodeLoc.DatasetPath = loc.DatasetPath

				if nodeLoc.ObjectIRI, err = getIdPropertyValue(m, nodeLoc); err != nil {
					return err
				}

				if err = findInJSONDocumentObject(ctxt, m, nodeLoc,
					objectLevel+1, srCollector); err != nil {
					return err
				}
			}

		}

	}

	return nil

}

func getIdPropertyValue(obj map[string]interface{}, loc StagingResourceLocation) (string, error) {

	if val, ok := hasEntryKey("@id", obj); ok {
		if valStr, ok := val.(string); ok {
			return valStr, nil
		} else {
			return "", errors.UnexpectedType.Newf(
				"expected Node object to have @id property with string "+
					"type value, found %T type: { loc = %s }",
				val, loc)
		}
	} else {
		return "", errors.NotFound.Newf(
			"expected Node object to have @id property: "+
				"{ loc = %s }",
			loc)

	}
}
