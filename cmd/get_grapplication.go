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

// Command to get list of dataset in a grapp
type dgrzGetDatasetCmd struct {
	Grapplication string `long:"grapplication" short:"g"  env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`
}

/*
func init() {
	// REGISTER THE 'get datatset ' COMMAND
	register(&dgrzGetDatasetCmd{})
}
*/

func (x *dgrzGetDatasetCmd) Execute(args []string) error {
	/*
		ctxt := getCmdContext()
		grapp := resource.GetGrapplicationResource(ctxt)
		var files []string
		var err error

			files, err = grapp.GetDataSets(ctxt, x.Grapplication)
			if err != nil {
				return err
			}

		printValues(files, file.DgrzDirName, os.Stdout)
		return err
	*/
	return nil

}

func (x *dgrzGetDatasetCmd) CommandName() string {
	return "get dataset"
}

func (x *dgrzGetDatasetCmd) ShortDescription() string {
	return "get listing of datasets in a grapplication"
}

func (x *dgrzGetDatasetCmd) LongDescription() string {
	return "get listing of datasets in a grapplication"
}
