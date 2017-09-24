// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/openrdap/rdap/bootstrap"
)

// Client implements an RDAP client.
//
// This client executes RDAP requests, and returns the responses as Go values.
//
// Quick usage:
//   client := &rdap.Client{}
//   domain, err := client.QueryDomain("google.cz")
//
//   if err == nil {
//     fmt.Printf("Handle=%s Domain=%s\n", domain.Handle, domain.LDHName)
//   }
// The QueryDomain(), QueryAutnum(), and QueryIP() methods all provide full contact information, and timeout after 30s.
//
// Normal usage:
//   // Query example.cz.
//   req := &rdap.Request{
//     Type: rdap.DomainRequest,
//     Query: "example.cz",
//   }
//
//   client := &rdap.Client{}
//   resp, err := client.Do(req)
//
//   if domain, ok := resp.Response.(*rdap.Domain); ok {
//     fmt.Printf("Handle=%s Domain=%s\n", domain.Handle, domain.LDHName)
//   }
//
// Advanced usage:
//
// This demonstrates custom FetchRoles, a custom Context, a custom HTTP client,
// a custom Bootstrapper, and a custom timeout.
//   // Nameserver query on rdap.nic.cz.
//   server, _ := url.Parse("https://rdap.nic.cz")
//   req := &rdap.Request{
//     Type: rdap.NameserverRequest,
//     Query: "a.ns.nic.cz",
//     FetchRoles: []string{"all"},
//     Timeout: time.Second * 45, // Custom timeout.
//
//     Server: server,
//   }
//
//   req = req.WithContext(ctx) // Custom context (see https://blog.golang.org/context).
//
//   client := &rdap.Client{}
//   client.HTTP = &http.Client{} // Custom HTTP client.
//   client.Bootstrap = &bootstrap.Client{} // Custom bootstapper.
//
//   resp, err := client.Do(req)
//
//   if ns, ok := resp.Response.(*rdap.Nameserver); ok {
//     fmt.Printf("Handle=%s Domain=%s\n", ns.Handle, ns.LDHName)
//   }
type Client struct {
	HTTP      *http.Client
	Bootstrap *bootstrap.Client

	ServiceProviderExperiment bool

	// Optional callback function for verbose messages.
	Verbose func(text string)
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
		c.Bootstrap = &bootstrap.Client{}
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

		question := &bootstrap.Question{
			RegistryType: *bootstrapType,
			Query:        req.Query,
		}
		question = question.WithContext(req.Context())

		var answer *bootstrap.Answer
		var err error

		answer, err = c.Bootstrap.Lookup(question)

		if err != nil {
			return nil, err
		}

		fmt.Printf("ok bootstrap ok %v\n", *answer)
	}

	// main issues are raw response, timeout working correctly, *Response or interface{}?
	return nil, nil
}

// QueryDomain makes an RDAP request for the |domain|.
//
// Full contact information (where available) is provided. The timeout is 30s.
func (c *Client) QueryDomain(domain string) (*Domain, error) {
	req := &Request{
		Type:  DomainRequest,
		Query: domain,
	}

	resp, err := c.doQuickRequest(req)
	if err != nil {
		return nil, err
	}

	if domain, ok := resp.Response.(*Domain); ok {
		return domain, nil
	}

	return nil, &ClientError{
		Type: WrongResponseType,
		Text: "The server didn't return an RDAP Domain response",
	}
}

func (c *Client) doQuickRequest(req *Request) (*Response, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	req = req.WithContext(ctx)
	resp, err := c.Do(req)

	return resp, err
}

// QueryAutnum makes an RDAP request for the Autonomous System Number (ASN) |autnum|.
//
// |autnum| is an ASN string, e.g. "AS2856" or "5400".
//
// Full contact information (where available) is provided. The timeout is 30s.
func (c *Client) QueryAutnum(autnum string) (*Autnum, error) {
	req := &Request{
		Type:  AutnumRequest,
		Query: autnum,
	}

	resp, err := c.doQuickRequest(req)
	if err != nil {
		return nil, err
	}

	if autnum, ok := resp.Response.(*Autnum); ok {
		return autnum, nil
	}

	return nil, &ClientError{
		Type: WrongResponseType,
		Text: "The server didn't return an RDAP Autnum response",
	}
}

// QueryIP makes an RDAP request for the IPv4/6 address |ip|, e.g. "192.0.2.0" or "2001:db8::".
//
// Full contact information (where available) is provided. The timeout is 30s.
func (c *Client) QueryIP(ip string) (*IPNetwork, error) {
	req := &Request{
		Type:  IPRequest,
		Query: ip,
	}

	resp, err := c.doQuickRequest(req)
	if err != nil {
		return nil, err
	}

	if ipNet, ok := resp.Response.(*IPNetwork); ok {
		return ipNet, nil
	}

	return nil, &ClientError{
		Type: WrongResponseType,
		Text: "The server didn't return an RDAP IPNetwork response",
	}
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
