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
	"github.com/datacequia/go-dogg3rz/resource"
	//	"github.com/datacequia/go-dogg3rz/util"
)

//"github.com/datacequia/go-dogg3rz/resource"

type dgrzInsertNode struct {
	Options struct {
		Type string `long:"type" short:"t" description:"RDF Schema Type (as an IRI or term)" required:"false"`
		ID   string `long:"id" short:"i" description:"Node ID (as an IRI or Term). If not supplied, will be a blank node" required:"false"`
	}

	Into dgrzSubCmdInto `command:"into" description:""`
}

type dgrzSubCmdInto struct {
	NamedGraph struct {
		Positional struct {
			GraphName string `positional-arg-name:"GRAPH_NAME" description:"the name of the graph" required:"true"`
		} `positional-args:"yes"`

		Dataset dgrzSubCmdDataset `command:"dataset" alias:"ds" `
	} `command:"named-graph" alias:"ng"`

	DefaultGraph struct {
		Dataset dgrzSubCmdDataset `command:"dataset" alias:"ds" `
	} `command:"default-graph" alias:"dg"`
}

type dgrzSubCmdDataset struct {
	Positional struct {
		Path string `positional-arg-name:"DATASET_PATH" description:"the repository path to the dataset" required:"true"`
	} `positional-args:"yes"`
	PropertyValues dgrzSubCmdNodePropertyValues `command:"property-values" alias:"pv" description:"describes the property names/values"`
}

type dgrzSubCmdNodePropertyValues struct {
	Properties []string `long:"property" short:"p" value-name:"IRI" description:"node property as an IRI or Term" required:"true"`
	Values     []string `long:"value" short:"v" value-name:"IRI" description:"node value (for corresponding property) as an IRI or Literal" required:"true"`
}

func (cmd *dgrzSubCmdNodePropertyValues) Execute(args []string) error {

	var err error
	insertDefaultGraph := &insertCmd.Node.Into.DefaultGraph.Dataset.PropertyValues
	insertNamedGraph := &insertCmd.Node.Into.NamedGraph.Dataset.PropertyValues

	repo := resource.GetRepositoryResource(getCmdContext())

	switch cmd {
	case insertDefaultGraph:
		//fmt.Printf("insert default graph")
		err = repo.InsertNode(insertCmd.Repository,
			insertCmd.Node.Into.DefaultGraph.Dataset.Positional.Path,
			insertCmd.Node.Options.Type,
			insertCmd.Node.Options.ID,
			"", // GRAPH NAME
			insertDefaultGraph.Properties,
			insertDefaultGraph.Values)

	case insertNamedGraph:
		//fmt.Printf("insert named graph ")
		err = repo.InsertNode(insertCmd.Repository,
			insertCmd.Node.Into.NamedGraph.Dataset.Positional.Path, //.DefaultGraph.Dataset.Positional.Path,
			insertCmd.Node.Options.Type,
			insertCmd.Node.Options.ID,
			insertCmd.Node.Into.NamedGraph.Positional.GraphName, // GRAPH NAME
			insertDefaultGraph.Properties,
			insertDefaultGraph.Values)

	default:
		panic("dgrzSubCmdNodePropertyValues.Execute: unhandled sub-command")
	}
	return err
}

///////////////////////////////////////////////////////////
// INSERT NODE  FUNCTIONS
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
// INSERT NODE SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
// CREATE TYPE DATATYPE SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////
