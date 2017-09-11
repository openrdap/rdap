// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/skip2/openrdap/rdap/test"
)

func TestNetRegistryLookupsIPv4(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/ipv4.json")

	var n *NetRegistry
	n, err := NewNetRegistry(bytes, 4)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"255.0.0.0",
			false,
			"",
			[]string{},
		},
		{
			"41.0.0.0",
			false,
			"41.0.0.0/8",
			[]string{
				"https://rdap.afrinic.net/rdap/",
				"http://rdap.afrinic.net/rdap/",
			},
		},
		{
			"41.255.255.255",
			false,
			"41.0.0.0/8",
			[]string{
				"https://rdap.afrinic.net/rdap/",
				"http://rdap.afrinic.net/rdap/",
			},
		},
		{
			"41.",
			true,
			"",
			[]string{},
		},
	}

	runRegistryTests(t, tests, n)
}

func TestNetRegistryLookupsIPv6(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/ipv6.json")

	var n *NetRegistry
	n, err := NewNetRegistry(bytes, 6)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"4000::",
			false,
			"",
			[]string{},
		},
		{
			"2001:1400::",
			false,
			"2001:1400::/23",
			[]string{
				"https://rdap.db.ripe.net/",
			},
		},
		{
			"2001:1400::5/128",
			false,
			"2001:1400::/23",
			[]string{
				"https://rdap.db.ripe.net/",
			},
		},
		{
			"2001:1400::/23",
			false,
			"2001:1400::/23",
			[]string{
				"https://rdap.db.ripe.net/",
			},
		},
		{
			"2001/129",
			true,
			"",
			[]string{},
		},
	}

	runRegistryTests(t, tests, n)
}
