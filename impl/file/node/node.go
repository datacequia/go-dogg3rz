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

package node

import (
	"os"

	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/impl/file/config"
)

type FileNodeResource struct {
}

func (node *FileNodeResource) InitNode() error {

	file.DotDirPath()

	createDirList := []string{file.DotDirPath(), file.DataDirPath(), file.RepositoriesDirPath()}

	for _, d := range createDirList {
		// CREATE DIR SO THAT ONLY USER CAN R/W
		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			return err
		}
	}

	err := config.SetConfigDefault()
	if err != nil {
		return err
	}

	return nil

}
