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

package repo

import (
	"context"
	"os"
	"path"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
)

const indexFileName = "index"

type FileRepositoryResource struct {
}

func (repo *FileRepositoryResource) InitRepo(ctxt context.Context, name string) error {

	repoDir := path.Join(file.RepositoriesDirPath(ctxt), name)
	// CREATE 'refs/heads' SUBDIR
	// CREATE '.dgrz' DIR AS SUBDI OF BASE REPO DIR
	dgrzDir := path.Join(repoDir, file.DgrzDirName)
	refsDir := path.Join(dgrzDir, file.RefsDirName)
	headsDir := path.Join(refsDir, file.HeadsDirName)

	dirsList := []string{repoDir, dgrzDir, refsDir, headsDir}

	for _, d := range dirsList {

		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			if os.IsNotExist(err) {
				// BASE REPO DIR DOES NOT EXIST
				return dgrzerr.NotFound.Wrapf(err, file.RepositoriesDirPath(ctxt))
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

// StageResources stages resources starting with and contained within the location
// specified by startList .
// Returns list of staged resources
func (repo *FileRepositoryResource) StageResources(ctxt context.Context, repoName string, startList []rescom.StagingResourceLocation) ([]rescom.StagingResource, error) {

	sr := &fileStageResource{}

	return sr.stageResources(ctxt, repoName, startList)
}

func (repo *FileRepositoryResource) CreateSnapshot(ctxt context.Context, repoName string) error {

	ss := &fileCreateSnapshot{}

	return ss.createSnapshot(ctxt, repoName)

}

func (repo *FileRepositoryResource) CreateDataset(ctxt context.Context, repoName string, datasetPath string) error {

	var fds *fileDataset
	var err error

	if fds, err = newFileDataset(ctxt, repoName, datasetPath); err != nil {
		return err
	}

	if err := fds.create(ctxt); err != nil {
		return err
	}

	return nil
}

func (repo *FileRepositoryResource) AddNamespaceDataset(ctxt context.Context, repoName string, datasetPath string, term string, iri string) error {

	if err := addNamespaceDataset(ctxt, repoName, datasetPath, term, iri); err != nil {
		return err
	}

	return nil
}

func (repo *FileRepositoryResource) AddNamespaceNode(ctxt context.Context, repoName string, datasetPath string, nodeID string, term string, iri string) error {

	o := &addNamespaceNode{}

	if err := o.execute(ctxt, repoName, datasetPath, nodeID, term, iri); err != nil {
		return err
	}

	return nil
}

func (repo *FileRepositoryResource) GetDataSets(ctxt context.Context, repoName string) ([]string, error) {

	repoDir := path.Join(file.RepositoriesDirPath(ctxt), repoName)
	var files []string
	var err error

	if file.DirExists(repoDir) {
		files, err = file.GetDirs(repoDir)
	}

	// Ignore any dogg3rz internal dirs and files.
	for i, v := range files {
		if strings.HasPrefix(v, ".") {
			files = append(files[:i], files[i+1:]...)
		}
	}

	return files, err
}
