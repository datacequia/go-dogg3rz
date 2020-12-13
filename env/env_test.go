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
	"testing"
)

func TestInitContextFromEnv(t *testing.T) {

	const (
		dgrzHomeValue       = "/home/test"
		dgrzRepoValue       = "myrepo"
		dgrzStateStoreValue = "file"
	)

	os.Setenv(EnvDogg3rzHome, dgrzHomeValue)
	os.Setenv(EnvDogg3rzRepo, dgrzRepoValue)
	os.Setenv(EnvDogg3rzStateStore, dgrzStateStoreValue)

	ctxt := InitContextFromEnv(context.Background())

	var i interface{}
	i = ctxt.Value(EnvDogg3rzHome)
	if value, ok := i.(string); !(ok && value == dgrzHomeValue) {
		t.Errorf("expected ctxt.Value() to return '%s' value for key '%s', found '%v' of type %T",
			dgrzHomeValue, EnvDogg3rzHome, i, i)
	}

	i = ctxt.Value(EnvDogg3rzRepo)
	if value, ok := i.(string); !(ok && value == dgrzRepoValue) {
		t.Errorf("expected ctxt.Value() to return '%s' value for key '%s', found '%v' of type %T",
			dgrzRepoValue, EnvDogg3rzRepo, i, i)
	}

	i = ctxt.Value(EnvDogg3rzStateStore)
	if value, ok := i.(string); !(ok && value == dgrzStateStoreValue) {
		t.Errorf("expected ctxt.Value() to return '%s' value for key '%s', found '%v' of type %T",
			dgrzStateStoreValue, EnvDogg3rzStateStore, i, i)
	}

}
