// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import "testing"

func TestWhoisStyleDNSSEC(t *testing.T) {
	yes, no := true, false

	tests := []struct {
		name      string
		secureDNS *SecureDNS
		want      []string
	}{
		{"no secureDNS", nil, nil},
		{"signed", &SecureDNS{DelegationSigned: &yes}, []string{"signedDelegation"}},
		{"unsigned", &SecureDNS{DelegationSigned: &no}, []string{"unsigned"}},
		{"unsigned wins over DS", &SecureDNS{DelegationSigned: &no, DS: []DSData{{}}}, []string{"unsigned"}},
		{"DS without delegationSigned", &SecureDNS{DS: []DSData{{}}}, []string{"signedDelegation"}},
		{"keyData without delegationSigned", &SecureDNS{Keys: []KeyData{{}}}, nil},
		{"empty secureDNS", &SecureDNS{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{Object: &Domain{
				LDHName:     "example.com",
				Nameservers: []Nameserver{{LDHName: "ns1.example.com"}},
				SecureDNS:   tt.secureDNS,
			}}
			w := r.ToWhoisStyleResponse()

			got := w.Data["DNSSEC"]
			if len(got) != len(tt.want) || (len(got) == 1 && got[0] != tt.want[0]) {
				t.Fatalf("DNSSEC = %v, want %v", got, tt.want)
			}

			// ICANN WHOIS ordering places DNSSEC directly after the name servers.
			if tt.want != nil {
				order := w.KeyDisplayOrder
				if order[len(order)-1] != "DNSSEC" || order[len(order)-2] != "Name Server" {
					t.Fatalf("KeyDisplayOrder = %v, want DNSSEC after Name Server", order)
				}
			}
		})
	}
}
