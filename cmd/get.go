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

//	"io"
//	"os"
//"strings"

//	"github.com/datacequia/go-dogg3rz/errors"

//	"github.com/xeipuuv/gojsonschema"

type dgrzGetCmd struct {
	Dataset dgrzGetDatasetCmd `command:"dataset" alias:"ds" description:"get datasets in a grapp" `

	Grapps dgrzGetGrappCmd `command:"grapplication" alias:"grapp" description:"get list of grapps" `
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzGetCmd{})
}

func (o *dgrzGetCmd) CommandName() string {
	return "get"
}

func (o *dgrzGetCmd) ShortDescription() string {
	return "get listing of grapplication resources"
}

func (o *dgrzGetCmd) LongDescription() string {
	return "get listing of grapplication resources"
}
