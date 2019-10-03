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
