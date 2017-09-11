// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestNetRegistryLookupsDNSNested(t *testing.T) {
	test.Start(test.BootstrapComplex)
	defer test.Finish()

	var bytes []byte = test.Get("https://rdap.example.org/dns.json")

	var d *DNSRegistry
	d, err := NewDNSRegistry(bytes)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"",
			false,
			"",
			[]string{"https://example.root", "http://example.root"},
		},
		{
			"example.com",
			false,
			"com",
			[]string{"https://example.com", "http://example.com"},
		},
		{
			"sub.example.com",
			false,
			"sub.example.com",
			[]string{"https://example.com/sub", "http://example.com/sub"},
		},
		{
			"sub.sub.example.com",
			false,
			"sub.example.com",
			[]string{"https://example.com/sub", "http://example.com/sub"},
		},
		{
			"example.xyz",
			false,
			"",
			[]string{"https://example.root", "http://example.root"},
		},
	}

	runRegistryTests(t, tests, d)
}

func TestNetRegistryLookupsDNS(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/dns.json")

	var d *DNSRegistry
	d, err := NewDNSRegistry(bytes)

	if err != nil {
		t.Fatal(err)
	}

	tests := []registryTest{
		{
			"",
			false,
			"",
			[]string{},
		},
		{
			"www.EXAMPLE.BR",
			false,
			"br",
			[]string{"https://rdap.registro.br/"},
		},
		{
			"example.xyz",
			false,
			"",
			[]string{},
		},
	}

	runRegistryTests(t, tests, d)
}
