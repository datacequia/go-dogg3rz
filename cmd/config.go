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

type dgrzConfigCmd struct {
	Init dgrzConfigInitCmd `command:"init" description:"initialize the user environment configuration" `
	//Grapp dgrzInitGrapp `command:"grapplication" alias:"grapp" description:"initialize a new grapplication" `
}

type dgrzConfigInitCmd struct {
	ActivityPubUserHandle string `long:"activitypub-handle" description:"user's ActivityPub handle (i.e. @<user>@<host>)" required:"true"`
	IPFSApiEndpoint       string `long:"ipfs-api-endpoint" description:"the http(s) url your IPFS node's api endpoint listener" default:"http://localhost:5001/"`
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzConfigCmd{})
}

func (x *dgrzConfigInitCmd) Execute(args []string) error {

	var c config.Dogg3rzConfig

	// ASSIGN REQUIIRED USER EMAIL
	c.User.ActivityPubUserHandle = x.ActivityPubUserHandle

	ctxt := getCmdContext()

	// INITIALIZE USER ENVIRONMENT
	if err := resource.GetConfigResource(ctxt).InitConfig(ctxt, c); err != nil {
		return err
	}

	// PULL IPFS IMAGE
	if err := ipfs.PullDefault(); err != nil {
		return err
	}

	return nil
}

// CONFIG CMD
func (o *dgrzConfigCmd) CommandName() string {
	return "config"
}

func (o *dgrzConfigCmd) ShortDescription() string {
	return "user configuration commands"
}

func (o *dgrzConfigCmd) LongDescription() string {
	return "user configuration commands"
}

// CONFIG INIT COMMAND
func (o *dgrzConfigInitCmd) CommandName() string {
	return "init"
}

func (o *dgrzConfigInitCmd) ShortDescription() string {
	return "initialize the user environment"
}

func (o *dgrzConfigInitCmd) LongDescription() string {
	return "initialize the user environment for use by this application"
}
