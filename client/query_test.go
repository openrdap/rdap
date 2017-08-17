// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package client

import (
	"net"
	"net/url"
	"testing"
)

const (
	ExampleDotCom = "https://example.com/rdap"
)

func TestNewIPQuery(t *testing.T) {
	q := NewIPQuery(net.ParseIP("192.0.2.0"))

	if q.requestURI() != "ip/192.0.2.0" {
		t.Errorf("Unexpected path: %s", q.requestURI())
	}
}

func TestNewIPNet4Query(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.0.2.0/24")
	q := NewIPNetQuery(ipNet)

	if q.requestURI() != "ip/192.0.2.0/24" {
		t.Errorf("Unexpected path: %s", q.requestURI())
	}
}

func TestNewIPNet6Query(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("2001:db8::/128")
	q := NewIPNetQuery(ipNet)

	if q.requestURI() != "ip/2001:db8::/128" {
		t.Errorf("Unexpected path: %s", q.requestURI())
	}
}

func TestNewDomainQuery(t *testing.T) {
	tests := []struct {
		QueryText    string
		ExpectedPath string
	}{
		{"example.com", "domain/example.com"},
		{"example/../com", "domain/example%2F..%2Fcom"},
	}

	for _, test := range tests {
		query := NewDomainQuery(test.QueryText)

		if query.requestURI() != test.ExpectedPath {
			t.Errorf("Domain path for %s is %s, expected %s\n",
				test.QueryText,
				query.requestURI(),
				test.ExpectedPath)
		}
	}
}

func TestNewEntityQuery(t *testing.T) {
	tests := []struct {
		QueryText    string
		ExpectedPath string
	}{
		{"MY-HANDLE", "entity/MY-HANDLE"},
		{"MY-HANDLE/../com", "entity/MY-HANDLE%2F..%2Fcom"},
	}

	for _, test := range tests {
		query := NewEntityQuery(test.QueryText)

		if query.requestURI() != test.ExpectedPath {
			t.Errorf("Entity path for %s is %s, expected %s\n",
				test.QueryText,
				query.requestURI(),
				test.ExpectedPath)
		}
	}
}

func TestNewURLQuery(t *testing.T) {
	urlString := "https://example.com/domain/example.com"
	url, _ := url.Parse(urlString)
	query := NewURLQuery(url)

	if query.URL().String() != urlString {
		t.Errorf("URL query for %s got %s, expected %s\n",
			urlString,
			query.URL().String(),
			urlString,
		)
	}
}

func TestSearchQueries(t *testing.T) {
	tests := []struct {
		SearchType   SearchType
		QueryText    string
		ExpectedPath string
	}{
		{
			DomainSearch,
			"example*.com&x=1",
			"domains?name=example%2A.com%26x%3D1",
		},
		{
			DomainSearchByNameserver,
			"example*.com&x=1",
			"domains?nsLdhName=example%2A.com%26x%3D1",
		},
		{
			DomainSearchByNameserverIP,
			"192.0.2.*.com&x=1",
			"domains?nsIp=192.0.2.%2A.com%26x%3D1",
		},
		{
			NameserverSearch,
			"ns1.example*.com&x=1",
			"nameservers?name=ns1.example%2A.com%26x%3D1",
		},
		{
			NameserverSearchByNameserverIP,
			"192.0.2.*.com&x=1",
			"nameservers?ip=192.0.2.%2A.com%26x%3D1",
		},
		{
			EntitySearch,
			"MY-FN*&x=1",
			"entities?fn=MY-FN%2A%26x%3D1",
		},
		{
			EntitySearchByHandle,
			"MY-HANDLE*&x=1",
			"entities?handle=MY-HANDLE%2A%26x%3D1",
		},
	}

	for _, test := range tests {
		expectedURL := ExampleDotCom + "/" + test.ExpectedPath

		query, _ := NewSearchQuery(test.SearchType, test.QueryText).UsingServerURL(ExampleDotCom)

		if query.URL().String() != expectedURL {
			t.Errorf("For query %s expected %s, got %s\n",
				test.QueryText,
				expectedURL,
				query.URL().String(),
			)
		}
	}
}

func TestUsingServer(t *testing.T) {
	tests := []struct {
		Server      string
		Domain      string
		ExpectedURL string
	}{
		{"http://example.com", "example.org", "http://example.com/domain/example.org"},
		{"http://example.com/", "example.org", "http://example.com/domain/example.org"},
		{"http://example.com/1/", "example.org", "http://example.com/1/domain/example.org"},
		{"http://example.com/1/", "example.org/", "http://example.com/1/domain/example.org%2F"},
		{"https://example.com/x/", "example.com*.&x=1?z=1", "https://example.com/x/domain/example.com%2A.&x=1%3Fz=1"},
	}

	for _, test := range tests {
		server, _ := url.Parse(test.Server)

		q := NewDomainQuery(test.Domain)
		q2, err := q.UsingServer(server)

		if err != nil {
			t.Errorf("Error %s for %s\n", err, test.Server)
		} else if q2.URL().String() != test.ExpectedURL {
			t.Errorf("Got URL %s, expected %s\n", q2.URL().String(), test.ExpectedURL)
		}

		if !q2.HasServer() {
			t.Errorf("Query %s not complete after UsingServer\n", test.Server)
		}
	}
}

func TestAutoQueryTypes(t *testing.T) {
	tests := []struct {
		QueryText    string
		ExpectedType string
	}{
		{"http://example.com/", "domain"},
		{"https://example.com", "domain"},

		{"http://example.com/domain/example.com", "url"},
		{"https://example.com/domain/example.com", "url"},

		{"192.0.2.0", "ip"},
		{"192.0.2.0/24", "ip"},
		{"2001:db8::", "ip"},
		{"2001:db8::/128", "ip"},

		{"AS1", "autnum"},
		{"as12", "autnum"},
		{"aS123", "autnum"},
		{"1234", "autnum"},

		{"example.com", "domain"},

		{"example", "entity"},
	}

	for _, test := range tests {
		query := NewAutoQuery(test.QueryText)

		if query.Type() != test.ExpectedType {
			t.Errorf("AutoQuery detected %s as type %s, expected %s\n",
				test.QueryText,
				query.Type(),
				test.ExpectedType)
		}
	}
}
