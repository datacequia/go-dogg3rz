package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource"
	"github.com/xeipuuv/gojsonschema"
)

type dgrzCreateCmd struct {
	Schema dgrzCreateSchema `command:"schema" description:"create a new schema object" `
}

type dgrzCreateSchema struct {
	Positional struct {
		RepoSchemaPath string `positional-arg-name:"REPO:SCHEMA_PATH" required:"yes" `
	} `positional-args:"yes"`

	// USED FOR SCHEMA PARSING
	jsonLoader    gojsonschema.JSONLoader
	wrappedReader io.Reader
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzCreateCmd{})
}

func (x *dgrzCreateSchema) Execute(args []string) error {

	repo := resource.GetRepositoryResource()

	repoSchemaPath := strings.SplitN(x.Positional.RepoSchemaPath, ":", 2)
	if len(repoSchemaPath) != 2 {
		return errors.UnexpectedValue.Newf("found '%s': want format REPO:SCHEMA_PATH",
			x.Positional.RepoSchemaPath)

	}

	repoName := repoSchemaPath[0]
	schemaPath := repoSchemaPath[1]

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

func (o *dgrzCreateCmd) CommandName() string {
	return "create"
}

func (o *dgrzCreateCmd) ShortDescription() string {
	return "create a new repository resource"
}

func (o *dgrzCreateCmd) LongDescription() string {
	return "create a new repository resource"
}
