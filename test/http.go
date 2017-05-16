// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package test

import (
	"net/http"
	"io/ioutil"
	"log"
	"path"
	"runtime"

	"gopkg.in/jarcoal/httpmock.v1"
)

type TestDataset int
const (
	Bootstrap TestDataset = iota
	BootstrapExperimental
	BootstrapMalformed
	BootstrapComplex
	BootstrapHTTPError
)

type response struct {
	Status int
	URL string
	Body string
}

var responses map[TestDataset][]response
var activatedURLs map[string]bool
var testDataPath string

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

func Finish() {
	activatedURLs = make(map[string]bool)
	httpmock.DeactivateAndReset()
}

func Get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	return data
}

func init() {
	responses = make(map[TestDataset][]response)
	activatedURLs = make(map[string]bool)

	testDataPath = findTestDataPath()

	loadTestDatasets()
}

func findTestDataPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller() failed")
	}

	dir, _ := path.Split(filename)

	return path.Join(dir, "testdata")
}

func loadTestDatasets() {
	// Valid snapshot of the IANA bootstrap files.
	load(Bootstrap, 200, "https://data.iana.org/rdap/asn.json", "bootstrap/asn.json")
	load(Bootstrap, 200, "https://data.iana.org/rdap/dns.json", "bootstrap/dns.json")
	load(Bootstrap, 200, "https://data.iana.org/rdap/ipv4.json", "bootstrap/ipv4.json")
	load(Bootstrap, 200, "https://data.iana.org/rdap/ipv6.json", "bootstrap/ipv6.json")

	// Experimental bootstrap file for service providers.
	// https://datatracker.ietf.org/doc/draft-hollenbeck-regext-rdap-object-tag/ .
	load(BootstrapExperimental, 200, "https://www.openrdap.org/rdap/service_provider.json", "bootstrap_experimental/service_provider.json")

	// Malformed bootstrap files.
	load(BootstrapMalformed, 200, "https://www.example.org/dns_bad_services.json", "bootstrap_malformed/dns_bad_services.json")
	load(BootstrapMalformed, 200, "https://www.example.org/dns_bad_url.json", "bootstrap_malformed/dns_bad_url.json")
	load(BootstrapMalformed, 200, "https://www.example.org/dns_empty.json", "bootstrap_malformed/dns_empty.json")
	load(BootstrapMalformed, 200, "https://www.example.org/dns_syntax_error.json", "bootstrap_malformed/dns_syntax_error.json")

	// Valid bootstrap files testing more features than yet used by IANA.
	load(BootstrapComplex, 200, "https://rdap.example.org/dns.json", "bootstrap_complex/dns.json")

	// Bootstrap HTTP errors
	load(BootstrapHTTPError, 404, "https://data.iana.org/rdap/asn.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, 404, "https://data.iana.org/rdap/dns.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, 404, "https://data.iana.org/rdap/ipv4.json", "bootstrap_http_error/404.html")
	load(BootstrapHTTPError, 404, "https://data.iana.org/rdap/ipv6.json", "bootstrap_http_error/404.html")
}

func load(set TestDataset, status int, url string, filename string) {
	var body []byte

	body, err := ioutil.ReadFile(path.Join(testDataPath, filename))

	if err != nil {
		log.Panic(err)
	}

	responses[set] = append(responses[set], response{status, url, string(body)})
}
