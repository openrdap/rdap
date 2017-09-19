// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"net/http"
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
}

func NewClient() *Client {
	return &Client{
		HTTP:       &http.Client{},
		Bootstrap:  bootstrap.NewClient(),
		FetchRoles: []string{},
	}
}

func (c *Client) Query(q *Query) (*Response, error) {
	if q == nil {
		return nil, &ClientError{Text: "nil query"}
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
