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

package cmd

import (
	"context"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

// STAGE  DEFAULT-GRAPH | GRAPH IRI | NODE IRI
type dgrzStageCmd struct {
	Grapplication string `long:"grapplication" short:"g" env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`

	All     dgrzStageAllCmd     `command:"all" description:"stage all resources in a grapplication"`
	Dataset dgrzStageDatasetCmd `command:"dataset" alias:"ds" description:"stage objects in a dataset"`
}

type dgrzStageAllCmd struct {
}

type dgrzStageDatasetCmd struct {
	Positional struct {
		DatasetPath string `positional-arg-name:"DATASET_PATH" required:"yes" `
	} `positional-args:"yes"`

	All        dgrzStageAllCmd     `command:"all" description:"stage all resources contained in a dataset"`
	Context    dgrzStageContext    `command:"context"  description:"stage the outermost context resource (only) in a dataset"`
	NamedGraph dgrzStageNamedGraph `command:"named-graph" alias:"ng" description:"stage all resources contained in a named graph"`
	Node       dgrzStageNode       `command:"node" description:"stage a single node resource in a dataset"`
}

type dgrzStageContext struct {
}

type dgrzStageNamedGraph struct {
	Positional struct {
		IRI string `positional-arg-name:"IRI" required:"yes" `
	} `positional-args:"yes"`
	ParentGraphIRI string `long:"parent" description:"Parent graph IRI" required:"false" default:""`
}

type dgrzStageNode struct {
	Positional struct {
		IRI string `positional-arg-name:"IRI" required:"yes" `
	} `positional-args:"yes"`
	ParentGraphIRI string `long:"parent" description:"Parent graph IRI" required:"false" default:""`
}

var stageCmd = dgrzStageCmd{}

func init() {
	// REGISTER THE 'stage' COMMAND
	register(&stageCmd)
}

//////////////////////////////////////
// EXECUTE ENTRYPOINTS
/////////////////////////////////////

// commaand entrypoint for staging either all 1) all datasets in a grapplication
// or 2) stage all resources in a dataset
func (cmd *dgrzStageAllCmd) Execute(args []string) error {

	ctxt, cancelFunc := context.WithCancel(getCmdContext())
	defer cancelFunc()

	stager, err := resource.GetGrapplicationResourceStager(ctxt, stageCmd.Grapplication)
	if err != nil {
		return err
	}
	defer stager.Close(ctxt)

	stageDataset := func(ctxt context.Context, datasetName string) error {

		srl := rescom.StagingResourceLocation{}
		srl.ContainerType = jsonld.DatasetResource
		srl.ContainerIRI = "" // n/a  for type
		srl.ObjectType = jsonld.DatasetResource
		srl.ObjectIRI = "" // n/a for type
		srl.DatasetPath = datasetName

		// STAGE DATASET
		if err := stager.Add(ctxt, srl); err != nil {
			if err2 := stager.Rollback(ctxt); err2 != nil {
				return errors.Wrap(err, err2.Error())
			}
			return err

		}

		return nil

	}

	switch cmd {
	case &stageCmd.All:
		// STAGE ALL DATASETS IN A GRAPPLICATION

		grapp := resource.GetGrapplicationResource(ctxt)

		datasets, err := grapp.GetDataSets(ctxt, stageCmd.Grapplication)
		if err != nil {
			return err
		}

		if len(datasets) < 1 {
			return errors.NotFound.New("no datasets to stage")
		}

		for _, ds := range datasets {
			if err := stageDataset(ctxt, ds); err != nil {
				return err
			}
		}

	case &stageCmd.Dataset.All:

		// STAGE A SPECIFIC DATASET WITHIN A GRAPPLICATION

		if err := stageDataset(ctxt, stageCmd.Dataset.Positional.DatasetPath); err != nil {
			return err
		}

	default:
		panic("unhandled 'stage all' command context")
	}

	return stager.Commit(ctxt)
}

// command entrypoint for staging outermost context in a json-ld dataset
func (c *dgrzStageContext) Execute(args []string) error {

	ctxt, cancelFunc := context.WithCancel(getCmdContext())
	defer cancelFunc()

	var datasetName = stageCmd.Dataset.Positional.DatasetPath

	srl := rescom.StagingResourceLocation{}
	srl.ContainerType = jsonld.DatasetResource
	srl.ContainerIRI = "" // n/a  for type
	srl.ObjectType = jsonld.ContextResource
	srl.ObjectIRI = "" // n/a for type
	srl.DatasetPath = datasetName

	return stageResourceLocation(ctxt, srl)

}

// command entrypoint for staging a named graph within a dataset and all
// the resources therein
func (c *dgrzStageNamedGraph) Execute(args []string) error {

	ctxt, cancelFunc := context.WithCancel(getCmdContext())
	defer cancelFunc()

	datasetName := stageCmd.Dataset.Positional.DatasetPath
	iri := stageCmd.Dataset.NamedGraph.Positional.IRI
	parentIRI := stageCmd.Dataset.NamedGraph.ParentGraphIRI

	srl := rescom.StagingResourceLocation{}
	if len(parentIRI) > 0 {
		srl.ContainerType = jsonld.NamedGraphResource
		srl.ContainerIRI = parentIRI
	} else {
		srl.ContainerType = jsonld.DatasetResource
		srl.ContainerIRI = "" // n/a for Datset type
	}

	srl.ObjectType = jsonld.NamedGraphResource
	srl.ObjectIRI = iri
	srl.DatasetPath = datasetName

	return stageResourceLocation(ctxt, srl)

}

// command entrypoint for staging a single node within a dataset
func (c *dgrzStageNode) Execute(args []string) error {

	ctxt, cancelFunc := context.WithCancel(getCmdContext())
	defer cancelFunc()

	datasetName := stageCmd.Dataset.Positional.DatasetPath
	iri := stageCmd.Dataset.Node.Positional.IRI
	parentIRI := stageCmd.Dataset.Node.ParentGraphIRI

	srl := rescom.StagingResourceLocation{}
	if len(parentIRI) > 0 {
		// NODE IS IN THE DEFAULT GRAPH
		srl.ContainerType = jsonld.NamedGraphResource
		srl.ContainerIRI = parentIRI
	} else {
		// NODE IS WITHIN A NAMED GRAPH
		srl.ContainerType = jsonld.DatasetResource
		srl.ContainerIRI = "" // n/a for Datset type
	}

	srl.ObjectType = jsonld.NodeResource
	srl.ObjectIRI = iri
	srl.DatasetPath = datasetName

	return stageResourceLocation(ctxt, srl)

}

func stageResourceLocation(ctxt context.Context, srl rescom.StagingResourceLocation) error {

	stager, err := resource.GetGrapplicationResourceStager(ctxt, stageCmd.Grapplication)
	if err != nil {
		return err
	}
	defer stager.Close(ctxt)

	// Stage outermost context
	if err := stager.Add(ctxt, srl); err != nil {
		if err2 := stager.Rollback(ctxt); err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
		return err

	}

	return stager.Commit(ctxt)

}

///////////

func (o *dgrzStageCmd) CommandName() string {
	return "stage"
}

func (o *dgrzStageCmd) ShortDescription() string {
	return "stage a grapplication resource"
}

func (o *dgrzStageCmd) LongDescription() string {
	return "stage a new JSON-LD grapplication resource from your workspace"
}
