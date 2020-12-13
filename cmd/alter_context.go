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

See: https://www.w3.org/TR/json-ld/#the-context

*/

package cmd

import (
	"github.com/datacequia/go-dogg3rz/resource"
)

type dgrzAlterContext struct {
	//	Repository string `long:"repo" short:"r" env:"DOGG3RZ_REPO" description:"repository name" required:"true"`

	//subcommands-optional:"true"
	Dataset dgrzAlterContextDataset `command:"dataset" alias:"ds" description:"the dataset object in which to alter the context"`

	// USED FOR SCHEMA PARSING

}

type dgrzAlterContextDataset struct {
	Positional struct {
		DatasetPath string `positional-arg-name:"DATASET_PATH" required:"yes" `
	} `positional-args:"yes"`

	Node dgrzAlterContextNode `command:"node" description:"the node object in which to alter the alter the context"`

	Add dgrzAlterContextAdd `command:"add" description:"add resource item to dataset"`
}

type dgrzAlterContextNode struct {
	// positional ID
	Positional struct {
		NodeID string `positional-arg-name:"NODE_ID" required:"yes" `
	} `positional-args:"yes"`

	Add dgrzAlterContextAdd `command:"add" description:"add resource item to node"`
}

type dgrzAlterContextAdd struct {
	Namespace dgrzAddNamespace `command:"namespace" alias:"ns" description:"the dataset object in which to alter the context"`
}

type dgrzAddNamespace struct {
	Positional struct {
		Term string `positional-arg-name:"TERM"  required:"yes" `
		IRI  string `postional-arg-name:"IRI" required:"yes"`
		//	Uri    string `positional-arg-name:"URI" required:"yes" `
	} `positional-args:"yes"`
}

func (o *dgrzAddNamespace) Execute(args []string) error {

	var err error

	ctxt := getCmdContext()

	repo := resource.GetRepositoryResource(ctxt)

	if len(alterCmd.Context.Dataset.Node.Positional.NodeID) > 0 {
		// alter context dataset <DATASET> node <NODEID> WAS CALLED
		// BECAUSE ABOVE WAS INITIALIZED
		//	AddNamespaceNode(repoName string, datasetPath string, nodeID string, term string, iri string) error

		err = repo.AddNamespaceNode(ctxt,
			alterCmd.Repository,
			alterCmd.Context.Dataset.Positional.DatasetPath,
			alterCmd.Context.Dataset.Node.Positional.NodeID,
			alterCmd.Context.Dataset.Node.Add.Namespace.Positional.Term,
			alterCmd.Context.Dataset.Node.Add.Namespace.Positional.IRI)

	} else {

		err = repo.AddNamespaceDataset(alterCmd.Repository,
			alterCmd.Context.Dataset.Positional.DatasetPath,
			alterCmd.Context.Dataset.Add.Namespace.Positional.Term,
			alterCmd.Context.Dataset.Add.Namespace.Positional.IRI)

	}
	return err
}
