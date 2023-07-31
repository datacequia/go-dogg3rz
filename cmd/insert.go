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

type dgrzInsertCmd struct {
	Grapplication string `long:"grapp" short:"r" env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`

	//	DirPath      string                 `long:"dirpath" description:"directory path" required:"true"`
	/*
		Dataset   dgrzCreateDataset   `command:"dataset" alias:"ds" description:"create a new dataset" `
		Namespace dgrzCreateNamespace `command:"namespace" alias:"ns" description:"create a new namespace (IRI) in a grapplication directory" `

		Snapshot dgrzCreateSnapshot `command:"snapshot" alias:"ss" description:"create a snapshot of the grapplication"`

		Type dgrzCreateType `command:"type" description:"create an instance of an RDF [Schema] type"`
	*/
	Node dgrzInsertNode `command:"node" description:"insert a JSON-LD Node into a grapplication dataset graph"`
}

///////////////////////////////////////////////////////////
// CREATE COMMAND FUNCTIONS
///////////////////////////////////////////////////////////

var insertCmd = dgrzInsertCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&insertCmd)
}

func (o *dgrzInsertCmd) CommandName() string {
	return "insert"
}

func (o *dgrzInsertCmd) ShortDescription() string {
	return "insert a new data resource into a grapplication"
}

func (o *dgrzInsertCmd) LongDescription() string {
	return "insert a new data resource into a grapplication"
}
