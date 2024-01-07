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
	"github.com/datacequia/go-dogg3rz/ipfs"
	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

/*
type dgrzInitCmd struct {
	Node  dgrzInitNode  `command:"node" description:"initialize this host as a dogg3rz node" `
	//Grapp dgrzInitGrapp `command:"grapplication" alias:"grapp" description:"initialize a new grapplication" `
}
*/

type dgrzInitCmd struct {
	ActivityPubUserHandle string `long:"activitypub-handle" description:"user's ActivityPub handle (i.e. @<user>@<host>)" required:"true"`
	IPFSApiEndpoint       string `long:"ipfs-api-endpoint" description:"the http(s) url your IPFS node's api endpoint listener" default:"http://localhost:5001/"`
}

/*
type dgrzInitGrapp struct {
	Positional struct {
		GrappName string `positional-arg-name:"GRAPP_NAME" required:"yes" `
	} `positional-args:"yes"`
}
*/

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzInitCmd{})
}

func (x *dgrzInitCmd) Execute(args []string) error {

	var c config.Dogg3rzConfig

	// ASSIGN REQUIIRED USER EMAIL
	c.User.ActivityPubUserHandle = x.ActivityPubUserHandle

	ctxt := getCmdContext()

	// INITIALIZE USER ENVIRONMENT
	if err := resource.GetNodeResource(ctxt).InitNode(ctxt, c); err != nil {
		return err
	}

	// PULL IPFS IMAGE
	if err := ipfs.PullDefault(); err != nil {
		return err
	}

	return nil
}

/*
func (x *dgrzInitGrapp) Execute(args []string) error {

	ctxt := getCmdContext()

	return resource.GetGrapplicationResource(ctxt).InitGrapp(ctxt, x.Positional.GrappName)

}
*/

// // IMPLEMENTS 'Commander' interface
/*
func (x *dgrzInitCmd) Execute(args []string) error {

	fmt.Printf("Grapp path is '%s'\n", "d")

	return nil
}
*/

func (o *dgrzInitCmd) CommandName() string {
	return "init-env"
}

func (o *dgrzInitCmd) ShortDescription() string {
	return "initialize the user environment"
}

func (o *dgrzInitCmd) LongDescription() string {
	return "initialize the user environment for use by this application"
}
