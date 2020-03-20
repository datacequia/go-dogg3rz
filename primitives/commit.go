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

package primitives

//const D_ATTR_NAME = "objectHeads"
const TYPE_DOGG3RZ_COMMIT = "dogg3rz.commit"

//const MD_ATTR_NAME = "name"
//const MD_ATTR_IPFS_PEER_ID = "ipfsPeerId"
const MD_ATTR_EMAIL_ADDR = "emailAddress"
const MD_ATTR_REPO_NAME = "repositoryName"
const MD_ATTR_REPO_ID = "repositoryId"

//const D_ATTR_ROOT_TREE = "rootTree"
const D_ATTR_TRIPLES = "triples"
const D_ATTR_IMPORTS = "imports"

var reservedMDAttrCommit = [...]string{MD_ATTR_REPO_NAME, MD_ATTR_EMAIL_ADDR}

var reservedDAttrCommit = [...]string{D_ATTR_TRIPLES, D_ATTR_IMPORTS}

type dgrzCommit struct {
	repositoryId string // GLOBALLY UNIQUE IMMUTABLE IDENTIFIER FOR REPOSITORY
	// THAT THIS COMMIT BELONGS TO.
	emailAddress string // EMAIL ADDRESS OF COMMITTER

	repositoryName string // REPOSITORY NAME FOR THIS COMMIT

	parents []string
}
