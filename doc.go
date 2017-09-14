// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package rdap implements a client for the Registration Data Access Protocol (RDAP).
//
// RDAP is a modern replacement for the text-based WHOIS (port 43) protocol. It provides registration data for domain names/IP addresses/AS numbers, and more, in a structured format.
//
// This client executes RDAP queries and returns the responses as Go values.
//
// Example quick usage:
//   client := rdap.NewClient()
//   domain, err := client.QueryDomain("google.cz")
//
//   if err != nil {
//     fmt.Printf("name=%s, address=%s\n", domain.Registrant.Name, domain.Registrant.Address)
//   }
//
// Manual query construction, with options to fetch specific data (if available):
//  client := rdap.NewClient()
//  client.Options = FetchRegistrant | FetchTechnical | FetchNOC
//
//  query := rdap.NewAutnumQuery(5400)
//  response, err := client.Query(query)
//
// The above examples If you are running lots of RDAP queries, enable the bootstrap data disk cache ($HOME/.openrdap or %UserData%\openrdap):
//
//  - text based query for google.cz
//  - client options
//  - use of bootstrap cache, custom http, timeout
//
//  - success/partial success
//  - timeouts
//
// As of June 2017, all five number registries (AFRINIC, ARIN, APNIC, LANIC,
// RIPE) run RDAP servers. A small number of TLDs (top level domains) support
// RDAP so far, listed on https://data.iana.org/rdap/dns.json.
//
// The RDAP protocol uses HTTP, with responses in a JSON format. A bootstrapping mechanism (http://data.iana.org/rdap/) is used to determine the server to query.
package rdap
