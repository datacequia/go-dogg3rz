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

package grapp

import (
	"context"
	"os"
	"path"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

const indexFileName = "index"

type FileGrapplicationResource struct {
}

func (grapp *FileGrapplicationResource) InitGrapp(ctxt context.Context, name string) error {

	grappDir := path.Join(file.GrapplicationsDirPath(ctxt), name)

	// CREATE DEFAULT BRANCH DIR
	mainBranchDir := path.Join(grappDir, file.MasterBranchName)

	// CREATE 'refs/heads' SUBDIR
	// CREATE '.dgrz' DIR AS SUBDI OF BASE GRAPP DIR
	dgrzDir := path.Join(grappDir, file.DgrzDirName)
	refsDir := path.Join(dgrzDir, file.RefsDirName)
	headsDir := path.Join(refsDir, file.HeadsDirName)

	dirsList := []string{grappDir, mainBranchDir, dgrzDir, refsDir, headsDir}

	for _, d := range dirsList {

		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			if os.IsNotExist(err) {
				// BASE GRAPP DIR DOES NOT EXIST
				return dgrzerr.NotFound.Wrapf(err, file.GrapplicationsDirPath(ctxt))
			}
			if os.IsExist(err) {

				return dgrzerr.AlreadyExists.Wrapf(err, name)
			}

			return err

		}

	}
	// WRITE THE HEAD FILE WITH A POINTER TO DEFAULT MASTER BRANCH
	err := file.WriteHeadFile(ctxt, name, file.MasterBranchName)
	if err != nil {
		return err
	}

	return nil

}

func (grapp *FileGrapplicationResource) CreateSnapshot(ctxt context.Context, grappName string) error {

	ss := &fileCreateSnapshot{}

	return ss.createSnapshot(ctxt, grappName)

}

func (grapp *FileGrapplicationResource) CreateDataset(ctxt context.Context, grappName string, datasetPath string) error {

	var fds *fileDataset
	var err error

	if fds, err = newFileDataset(ctxt, grappName, datasetPath); err != nil {
		return err
	}

	if err := fds.create(ctxt); err != nil {
		return err
	}

	return nil
}

func (grapp *FileGrapplicationResource) AddNamespaceDataset(ctxt context.Context, grappName string, datasetPath string, term string, iri string) error {

	if err := addNamespaceDataset(ctxt, grappName, datasetPath, term, iri); err != nil {
		return err
	}

	return nil
}

func (grapp *FileGrapplicationResource) AddNamespaceNode(ctxt context.Context, grappName string, datasetPath string, nodeID string, term string, iri string) error {

	o := &addNamespaceNode{}

	if err := o.execute(ctxt, grappName, datasetPath, nodeID, term, iri); err != nil {
		return err
	}

	return nil
}

func (grapp *FileGrapplicationResource) GetDataSets(ctxt context.Context, grappName string) ([]string, error) {

	grappDir := path.Join(file.GrapplicationsDirPath(ctxt), grappName)
	var files []string
	var err error

	if file.DirExists(grappDir) {
		files, err = file.GetDirs(grappDir)
	}

	// Ignore any dogg3rz internal dirs and files.
	for i, v := range files {

		if strings.HasPrefix(v, ".") {

			files = append(files[:i], files[i+1:]...)
		}
	}

	return files, err
}

func (grapp *FileGrapplicationResource) Add(ctxt context.Context, grappName string, path string) error {

	return add(ctxt, grappName, path)

}
