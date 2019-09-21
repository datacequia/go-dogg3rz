package primitives

import (
	"testing"
)

func TestTreeNew(t *testing.T) {

	if tree, err := Dogg3rzTreeNew("abc123 -_XYZ"); err != nil {
		t.Errorf("failed to create new dogg3rz tree object with valid tree name: { error = %v, tree = %v }", err, tree)
	}

	if tree, err := Dogg3rzTreeNew("abc123 -_XY#Z"); err == nil {

		t.Errorf("failed to create new dogg3rz tree object with valid tree name: { error = %v, tree = %v }", err, tree)
	}

}

func TestPutGetTree(t *testing.T) {

	var (
		parent *dgrzTree
		child  *dgrzTree
		err    error
	)
	if parent, err = Dogg3rzTreeNew("parent"); err != nil {
		t.Errorf("failed to create new dogg3rz tree object with valid tree name: { error = %v, tree = %v }", err, parent)
	}

	if child, err = Dogg3rzTreeNew("child"); err != nil {
		t.Errorf("failed to create new dogg3rz tree object with valid tree name: { error = %v, tree = %v }", err, child)
	}

	err = parent.PutTree(child)
	if err != nil {
		t.Errorf("error putting child tree in parent tree: { error = %v, parent = %v, child = %v}", err, parent, child)
	}

	child2, err2 := parent.GetEntry("child")
	if err2 != nil {
		t.Errorf("error getting child tree entry from parent tree that was just put: { error = %v, parent = %v, child = %v}", err2, parent, child2)
	}

	if c2, ok := child2.(*dgrzTree); !ok {
		t.Errorf("did not retrieve child as tree type: { c2 = %v}", c2)
	} else {

		if c2 != child {
			// IF WE PUT A POINTER TO dgrzTree, should receive the same pointer is assumption

			t.Errorf("child tree entry put != child tree get: { child put = %p, child get = %p}", child, c2)
		}
	}

}
