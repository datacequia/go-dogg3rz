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

package config

import (
	"bytes"
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
	GetConfig() (*Dogg3rzConfig, error)
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
