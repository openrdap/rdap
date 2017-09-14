// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"net/http"

	"github.com/openrdap/rdap/bootstrap"
)

type Client struct {
	HTTP      *http.Client
	Bootstrap *bootstrap.Client

	FetchRoles []string
}

func NewClient() *Client {
	return &Client{
		HTTP:       &http.Client{},
		Bootstrap:  bootstrap.NewClient(),
		FetchRoles: []string{"all"},
	}
}

func (c *Client) Query(q *Query) (*Response, error) {
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
