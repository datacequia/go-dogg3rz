package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {

	errStr := "<my var name>"

	err := InvalidArg.New(errStr)
	bd, ok := err.(badDogg3rz)
	if !ok {
		t.Errorf("returned %v", bd)
	}

	fmt.Println("bd ", bd.errorType, NoType)
	if bd.errorType != InvalidArg {
		//fmt.Println("failed")
		t.Errorf("expected ErrorType %v, found %v", InvalidArg, bd.errorType)
	}

	if !strings.Contains(err.Error(), errStr) {
		t.Errorf("expected error (sub)string '%s' in '%s'. Not found", errStr, err.Error())

	}
}
