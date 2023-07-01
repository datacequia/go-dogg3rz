package cmd

import (
	"github.com/datacequia/go-dogg3rz/resource"
)

type dgrzAddCmd struct {
	Repository string `long:"repo" short:"r" env:"DOGG3RZ_REPO" description:"repository name" required:"true"`

	Positional struct {
		Path string `positional-arg-name:"FILE" description:"path to resource file (.jsonld)" required:"yes" `
	} `positional-args:"yes"`

	// USED FOR SCHEMA PARSING

}

var addCmd = dgrzAddCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&addCmd)
}

func (o *dgrzAddCmd) CommandName() string {
	return "add"
}

func (o *dgrzAddCmd) ShortDescription() string {
	return "add a new data(set) resource into a repository"
}

func (o *dgrzAddCmd) LongDescription() string {
	return "add a new data(set) resource into a repository"
}

func (x *dgrzAddCmd) Execute(args []string) error {

	//	fmt.Printf("hello snapshot: { repo = %s }\n", x.Positional.Repository)

	ctxt := getCmdContext()

	repo := resource.GetRepositoryResource(ctxt)

	//return repo.CreateSnapshot(ctxt, createCmd.Repository)

	return repo.Add(ctxt, x.Repository, x.Positional.Path)

	
}
