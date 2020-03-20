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
