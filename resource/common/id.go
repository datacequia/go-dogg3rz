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

package common

/*

Unix format
/dgrz/[user@]COMMIT_MULTIHASH[]/RESOURCE_KINDs/RESOURCE_SUBPATH]

http format


*/

// THIS INTERFACE IS IMPLEMENTED BY DOGG3RZ REPO OBJECTS
// TO PROVIDE A GLOBALLY UNIQUE IDENTIFER TO THAT RESOURCE
//
// ANY ID'S THAT HAVE SAME VALUES FOR ALL ATTRIBUTES
//  ARE THOUGHT TO IDENTIFY THE SAME
// RESOURCE. THIS IS DICTATED BY THE IMPLEMENTER OF THIS INTERFACE

const RootPathElementName = "dgrz"

// INTERFACE THAT IDENTIFIES A SPECIFIC
//  USER CREATED REPOSITORY RESOURCE (PRIMITIVE)
type RepositoryResourceId interface {
	User() string            // THE USER CONTEXT. EMPTY STRIING IF NO USER CONTEXT
	CommitMultiHash() string // THE COMMIT HASH FOR THE OBJECT. EMPTY STRING
	// IF OBJECT NOT YET STAGED
	Kind() string // THE KIND OF OBJECT PRIMITIVE

	Subpath() string // THE SUBPATH THAT PROVIDES UNIQUENESS TO THE OBJECT
	// CAN BE THE OBJECTS NAME IF RELEVANT OR SOME OTHER
	// AGGREGATION OF IT'S ATTIBUTES THAT PROVIDE
	// UNIQUENESS

}
