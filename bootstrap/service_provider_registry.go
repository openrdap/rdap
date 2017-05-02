// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"fmt"
	"net/url"
	"strings"
)

type ServiceProviderRegistry struct {
	Services map[string][]*url.URL // Map of service tag (e.g. "VRSN") to RDAP URLs.
}

// NewServiceProviderRegistry creates a queryable Service Provider registry
// from a Service Provider JSON document.
//
// The JSON bootstrap format is described in
// https://datatracker.ietf.org/doc/draft-hollenbeck-regext-rdap-object-tag/.
func NewServiceProviderRegistry(json []byte) (*ServiceProviderRegistry, error) {
	var r *registryFile
	r, err := parse(json)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Service Provider bootstrap: %s", err)
	}

	return &ServiceProviderRegistry{
		Services: r.Entries,
	}, nil
}

// Lookup returns a list of RDAP server URLs for an entity handle.
//
// e.g. for the handle "53774930~VRSN", the RDAP URLs for "VRSN" are returned.
//
// Missing/malformed/unknown service tags are not treated as errors. No URLs
// are returned in these cases.
func (s *ServiceProviderRegistry) Lookup(input string) (*Result, error) {
	// Valid input looks like 12345-VRSN.
	offset := strings.IndexByte(input, '~')

	if offset == -1 || offset == len(input)-1 {
		return &Result{
			Query: input,
		}, nil
	}

	service := input[offset+1:]

	urls, ok := s.Services[service]

	if !ok {
		service = ""
	}

	return &Result{
		URLs:  urls,
		Query: input,
		Entry: service,
	}, nil
}
