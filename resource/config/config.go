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

package config

import (
	"bytes"
	"context"
	"text/template"
)

const CONFIG_JSON_SCHEMA = `
{
	"$id": "https://www.datacequia.com/dogg3rz.config.schema.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Dogg3rz configuration",
  "description": "Configuration schema for dogg3rz",
  "type": "object",
  "properties": {
    "ipfs": {
      "description": "IPFS Node Configuration Section",
      "type": "object",
			"properties": {
				"apiEndpoint": {
					"description":"IPFS Node REST API Endpoint",
					"type": "string"
				}
			},
			"required": ["apiEndpoint"]
    },
		"user": {
			"description":"User Information Section",
			"type": "object",
			"properties": {
				"email": {
					"description":"User's email address",
					"type":"string"
				},
				"firstName": {
					"description":"User's first name",
					"type": "string"
				},
				"lastName": {
					"description":"User's last name",
					"description":"string"
				}
			},
			"required":["email"]
		}
  },
  "required": [ "ipfs","user" ]
}
`

// use CONFIG_JSON_DEFAAULT 	with text/template to generraate default
const CONFIG_JSON_DEFAULT_TEMPLATE = `
{
    "ipfs": {
      "apiEndpoint":"{{ .IPFS.ApiEndpoint }}"
    },
		"user": {
			"email":"{{ .User.Email }}",
			"firstName":"{{ .User.FirstName }}",
			"lastName":"{{ .User.LastName  }}"
		}

}
`

type Dogg3rzConfig struct {
	IPFS IPFSConfig `json:"ipfs"`
	User UserConfig `json:"user"`
}

type IPFSConfig struct {
	ApiEndpoint string `json:"apiEndpoint"`
}

type UserConfig struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ConfigResource interface {
	GetConfig(ctxt context.Context) (*Dogg3rzConfig, error)
}

func GenerateDefault(config Dogg3rzConfig) (string, error) {

	//var user config.UserConfig
	var buf bytes.Buffer
	var tmpl *template.Template
	var err error

	// DEFAULT TO localhost:5001 if not provided
	if len(config.IPFS.ApiEndpoint) == 0 {
		config.IPFS.ApiEndpoint = "http://localhost:5001/"
	}

	tmpl, err = template.New("GenerateDefault").Parse(
		CONFIG_JSON_DEFAULT_TEMPLATE)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, config)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil

}
