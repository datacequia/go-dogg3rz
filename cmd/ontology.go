package cmd

import (
	"fmt"
	"os"

	"github.com/datacequia/go-dogg3rz/ontology"
)

type dgrzOntologyCmd struct {
	Get dgrzOntologyGetCmd `command:"get"  description:"display the dogg3rz ontology" `
}

type dgrzOntologyGetCmd struct {
}

var ontologyCmd = dgrzOntologyCmd{}
var ontologyGetCmd = dgrzOntologyGetCmd{}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&ontologyCmd)
}

// ONTOLOGY CMD
func (o *dgrzOntologyCmd) CommandName() string {
	return "ontology"
}

func (o *dgrzOntologyCmd) ShortDescription() string {
	return "ontology related commands"
}

func (o *dgrzOntologyCmd) LongDescription() string {
	return "ontology related commands"
}

func (x *dgrzOntologyCmd) Execute(args []string) error {

	//ctxt := getCmdContext()

	fmt.Println("got ontology?")

	return nil
}

// GET CMD
func (o *dgrzOntologyGetCmd) CommandName() string {
	return "ontology"
}

func (o *dgrzOntologyGetCmd) ShortDescription() string {
	return "ontology related commands"
}

func (o *dgrzOntologyGetCmd) LongDescription() string {
	return "ontology related commands"
}

func (x *dgrzOntologyGetCmd) Execute(args []string) error {

	//ctxt := getCmdContext()

	ontology.Get(os.Stdout)

	return nil
}
