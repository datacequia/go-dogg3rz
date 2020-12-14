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

package repo

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	"github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
	cid "github.com/ipfs/go-cid"
)

const (
	//	FILE_INDEX_VERSION = uint32(1)
	HeaderLength   = 12
	ChecksumLength = sha1.Size

//	HeaderSignature    = [4]byte{'R', 'E', 'S', 'C'}
)

type fileRepositoryIndex struct {

	// O/S PATH TO BASE REPOSITORY
	// DIRECTORY
	repoDir string
	// THE REPOSITORY NAME
	repoName string
	// THE O/S PATH TO THE INDEX FILE
	path string
}

/*
type indexEntry struct {
	datasetPath    string
	objectType     jsonld.JSONLDResourceType
	objectID       string // IRI or TERM
	containerType  jsonld.JSONLDResourceType
	containerID    string
	lastModifiedNs int64
	objectCID      string
}
*/

type indexEntryInternal struct {
	datasetPathLength int16
	datasetPath       []byte
	objectType        jsonld.JSONLDResourceType
	objectIDLength    int16
	objectID          []byte
	containerType     jsonld.JSONLDResourceType
	containerIDLength int16
	containerID       []byte
	lastModifiedNs    int64
	objectCIDLength   int8
	objectCID         []byte
}

/*

INDEX FILE BINARY FORMAT

HEADER
  SIGNATURE - 4 BYTES . ALWAYS 'RESC' (RESOURCE CACHE)
  VERSION   - 4 BYTES. INDEX FORMAT VERSION
  NUM ENTRIES - 4 BYTES, NUMBER OF INDEX ENTRIES
ENTRY
  DATASET PATH LENGTH - 2 BYTES
  DATASET PATH (UNIX STYLE) - VARIABLE LENGTH (UP TO 2^16 )
	(JSONLD) OBJECT TYPE - 1 BYTE
	OBJECT ID/IRI LENGTH - 2 BYTES
	(JSONLD) OBJECT ID/IRI - VARIABLE LENGTH (UP TO 2^16)
	(JSONLD) CONTAINER TYPE - 1 BYTE
	(JSONLD) CONTAINER ID/IRI LENGTH - 2 BYTES
	(JSONLD) CONTAINER ID/IRI - VARIABLE LENGTH (UP TO 2^16)

	(JSONLD) OBJECT LAST MODIFIED (EPOCH TIME) NANO - 8 BYTES
	(JSONLD) OBJECT CID LENGTH -  1 BYTE
	(JSONLD) OBJECT CID - VARIABLE LENGTH (UP TO 2^8)

CHECKSUM
  SHA-1 INDEX CHECKSUM - 160 BIT (8 BYTES) OVER CONTENT OF INDEX BEFORE THIS
                         CHECKSUM
*/

func newFileRepositoryIndex(ctxt context.Context, repoName string) (*fileRepositoryIndex, error) {

	index := &fileRepositoryIndex{}

	index.repoName = repoName

	index.repoDir = filepath.Join(file.RepositoriesDirPath(ctxt), repoName)

	index.path = filepath.Join(index.repoDir, file.IndexFileName)

	if !file.DirExists(index.repoDir) {
		return nil, errors.NotFound.Newf("repository directory at %s does not exist",
			index.repoDir)
	}

	return index, nil

}

// Adds a new resource (resId) to the index
func (index *fileRepositoryIndex) update(entry1 common.StagingResource) error {

	var internalEntry indexEntryInternal
	var err error

	if internalEntry, err = ValidateIndexEntry(entry1); err != nil {
		return err
	}
	// CREATE CALLBACK TO WRITE NEW ENTRY

	addFunc := func() (io.Reader, error) {
		// GET EXISTING ENTRIES IN INDEX FILE//

		var indexEntriesInternal []indexEntryInternal

		if ie, err := index.readIndexFileInternal(); err != nil {
			if os.IsNotExist(err) {
				// DOESN'T EXIST YET. ASSUME FOR NOW
				// THAT REPO DIR EXISTS BUT INDEX FILE DOESN'T
				indexEntriesInternal = make([]indexEntryInternal, 0)
			} else {
				// SOME OTHER ERROR OCCURRED OTHER THAN
				// INDEX FILE DOES NOT EXIIT. RETURN THE ERROR
				return nil, err
			}
		} else {
			// INDEX READ SUCCESSFUL
			//	fmt.Println("index read returned entries",len(ie))
			indexEntriesInternal = ie
		}
		var updatedExistingEntry bool = false

		for i, e := range indexEntriesInternal {
			// CONVERT STRING FORMATTEED UUID BEFORE COMPARE
			// AS FORMATTING COULD VARY . COMPARE BYTE[16] ARRAYS
			//	fmt.Printf("entry read type=%d, uuid=%v,fs=%d,mtime=%d,spl=%d",e.Type,e.Uuid,e.FileSize,e.MtimeNs,e.SubpathLength)

			// COMPARE BY OBJECT IDENTIFIERS (IRI) WHICH SHOULD BE UNIQUE
			// FOR A GIVEN DATASET
			if bytes.Equal(e.datasetPath, internalEntry.datasetPath) &&
				bytes.Equal(e.objectID, internalEntry.objectID) {
				indexEntriesInternal[i] = internalEntry
				updatedExistingEntry = true
			}
		}
		if !updatedExistingEntry {
			indexEntriesInternal = append(indexEntriesInternal, internalEntry)

		}

		buf := &bytes.Buffer{}
		///////////////////////////////////
		// WRITE HEADER TO BUFFFER
		//////////////////////////////////
		if err := binary.Write(buf, binary.BigEndian, []byte(file.ResourceCacheSignature)); err != nil {
			return nil, err
		}

		if err := binary.Write(buf, binary.BigEndian, file.IndexFormatVersion); err != nil {
			return nil, err
		}

		var numIndexEntries uint32 = uint32(len(indexEntriesInternal))
		if err := binary.Write(buf, binary.BigEndian, numIndexEntries); err != nil {
			return nil, err
		}
		// WRITE ENTRIES TO BUFFER
		if err := writeIndexEntriesToBuffer(buf, indexEntriesInternal); err != nil {
			return nil, err
		}

		// WRITE SHA1 HASH

		var checkSum [ChecksumLength]byte = sha1.Sum(buf.Bytes())

		if err := binary.Write(buf, binary.BigEndian, checkSum); err != nil {
			return buf, err
		}

		return buf, nil

	}

	//log.Printf("index.path = %s",index.path)

	if _, err := file.WriteToFileAtomic(addFunc, index.path); err != nil {
		return err
	}

	return nil
}

func writeIndexEntriesToBuffer(buf *bytes.Buffer, indexEntries []indexEntryInternal) error {

	for _, e := range indexEntries {

		var writeList []interface{} = []interface{}{
			e.datasetPathLength,
			e.datasetPath,
			e.objectType,
			e.objectIDLength,
			e.objectID,
			e.containerType,
			e.containerIDLength,
			e.containerID,
			e.lastModifiedNs,
			e.objectCIDLength,
			e.objectCID,
		}

		for _, dataElement := range writeList {
			if err := binary.Write(buf, binary.BigEndian, dataElement); err != nil {
				return err
			}
		}

		//	fmt.Printf("writeIndexEntriesToBuffer[%d]=%s\n",i,e.String())
	}

	return nil

}

func (index *fileRepositoryIndex) readIndexFile() ([]common.StagingResource, error) {

	var indexEntriesInternal []indexEntryInternal = []indexEntryInternal{}
	var indexEntries []common.StagingResource = []common.StagingResource{}

	if i, err := index.readIndexFileInternal(); err != nil {
		return indexEntries, err
	} else {
		indexEntriesInternal = i
	}

	indexEntries = make([]common.StagingResource, len(indexEntriesInternal))
	for i, indexEntryInternal := range indexEntriesInternal {
		indexEntries[i] = indexEntryInternal.ToStagingResource()
	}

	return indexEntries, nil

}

func (index *fileRepositoryIndex) readIndexFileInternal() ([]indexEntryInternal, error) {

	indexEntries := []indexEntryInternal{}

	/////////////////////////
	// OPEN INDEX FILE
	////////////////////////

	f, err := os.Open(index.path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// CHECK SIZE OF FILE TO ENSURE IT MEETS
	// EXPEECT MIN SIZE REQUIREMNETS
	minSize := int64(len(file.ResourceCacheSignature) +
		4 + // SIZEOF UINT32 FOR HEADER VERSION
		4 + // SIZEOF UIN32 FOR # OF INDEX ENTRIES
		ChecksumLength) // SIZEOF SHA1 HASH AT END

	var indexPathFileInfo os.FileInfo

	if fileInfo, err := f.Stat(); err != nil {
		return indexEntries, err
	} else {

		indexPathFileInfo = fileInfo

		if fileInfo.Size() < minSize {
			return indexEntries, errors.OutOfRange.Newf(
				"expected index file at %s to be at least %d bytes, actual size is %d bytes",
				index.path, minSize, fileInfo.Size())

		}
	}

	var offset1 int64 = indexPathFileInfo.Size() - int64(ChecksumLength)

	///////////////////////////////////////////////
	// COMPUTE SHA1 CHECKSUM OF INDEX FILE
	//////////////////////////////////////////////
	h := sha1.New()
	if _, err := io.CopyN(h, f, offset1); err != nil {
		return indexEntries, err
	}

	///////////////////////////////////////////////////
	// READ CHECKSUM WRITTEN AT END OF INDEX FILE
	//////////////////////////////////////////////////
	fileChecksum := make([]byte, ChecksumLength)
	if readSize, err := f.Read(fileChecksum); err != nil {
		return indexEntries, err
	} else {
		if readSize != ChecksumLength {
			return indexEntries, errors.UnexpectedValue.Newf(
				"expected to read %d byte index hash at "+
					"end of index file at %s, only was able to read %d bytes",
				ChecksumLength, index.path, readSize)
		} else {

			// COMPARE HASHES
			computedChecksum := h.Sum(nil)
			if !bytes.Equal(computedChecksum, fileChecksum) {
				return indexEntries, errors.UnexpectedValue.Newf(
					"index file at %s checksum verification failed:"+
						" { index file checksum = %x, computed checksum = %x}",
					index.path, fileChecksum, computedChecksum)
			}

		}
	}

	// REWIND FILE POINTER BACK TO BEGINNNING NOW THAT
	// INDEX FILE CHECKSUM VERIFIED
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return indexEntries, err
	}

	// READ HEADER SIGNATURE
	sig := make([]byte, len(file.ResourceCacheSignature))
	if err := binary.Read(f, binary.BigEndian, sig); err != nil {
		return indexEntries, err
	}

	if string(sig) != file.ResourceCacheSignature {
		return indexEntries, errors.UnexpectedValue.Newf(
			"%s: expected header signature '%s' , found '%s'",
			index.path, file.ResourceCacheSignature, string(sig))

	}

	// READ INDEX VERSION
	var indexVersion uint32
	if err := binary.Read(f, binary.BigEndian, &indexVersion); err != nil {
		return indexEntries, err
	}
	if indexVersion != file.IndexFormatVersion {
		return indexEntries, errors.UnexpectedValue.Newf("expected index version %d, found %d",
			file.IndexFormatVersion, indexVersion)
	}

	// READ NUMBER OF INDEX ENTRIES
	var numIndexEntries uint32
	if err := binary.Read(f, binary.BigEndian, &numIndexEntries); err != nil {
		return indexEntries, err
	}

	if numIndexEntries < 1 {
		// EMPTY INDEX. NO ENTRIES...
		return indexEntries, nil
	}

	// READ INDEX ENTRIES INTO MEMORY
	indexEntries = make([]indexEntryInternal, numIndexEntries)
	for i := uint32(0); i < numIndexEntries; i++ {

		var e indexEntryInternal

		// READ DATASETPATH
		if err := binary.Read(f, binary.BigEndian, &e.datasetPathLength); err != nil {
			return indexEntries, err
		}

		e.datasetPath = make([]byte, e.datasetPathLength)
		if err := binary.Read(f, binary.BigEndian, &e.datasetPath); err != nil {
			return indexEntries, err
		}

		// READ OBJECTTYPE
		if err := binary.Read(f, binary.BigEndian, &e.objectType); err != nil {
			return indexEntries, err
		}

		// READ OBJECTID
		if err := binary.Read(f, binary.BigEndian, &e.objectIDLength); err != nil {
			return indexEntries, err
		}

		e.objectID = make([]byte, e.objectIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.objectID); err != nil {
			return indexEntries, err
		}

		// READ CONTAINERTYPE
		if err := binary.Read(f, binary.BigEndian, &e.containerType); err != nil {
			return indexEntries, err
		}

		// READ CONTAINERID
		if err := binary.Read(f, binary.BigEndian, &e.containerIDLength); err != nil {
			return indexEntries, err
		}

		e.containerID = make([]byte, e.containerIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.containerID); err != nil {
			return indexEntries, err
		}

		// READ LASTMODIFIEDNS
		if err := binary.Read(f, binary.BigEndian, &e.lastModifiedNs); err != nil {
			return indexEntries, err
		}

		// READ OBJECTCID
		if err := binary.Read(f, binary.BigEndian, &e.objectCIDLength); err != nil {
			return indexEntries, err
		}

		e.objectCID = make([]byte, e.objectCIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.objectCID); err != nil {
			return indexEntries, err
		}

		indexEntries[i] = e

	}

	return indexEntries, nil

}

/*
func (e indexEntry) String() string {

	return fmt.Sprintf("indexEntry: { datasetPath=%s, objectType = %s, "+
		"objectID=%s, containerType=%s, containerID = %s, lastModifiedNs = %d objectCID = %s }",
		e.datasetPath, e.objectType, e.objectID, e.containerType, e.containerID,
		e.lastModifiedNs, e.objectCID)

}
*/

func ValidateIndexEntry(e common.StagingResource) (indexEntryInternal, error) {

	var internal indexEntryInternal

	// VALIDATE DATASET PATH
	var rp *common.RepositoryPath
	var err error

	if err = e.AssertValid(); err != nil {
		return internal, err
	}

	if rp, err = common.RepositoryPathNew(e.DatasetPath); err != nil {
		return internal, err
	}

	internal.datasetPath = []byte(rp.ToString())

	if len(internal.datasetPath) > math.MaxInt16 {
		return internal, errors.OutOfRange.Newf(
			"datasetPath length too long: { datasetPath length %d, max allowed length %d }",
			len(internal.datasetPath), math.MaxInt16)

	}
	internal.datasetPathLength = int16(len(internal.datasetPath))

	// VALIDATE OBJECT TYPE
	/*
		if err = e.objectType.AssertValid(); err != nil {
			return internal, err
		}
	*/

	internal.objectType = e.ObjectType

	// VALIDATE OBJECT ID
	internal.objectID = []byte(e.ObjectIRI)
	if len(internal.objectID) > math.MaxInt16 {
		return internal, errors.OutOfRange.Newf("objectID length too long: { objectID length %d, max allowed length %d}",
			len(internal.objectID), math.MaxInt16)

	}
	internal.objectIDLength = int16(len(internal.objectID))

	// VALIDATE CONTAINER TYPE
	/*
		if err := e.containerType.AssertValid(); err != nil {
			return internal, err
		}
	*/

	internal.containerType = e.ContainerType

	// 	VALIDATE CONTAINER ID
	internal.containerID = []byte(e.ContainerIRI)
	if len(internal.containerID) > math.MaxInt16 {
		return internal, errors.OutOfRange.Newf("containerID length too long: { containerID length %d, max allowed length %d}",
			len(internal.containerID), math.MaxInt16)

	}
	internal.containerIDLength = int16(len(internal.containerID))

	//  VALIDATE LAST MODIFIED TIME
	if e.LastModifiedNs < 0 {
		// EPOCH TIME SHOULD NEVER BE LESS THAN ZERO
		return internal, errors.OutOfRange.Newf("lastModifiedNs (epoch time) should be between 0 and %d. found %d",
			math.MaxInt64, e.LastModifiedNs)
	}
	internal.lastModifiedNs = e.LastModifiedNs

	// VALIDATE CID (IF StagingResource Wants CID assigned)
	if e.WantsCID() {
		if _, err = cid.Decode(e.ObjectCID); err != nil {
			return internal, errors.InvalidValue.Wrapf(err, ""+
				"bad content identifier value '%s' assigned to objectCID", e.ObjectCID)

		}
	}

	internal.objectCID = []byte(e.ObjectCID)
	if len(internal.objectCID) > math.MaxInt8 {
		return internal, errors.OutOfRange.Newf("objectCID length too long: { objectCID length %d, max allowed length %d}",
			len(internal.objectCID), math.MaxInt8)
	}
	internal.objectCIDLength = int8(len(internal.objectCID))

	return internal, nil

}

func (i indexEntryInternal) ToStagingResource() common.StagingResource {

	var e common.StagingResource

	e.DatasetPath = string(i.datasetPath)
	e.ObjectType = i.objectType
	e.ObjectIRI = string(i.objectID)
	e.ContainerType = i.containerType
	e.ContainerIRI = string(i.containerID)
	e.LastModifiedNs = i.lastModifiedNs
	e.ObjectCID = string(i.objectCID)

	return e
}

func (i indexEntryInternal) String() string {

	return fmt.Sprintf("indexEntryInternal = { DatasetPathLength=%d, DatasetPath=%s, "+
		"objectType=%s, objectIDLength=%d, objectID=%s,"+
		" containerType=%s, containerIDLength=%d, containerID=%s"+
		" lastModifiedNs=%d, objectCIDLength=%d, objectCID=%s }",
		i.datasetPathLength, string(i.datasetPath),
		i.objectType,
		i.objectIDLength, string(i.objectID),
		i.containerType,
		i.containerIDLength, string(i.containerID),
		i.lastModifiedNs,
		i.objectCIDLength, string(i.objectCID))

}
