// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestServiceProviderRegistryLookups(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/object-tags.json")

	var s *ServiceProviderRegistry
	s, err := NewServiceProviderRegistry(bytes)

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
			"12345-FRNIC",
			false,
			"FRNIC",
			[]string{"https://rdap.nic.fr/"},
		},
		{
			"*-FRNIC",
			false,
			"FRNIC",
			[]string{"https://rdap.nic.fr/"},
		},
		{
			"-FRNIC",
			false,
			"FRNIC",
			[]string{"https://rdap.nic.fr/"},
		},
		{
			"A-B-FRNIC",
			false,
			"FRNIC",
			[]string{"https://rdap.nic.fr/"},
		},
	}

	runRegistryTests(t, tests, s)
}
