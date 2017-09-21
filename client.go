// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/openrdap/rdap/bootstrap"
)

type Client struct {
	HTTP      *http.Client
	Bootstrap *bootstrap.Client

	// A list of contact roles to additinoal fetch & merge into resuts.
	// Default is no extra fetches, set special string "all" to try fetching all of them
	FetchRoles []string

	// Timeout to complete a fetch.
	// Includes bootstrapping, sub-fetches...
	// Default is no timeout
	Timeout time.Duration

	// Enable experimental Service Provider bootstrapping.
	//
	// Default is disabled.
	ServiceProviderExperiment bool

	// Optional callback function for verbose messages.
	Verbose func(text string)
}

func NewClient() *Client {
	return &Client{
		HTTP:       &http.Client{},
		Bootstrap:  bootstrap.NewClient(),
		FetchRoles: []string{},
		Verbose:    defaultVerboseFunc,
	}
}

func (c *Client) Do(req *Request) (*Response, error) {
	// Bad query?
	if req == nil {
		return nil, &ClientError{Text: "nil Request"}
	}

	// Init HTTP client?
	if c.HTTP == nil {
		c.HTTP = &http.Client{}
	}

	// Init Bootstrap client?
	if c.Bootstrap == nil {
		c.Bootstrap = bootstrap.NewClient()
	}

	// Init Verbose callback?
	if c.Verbose == nil {
		c.Verbose = defaultVerboseFunc
	}

	c.Verbose(fmt.Sprintf("client: running request type %s (client: text=%s url=%s)",
		req.Type,
		req.Query,
		req.URL()))

	// Need to bootstrap the query?
	if req.Server == nil {
		var bootstrapType *bootstrap.RegistryType = bootstrapTypeFor(req)

		if bootstrapType == nil || (*bootstrapType == bootstrap.ServiceProvider && !c.ServiceProviderExperiment) {
			return nil, &ClientError{
				Type: BootstrapNotSupported,
				Text: fmt.Sprintf("Cannot run query type '%s' without a server URL, "+
					"the server must be specified",
					req.Type),
			}
		}

		c.Verbose(fmt.Sprintf("client: bootstrap required, running..."))

		var result *bootstrap.Result
		var err error

		result, err = c.Bootstrap.Lookup(*bootstrapType, req.Query)

		if err != nil {
			return nil, err
		}

		fmt.Printf("ok bootstrap ok %v\n", *result)
	}

	// main issues are raw response, timeout working correctly, *Response or interface{}?
	return nil, nil
}

func (c *Client) QueryDomain(domain string) (*Domain, error) {
	return nil, nil
}
func (c *Client) QueryAutnum(autnum string) (*Response, error) {
	return nil, nil
}
func (c *Client) QueryIP(ip string) (*Domain, error) {
	return nil, nil
}

func defaultVerboseFunc(text string) {
}

func bootstrapTypeFor(req *Request) *bootstrap.RegistryType {
	var b *bootstrap.RegistryType

	switch req.Type {
	case DomainRequest:
		*b = bootstrap.DNS
	case AutnumRequest:
		*b = bootstrap.ASN
	case EntityRequest:
		*b = bootstrap.ServiceProvider
	case IPRequest:
		if strings.Contains(req.Query, ":") {
			*b = bootstrap.IPv6
		} else {
			*b = bootstrap.IPv4
		}
	default:
		b = nil
	}

	return b
}
