// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package rdap implements decoding of RDAP responses.
//
// RDAP responses (which are JSON documents, defined in RFC7483) are decoded to Go values.
//
// To quickly decode an RDAP response:
//
//  jsonBlob := []byte("...")
//  result, err := rdap.Decode(jsonBlob)
//
//  if err != nil {
//    if domain, ok := result.(*rdap.Domain); ok {
//      fmt.Printf("Domain name = %s\n", domain.LDHName)
//    }
//  }
//
// There are several RDAP response types. On success, |result| is one of:
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
// This package does not perform any network connections.
//
// For RFC7483, see https://tools.ietf.org/html/rfc7483.
package rdap
