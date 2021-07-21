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

package resource

import (
	"context"
	"fmt"

	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"
	resourcenode "github.com/datacequia/go-dogg3rz/resource/node"
	resourcerepo "github.com/datacequia/go-dogg3rz/resource/repo"

	//fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	filerepo "github.com/datacequia/go-dogg3rz/impl/file/repo"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/util"
)

/*
var configResource resourceconfig.ConfigResource
var nodeResource resourcenode.NodeResource
var repoResource resourcerepo.RepositoryResource
*/

// INITIALIZES THE DESIRED
// PERSISTENCE IMPLEMENTATION TO USE FOR RESOURCES
// AS SPECIFIED IN  ENV VAR 'DOGG3RZ_STATE_STORE',
// AND ASSIGNS IT TO IT'S CORRESPONDING RESOURCE TYPE INTERFACE
// FOR USE IN REST OF THE CODE WHEN INTERACTING WITH THESE RESOURCE TYPES

//const EnvDogg3rzStateStore = "DOGG3RZ_STATE_STORE"
const StateStoreTypeFile = "file"

/// GetConfigResource returns a ConfigResource which allows the caller
// to interact with the node configuration at runtime for all supported
// configuration type operations
func GetConfigResource(ctxt context.Context) resourceconfig.ConfigResource {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) // , defaultValue)appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:
		return &fileconfig.FileConfigResource{}

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}

	//	return configResource
}

// GetNodeResource returns a NodeResource which allows the caller
// to interact with the node at runtime for all supported
// node type operations
func GetNodeResource(ctxt context.Context) resourcenode.NodeResource {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) // appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:
		return &filenode.FileNodeResource{}

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}
}

// GetRepositoryResource returns a RepositoryResource which allows the caller
// to interact with the configured repository type at runtime for most
// repository operations
func GetRepositoryResource(ctxt context.Context) resourcerepo.RepositoryResource {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) //appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:
		return &filerepo.FileRepositoryResource{}

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}
}

// GetRepositoryResourceStager returns a RepositoryResourceStager which allows the caller
// to interact with the configured repository type at runtime for staging type
// repository operations
func GetRepositoryResourceStager(ctxt context.Context, repoName string) (resourcerepo.RepositoryResourceStager, error) {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) //appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:

		return filerepo.NewFileRepositoryResourceStager(ctxt, repoName)

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}
}
