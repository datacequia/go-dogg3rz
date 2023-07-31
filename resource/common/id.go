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
//
//	USER CREATED REPOSITORY RESOURCE (PRIMITIVE)
type GrapplicationResourceId interface {
	User() string            // THE USER CONTEXT. EMPTY STRIING IF NO USER CONTEXT
	CommitMultiHash() string // THE COMMIT HASH FOR THE OBJECT. EMPTY STRING
	// IF OBJECT NOT YET STAGED
	Kind() string // THE KIND OF OBJECT PRIMITIVE

	Subpath() string // THE SUBPATH THAT PROVIDES UNIQUENESS TO THE OBJECT
	// CAN BE THE OBJECTS NAME IF RELEVANT OR SOME OTHER
	// AGGREGATION OF IT'S ATTIBUTES THAT PROVIDE
	// UNIQUENESS

}
