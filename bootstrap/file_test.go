// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestParseValid(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	var bytes []byte = test.Get("https://data.iana.org/rdap/dns.json")

	var r *File
	r, err := NewFile(bytes)

	if err != nil {
		t.Fatal(err)
	}

	if len(r.Entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d: %v\n", len(r.Entries), r)
	}
}

func TestParseEmpty(t *testing.T) {
	test.Start(test.BootstrapMalformed)
	defer test.Finish()

	var bytes []byte = test.Get("https://www.example.org/dns_empty.json")

	_, err := NewFile(bytes)

	if err == nil {
		t.Fatal("Unexpected success parsing empty file")
	}
}

func TestParseSyntaxError(t *testing.T) {
	test.Start(test.BootstrapMalformed)
	defer test.Finish()

	var bytes []byte = test.Get("https://www.example.org/dns_syntax_error.json")

	_, err := NewFile(bytes)

	if err == nil {
		t.Fatal("Unexpected success parsing file with syntax error")
	}
}

func TestParseBadServices(t *testing.T) {
	test.Start(test.BootstrapMalformed)
	defer test.Finish()

	var bytes []byte = test.Get("https://www.example.org/dns_bad_services.json")

	_, err := NewFile(bytes)

	if err == nil {
		t.Fatal("Unexpected success parsing file with bad services array")
	}
}

func TestParseBadURL(t *testing.T) {
	test.Start(test.BootstrapMalformed)
	defer test.Finish()

	var bytes []byte = test.Get("https://www.example.org/dns_bad_url.json")

	var r *File
	r, err := NewFile(bytes)

	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal("Unexpected error parsing file with bad URL")
	}

	if len(r.Entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d: %v\n", len(r.Entries), r)
	}
}
