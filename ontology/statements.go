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

type Snapshot struct {
}

type GrapplicationImage struct {
}

type GrapplicationRuntimeImage struct {
}

type GrapplicationMetadata struct {
}

type DataFile struct {
}
type ParquetFile struct {
	DataFile
}

type Namespace struct {
}

type DeployableNamespace struct {
}

var context = Object{
	"@base": dev.GitRemoteURL,
	"rdf":   "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
}

var graph = []any{

	snapshotClassDecl(),
	imagePropertyDecl(),
	signaturePropertyDecl(),
	namespacePropertyDecl(),
	grapplicationImageClassDecl(),
	grapplicationRuntimeImageClassDecl(),
	dataPropertyDecl(),
	metadataPropertyDecl(),
	dataFileClassDecl(),
	parquetFileClassDecl(),
	grapplicationMetadataClassDecl(),
	namespacePropertyDecl(),
	namespaceClassDecl(),
	iriPropertyDecl(),
	prefixPropertyDecl(),
	deployableNamespaceClassDecl(),
	containerImageIdPropertyDecl(),
	gitCommitHashPropertyDecl(),
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

/*
func applicationRuntimeFile() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(Sn{}).Name(),
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
*/
// ////////////////////////////////////////////////////////////////////////
// RDFS Datatype Subclasses used by this ontology
// ////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////
// Snapshot Class Declaration and its Properties
// ////////////////////////////////////////////////////////////////////////
func snapshotClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(Snapshot{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "A point in time capture (image) of a Grapplication",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(Snapshot{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Resource") // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func imagePropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "image",
			},
			Type:        "rdfs:Property",
			Comment:     "Grapplication image object",
			IsDefinedBy: "",
			Label:       "image",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(Snapshot{}).Name(),
		Range:  reflect.TypeOf(GrapplicationImage{}).Name(),
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func signaturePropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "signature",
			},
			Type:        "rdfs:Property",
			Comment:     "Digital signature of the GrapplicationImage object's  Content Identifer (CID)",
			IsDefinedBy: "",
			Label:       "signature",
			Member:      "",
			SeeAlso:     "dgrz:image",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(Snapshot{}).Name(),
		Range:  "xsd:string",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

//////////////////////////////////////////////////////////////////////////
// GrapplicationImage Class Declaration and its properties
//////////////////////////////////////////////////////////////////////////

func grapplicationImageClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(GrapplicationImage{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Represents a point in design time of a Grapplication",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(GrapplicationImage{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Resource") // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func dataPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "data",
			},
			Type:        "rdfs:Property",
			Comment:     "Represents the current sum of all datasets added to the grapplication",
			IsDefinedBy: "",
			Label:       "data",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationImage{}).Name(),
		Range:  reflect.TypeOf(DataFile{}).Name(),
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func metadataPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "metadata",
			},
			Type:        "rdfs:Property",
			Comment:     "Represents all metadata describing the Grapplication",
			IsDefinedBy: "",
			Label:       "metadata",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationImage{}).Name(),
		Range:  reflect.TypeOf(GrapplicationMetadata{}).Name(),
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func grapplicationRuntimeImageClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(GrapplicationRuntimeImage{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Represents a point-in-time snapshot of a Grapplication runtime",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(GrapplicationRuntimeImage{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId(reflect.TypeOf(GrapplicationImage{}).Name()) // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

// ////////////////////////////////////////////////////////////////////////
// Datafile class/subclasses and their properties
// ////////////////////////////////////////////////////////////////////////
func dataFileClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(DataFile{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Where a Grapplication's RDF triples are stored and retrieved",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(DataFile{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Resource") // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func parquetFileClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(ParquetFile{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Where a Grapplication's RDF triples are stored and retrieved using parquet file format.",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(ParquetFile{}).Name(),
			SeeAlso:     "https://parquet.apache.org/",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId(reflect.TypeOf(DataFile{}).Name())

	return &c

}

// ////////////////////////////////////////////////////////////////////////
// GrapplicationMetadata Class and its properties
// ////////////////////////////////////////////////////////////////////////
func grapplicationMetadataClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(GrapplicationMetadata{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "Describes Grapplication",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(GrapplicationMetadata{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Class") // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func namePropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "name",
			},
			Type:        "rdfs:Property",
			Comment:     "name of Grapplication",
			IsDefinedBy: "",
			Label:       "name",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationMetadata{}).Name(),
		Range:  "xsd:string",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func namespacePropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "namespace",
			},
			Type:        "rdfs:Property",
			Comment:     "namespace (IRIs) imported into Grapplication and their associated properties and state",
			IsDefinedBy: "",
			Label:       "namespace",
			Member:      "",
			//SeeAlso:     "",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(GrapplicationMetadata{}).Name(),
		Range:  reflect.TypeOf(Namespace{}).Name(),
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

// ////////////////////////////////////////////////////////////////////////
// Namespace Class and its properties
// ////////////////////////////////////////////////////////////////////////

func namespaceClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(Namespace{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "A namespace used by the Grapplication",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(Namespace{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId("rdfs:Class") // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func iriPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "iri",
			},
			Type:        "rdfs:Property",
			Comment:     "International Resource Identifier address of ontology",
			IsDefinedBy: "",
			Label:       "iri",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(Namespace{}).Name(),
		Range:  "xsd:anyURI",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func prefixPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "prefix",
			},
			Type:        "rdfs:Property",
			Comment:     "prefix for associated 'iri' property ",
			IsDefinedBy: "",
			Label:       "prefix",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(Namespace{}).Name(),
		Range:  "xsd:string",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

// ////////////////////////////////////////////////////////////////////////
// DeployableNamespace Class and its properties
// ////////////////////////////////////////////////////////////////////////
func deployableNamespaceClassDecl() *RDFSClass {

	c := RDFSClass{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: reflect.TypeOf(DeployableNamespace{}).Name(),
			},
			Type:        "rdfs:Class",
			Comment:     "A namespace iri that points to a remote git project.  When built as a container image,  it will provide its ontology and behavior to the grapplication",
			IsDefinedBy: "",
			Label:       "",
			Member:      reflect.TypeOf(DeployableNamespace{}).Name(),
			SeeAlso:     "",
		},
		//SubClassOf: &ResourceIdentifier{},
	}
	c.SubClassOf = resourceId(reflect.TypeOf(Namespace{}).Name()) // allow other properties with rdfs:Class domain to be assigned to this class

	return &c

}

func containerImageIdPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "containerImageId",
			},
			Type:        "rdfs:Property",
			Comment:     "A container image identifier (hash) which identifies the deployable namespace",
			IsDefinedBy: "",
			Label:       "containerImageId",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(DeployableNamespace{}).Name(),
		Range:  "xsd:string",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}

func gitCommitHashPropertyDecl() *RDFProperty {

	p := &RDFProperty{
		RDFSResource: RDFSResource{
			ResourceIdentifier: ResourceIdentifier{
				Id: "gitCommitHash",
			},
			Type:        "rdfs:Property",
			Comment:     "The git commit hash of a deployable IRI namespace that refers to a remote git repository",
			IsDefinedBy: "",
			Label:       "gitCommitHash",
			Member:      "",
			//SeeAlso:     "https://parquet.apache.org/",
		},
		//SubPropertyOf: resourceId(""),
		Domain: reflect.TypeOf(DeployableNamespace{}).Name(),
		Range:  "xsd:string",
	}

	//p.SubPropertyOf = resourceId("rdfs:Property")

	return p
}
