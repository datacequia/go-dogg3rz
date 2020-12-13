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

package errors

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {

	errStr := "<my var name>"

	err := InvalidValue.New(errStr)
	bd, ok := err.(badDogg3rz)
	if !ok {
		t.Errorf("returned %v", bd)
	}

	//	fmt.Println("bd ", bd.errorType, NoType)
	if bd.errorType != InvalidValue {
		//fmt.Println("failed")
		t.Errorf("expected ErrorType %v, found %v", InvalidValue, bd.errorType)
	}

	if !strings.Contains(err.Error(), errStr) {
		t.Errorf("expected error (sub)string '%s' in '%s'. Not found", errStr, err.Error())

	}
}
