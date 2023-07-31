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

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/datacequia/go-dogg3rz/errors"
)

// PARSE GRAPP / SCHEMA PATH IN  GRAPP:SCHEMA_PATH FORMAT
func parseGrappSchemaPath(grappSchemaPath string) (string, string, error) {

	elements := strings.SplitN(grappSchemaPath, ":", 2)
	if len(elements) != 2 {
		return "", "", errors.UnexpectedValue.Newf("found '%s': want format GRAPP:SCHEMA_PATH",
			grappSchemaPath)

	}

	return elements[0], elements[1], nil

}

func parseGrappAndPathMaybe(grappAndPath string) (string, string, error) {

	elements := strings.SplitN(grappAndPath, ":", 2)
	switch len(elements) {
	case 1:
		return elements[0], "", nil
	case 2:
		return elements[0], elements[1], nil
	}

	return "", "", errors.UnexpectedValue.Newf("found '%s': want format GRAPP[:PATH]",
		grappAndPath)

}

func printValues(values []string, ignoreValue string, out io.Writer) {

	for _, v := range values {
		fmt.Fprintln(out, v)

	}
}
