// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"io"
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestPrintDomain(t *testing.T) {
	obj := loadObject("rdap/rdap.nic.cz/domain-example.cz.json")

	printer := &Printer{
		BriefLinks: true,
		Writer:     io.Discard,
	}

	printer.Print(obj)
}

func loadObject(filename string) RDAPObject {
	d := NewDecoder(test.LoadFile(filename))

	result, err := d.Decode()
	if err != nil {
		panic("Decode unexpectedly failed: " + err.Error())
	}

	return result
}

// cleanString runs on every printed heading and value. The clean case (no bad
// runes) is the common case and should avoid the rune-by-rune strings.Map scan.
var (
	benchCleanStringInputs = []string{
		"Domain Name",
		"EXAMPLE.COM",
		"2021-01-01T00:00:00Z",
		"client transfer prohibited",
		"https://rdap.verisign.com/com/v1/domain/EXAMPLE.COM",
		"ns1.example.com",
		"registrant",
		"Registrar Abuse Contact Email",
	}
	benchCleanStringDirty = "line one\nline two\r\x00trailing"
)

func BenchmarkCleanStringClean(b *testing.B) {
	p := &Printer{}
	b.ReportAllocs()
	for range b.N {
		for _, s := range benchCleanStringInputs {
			_ = p.cleanString(s)
		}
	}
}

func BenchmarkCleanStringDirty(b *testing.B) {
	p := &Printer{}
	b.ReportAllocs()
	for range b.N {
		_ = p.cleanString(benchCleanStringDirty)
	}
}
