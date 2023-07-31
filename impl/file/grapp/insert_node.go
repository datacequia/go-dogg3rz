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

package grapp

import (
	"context"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

func (grapp *FileGrapplicationResource) InsertNode(ctxt context.Context,
	grappName string, datasetPath string,
	nodeType string, nodeID string, graphName string,
	nodeProperties []string, nodeValues []string) error {

	var fds *fileDataset
	var err error
	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	node := make(map[string]interface{})

	if len(nodeType) > 0 {
		node["@type"] = nodeType
	}

	if len(nodeID) > 0 {
		node["@id"] = nodeID
	} else {
		// THIS IS A BLANK NODE. GENERATE A BLANK NODE ID
		node["@id"] = jsonld.NewBlankNodeID().String()

	}

	if len(nodeProperties) != len(nodeValues) {
		return errors.OutOfRange.Newf("the number of node properties/values "+
			"provided must be equal. %d node properties provided, %d node values provided",
			len(nodeProperties), len(nodeValues))

	}

	for i, _ := range nodeProperties {
		prop := nodeProperties[i]
		val := nodeValues[i]
		node[prop] = val

	}

	if len(graphName) > 0 {
		if err = fds.appendNodeToGraph(ctxt, node, graphName); err != nil {
			return err
		}
	} else {

		if err = fds.appendNodeToDefaultGraph(ctxt, node); err != nil {
			return err
		}
	}
	return nil
}
