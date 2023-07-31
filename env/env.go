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

package env

import (
	"context"
	"os"
	"path/filepath"
)

const (
	// PREFIX FOR ALL DOGG3RZ O/S ENVIRONMENT VARIABLES
	EnvDogg3rzPrefix = "DOGG3RZ_"
	// SPECIFIES PATH TO CURRENT DOGG3RZ GRAPPLICATION DIRECTORY CONTEXT  (OPTIONAL)
	EnvDogg3rzGrapp = EnvDogg3rzPrefix + "GRAPP"
	// SPECIFIES PATH TO BASE/HOME DIECTORY FOR DOGG3RZ FILES (OPTIONAL)
	EnvDogg3rzHome = EnvDogg3rzPrefix + "HOME"
	// SPECIFIES THE PERSISTENCE TYPE FOR PERSISTING STATE IN DOGG3RZ
	// (CURRENTLY DEFAULTS TO 'file' IF NOT SET)
	EnvDogg3rzStateStore = EnvDogg3rzPrefix + "STATE_STORE"
)

var (
	applicationName string
)

func init() {
	applicationName = filepath.Base(os.Args[0])

}

var envNames = []string{
	EnvDogg3rzGrapp,
	EnvDogg3rzHome,
	EnvDogg3rzStateStore,
}

// InitContextFromEnv sets and returns  a new context initialized from
// designated Dogg3rz environment variables set for the running dogg3rz process
func InitContextFromEnv(ctxt context.Context) context.Context {

	for _, envName := range envNames {

		if envValue, ok := os.LookupEnv(envName); ok {
			ctxt = context.WithValue(ctxt, envName, envValue)
		}
	}

	return ctxt

}

func ApplicationName() string {
	return applicationName
}
