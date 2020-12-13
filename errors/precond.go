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
	"fmt"
	"reflect"
	"runtime"
	"unicode"

	"github.com/pkg/errors"
)

func NotNil(value interface{}) (bool, error) {

	x := reflect.ValueOf(value).Kind()
	if x != reflect.Ptr {
		// it's not a pointer so can't be nil
		return true, nil
	}

	//fmt.Println("ValueOf=", reflect.ValueOf(value))
	// GET VALUE OF POINTER
	val2 := reflect.ValueOf(value).Pointer()
	//fmt.Println("val2=", val2)
	if val2 == 0 {
		// POINTER IS NULL (ZERO)
		fpcs := make([]uintptr, 1)
		n := runtime.Callers(2, fpcs)

		if n == 0 {
			return false, errors.New("nil pointer")
		} else {
			caller := runtime.FuncForPC(fpcs[0] - 1)
			if caller == nil {
				return false, errors.New("nil pointer")
			} else {
				fileName, lineNo := caller.FileLine(fpcs[0] - 1)

				return false, errors.Errorf("nil pointer - {func=%s, file=%s, line=%d}", caller.Name(), fileName, lineNo)

			}
		}
	}

	return true, nil

}

func StrlenGtZero(str string) (bool, error) {
	if len(str) < 1 {
		return false, errors.New("string length less than one")
	}

	return true, nil
}

func StrAlpha(str string) (bool, error) {
	if ok, err := StrlenGtZero(str); !ok {
		return ok, err
	}

	for i, c := range str {
		if !unicode.IsLetter(c) {
			return false, fmt.Errorf("non-alpha character found at position %d", i)
		}
	}

	return true, nil

}

func PathElementValid(pathElement string) (bool, error) {
	if ok, err := StrlenGtZero(pathElement); !ok {
		return ok, err

	}

	for i, c := range pathElement {

		switch true {
		case unicode.IsLetter(c):
		case unicode.IsDigit(c):
		case unicode.IsSpace(c):
		case c == '-':
		case c == '_':
		default:
			return false, InvalidPathElement.Newf("non-allowable path element character found in path element '%s': {character = '%s', position = %d, ascii value = %d}", pathElement, string(c), i, int(c))
		}
	}
	return true, nil
}
