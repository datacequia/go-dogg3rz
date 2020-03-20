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
	"io"
)

type RepositoryIndexEntry struct {
}

type RepositoryResource interface {
	InitRepo(repoName string) error
	//	GetRepoIndex(repoName string) (RepositoryIndex, error)

	// REPOSITORY COMMANDS
	CreateSchema(repoName string, schemaSubpath string, schemaReader io.Reader) error
	StageResource(repoName string, schemaSubpath string) error
}

/*
type RepositoryIndex interface {
	Stage(resId rescom.RepositoryResourceId, multiHash string) error
	//AddMap(resIdMap map[rescom.ResourceId]string) error
	Unstage(resId rescom.RepositoryResourceId) error
}
*/
