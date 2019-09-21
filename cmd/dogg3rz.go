package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jessevdk/go-flags"
)

type dgrzInitCmd struct {
	RepoPath string `short:"r" long:"repo-dir" description:"Specify a custom repository path."`
}

type dgrzAddCmd struct {
	//xFileName string ` value-name:"FILE"`
}

type options struct {
}

var errMissingFilePath = errors.New("missing file path")

var parser = flags.NewParser(&options{}, flags.Default)

func init() {
	_, _ = parser.AddCommand("init",
		"Init the node repo and exit",
		"Initialize the node repository and exit.",
		&dgrzInitCmd{})
	_, _ = parser.AddCommand("add",
		"add file to repo and exit",
		"Add file to repo and exit",
		&dgrzAddCmd{})
}

func Run() {

	_, _ = parser.Parse()
}

func (x *dgrzInitCmd) Execute(args []string) error {

	return dgrzInit(x)
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
