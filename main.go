package main

import (
	"fmt"

	cmd "github.com/adpadilla/go-dogg3rz/cmd"
	"github.com/adpadilla/go-dogg3rz/errors"
)

func testErr() error {

	err1 := errors.InvalidArg.New("ddd")

	e2 := errors.AddErrorContext(err1, "test1", "test1value")
	e3 := errors.AddErrorContext(e2, "test2", "test2value")

	//var ff int = 7
	var y *int = nil //&ff

	i, e5 := errors.NotNil(y)
	if !i {

		return errors.InvalidArg.Wrap(e5, "y")
	} else {
		fmt.Println("not nil")
	}

	e4 := errors.Range.Wrapf(e3, "sss")
	return e4
}

func main() {

	e := testErr()

	fmt.Println("the err:", e)

	cmd.Run()
}
