// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/openrdap/rdap/bootstrap/cache"
	"github.com/openrdap/rdap/test"
)

func TestDownload(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	c := &Client{}

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
			"12345-FRNIC",
			true,
			[]string{"https://rdap.nic.fr/"},
		},
	}

	test.Start(test.Bootstrap)
	test.Start(test.BootstrapExperimental)
	defer test.Finish()

	c := &Client{}

	for _, test := range tests {
		var r *Answer

		question := &Question{
			RegistryType: test.Registry,
			Query:        test.Input,
		}

		r, err := c.Lookup(question)

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

// A pre-populated cache should be used instead of downloading. The HTTP
// responders here all return 404, so a download attempt fails the test.
func TestLookupUsesPrePopulatedCache(t *testing.T) {
	test.Start(test.BootstrapHTTPError)
	defer test.Finish()

	memCache := cache.NewMemoryCache()
	if err := memCache.Save("dns.json", test.LoadFile("bootstrap/dns.json")); err != nil {
		t.Fatalf("Save() error: %s", err)
	}

	c := &Client{Cache: memCache}

	answer, err := c.Lookup(&Question{RegistryType: DNS, Query: "example.br"})
	if err != nil {
		t.Fatalf("Lookup() error: %s", err)
	}

	if len(answer.URLs) != 1 || answer.URLs[0].String() != "https://rdap.registro.br/" {
		t.Errorf("Lookup() got %v, want [https://rdap.registro.br/]", answer.URLs)
	}

	if c.DNS() == nil {
		t.Error("DNS() returned nil after loading the registry from the cache")
	}
}

// An expired cached file is used when the refresh download fails, rather than
// failing the lookup.
func TestLookupFallsBackToExpiredCache(t *testing.T) {
	test.Start(test.BootstrapHTTPError)
	defer test.Finish()

	memCache := cache.NewMemoryCache()
	if err := memCache.Save("dns.json", test.LoadFile("bootstrap/dns.json")); err != nil {
		t.Fatalf("Save() error: %s", err)
	}

	memCache.SetTimeout(-1 * time.Second)

	if got := memCache.State("dns.json"); got != cache.Expired {
		t.Fatalf("cache state = %v, want Expired", got)
	}

	c := &Client{Cache: memCache}

	answer, err := c.Lookup(&Question{RegistryType: DNS, Query: "example.br"})
	if err != nil {
		t.Fatalf("Lookup() error: %s", err)
	}

	if len(answer.URLs) != 1 || answer.URLs[0].String() != "https://rdap.registro.br/" {
		t.Errorf("Lookup() got %v, want [https://rdap.registro.br/]", answer.URLs)
	}
}

// An expired file is re-downloaded even though it's already parsed into memory.
func TestLookupRefreshesExpiredCache(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	c := &Client{}

	if _, err := c.Lookup(&Question{RegistryType: DNS, Query: "example.br"}); err != nil {
		t.Fatalf("Lookup() error: %s", err)
	}

	downloads := httpmock.GetTotalCallCount()

	c.Cache.SetTimeout(-1 * time.Second)

	if _, err := c.Lookup(&Question{RegistryType: DNS, Query: "example.br"}); err != nil {
		t.Fatalf("Lookup() error: %s", err)
	}

	if got := httpmock.GetTotalCallCount() - downloads; got != 1 {
		t.Errorf("expired registry triggered %d downloads, want 1", got)
	}
}

// An unparseable cached file plus a failed download reports the download error,
// rather than panicking on the unpopulated registry.
func TestLookupWithCorruptCacheAndDownloadError(t *testing.T) {
	test.Start(test.BootstrapHTTPError)
	defer test.Finish()

	memCache := cache.NewMemoryCache()
	if err := memCache.Save("dns.json", []byte("{{{ not json")); err != nil {
		t.Fatalf("Save() error: %s", err)
	}

	c := &Client{Cache: memCache}

	if _, err := c.Lookup(&Question{RegistryType: DNS, Query: "example.br"}); err == nil {
		t.Error("Lookup() unexpectedly succeeded")
	}
}

func TestLookupWithDownloadError(t *testing.T) {
	test.Start(test.BootstrapHTTPError)
	defer test.Finish()

	c := &Client{}

	question := &Question{
		RegistryType: DNS,
		Query:        "example.br",
	}

	_, err := c.Lookup(question)

	if err == nil {
		t.Errorf("Unexpected success")
	}

	t.Logf("Error was: %s", err)
}
