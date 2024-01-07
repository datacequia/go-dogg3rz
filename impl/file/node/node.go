/*
 * Copyright (c) 2019-2020 Datacequia LLC. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

package node

import (
	"context"
	"os"

	//"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"

	//	"github.com/datacequia/go-dogg3rz/impl/file/config"
	"github.com/datacequia/go-dogg3rz/impl/file/config"
	conf "github.com/datacequia/go-dogg3rz/resource/config"
)

type FileNodeResource struct {
}

func (node *FileNodeResource) InitNode(ctxt context.Context, c conf.Dogg3rzConfig) error {

	file.DotDirPath(ctxt)

	createDirList := []string{file.DotDirPath(ctxt), file.DataDirPath(ctxt) /*, file.GrapplicationsDirPath(ctxt)*/}

	for _, d := range createDirList {
		// CREATE DIR SO THAT ONLY USER CAN R/W
		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			return err
		}
	}

	err := config.SetConfigDefault(ctxt, c)
	if err != nil {
		return err
	}

	return nil

}

func (node *FileNodeResource) GetGrapps(ctxt context.Context) ([]string, error) {
	/*
		root := file.GrapplicationsDirPath(ctxt)

		var files []string
		var err error

		if file.DirExists(root) {
			files, err = file.GetDirs(root)
		}
	*/

	return nil, errors.GrappError.New("not implemented")
}
