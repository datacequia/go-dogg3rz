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

	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

type dgrzInitCmd struct {
	Node  dgrzInitNode  `command:"node" description:"initialize this host as a dogg3rz node" `
	Grapp dgrzInitGrapp `command:"grapplication" alias:"grapp" description:"initialize a new grapplication" `
}

type dgrzInitNode struct {
	UserEmail       string `long:"user-email" description:"user's email address" required:"true"`
	UserFirstName   string `long:"user-firstname" description:"user's first name"`
	UserLastName    string `long:"user-lastname" description:"user's last name"`
	IPFSApiEndpoint string `long:"ipfs-api-endpoint" description:"the http(s) url your IPFS node's api endpoint listener" default:"http://localhost:5001/"`
}

type dgrzInitGrapp struct {
	Positional struct {
		GrappName string `positional-arg-name:"GRAPP_NAME" required:"yes" `
	} `positional-args:"yes"`
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzInitCmd{})
}

func (x *dgrzInitNode) Execute(args []string) error {

	var c config.Dogg3rzConfig

	// ASSIGN REQUIIRED USER EMAIL
	c.User.Email = x.UserEmail
	c.User.FirstName = x.UserFirstName
	c.User.LastName = x.UserLastName

	ctxt := getCmdContext()

	return resource.GetNodeResource(ctxt).InitNode(ctxt, c)

}

func (x *dgrzInitGrapp) Execute(args []string) error {

	ctxt := getCmdContext()

	return resource.GetGrapplicationResource(ctxt).InitGrapp(ctxt, x.Positional.GrappName)

}

// // IMPLEMENTS 'Commander' interface
func (x *dgrzInitCmd) Execute(args []string) error {

	fmt.Printf("Grapp path is '%s'\n", "d")

	return nil
}

func (o *dgrzInitCmd) CommandName() string {
	return "init"
}

func (o *dgrzInitCmd) ShortDescription() string {
	return "initialize a new grapplication"
}

func (o *dgrzInitCmd) LongDescription() string {
	return "initialize a new grapplication"
}
