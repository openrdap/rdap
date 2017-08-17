// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package client

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// A SearchType specifies an RDAP search query type. Used with NewSearchQuery().
//
// Some servers may choose not to support these advanced query types.
type SearchType int

const (
	DomainSearch                   SearchType = iota // Example query: example*.com
	DomainSearchByNameserver                         // Example query: ns1.example*.com
	DomainSearchByNameserverIP                       // Example query: 192.0.2.0
	NameserverSearch                                 // Example query: ns1.example*.com
	NameserverSearchByNameserverIP                   // Example query: 2001:db8::
	EntitySearch                                     // Example query: Bobby%20Joe*
	EntitySearchByHandle                             // Example query: CID-40*
)

// A Query represents an RDAP query. Queries are executed by a Client.
//
// To execute a Query, an RDAP server is required. The RDAP servers for Autnum,
// IP, IPNet, and Domain queries are determined automatically by the Client, in
// a process called bootstrapping. For other query types, you must explicity
// specify an RDAP server.
//
// See https://tools.ietf.org/html/rfc7482 for more information on RDAP query types.
type Query struct {
	url *url.URL

	queryType   string
	queryPath   string
	queryValues url.Values
	queryText   string
}

// NewIPQuery creates a new Query for the IP address ip.
//
// The RDAP server to query will be automatically determined during query execution.
// To query a specific server, use UsingServer().
func NewIPQuery(ip net.IP) *Query {
	queryText := ip.String()

	return &Query{
		queryType: queryTypeForIP(ip),
		queryPath: fmt.Sprintf("ip/%s", queryText),
		queryText: queryText,
	}
}

// NewIPNetQuery creates a new Query for the IP network net.
//
// The RDAP server to query will be automatically determined during query execution.
// To query a specific server, use UsingServer().
func NewIPNetQuery(net *net.IPNet) *Query {
	queryText := net.String()

	return &Query{
		queryType: queryTypeForIP(net.IP),

		queryPath: fmt.Sprintf("ip/%s", queryText),
		queryText: queryText,
	}
}

// NewAutnumQuery creates a new Query for the Autonomous System (AS) number autnum, e.g. 5400.
//
// The RDAP server to query will be automatically determined during query execution.
// To query a specific server, use UsingServer().
func NewAutnumQuery(autnum uint32) *Query {
	return &Query{
		queryType: "autnum",

		queryPath: fmt.Sprintf("autnum/%d", autnum),
		queryText: fmt.Sprintf("%d", autnum),
	}
}

// NewDomainQuery creates a new Query for a domain name, e.g. "google.cz".
//
// The RDAP server to query will be automatically determined during query execution.
// To query a specific server, use UsingServer().
func NewDomainQuery(domain string) *Query {
	return &Query{
		queryType: "dns",

		queryPath: fmt.Sprintf("domain/%s", escapePath(domain)),
		queryText: domain,
	}
}

// NewNameserverQuery creates a new Query for a nameserver, e.g. "a.ns.nic.cz".
//
// To execute this query type, you must specify the RDAP server using UsingServer().
func NewNameserverQuery(nameserver string) *Query {
	return &Query{
		queryType: "nameserver",

		queryPath: fmt.Sprintf("nameserver/%s", escapePath(nameserver)),
		queryText: nameserver,
	}
}

// NewEntityQuery creates a new Query for an entity name, e.g. "TEST-NET-1".
//
// The RDAP server to query will be determined automatically, but only if the
// handle contains a service provider tag:
//
// The RDAP server for a handle with a tag (e.g. the handle 86413629~VRSN
// contains VeriSign's VRSN tag) is determined automatically. Otherwise, you
// must specify the RDAP server using UsingServer().
func NewEntityQuery(handle string) *Query {
	return &Query{
		queryType: "entity",

		queryPath: fmt.Sprintf("entity/%s", escapePath(handle)),
		queryText: handle,
	}
}

// NewHelpQuery creates a new Query for a server's help information.
//
// To execute this query type, you must specify the RDAP server using UsingServer().
func NewHelpQuery() *Query {
	return &Query{
		queryType: "help",

		queryPath: "help",
	}
}

// NewURLQuery creates a new Query, to query a known RDAP server URL, e.g. "https://rdap.nic.cz/domain/google.cz".
//
// The URL can return any RDAP response type.
func NewURLQuery(url *url.URL) *Query {
	return &Query{
		url:       &*url,
		queryType: "url",
	}
}

func newDomainSearchQuery(domainSearchPattern string) *Query {
	v := url.Values{}
	v.Add("name", domainSearchPattern)

	return &Query{
		queryType: "domain-search",

		queryPath:   "domains",
		queryValues: v,
		queryText:   domainSearchPattern,
	}
}

func newDomainSearchByNameserverQuery(nameserverSearchPattern string) *Query {
	v := url.Values{}
	v.Add("nsLdhName", nameserverSearchPattern)

	return &Query{
		queryType: "domain-search-by-nameserver",

		queryPath:   "domains",
		queryValues: v,
		queryText:   nameserverSearchPattern,
	}
}

func newDomainSearchByNameserverIPQuery(ipSearchPattern string) *Query {
	v := url.Values{}
	v.Add("nsIp", ipSearchPattern)

	return &Query{
		queryType: "domain-search-by-nameserver-ip",

		queryPath:   "domains",
		queryValues: v,
		queryText:   ipSearchPattern,
	}
}

// NewSearchQuery creates a new search Query of type searchType, for the search pattern searchPattern.
//
// See https://tools.ietf.org/html/rfc7482#section-3.2 for search pattern examples.
//
// To execute this query type, you must specify the RDAP server using UsingServer().
func NewSearchQuery(searchType SearchType, searchPattern string) *Query {
	switch searchType {
	case DomainSearch:
		return newDomainSearchQuery(searchPattern)
	case DomainSearchByNameserver:
		return newDomainSearchByNameserverQuery(searchPattern)
	case DomainSearchByNameserverIP:
		return newDomainSearchByNameserverIPQuery(searchPattern)
	case NameserverSearch:
		return newNameserverSearchQuery(searchPattern)
	case NameserverSearchByNameserverIP:
		return newNameserverSearchByNameserverIPQuery(searchPattern)
	case EntitySearch:
		return newEntitySearchQuery(searchPattern)
	case EntitySearchByHandle:
		return newEntitySearchByHandleQuery(searchPattern)
	}

	return nil
}

func newNameserverSearchQuery(nameserverSearchPattern string) *Query {
	v := url.Values{}
	v.Add("name", nameserverSearchPattern)

	return &Query{
		queryType: "nameserver-search",

		queryPath:   "nameservers",
		queryValues: v,
		queryText:   nameserverSearchPattern,
	}
}

func newNameserverSearchByNameserverIPQuery(ipSearchPattern string) *Query {
	v := url.Values{}
	v.Add("ip", ipSearchPattern)

	return &Query{
		queryType: "nameserver-search-by-ip",

		queryPath:   "nameservers",
		queryValues: v,
		queryText:   ipSearchPattern,
	}
}

func newEntitySearchQuery(entitySearchPattern string) *Query {
	v := url.Values{}
	v.Add("fn", entitySearchPattern)

	return &Query{
		queryType: "entity-search",

		queryPath:   "entities",
		queryValues: v,
		queryText:   entitySearchPattern,
	}
}

func newEntitySearchByHandleQuery(handleSearchPattern string) *Query {
	v := url.Values{}
	v.Add("handle", handleSearchPattern)

	u := &url.URL{}
	u.Path = "entities"
	u.RawQuery = v.Encode()

	return &Query{
		queryType: "entity-search-by-handle",

		queryPath:   "entities",
		queryValues: v,
		queryText:   handleSearchPattern,
	}
}

func queryTypeForIP(ip net.IP) string {
	if ip.To16() != nil {
		return "ipv6"
	}

	return "ipv4"
}

// NewAutoQuery creates a Query by guessing the Query type required for queryText.
//
// The types of queries guessed, and example inputs are:
//
//  - NewDomainQuery() - example.com, https://example.com, http://example.com/
//  - NewURLQuery()    - https://example.com/domain/example2.com
//  - NewIPQuery()     - 192.0.2.0, 2001:db8::
//  - NewIPNetQuery()  - 192.0.2.0/24, 2001:db8::/128
//  - NewAutnumQuery() - 5400
//  - NewEntityQuery() - all other queries.
//
// Use Type() to determine the Query Type chosen.
func NewAutoQuery(queryText string) *Query {
	// Full RDAP URL?
	fullURL, err := url.Parse(queryText)
	if err == nil && (fullURL.Scheme == "http" || fullURL.Scheme == "https") {
		// Parse "http://example.com/" as a domain query for convenience.
		if fullURL.Path == "" || fullURL.Path == "/" {
			return NewDomainQuery(fullURL.Host)
		}

		return NewURLQuery(fullURL)
	}

	// IP address?
	ip := net.ParseIP(queryText)
	if ip != nil {
		return NewIPQuery(ip)
	}

	// IP network?
	_, ipNet, err := net.ParseCIDR(queryText)
	if ipNet != nil {
		return NewIPNetQuery(ipNet)
	}

	// AS number? (formats: AS1234, as1234, 1234).
	autnum, err := parseAutnum(queryText)
	if err == nil {
		return NewAutnumQuery(autnum)
	}

	// Looks like a domain name?
	if strings.Contains(queryText, ".") {
		return NewDomainQuery(queryText)
	}

	// Otherwise call it an entity query.
	return NewEntityQuery(queryText)
}

func parseAutnum(autnum string) (uint32, error) {
	autnum = strings.ToUpper(autnum)
	autnum = strings.TrimPrefix(autnum, "AS")
	result, err := strconv.ParseUint(autnum, 10, 32)

	if err != nil {
		return 0, err
	}

	return uint32(result), nil
}

func escapePath(text string) string {
	var escaped []byte

	for i := 0; i < len(text); i++ {
		b := text[i]

		if !shouldPathEscape(b) {
			escaped = append(escaped, b)
		} else {
			escaped = append(escaped, '%',
				"0123456789ABCDEF"[b>>4],
				"0123456789ABCDEF"[b&0xF],
			)
		}
	}

	return string(escaped)
}

func shouldPathEscape(b byte) bool {
	if ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || ('0' <= b && b <= '9') {
		return false
	}

	switch b {
	case '-', '_', '.', '~', '$', '&', '+', ':', '=', '@':
		return false
	}

	return true
}

// HasServer returns true if the Query has a server specified.
//
// If HasServer() returns true, then the URL() is non-nil, otherwise nil.
//
// Bootstrappable queries can be executed without a server set. Otherwise, the
// server must be specified with UsingServer().
func (q *Query) HasServer() bool {
	return q.url != nil
}

func (q *Query) bootstrapInfo() (string, string) {
	return q.queryType, q.queryText
}

// URL returns the Query URL.
//
// Returns nil if the Query doesn't have a server specified, see HasServer().
func (q *Query) URL() *url.URL {
	if q.url == nil {
		return nil
	}

	return &*q.url
}

// Type returns a string describing the RDAP query type. This is one of:
//
//  - ip
//  - domain
//  - autnum
//  - nameserver
//  - entity
//  - help
//  - url
//  - entity-search-by-handle
//  - domain-search
//  - domain-search-by-nameserver
//  - domain-search-by-nameserver-ip
//  - nameserver-search
//  - nameserver-search-by-ip
//  - entity-search
//  - entity-search-by-handle
//
// Queries created by NewIPQuery()/NewIPNetQuery() both have the type "ip".
//
// For URL queries (see NewURLQuery()), no attempt is made to parse the URL to
// determine the query type, so these return type "url".
func (q *Query) Type() string {
	switch q.queryType {
	case "ipv4", "ipv6":
		return "ip"
	case "dns":
		return "domain"
	default:
		return q.queryType
	}
}

// Text returns the query text, (the domain name, IP address, IP
// network, AS number, nameserver, or search pattern being queried) as a string.
//
// For URL queries, no attempt is made to parse the query URL, and thus Text()
// returns an empty string.
func (q *Query) Text() string {
	return q.queryText
}

func (q *Query) requestURI() string {
	if q.HasServer() {
		return q.requestURI()
	}

	var requestURI string = q.queryPath
	var queryString string = q.queryValues.Encode()

	if queryString != "" {
		requestURI += "?" + queryString
	}

	return requestURI
}

// UsingServer returns a clone of the Query, with the server set to
// server. The original query is left unchanged.
//
// The created Query retains the same Type() and Text() values, has a URL(),
// and HasServer() == true.
//
// Returns an error if the server is already specified.
func (q *Query) UsingServer(server *url.URL) (*Query, error) {
	if server == nil {
		return nil, errors.New("Cannot use a nil URL")
	} else if q.HasServer() {
		return nil, errors.New("Cannot set Query server, already set.")
	}

	u := &*server
	u.RawQuery = ""
	u.Fragment = ""
	urlString := u.String()

	if len(urlString) == 0 || urlString[len(urlString)-1] != '/' {
		urlString += "/"
	}

	urlString += q.requestURI()

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	return &Query{
		url:       u,
		queryType: q.queryType,
		queryText: q.queryText,
	}, nil
}

// UseServerURL is alternative to UsingServer(), with the server URL
// specified as a string (e.g. "https://rdap.nic.cz").
//
// Internally the server string is parsed, then UsingServer() is called.
func (q *Query) UsingServerURL(server string) (*Query, error) {
	u, err := url.Parse(server)

	if err != nil {
		return nil, err
	}

	return q.UsingServer(u)
}
