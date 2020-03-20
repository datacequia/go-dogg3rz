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
