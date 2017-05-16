// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/skip2/rdap/test"
)

func TestDownloadAll(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	c := NewClient()

	if c.ASN() != nil || c.DNS() != nil || c.IPv4() != nil || c.IPv6() != nil {
		t.Fatalf("Registr(y|ies) not nil initially")
	}

	err := c.DownloadAll()

	if err != nil {
		t.Fatalf("DownloadAll() error: %s", err)
	}

	if c.ASN() == nil || c.DNS() == nil || c.IPv4() == nil || c.IPv6() == nil {
		t.Fatalf("Registr(y|ies) nil after DownloadAll()")
	}
}

func TestDownload(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	c := NewClient()

	err := c.Download(DNS)

	if err != nil {
		t.Fatalf("Download() error: %s", err)
	}

	if c.ASN() != nil || c.DNS() == nil || c.IPv4() != nil || c.IPv6() != nil {
		t.Fatalf("Download() bad")
	}
}

func TestLookups(t *testing.T) {
	tests := []struct {
		Registry RegistryType
		Input    string
		Success  bool
		URLs     []string
	}{
		{
			ASN,
			"as1768",
			true,
			[]string{"https://rdap.apnic.net/"},
		},
		{
			DNS,
			"example.br",
			true,
			[]string{"https://rdap.registro.br/"},
		},
		{
			IPv4,
			"41.0.0.0",
			true,
			[]string{
				"https://rdap.afrinic.net/rdap/",
				"http://rdap.afrinic.net/rdap/",
			},
		},
		{
			IPv6,
			"2001:1400::",
			true,
			[]string{
				"https://rdap.db.ripe.net/",
			},
		},
		{
			ServiceProvider,
			"12345~VRSN",
			true,
			[]string{"https://rdap.verisignlabs.com/rdap/v1"},
		},
	}

	test.Start(test.Bootstrap)
	test.Start(test.BootstrapExperimental)
	defer test.Finish()

	c := NewClient()

	for _, test := range tests {
		var r *Result

		r, err := c.Lookup(test.Registry, test.Input)

		if test.Success != (err == nil) {
			t.Errorf("Lookup %s: expected success=%v, got opposite, err=%v", test.Input, test.Success, err)
			continue
		}

		if r == nil {
			t.Errorf("Lookup %s: unexpected nil result", test.Input)
			continue
		}

		for i, url := range test.URLs {
			if r.URLs[i].String() != url {
				t.Errorf("Lookup %s, URL #%d, expected %s, got %s\n", test.Input, i, url, r.URLs[i])
				continue
			}
		}

	}
}
