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

package repo

import (
	"context"

	rescom "github.com/datacequia/go-dogg3rz/resource/common"
)

// RepositoryResource is an interface the provides all the interactions
// a user can leverage against a repository
type RepositoryResource interface {
	InitRepo(ctxt context.Context, repoName string) error
	//	GetRepoIndex(repoName string) (RepositoryIndex, error)

	// StageResource stages  repository object resources in repository repoName
	// starting at startList locations
	// and all the repository object  resources contained within each specified
	// location, if any, and returns a slice of all the resources staged  within,
	// and including, startList locations
	StageResources(ctxt context.Context, repoName string, startList []rescom.StagingResourceLocation) ([]rescom.StagingResource, error)

	CreateDataset(ctxt context.Context, repoName string, datasetPath string) error

	AddNamespaceDataset(ctxt context.Context, repoName string, datasetPath string, term string, iri string) error

	AddNamespaceNode(ctxt context.Context, repoName string, datasetPath string, nodeID string, term string, iri string) error

	CreateSnapshot(ctxt context.Context, repoName string) error

	CreateTypeClass(ctxt context.Context, repoName string, datasetPath string, typeID string, subclassOf string,
		label string, comment string) error

	CreateTypeDatatype(ctxt context.Context, repoName string, datasetPath string, typeID string, subclassOf string,
		label string, comment string) error

	CreateTypeProperty(ctxt context.Context, repoName string, datasetPath string, typeID string, subPropertyOf string,
		domain string, _range string, label string, comment string) error

	// INSERT NEW NODE INTO DEFAULT-GRAPH (graphName="") or NAMED-GRAPH (graphName!="")
	// FOR DATASET datasetPath REPO repoName
	InsertNode(ctxt context.Context, repoName string, datasetPath string,
		nodeType string, nodeID string, graphName string,
		nodeProperties []string, nodeValues []string) error
}

type GetResourceItem interface {
	GetPath() string // PATH TO RESOURCE
	GetStatus() string
}
