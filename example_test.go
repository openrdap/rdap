// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap_test

import (
	"fmt"

	"github.com/openrdap/rdap"
)

func Example() {
	var jsonBlob = []byte(`
		{
			"objectClassName": "domain",
			"rdapConformance": ["rdap_level_0"],
			"handle":          "EXAMPLECOM",
			"ldhName":         "example.com",
			"status":          ["active"],
			"entities":        [
				{
					"objectClassName": "entity",
					"handle": "EXAMPLECOMREG",
					"roles": ["registrant"],
					"vcardArray": [
						"vcard",
						[
							["version", {}, "text", "4.0"],
							["fn", {}, "text", "John Smith"],
							["adr", {}, "text",
								[
									"Box 1",
									"Suite 29",
									"1234 Fake St",
									"Toronto",
									"ON",
									"M5E 1W5",
									"Canada"
								]
							],
							["tel", {}, "uri", "tel:+1-555-555-5555"],
							["email", {}, "text", "hi@example.com"]
						]
					]
				}
			]
		}
	`)

	// Decode the response.
	d := rdap.NewDecoder(jsonBlob)
	result, err := d.Decode()

	// Print the response.
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		domain := result.(*rdap.Domain)
		PrintDomain(domain)
	}

	// Output:
	// Handle=EXAMPLECOM
	// LDHName=example.com
	// Status=[active]
	// Contact 0, roles: [registrant]:
	//   Name      : 'John Smith'
	//   POBox     : 'Box 1'
	//   Ext       : 'Suite 29'
	//   Street    : '1234 Fake St'
	//   Locality  : 'Toronto'
	//   Region    : 'ON'
	//   PostalCode: 'M5E 1W5'
	//   Country   : 'Canada'
	//   Tel       : '+1-555-555-5555'
	//   Fax       : ''
	//   Email     : 'hi@example.com'
}

// PrintDomain prints some basic rdap.Domain fields.
func PrintDomain(d *rdap.Domain) {
	// Registry unique identifier for the domain. Here, "example.cz".
	fmt.Printf("Handle=%s\n", d.Handle)

	// Domain name (LDH = letters, digits, hyphen). Here, "example.cz".
	fmt.Printf("LDHName=%s\n", d.LDHName)

	// Domain registration status. Here, "active".
	// See https://tools.ietf.org/html/rfc7483#section-10.2.2.
	fmt.Printf("Status=%v\n", d.Status)

	// Contact information.
	for i, e := range d.Entities {
		// Contact roles, such as "registrant", "administrative", "billing".
		// See https://tools.ietf.org/html/rfc7483#section-10.2.4.
		fmt.Printf("Contact %d, roles: %v:\n", i, e.Roles)

		// RDAP uses VCard for contact information, including name, address,
		// telephone, and e-mail address.
		if e.VCard != nil {
			v := e.VCard

			// Name.
			fmt.Printf("  Name      : '%s'\n", v.Name())

			// Address.
			fmt.Printf("  POBox     : '%s'\n", v.POBox())
			fmt.Printf("  Ext       : '%s'\n", v.ExtendedAddress())
			fmt.Printf("  Street    : '%s'\n", v.StreetAddress())
			fmt.Printf("  Locality  : '%s'\n", v.Locality())
			fmt.Printf("  Region    : '%s'\n", v.Region())
			fmt.Printf("  PostalCode: '%s'\n", v.PostalCode())
			fmt.Printf("  Country   : '%s'\n", v.Country())

			// Phone numbers.
			fmt.Printf("  Tel       : '%s'\n", v.Tel())
			fmt.Printf("  Fax       : '%s'\n", v.Fax())

			// Email address.
			fmt.Printf("  Email     : '%s'\n", v.Email())

			// The raw VCard fields are also accessible.
		}
	}
}
