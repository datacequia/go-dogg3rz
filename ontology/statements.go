package ontology

import (
	"reflect"

	"github.com/datacequia/go-dogg3rz/env/dev"
)

type Term string
type IRI string // TODO: convert to object

type Object map[string]interface{}

type ResourceIdentifier struct {
	Id string `json:"@id"`
}

type RDFProperty struct {
	RDFSResource

	SubPropertyOf *ResourceIdentifier `json:"rdfs:subPropertyOf,omitempty"`
	Domain        string              `json:"rdfs:domain,omitempty"`
	Range         string              `json:"rdfs:range,omitempty"`
}

type RDFSResource struct {
	ResourceIdentifier
	Type        string `json:"@type"`
	Comment     string `json:"rdfs:comment,omitempty"`
	IsDefinedBy string `json:"rdfs:isDefinedBy,omitempty"`
	Label       string `json:"rdfs:label,omitempty"`
	Member      string `json:"rdfs:member,omitempty"`
	SeeAlso     string `json:"rdfs:seeAlso,omitempty"`
}

type RDFSClass struct {
	RDFSResource

	SubClassOf *ResourceIdentifier `json:"rdfs:subClassOf,omitempty"`
}

type GrapplicationMetadata struct {
}

type ParquetFile struct {
}

var context = Object{
	"@base": dev.GitRemoteURL,
	"rdf":   "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
}

var graph = []any{

	grapplicationMetadataClassDecl(),

	namespacePropertyDecl(),
}

type Dataset struct {
	Context Object `json:"@context"`
	Graph   []any  `json:"@graph"`
}

var ontology = &Dataset{
	Context: context,
	Graph:   graph,
}

func resourceId(id string) *ResourceIdentifier {
	if len(id) > 0 {
		return &ResourceIdentifier{
			Id: id,
		}
	}
	return nil

}

func applicationRuntimeFile() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(GrapplicationMetadata{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "A file used by ",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(GrapplicationMetadata{}).Name(),
			SeeAlso:     "",
		},
		SubClassOf: &ResourceIdentifier{
			Id: "rdfs:Resource",
		},
	}
	c.SubClassOf = resourceId("rdfs:Class")

	return &c

}

func grapplicationMetadataClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(GrapplicationMetadata{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Describes Grappliction Resources",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(GrapplicationMetadata{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Resource")

	return &c

}

func namespacePropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "namespace",
			},
			Type:        "rdfs:Property",
			Comment:     "Namespace(s) imported into the Grapplication",
			IsDefinedBy: "",
			Label:       "namespace",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationMetadata{}).Name(),
		Range:  "xsd:anyURI",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func stagedDataPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "stagedDataSet",
			},
			Type:        "rdfs:Property",
			Comment:     "Datasets added but not comitted to Grapplication",
			IsDefinedBy: "",
			Label:       "stagedDataSet",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationMetadata{}).Name(),
		Range:  reflect.TypeOf(ParquetFile{}).Name(),
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}
