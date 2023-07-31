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

type dgrzCreateCmd struct {
	Grapplication string `long:"grapplication" short:"g" env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`

	//	DirPath      string                 `long:"dirpath" description:"directory path" required:"true"`
	Dataset dgrzCreateDataset `command:"dataset" alias:"ds" description:"create a new dataset" `
	//	Namespace dgrzCreateNamespace `command:"namespace" alias:"ns" description:"create a new namespace (IRI) in a grapplication directory" `

	Snapshot dgrzCreateSnapshot `command:"snapshot" alias:"ss" description:"create a snapshot of the grapplication"`

	Type dgrzCreateType `command:"type" description:"create an instance of an RDF [Schema] type"`

	NamedGraph dgrzCreateNamedGraph `command:"namedgraph" alias:"ng" description:"create a new named graph in the dataset"`
}

///////////////////////////////////////////////////////////
// CREATE COMMAND FUNCTIONS
///////////////////////////////////////////////////////////

var createCmd = dgrzCreateCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&createCmd)
}

func (o *dgrzCreateCmd) CommandName() string {
	return "create"
}

func (o *dgrzCreateCmd) ShortDescription() string {
	return "create a new schema/non-data grapplication resource"
}

func (o *dgrzCreateCmd) LongDescription() string {
	return "create a new schema/non-data grapplication resource"
}
