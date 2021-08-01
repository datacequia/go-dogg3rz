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
	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

const defaultGraphID = "default"

func GetDefaultGraph(m map[string]interface{}) (map[string]interface{}, []interface{}, error) {

	// GET DEFAULT GRAPH OBJECT FROM DOC
	var defaultGraphRaw interface{}
	var success bool

	// ENSURE DEFAULT GRAPH EXISTS
	if defaultGraphRaw, success = m["@graph"]; !success {
		return nil, nil, errors.NotFound.New(
			"default graph not found as outermost " +
				"attribute '@graph' in JSON-LD document")

	}
	var defaultGraph []interface{}
	if defaultGraphRaw == nil || defaultGraphRaw == "" {
		defaultGraph = make([]interface{}, 0)
	} else {

		if defaultGraph, success = defaultGraphRaw.([]interface{}); !success {
			return nil, nil, errors.InvalidValue.New(
				"Connot convert default graph to list of interface")

		}
	}

	return m, defaultGraph, nil
}

func GetNodeID(newNode map[string]interface{}) (string, error) {
	// CHECK IF NEW NODE HAS ID
	var newNodeIDValue string

	if newNodeID, ok := newNode["@id"]; ok {
		// NEW NODE HAS ID. EXTRACT AS TYPE string
		if newNodeIDStr, ok := newNodeID.(string); ok {
			newNodeIDValue = newNodeIDStr

		} else {
			return "", errors.UnexpectedType.Newf(
				"expected @id value of node to append to be type %T, found type %T",
				newNodeIDValue, newNodeID)
		}
	} else {
		newNodeIDValue = ""
	}
	return newNodeIDValue, nil
}

func AddNodeToGraph(defaultGraphMap *map[string]interface{}, newNode map[string]interface{}, parentGraphID string) error {
	//find the parent nodeID

	var err error
	if _, err = GetNodeID(newNode); err != nil {
		return errors.InvalidValue.Newf("Node does not have @id attribute, cannot add node to graph. details: '%s'", err.Error())
	}

	var err2 error
	if err2 = updateDefaultGraph(defaultGraphMap, parentGraphID, newNode); err2 != nil {
		return err2
	}

	return nil

}

func GetGraph(graph []interface{}, graphID string) (interface{}, error) {

	if len(graph) == 0 && graphID == defaultGraphID {
		return make([]interface{}, 0), nil
	}
	for _, node := range graph {

		var nodeAsMap map[string]interface{}
		var isMap bool
		if nodeAsMap, isMap = node.(map[string]interface{}); isMap {
			// if graphID is blank then check that @id is not there

			curID, ok := nodeAsMap["@id"]

			if (graphID == defaultGraphID && !ok) || (ok && getIDValue(curID) == graphID) {

				childGraphRaw, _ := nodeAsMap["@graph"]
				return childGraphRaw, nil
			}

			childGraphRaw, success1 := nodeAsMap["@graph"]
			if !success1 {
				return nil, errors.NotFound.New(
					"Graph node found in Graph: " + graphID)

			}

			childGraph, success2 := childGraphRaw.([]interface{})
			if success2 {
				if len(childGraph) != 0 {
					value, _ := GetGraph(childGraph, graphID)
					if value != nil {
						return value, nil
					}
				}

			}

		}

	}
	return nil, errors.NotFound.New(
		"Graph: " + graphID + " not found in the map")

}

func updateDefaultGraph(defaultGraphMap *map[string]interface{}, graphIDToUpdate string, newNode map[string]interface{}) error {

	if graphIDToUpdate == defaultGraphID {
		err := appendToGraph(defaultGraphMap,defaultGraphMap, newNode)
		return err
	}

	var graphRaw interface{}
	var success bool

	if graphRaw, success = (*defaultGraphMap)["@graph"]; !success {
		return errors.NotFound.New(
			" graph not found as outermost " +
				"attribute '@graph' in JSON-LD document")

	}
	var graph []interface{}
	if graphRaw == nil || graphRaw == "" {
		graph = make([]interface{}, 0)
	} else {

		if graph, success = graphRaw.([]interface{}); !success {
			return errors.InvalidValue.New(
				"Connot convert  graph to list of interface")

		}
	}
	for _, node := range graph {

		var nodeAsMap map[string]interface{}
		var isMap bool
		if nodeAsMap, isMap = node.(map[string]interface{}); isMap {
			// if graphID is blank then check that @id is not there

			childGraphRaw, success1 := nodeAsMap["@graph"]
			curID, ok := nodeAsMap["@id"]

			if success1 && ok {

				childGraph, success2 := childGraphRaw.([]interface{})
				if !success2 {
					return errors.NotFound.New(
						"Parent Graph: " + graphIDToUpdate + " graph node cannot be converted to list")

				}

				if ok && getIDValue(curID) == graphIDToUpdate {
					return  appendToGraph(defaultGraphMap,&nodeAsMap, newNode)



				}

				if len(childGraph) != 0 {
					if err1 := updateDefaultGraph(&nodeAsMap, graphIDToUpdate, newNode); err1 == nil {
						return nil
					}
				}
			}
		}

	}
	return errors.NotFound.New(
		"Parent Graph: " + graphIDToUpdate + " not found in the map")

}

func appendToGraph(defaultGraphMap *map[string]interface{},parentNodeMap *map[string]interface{}, newNode map[string]interface{}) error {

	var graphRaw interface{}
	var success bool

	if graphRaw, success = (*parentNodeMap)["@graph"]; !success {
		return errors.NotFound.New(
			" graph not found as outermost " +
				"attribute '@graph' in JSON-LD document")

	}
	var graph []interface{}
	if graphRaw == nil || graphRaw == "" {
		graph = make([]interface{}, 0)
	} else {

		if graph, success = graphRaw.([]interface{}); !success {
			return errors.InvalidValue.New(
				"Connot convert  graph to list of interface")

		}
	}

	(*parentNodeMap)["@graph"] = append(graph, newNode)

	// update the MTIME
	nodeID, err2 := GetNodeID(newNode)
	if err2 != nil {
		return err2
	}
	return updateMTIME(defaultGraphMap, nodeID)

	
}

func getIDValue(id interface{}) string {
	if value, ok := id.(string); ok {
		return value
	}
	return ""

}

func updateMTIME(m *map[string]interface{}, newNodeIDValue string) error {

	if _, ok := (*m)[jsonld.MtimesEntryKeyName]; !ok {
		(*m)[jsonld.MtimesEntryKeyName] = make(map[string]interface{})
	}

	// UPDATE MTIME FOR THIS RESOURCE (NODE)
	var loc JSONLDDocumentLocation

	loc.ContainerType = jsonld.DatasetResource // SINCE IT'S DEFAULT GRAPH, CONTAINER IS DATASET
	loc.ContainerIRI = ""
	loc.ObjectType = jsonld.NodeResource
	loc.ObjectIRI = newNodeIDValue

	if err := UpdateResourceMtimeToNow(*m, loc); err != nil {
		return err
	}
	return nil
}
