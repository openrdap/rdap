// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package bootstrap implements Registration Data Access Protocol (RDAP) bootstrapping.
//
// All RDAP queries are handled by an RDAP server. To help clients discover
// RDAP servers, IANA publishes Service Registry files
// (https://data.iana.org/rdap) for several query types: Domain names, IP
// addresses, and Autonomous Systems.
//
// Given an RDAP query, this package finds the list of RDAP server URLs which
// can answer it. This includes downloading & parsing the Service Registry
// files.
//
// Basic usage:
//   var result *bootstrap.Result
//
//   b := bootstrap.NewClient()
//   result, err := b.Lookup(bootstrap.DNS, "google.cz") // Downloads https://data.iana.org/rdap/dns.json automatically.
//
//   if err == nil {
//     for _, url := range result.URLs {
//       fmt.Println(url)
//     }
//   }
//
// Download and list the contents of the DNS Service Registry:
//   b := bootstrap.NewClient()
//
//   // Before you can use a Registry, you need to download it first.
//   err := b.Download(bootstrap.DNS) // Downloads https://data.iana.org/rdap/dns.json.
//
//   if err == nil {
//     var dns *DNSRegistry = b.DNS()
//
//     // Print TLDs with RDAP service.
//     for _, tld := range dns.DNS {
//       fmt.Println(tld)
//     }
//   }
//
// A bootstrap.Client caches the Service Registry files in memory for both performance, and courtesy to data.iana.org. The functions which make network requests are:
//   - Download()      - download one Service Registry file.
//   - DownloadAll()   - download all four Service Registry files.
//
//   - Lookup()        - download one Service Registry file if missing, or if the cached file is over (by default) 24 hours old.
//
// Lookup() is intended for repeated usage: A long lived bootstrap.Client will
// download each of {asn,dns,ipv4,ipv6}.json once per 24 hours only, regardless
// of the number of calls made to Lookup(). You can still refresh them manually
// using Download() if required.
//
// As well as the default memory cache, bootstrap.Client also supports caching
// the Service Registry files on disk. The default cache location is
// $HOME/.openrdap/.
//
// Disk cache usage:
//
//   b := bootstrap.NewClient()
//   b.Cache = cache.NewDiskCache()
//
//   dsr := b.DNS()  // Tries to load dns.json from disk cache, doesn't exist yet, so returns nil.
//   b.Download(bootstrap.DNS) // Downloads dns.json, saves to disk cache.
//
//   b2 := bootstrap.NewClient()
//   b2.Cache = cache.NewDiskCache()
//
//   dsr2 := b.DNS()  // Loads dns.json from disk cache.
//
// This package also implements the experimental Service Provider registry. Due
// to the experimental nature, no Service Registry file exists on data.iana.org
// yet. Instead, an unofficial one is downloaded from
// https://www.openrdap.org/rdap/.
//
// RDAP bootstrapping is defined in https://tools.ietf.org/html/rfc7484.
package bootstrap

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/skip2/rdap/bootstrap/cache"
)

// A RegistryType represents a bootstrap registry type.
type RegistryType int

const (
	DNS RegistryType = iota
	IPv4
	IPv6
	ASN
	ServiceProvider
)

const (
	// Default URL of the Service Registry files.
	DefaultBaseURL = "https://data.iana.org/rdap/"

	// Default cache timeout of Service Registries.
	DefaultCacheTimeout = time.Hour * 24

	// Location of the experimental service_provider.json.
	experimentalBaseURL = "https://www.openrdap.org/rdap/"
)

// Client implements an RDAP bootstrap client.
//
// Create a Client using NewClient().
type Client struct {
	HTTP    *http.Client        // HTTP client.
	BaseURL *url.URL            // Base URL of the Service Registry files. Default is DefaultBaseURL.
	Cache   cache.RegistryCache // Service Registry cache. Default is a MemoryCache.

	registries map[RegistryType]Registry
}

// A Registry implements bootstrap lookups.
type Registry interface {
	Lookup(input string) (*Result, error)
}

// Result represents the result of bootstrapping a single query.
type Result struct {
	// Query looked up in the registry.
	//
	// This includes any canonicalisation performed to match the Service
	// Registry's data format. e.g. lowercasing of domain names, and removal of
	// "AS" from AS numbers.
	Query string

	// Matching service entry. Empty string if no match.
	Entry string

	// List of RDAP base URLs.
	URLs []*url.URL
}

// NewClient creates a new bootstrap.Client.
func NewClient() *Client {
	c := &Client{
		HTTP:  &http.Client{},
		Cache: cache.NewMemoryCache(),

		registries: make(map[RegistryType]Registry),
	}

	c.BaseURL, _ = url.Parse(DefaultBaseURL)
	c.Cache.SetTimeout(DefaultCacheTimeout)

	c.registries[ASN] = nil
	c.registries[DNS] = nil
	c.registries[IPv4] = nil
	c.registries[IPv6] = nil
	c.registries[ServiceProvider] = nil

	return c
}

// Download downloads a single bootstrap registry file.
//
// On success, the relevant Registry is refreshed. Use the matching accessor (ASN(), DNS(), IPv4(), or IPv6()) to access it.
//
// For example, to download and list the DNS bootstrap registry file:
//   b := bootstrap.NewClient()
//   err := b.Download(bootstrap.DNS)
//
//   if err == nil {
//     dns := b.DNS()
//
//     for _, tld := range dns.DNS {
//       fmt.Println(tld)
//     }
//   }
func (c *Client) Download(registry RegistryType) error {
	var json []byte
	var s Registry

	json, s, err := c.download(registry)

	if err != nil {
		return err
	}

	err = c.Cache.Save(registry.Filename(), json)
	if err != nil {
		return err
	}

	c.registries[registry] = s

	return nil
}

func (c *Client) download(registry RegistryType) ([]byte, Registry, error) {
	u, err := url.Parse(registry.Filename())
	if err != nil {
		return nil, nil, err
	}

	var fetchURL *url.URL

	if registry == ServiceProvider && c.BaseURL.String() == DefaultBaseURL {
		experimentalURL, _ := url.Parse(experimentalBaseURL)
		fetchURL = experimentalURL.ResolveReference(u)
	} else {
		fetchURL = c.BaseURL.ResolveReference(u)
	}

	resp, err := c.HTTP.Get(fetchURL.String())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var s Registry
	s, err = newRegistry(registry, json)

	if err != nil {
		return json, nil, err
	}

	return json, s, nil
}

func (c *Client) freshenFromCache(registry RegistryType) {
	if c.Cache.State(registry.Filename()) == cache.ShouldReload {
		c.reloadFromCache(registry)
	}
}

func (c *Client) reloadFromCache(registry RegistryType) error {
	json, err := c.Cache.Load(registry.Filename())

	if err != nil {
		return err
	}

	var s Registry
	s, err = newRegistry(registry, json)

	if err != nil {
		return err
	}

	c.registries[registry] = s

	return nil
}

func newRegistry(registry RegistryType, json []byte) (Registry, error) {
	var s Registry
	var err error

	switch registry {
	case ASN:
		s, err = NewASNRegistry(json)
	case DNS:
		s, err = NewDNSRegistry(json)
	case IPv4:
		s, err = NewNetRegistry(json, 4)
	case IPv6:
		s, err = NewNetRegistry(json, 6)
	case ServiceProvider:
		s, err = NewServiceProviderRegistry(json)
	default:
		panic("Unknown Registrytype")
	}

	return s, err
}

// DownloadAll downloads all four bootstrap registry files ({asn,dns,ipv4,ipv6}.json).
//
// On success, all four Registries are refreshed. Use ASN(), DNS(), IPv4(), and IPv6() to access them.
//
// This does not download the experimental ServiceProvider registry yet.
func (c *Client) DownloadAll() error {
	registryTypes := []RegistryType{ASN, DNS, IPv4, IPv6}

	for _, registryType := range registryTypes {
		err := c.Download(registryType)
		if err != nil {
			return err
		}
	}
	return nil
}

// Lookup returns the RDAP base URLs for the query |input| in the registry type |registry|.
//
// This function will download a Service Registry file if necessary, with each
// file cached for 24 hours by default. To adjust the cache duration, use
// c.Cache.SetTimeout().
func (c *Client) Lookup(registry RegistryType, input string) (*Result, error) {
	var forceDownload bool = false
	if c.Cache.State(registry.Filename()) == cache.ShouldReload {
		if err := c.reloadFromCache(registry); err != nil {
			forceDownload = true
		}
	}

	if c.registries[registry] == nil || forceDownload {
		err := c.Download(registry)
		if err != nil {
			return nil, err
		}
	}

	var result *Result
	result, err := c.registries[registry].Lookup(input)

	return result, err
}

// ASN returns the current ASN Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) ASN() *ASNRegistry {
	c.freshenFromCache(ServiceProvider)

	s, _ := c.registries[ASN].(*ASNRegistry)
	return s
}

//
// DNS returns the current DNS Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) DNS() *DNSRegistry {
	c.freshenFromCache(ServiceProvider)

	s, _ := c.registries[DNS].(*DNSRegistry)
	return s
}

// IPv4 returns the current IPv4 Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) IPv4() *NetRegistry {
	c.freshenFromCache(ServiceProvider)

	s, _ := c.registries[IPv4].(*NetRegistry)
	return s
}

// IPv6 returns the current IPv6 Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) IPv6() *NetRegistry {
	c.freshenFromCache(ServiceProvider)

	s, _ := c.registries[IPv6].(*NetRegistry)
	return s
}

// ServiceProvider returns the current ServiceProvider Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) ServiceProvider() *ServiceProviderRegistry {
	c.freshenFromCache(ServiceProvider)

	s, _ := c.registries[ServiceProvider].(*ServiceProviderRegistry)
	return s
}

// Filename returns the JSON document filename: One of {asn,dns,ipv4,ipv6,service_provider}.json.
func (r RegistryType) Filename() string {
	switch r {
	case ASN:
		return "asn.json"
	case DNS:
		return "dns.json"
	case IPv4:
		return "ipv4.json"
	case IPv6:
		return "ipv6.json"
	case ServiceProvider:
		// This is a guess and will need fixing to match whatever IANA chooses.
		return "service_provider.json"
	default:
		panic("Unknown RegistryType")
	}
}
