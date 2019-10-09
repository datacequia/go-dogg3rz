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

package file

import (
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
)

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
