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

package resource

import (
	"fmt"
	"os"

	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"
	resourcenode "github.com/datacequia/go-dogg3rz/resource/node"
	resourcerepo "github.com/datacequia/go-dogg3rz/resource/repo"

	//fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	filerepo "github.com/datacequia/go-dogg3rz/impl/file/repo"
)

var configResource resourceconfig.ConfigResource
var nodeResource resourcenode.NodeResource
var repoResource resourcerepo.RepositoryResource

func init() {

	// DETERRMINE STORE TYPE
	storeType := os.Getenv("DOGG3RZ_STATE_STORE")
	if len(storeType) < 1 {
		// DEFAULT TO FILE
		storeType = "file"
	}

	switch storeType {
	case "file":
		nodeResource = &filenode.FileNodeResource{}
		configResource = &fileconfig.FileConfigResource{}
		repoResource = &filerepo.FileRepositoryResource{}
	default:

		fmt.Fprintf(os.Stderr,
			"unknown store type assigned to DOGG3RZ_STATE_STORE environment variable: %s", storeType)
		os.Exit(1)
	}

}

func GetConfigResource() resourceconfig.ConfigResource {

	return configResource
}

func GetNodeResource() resourcenode.NodeResource {
	return nodeResource
}

func GetRepositoryResource() resourcerepo.RepositoryResource {
	return repoResource
}
