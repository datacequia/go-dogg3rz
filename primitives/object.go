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

package primitives

import (
	"io"
)

const (
	TYPE_DOGG3RZ_OBJECT          Dogg3rzObjectType = 1 << 3
	DOGG3RZ_OBJECT_ATTR_TYPE                       = "type"
	DOGG3RZ_OBJECT_ATTR_ID                         = "id"
	DOGG3RZ_OBJECT_ATTR_METADATA                   = "meta"
	DOGG3RZ_OBJECT_ATTR_BODY                       = "body"
	DOGG3RZ_OBJECT_ATTR_PARENT                     = "parent"
)

type dgrzObject struct {
	ObjectType string                 `structs:"type" json:"type"`
	Metadata   map[string]string      `structs:"metadata" json:"metadata"`
	Data       map[string]interface{} `structs:"data,omitempty" json:"data,omitempty"`
	Parents    []string               `structs:"parent,omitempty" json:"parent,omitempty"`
}

type Dogg3rzObject interface {
	//ToDogg3rzObject() *dgrzObject

	Type() string
	Id() string

	LastModified()

	JSONReadCloser() io.ReadCloser
}
