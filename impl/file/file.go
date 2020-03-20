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
const IndexFileName = ".index"
const DirLockFileName = ".__dirlock__"
const ResourceCacheSignature = "RESC"
const IndexFormatVersion = uint32(1)

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

	var dogg3rzHomeDir string

	if envVal, envIsSet := os.LookupEnv("DOGG3RZ_HOME"); envIsSet {
		dogg3rzHomeDir = envVal
	} else {
		// DEFAULT TO DOT DIR PATH IF DOGG3RZ_HOME IS NOT SET
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// CAN'T FETCH THE HOMEDIR???
			// BAIL!
			log.Panicf("can't find user home directory: %s", err)
		}
		dogg3rzHomeDir = filepath.Join(homeDir, dotDirName)
	}

	return dogg3rzHomeDir

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
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panicf("can't stat path %s: %v", path, err)
	}
	return !info.IsDir()

}

func DirExists(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panicf("can't stat path %s: %v", path, err)
	}

	return info.IsDir()

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

func GetResourceAttributeCB(resPath string, attrName string, cb func(io.Reader,
	os.FileInfo) error) error {

	dotFile := "." + attrName

	attrPath := filepath.Join(resPath, dotFile)

	if fp, err := os.Open(attrPath); err != nil {
		return err
	} else {

		defer fp.Close()

		// * call stat on open file to get metadata
		if fileInfo, err := fp.Stat(); err != nil {
			return err
		} else {
			if err := cb(fp, fileInfo); err != nil {
				return err
			}
		}

	}

	return nil
}

func GetResourceAttributeS(resPath string, attrName string) (string, error) {

	var attrValue []byte
	var err error

	cb := func(r io.Reader, fileInfo os.FileInfo) error {
		attrValue, err = ioutil.ReadAll(r)

		return err
	}

	err = GetResourceAttributeCB(resPath, attrName, cb)

	return string(attrValue), err

}

// WRITE ATTRIBUTE FILE TO RESOURCE AT 'resPath' (a directory)
func PutResourceAttribute(resPath string, attrName string, attrValue io.Reader) (string, error) {

	dotFile := "." + attrName

	attrPath := filepath.Join(resPath, dotFile)

	cbFunc := func() (io.Reader, error) {
		return attrValue, nil
	}

	if _, err := WriteToFileAtomic(cbFunc, attrPath); err != nil {
		return attrPath, err
	}

	return attrPath, nil

}

func PutResourceAttributeS(resPath string, attrName string,
	attrValue string) (string, error) {

	attrValueAsReader := strings.NewReader(attrValue)

	return PutResourceAttribute(resPath, attrName, attrValueAsReader)

}
