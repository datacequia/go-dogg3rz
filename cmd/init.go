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
	//config "github.com/ipfs/go-ipfs-config"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	//	"os"

	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/datacequia/go-dogg3rz/resource/config"
	//"github.com/datacequia/go-dogg3rz/util"
)

type dgrzInitCmd struct {
	Node dgrzInitNode `command:"node" description:"initialize this host as a dogg3rz node" `
	Repo dgrzInitRepo `command:"repository" alias:"repo" description:"initialize a new repository" `
}

type dgrzInitNode struct {
	UserEmail       string `long:"user-email" description:"user's email address" required:"true"`
	UserFirstName   string `long:"user-firstname" description:"user's first name"`
	UserLastName    string `long:"user-lastname" description:"user's last name"`
	IPFSApiEndpoint string `long:"ipfs-api-endpoint" description:"the http(s) url your IPFS node's api endpoint listener" default:"http://localhost:5001/"`
}

type dgrzInitRepo struct {
	Positional struct {
		RepoName string `positional-arg-name:"REPO_NAME" required:"yes" `
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

	return resource.GetNodeResource(getCmdContext()).InitNode(c)

}

func (x *dgrzInitRepo) Execute(args []string) error {

	return resource.GetRepositoryResource(getCmdContext()).InitRepo(x.Positional.RepoName)

}

// // IMPLEMENTS 'Commander' interface
func (x *dgrzInitCmd) Execute(args []string) error {

	fmt.Printf("Repo path is '%s'\n", "d")

	return nil
}

func (o *dgrzInitCmd) CommandName() string {
	return "init"
}

func (o *dgrzInitCmd) ShortDescription() string {
	return "initialize a new repository"
}

func (o *dgrzInitCmd) LongDescription() string {
	return "initialize a new repository"
}
