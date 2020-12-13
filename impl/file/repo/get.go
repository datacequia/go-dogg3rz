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

package repo

import (
	"github.com/datacequia/go-dogg3rz/resource/repo"
)

type fileGetResourcesContext struct {
	repoName     string
	resourcePath string
	//schemaSubpath  common.RepositoryPath
	//fileSystemPath string
	resourceList []repo.GetResourceItem
}

type fileGetResourceItem struct {
	path   string
	status string
}

func (ctxt *fileGetResourcesContext) getResources(repoName string, resourcePath string) ([]repo.GetResourceItem, error) {

	return ctxt.resourceList, nil
}

func (i fileGetResourceItem) GetPath() string {
	return i.path
}

func (i fileGetResourceItem) GetStatus() string {
	return i.status

}
