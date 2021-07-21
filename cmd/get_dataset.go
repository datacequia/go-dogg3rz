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
	//"github.com/datacequia/go-dogg3rz/impl/file/repo"
	"github.com/datacequia/go-dogg3rz/resource"
  "os"
)

// Command to get list of dataset in a repo
type dgrzGetDatasetCmd struct {

	Repository string `long:"repo" short:"r"  env:"DOGG3RZ_REPO" description:"repository name" required:"true"`


}

/*
func init() {
	// REGISTER THE 'get datatset ' COMMAND
	register(&dgrzGetDatasetCmd{})
}
*/


func (x *dgrzGetDatasetCmd) Execute(args []string) error {


	ctxt := getCmdContext()
	repo := resource.GetRepositoryResource(ctxt)
	var files []string
	var err error
	files, err = repo.GetDataSets(ctxt, x.Repository)
	if err != nil {
     return err
    }
    printValues(files, file.DgrzDirName, os.Stdout)
	return err

}



func (x *dgrzGetDatasetCmd) CommandName() string {
	return "get dataset"
}

func (x *dgrzGetDatasetCmd) ShortDescription() string {
	return "get listing of datasets in a repository"
}

func (x *dgrzGetDatasetCmd) LongDescription() string {
	return "get listing of datasets in a repository"
}
