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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/datacequia/go-dogg3rz/env"
	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/util"
	//	"github.com/datacequia/go-dogg3rz/impl/file/config"
)

// FILESTORE CONSTANTS
const dotDirName = ".dogg3rz"
const LOCK_FILE_SUFFIX = ".lock"
const dataDirName = "data"
const repositoriesDirName = "repositories"

// BASE REPO DIR FOR ALL STATE FILES
const DgrzDirName = ".dgrz"
const RefsDirName = "refs"
const HeadsDirName = "heads"
const MasterBranchName = "master"
const IndexFileName = ".index"
const DirLockFileName = ".__dirlock__"
const ResourceCacheSignature = "RESC"
const IndexFormatVersion = uint32(1)
const HeadFileName = "HEAD"
const RepositoryIdFileName = "ID"
const JSONLDDocumentName = ".document.jsonld"

var validPathElementRegex = regexp.MustCompilePOSIX("^[a-z][-a-z0-9]*$")

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
		if errRm := os.Remove(lockFile); errRm != nil {
			log.Printf("failed to remove lockfile at %s: %s", lockFile, errRm)
		}

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

func DotDirPath(ctxt context.Context) string {

	var dogg3rzHomeDir string
	var ok bool

	if dogg3rzHomeDir, ok = util.ContextValueAsString(ctxt, env.EnvDogg3rzHome); (!ok) || dogg3rzHomeDir == "" {
		//	if envVal, envIsSet := os.LookupEnv("DOGG3RZ_HOME"); envIsSet {

		//	} else {
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

func DataDirPath(ctxt context.Context) string {
	return path.Join(DotDirPath(ctxt), dataDirName)
}

func RepositoriesDirPath(ctxt context.Context) string {
	return path.Join(DataDirPath(ctxt), repositoriesDirName)

}
func RepositoriesDgrzDirPath(repoName string, ctxt context.Context) string {
	return path.Join(RepositoriesDirPath(ctxt), repoName, DgrzDirName)
}

func RepositoriesRefsDirPath(repoName string, ctxt context.Context) string {
	return path.Join(RepositoriesDgrzDirPath(repoName, ctxt), RefsDirName)
}

func RepositoriesRefsHeadsDirPath(repoName string, ctxt context.Context) string {
	return path.Join(RepositoriesRefsDirPath(repoName, ctxt), HeadsDirName)
}

// returns list of directory nams that are repository dirs
func RepositoryDirList(ctxt context.Context) ([]string, error) {

	//	baseRepoDir := RepositoryDirPath()

	var dirListing []os.FileInfo
	var err error

	if dirListing, err = ioutil.ReadDir(RepositoriesDirPath(ctxt)); err != nil {
		return nil, err
	}

	repoDirList := make([]string, 0)
	for _, fileInfo := range dirListing {
		// if directory entry is a directory and has a subdirectory entry itself named '.dgrz'
		// then we can assume it's a repository dir
		if fileInfo.IsDir() && DirExists(filepath.Join(RepositoriesDirPath(ctxt), DgrzDirName)) {
			// IS A REPO DIR ADD IT
			repoDirList = append(repoDirList, fileInfo.Name())
		}
	}

	return repoDirList, nil

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

func WriteHeadFile(repoName string, branchName string, ctxt context.Context) error {

	content := fmt.Sprintf("ref: %s\n", filepath.Join(RefsDirName, HeadsDirName, branchName))

	_, err := WriteToFileAtomic(func() (io.Reader, error) { return strings.NewReader(content), nil },
		path.Join(RepositoriesDirPath(ctxt), repoName, DgrzDirName, HeadFileName))

	return err
}

func WriteCommitHashToCurrentBranchHeadFile(repoName string, commitHash string, ctxt context.Context) error {

	headFile := path.Join(RepositoriesDirPath(ctxt), repoName, DgrzDirName, HeadFileName)

	if f, err := os.Open(headFile); err != nil {
		return err
	} else {
		defer f.Close()

		var headFileSubPath string
		if n, err := fmt.Fscanf(f, "ref: %s", &headFileSubPath); err != nil {
			return err
		} else {
			if n != 1 {
				return dgrzerr.UnexpectedValue.Newf("%s: expected to scan 1 relative path, "+
					", scanned %d", headFile, n)
			}
		}

		commitHashFile := path.Join(RepositoriesDirPath(ctxt), repoName, DgrzDirName, headFileSubPath)

		myFunc := func() (io.Reader, error) {
			return strings.NewReader(commitHash), nil
		}

		if _, err := WriteToFileAtomic(myFunc, commitHashFile); err != nil {
			return err
		}

	}
	return nil

}

func RepositoryExist(repoName string, ctxt context.Context) bool {
	// TODO: have this method return bool, error
	repoPath := filepath.Join(RepositoriesDirPath(ctxt), repoName)

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
