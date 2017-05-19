// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"encoding/json"
	"errors"
	"net/url"
)

// RegistryFile represents a bootstrap registry file (i.e. {asn,dns,ipv4,ipv6}.json).
type RegistryFile struct {
	// Fields from the JSON document.
	Description string
	Publication string
	Version     string

	// Map of service entries to RDAP base URLs.
	//
	// e.g. in ipv6.json, the following mapping:
	// "2c00::/12" => https://rdap.afrinic.net/rdap/,
	//                http://rdap.afrinic.net/rdap/.
	Entries map[string][]*url.URL

	// The file's JSON document.
	JSON []byte
}

func parse(jsonDocument []byte) (*RegistryFile, error) {
	var doc struct {
		Description string
		Publication string
		Version     string

		Services [][][]string
	}

	err := json.Unmarshal(jsonDocument, &doc)
	if err != nil {
		return nil, err
	}

	b := &RegistryFile{}
	b.Description = doc.Description
	b.Publication = doc.Publication
	b.Version = doc.Version
	b.JSON = jsonDocument

	b.Entries = make(map[string][]*url.URL)

	for _, s := range doc.Services {
		if len(s) != 2 {
			return nil, errors.New("Malformed bootstrap (bad services array)")
		}

		entries := s[0]
		rawURLs := s[1]

		var urls []*url.URL

		for _, rawURL := range rawURLs {
			url, err := url.Parse(rawURL)

			// Ignore unparsable URLs.
			if err != nil {
				continue
			}

			urls = append(urls, url)
		}

		if len(urls) > 0 {
			for _, entry := range entries {
				b.Entries[entry] = urls
			}
		}
	}

	return b, nil
}
