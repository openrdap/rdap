// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package test

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/jarcoal/httpmock"
)

type TestDataset int

const (
	Bootstrap TestDataset = iota
	BootstrapExperimental
	BootstrapMalformed
	BootstrapComplex
	BootstrapHTTPError

	Responses
)

type response struct {
	Status int
	URL    string
	Body   string
}

var (
	responses     = buildResponses()
	activatedURLs = map[string]bool{}
)

// Start activates HTTP mocking and registers the responders for the given test
// dataset. It panics if two datasets register the same URL.
func Start(set TestDataset) {
	httpmock.Activate()

	for _, r := range responses[set] {
		if _, ok := activatedURLs[r.URL]; ok {
			log.Panicf("Test sets conflict on URL %s\n", r.URL)
		}

		activatedURLs[r.URL] = true

		httpmock.RegisterResponder("GET", r.URL,
			httpmock.NewStringResponder(r.Status, r.Body))
	}
}

// Finish deactivates HTTP mocking and clears the registered URLs, undoing
// Start.
func Finish() {
	activatedURLs = make(map[string]bool)
	httpmock.DeactivateAndReset()
}

// Get performs an HTTP GET and returns the response body, panicking on error.
func Get(url string) []byte {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	return data
}

// buildResponses loads every test dataset's mock HTTP responses from testdata.
func buildResponses() map[TestDataset][]response {
	r := map[TestDataset][]response{}

	load := func(set TestDataset, status int, url, filename string) {
		body := LoadFile(filename)
		r[set] = append(r[set], response{status, url, string(body)})
	}

	// Valid snapshot of the IANA bootstrap files.
	load(Bootstrap, http.StatusOK, "https://data.iana.org/rdap/asn.json", "bootstrap/asn.json")
	load(Bootstrap, http.StatusOK, "https://data.iana.org/rdap/dns.json", "bootstrap/dns.json")
	load(Bootstrap, http.StatusOK, "https://data.iana.org/rdap/ipv4.json", "bootstrap/ipv4.json")
	load(Bootstrap, http.StatusOK, "https://data.iana.org/rdap/ipv6.json", "bootstrap/ipv6.json")
	load(Bootstrap, http.StatusOK, "https://data.iana.org/rdap/object-tags.json", "bootstrap/object-tags.json")

	// Malformed bootstrap files.
	load(BootstrapMalformed, http.StatusOK, "https://www.example.org/dns_bad_services.json", "bootstrap_malformed/dns_bad_services.json")
	load(BootstrapMalformed, http.StatusOK, "https://www.example.org/dns_bad_url.json", "bootstrap_malformed/dns_bad_url.json")
	load(BootstrapMalformed, http.StatusOK, "https://www.example.org/dns_empty.json", "bootstrap_malformed/dns_empty.json")
	load(BootstrapMalformed, http.StatusOK, "https://www.example.org/dns_syntax_error.json", "bootstrap_malformed/dns_syntax_error.json")

	// Valid bootstrap files testing more features than yet used by IANA.
	load(BootstrapComplex, http.StatusOK, "https://rdap.example.org/dns.json", "bootstrap_complex/dns.json")

	// Bootstrap HTTP errors.
	load(BootstrapHTTPError, http.StatusNotFound, "https://data.iana.org/rdap/asn.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, http.StatusNotFound, "https://data.iana.org/rdap/dns.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, http.StatusNotFound, "https://data.iana.org/rdap/ipv4.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, http.StatusNotFound, "https://data.iana.org/rdap/ipv6.json", "bootstrap_http_error/404.html")

	// RDAP responses.
	load(Responses, http.StatusOK, "https://rdap.nic.cz/domain/example.cz", "rdap/rdap.nic.cz/domain-example.cz.json")
	load(Responses, http.StatusNotFound, "https://rdap.nic.cz/domain/non-existent.cz", "misc/empty.html")
	load(Responses, http.StatusOK, "https://rdap.nic.cz/domain/wrong-response-type.cz", "rdap/rdap.nic.cz/nameserver-ns2.pipni.cz.json")
	load(Responses, http.StatusOK, "https://rdap.nic.cz/domain/malformed.cz", "misc/malformed.json")

	return r
}
