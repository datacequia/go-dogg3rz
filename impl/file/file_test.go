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

package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/resource/config"
)

func nodeGrappSetup(t *testing.T, prefix string) string {

	dogg3rzHome := filepath.Join(os.TempDir(),
		fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()))

	var dgrzConf config.Dogg3rzConfig

	// REQUIRED CONF
	dgrzConf.User.Email = "test@dogg3rz.com"

	t.Logf("created DOGG3RZ_HOME at %s", dogg3rzHome)

	return dogg3rzHome

}

func nodeGrappTeardown(t *testing.T, dogg3rzHome string) {

	os.RemoveAll(dogg3rzHome)

}

func TestWriteToFileAtomic(t *testing.T) {

	theFile := path.Join(os.TempDir(), "HEAD")
	theFileLock := theFile + LOCK_FILE_SUFFIX

	s := "test data"

	// REMOVE ANY DANGLING LOCK FILE IF IT EXISTS
	// BEFORE THE TEST
	os.Remove(theFileLock)

	// TEST
	reader := strings.NewReader(s)
	bytesWritten, err := WriteToFileAtomic(func() (io.Reader, error) { return reader, nil }, theFile)
	if err != nil {

		if FileExists(theFileLock) {
			// the lock file, if created, should have been removed after fail
			t.Errorf("lock file exists after failed WriteFileToAtomic: %s", theFileLock)
		}

		t.Errorf("WriteFileAtomic failed: %s", err)

	}

	// CHECK THAT BYTES WRITTEN EQUALS BYTES SEND OF
	// TEST FILE CONTENT
	if bytesWritten != int64(len(s)) {
		t.Errorf("WriteToFileAtomic failed: bytes written != bytes sent: { bytes written = %d, bytes sent %d}", bytesWritten, len(s))

	}

	// NOW CREATE LOCK FILE ARTIFICIALLY
	// SHOULD EXPECT A 'TryAgain' error to return
	err = Touch(theFileLock)
	if err != nil {
		t.Fail()
	}

	bytesWritten, err = WriteToFileAtomic(func() (io.Reader, error) { return strings.NewReader(s), nil }, theFile)
	if err != nil {
		if dgrzerr.GetType(err) != dgrzerr.TryAgain {
			// Error Type is not a TryAgain Errors
			t.Errorf("Expected TryAgain error type, got %s: %s", reflect.TypeOf(err), err)
		}
	} else {
		// THERE WAS AN ERROR BUT..
		// SHOULD HAVE FAILED WITH 'TryAgain' exception

		// Error Type is not a TryAgain Errors
		t.Error("Expected TryAgain error type, got no error!")

	}

}

func TestGrappDirList(t *testing.T) {

	// CREATE A NEW NODE/GRAPP SANDBOX FOR TESTING
	dogg3rzHome := nodeGrappSetup(t, "file_test")
	defer nodeGrappTeardown(t, dogg3rzHome)

}
