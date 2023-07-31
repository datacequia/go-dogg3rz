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
	//"github.com/datacequia/go-dogg3rz/util"
)

type dgrzCreateNamedGraph struct {
	Positional struct {
		GraphID string `positional-arg-name:"GRAPH_NAME" description:"the name/ID of the graph" required:"true"`
	} `positional-args:"yes"`

	Dataset       string `long:"dataset"   description:"Dataset path" required:"true" `
	ParentGraphID string `long:"parent"  description:"Parent graph ID" required:"false" default:"default"`
}

///////////////////////////////////////////////////////////
// CREATE CLASS SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////

func (x *dgrzCreateNamedGraph) Execute(args []string) error {

	ctxt := getCmdContext()
	grapp := resource.GetGrapplicationResource(ctxt)

	err := grapp.CreateNamedGraph(ctxt,
		createCmd.Grapplication,
		x.Dataset,
		x.Positional.GraphID,
		x.ParentGraphID,
	)

	return err
}

func (x *dgrzCreateNamedGraph) CommandName() string {
	return "create namedgraph"
}

func (x *dgrzCreateNamedGraph) ShortDescription() string {
	return "create named graph in dataset at path"
}

func (x *dgrzCreateNamedGraph) LongDescription() string {
	return "create named graph in dataset at path"
}
