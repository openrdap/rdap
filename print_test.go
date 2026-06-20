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
