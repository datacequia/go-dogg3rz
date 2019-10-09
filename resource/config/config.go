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
    }

  },
  "required": [ "ipfs" ]
}
`

const CONFIG_JSON_DEFAULT = `
{
    "ipfs": {
      "apiEndpoint":"http://localhost:5001/"
    }

}
`

type Dogg3rzConfig struct {
	IPFS IPFSConfig `json:"ipfs"`
}

type IPFSConfig struct {
	ApiEndpoint string `json:"apiEndpoint"`
}

type ConfigResource interface {
	GetConfig() (*Dogg3rzConfig, error)
}
