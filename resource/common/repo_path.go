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


package common

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/datacequia/go-dogg3rz/errors"
)

type RepositoryPath struct {
	pathElements        []string
	lastCharPathElement bool
}

var REPO_PATH_SEPARATOR = '/'

var VALID_PATH_ELEMENT_SPECIAL_CHARS = []rune{'.', '_', '-'}

func RepositoryPathNew(path string) (*RepositoryPath, error) {

	if len(path) < 1 {
		return nil, errors.InvalidValue.New("repository path is zero length (empty) string")
	}

	rp := &RepositoryPath{}

	if rune(path[len(path)-1]) == REPO_PATH_SEPARATOR {
		// PATH LAST CHARACTER ENDS WITH PATH SEP. MUST BE REFERRING TO DIR
		rp.lastCharPathElement = true

	}

	var standardizedPathElements []string

	paths := strings.Split(path, string(REPO_PATH_SEPARATOR))

	for _, path := range paths {

		if len(path) < 1 {
			// EMPTY PATH ELEMENT. SKIP
			//fmt.Println("skip")
			continue
		}

		// ITERATE EACH RUNE IN CURRENT PATH ELEMENT
		// TO MAKE SURE IT CONTAINS VALID PATH ELEMENT CHARACTERS

		for offset, c := range path {
			switch offset {
			case 0:
				// FIRST CHARACTER IN PATH ELEMENT EVAL HERE
				if !(unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c))) {
					return nil, errors.InvalidPathElement.Newf(
						"%s: expecting path element that begins "+
							"with alphanumeric character, found '%s'", path, string(c))
				}
			default:
				// ALL OTHER PATH ELEMENT CHARS EVAL HERE
				if !(unicode.IsLetter(rune(c)) || unicode.IsNumber(rune(c)) ||
					isPathElementSpecialChar(rune(c))) {
					return nil, errors.InvalidPathElement.Newf(
						"expecting path elements that contain "+
							"alphanumeric characters or the following special characters only: "+
							validPathElementSpecialCharsToString()+
							": %s: found '%s' character at offset %d", path, string(c), offset)
				}
			}
		}

		standardizedPathElements = append(standardizedPathElements, path)

	}

	if len(standardizedPathElements) < 1 {
		return nil, errors.InvalidValue.New("repository path has zero path elements")

	}

	rp.pathElements = standardizedPathElements

	return rp, nil

}

func (rp *RepositoryPath) ToString() string {

	pj := strings.Join(rp.pathElements, string(REPO_PATH_SEPARATOR))

	if rp.lastCharPathElement {
		// TACK ON A TRAILING PATH SEPARATOR
		// SINCE IT WAS SPECIFIED WHEN THIS OBJECT WAS
		// CONSTRUCTED

		// I.E. INTENDED TO BE A DIR
		return pj + string(REPO_PATH_SEPARATOR)
	}

	return pj

}

func (rp *RepositoryPath) ToOperatingSystemPath() string {
	pj := strings.Join(rp.pathElements, string(os.PathSeparator))

	if rp.lastCharPathElement {
		// TACK ON A TRAILING PATH SEPARATOR
		// SINCE IT WAS SPECIFIED WHEN THIS OBJECT WAS
		// CONSTRUCTED

		// I.E. INTENDED TO BE A DIR
		return pj + string(os.PathSeparator)
	}

	return pj

}

func (rp *RepositoryPath) Size() int {
	return len(rp.pathElements)
}

func (rp *RepositoryPath) EndsWithPathSeparator() bool {
	return rp.lastCharPathElement

}

func (rp *RepositoryPath) PathElements() []string {
	return rp.pathElements

}

func isPathElementSpecialChar(r rune) bool {

	for _, r2 := range VALID_PATH_ELEMENT_SPECIAL_CHARS {
		if r == r2 {
			return true
		}
	}

	return false
}

func validPathElementSpecialCharsToString() string {

	var newList []string

	for _, y := range VALID_PATH_ELEMENT_SPECIAL_CHARS {
		newList = append(newList, fmt.Sprintf("'%s'", string(y)))
	}

	return strings.Join(newList, ", ")

}
