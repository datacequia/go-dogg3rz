package iricache

import (
	"testing"
)

func TestGrapplicationPathNew(t *testing.T) {

	for prefix, _ := range embeddedPrefixMap {

		if !IsCached(prefix) {
			t.Errorf("IsCached(\"%s\"): expected true", prefix)
		}

		if _, loc, err := Get(prefix); err != nil {
			t.Errorf("Get: %s: %s", prefix, err)

		} else {
			if len(loc) < 1 {
				t.Errorf("embedded location for prefix iri %s not set", prefix)
			}
		}

	}

}
