// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"fmt"
	"net/url"
	"strings"
)

type DNSRegistry struct {
	DNS map[string][]*url.URL
}

// NewDNSRegistry creates a queryable DNS registry from a DNS registry JSON document.
//
// The document format is specified in https://tools.ietf.org/html/rfc7484#section-4.
func NewDNSRegistry(json []byte) (*DNSRegistry, error) {
	var r *registryFile
	r, err := parse(json)

	if err != nil {
		return nil, fmt.Errorf("Error parsing DNS bootstrap: %s", err)
	}

	return &DNSRegistry{
		DNS: r.Entries,
	}, nil
}

func (d *DNSRegistry) Lookup(input string) (*Result, error) {
	input = strings.TrimSuffix(input, ".")
	input = strings.ToLower(input)
	fqdn := input

	// Lookup the FQDN.
	// e.g. for an.example.com, the following lookups could occur:
	// - "an.example.com"
	// - "example.com"
	// - "com"
	// - "" (the root zone)
	var urls []*url.URL
	for {
		var ok bool
		urls, ok = d.DNS[fqdn]

		if ok {
			break
		} else if fqdn == "" {
			break
		}

		index := strings.IndexByte(fqdn, '.')
		if index == -1 {
			fqdn = ""
		} else {
			fqdn = fqdn[index+1:]
		}
	}

	return &Result{
		URLs:  urls,
		Query: input,
		Entry: fqdn,
	}, nil
}
