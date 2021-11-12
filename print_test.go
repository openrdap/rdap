// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestPrintDomain(t *testing.T) {
	obj := loadObject("rdap/rdap.nic.cz/domain-example.cz.json")

	printer := &Printer{
		BriefLinks: true,
	}

	_ = obj
	_ = printer
	//printer.Print(obj)
}

func loadObject(filename string) RDAPObject {
	jsonBlob := test.LoadFile(filename)

	d := NewDecoder([]byte(jsonBlob))
	result, err := d.Decode()

	if err != nil {
		panic("Decode unexpectedly failed")
	}

	return result
}
