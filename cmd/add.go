package cmd

import (
	"github.com/datacequia/go-dogg3rz/resource"
)

type dgrzAddCmd struct {
	Grapplication string `long:"grapp" short:"r" env:"DOGG3RZ_GRAPP" description:"grapplication name" required:"true"`

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
	return "add a new data(set) resource into a grapplication"
}

func (o *dgrzAddCmd) LongDescription() string {
	return "add a new data(set) resource into a grapplication"
}

func (x *dgrzAddCmd) Execute(args []string) error {

	//	fmt.Printf("hello snapshot: { grapplication = %s }\n", x.Positional.Grapplication)

	ctxt := getCmdContext()

	grapp := resource.GetGrapplicationResource(ctxt)

	return grapp.Add(ctxt, x.Grapplication, x.Positional.Path)

}
