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

// ipfs package ipfs provides wrapper functions to IPFS interactions that
// are used by dogg3rz
package ipfs

import (
	"bytes"
	"encoding/json"
	//	"fmt"

	shell "github.com/ipfs/go-ipfs-api"
)

// DagPut takes a data construct  'data' and commits to IPFS as
// an IPLD object. Returns CID
func DagPut(data interface{}) (string, error) {

	sh := newShell()

	// handle error coming from sh.DagPut() where it does not
	// accept map[string] interface {} by converting to a byte buffer (io.Reader)
	// NOTE: error is "cannot current handle putting values of type map[string]interface {}"
	if _, ok := (data).(map[string]interface{}); ok {

		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)

		if err := encoder.Encode(data); err != nil {
			return "", err
		}

		//fmt.Println("%s", string(buf.Bytes()))
		// SET DATA TO BE AN io.Reader (via bytes.Buffer)
		data = buf
	}

	return sh.DagPut(data, "json", "cbor")

}

func newShell() *shell.Shell {

	// TODO GET IPFS HOST FROM CONFIG
	return shell.NewShell("localhost:5001")
}
