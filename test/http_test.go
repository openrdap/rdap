// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package test

import (
	"testing"
	"strings"
)

func TestSmoke(t *testing.T) {
	Start(Bootstrap)
	defer Finish()

	var bytes []byte
	bytes = Get("https://data.iana.org/rdap/asn.json")

	if !strings.Contains(string(bytes), "ripe.net") {
		t.Fatalf("ASN doesn't contain ripe.net: %s\n", string(bytes))
	}
}
