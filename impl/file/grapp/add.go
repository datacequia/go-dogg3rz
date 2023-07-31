package grapp

import (
	"context"

	"fmt"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

func add(ctxt context.Context, grappName string, path string) error {

	if !file.GrapplicationExist(ctxt, grappName) {
		return errors.NotFound.Newf("grapplication '%s' does not exist",
			grappName)
	}

	grappFsPath := file.GrapplicationDgrzDirPath(ctxt, grappName)

	fmt.Println(grappFsPath)

	return nil
}
