// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/skip2/openrdap/rdap/test"
)

func TestNetRegistryLookupsASN(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/asn.json")

	var n *ASNRegistry
	n, err := NewASNRegistry(bytes)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"as287",
			false,
			"AS287",
			[]string{"https://rdap.arin.net/registry", "http://rdap.arin.net/registry"},
		},
		{
			"As1768",
			false,
			"AS1768-AS1769",
			[]string{"https://rdap.apnic.net/"},
		},
		{
			"266652",
			false,
			"AS265629-AS266652",
			[]string{"https://rdap.lacnic.net/rdap/"},
		},
		{
			"not-a-number",
			true,
			"",
			[]string{},
		},
		{
			"999999",
			false,
			"",
			[]string{},
		},
	}

	runRegistryTests(t, tests, n)
}
