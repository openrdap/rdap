// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"fmt"

	"github.com/skip2/openrdap/rdap"
)

func Example() {
	var jsonBlob = []byte(`
	`)

	// Decode the response.
	result, err := rdap.Decode(jsonBlob)

	// Print the response.
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		domain := result.(*rdap.Domain)
		PrintDomain(domain)
	}
}

// PrintDomain prints some basic rdap.Domain fields.
func PrintDomain(d *rdap.Domain) {
	// Registry unique identifier for the domain. Here, "google.cz".
	fmt.Printf("Handle=%s\n", d.Handle)

	// Domain name (LDH = letters, digits, hyphen). Here, "google.cz".
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
			fmt.Printf("  Name    : %s\n", v.Name())

			// Address.
			fmt.Printf("  POBox   : %s\n", v.POBox())
			fmt.Printf("  Ext     : %s\n", v.Ext())
			fmt.Printf("  Street  : %s\n", v.Street())
			fmt.Printf("  Locality: %s\n", v.Locality())
			fmt.Printf("  Region  : %s\n", v.Region())
			fmt.Printf("  Code    : %s\n", v.Code())
			fmt.Printf("  Country : %s\n", v.Country())

			// Phone numbers.
			fmt.Printf("  Tel     : %s\n", v.Tel())
			fmt.Printf("  Fax     : %s\n", v.Fax())

			// Email address.
			fmt.Printf("  Email   : %s\n", v.Email())

			// The raw VCard fields are also accessible.
		}
	}
}
