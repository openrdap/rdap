// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import "net/url"

// Answer represents the result of bootstrapping a single query.
type Answer struct {
	// Matching service entry. Empty string if no match.
	Entry string

	// Query looked up in the registry.
	//
	// This includes any canonicalization performed to match the Service
	// Registry's data format. e.g., lowercasing of domain names, and removal of
	// "AS" from AS numbers.
	Query string

	// List of RDAP base URLs.
	URLs []*url.URL
}
