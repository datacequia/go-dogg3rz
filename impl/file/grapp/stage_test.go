//go:build nothing

package grapp

import (
	"context"
	"testing"

	"github.com/datacequia/go-dogg3rz/env"
	"github.com/datacequia/go-dogg3rz/errors"
	rescom "github.com/datacequia/go-dogg3rz/resource/common"
	"github.com/datacequia/go-dogg3rz/resource/jsonld"
)

const (
	testDatasetName     = "testDataset"
	createTypeClassName = "class1"
	insertNode1Name     = "node1"
)

func TestFileGrapplicationResourceStager(t *testing.T) {

	//ctxt, _ := context.WithCancel(context.Background())

	ctxt, cancelFunc, err := initTestNode("FileGrapplicationResourceStager")
	if err != nil {
		t.Error(err)
		return
	}
	defer cancelFunc()

	//dogg3rzHome, _ := ctxt.Value(env.EnvDogg3rzHome).(string)
	grappName, _ := ctxt.Value(env.EnvDogg3rzGrapp).(string)

	//fmt.Println(dogg3rzHome)

	//defer os.RemoveAll(dogg3rzHome)

	stager, err := NewFileGrapplicationResourceStager(ctxt, grappName)
	if err != nil {
		t.Error(stager, err)
		return
	}

	defer stager.Close(ctxt)

	err = addTestResourcesToJSONLDDoc(ctxt)
	if err != nil {
		t.Error(err)
		return
	}

	srl := rescom.StagingResourceLocation{}

	srl.ContainerIRI = ""
	srl.ContainerType = jsonld.ContextResource
	srl.ObjectIRI = ""
	srl.ObjectType = jsonld.DatasetResource

	// add bad SRL datasets can't have a container
	err = stager.Add(ctxt, srl)
	if err == nil {
		t.Errorf("expected error when passing bad StagingResourceLocation")
	}
	if errors.GetType(err) != errors.InvalidValue {
		t.Errorf("expected error InvalidValue, found %s", err)
	}

	// STAGE OUTERMOST CONTEXT THAT WAS ALREADY POPULATED IN WORKSPACE
	srl.ContainerIRI = ""
	srl.ContainerType = jsonld.DatasetResource
	srl.ObjectIRI = ""
	srl.ObjectType = jsonld.ContextResource
	srl.DatasetPath = testDatasetName

	//srl.DatasetPath = testDatasetName

	err = stager.Add(ctxt, srl)
	if err != nil {
		t.Errorf("stage outermost context failed: %s", err)
		return
	}

	// STAGE NEW SCHEMA TYPE
	// STAGE OUTERMOST CONTEXT THAT WAS ALREADY POPULATED IN WORKSPACE
	srl.ContainerIRI = ""
	srl.ContainerType = jsonld.DatasetResource
	srl.ObjectIRI = createTypeClassName
	srl.ObjectType = jsonld.NodeResource
	srl.DatasetPath = testDatasetName

	err = stager.Add(ctxt, srl)
	if err != nil {
		t.Errorf("stage schema class type failed failed: %s", err)
		return
	}

	// STAGE NODE
	srl.ContainerIRI = ""
	srl.ContainerType = jsonld.DatasetResource
	srl.ObjectIRI = insertNode1Name
	srl.ObjectType = jsonld.NodeResource
	srl.DatasetPath = testDatasetName

	err = stager.Add(ctxt, srl)
	if err != nil {
		t.Errorf("stage node failed: %s", err)
		return
	}

	// ATTEMPT TO STAGE NODE THAT DOES NOT EXIST IN WORKSPACE
	srl.ContainerIRI = ""
	srl.ContainerType = jsonld.DatasetResource
	srl.ObjectIRI = "DOES_NOT_EXIST"
	srl.ObjectType = jsonld.NodeResource
	srl.DatasetPath = testDatasetName

	err = stager.Add(ctxt, srl)
	if err == nil {
		t.Errorf("stage non existent node did not fail")
		return
	}

	err = stager.Commit(ctxt)
	if err != nil {
		t.Errorf("stager commit failed: %s", err)
	}

}

// //////////////////////////////
// test helper funcs
// /////////////////////////////
func addTestResourcesToJSONLDDoc(ctxt context.Context) error {

	frr := &FileGrapplicationResource{}

	grappName, _ := ctxt.Value(env.EnvDogg3rzGrapp).(string)

	err := frr.CreateDataset(ctxt, grappName, testDatasetName)
	if err != nil {
		return err
	}

	// create outermost context prefix
	err = frr.AddNamespaceDataset(ctxt, grappName, testDatasetName, "rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	if err != nil {
		return err
	}

	// create new schema type
	err = frr.CreateTypeClass(ctxt, grappName, testDatasetName, createTypeClassName, "", "my first class", "my first class")
	if err != nil {
		return err
	}

	// create new node
	nodeProps := []string{"a", "b", "c"}
	nodeValues := []string{"1", "2", "3"}
	err = frr.InsertNode(ctxt, grappName, testDatasetName, createTypeClassName, insertNode1Name, "", nodeProps, nodeValues)
	if err != nil {
		return err
	}
	return nil

}
