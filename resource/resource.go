package resource

import (
	"fmt"
	"os"

	resourceconfig "github.com/datacequia/go-dogg3rz/resource/config"
	resourcenode "github.com/datacequia/go-dogg3rz/resource/node"
	resourcerepo "github.com/datacequia/go-dogg3rz/resource/repo"

	//fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	fileconfig "github.com/datacequia/go-dogg3rz/impl/file/config"
	filenode "github.com/datacequia/go-dogg3rz/impl/file/node"
	filerepo "github.com/datacequia/go-dogg3rz/impl/file/repo"
)

var configResource resourceconfig.ConfigResource
var nodeResource resourcenode.NodeResource
var repoResource resourcerepo.RepositoryResource

func init() {

	// DETERRMINE STORE TYPE
	storeType := os.Getenv("DOGG3RZ_STATE_STORE")
	if len(storeType) < 1 {
		// DEFAULT TO FILE
		storeType = "file"
	}

	switch storeType {
	case "file":
		nodeResource = &filenode.FileNodeResource{}
		configResource = &fileconfig.FileConfigResource{}
		repoResource = &filerepo.FileRepositoryResource{}
	default:

		fmt.Fprintf(os.Stderr,
			"unknown store type assigned to DOGG3RZ_STATE_STORE environment variable: %s", storeType)
		os.Exit(1)
	}

}

func GetConfigResource() resourceconfig.ConfigResource {

	return configResource
}

func GetNodeResource() resourcenode.NodeResource {
	return nodeResource
}

func GetRepositoryResource() resourcerepo.RepositoryResource {
	return repoResource
}
