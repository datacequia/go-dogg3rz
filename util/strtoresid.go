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

package util

import (
	"fmt"
	"strings"

	dgrzerr "github.com/datacequia/go-dogg3rz/errors"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
)

type resId struct {
	user            string
	commitMultiHash string
	kind            string
	subPath         string
}

func UnixStylePathToResourceId(path string) (rescom.GrapplicationResourceId, error) {

	pathElements := strings.SplitN(path, "/", 5)

	if len(pathElements) < 5 {
		return nil, dgrzerr.OutOfRange.Wrap(dgrzerr.UnexpectedValue.Newf("expected at least 4 path elements, found %d", len(pathElements)-1), "bad resource id")
	}

	dgrzNs := pathElements[1]
	userCommitMultiHash := pathElements[2]
	kind := pathElements[3]
	subPath := pathElements[4]

	if dgrzNs != rescom.RootPathElementName {
		return nil, dgrzerr.InvalidValue.Wrap(
			dgrzerr.InvalidPathElement.Newf("expected first path element to be '%s', found '%s'",
				rescom.RootPathElementName, dgrzNs), "bad resource id")

	}

	var user string
	var commitMultiHash string

	pathElements = strings.SplitN(userCommitMultiHash, "@", 2)
	switch len(pathElements) {
	case 1:
		commitMultiHash = pathElements[0]
	case 2:
		user = pathElements[0]
		commitMultiHash = pathElements[1]
	default:
		//
		panic("commitMultiHash should only have 2 components")

	}
	return resId{user: user, commitMultiHash: commitMultiHash, kind: kind, subPath: subPath}, nil

}

func (o resId) User() string {
	return o.user
}

func (o resId) CommitMultiHash() string {
	return o.commitMultiHash
}

func (o resId) Kind() string {
	return o.kind
}

func (o resId) Subpath() string {
	return o.subPath
}

func (o resId) UnixStylePath() string {

	return fmt.Sprintf("/%s/%s/%s/%s", rescom.RootPathElementName,
		o.userCommitMultiHash(), o.kind, o.subPath)

}

func (o resId) HttpUrl(host string, port uint16) string {

	return fmt.Sprintf("http://%s:%d/%s/%s/%s/%s", host, port,
		rescom.RootPathElementName,
		o.userCommitMultiHash(), o.kind, o.subPath)
}

func (o resId) userCommitMultiHash() string {

	var elem1 string
	if len(o.user) > 0 {
		elem1 = o.user + "@" + o.commitMultiHash
	} else {
		elem1 = o.commitMultiHash
	}

	return elem1

}
