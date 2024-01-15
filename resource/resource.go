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

	filegrapp "github.com/datacequia/go-dogg3rz/impl/file/grapp"
	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"
	resourcegrapp "github.com/datacequia/go-dogg3rz/resource/grapp"

	//fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/util"
)

/*
var configResource resourceconfig.ConfigResource
var nodeResource resourcenode.NodeResource
var grappResource resourcegrapp.GrapplicationResource
*/

// INITIALIZES THE DESIRED
// PERSISTENCE IMPLEMENTATION TO USE FOR RESOURCES
// AS SPECIFIED IN  ENV VAR 'DOGG3RZ_STATE_STORE',
// AND ASSIGNS IT TO IT'S CORRESPONDING RESOURCE TYPE INTERFACE
// FOR USE IN REST OF THE CODE WHEN INTERACTING WITH THESE RESOURCE TYPES

// const EnvDogg3rzStateStore = "DOGG3RZ_STATE_STORE"
const StateStoreTypeFile = "file"

// / GetConfigResource returns a ConfigResource which allows the caller
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
/*
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
*/

// GetGrapplicationResource returns a GrapplicationResource which allows the caller
// to interact with the configured grapplication type at runtime for most
// grapplication operations
func GetGrapplicationResource(ctxt context.Context) resourcegrapp.GrapplicationResource {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) //appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:
		return &filegrapp.FileGrapplicationResource{}

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}
}

// GetGrapplicationResourceStager returns a GrapplicationResourceStager which allows the caller
// to interact with the configured grapplication type at runtime for staging type
// grapplication operations
/*
func GetGrapplicationResourceStager(ctxt context.Context, grappName string) (resourcegrapp.GrapplicationResourceStager, error) {

	storeType := util.ContextValueAsStringOrDefault(ctxt, env.EnvDogg3rzStateStore, StateStoreTypeFile) //appCtxt.GetOrDefault("DOGG3RZ_STATE_STORE", StateStoreTypeFile)
	switch storeType {
	case StateStoreTypeFile:

		return filegrapp.NewFileGrapplicationResourceStager(ctxt, grappName)

	default:
		panic(fmt.Sprintf(
			"unknown store type assigned to '%s' app context variable: %s",
			env.EnvDogg3rzStateStore,
			storeType))

	}
}
*/
