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

package grapp

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"

	"math"
	"path"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
)

const indexFileName = "index"

type FileGrapplicationResource struct {
}

func (grapp *FileGrapplicationResource) Create(ctxt context.Context, grappDir string) error {

	// check to see if grapp dir exists. if not try to create it
	if !file.DirExists(grappDir) {
		if err := os.MkdirAll(grappDir, 0750); err != nil {
			return err
		}

	}

	//grappDir := path.Join(file.GrapplicationsDirPath(ctxt), path)
	//grappDir := path

	// CREATE DEFAULT BRANCH DIR
	mainBranchDir := path.Join(grappDir, file.MasterBranchName)

	// CREATE 'refs/heads' SUBDIR
	// CREATE '.dgrz' DIR AS SUBDI OF BASE GRAPP DIR
	dgrzDir := path.Join(grappDir, file.DgrzDirName)
	refsDir := path.Join(dgrzDir, file.RefsDirName)
	headsDir := path.Join(refsDir, file.HeadsDirName)

	// CREATE GRAPP DIRS IN THE FOLLOWING ORDER
	dirsList := []string{mainBranchDir, dgrzDir, refsDir, headsDir}

	for _, d := range dirsList {

		err := os.Mkdir(d, os.FileMode(0700))

		if err != nil {
			if os.IsNotExist(err) {
				// BASE GRAPP DIR DOES NOT EXIST
				return dgrzerr.NotFound.Wrapf(err, grappDir)
			}
			if os.IsExist(err) {

				return dgrzerr.AlreadyExists.Wrapf(err, grappDir)
			}

			return err

		}

	}
	// WRITE THE HEAD FILE WITH A POINTER TO DEFAULT MASTER BRANCH
	err := file.WriteHeadFile(ctxt, grappDir, file.MasterBranchName)
	if err != nil {
		return err
	}

	// CREATE IPFS CONTAINER FOR THIS GRAPPLICATION
	_, err = allocateIPFSAPIPort(ctxt, grappDir)

	return err

}

func (grapp *FileGrapplicationResource) Add(ctxt context.Context, grappName string, path string) error {

	return dgrzerr.GrappError.New("not implemented") //add(ctxt, grappName, path)

}

func allocateIPFSAPIPort(ctxt context.Context, dirPath string) (int, error) {

	return nextIPFSAPIPort(ctxt, 50001, dirPath)
}

// Allocate next available IPFS API Listen Port based on counter file
// in dirPath
func nextIPFSAPIPort(ctxt context.Context, basePortNum int, dirPath string) (int, error) {

	if !file.DirExists(dirPath) {
		return -1, dgrzerr.InvalidValue.Newf("%s is not a directory or does not exist", dirPath)
	}

	baseGrappDirPath := file.DataDirPath(ctxt) //file.GrapplicationsDirPath(ctxt)
	//grappDirPath := path.Join(file.GrapplicationsDirPath(ctxt), grappDirName)
	counterFilePath := path.Join(baseGrappDirPath, file.IPFSAPIPortCounterFileName)

	var nextPort int = -1

	const minTCPPort = 1025
	const maxTCPPort = int(math.MaxUint16) // TCP PORTS ARE unsigned 16 bit integers

	if basePortNum > maxTCPPort || basePortNum < minTCPPort {
		// supplied base port is outside of user tcp port range
		return -1, dgrzerr.OutOfRange.Newf("supplied base port number '%d' is outside of allowable tcp port range %d->%d",
			basePortNum, minTCPPort, maxTCPPort)

	}

	counterFunc := func() (io.Reader, error) {
		// READ NEXT PORT VALUE FROM PORT COUNTER FILE

		nextPort = basePortNum

		if file.FileExists(counterFilePath) {
			// COUNTER FILE EXISTS. READ CURRENT VALUE

			nextPortAsBytes, err := os.ReadFile(counterFilePath)

			if err != nil {
				return nil, err
			}

			if nextPort, err = strconv.Atoi(string(nextPortAsBytes)); err != nil {
				return nil, err
			}
			nextPort++ // INCREMENT FROM EXISTING VALUE
		}

		// INCREMENT PORT COUNTER AND WRITE TO COUNTER FILE
		return strings.NewReader(strconv.Itoa(nextPort)), nil

	}

	if _, err := file.WriteToFileAtomic(counterFunc, counterFilePath); err != nil {
		return -1, err
	}

	if nextPort < minTCPPort {
		// never initialized. programmer error
		panic("nextPort never assigned or assigned out of range min value")
	}

	if nextPort > maxTCPPort {
		return -1, dgrzerr.OutOfRange.Newf("exhausted port range %d-%d looking for unused IPFS API port",
			basePortNum, maxTCPPort)
	}

	return nextPort, nil

}
