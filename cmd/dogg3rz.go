package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	//fmt.Printf("you typed filename '%s'\n", args[0])
	/*
		r, err := os.Open(args[0])
		if err != nil {

			return err
		}
	*/
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
