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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/primitives"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/google/uuid"
	//	"github.com/datacequia/go-dogg3rz/impl/file/config"
)

// FILESTORE CONSTANTS
const dotDirName = ".dogg3rz"
const LOCK_FILE_SUFFIX = ".lock"
const dataDirName = "data"
const repositoriesDirName = "repositories"
const RefsDirName = "refs"
const HeadsDirName = "heads"
const MasterBranchName = "master"
const IndexFileName = "index"

const ResourceCacheSignature = "RESC"

// Writes contents of Reader object to 'path' atomically
// i.e. no other writers can write at the same time.
// An attempt for other writers to do so simultaneously
// will result inn a 'TryAgain' error being returned
// RETURNS PathError or TryAgain error types
func WriteToFileAtomic(readerFunc func() (io.Reader, error), path string) (int64, error) {

	var bytesWritten int64 = 0

	// CREATE TEMP FILE IN SYSTEM TEMP DIR
	// BY ADDING .lock SUFFIX
	lockFile := path + ".lock"

	var lf *os.File
	var err error
	var r io.Reader
	// OPEN LOCK FILE EXCLUSIVELY

	lf, err = os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.FileMode(0600))
	if err != nil {
		if os.IsExist(err) {
			// Lock file exists!
			// NOTIFY USER TO TRY AGAIN.
			// NOTE, IT COULD BE THE CASE THAT THE LOCK FILE WAS ORPHANED BY ANOTHER PROCESS/THREAD
			// AND IT'S PREVENTING SUBSQUENT OPERATIONS ON THE RESOURCE UNNECESSARILY
			// IN THIS CASE ONLY RECOURSE IS TO SHUTDOWN dogg3rz and MANUUALLY
			// REMOVE LOCK FILE
			return 0, dgrzerr.TryAgain.Wrapf(err, "resource is temporarily unavailable. try operation again later...")
			// ANOTHER PROCESS/THREAD IS TRYING TO WRITE TO THIS FILE

		}
		// OTHERWISE RETURN ORIGINAL Errors
		return bytesWritten, err
	}

	r, err = readerFunc()
	if err != nil {
		// FAILED TO GET READER.
		// RM LOCK FILE AND EXIT
		err = os.Remove(lockFile)
		return bytesWritten, err

	}

	// OPENING LOCK FILE SUCCEEDED. COPY DATA FROM Reader
	bytesWritten, err = io.Copy(lf, r)

	// CLOSE THE LOCK FILE BEFORE DOING ANYTHING ELSE
	lf.Close()

	if err != nil {
		// COPY FAILED. REMOVE LOCK FILE
		// OTHER  WRITERS CAN'T CREATE LOCK FILE
		// UNTIL IT IS REEMOVED
		errRm := os.Remove(lockFile)
		if errRm != nil {
			panic(fmt.Sprintf("failed to remove lock file on error: %s", errRm))
		}

	} else {
		// RENAME COPIED CONTENTS TO TARGET PATH (ATOMIC UPDATE OF CONTENT)
		err = os.Rename(lockFile, path)
	}

	return bytesWritten, err

}

// Creates an empty file at 'path' similar to Unix touch command
func Touch(path string) error {

	newFile, err := os.Create(path)
	if err != nil {
		return err
	}

	newFile.Close()

	return nil
}

func DotDirPath() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// CAN'T FETCH THE HOMEDIR???
		// BAIL!
		log.Panicf("can't find user home directory: %s", err)
	}

	return path.Join(homeDir, dotDirName)

}

func DataDirPath() string {
	return path.Join(DotDirPath(), dataDirName)
}

func RepositoriesDirPath() string {
	return path.Join(DataDirPath(), repositoriesDirName)

}

func RepositoriesRefsDirPath() string {
	return path.Join(RepositoriesDirPath(), RefsDirName)
}

func RepositoriesRefsHeadsDirPath() string {
	return path.Join(RepositoriesRefsDirPath(), HeadsDirName)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()

}

func WriteHeadFile(repoName string, branchName string) error {

	content := fmt.Sprintf("ref: %s\n", strings.Join([]string{RefsDirName, HeadsDirName, branchName},
		string(os.PathSeparator)))

	_, err := WriteToFileAtomic(func() (io.Reader, error) { return strings.NewReader(content), nil },
		path.Join(RepositoriesDirPath(), repoName, "HEAD"))

	return err
}

func RepositoryExist(repoName string) bool {

	repoPath := filepath.Join(RepositoriesDirPath(), repoName)

	if info, err := os.Stat(repoPath); err != nil {
		// PATH DOES NOT EXIST

		return false

	} else {
		// REPO DOES EXIST. IS IT A DIR?
		// COULD BE A FILE BUT THAT WOULD BE AN ERROR
		// GIVEN THAT NO NON DIR OBJECTS SHOULD EXIST IN BASE REPO DIR
		return info.IsDir()
	}

}

// CREATES THE RESOURCE PATH IN THE DESIGNATED REPOSITORY OF A SPECIFIC
// RESOURCE TYPE
func CreateRepositoryResourcePath(resPath *rescom.RepositoryPath, repoName string,
	resType string, bodyReader io.Reader) (string, error) {

	if !RepositoryExist(repoName) {
		return "", dgrzerr.NotFound.Newf("repository '%s' does not exist. please create it first", repoName)
	}
	// REPO DOES EXIST. CREATE EACH PATH ELEMENT IF NECESSARY
	curPath := filepath.Join(RepositoriesDirPath(), repoName)
	curResType := primitives.TYPE_DOGG3RZ_TREE

	var success bool = false

	// SETUP CALLBACK TO REMOVE FILESYSTEM PATH using
	// DEFER IN EVENT THAT THIS FUNCTION FAILS
	cleanupResourceOnErrFunc := func(path string) {
		if !success {
			err := os.RemoveAll(path)
			if err != nil {
				log.Printf("failed to remove repository resource %s on error.  "+
					"Please remove manually: %s", path, err)
			}
		}
	}

	for pathElementIndex, path := range resPath.PathElements() {
		curPath = filepath.Join(curPath, path)

		if pathElementIndex == (resPath.Size() - 1) {
			// LAST ELEMENT. MAKE CUR RESOURCE TYPE
			// THE DESIRED RESOURCE TYPE
			curResType = resType
		}

		// EVAL CURRENT REPO PATH TO ENSURE IT'S A DIRECTORY AND
		// A DOGG3RZ TREE OBJECT
		if _, err := os.Stat(curPath); err != nil {
			if os.IsNotExist(err) {

				// CURRENT PATH DOES NOT EXIST. CREATE IT
				if err := os.Mkdir(curPath, os.FileMode(0700)); err != nil {
					return "", err
				}
				// DELETE THIS DIRECTORY IF THIS FUNCTION FAILS
				defer cleanupResourceOnErrFunc(curPath)

				// NOW CREATE '.type' attr file
				cbFunc := func() (io.Reader, error) {
					return strings.NewReader(curResType), nil
				}

				if _, err := WriteToFileAtomic(cbFunc,
					filepath.Join(curPath, ".type")); err != nil {
					return "", err
				}

				// NOW CREATE '.id' attribute file containing unique uuid
				cbFunc = func() (io.Reader, error) {
					treeUUID := uuid.New().String()
					return strings.NewReader(treeUUID), nil
				}

				if _, err := WriteToFileAtomic(cbFunc,
					filepath.Join(curPath, ".id")); err != nil {
					return "", err
				}

			} else {
				// SOME OTHER (SYSTEM?) ERROR OCCURRED. RETURN IT
				return "", err
			}
		} else {
			// PATH EXISTS. IS IT A DOGG3RZ TREE OBJECT?
			typeAttrPath := filepath.Join(curPath, ".type")
			if content, err := ioutil.ReadFile(typeAttrPath); err != nil {
				// ERROR READING .type ATTR FILE. ALL DIRECTORIES
				// IN A REPO SHOULD HAVE ONE
				return "", err
			} else {
				// READ CONTENT. IS IT A DOGGERZ TREE OBJECT?
				if string(content) != primitives.TYPE_DOGG3RZ_TREE {

					if pathElementIndex == (resPath.Size() - 1) {
						return "", dgrzerr.AlreadyExists.Newf(
							"%s: type = %s",
							resPath.ToString(),
							content)

					} else {

						return "", dgrzerr.InvalidPathElement.Newf(
							"encountered invalid base path '%s' for creation of "+
								"repository resource '%s': want type '%s', found type '%s'",
							curPath,
							resPath.ToString(),
							primitives.TYPE_DOGG3RZ_TREE,
							content)

					}

				}
				// IS A TREE OBJECT. GTG...
			}

		}
	}
	// NOW WRITE THE BODY
	bodyFunc := func() (io.Reader, error) {

		return bodyReader, nil
	}

	if _, err := WriteToFileAtomic(bodyFunc,
		filepath.Join(curPath, ".body")); err != nil {
		return "", err
	}

	// FLAG AS SUCCESSFUL SO DEFER FUNC WILL NOT
	// REMOVE DIRECTORIES CREATED
	success = true

	return curPath, nil

}
