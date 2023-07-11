package dev

import (
	"reflect"
	"strings"
)

// Populated by linker with current hash (HEAD)
// of this project. See Makefile
var GitCommitHash string
var GitRemoteName string
var GitRemoteURL string

type dummy struct {
}

func PackageName() string {

	return strings.TrimSuffix(reflect.TypeOf(dummy{}).PkgPath(), "/")

}
