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
	"fmt"
	"context"
	"github.com/datacequia/go-dogg3rz/resource"
)


// STAGE  DEFAULT-GRAPH | GRAPH IRI | NODE IRI
type dgrzStageCmd struct {
	Repository string `long:"repo" short:"r" env:"DOGG3RZ_REPO" description:"repository name" required:"true"`

	All     dgrzStageAllCmd     `command:"all" description:"stage all resources in a repository"`
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
}

type dgrzStageNode struct {
	Positional struct {
		IRI string `positional-arg-name:"IRI" required:"yes" `
	} `positional-args:"yes"`
}

var stageCmd = dgrzStageCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&stageCmd)
}

//////////////////////////////////////
// EXECUTE ENTRYPOINTS
/////////////////////////////////////

func (cmd *dgrzStageAllCmd) Execute(args []string) error {

	ctxt,cancelFunc := context.WithCancel(getCmdContext()) 

	defer cancelFunc()


	stager := resource.GetRepositoryStagerResource(ctxt,stageCmd.Repository)

	fmt.Println("stager",stager)

/*
	ctxt := getCmdContext()

	repo := resource.GetRepositoryResource(ctxt)

	var err error
	var stageAllCmd *dgrzStageAllCmd = &stageCmd.All
	var stageDatasetAllCmd *dgrzStageAllCmd = &stageCmd.Dataset.All

	// DIFFERENTIATE THE CONTEXT IN WHICH THE 'ALL' COMMAND WAS USED
	// BY COMPARING IT'S (RECEIVER) POINTER VALUE TO THE TWO CONTEXTS IN WHICH IT
	// IS KNOWN TO BE USED
	switch cmd {

	case stageAllCmd:

	case stageDatasetAllCmd:

		var stagingList = []common.StagingResourceLocation{
			common.StagingResourceLocation{
				JSONLDDocumentLocation: common.JSONLDDocumentLocation{
					ObjectType:    jsonld.DatasetResource,
					ObjectIRI:     "",
					ContainerType: jsonld.DatasetResource,
					ContainerIRI:  "",
				},
				DatasetPath: stageCmd.Dataset.Positional.DatasetPath,
			},
		}

		var srl []common.StagingResource
		if srl, err = repo.StageResources(ctxt, stageCmd.Repository, stagingList); err != nil {
			return err
		}

		for _, sr := range srl {

			fmt.Printf("staged %s\n", sr)
		}

	default:
		panic(fmt.Sprintf("unhandled context in which 'all' sub-command is used"))
	}
*/

	return nil
}

// function dgrzStageContext.Execute is the entrypoint function to
// stage the outermost context of a JSON-LD dataset
func (c *dgrzStageContext) Execute(args []string) error {

	return nil

}

func (c *dgrzStageNamedGraph) Execute(args []string) error {

	return nil
}

func (c *dgrzStageNode) Execute(args []string) error {

	return nil
}

///////////

func (o *dgrzStageCmd) CommandName() string {
	return "stage"
}

func (o *dgrzStageCmd) ShortDescription() string {
	return "stage a repository resource"
}

func (o *dgrzStageCmd) LongDescription() string {
	return "stage a new JSON-LD repository resource from your working tree"
}
