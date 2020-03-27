/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"fmt"
	//config "github.com/ipfs/go-ipfs-config"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	//	"os"

	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/datacequia/go-dogg3rz/resource/config"
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

	return resource.GetNodeResource().InitNode(c)

}

func (x *dgrzInitRepo) Execute(args []string) error {

	return resource.GetRepositoryResource().InitRepo(x.Positional.RepoName)

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
