package iricache

import (
	"embed"

	"github.com/datacequia/go-dogg3rz/errors"
)

//go:embed documents/*.jsonld
var jsonld_docs embed.FS

const RDF_PREFIX = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
const RDFS_PREFIX = "http://www.w3.org/2000/01/rdf-schema#"
const OWL_PREFIX = "http://www.w3.org/2002/07/owl#"
const DCAM_PREFIX = "http://purl.org/dc/dcam/"
const DC_PREFIX = "http://purl.org/dc/elements/1.1/"
const PROF_PREFIX = "http://www.w3.org/ns/dx/prof/"
const XSD_PREFIX = "http://www.w3.org/2001/XMLSchema#"
const PROV_PREFIX = "http://www.w3.org/ns/prov#"
const DCTERMS_PREFIX = "http://purl.org/dc/terms/"
const DOAP_PREFIX = "http://usefulinc.com/ns/doap#"
const QB_PREFIX = "http://purl.org/linked-data/cube#"
const VOID_PREFIX = "http://rdfs.org/ns/void#"
const DCAT_PREFIX = "http://www.w3.org/ns/dcat#"
const SH_PREFIX = "http://www.w3.org/ns/shacl#"
const SKOS_PREFIX = "http://www.w3.org/2004/02/skos/core#"

const BRICK_PREFIX = "https://brickschema.org/schema/Brick#"
const CSVW_PREFIX = "http://www.w3.org/ns/csvw#"
const SOSA_PREFIX = "http://www.w3.org/ns/sosa/"
const TIME_PREFIX = "http://www.w3.org/2006/time#"
const VANN_PREFIX = "http://purl.org/vocab/vann/"
const DCMITYPE_PREFIX = "http://purl.org/dc/dcmitype/"
const FOAF_PREFIX = "http://xmlns.com/foaf/0.1/"
const ODRL_PREFIX = "http://www.w3.org/ns/odrl/2/"
const ORG_PREFIX = "http://www.w3.org/ns/org#"
const SSN_PREFIX = "http://www.w3.org/ns/ssn/"
const GRDDL_PREFIX = "http://www.w3.org/2003/g/data-view#"
const ADMS_PREFIX = "http://www.w3.org/ns/adms#"

const documentsDirName = "documents" // subdirectory from this path

// for any [X]_PREFIX const defined above put in map here
var embeddedPrefixMap = map[string]string{
	RDF_PREFIX:      documentsDirName + "/rdf.jsonld",
	RDFS_PREFIX:     documentsDirName + "/rdfs.jsonld",
	OWL_PREFIX:      documentsDirName + "/owl.jsonld",
	DCAM_PREFIX:     documentsDirName + "/dcam.jsonld",
	DC_PREFIX:       documentsDirName + "/dc.jsonld",
	PROF_PREFIX:     documentsDirName + "/prof.jsonld",
	XSD_PREFIX:      documentsDirName + "/xsd.jsonld",
	PROV_PREFIX:     documentsDirName + "/prov.jsonld",
	DCTERMS_PREFIX:  documentsDirName + "/dcterms.jsonld",
	DOAP_PREFIX:     documentsDirName + "/doap.jsonld",
	QB_PREFIX:       documentsDirName + "/qb.jsonld",
	VOID_PREFIX:     documentsDirName + "/void.jsonld",
	DCAT_PREFIX:     documentsDirName + "/dcat.jsonld",
	SH_PREFIX:       documentsDirName + "/sh.jsonld",
	SKOS_PREFIX:     documentsDirName + "/skos.jsonld",
	BRICK_PREFIX:    documentsDirName + "/brick.jsonld",
	CSVW_PREFIX:     documentsDirName + "/csvw.jsonld",
	SOSA_PREFIX:     documentsDirName + "/sosa.jsonld",
	TIME_PREFIX:     documentsDirName + "/time.jsonld",
	VANN_PREFIX:     documentsDirName + "/time.jsonld",
	DCMITYPE_PREFIX: documentsDirName + "/dcmitype.jsonld",
	FOAF_PREFIX:     documentsDirName + "/foaf.jsonld",
	ODRL_PREFIX:     documentsDirName + "/odrl.jsonld",
	ORG_PREFIX:      documentsDirName + "/org.jsonld",
	SSN_PREFIX:      documentsDirName + "/ssn.jsonld",
	GRDDL_PREFIX:    documentsDirName + "/grddl.jsonld",
	ADMS_PREFIX:     documentsDirName + "/adms.jsonld",
}

// returns true if  supplied prefix is in embedded cache
func IsCached(iriPrefix string) bool {

	if _, ok := embeddedPrefixMap[iriPrefix]; ok {
		return true
	}

	return false

}

func Get(iriPrefix string) ([]byte, string, error) {

	if v, ok := embeddedPrefixMap[iriPrefix]; ok {
		var embed_location = "urn:go:embed:" + v
		data, err := jsonld_docs.ReadFile(v)
		if err != nil {
			return nil, "", err
		}
		return data, embed_location, nil

	}

	return nil, "", errors.NotFound.Newf("prefix '%s' not found in embedded cache", iriPrefix)

}
