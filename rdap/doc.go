// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package rdap implements decoding of RDAP responses.
//
// An RDAP response describes an object such as a Domain (data resembling "whois
// example.com"), or an IP Network (data resembling "whois 2001:db8::"). For a
// live example, see https://rdap.nic.cz/domain/ctk.cz.
//
// RDAP responses are JSON documents, as defined in RFC7483. This package
// decodes RDAP responses into Go values.
//
// To decode an RDAP response:
//
//  jsonBlob := []byte(`
//    {
//      "objectClassName": "domain",
//      "rdapConformance": ["rdap_level_0"],
//      "handle":          "EXAMPLECOM",
//      "ldhName":         "example.com",
//      "entities":        []
//    }
//  `)
//
//  d := rdap.NewDecoder(jsonBlob)
//  result, err := d.Decode()
//
//  if err != nil {
//    if domain, ok := result.(*rdap.Domain); ok {
//      fmt.Printf("Domain name = %s\n", domain.LDHName)
//    }
//  }
//
// RDAP responses are decoded into the following types:
//  &rdap.Error{}                   - Responses with an errorCode value.
//  &rdap.Autnum{}                  - Responses with objectClassName="autnum".
//  &rdap.Domain{}                  - Responses with objectClassName="domain".
//  &rdap.Entity{}                  - Responses with objectClassName="entity".
//  &rdap.IPNetwork{}               - Responses with objectClassName="ip network".
//  &rdap.Nameserver{}              - Responses with objectClassName="nameserver".
//  &rdap.DomainSearchResults{}     - Responses with a domainSearchResults array.
//  &rdap.EntitySearchResults{}     - Responses with a entitySearchResults array.
//  &rdap.NameserverSearchResults{} - Responses with a nameserverSearchResults array.
//  &rdap.Help{}                    - All other valid JSON responses.
//
// Note that an RDAP server may return a different response type than expected.
//
// The decoder supports unknown RDAP fields (such as "fred_nsset" in the
// rdap.nic.cz example above). These are stored in each result struct, see the
// DecodeData documentation for accessing them.
//
// Decoding is performed on a best-effort basis, with "minor error"s ignored.
// This avoids minor errors rendering a response undecodable.
//
// This package does not perform any network connections.
//
// For RFC7483, see https://tools.ietf.org/html/rfc7483.
package rdap
