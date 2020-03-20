package repo

import (
	"io"
	"os"

	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/datacequia/go-dogg3rz/errors"
	"github.com/datacequia/go-dogg3rz/primitives"
	"github.com/datacequia/go-dogg3rz/impl/file"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"

	cid "github.com/ipfs/go-cid"
//	"log"
//	"log"
	"fmt"
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

type indexEntry struct {
	Type      string
	Uuid      string
	MtimeNs   int64
	FileSize  int64
	Multihash string
	Subpath   string
}

type indexEntryInternal struct {
	Type uint32
	Uuid uuid.UUID
	MtimeNs int64
	FileSize int64
	MultihashLength int8
	Multihash []byte
	SubpathLength int16
	Subpath []byte
}

// MAP DOGG3RZOBJECT TYPES TO UIN32 EQUIVLANT
var dogg3rzObjectTypesToUint32Map map[primitives.Dogg3rzObjectType]uint32
// MAP UINT32 TO DOGG3RZOBJECTYPE
var uint32ToDogg3rzObjectTypesMap map[uint32]primitives.Dogg3rzObjectType

func init() {

	 dogg3rzObjectTypesToUint32Map = make(map[primitives.Dogg3rzObjectType]uint32)
	 uint32ToDogg3rzObjectTypesMap = make(map[uint32]primitives.Dogg3rzObjectType)

	 for _,t := range primitives.Dogg3rzObjectTypes() {

		 dogg3rzObjectTypesToUint32Map[t] = uint32(t)
		 uint32ToDogg3rzObjectTypesMap[uint32(t)] = t

	 }
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

func newFileRepositoryIndex(repoName string) (*fileRepositoryIndex, error) {

	index := &fileRepositoryIndex{}



	index.repoName = repoName

	index.repoDir = filepath.Join(file.RepositoriesDirPath(), repoName)

	index.path = filepath.Join(index.repoDir, file.IndexFileName)



	if !file.DirExists(index.repoDir) {
		return nil, errors.NotFound.Newf("repository directory at %s does not exist",
			index.repoDir)
	}

	return index, nil

}

// Adds a new resource (resId) to the index
func (index *fileRepositoryIndex) update(entry1 indexEntry) error {

	var internalEntry indexEntryInternal

	if i,err := ValidateIndexEntry(entry1); err != nil {
		return err
	} else {
		internalEntry = i
	}
	//fmt.Println("ValidateIndexEntry returned",internalEntry.String())

	addFunc := func() (io.Reader, error) {
		// GET EXISTING ENTRIES IN INDEX FILE//

		var indexEntriesInternal []indexEntryInternal

		if ie, err := index.readIndexFileInternal();err != nil {
			if os.IsNotExist(err) {
				// DOESN'T EXIST YET. ASSUME FOR NOW
				// THAT REPO DIR EXISTS BUT INDEX FILE DOESN'T
				indexEntriesInternal = make([]indexEntryInternal,0)
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
			if  e.Uuid == internalEntry.Uuid {
				 indexEntriesInternal[i] = internalEntry
				 updatedExistingEntry=true
			}
		}


		if ! updatedExistingEntry {
			indexEntriesInternal = append(indexEntriesInternal, internalEntry)

		}

		buf := &bytes.Buffer{}
		///////////////////////////////////
		// WRITE HEADER TO BUFFFER
		//////////////////////////////////
		if err := binary.Write(buf, binary.BigEndian, []byte(file.ResourceCacheSignature));
		err != nil {
			return nil, err
		}

		if err := binary.Write(buf,binary.BigEndian,file.IndexFormatVersion); err != nil {
			return nil, err
		}

		var numIndexEntries uint32 = uint32(len(indexEntriesInternal))
		if err := binary.Write(buf,binary.BigEndian,numIndexEntries); err!=nil {
			return nil,err
		}

		// WRITE ENTRIES TO BUFFER
		if err := writeIndexEntriesToBuffer(buf, indexEntriesInternal); err != nil {
			return nil, err
		}

		// WRITE SHA1 HASH

		var checkSum [ChecksumLength]byte = sha1.Sum(buf.Bytes())

		if err := binary.Write(buf,binary.BigEndian, checkSum); err != nil {
			return buf, err
		}


		return buf,nil

	}

	//log.Printf("index.path = %s",index.path)

	if _, err := file.WriteToFileAtomic(addFunc, index.path);err != nil {
		return err
	}


	return nil
}

func writeIndexEntriesToBuffer(buf *bytes.Buffer,indexEntries []indexEntryInternal) error {

	for _,e := range indexEntries {


		if err := binary.Write(buf, binary.BigEndian, e.Type); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, e.Uuid); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, e.MtimeNs); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, e.FileSize); err != nil {
			return  err
		}

	//fmt.Println("writeIndexEntriesToBuffer: multihashlength:",e.MultihashLength,e.FileSize)
		if err := binary.Write(buf, binary.BigEndian, e.MultihashLength); err != nil {
			return  err
		}


		if err := binary.Write(buf, binary.BigEndian, e.Multihash); err != nil {
			return  err
		}


		if err := binary.Write(buf, binary.BigEndian, e.SubpathLength); err != nil {
			return  err
		}

		if err := binary.Write(buf, binary.BigEndian, e.Subpath); err != nil {
			return  err
		}
	//	fmt.Printf("writeIndexEntriesToBuffer[%d]=%s\n",i,e.String())
	}

	return nil

}

func (index *fileRepositoryIndex) readIndexFile() ([]indexEntry,error) {

	  var indexEntriesInternal []indexEntryInternal = []indexEntryInternal{}
		var indexEntries []indexEntry = []indexEntry{}

		if i,err := index.readIndexFileInternal();err != nil  {
			return indexEntries,err
		} else {
			indexEntriesInternal = i
		}

		indexEntries = make([]indexEntry,len(indexEntriesInternal))
		for i, indexEntryInternal := range indexEntriesInternal {
				indexEntries[i] = indexEntryInternal.ToIndexEntry()
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

	//	var dgrzType uint32

		if err := binary.Read(f, binary.BigEndian, &e.Type); err != nil {
			return indexEntries, err
		}

	//	var id = make([]byte, 16)

		if err := binary.Read(f, binary.BigEndian, &e.Uuid); err != nil {
			return indexEntries, err
		}

	//	var mTimeNs int64

		if err := binary.Read(f, binary.BigEndian, &e.MtimeNs); err != nil {
			return indexEntries, err
		}

//		var fileSize int64

		if err := binary.Read(f, binary.BigEndian, &e.FileSize); err != nil {
			return indexEntries, err
		}

	//	var multiHashLength int8

		if err := binary.Read(f, binary.BigEndian, &e.MultihashLength); err != nil {
			return indexEntries, err
		}
		if e.MultihashLength < 1 {
			return indexEntries, errors.UnexpectedValue.Newf(
				"expected multihash length to be greater than zero, found length %d: "+
					"{ file = %s, entry = %d }", e.MultihashLength, index.path, i)

		}

		var multiHash []byte = make([]byte,e.MultihashLength)
		if err := binary.Read(f, binary.BigEndian, multiHash); err != nil {
			return indexEntries, err
		} else {
			e.Multihash = multiHash
		}

//		var subPathLength int16

		if err := binary.Read(f, binary.BigEndian, &e.SubpathLength); err != nil {
			return indexEntries, err
		}

		if e.SubpathLength  < 1 {
			return indexEntries, errors.UnexpectedValue.Newf(
				"expected sub path length to be greater than zero, found length %d: "+
					"{ file = %s, index entry offset = %d }",
					e.SubpathLength,
				index.path,
				i)

		}

		var subPath = make([]byte, e.SubpathLength)
		if err := binary.Read(f, binary.BigEndian, &subPath); err != nil {
			return indexEntries, err
		} else {
			e.Subpath = subPath
		}

		indexEntries[i] = e

	}

	return indexEntries, nil
}



func (e indexEntry) String() string {

	return fmt.Sprintf("indexEntry: { Type=%s, Uuid = %s, MtimeNs=%d, FileSize=%d, Multihash = %s, Subpath = %s }",
	e.Type,e.Uuid,e.MtimeNs,e.FileSize,e.Multihash,e.Subpath)

}

func ValidateIndexEntry(e indexEntry) (indexEntryInternal,error) {

	var internal indexEntryInternal

	if  dot,err := primitives.Dogg3rzObjectTypeFromString(e.Type);err != nil  {
		return internal, errors.UnexpectedValue.Newf("ValidateIndexEntry(): indexEntry.Type = '%s': " +
		"This identifier does not map to one of following defined in the 'primitives' package: %v",
		e.Type, primitives.Dogg3rzObjectTypes())
	}  else {
		internal.Type = uint32(dot)
	}

	internal.MtimeNs = e.MtimeNs


	if e.FileSize < 0 {
		return internal, errors.OutOfRange.Newf("ValidateIndexEntry(): indexEntry.FileSize = %d. " +
		"Must be greater than zero",e.FileSize)
	} else {
		internal.FileSize = e.FileSize
	}

	if _,err := cid.Decode(e.Multihash); err != nil {
		 return internal, errors.InvalidValue.Wrapf(err,"ValidateIndexEntry(): "+
		 "indexEntry.Multihash = %s. Expected valid multihash",e.Multihash)
	} else {
		internal.Multihash = []byte(e.Multihash)
		internal.MultihashLength = int8(len (internal.Multihash) )
	//	fmt.Println("ValidateIIndexEntry: multihashlength:",internal.MultihashLength)
	}

	if repoPath,err := rescom.RepositoryPathNew(e.Subpath); err != nil {
		return internal, errors.InvalidValue.Wrapf(err,"ValidateIndexEntry(): " +
	  "indexEntry.Subpath: %s",err)
	} else {
		 internal.Subpath = []byte(repoPath.ToString()) // used normalized path
		 internal.SubpathLength = int16(len(internal.Subpath))
	}

	if myUUID, err := uuid.Parse(e.Uuid); err != nil {
		return internal, errors.InvalidValue.Wrapf(err,"ValidateIndexEntry(): " +
		"indexEntry.Uuid: %s",err)

	} else {
		internal.Uuid = myUUID
	}

	return internal,nil

}

func (i indexEntryInternal) ToIndexEntry() indexEntry {

	var e indexEntry

	e.Type = uint32ToDogg3rzObjectTypesMap[i.Type].String()
	e.Uuid = i.Uuid.String()
	e.MtimeNs = i.MtimeNs
	e.FileSize = i.FileSize
	e.Subpath = string(i.Subpath)


	e.Multihash = string(i.Multihash)

	return e
}

func (i indexEntryInternal) String() string {

	return fmt.Sprintf("indexEntryInternal={Type=%s,Uuid=%s,FileSize=%d,MtimeNs=%d,MultiHashLength=%d,Multihash=%s,SubpathLength=%d,Subpath=%s}",
		uint32ToDogg3rzObjectTypesMap[i.Type].String(),
		i.Uuid.String(),
		i.FileSize,
		i.MtimeNs,
		i.MultihashLength,
		string(i.Multihash),
		i.SubpathLength,
		string(i.Subpath))
}
