<img src="https://www.openrdap.org/public/img/logo.png">

OpenRDAP is an command line [RDAP](https://datatracker.ietf.org/wg/weirds/documents/) client implementation in Go.
[![Build Status](https://travis-ci.org/openrdap/rdap.svg?branch=master)](https://travis-ci.org/openrdap/rdap)

https://www.openrdap.org - homepage

https://www.openrdap.org/demo - live demo

## Features
* Command line RDAP client
* Output formats: text, JSON, WHOIS style
* Query types supported:
    * ip
    * domain
    * autnum
    * nameserver
    * entity
    * help
    * url
    * domain-search
    * domain-search-by-nameserver
    * domain-search-by-nameserver-ip
    * nameserver-search
    * nameserver-search-by-ip
    * entity-search
    * entity-search-by-handle
* Automatic server detection for ip/domain/autnum/entities
* Object tags support
* Bootstrap cache (optional, uses ~/.openrdap by default)
* X.509 client authentication

## Installation

This program uses Go. The Go compiler is available from https://golang.org/.

To install:

    go install github.com/openrdap/rdap/cmd/rdap@master

This will install the "rdap" binary in your $GOPATH/go/bin directory. Try running:

    ~/go/bin/rdap google.com

## Usage

| Query type                | Usage                                                                    |
| ------------------------- | ------------------------------------------------------------------------ |
| Domain (.com)             | rdap -v example.com                                                      |
| IPv4 Address              | rdap -v 192.0.2.0                                                        |
| IPv6 Address              | rdap -v 2001:db8::                                                       |
| Autonomous System (ASN)   | rdap -v AS15169                                                          |
| Entity (with object tag)  | rdap -v OPS4-RIPE                                                        |

## Advanced usage (server must be specified using -s; not all servers support all query types)
| Query type                | Usage                                                                    |
| ------------------------- | ------------------------------------------------------------------------ |
| Nameserver                | rdap -v -t nameserver -s https://rdap.verisign.com/com/v1 ns1.google.com |
| Help                      | rdap -v -t help -s https://rdap.verisign.com/com/v1                      |
| Domain Search             | rdap -v -t domain-search -s $SERVER_URL example*.gtld                    |
| Domain Search (by NS)     | rdap -v -t domain-search-by-nameserver -s $SERVER_URL ns1.example.gtld   |
| Domain Search (by NS IP)  | rdap -v -t domain-search-by-nameserver-ip -s $SERVER_URL 192.0.2.0       |
| Nameserver Search         | rdap -v -t nameserver-search -s $SERVER_URL ns1.example.gtld             |
| Nameserver Search (by IP) | rdap -v -t nameserver-search-by-ip -s $SERVER_URL 192.0.2.0              |
| Entity Search             | rdap -v -t entity-search -s $SERVER_URL ENTITY-TAG                       |
| Entity Search (by handle) | rdap -v -t entity-search-by-handle -s $SERVER_URL ENTITY-TAG             |

See https://www.openrdap.org/docs.

## Example output

Click the examples to see the output:

<details>
<summary>rdap example.com</summary>

```Domain:
  Domain Name: EXAMPLE.COM
  Handle: 2336799_DOMAIN_COM-VRSN
  Status: client delete prohibited
  Status: client transfer prohibited
  Status: client update prohibited
  Conformance: rdap_level_0
  Conformance: icann_rdap_technical_implementation_guide_0
  Conformance: icann_rdap_response_profile_0
  Notice:
    Title: Terms of Use
    Description: Service subject to Terms of Use.
    Link: https://www.verisign.com/domain-names/registration-data-access-protocol/terms-service/index.xhtml
  Notice:
    Title: Status Codes
    Description: For more information on domain status codes, please visit https://icann.org/epp
    Link: https://icann.org/epp
  Notice:
    Title: RDDS Inaccuracy Complaint Form
    Description: URL of the ICANN RDDS Inaccuracy Complaint Form: https://icann.org/wicf
    Link: https://icann.org/wicf
  Link: https://rdap.verisign.com/com/v1/domain/EXAMPLE.COM
  Event:
    Action: registration
    Date: 1995-08-14T04:00:00Z
  Event:
    Action: expiration
    Date: 2023-08-13T04:00:00Z
  Event:
    Action: last changed
    Date: 2023-05-12T15:13:35Z
  Event:
    Action: last update of RDAP database
    Date: 2023-05-16T20:36:06Z
  Secure DNS:
    Delegation Signed: true
    DSData:
      Key Tag: 370
      Algorithm: 13
      Digest: BE74359954660069D5C63D200C39F5603827D7DD02B56F120EE9F3A86764247C
      DigestType: 2
  Entity:
    Handle: 376
    Public ID:
      Type: IANA Registrar ID
      Identifier: 376
    Role: registrar
    vCard version: 4.0
    vCard fn: RESERVED-Internet Assigned Numbers Authority
    Entity:
      Role: abuse
      vCard version: 4.0
  Nameserver:
    Nameserver: A.IANA-SERVERS.NET
  Nameserver:
    Nameserver: B.IANA-SERVERS.NET
```

</details>

<details>
<summary>rdap 8.8.8.8</summary>

```IP Network:
  Handle: NET-8-8-8-0-1
  Start Address: 8.8.8.0
  End Address: 8.8.8.255
  IP Version: v4
  Name: LVLT-GOGL-8-8-8
  Type: ALLOCATION
  ParentHandle: NET-8-0-0-0-1
  Status: active
  Port43: whois.arin.net
  Notice:
    Title: Terms of Service
    Description: By using the ARIN RDAP/Whois service, you are agreeing to the RDAP/Whois Terms of Use
    Link: https://www.arin.net/resources/registry/whois/tou/
  Notice:
    Title: Whois Inaccuracy Reporting
    Description: If you see inaccuracies in the results, please visit: 
    Link: https://www.arin.net/resources/registry/whois/inaccuracy_reporting/
  Notice:
    Title: Copyright Notice
    Description: Copyright 1997-2023, American Registry for Internet Numbers, Ltd.
  Entity:
    Handle: GOGL
    Port43: whois.arin.net
    Remark:
      Title: Registration Comments
      Description: Please note that the recommended way to file abuse complaints are located in the following links. 
      Description: To report abuse and illegal activity: https://www.google.com/contact/
      Description: For legal requests: http://support.google.com/legal 
      Description: Regards, 
      Description: The Google Team
    Link: https://rdap.arin.net/registry/entity/GOGL
    Link: https://whois.arin.net/rest/org/GOGL
    Event:
      Action: last changed
      Date: 2019-10-31T15:45:45-04:00
    Event:
      Action: registration
      Date: 2000-03-30T00:00:00-05:00
    Role: registrant
    vCard version: 4.0
    vCard fn: Google LLC
    vCard kind: org
    Entity:
      Handle: ABUSE5250-ARIN
      Status: validated
      Port43: whois.arin.net
      Remark:
        Title: Registration Comments
        Description: Please note that the recommended way to file abuse complaints are located in the following links.
        Description: To report abuse and illegal activity: https://www.google.com/contact/
        Description: For legal requests: http://support.google.com/legal 
        Description: Regards,
        Description: The Google Team
      Link: https://rdap.arin.net/registry/entity/ABUSE5250-ARIN
      Link: https://whois.arin.net/rest/poc/ABUSE5250-ARIN
      Event:
        Action: last changed
        Date: 2022-10-24T08:43:11-04:00
      Event:
        Action: registration
        Date: 2015-11-06T15:36:35-05:00
      Role: abuse
      vCard version: 4.0
      vCard fn: Abuse
      vCard org: Abuse
      vCard kind: group
      vCard email: network-abuse@google.com
      vCard tel: +1-650-253-0000
    Entity:
      Handle: ZG39-ARIN
      Status: validated
      Port43: whois.arin.net
      Link: https://rdap.arin.net/registry/entity/ZG39-ARIN
      Link: https://whois.arin.net/rest/poc/ZG39-ARIN
      Event:
        Action: last changed
        Date: 2022-11-10T07:12:44-05:00
      Event:
        Action: registration
        Date: 2000-11-30T13:54:08-05:00
      Role: technical
      Role: administrative
      vCard version: 4.0
      vCard fn: Google LLC
      vCard org: Google LLC
      vCard kind: group
      vCard email: arin-contact@google.com
      vCard tel: +1-650-253-0000
  Link: https://rdap.arin.net/registry/ip/8.8.8.0
  Link: https://whois.arin.net/rest/net/NET-8-8-8-0-1
  Link: https://rdap.arin.net/registry/ip/8.0.0.0/9
  Event:
    Action: last changed
    Date: 2014-03-14T16:52:05-04:00
  Event:
    Action: registration
    Date: 2014-03-14T16:52:05-04:00
  cidr0_cidrs:
    v4prefix: 8.8.8.0
    length: 24
```

</details>

<details>
<summary>rdap --json AS15169</summary>

```
{
  "rdapConformance": [
    "nro_rdap_profile_0",
    "rdap_level_0",
    "nro_rdap_profile_asn_flat_0"
  ],
  "notices": [
    {
      "title": "Terms of Service",
      "description": [
        "By using the ARIN RDAP/Whois service, you are agreeing to the RDAP/Whois Terms of Use"
      ],
      "links": [
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "terms-of-service",
          "type": "text/html",
          "href": "https://www.arin.net/resources/registry/whois/tou/"
        }
      ]
    },
    {
      "title": "Whois Inaccuracy Reporting",
      "description": [
        "If you see inaccuracies in the results, please visit: "
      ],
      "links": [
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "inaccuracy-report",
          "type": "text/html",
          "href": "https://www.arin.net/resources/registry/whois/inaccuracy_reporting/"
        }
      ]
    },
    {
      "title": "Copyright Notice",
      "description": [
        "Copyright 1997-2023, American Registry for Internet Numbers, Ltd."
      ]
    }
  ],
  "handle": "AS15169",
  "startAutnum": 15169,
  "endAutnum": 15169,
  "name": "GOOGLE",
  "events": [
    {
      "eventAction": "last changed",
      "eventDate": "2012-02-24T09:44:34-05:00"
    },
    {
      "eventAction": "registration",
      "eventDate": "2000-03-30T00:00:00-05:00"
    }
  ],
  "links": [
    {
      "value": "https://rdap.arin.net/registry/autnum/15169",
      "rel": "self",
      "type": "application/rdap+json",
      "href": "https://rdap.arin.net/registry/autnum/15169"
    },
    {
      "value": "https://rdap.arin.net/registry/autnum/15169",
      "rel": "alternate",
      "type": "application/xml",
      "href": "https://whois.arin.net/rest/asn/AS15169"
    }
  ],
  "entities": [
    {
      "handle": "GOGL",
      "vcardArray": [
        "vcard",
        [
          [
            "version",
            {},
            "text",
            "4.0"
          ],
          [
            "fn",
            {},
            "text",
            "Google LLC"
          ],
          [
            "adr",
            {
              "label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"
            },
            "text",
            [
              "",
              "",
              "",
              "",
              "",
              "",
              ""
            ]
          ],
          [
            "kind",
            {},
            "text",
            "org"
          ]
        ]
      ],
      "roles": [
        "registrant"
      ],
      "remarks": [
        {
          "title": "Registration Comments",
          "description": [
            "Please note that the recommended way to file abuse complaints are located in the following links. ",
            "",
            "To report abuse and illegal activity: https://www.google.com/contact/",
            "",
            "For legal requests: http://support.google.com/legal ",
            "",
            "Regards, ",
            "The Google Team"
          ]
        }
      ],
      "links": [
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "self",
          "type": "application/rdap+json",
          "href": "https://rdap.arin.net/registry/entity/GOGL"
        },
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "alternate",
          "type": "application/xml",
          "href": "https://whois.arin.net/rest/org/GOGL"
        }
      ],
      "events": [
        {
          "eventAction": "last changed",
          "eventDate": "2019-10-31T15:45:45-04:00"
        },
        {
          "eventAction": "registration",
          "eventDate": "2000-03-30T00:00:00-05:00"
        }
      ],
      "entities": [
        {
          "handle": "ABUSE5250-ARIN",
          "vcardArray": [
            "vcard",
            [
              [
                "version",
                {},
                "text",
                "4.0"
              ],
              [
                "adr",
                {
                  "label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"
                },
                "text",
                [
                  "",
                  "",
                  "",
                  "",
                  "",
                  "",
                  ""
                ]
              ],
              [
                "fn",
                {},
                "text",
                "Abuse"
              ],
              [
                "org",
                {},
                "text",
                "Abuse"
              ],
              [
                "kind",
                {},
                "text",
                "group"
              ],
              [
                "email",
                {},
                "text",
                "network-abuse@google.com"
              ],
              [
                "tel",
                {
                  "type": [
                    "work",
                    "voice"
                  ]
                },
                "text",
                "+1-650-253-0000"
              ]
            ]
          ],
          "roles": [
            "abuse"
          ],
          "remarks": [
            {
              "title": "Registration Comments",
              "description": [
                "Please note that the recommended way to file abuse complaints are located in the following links.",
                "",
                "To report abuse and illegal activity: https://www.google.com/contact/",
                "",
                "For legal requests: http://support.google.com/legal ",
                "",
                "Regards,",
                "The Google Team"
              ]
            }
          ],
          "links": [
            {
              "value": "https://rdap.arin.net/registry/autnum/15169",
              "rel": "self",
              "type": "application/rdap+json",
              "href": "https://rdap.arin.net/registry/entity/ABUSE5250-ARIN"
            },
            {
              "value": "https://rdap.arin.net/registry/autnum/15169",
              "rel": "alternate",
              "type": "application/xml",
              "href": "https://whois.arin.net/rest/poc/ABUSE5250-ARIN"
            }
          ],
          "events": [
            {
              "eventAction": "last changed",
              "eventDate": "2022-10-24T08:43:11-04:00"
            },
            {
              "eventAction": "registration",
              "eventDate": "2015-11-06T15:36:35-05:00"
            }
          ],
          "status": [
            "validated"
          ],
          "port43": "whois.arin.net",
          "objectClassName": "entity"
        },
        {
          "handle": "ZG39-ARIN",
          "vcardArray": [
            "vcard",
            [
              [
                "version",
                {},
                "text",
                "4.0"
              ],
              [
                "adr",
                {
                  "label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"
                },
                "text",
                [
                  "",
                  "",
                  "",
                  "",
                  "",
                  "",
                  ""
                ]
              ],
              [
                "fn",
                {},
                "text",
                "Google LLC"
              ],
              [
                "org",
                {},
                "text",
                "Google LLC"
              ],
              [
                "kind",
                {},
                "text",
                "group"
              ],
              [
                "email",
                {},
                "text",
                "arin-contact@google.com"
              ],
              [
                "tel",
                {
                  "type": [
                    "work",
                    "voice"
                  ]
                },
                "text",
                "+1-650-253-0000"
              ]
            ]
          ],
          "roles": [
            "technical",
            "administrative"
          ],
          "links": [
            {
              "value": "https://rdap.arin.net/registry/autnum/15169",
              "rel": "self",
              "type": "application/rdap+json",
              "href": "https://rdap.arin.net/registry/entity/ZG39-ARIN"
            },
            {
              "value": "https://rdap.arin.net/registry/autnum/15169",
              "rel": "alternate",
              "type": "application/xml",
              "href": "https://whois.arin.net/rest/poc/ZG39-ARIN"
            }
          ],
          "events": [
            {
              "eventAction": "last changed",
              "eventDate": "2022-11-10T07:12:44-05:00"
            },
            {
              "eventAction": "registration",
              "eventDate": "2000-11-30T13:54:08-05:00"
            }
          ],
          "status": [
            "validated"
          ],
          "port43": "whois.arin.net",
          "objectClassName": "entity"
        }
      ],
      "port43": "whois.arin.net",
      "objectClassName": "entity"
    },
    {
      "handle": "ZG39-ARIN",
      "vcardArray": [
        "vcard",
        [
          [
            "version",
            {},
            "text",
            "4.0"
          ],
          [
            "adr",
            {
              "label": "1600 Amphitheatre Parkway\nMountain View\nCA\n94043\nUnited States"
            },
            "text",
            [
              "",
              "",
              "",
              "",
              "",
              "",
              ""
            ]
          ],
          [
            "fn",
            {},
            "text",
            "Google LLC"
          ],
          [
            "org",
            {},
            "text",
            "Google LLC"
          ],
          [
            "kind",
            {},
            "text",
            "group"
          ],
          [
            "email",
            {},
            "text",
            "arin-contact@google.com"
          ],
          [
            "tel",
            {
              "type": [
                "work",
                "voice"
              ]
            },
            "text",
            "+1-650-253-0000"
          ]
        ]
      ],
      "roles": [
        "technical"
      ],
      "links": [
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "self",
          "type": "application/rdap+json",
          "href": "https://rdap.arin.net/registry/entity/ZG39-ARIN"
        },
        {
          "value": "https://rdap.arin.net/registry/autnum/15169",
          "rel": "alternate",
          "type": "application/xml",
          "href": "https://whois.arin.net/rest/poc/ZG39-ARIN"
        }
      ],
      "events": [
        {
          "eventAction": "last changed",
          "eventDate": "2022-11-10T07:12:44-05:00"
        },
        {
          "eventAction": "registration",
          "eventDate": "2000-11-30T13:54:08-05:00"
        }
      ],
      "status": [
        "validated"
      ],
      "port43": "whois.arin.net",
      "objectClassName": "entity"
    }
  ],
  "port43": "whois.arin.net",
  "status": [
    "active"
  ],
  "objectClassName": "autnum"
}
```

</details>

## Go docs
[![godoc](https://godoc.org/github.com/openrdap/rdap?status.png)](https://godoc.org/github.com/openrdap/rdap)

## Uses
Go 1.20+

## Links
- Wikipedia - [Registration Data Access Protocol](https://en.wikipedia.org/wiki/Registration_Data_Access_Protocol)
- [ICANN RDAP pilot](https://www.icann.org/rdap)

- [OpenRDAP](https://www.openrdap.org)

- https://data.iana.org/rdap/ - Official IANA bootstrap information

- [RFC 7480 HTTP Usage in the Registration Data Access Protocol (RDAP)](https://tools.ietf.org/html/rfc7480)
- [RFC 7481 Security Services for the Registration Data Access Protocol (RDAP)](https://tools.ietf.org/html/rfc7481)
- [RFC 7482 Registration Data Access Protocol (RDAP) Query Format](https://tools.ietf.org/html/rfc7482)
- [RFC 7483 JSON Responses for the Registration Data Access Protocol (RDAP)](https://tools.ietf.org/html/rfc7483)
- [RFC 7484 Finding the Authoritative Registration Data (RDAP) Service](https://tools.ietf.org/html/rfc7484)
- [RFC 8521 Registration Data Access Protocol (RDAP) Object Tagging] (https://datatracker.ietf.org/doc/rfc8521/)
