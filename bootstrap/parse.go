// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"encoding/json"
	"errors"
	"net/url"
)

type registryFile struct {
	Description string
	Publication string
	Version     string

	Entries map[string][]*url.URL
}

func parse(jsonDocument []byte) (*registryFile, error) {
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

	b := &registryFile{}
	b.Description = doc.Description
	b.Publication = doc.Publication
	b.Version = doc.Version

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
