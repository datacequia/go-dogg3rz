package primitives

import (
	"encoding/json"
	"io"
)

const (
	TYPE_DOGG3RZ_OBJECT             = "dogg3rz.object"
	DOGG3RZ_OBJECT_ATTR_OBJECT_TYPE = "type"
	DOGG3RZ_OBJECT_ATTR_METADATA    = "metadata"
	DOGG3RZ_OBJECT_ATTR_DATA        = "data"
	DOGG3RZ_OBJECT_ATTR_PARENT      = "parent"
)

type dgrzObject struct {
	ObjectType string                 `structs:"type" json:"type"`
	Metadata   map[string]string      `structs:"metadata" json:"metadata"`
	Data       map[string]interface{} `structs:"data,omitempty" json:"data,omitempty"`
	Parent     string                 `structs:"parent,omitempty" json:"parent,omitempty"`
}

type Dogg3rzObjectifiable interface {
	ToDogg3rzObject() *dgrzObject
}

func Dogg3rzObjectDeserializeFromJson(reader io.Reader) (*dgrzObject, error) {

	decoder := json.NewDecoder(reader)

	obj := dgrzObjectNew()

	err := decoder.Decode(obj)
	if err != nil {
		return nil, err
	}

	return obj, err

}

func Dogg3rzObjectSerializeToJson(obj *dgrzObject, writer io.Writer) error {

	encoder := json.NewEncoder(writer)

	err := encoder.Encode(obj)

	return err

}

func Dogg3rzObjectNew(objectType string) *dgrzObject {

	o := dgrzObjectNew()

	o.ObjectType = objectType

	return o

}

func dgrzObjectNew() *dgrzObject {
	return &dgrzObject{ObjectType: "", Metadata: make(map[string]string),
		Data: make(map[string]interface{})}
}
