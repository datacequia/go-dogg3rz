package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
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
const IndexFileName = "index"

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
		err = os.Remove(lockFile)
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
