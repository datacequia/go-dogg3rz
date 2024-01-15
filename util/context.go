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

import "context"

func ContextValueAsString(ctxt context.Context, key interface{}) (string, bool) {

	var i interface{}
	var s string

	if i = ctxt.Value(key); i != nil {
		// has value
		if s, ok := i.(string); ok {
			// is string value
			return s, true
		}
	}

	return s, false

}

func ContextValueAsStringOrDefault(ctxt context.Context, key interface{}, defaultValue string) string {

	var s string
	var ok bool

	if s, ok = ContextValueAsString(ctxt, key); ok {
		return s
	}

	return defaultValue

}
