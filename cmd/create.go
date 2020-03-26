package cmd

import (
	"io"
	"os"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/xeipuuv/gojsonschema"
)

type dgrzCreateCmd struct {
	Schema   dgrzCreateSchema   `command:"schema" description:"create a new schema object" `
	Snapshot dgrzCreateSnapshot `command:"snapshot" alias:"ss" description:"create a snapshot of the repository"`
}

type dgrzCreateSchema struct {
	Positional struct {
		RepoSchemaPath string `positional-arg-name:"REPOSITORY:SCHEMA_PATH" required:"yes" `
	} `positional-args:"yes"`

	// USED FOR SCHEMA PARSING
	jsonLoader    gojsonschema.JSONLoader
	wrappedReader io.Reader
}

type dgrzCreateSnapshot struct {
	Positional struct {
		Repository string `positional-arg-name:"REPOSITORY" required:"yes"`
	} `positional-args:"yes"`
}

///////////////////////////////////////////////////////////
// CREATE COMMAND FUNCTIONS
///////////////////////////////////////////////////////////

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzCreateCmd{})
}

func (o *dgrzCreateCmd) CommandName() string {
	return "create"
}

func (o *dgrzCreateCmd) ShortDescription() string {
	return "create a new repository resource"
}

func (o *dgrzCreateCmd) LongDescription() string {
	return "create a new repository resource"
}

///////////////////////////////////////////////////////////
// CREATE SCHEMA SUBCOMMAND FUNCTIONS
///////////////////////////////////////////////////////////

func (x *dgrzCreateSchema) Execute(args []string) error {

	repo := resource.GetRepositoryResource()

	repoName, schemaPath, err := parseRepoSchemaPath(x.Positional.RepoSchemaPath)
	if err != nil {
		return err
	}

	x.jsonLoader, x.wrappedReader = gojsonschema.NewReaderLoader(os.Stdin)

	yy := repo.CreateSchema(repoName, schemaPath, x)

	return yy
}

func (o *dgrzCreateSchema) Read(p []byte) (n int, err error) {

	// PASS ON TO UNDERLYING READER
	bytesRead, err := o.wrappedReader.Read(p)
	if err == io.EOF {
		// UNDERLYING READER IS FINISHED
		// READING GRACEFULLY. NOW TRY TO VALIDATE
		// THE SCHEMA USING THE JSON LOADER OBJECT
		if _, err := gojsonschema.NewSchema(o.jsonLoader); err != nil {
			// RETURN SCHEMA RELATED ERROR INSTEAD OF EOF
			return bytesRead, errors.UnexpectedValue.Wrap(err, "schema error")
		}

	}

	return bytesRead, err

}

///////////////////////////////////////////////////////////
// CREATE SNAPSHOT SUBCOMMAND  FUNCTIONS
///////////////////////////////////////////////////////////

func (x *dgrzCreateSnapshot) Execute(args []string) error {

	//	fmt.Printf("hello snapshot: { repo = %s }\n", x.Positional.Repository)

	repo := resource.GetRepositoryResource()

	return repo.CreateSnapshot(x.Positional.Repository)

}

func (o *dgrzCreateSnapshot) CommandName() string {
	return "create snapshot"
}

func (o *dgrzCreateSnapshot) ShortDescription() string {
	return "create a repository snapshot"
}

func (o *dgrzCreateSnapshot) LongDescription() string {
	return "create a repository snapshot"
}
