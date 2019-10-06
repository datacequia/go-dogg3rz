package cmd

import (
	"fmt"
	//config "github.com/ipfs/go-ipfs-config"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	//	"os"

	"github.com/datacequia/go-dogg3rz/resource"
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

	return resource.GetNodeResource().InitNode()

}

func (x *dgrzInitRepo) Execute(args []string) error {

	return resource.GetRepositoryResource().InitRepo(x.Positional.RepoName)

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
