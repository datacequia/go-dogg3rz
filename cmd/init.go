package cmd

import (
	"fmt"
	//config "github.com/ipfs/go-ipfs-config"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	//	"os"

	implnode "github.com/datacequia/go-dogg3rz/impl/file/node"

	resnode "github.com/datacequia/go-dogg3rz/resource/node"
)

type dgrzInitCmd struct {
	Node dgrzInitNode `command:"node" description:"initialize this host as a dogg3rz node" `
	Repo dgrzInitRepo `command:"repository" alias:"repo" description:"initialize a new repository" `
}

type dgrzInitNode struct {
}

type dgrzInitRepo struct {
	Positional struct {
		RepoName string `positional-arg-name:"REPO_NAME" required:"yes" `
	} `positional-args:"yes"`
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzInitCmd{})
}

func (x *dgrzInitNode) Execute(args []string) error {

	var initNode resnode.NodeResource = &implnode.FileNodeResource{}

	err := initNode.InitNode()
	if err != nil {
		return err
	}

	return nil
}

func (x *dgrzInitRepo) Execute(args []string) error {
	fmt.Printf("init repo here ''%s'", x.Positional.RepoName)

	return nil
}

// // IMPLEMENTS 'Commander' interface
func (x *dgrzInitCmd) Execute(args []string) error {

	fmt.Printf("Repo path is '%s'\n", "d")

	return nil
}

func (o *dgrzInitCmd) CommandName() string {
	return "init"
}

func (o *dgrzInitCmd) ShortDescription() string {
	return "initialize a resource"
}

func (o *dgrzInitCmd) LongDescription() string {
	return "initialize a resource"
}
