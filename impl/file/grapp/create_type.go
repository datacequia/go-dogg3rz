//go:build nothing

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
	"context" //"github.com/datacequia/go-dogg3rz/util"
)

func (grapp *FileGrapplicationResource) CreateTypeClass(ctxt context.Context, grappName string, datasetPath string,
	typeID string, subclassOf string,
	label string, comment string) error {

	var fds *fileDataset
	var err error
	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	node := make(map[string]interface{})

	node["@id"] = typeID
	node["@type"] = "rdfs:Class"
	if len(subclassOf) > 0 {
		node["rdfs:subClassOf"] = subclassOf
	}
	if len(label) > 0 {
		node["rdfs:label"] = label
	}
	if len(comment) > 0 {
		node["rdfs:comment"] = comment
	}

	if err = fds.appendNodeToDefaultGraph(ctxt, node); err != nil {
		return err
	}

	return nil

}

func (grapp *FileGrapplicationResource) CreateTypeDatatype(ctxt context.Context,
	grappName string, datasetPath string,
	typeID string, subclassOf string,
	label string, comment string) error {

	var fds *fileDataset
	var err error
	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	node := make(map[string]interface{})

	node["@id"] = typeID
	node["@type"] = "rdfs:Datatype"
	if len(subclassOf) > 0 {
		node["rdfs:subClassOf"] = subclassOf
	}
	if len(label) > 0 {
		node["rdfs:label"] = label
	}
	if len(comment) > 0 {
		node["rdfs:comment"] = comment
	}

	if err = fds.appendNodeToDefaultGraph(ctxt, node); err != nil {
		return err
	}

	return nil
}

func (grapp *FileGrapplicationResource) CreateTypeProperty(ctxt context.Context,
	grappName string, datasetPath string,
	typeID string, subPropertyOf string,
	domain string, _range string, label string, comment string) error {

	var fds *fileDataset
	var err error
	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	node := make(map[string]interface{})

	node["@id"] = typeID
	node["@type"] = "rdf:Property"
	if len(subPropertyOf) > 0 {
		node["rdfs:subPropertyOf"] = subPropertyOf
	}
	if len(domain) > 0 {
		node["rdfs:domain"] = domain
	}
	if len(_range) > 0 {
		node["rdfs:range"] = _range
	}

	if len(label) > 0 {
		node["rdfs:label"] = label
	}
	if len(comment) > 0 {
		node["rdfs:comment"] = comment
	}

	if err = fds.appendNodeToDefaultGraph(ctxt, node); err != nil {
		return err
	}

	return nil
}
