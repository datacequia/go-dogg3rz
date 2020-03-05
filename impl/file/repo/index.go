package repo

import (
	"os"

	"bytes"
	"crypto/sha1"
	"encoding/binary"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/impl/file"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/util"
)

const FILE_INDEX_VERSION = uint32(1)

type FileRepositoryIndex struct {

	// O/S PATH TO BASE REPOSITORY
	// DIRECTORY
	repoDir string

	// THE REPOSITORY NAME
	repoName string

	// THE O/S PATH TO THE INDEX FILE
	path string
}

/*
Index format

HEADER
  SIGNATURE - 4 BYTES . ALWAYS 'RESC' (RESOURCE CACHE)
  VERSION   - 4 BYTES. INDEX FORMAT VERSION
  NUM ENTRIES - 4 BYTES, NUMBER OF INDEX ENTRIES
ENTRY
  DGRZ RESOURCE PATH (UNIX STYLE)
  NUL BYTE - 1 BYTE
  MULTIASH - THE (IPFS) HASH FINGERPRINT FOR THE RESOURCE CONTENT
    HASH FUNCTION TYPE - 1 BYTE,
    DIGEST LENGTH - 1 BYTE, THE LENGTH OF THE DIGEST
    DIGEST VALUE - THE CONTENT OF DIGEST
  NUL BYTE - 1 BYTE
CHECKSUM
  SHA-1 INDEX CHECKSUM - 160 BIT (8 BYTES) OVER CONTENT OF INDEX BEFORE THIS
                         CHECKSUM
*/

// Adds a new resource (resId) to the index
func (index *FileRepositoryIndex) Stage(resId rescom.RepositoryResourceId, multiHash string) error {

	// THE MULTIHASH AS AN ARRAY OF BYTES
	/*
		var mhBytes []byte

		var decodedMultihash *multihash.DecodedMultihash
		var err error

		if mhBytes, err = hex.DecodeString(multiHash); err != nil {

			return errors.InvalidArg.Wrapf(err, "FileRepositoryIndex.Add: multiHash")
		}

		if decodedMultihash, err = multihash.Decode(mhBytes); err != nil {
			return errors.InvalidArg.Wrapf(err, "FileRepositoryIndex.Add: multiHash")
		}
	*/

	//	addFunc := func() (io.Reader, error) {

	// GET EXISTING ENTRIES IN INDEX FILE//
	//		indexEntries, err := readIndexFile(index.path)
	//		if err != nil {
	//			return nil, err
	//		}

	// CREATE NEW ENTRY
	//		entry := indexEntry{resourceId: resId}

	// CHECK IF ENTRY EXISTS IN INDEX ALREADY
	//	var entryExists bool
	//	var entryIndex uint
	// APPEND ENTRY
	//		for i, e := range indexEntries {
	//	if e.resourceId.Kind() == entry.resourceId.Kind() {
	//			if e.resourceId.Subpath() == entry.resourceId.Subpath() {
	//entryExists = true
	//	entryIndex = uint(i)
	//			}
	//	}
	//		}

	//indexMap[resId.UnixStylePath()] = multiHash
	// TODO RETURN SOME VALUE
	//		return nil, nil
	//	}
	/*
		bytesWritten, err := file.WriteToFileAtomic(addFunc, index.path)
		if err != nil {

		}
	*/

	return nil // TODO return some value
}

func (index *FileRepositoryIndex) Unstage(resId rescom.RepositoryResourceId) error {

	return nil
}

func readIndexFile(indexPath string) ([]indexEntry, error) {

	/////////////////////////
	// OPEN INDEX FILE
	////////////////////////
	f, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}

	////////////////////
	// READ HEADER
	////////////////////
	const HeaderLength = 12
	const ChecksumLength = sha1.Size
	var indexFileSize = int64(0)

	// VALIDATE FILE SIZE OF INDEX IS LARGE
	// ENOUGH TO HOLD A MINIMUM SIZE INDEX
	if fileInfo, err := f.Stat(); err != nil {
		return nil, err

	} else {

		const minIndexFileSize = HeaderLength + ChecksumLength
		if fileInfo.Size() < minIndexFileSize {
			return nil, errors.UnexpectedValue.Newf("index: %s: expected at least a "+
				"%d byte file (header+checksum), found %d byte file",
				indexPath, minIndexFileSize, fileInfo.Size())

		}

		indexFileSize = fileInfo.Size()

	}

	// MAKE READ BUFFER TO READ WHOLE INDEX
	// INTO MEMORY
	var fileData = make([]byte, indexFileSize)

	// READ WHOLE FILE INTO BYTE ARRAY
	if numBytesRead, err := f.Read(fileData); err != nil {
		return nil, err
	} else {
		if int64(numBytesRead) != indexFileSize {
			return nil, errors.UnexpectedValue.Newf(
				"index: %s: expected to read %d bytes (file size), read only %d bytes",
				indexPath, indexFileSize, numBytesRead)
		}
	}

	// GET CHECKSUM AT END OF FILE
	checkSum := fileData[indexFileSize-ChecksumLength:]

	// RECOMPUTE CHECKSUM ON FILE CONTENTS
	verifyCheckSum := sha1.Sum(fileData[0 : indexFileSize-ChecksumLength])

	///////////////////////////
	// VERIFY CHECKSUM
	///////////////////////////

	// NOTE:  USE [:] TO CONVERT fixed byte[sha1.Size] array  to SLICE
	if !bytes.Equal(checkSum, verifyCheckSum[:]) {
		return nil, errors.UnexpectedValue.Newf("index: %s: checksum failed: "+
			"index checksum %v , computed checksum %v",
			indexPath, checkSum, verifyCheckSum)

	}

	////////////////////////////
	// VERIFY HEADER
	///////////////////////////

	// EVAL HEADER SIG
	headerSig := fileData[0:4]
	sig := []byte(file.ResourceCacheSignature)
	if !bytes.Equal(sig, headerSig) {
		return nil, errors.UnexpectedValue.Newf("index: %s: expected resource signature value %#v, found %#v",
			indexPath, sig, headerSig)
	}

	// EVAL HEADER VERSION
	headerVersion := binary.BigEndian.Uint32(fileData[4:8])
	if headerVersion != FILE_INDEX_VERSION {
		return nil, errors.UnexpectedValue.Newf("index: %s: expected header version  %d, found %d",
			indexPath, FILE_INDEX_VERSION, headerVersion)

	}

	// EVAL NUMBER OF ENTRIES
	numIndexEntries := binary.BigEndian.Uint32(fileData[8:12])

	indexEntries := make([]indexEntry, numIndexEntries)

	// COUNTER FOR ACTUAL INDEX ENTRIES ITERATED
	var actualIndexEntries = uint32(0)

	for nextOffset := HeaderLength; nextOffset < len(fileData); {

		var dgrzPath string
		var multiHash string

		dgrzPath, nextOffset = readNullTermString(fileData, nextOffset)
		multiHash, nextOffset = readNullTermString(fileData, nextOffset)

		if len(dgrzPath) > 0 && len(multiHash) > 0 {

			var resId rescom.RepositoryResourceId
			//	var decodedMultiHash *multihash.DecodedMultihash
			var err error

			//CONVERT STRING TO RESOURCE ID
			if resId, err = util.UnixStylePathToResourceId(dgrzPath); err != nil {
				return nil, err
			}

			// VALIDATE multiHash

			/*
				if hexDecodedMultiHash, err := hex.DecodeString(multiHash); err != nil {
					return nil, err
				} else {
					if decodedMultiHash, err := multihash.Decode(hexDecodedMultiHash); err != nil {
						return nil, err
					}
				}
			*/

			entry := indexEntry{resourceId: resId}

			// ADD TO SLICE OF INDEX ENTRIES
			indexEntries[actualIndexEntries] = entry
			// INCREMENT COUNT OF ACTUAL INDEX ENTRIES AFTER ADD
			actualIndexEntries++

		}

	}

	// MAKE SURE EXPECTED INDEX entries
	// AND ACTUAL MATCH UP
	if actualIndexEntries != numIndexEntries {
		return nil, errors.UnexpectedValue.Newf("index: %s: expected  %d index entries, counted %d",
			indexPath, numIndexEntries, actualIndexEntries)
	}

	// SUCCESS!
	return indexEntries, nil

}

// READ NULL TERMINATED STRING FROM BYTE ARRAY STARTING
// AT offset AND RETURN NEXT OFFSET TO READ OR len(buf)+1
func readNullTermString(buf []byte, offset int) (string, int) {

	var s string
	var lenBuf = len(buf)

	for i := offset; i < lenBuf; i++ {
		if buf[i] == 0x00 {
			// RETURN STRING AND NEXT OFFSET TO READ
			return string(buf[offset:i]), i + 1
		}
	}

	return s, lenBuf + 1

}
