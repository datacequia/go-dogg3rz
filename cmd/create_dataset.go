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

/*

Excerpts on RDF Datasets from  https://www.w3.org/TR/rdf11-concepts/#data-model

*  A new concept in RDF 1.1 is the notion of an RDF dataset to represent
   multiple graphs."

*  RDF datasets are used to organize collections of RDF graphs, and comprise a
   default graph and zero or more named graphs"

*  RDF datasets are used to organize collections of RDF graphs, and comprise a
   default graph and zero or more named graphs"

 * There are many possible uses for RDF datasets. One such use is to hold snapshots
   of multiple RDF sources."

 * An RDF dataset is a collection of RDF graphs, and comprises:

    * Exactly one default graph, being an RDF graph.
	    The default graph does not have a name and may be empty.
    * Zero or more named graphs. Each named graph is a pair consisting of an IRI
	    or a blank node (the graph name), and an RDF graph. Graph names are unique
	  	within an RDF dataset.

*/

package cmd

import (
	"github.com/datacequia/go-dogg3rz/resource"
	//"github.com/datacequia/go-dogg3rz/util"
)

type dgrzCreateDataset struct {
	//	Repository string `long:"repo" short:"r" env:"DOGG3RZ_REPO" description:"repository name" required:"true"`

	Positional struct {
		DatasetPath string `positional-arg-name:"DATASET_PATH" description:"repository path to dataset" required:"yes" `
	} `positional-args:"yes"`
	//	SubclassOf string `long:"subclass-of" description:"RDF subclass" required:"false"`

	// USED FOR SCHEMA PARSING

}

///////////////////////////////////////////////////////////
// CREATE CLASS SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////

func (x *dgrzCreateDataset) Execute(args []string) error {

	ctxt := getCmdContext()
	repo := resource.GetRepositoryResource(ctxt)

	//fmt.Println("create dataset", createCmd.Repository, x.Positional.DatasetPath)

	if err := repo.CreateDataset(ctxt,
		createCmd.Repository,
		x.Positional.DatasetPath); err != nil {
		return err
	}
	//	fmt.Println("success")

	return nil
}

func (o *dgrzCreateDataset) CommandName() string {
	return "create dataset"
}

func (o *dgrzCreateDataset) ShortDescription() string {
	return "create RDF dataset at path"
}

func (o *dgrzCreateDataset) LongDescription() string {
	return "create RDF dataset at path"
}
