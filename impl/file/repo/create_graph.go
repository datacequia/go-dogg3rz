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
	"context"
	"fmt"
)

func (repo *FileRepositoryResource) CreateNamedGraph(ctxt context.Context, repoName string, datasetPath string, graphName string,
	parentGraphName string) error {
	var fds *fileDataset
	var err error
	if fds, err = newFileDataset(ctxt, repoName, datasetPath); err != nil {
		return err
	}

	if err = fds.createNamedGraph(ctxt, graphName, parentGraphName); err != nil {
		fmt.Println(err)
		return err

	}
	return nil
}
