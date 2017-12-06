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
		question := &Question{
			Query: test.Query,
		}

		var r *Answer
		r, err := reg.Lookup(question)

		if test.Error && err == nil {
			t.Errorf("Query: %s, expected error, didn't get one\n", test.Query)
			continue
		} else if !test.Error && err != nil {
			t.Errorf("Query: %s, unexpected error: %s\n", test.Query, err)
			continue
		}

		if test.Error {
			continue
		}

		if r == nil {
			t.Errorf("Query: %s, unexpected nil Answer, err=%v\n", test.Query, err)
			continue
		}

		if r.Entry != test.Entry {
			t.Errorf("Query: %s, expected Entry %s, got %s\n", test.Query, test.Entry, r.Entry)
			continue
		}

		if len(r.URLs) != len(test.URLs) {
			t.Errorf("Query: %s, expected %d urls, got %d\n", test.Query, len(test.URLs), len(r.URLs))
			continue
		}

		for i, url := range test.URLs {
			if r.URLs[i].String() != url {
				t.Errorf("Query %s, URL #%d, expected %s, got %s\n", test.Query, i, url, r.URLs[i])
				continue
			}
		}
	}
}
