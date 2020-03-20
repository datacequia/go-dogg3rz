package cmd

//	"io"
//	"os"
//"strings"
import (
	"github.com/datacequia/go-dogg3rz/resource"
)

//	"github.com/datacequia/go-dogg3rz/errors"

//	"github.com/xeipuuv/gojsonschema"

type dgrzStageCmd struct {
	Positional struct {
		RepoSchemaPath string `positional-arg-name:"REPO:RESOURCE_PATH" required:"yes" `
	} `positional-args:"yes"`
}

func init() {
	// REGISTER THE 'init' COMMAND
	register(&dgrzStageCmd{})
}

func (x *dgrzStageCmd) Execute(args []string) error {

	repoName, schemaSubpath, err := parseRepoSchemaPath(x.Positional.RepoSchemaPath)
	if err != nil {
		return err
	}

	repo := resource.GetRepositoryResource()

	return repo.StageResource(repoName, schemaSubpath)

}

func (o *dgrzStageCmd) CommandName() string {
	return "stage"
}

func (o *dgrzStageCmd) ShortDescription() string {
	return "stage a repository resource"
}

func (o *dgrzStageCmd) LongDescription() string {
	return "stage a new repository resource in your working tree"
}
