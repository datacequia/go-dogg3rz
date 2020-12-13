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
	"github.com/datacequia/go-dogg3rz/resource"
	//"github.com/datacequia/go-dogg3rz/util"
)

type dgrzCreateSnapshot struct {
	//	Repository string `long:"repo" short:"r" env:"DOGG3RZ_REPO" description:"repository name" required:"true"`
}

///////////////////////////////////////////////////////////
// CREATE CLASS SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
// CREATE SNAPSHOT SUBCOMMAND  FUNCTIONS
///////////////////////////////////////////////////////////

func (x *dgrzCreateSnapshot) Execute(args []string) error {

	//	fmt.Printf("hello snapshot: { repo = %s }\n", x.Positional.Repository)

	repo := resource.GetRepositoryResource(getCmdContext())

	return repo.CreateSnapshot(createCmd.Repository)

}

func (o *dgrzCreateSnapshot) CommandName() string {
	return "create snapshot"
}

func (o *dgrzCreateSnapshot) ShortDescription() string {
	return "create a repository snapshot"
}

func (o *dgrzCreateSnapshot) LongDescription() string {
	return "create a repository snapshot"
}
