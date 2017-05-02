// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package bootstrap implements Registration Data Access Protocol (RDAP) bootstrapping.
//
// All RDAP queries are handled by an RDAP server. To help clients discover
// authoriative RDAP servers, IANA publishes a Service Registry
// (https://data.iana.org/rdap) for three query types: Domain names, IP
// addresses, and Autonomous Systems.
//
// This module implements RDAP bootstrapping: The process of taking a query and
// returning a list of authoriative RDAP servers which may answer it.
//
// Basic usage:
//   var urls []*url.URL
//
//   b := bootstrap.NewClient()
//   urls, err := b.Lookup(bootstrap.DNS, "google.cz")  // Downloads https://data.iana.org/rdap/dns.json automatically.
//
//   if err == nil {
//     for _, url := range urls {
//       fmt.Println(url.String()) // Prints https://rdap.nic.cz.
//     }
//   }
//
// Download and list a RDAP Service Registry:
//   b := bootstrap.NewClient()
//
//   // Before you can use a Registry, you need to download it.
//   err := b.Download(bootstrap.IPv6) // Downloads https://data.iana.org/rdap/ipv6.json.
//
//   if err == nil {
//     ipv6 := b.IPv6()
//
//     // Print IPv6 networks listed in the IPv6 service registry.
//     for _, net := range ipv6.Nets() {
//       fmt.Println(net.String()) // e.g. prints "2001:4200::/23".
//     }
//   }
//
// A bootstrap.Client caches the Service Registry files in memory for both performance, and courtesy to data.iana.org. The functions which make network requests are:
//   - Download()      - download one Service Registry file.
//   - DownloadAll()   - download all four Service Registry files.
//
//   - Lookup() - download one Service Registry file if missing, or if the cached file is over (by default) 24 hours old.
// Lookup() is intended for repeated usage: A long lived
// bootstrap.Client will download each of {asn,dns,ipv4,ipv6}.json once per 24 hours only,
// regardless of the number of calls made to Lookup(). You can still refresh them manually using Download(), if required.
//
// As well as the default memory cache, bootstrap.Client also supports caching
// the Service Registry files on disk. The default cache location is
// $HOME/.openrdap/{asn,dns,ipv4,ipv6}.json.
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
	// Default location of the IANA bootstrap files.
	DefaultBaseURL      = "https://data.iana.org/rdap/"

	// Default cache timeout for bootstrap files.
	DefaultCacheTimeout = time.Hour * 24

	// Location of the experimental service_provider.json.
	experimentalBaseURL = "https://www.openrdap.org/rdap/"
)

// Client implements an RDAP bootstrap client.
type Client struct {
	HTTP    *http.Client        // HTTP client.
	BaseURL *url.URL            // Base URL of the Service Registry directory. Default is DefaultBaseURL.
	Cache   cache.RegistryCache // Service Registry cache. Default is a MemoryCache.

	registries map[RegistryType]Registry
}

// A Registry performs RDAP bootstrapping.
type Registry interface {
	Lookup(input string) (*Result, error)
}

// Result represents the result of bootstrapping a single query.
type Result struct {
	// Query looked up in the bootstrap registry.
	//
	// This includes any canonicalisation to match the bootstrap registry
	// format. e.g. lowercasing of domain names, and removal of "AS" from AS
	// numbers.
	Query string

	// Matching bootstrap entry. Empty string if no match.
	Entry string

	// List of base RDAP URLs.
	URLs  []*url.URL
}

// NewClient creates a new bootstrap Client.
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
//     for _, tld := range dns.TLDs() {
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

func (c *Client) reloadFromCache(registry RegistryType) error {
	json, isNew, err := c.Cache.Load(registry.Filename())

	if err != nil {
		return err
	} else if !isNew {
		return nil
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
// This does not refresh the experimental ServiceProvider registry yet.
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

// Returns true if the RegistryType registry is stale (cache time expired) or
// missing, and thus should be downloaded again.
//
// Stale registries remain accessible for use.
func (c *Client) IsStale(registry RegistryType) bool {
	c.reloadFromCache(ASN)
	return c.Cache.IsStale(registry.Filename())
}

func (c *Client) Lookup(registry RegistryType, query string) (*Result, error) {
	if c.IsStale(registry) {
		err := c.Download(registry)
		if err != nil {
			return nil, err
		}
	}

	var result *Result
	result, err := c.registries[registry].Lookup(query)

	return result, err
}

// ASN returns the current ASN Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) ASN() *ASNRegistry {
	c.reloadFromCache(ASN)

	s, _ := c.registries[ASN].(*ASNRegistry)
	return s
}

//
// DNS returns the current DNS Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) DNS() *DNSRegistry {
	c.reloadFromCache(DNS)

	s, _ := c.registries[DNS].(*DNSRegistry)
	return s
}

// IPv4 returns the current IPv4 Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) IPv4() *NetRegistry {
	c.reloadFromCache(IPv4)

	s, _ := c.registries[IPv4].(*NetRegistry)
	return s
}

// IPv6 returns the current IPv6 Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) IPv6() *NetRegistry {
	c.reloadFromCache(IPv6)

	s, _ := c.registries[IPv6].(*NetRegistry)
	return s
}

// ServiceProvider returns the current ServiceProvider Registry (or nil if the registry file hasn't been Download()ed).
//
// This function never initiates a network transfer.
func (c *Client) ServiceProvider() *ServiceProviderRegistry {
	c.reloadFromCache(ServiceProvider)

	s, _ := c.registries[ServiceProvider].(*ServiceProviderRegistry)
	return s
}

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
