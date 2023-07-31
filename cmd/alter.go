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

type dgrzAlterCmd struct {
	Grapplication string `long:"grapp" short:"r" env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`

	Context dgrzAlterContext `command:"context" alias:"ctxt" description:"alter a JSON-LD object's context"`
}

///////////////////////////////////////////////////////////
// ALTER COMMAND FUNCTIONS
///////////////////////////////////////////////////////////

var alterCmd = dgrzAlterCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&alterCmd)
}

func (o *dgrzAlterCmd) CommandName() string {
	return "alter"
}

func (o *dgrzAlterCmd) ShortDescription() string {
	return "alter a JSON-LD schema/non-data resource"
}

func (o *dgrzAlterCmd) LongDescription() string {
	return "alter a JSON-LD schema/non-data resource"
}
