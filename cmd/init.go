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
	"os"

	"github.com/datacequia/go-dogg3rz/resource"
)

type dgrzInitCmd struct {
	Positional struct {
		GrapplicationDirectory string `positional-arg-name:"DIRECTORY" description:"grapplication project directory path"  `
	} `positional-args:"yes"`
}

///////////////////////////////////////////////////////////
// CREATE COMMAND FUNCTIONS
///////////////////////////////////////////////////////////

var initCmd = dgrzInitCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&initCmd)
}

func (o *dgrzInitCmd) CommandName() string {
	return "init"
}

func (o *dgrzInitCmd) ShortDescription() string {
	return "initialize directory as grapplication project"
}

func (o *dgrzInitCmd) LongDescription() string {
	return "initialize directory as grapplication project"
}

func (x *dgrzInitCmd) Execute(args []string) error {

	ctxt := getCmdContext()

	grapp := resource.GetGrapplicationResource(ctxt)

	// DIRECTORY PATH POSITIONAL ARGUMENT NOT SUPPLIED
	//
	if len(x.Positional.GrapplicationDirectory) < 1 {
		if d, err := os.Getwd(); err != nil {
			return err
		} else {
			x.Positional.GrapplicationDirectory = d
		}
	}
	//	fmt.Printf("INIT: { grapplication dir = %s }\n", x.Positional.GrapplicationDirectory)

	return grapp.Init(ctxt, x.Positional.GrapplicationDirectory) //grapp.Add(ctxt, x.Grapplication, x.Positional.Path)

}
