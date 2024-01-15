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
	"io"
	"os"

	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource"
)

type dgrzValidateCmd struct {
	//Init dgrzConfigInitCmd `command:"init" description:"initialize the user environment configuration" `
	//Grapp dgrzInitGrapp `command:"grapplication" alias:"grapp" description:"initialize a new grapplication" `
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose validate information"`
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzValidateCmd{})
}

func (x *dgrzValidateCmd) Execute(args []string) error {

	// INITIALIZE USER ENVIRONMENT
	ctxt := getCmdContext()

	//fmt.Println("file vaalidate ", grappDir)
	grappDir, err := file.GrapplicationDirPath(ctxt)
	if err != nil {
		return err
	}
	var verboseWriter io.Writer

	if len(x.Verbose) > 0 && x.Verbose[0] {
		verboseWriter = os.Stdout
		//fmt.Println("chose verbose option", len(x.Verbose), x.Verbose[0])

	}

	if err := resource.GetGrapplicationResource(ctxt).Validate(ctxt, grappDir, verboseWriter); err != nil {
		return err
	}

	return nil
}

func (o *dgrzValidateCmd) CommandName() string {
	return "validate"
}

func (o *dgrzValidateCmd) ShortDescription() string {
	return "validate grapplication project files"
}

func (o *dgrzValidateCmd) LongDescription() string {
	return "validate grapplication project files"
}
