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
	"github.com/datacequia/go-dogg3rz/impl/file"

	"os"

	"github.com/datacequia/go-dogg3rz/resource"
)

// Command to get list of grapplications
type dgrzGetGrappCmd struct {
}

/*
func init() {
	// REGISTER THE 'get grapp ' COMMAND
	register(&dgrzGetGrappCmd{})
}
*/

func (o *dgrzGetGrappCmd) Execute(args []string) error {

	ctxt := getCmdContext()
	node := resource.GetNodeResource(ctxt)
	var files []string
	var err error
	files, err = node.GetGrapps(ctxt)
	if err != nil {
		return err
	}
	printValues(files, file.DgrzDirName, os.Stdout)
	return err

}

func (o *dgrzGetGrappCmd) CommandName() string {
	return "get grapplication"
}

func (o *dgrzGetGrappCmd) ShortDescription() string {
	return "get listing of grapplications"
}

func (o *dgrzGetGrappCmd) LongDescription() string {
	return "get listing of grapplications"
}
