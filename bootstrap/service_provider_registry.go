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
	// Map of service tag (e.g. "VRSN") to RDAP base URLs.
	services map[string][]*url.URL

	// The registry's JSON document.
	file *RegistryFile
}

// NewServiceProviderRegistry creates a ServiceProviderRegistry from a Service
// Provider JSON document.
//
// The document format is specified in
// https://datatracker.ietf.org/doc/draft-hollenbeck-regext-rdap-object-tag/.
func NewServiceProviderRegistry(json []byte) (*ServiceProviderRegistry, error) {
	var r *RegistryFile
	r, err := parse(json)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Service Provider bootstrap: %s", err)
	}

	return &ServiceProviderRegistry{
		services: r.Entries,
		file:     r,
	}, nil
}

// Lookup returns a list of RDAP base URLs for the |input| entity handle.
//
// e.g. for the handle "53774930~VRSN", the RDAP base URLs for "VRSN" are returned.
//
// Missing/malformed/unknown service tags are not treated as errors. An empty
// list of URLs is returned in these cases.
func (s *ServiceProviderRegistry) Lookup(input string) (*Result, error) {
	// Valid input looks like 12345-VRSN.
	offset := strings.IndexByte(input, '~')

	if offset == -1 || offset == len(input)-1 {
		return &Result{
			Query: input,
		}, nil
	}

	service := input[offset+1:]

	urls, ok := s.services[service]

	if !ok {
		service = ""
	}

	return &Result{
		URLs:  urls,
		Query: input,
		Entry: service,
	}, nil
}

// File returns a struct describing the registry's JSON document.
func (s *ServiceProviderRegistry) File() *RegistryFile {
	return s.file
}
