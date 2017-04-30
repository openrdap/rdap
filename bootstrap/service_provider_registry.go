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
	Services map[string][]*url.URL
}

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
