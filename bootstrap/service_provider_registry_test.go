// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/keperry/rdap/test"
)

func TestServiceProviderRegistryLookups(t *testing.T) {
	test.Start(test.BootstrapExperimental)
	defer test.Finish()

	var bytes []byte = test.Get("https://test.rdap.net/rdap/serviceprovider-draft-03.json")

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
			"~",
			false,
			"",
			[]string{},
		},
		{
			"X~VRSN~",
			false,
			"",
			[]string{},
		},
		{
			"12345~VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"*~VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"~VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"12345-VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"*-VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"-VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
		{
			"A-B-VRSN",
			false,
			"VRSN",
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
	}

	runRegistryTests(t, tests, s)
}
