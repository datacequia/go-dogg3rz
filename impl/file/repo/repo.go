/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package repo

import (
	"os"
	"path"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

type FileRepositoryResource struct {
}

func (repo *FileRepositoryResource) InitRepo(name string) error {

	repoDir := path.Join(file.RepositoriesDirPath(), name)
	// CREATE 'refs/heads' SUBDIR
	refsDir := path.Join(repoDir, "refs")
	headsDir := path.Join(refsDir, "heads")

	dirsList := []string{repoDir, refsDir, headsDir}

	for _, d := range dirsList {

		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			if os.IsNotExist(err) {
				// BASE REPO DIR DOES NOT EXIST
				return dgrzerr.NotFound.Wrapf(err, file.RepositoriesDirPath())
			}
			if os.IsExist(err) {

				return dgrzerr.AlreadyExists.Wrapf(err, name)
			}

			return err

		}
		// WRITE THE HEAD FILE WITH A POINTER TO DEFAULT MASTER BRANCH
		err = file.WriteHeadFile(name, file.MasterBranchName)
		if err != nil {
			return err
		}

	}

	return nil

}
