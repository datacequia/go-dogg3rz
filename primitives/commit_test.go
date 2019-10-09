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

package primitives

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {

	commit, err := Dogg3rzCommitNew("test", "myPeerId", "dogg3rz@datacequia.com")

	if err != nil {
		t.Errorf("failed to create Dogg3rzCommit object: { error = %v }", err)
	}
	var b strings.Builder
	encoder := json.NewEncoder(&b)

	encoder.Encode(commit.Dogg3rzObject())

	decoder := json.NewDecoder(strings.NewReader(b.String()))

	do := Dogg3rzObjectNew("objectType")
	err2 := decoder.Decode(do)
	if err2 != nil {
		t.Errorf("failed to decode commit object from json to Dogg3rzObject: { error = %v }", err2)
	}

	if do.ObjectType != TYPE_DOGG3RZ_COMMIT {
		t.Errorf("bad object type. expected %s, got %s", TYPE_DOGG3RZ_COMMIT, do.ObjectType)
	}

	// CHECK THAT PEERID VALUE WAS SET
	if do.Metadata[MD_ATTR_IPFS_PEER_ID] != "myPeerId" {
		t.Errorf("bad metadata attr %s. expected %s, got %s", MD_ATTR_IPFS_PEER_ID, "myPeerId", do.Metadata[MD_ATTR_IPFS_PEER_ID])
	}
	// CHECK THAT EMAIL ADDR VALUE WAS SET
	if do.Metadata[MD_ATTR_EMAIL_ADDR] != "dogg3rz@datacequia.com" {
		t.Errorf("bad metadata attr %s. expected %s, got %s", MD_ATTR_EMAIL_ADDR, "dogg3rz@datacequia.com", do.Metadata[MD_ATTR_EMAIL_ADDR])

	}
	// CHECK THAT THE METADATA NAME ATTR WAS SET
	if rootTree, ok := do.Data[D_ATTR_ROOT_TREE].(map[string]interface{}); ok {

		if metadata, ok := rootTree["metadata"].(map[string]interface{}); ok && metadata[MD_ATTR_NAME] != "test" {
			t.Errorf("unexpected root tree name. expected %s, got %s", "test", metadata[MD_ATTR_NAME])
		}
	} else {
		t.Errorf("commit attr 'rootTree' does not have expected type")
	}

}
