package cmd

import (
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
)

// PARSE REPO / SCHEMA PATH IN  REPO:SCHEMA_PATH FORMAT
func parseRepoSchemaPath(repoSchemaPath string) (string, string, error) {

	elements := strings.SplitN(repoSchemaPath, ":", 2)
	if len(elements) != 2 {
		return "", "", errors.UnexpectedValue.Newf("found '%s': want format REPO:SCHEMA_PATH",
			repoSchemaPath)

	}

	return elements[0], elements[1], nil

}
