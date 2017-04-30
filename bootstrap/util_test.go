// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import "testing"

type registryTest struct {
	Query string

	Error bool

	Entry string
	URLs  []string
}

func runRegistryTests(t *testing.T, tests []registryTest, reg Registry) {
	for _, test := range tests {
		var r *Result
		r, err := reg.Lookup(test.Query)

		if test.Error && err == nil {
			t.Fatalf("Query: %s, expected error, didn't get one\n", test.Query)
		} else if !test.Error && err != nil {
			t.Fatalf("Query: %s, unexpected error: %s\n", test.Query, err)
		}

		if test.Error {
			continue
		}

		if r == nil {
			t.Fatalf("Query: %s, unexpected nil Result, err=%v\n", test.Query, err)
		}

		if r.Entry != test.Entry {
			t.Fatalf("Query: %s, expected Entry %s, got %s\n", test.Query, test.Entry, r.Entry)
		}

		if len(r.URLs) != len(test.URLs) {
			t.Fatalf("Query: %s, expected %d urls, got %d\n", test.Query, len(test.URLs), len(r.URLs))
		}

		for i, url := range test.URLs {
			if r.URLs[i].String() != url {
				t.Fatalf("Query %s, URL #%d, expected %s, got %s\n", test.Query, i, url, r.URLs[i])
			}
		}
	}
}
