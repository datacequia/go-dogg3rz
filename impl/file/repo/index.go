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
	"time"

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

	channelWriteTimeoutSeconds = 5

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

	// CHANNEL USED BY WORKER GO-ROUTINE TO RECV REQUESTS
	requestChannel chan msgIndexRequest
}

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

type fileRepositoryIndexWorker struct {
	index fileRepositoryIndex

	// CACHED ENTRIES
	indexCache []common.StagingResource

	// TIMESTAMP THAT TRACKS LAST TIME INDEX FILE WAS LOADED INTO CACHE
	indexCacheLastUpdated time.Time

	// COUNTER THAT INCREMENTS EACH TIME AN OPERATION MODIFIES CACHE (DIRTY CACHE)
	modCount int
}

type msgIndexRequest struct {
	value interface{} // MUST BE ASSIGNED ONE OF msg[X]Request structs below
}

type msgIndexResponse struct {
	request interface{}
	err     error
}

// REQUEST TYPES

// REQUEST TYPE ONLY RETURNS A RESPONSE VIA SINGLE CHANNEL
type msgIndexBaseRequest struct {
	responseChan chan msgIndexResponse
}

// REQUEST TYPE RETURNS A RESPONSE AND A RESULT SET OF StagingResource INDEX ENTRIES
type msgIndexResultSetRequest struct {
	msgIndexBaseRequest
	resultSetChan chan common.StagingResource
}

type msgIndexScanRequest struct {
	msgIndexResultSetRequest
	filterFunc func(common.StagingResource) bool
}

type msgIndexStageRequest struct {
	msgIndexBaseRequest
	value common.StagingResource
}

type msgIndexRemoveRequest struct {
	msgIndexResultSetRequest
	value common.StagingResource
}

type msgIndexCommitRequest struct {
	msgIndexBaseRequest
}

type msgIndexRollbackRequest struct {
	msgIndexBaseRequest
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

	index := fileRepositoryIndex{}

	index.repoName = repoName

	index.repoDir = filepath.Join(file.RepositoriesDirPath(ctxt), repoName)

	index.path = file.IndexFilePath(ctxt, repoName)

	if !file.DirExists(index.repoDir) {
		return nil, errors.NotFound.Newf("repository directory at %s does not exist",
			index.repoDir)
	}

	// CREATE REQUEST CHANNEL
	index.requestChannel = make(chan msgIndexRequest, 1)

	//  CREATE INDEX WORKER AND SPAWN IT

	worker, err := newFileRepositoryIndexWorker(index)
	if err != nil {
		return nil, err

	}

	go worker.run(ctxt)

	return &index, nil

}

// close closes the request channel, which releases resources allocated
// by fileRepositoryIndex
//
// HINT: You probably want to call  defer <index object>close()
//       after calling newFileRepositoryIndex()
func (index *fileRepositoryIndex) close() {
	// CLOSE THE REQUEST CHANNEL. THIS WILL RELEASE THE WORKER
	close(index.requestChannel)
	// RESET requestChannel to nil to indicate it's cloosed
	index.requestChannel = nil

}

// scan scans all entries in the index file and returns matching entries.
// Returns two channels via msgIndexResultSetRequest members that must be read by the caller until they are closed.
// The common.StagingResource channel will return those index entries that cause the filterFunc
// parameter to return true when invoked.
// The msgIndexResponse channel will  return a single response upon read by the caller
// which delivers the final outcome of this call.
func (index *fileRepositoryIndex) scan(filterFunc func(common.StagingResource) bool) (*msgIndexResultSetRequest, error) {

	if index.requestChannel == nil {
		return nil, errors.ChannelClosed.New("request channel closed")
	}

	// construct request msg
	scanReq := msgIndexScanRequest{}
	scanReq.responseChan = make(chan msgIndexResponse)
	scanReq.resultSetChan = make(chan common.StagingResource)
	scanReq.filterFunc = filterFunc

	request := msgIndexRequest{}
	request.value = scanReq

	// send the request
	index.requestChannel <- request
	//fmt.Println("sent scan request to req chan")
	// return response channels to caller for read
	return &scanReq.msgIndexResultSetRequest, nil

}

// stage stages a single resource to the index file. Returns a single channel via
// msgIndexBaseRequest member
// that must be read until it is closed. This channel upon a single read will
// return the final outcome of this call.
// Note: A subsequent call to either commit or rollback must occur so that changes
// will take effect or be undone. Failure to call either will result in
// index file being locked indefinitely
func (index *fileRepositoryIndex) stage(resource common.StagingResource) error {

	if index.requestChannel == nil {
		return errors.ChannelClosed.New("request channel closed")
	}

	stageReq := msgIndexStageRequest{}
	stageReq.responseChan = make(chan msgIndexResponse)
	stageReq.value = resource

	request := msgIndexRequest{}
	request.value = stageReq

	index.requestChannel <- request

	var err error

	select {
	case response := <-stageReq.responseChan:
		err = response.err
	case <-time.After(time.Second * channelWriteTimeoutSeconds):
		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to receive response message: { request = %s }",
			channelWriteTimeoutSeconds, request)
	}

	return err

}

// remove removes  any resources from the index file that match filterFunc where
// an index entry passed to filterFunc returns true.  Returns two channels via
// msgIndexResultSetRequest return value which
// must be read by the caller until they closed.
// The first channel  will return those matching resources upon read
// by the caller.
// The second channel will  return a single response upon read by the caller
// which delivers the final outcome of this call.
//
// Note: A subsequent call to either commit or rollback MUST occur so that changes
// will take effect or be undone. Failure to call either will result in
// unreleased resources.
func (index *fileRepositoryIndex) remove(resource common.StagingResource) (*msgIndexResultSetRequest, error) {

	if index.requestChannel == nil {
		return nil, errors.ChannelClosed.New("request channel closed")
	}

	removeReq := msgIndexRemoveRequest{}
	removeReq.value = resource
	removeReq.responseChan = make(chan msgIndexResponse)
	removeReq.resultSetChan = make(chan common.StagingResource)

	request := msgIndexRequest{}
	request.value = removeReq

	index.requestChannel <- request

	//fmt.Println("sent scan request to req chan")
	// return response channels to caller for read
	return &removeReq.msgIndexResultSetRequest, nil

}

// commit commits pending changes made from prior calls to stage or remove
// from the receiver object.
// Returns a single channel via msgIndexBaseRequest member  that is subsequently
// read by caller to retrieve a single
// value of type msgIndexResponse that provides the success or failure status of
// this call.
func (index *fileRepositoryIndex) commit() error {

	if index.requestChannel == nil {
		return errors.ChannelClosed.New("request channel closed")
	}

	commitReq := msgIndexCommitRequest{}
	commitReq.responseChan = make(chan msgIndexResponse)

	request := msgIndexRequest{}
	request.value = commitReq

	// SEND REQUEST
	index.requestChannel <- request

	var err error

	select {
	case response := <-commitReq.responseChan:
		err = response.err
	case <-time.After(time.Second * channelWriteTimeoutSeconds):
		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to receive response message: { request = %s }",
			channelWriteTimeoutSeconds, request)
	}

	return err

}

// rollback rolls back pending changes made from prior calls to stage or remove
// from the receiver object.
// Returns a channel that is subsequently read by caller to retrieve a single
// value of type msgIndexResponse that determines the success or failure status of
// this call.
func (index *fileRepositoryIndex) rollback() error {

	if index.requestChannel == nil {
		return errors.ChannelClosed.New("request channel closed")
	}

	rollbackReq := msgIndexRollbackRequest{}
	rollbackReq.responseChan = make(chan msgIndexResponse)

	request := msgIndexRequest{}
	request.value = rollbackReq

	// SEND REQUEST
	index.requestChannel <- request

	var err error

	select {
	case response := <-rollbackReq.responseChan:
		err = response.err
	case <-time.After(time.Second * channelWriteTimeoutSeconds):
		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to receive response message: { request = %s }",
			channelWriteTimeoutSeconds, request)
	}

	return err

}

// Adds a new resource (resId) to the index

func writeIndexToBuffer(entries []common.StagingResource) (*bytes.Buffer, error) {

	var internalEntry indexEntryInternal
	var err error

	// CONVERT StaginResources array to internal entries
	indexEntriesInternal := make([]indexEntryInternal, len(entries))
	for i, e := range entries {
		if internalEntry, err = ValidateIndexEntry(e); err != nil {
			return nil, err
		}

		indexEntriesInternal[i] = internalEntry
	}

	buf := &bytes.Buffer{}

	///////////////////////////////////
	// WRITE HEADER TO BUFFFER
	//////////////////////////////////

	if err := binary.Write(buf, binary.BigEndian, []byte(file.ResourceCacheSignature)); err != nil {
		return nil, err
	}

	//binary.Size

	if err := binary.Write(buf, binary.BigEndian, file.IndexFormatVersion); err != nil {
		return nil, err
	}

	var numIndexEntries uint32 = uint32(len(indexEntriesInternal))
	if err := binary.Write(buf, binary.BigEndian, numIndexEntries); err != nil {
		return nil, err
	}

	// WRITE ENTRIES TO BUFFER
	if err := writeIndexEntriesInternalToBuffer(buf, indexEntriesInternal); err != nil {
		return nil, err
	}

	// WRITE SHA1 HASH

	var checkSum [ChecksumLength]byte = sha1.Sum(buf.Bytes())

	if err := binary.Write(buf, binary.BigEndian, checkSum); err != nil {
		return buf, err
	}

	return buf, nil

}

func writeIndexEntriesInternalToBuffer(buf *bytes.Buffer, indexEntries []indexEntryInternal) error {

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

func (index *fileRepositoryIndex) readIndexFile() ([]common.StagingResource, os.FileInfo, error) {

	var indexEntriesInternal []indexEntryInternal = []indexEntryInternal{}
	var indexEntries []common.StagingResource = []common.StagingResource{}

	var fileInfo os.FileInfo

	if i, fInfo, err := index.readIndexFileInternal(); err != nil {
		return indexEntries, fInfo, err
	} else {
		indexEntriesInternal = i
		fileInfo = fInfo
	}

	indexEntries = make([]common.StagingResource, len(indexEntriesInternal))
	for i, indexEntryInternal := range indexEntriesInternal {
		indexEntries[i] = indexEntryInternal.ToStagingResource()
	}

	return indexEntries, fileInfo, nil

}

func (index *fileRepositoryIndex) readIndexFileInternal() ([]indexEntryInternal, os.FileInfo, error) {

	indexEntries := []indexEntryInternal{}

	/////////////////////////
	// OPEN INDEX FILE
	////////////////////////

	f, err := os.Open(index.path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var fileInfo os.FileInfo

	fileInfo, err = f.Stat()
	if err != nil {
		return nil, nil, err
	}

	// CHECK SIZE OF FILE TO ENSURE IT MEETS
	// EXPEECT MIN SIZE REQUIREMNETS
	minSize := int64(len(file.ResourceCacheSignature) +
		4 + // SIZEOF UINT32 FOR HEADER VERSION
		4 + // SIZEOF UIN32 FOR # OF INDEX ENTRIES
		ChecksumLength) // SIZEOF SHA1 HASH AT END

	var indexPathFileInfo os.FileInfo

	if fileInfo, err := f.Stat(); err != nil {
		return indexEntries, fileInfo, err
	} else {

		indexPathFileInfo = fileInfo

		if fileInfo.Size() < minSize {
			return indexEntries, fileInfo, errors.OutOfRange.Newf(
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
		return indexEntries, fileInfo, err
	}

	///////////////////////////////////////////////////
	// READ CHECKSUM WRITTEN AT END OF INDEX FILE
	//////////////////////////////////////////////////
	fileChecksum := make([]byte, ChecksumLength)
	if readSize, err := f.Read(fileChecksum); err != nil {
		return indexEntries, fileInfo, err
	} else {
		if readSize != ChecksumLength {
			return indexEntries, fileInfo, errors.UnexpectedValue.Newf(
				"expected to read %d byte index hash at "+
					"end of index file at %s, only was able to read %d bytes",
				ChecksumLength, index.path, readSize)
		} else {

			// COMPARE HASHES
			computedChecksum := h.Sum(nil)
			if !bytes.Equal(computedChecksum, fileChecksum) {
				return indexEntries, fileInfo, errors.UnexpectedValue.Newf(
					"index file at %s checksum verification failed:"+
						" { index file checksum = %x, computed checksum = %x}",
					index.path, fileChecksum, computedChecksum)
			}

		}
	}

	// REWIND FILE POINTER BACK TO BEGINNNING NOW THAT
	// INDEX FILE CHECKSUM VERIFIED
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return indexEntries, fileInfo, err
	}

	// READ HEADER SIGNATURE
	sig := make([]byte, len(file.ResourceCacheSignature))
	if err := binary.Read(f, binary.BigEndian, sig); err != nil {
		return indexEntries, fileInfo, err
	}

	if string(sig) != file.ResourceCacheSignature {
		return indexEntries, fileInfo, errors.UnexpectedValue.Newf(
			"%s: expected header signature '%s' , found '%s'",
			index.path, file.ResourceCacheSignature, string(sig))

	}

	// READ INDEX VERSION
	var indexVersion uint32
	if err := binary.Read(f, binary.BigEndian, &indexVersion); err != nil {
		return indexEntries, fileInfo, err
	}
	if indexVersion != file.IndexFormatVersion {
		return indexEntries, fileInfo, errors.UnexpectedValue.Newf("expected index version %d, found %d",
			file.IndexFormatVersion, indexVersion)
	}

	// READ NUMBER OF INDEX ENTRIES
	var numIndexEntries uint32
	if err := binary.Read(f, binary.BigEndian, &numIndexEntries); err != nil {
		return indexEntries, fileInfo, err
	}

	if numIndexEntries < 1 {
		// EMPTY INDEX. NO ENTRIES...
		return indexEntries, fileInfo, nil
	}

	// READ INDEX ENTRIES INTO MEMORY
	indexEntries = make([]indexEntryInternal, numIndexEntries)
	for i := uint32(0); i < numIndexEntries; i++ {

		var e indexEntryInternal

		// READ DATASETPATH
		if err := binary.Read(f, binary.BigEndian, &e.datasetPathLength); err != nil {
			return indexEntries, fileInfo, err
		}

		e.datasetPath = make([]byte, e.datasetPathLength)
		if err := binary.Read(f, binary.BigEndian, &e.datasetPath); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ OBJECTTYPE
		if err := binary.Read(f, binary.BigEndian, &e.objectType); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ OBJECTID
		if err := binary.Read(f, binary.BigEndian, &e.objectIDLength); err != nil {
			return indexEntries, fileInfo, err
		}

		e.objectID = make([]byte, e.objectIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.objectID); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ CONTAINERTYPE
		if err := binary.Read(f, binary.BigEndian, &e.containerType); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ CONTAINERID
		if err := binary.Read(f, binary.BigEndian, &e.containerIDLength); err != nil {
			return indexEntries, fileInfo, err
		}

		e.containerID = make([]byte, e.containerIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.containerID); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ LASTMODIFIEDNS
		if err := binary.Read(f, binary.BigEndian, &e.lastModifiedNs); err != nil {
			return indexEntries, fileInfo, err
		}

		// READ OBJECTCID
		if err := binary.Read(f, binary.BigEndian, &e.objectCIDLength); err != nil {
			return indexEntries, fileInfo, err
		}

		e.objectCID = make([]byte, e.objectCIDLength)
		if err := binary.Read(f, binary.BigEndian, &e.objectCID); err != nil {
			return indexEntries, fileInfo, err
		}

		indexEntries[i] = e

	}

	return indexEntries, fileInfo, nil

}

func (index *fileRepositoryIndex) String() string {
	return fmt.Sprintf("fileRepositoryIndex { path: %s, repoDir: %s, repoName: %s, requestChannel: %T at %v }",
		index.path, index.repoDir, index.repoName, index.requestChannel, index.requestChannel)
}

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

///////////////////////////////////////////////
// fileRepositoryIndexWorker functions
//////////////////////////////////////////////

// newFileRepositoryIndexWorker creates a new instance
func newFileRepositoryIndexWorker(index fileRepositoryIndex) (*fileRepositoryIndexWorker, error) {

	worker := fileRepositoryIndexWorker{}

	worker.index = index

	// CREATE EMPTY CACHE
	worker.indexCache = make([]common.StagingResource, 0)

	return &worker, nil

}

func (worker *fileRepositoryIndexWorker) run(ctxt context.Context) {

	for ok := worker.handleRequest(ctxt); ok; ok = worker.handleRequest(ctxt) {

	}

}

// run is the entry point for the worker. This should be invoke by go keyword
// so that it's run in its own execution context
func (worker *fileRepositoryIndexWorker) handleRequest(ctxt context.Context) bool {

	//var ok bool

	// Returns true if msg is of type msgIndexScanRequest
	// which is the only read request operation for the
	// index
	isReadRequest := func(msg msgIndexRequest) bool {
		_, ok := msg.value.(msgIndexScanRequest)

		return ok
	}

	// Returns a pointer to a request or nil if context was cancelled
	// or request channel closed
	getMsg := func() (*msgIndexRequest, error) {

		select {
		case <-ctxt.Done():
			// CONTEXT CANCELLED
			return nil, errors.Cancelled.Wrapf(ctxt.Err(),
				"worker received cancellation event: %s", worker)
		case msg, ok := <-worker.index.requestChannel:
			if ok {
				//fmt.Printf("received request on channel %p: %s\n", worker.index.requestChannel, msg)
				return &msg, nil
			}
			// REQUEST CHANNEL WAS CLOSED
			// RETURN TO CALLER
			return nil, errors.ChannelClosed.Newf("request channel: %s", worker)
		}
	} // end getMsg

	var pMsg *msgIndexRequest
	var err error

	// Keep reading index request messages while
	// they are read-only operations
	for pMsg, err = getMsg(); pMsg != nil && isReadRequest(*pMsg); pMsg, err = getMsg() {

		switch v := pMsg.value.(type) {
		case msgIndexScanRequest:
			worker.handleIndexScanRequest(ctxt, v)
		default:
			// IF THIS PANICS, THEN DIDN'T HANDLE ANY OTHER TYPE OF READ
			// REQUESTS THAT WERE DEFINED in isReadRequest
			panic(fmt.Sprintf("msg type not handled in isReadRequest(): %T", v))
		}

	}

	if pMsg == nil {
		// context cancelled or request channel closed BEFORE
		// any non-read requests were issued

		if err != nil &&
			(errors.GetType(err) == errors.ChannelClosed ||
				errors.GetType(err) == errors.Cancelled) {
			return false // EXIT EXECUTION CONTEXT
		}
		// TODO: log warning msg here before exit
		//fmt.Println("here ", err)
		panic("how count err not be one of above and pMsg is still nil?")
	}

	// THIS CALLBACK IS CALLED WHEN ATOMICALLY WRITING INDEX FILE TO DISK
	// IT PROVIDES THE CONTENTS OFTHE INDEX FILE. IS CALLED IMMEDIATELY
	// FOLLOWING THIS DECLARATION
	txFunc := func() (io.Reader, error) {

		// INDEX FILE IS LOCKED WHILE THIS FUNC IS RUNNING

		// IF CURRENT pMsg VALUE IS A NON READ REQUEST THE FOLLOWING
		// LOOP WILL RUN UNTIL COMMIT OR ROLLBACK MESSAGE IS RECEIVED
		// , CONTEXT IS CANCELLED, OR REQUEST CHANNEL IS CLOSED, WHICHEVER
		// COMES FIRST
		for ; pMsg != nil; pMsg, err = getMsg() {

			switch v := pMsg.value.(type) {
			case msgIndexScanRequest:
				worker.handleIndexScanRequest(ctxt, v)

			case msgIndexStageRequest:
				worker.handleIndexStageRequest(ctxt, v)

			case msgIndexRemoveRequest:
				worker.handleIndexRemoveRequest(ctxt, v)

			case msgIndexCommitRequest:

				// CHECK TO SEE IF CACHE HAS BEEN MODIFIED AT ALL
				if !worker.hasDirtyCache() {
					// NOTHING TO COMMIT
					response := msgIndexResponse{}
					response.request = v
					response.err = errors.EmptyCommit.Newf("nothing to commit. worker.modCount equals %d", worker.modCount)

					_ = v.sendResponse(response)

					return nil, response.err
				}

				// WRITE INDEX CACHE TO IO BUFFER
				// AND RETURN TO CALLER.
				var buf *bytes.Buffer

				if buf, err = writeIndexToBuffer(worker.indexCache); err != nil {
					return nil, err
				}

				// reset modCount
				worker.modCount = 0

				return buf, nil // INDEX ENTRIES WRITTEN TO BUFFER (io.Reader)

			case msgIndexRollbackRequest:
				worker.handleIndexRollbackRequest(ctxt, v)
				// NEED TO RETURN AN ERROR IN THIS CALLBACK SO THAT INDEX WRITE
				// IS ABORTED UPON RECEIVING ROLLBACK REQUEST MSG
				return nil, errors.RollbackRequested.Newf(v.String())
			default:
				// CALLER SENT UNKNOWN MSG TYPE TO REQUEST CHANNEL
				panic(fmt.Sprintf("unknown msg type (%T): expected one of the following types: %s,%s,%s,%s,%s",
					v,
					"msgIndexScanRequest", "msgIndexStageRequest", "msgIndexRemoveRequest",
					"msgIndexCommitRequest", "msgIndexRollbackRequest"))

			}
		}
		// CONTEXT CANCELLED OR REQUEST CHANNEL CLOSED
		//  BEFORE COMMIT OR ROLLBACK
		return nil, err

	} // end txFunc

	if _, err = file.WriteToFileAtomic(txFunc, worker.index.path); err != nil {

		switch errors.GetType(err) {
		case errors.Cancelled, errors.ChannelClosed:
			// NOT AN ERROR JUST A NON COMMIT OUTCOME
			return false // EXIT EXECUTION CONTEXT
		case errors.RollbackRequested, errors.EmptyCommit:
			return true // HANDLE ANOTHR REQUEST
		case errors.TryAgain:
			// could not get lock on index file !
			switch v := pMsg.value.(type) {
			case msgIndexStageRequest:
				// IF ERROR WAS TRY AGAIN, THE ASSUMPTION WAS THAT
				// MODIFY OPERATION COOULD NOT GET LOCK ON INDEX FILE

				response := msgIndexResponse{}
				response.err = err
				response.request = pMsg
				v.sendResponse(response)
				return true // HANDLE ANOTHER REQUEST
			case msgIndexRemoveRequest:
				// IF ERROR WAS TRY AGAIN, THE ASSUMPTION WAS THAT
				// MODIFY OPERATION COOULD NOT GET LOCK ON INDEX FILE

				response := msgIndexResponse{}
				response.err = err
				response.request = pMsg
				v.sendResponse(response)
				return true // HANDLE ANOTHER REQUEST

			default:
				// if here then there was likely an unhandled modify
				// request msg type
				panic(fmt.Sprintf("expected an index modify request type, found %T", pMsg.value))
			}
			return true
		}

		// TODO: SHOULD LOG UNHANDLED ERROR MSG HERE

	}

	if pMsg == nil {
		panic("pMsg is nil")
	}

	// if here pMsg should be a commit msg
	var req msgIndexCommitRequest
	var ok bool

	if req, ok = pMsg.value.(msgIndexCommitRequest); !ok {
		panic(fmt.Sprintf("expected msgIndexBaseRequest, found %T: %s", pMsg.value, err))
	}

	// send back response for sucessful commit
	response := msgIndexResponse{}
	response.err = err
	response.request = pMsg
	req.sendResponse(response)

	return true // HANDLE ANOTHER REQUEST

}

// String implements the fmt.Stringer interface to represent
// fileRepositoryIndexWorker  as a string
func (worker *fileRepositoryIndexWorker) String() string {
	return fmt.Sprintf(
		"fileRepositoryIndexWorker { index: %v, indexCacheLastUpdated: %s, "+
			"len(worker.indexCache): %d, modCount: %d }", worker.index,
		worker.indexCacheLastUpdated,
		len(worker.indexCache), worker.modCount)

}

// handleIndexScanRequest processes a msgIndexScanRequest message
func (worker *fileRepositoryIndexWorker) handleIndexScanRequest(ctxt context.Context, request msgIndexScanRequest) {

	// ONLY UPDATE CACHE IF IT HAS NOT BEEN MODIFIED AND UNDERLYING INDEX FILE
	// HAS CHANGED SINCE LAST CACHE REFRESH
	if !worker.hasDirtyCache() {
		//fmt.Println("dirty cache ")
		if err := worker.updateIndexCacheIfInvalidated(); err != nil {
			//fmt.Println("invalidate failed ", err)
			// CAN'T STAT INDEX FILE.
			// SEND BACK ERROR
			response := msgIndexResponse{}
			response.request = request

			// DID CACHE CHECK FAIL BECAUSE INDEX FILE DOES NOT EXIST?
			if pathErr, ok := err.(*os.PathError); ok {

				// COMPARE PATH THAT DOESN'T EXIST WITH INDEX FILE PATH
				// (ASSUMING THAT BOTH ARE ABSOLUTE PATHS)
				if pathErr.Path == worker.index.path {
					response.err = nil // NO INDEX FILE, NOTHING TO SCAN.
				} else {
					response.err = err
				}
			} else {
				response.err = err
			}
			_ = request.sendResponse(response)

			// TODO: IF TIMEOUT OCCURS SHOULD LOG A WARNING FROM RETURN VALUE ABOVE
			return
		}
		//	fmt.Println("updated invalidated cache  ")
	}

	//fmt.Println("indexCache len", len(worker.indexCache))

	for _, result := range worker.indexCache {
		//fmt.Println("iterate worker.indexCache: ", result, request.filterFunc(result))
		//fmt.Println("range worker.indexCache", result)
		if request.filterFunc(result) {
			// MATCHES FILTER FUNC. SEND BACK
			if err := request.sendResult(result); err != nil {
				//fmt.Println("error sending result", err)
				_ = request.sendResponse(msgIndexResponse{request: request, err: err})
				return
			}
			//fmt.Println("sent result", result)
		}
	}
	//fmt.Println("iterate done!")

	// SEND (SUCDESS) RESPONSE BACK
	_ = request.sendResponse(msgIndexResponse{request: request, err: nil})

}

func (worker *fileRepositoryIndexWorker) hasDirtyCache() bool {
	return worker.modCount > 0

}

func (worker *fileRepositoryIndexWorker) handleIndexStageRequest(ctxt context.Context, request msgIndexStageRequest) {

	// VALIDATE REQUEST VALUE
	if _, err := ValidateIndexEntry(request.value); err != nil {
		response := msgIndexResponse{}
		response.err = err

		_ = request.sendResponse(response)
		// TODO: IF TIMEOUT OCCURS SHOULD LOG A WARNING FROM RETURN VALUE ABOVE
		return
	}

	// ONLY UPDATE CACHE IF IT HAS NOT BEEN MODIFIED AND UNDERLYING INDEX FILE
	// HAS CHANGED SINCE LAST CACHE REFRESH
	if !worker.hasDirtyCache() {
		if err := worker.updateIndexCacheIfInvalidated(); err != nil {
			// CAN'T STAT INDEX FILE.
			// SEND BACK ERROR
			response := msgIndexResponse{}
			response.err = err

			_ = request.sendResponse(response)
			// TODO: IF TIMEOUT OCCURS SHOULD LOG A WARNING FROM RETURN VALUE ABOVE
			return
		}
	}

	oldModCount := worker.modCount

	//fmt.Println("update (oldModeCount)", oldModCount)

	//fmt.Println("before worker.indexCache ", worker.indexCache)

	for i, item := range worker.indexCache {

		// IF CURRENT ITEM AND REQUEST HAVE SAME LOCATION, UPDATE CACHE ENTRY
		if item.StagingResourceLocation == request.value.StagingResourceLocation {
			//fmt.Println("cmp", item.StagingResourceLocation, request.value.StagingResourceLocation)
			//	fmt.Println("update index.workerCache ")
			worker.indexCache[i] = request.value
			worker.modCount++

			//fmt.Println("update ", request.value, i, worker.modCount)
			break
		}
	}

	if worker.modCount == oldModCount {
		//fmt.Println("append index.workerCache ")
		// NOT IN CACHE. ADD IT TO CACHE
		worker.indexCache = append(worker.indexCache, request.value)
		//fmt.Println("append ", request.value)
		worker.modCount++
	}
	//fmt.Println("worker.modCount after stage: ", worker.modCount)

	// SEND (SUCDESS) RESPONSE BACK
	_ = request.sendResponse(msgIndexResponse{request: request, err: nil})

}

func (worker *fileRepositoryIndexWorker) handleIndexRemoveRequest(ctxt context.Context, request msgIndexRemoveRequest) {
	// ONLY UPDATE CACHE IF IT HAS NOT BEEN MODIFIED AND UNDERLYING INDEX FILE
	// HAS CHANGED SINCE LAST CACHE REFRESH
	if !worker.hasDirtyCache() {
		if err := worker.updateIndexCacheIfInvalidated(); err != nil {
			// CAN'T STAT INDEX FILE.
			// SEND BACK ERROR
			response := msgIndexResponse{}
			response.err = err

			_ = request.sendResponse(response)
			// TODO: IF TIMEOUT OCCURS SHOULD LOG A WARNING FROM RETURN VALUE ABOVE
			return
		}
	}

	childrenContainers := worker.getChildrenContainers(request.value)
	//fmt.Printf("children containers of %s: %s\n",request.value,childrenContainers)

	//childContainers := worker.getChildContainers(request.value)
	//fmt.Println("childContainers:",childContainers)
	childOfRequest := func(sr common.StagingResource) bool {

		//	childContainers := worker.getChildContainers(sr)
		for _, child := range childrenContainers {

			if sr == child {
				//	  sr.ContainerType == container.ObjectType)  {
				// CURRENT STAGING RESOURCE IS EITHER
				// 1. WITHIN THE LIST OF CHILD CONTAINERS WITHIN request.value
				// OR ...
				// 2. A CHILD OF ONE OF THE CONTAINERS IN TH LIST OF CHILD CONTAINERS
				return true
			}
		}

		return false

	}

	delCount := 0
	// CREATE A NEW CACHE SAME SIZE AS EXISTING CACHE
	// TO HOLD ANY ITEMS NOT DELETED
	newCache := make([]common.StagingResource, len(worker.indexCache))

	// DELETE REQUESTED ENTRY OR A CHILD OF REQUESTED ENTRY FROM INDEX
	for i, item := range worker.indexCache {

		if request.value == item || childOfRequest(item) {
			// DELETE ITEM
			delCount++

			if err := request.sendResult(item); err != nil {
				//fmt.Printf("error sending result for item %s: %s\n",item,err)
				return
			}
			//	fmt.Println("index remove: sent result: ",item )
			//fmt.Println("result",item)
		} else {
			// KEEP THIS ITEM. COPY ITEM OVER TO NEWCACHE
			newCache[i-delCount] = item

		}

	}

	if delCount > 0 {
		// TRUNCATE UNUSED NEWCACHE ELEMENTS BY  THE NUMBER
		// OF DELETED ITEMS
		worker.indexCache = newCache[:len(worker.indexCache)-delCount]
		worker.modCount += delCount
	}

	// SEND (SUCDESS) RESPONSE BACK
	//fmt.Printf("handleIndexRemoveRequest: before send response: request = %s\n",request)
	_ = request.sendResponse(msgIndexResponse{request: request, err: nil})

}

// getChildContainers retrieves a list of child container resources of  sr
// within the index (cache)
func (worker *fileRepositoryIndexWorker) getChildrenContainers(sr common.StagingResource) []common.StagingResource {

	childContainers := make([]common.StagingResource, 0)

	for _, entry := range worker.indexCache {

		if entry.DatasetPath == sr.DatasetPath &&
			entry.ContainerType == sr.ObjectType &&
			entry.ContainerIRI == sr.ObjectIRI {

			if entry.ObjectType.IsContainer() {
				// APPEND CHILDREN OF CURENT ENTRY TO LIST
				childContainers = append(childContainers, worker.getChildrenContainers(entry)...)
			}
			// APPEND ENTRY TO LIST
			childContainers = append(childContainers, entry)

		}
	}

	return childContainers

}

func (worker *fileRepositoryIndexWorker) handleIndexRollbackRequest(ctxt context.Context, request msgIndexRollbackRequest) {

	response := msgIndexResponse{}
	response.request = request
	response.err = nil

	if worker.modCount > 0 {
		// CLEAR CACHE TO ZERO ELEMENTS
		worker.indexCache = make([]common.StagingResource, 0)
		// RESET LAST CACHE UPDATE TO ZERO TIME
		worker.indexCacheLastUpdated = time.Time{}

		worker.modCount = 0

	}

	_ = request.sendResponse(response)

}

// updateIndexCacheIfInvalidated updates the local memory index cache of the
// receiver if the associated index file has changed since last cache refresh
func (worker *fileRepositoryIndexWorker) updateIndexCacheIfInvalidated() error {

	// CHECK IF INDEX FILE HAS BEEN UPDATED SINCE
	// LAST TIME CACHE WAS LOADED
	fileInfo, err := os.Stat(worker.index.path)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			if pathErr.Path == worker.index.path {
				// THE INDEX FILE DOESN'T EXIST.
				// THIS MAY HAPPEN IF NO COMMITS HAVE
				// OCCURRED YET
				return nil
			}
		}
		return err
	}

	if worker.indexCacheLastUpdated.Before(fileInfo.ModTime()) {
		// CACHE IS INVALIDATED. RELOAD
		var fileInfo os.FileInfo

		if worker.indexCache, fileInfo, err = worker.index.readIndexFile(); err != nil {

			return err
		}

		// UPDATE LAST CACHE UPDATED TIME
		worker.indexCacheLastUpdated = fileInfo.ModTime()

	}

	return nil
}

///////////////////////////////////////////////
// msgIndex*Request methods
//////////////////////////////////////////////

func (request msgIndexBaseRequest) sendResponse(response msgIndexResponse) error {

	var err error

	// SEND RESPONSE

	select {
	case request.responseChan <- response:

	case <-time.After(channelWriteTimeoutSeconds * time.Second):
		//close(msg.resultSetChan)
		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to send response message: { request = %s, response = %s }",
			channelWriteTimeoutSeconds, request, response)
	}

	// CLOSE RESPONSE CHANNEL AFTER WRITING RESPONSE OR TIMEOUT
	close(request.responseChan)

	return err
}

func (request msgIndexResultSetRequest) sendResponse(response msgIndexResponse) error {

	var err error

	// SEND RESPONSE

	select {
	case request.responseChan <- response:
	case <-time.After(channelWriteTimeoutSeconds * time.Second):
		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to send response message: { request = %s, response = { %s }",
			channelWriteTimeoutSeconds, request, response)
	}

	// CLOSE RESPONSE CHANNEL AFTER WRITING RESPONSE OR TIMEOUT

	close(request.resultSetChan)
	close(request.responseChan)

	return err

}

func (request msgIndexResultSetRequest) sendResult(result common.StagingResource) error {
	var err error

	// SEND RESPONSE

	select {
	case request.resultSetChan <- result:
	case <-time.After(channelWriteTimeoutSeconds * time.Second):
		// WRITE TIMEOUT. CLOSE MSG CHANNELS
		close(request.resultSetChan)
		close(request.responseChan)

		err = errors.TimedOut.Newf("timeout occurred after %d seconds waiting to send result message: { request = %s, result = { %s }",
			channelWriteTimeoutSeconds, request, result)
	}

	return err
}

///////////////////////////////////////////
// msgIndexBaseRequest functions
///////////////////////////////////////////

func (request msgIndexBaseRequest) String() string {
	return fmt.Sprintf("msgIndexBaseRequest { responseChan: %v }", request.responseChan)
}

func (response msgIndexResponse) String() string {
	return fmt.Sprintf("msgIndexResponse { err: %s }", response.err)
}

///////////////////////////////////////////
// msgIndexResultSetRequest functions
///////////////////////////////////////////

func (request msgIndexResultSetRequest) String() string {
	return fmt.Sprintf("msgIndexBaseRequest { responseChan: %v, resultSetChan: %v }",
		request.responseChan, request.resultSetChan)
}

func (req msgIndexResultSetRequest) next() (*common.StagingResource, error) {

	select {
	case result := <-req.resultSetChan:
		//	fmt.Println("next(): return result", result)
		return &result, nil
	case response := <-req.responseChan:
		//	fmt.Println("next(): return response")
		return nil, response.err
	}

}
