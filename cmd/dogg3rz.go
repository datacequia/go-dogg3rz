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
	"context"
	"errors"
	"log"
	//	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/env"
	"github.com/jessevdk/go-flags"
)

type Dogg3rzCmd interface {
	CommandName() string
	ShortDescription() string
	LongDescription() string
}

/*
type dgrzAddCmd struct {
	//xFileName string ` value-name:"FILE"`
}
*/

type options struct {
}

// HOLDS LIST OF COMMANDS REGISTERED BY VARIOUS COMMAND SRC FILE
var dgrzCmds []Dogg3rzCmd

// HOLDS LIST OF INITIALIZATION FUNCS FOR EACH IMPLEMENTED
// STATE-STORE SUPPORTED
// KEY = STORE TYPE
// VALLUE = LIST OF INIT FUNCS FOR STORE TYPE
var stateStoreRegistry = make(map[string][]func() error)

var errMissingFilePath = errors.New("missing file path")
var parser = flags.NewParser(&options{}, flags.Default)

func register(cmd Dogg3rzCmd) {
	// APPEND COMMAND TO LIST OF ALRAEDY REGISTERED DOGG3RZ COMMANDS
	dgrzCmds = append(dgrzCmds, cmd)
}

func getCmdContext() context.Context {

	// INIT A NEW CONTEXT BASED ON DOGG3RZ ENV. VARS
	// USING A (ROOT) BACKGROUND CONTEXT AS PARENT CONTEXT
	return env.InitContextFromEnv(context.Background())

}

func Run() {

	// ADD ALL REGISTERED COMMANDS TO COMMAND PARSER
	for _, c := range dgrzCmds {
		_, err := parser.AddCommand(c.CommandName(), c.ShortDescription(), c.LongDescription(), c)
		if err != nil {
			log.Fatalf("failed to add command to command parser: { Command Name = '%s'}: %s", c.CommandName(), err)
		}
	}

	// PARSE COMMAND LINE ARGS
	_, err := parser.Parse()
	if err != nil {
		//log.Fatalf("failed to parse command line arguments: %s", err)
	}
}
