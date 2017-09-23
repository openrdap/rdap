// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"net"
	"net/url"
	"testing"
)

const (
	ExampleServer = "https://test.rdap.example/rdap"
)

func testRequestURL(t *testing.T, r *Request, path string) {
	expectedURL := ExampleServer + "/" + path

	server, _ := url.Parse(ExampleServer)
	r2 := r.WithServer(server)

	actualURL := r2.URL()

	if actualURL == nil {
		t.Errorf("Nil URL\n")
		return
	}

	if actualURL.String() != expectedURL {
		t.Errorf("Got URL %s, expected %s\n", actualURL.String(), expectedURL)
		return
	}
}

func TestNewAutnumRequest(t *testing.T) {
	r := NewAutnumRequest(123456)

	testRequestURL(t, r, "autnum/123456")
}

func TestNewIPv4Request(t *testing.T) {
	r := NewIPRequest(net.ParseIP("192.0.2.0"))

	testRequestURL(t, r, "ip/192.0.2.0")
}

func TestNewIPv6Request(t *testing.T) {
	r := NewIPRequest(net.ParseIP("2001:DB8::a"))

	testRequestURL(t, r, "ip/2001:db8::a")
}

func TestNewIPv4NetRequest(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.0.2.0/24")
	r := NewIPNetRequest(ipNet)

	testRequestURL(t, r, "ip/192.0.2.0/24")
}

func TestNewIPv6NetRequest(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("2001:DB8::1/128")
	r := NewIPNetRequest(ipNet)

	testRequestURL(t, r, "ip/2001:db8::1/128")
}

func TestNewNameserverRequest(t *testing.T) {
	r := NewNameserverRequest("ns.example")

	testRequestURL(t, r, "nameserver/ns.example")
}

func TestNewDomainRequest(t *testing.T) {
	tests := []struct {
		Query        string
		ExpectedPath string
	}{
		{"example.com", "domain/example.com"},
		{"example/../com", "domain/example%2F..%2Fcom"},
	}

	for _, test := range tests {
		r := NewDomainRequest(test.Query)

		testRequestURL(t, r, test.ExpectedPath)
	}
}

func TestNewEntityRequest(t *testing.T) {
	tests := []struct {
		Query        string
		ExpectedPath string
	}{
		{"MY-HANDLE", "entity/MY-HANDLE"},
		{"MY-HANDLE/../com", "entity/MY-HANDLE%2F..%2Fcom"},
	}

	for _, test := range tests {
		r := NewEntityRequest(test.Query)

		testRequestURL(t, r, test.ExpectedPath)
	}
}

func TestNewHelpRequest(t *testing.T) {
	r := NewHelpRequest()

	testRequestURL(t, r, "help")
}

func TestNewRawRequest(t *testing.T) {
	urlString := "https://example.com/domain/example.com"
	url, _ := url.Parse(urlString)
	r := NewRawRequest(url)

	actualURL := r.URL()
	if actualURL.String() != urlString {
		t.Errorf("Raw query for %s got %s, expected %s\n",
			urlString,
			actualURL.String(),
			urlString,
		)
	}
}

func TestNewSearchRequests(t *testing.T) {
	tests := []struct {
		RequestType  RequestType
		Query        string
		ExpectedPath string
	}{
		{
			DomainSearchRequest,
			"example*.com&x=1",
			"domains?name=example%2A.com%26x%3D1",
		},
		{
			DomainSearchByNameserverRequest,
			"example*.com&x=1",
			"domains?nsLdhName=example%2A.com%26x%3D1",
		},
		{
			DomainSearchByNameserverIPRequest,
			"192.0.2.*.com&x=1",
			"domains?nsIp=192.0.2.%2A.com%26x%3D1",
		},
		{
			NameserverSearchRequest,
			"ns1.example*.com&x=1",
			"nameservers?name=ns1.example%2A.com%26x%3D1",
		},
		{
			NameserverSearchByNameserverIPRequest,
			"192.0.2.*.com&x=1",
			"nameservers?ip=192.0.2.%2A.com%26x%3D1",
		},
		{
			EntitySearchRequest,
			"MY-FN*&x=1",
			"entities?fn=MY-FN%2A%26x%3D1",
		},
		{
			EntitySearchByHandleRequest,
			"MY-HANDLE*&x=1",
			"entities?handle=MY-HANDLE%2A%26x%3D1",
		},
	}

	for _, test := range tests {
		r := NewRequest(test.RequestType, test.Query)
		testRequestURL(t, r, test.ExpectedPath)
	}
}

func TestRequestURLConstruction(t *testing.T) {
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

		r := NewDomainRequest(test.Domain)
		r2 := r.WithServer(server)
		actualURL := r2.URL()

		if actualURL == nil {
			t.Errorf("nil url")
		} else if actualURL.String() != test.ExpectedURL {
			t.Errorf("Got URL %s, expected %s\n", actualURL.String(), test.ExpectedURL)
		}
	}
}

func TestNewAutoRequest(t *testing.T) {
	tests := []struct {
		Query        string
		ExpectedType RequestType
	}{
		{"http://example.com/", DomainRequest},
		{"https://example.com", DomainRequest},

		{"http://example.com/domain/example.com", RawRequest},
		{"https://example.com/domain/example.com", RawRequest},

		{"192.0.2.0", IPRequest},
		{"192.0.2.0/24", IPRequest},
		{"2001:db8::", IPRequest},
		{"2001:db8::/128", IPRequest},

		{"AS1", AutnumRequest},
		{"as12", AutnumRequest},
		{"aS123", AutnumRequest},
		{"1234", AutnumRequest},

		{"example.com", DomainRequest},

		{"example", EntityRequest},
	}

	for _, test := range tests {
		r := NewAutoRequest(test.Query)

		if r.Type != test.ExpectedType {
			t.Errorf("AutoQuery detected %s as type %s, expected %s\n",
				r.Query,
				r.Type,
				test.ExpectedType)
		}
	}
}
