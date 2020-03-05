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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jessevdk/go-flags"
)

type Dogg3rzCmd interface {
	CommandName() string
	ShortDescription() string
	LongDescription() string
}

type dgrzAddCmd struct {
	//xFileName string ` value-name:"FILE"`
}

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

func init() {
	/*
		_, _ = parser.AddCommand("init",
			"Init the node repo and exit",
			"Initialize the node repository and exit.",
			&dgrzInitCmd{})

	*/

	_, _ = parser.AddCommand("add",
		"add file to repo and exit",
		"Add file to repo and exit",
		&dgrzAddCmd{})
}

func register(cmd Dogg3rzCmd) {
	// APPEND COMMAND TO LIST OF ALRAEDY REGISTERED DOGG3RZ COMMANDS
	dgrzCmds = append(dgrzCmds, cmd)
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

func (x *dgrzAddCmd) Execute(args []string) error {

	if len(args) < 1 {
		return errMissingFilePath
	}

	filePath := args[0]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	filedata, err := ioutil.ReadFile(args[0])

	if !json.Valid(filedata) {

		return fmt.Errorf("Not a json file: %s", args[0])
	}

	sh := shell.NewShell("localhost:5001")

	//fmt.Printf("is up = %v\n", sh)

	cid, err := sh.DagPut(string(filedata), "json", "cbor")
	if err != nil {
		return err
	}

	fmt.Printf("Added %s for file %s", cid, args[0])

	return nil
}
