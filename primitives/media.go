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

//	"reflect"

//"github.com/adpadilla/go-dogg3rz/errors"
//"github.com/fatih/structs"

//const TYPE_DOGG3RZ_MEDIA = "dogg3rz.media"

const TYPE_DOGG3RZ_MEDIA Dogg3rzObjectType = 1 << 4

//const D_ATTR_ENTRIES = "entries"

type dgrzMedia struct {
	name string

	parent string
}
