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
)

// GrapplicationResource is an interface the provides all the non-iterative interactions
// a user can perform against a grapplication
type GrapplicationResource interface {
	// CREATE A NEW GRAPPLICATION
	Init(ctxt context.Context, grappDirPath string) error
	Validate(ctxt context.Context) error
	//CreateDataset(ctxt context.Context, grappName string, datasetPath string) error

	//AddNamespaceDataset(ctxt context.Context, grappName string, datasetPath string, term string, iri string) error

	//AddNamespaceNode(ctxt context.Context, grappName string, datasetPath string, nodeID string, term string, iri string) error

	//CreateSnapshot(ctxt context.Context, grappName string) error

	//CreateTypeClass(ctxt context.Context, grappName string, datasetPath string, typeID string, subclassOf string,
	//	label string, comment string) error

	//CreateTypeDatatype(ctxt context.Context, grappName string, datasetPath string, typeID string, subclassOf string,
	//	label string, comment string) error

	//CreateTypeProperty(ctxt context.Context, grappName string, datasetPath string, typeID string, subPropertyOf string,
	//	domain string, _range string, label string, comment string) error

	// INSERT NEW NODE INTO DEFAULT-GRAPH (graphName="") or NAMED-GRAPH (graphName!="")
	// FOR DATASET datasetPath GRAPP grappName
	//InsertNode(ctxt context.Context, grappName string, datasetPath string,
	//	nodeType string, nodeID string, graphName string,
	//	nodeProperties []string, nodeValues []string) error
	// Returns list of data sets in the grapp
	//GetDataSets(ctxt context.Context, grappName string) ([]string, error)

	// INSERT NEW NODE INTO DEFAULT-GRAPH (graphName="") or NAMED-GRAPH (graphName!="")
	// FOR DATASET datasetPath GRAPP grappName
	//CreateNamedGraph(ctxt context.Context, grappName string, datasetPath string, graphName string,
	//	parentGraphName string) error

	// ADD DATASET TO GRAPPLICATION
	//Add(ctxt context.Context, grappName string, path string) error
}

// ALLOWS USER TO STAGE/UNSTAGE (i.e. .Add(), Remove() )  EXISTING (JSON-LD) WORKSPACE RESOURCES ITERATIVELY
// TO GRAPPLICATION IDENTIFIED BY VALUE RETURNED FROM  .Grapplication()
// BEFORE FLUSHING STAGED RESOURCES USING .Commit()
//
// NOTE: CALLER SHOULD PASS THE SAME CANCELLABLE (context.Context.WithCancel()) CONTEXT OBJECT INSTANCE
// TO ResourceStager METHODS WHICH REQUIRE A CONTEXT AND CALL THE CANCEL FUNCTION WHEN DONE TO ENSURE
// ANY ALLOCATED GO-ROUTINES ALLOCATED DURING THE INTERACTION WITH THIS INTERFACE ARE DEALLOCATED

/*
type GrapplicationResourceStager interface {
	Add(ctxt context.Context, sr rescom.StagingResourceLocation) error    // stage an new/existing resource (from workspace)
	Remove(ctxt context.Context, sr rescom.StagingResourceLocation) error // remove resource from staging
	Commit(ctxt context.Context) error                                    // save changes to staging
	Rollback(ctxt context.Context) error                                  // undo changes since last commit
	Close(ctxt context.Context) error                                     // release all resources
	Grapplication() string                                                // return grapplication context for staging operations
}
*/
