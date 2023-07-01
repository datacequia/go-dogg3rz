package ontology

import (
	"encoding/json"
	"io"
)

func Get(out io.Writer) {

	if b, e := json.MarshalIndent(ontology, "", "    "); e == nil {
		out.Write(b)
	}

}
