
BUILD_TAGS= #ipfs_embed
GIT_COMMIT_HASH := $$(git rev-parse HEAD)
GIT_REMOTE_NAME := $$(git remote)
GIT_REMOTE_URL := $$(git remote get-url $(GIT_REMOTE_NAME)) 
GO_PKG_BASE := "github.com/datacequia/go-dogg3rz"
build:
	go build -tags "$(BUILD_TAGS)" -ldflags "-X $(GO_PKG_BASE)/env/dev.GitCommitHash=$(GIT_COMMIT_HASH) -X $(GO_PKG_BASE)/env/dev.GitRemoteName=$(GIT_REMOTE_NAME) -X $(GO_PKG_BASE)/env/dev.GitRemoteURL=$(GIT_REMOTE_URL)"  -o dogg3rz main.go 
	mkdir -p dist/
	mv dogg3rz dist/

install:
	mv dist/dogg3rz $$GOPATH/bin

