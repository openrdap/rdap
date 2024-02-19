// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"reflect"
	"testing"

	"github.com/openrdap/rdap/test"
)

func TestVCardErrors(t *testing.T) {
	filenames := []string{
		"jcard/error_invalid_json.json",
		"jcard/error_bad_top_type.json",
		"jcard/error_bad_vcard_label.json",
		"jcard/error_bad_properties_array.json",
		"jcard/error_bad_property_size.json",
		"jcard/error_bad_property_name.json",
		"jcard/error_bad_property_type.json",
		"jcard/error_bad_property_parameters.json",
		"jcard/error_bad_property_parameters_2.json",
		"jcard/error_bad_property_nest_depth.json",
	}

	for _, filename := range filenames {
		j, err := NewVCard(test.LoadFile(filename))

		if j != nil || err == nil {
			t.Errorf("jCard with error unexpectedly parsed %s %v %s\n", filename, j, err)
		}
	}
}

func TestVCardIgnoreInvalidProperties(t *testing.T) {
	json := test.LoadFile("jcard/error_invalid_properties.json")

	j1, err1 := NewVCardWithOptions(json, VCardOptions{IgnoreInvalidProperties: true})
	if j1 == nil || len(j1.Properties) != 4 || err1 != nil {
		t.Errorf("jCard with ignored errors not parsed correctly\n")
	}

	j2, err2 := NewVCardWithOptions(json, VCardOptions{IgnoreInvalidProperties: false})
	if j2 != nil || err2 == nil {
		t.Errorf("jCard with errors unexpectedly parsed\n")
	}
}

func TestVCardExample(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/example.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	numProperties := 17
	if len(j.Properties) != numProperties {
		t.Errorf("Got %d properties expected %d", len(j.Properties), numProperties)
	}

	expectedVersion := &VCardProperty{
		Name:       "version",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      "4.0",
	}

	if !reflect.DeepEqual(j.Get("version")[0], expectedVersion) {
		t.Errorf("version field incorrect")
	}

	expectedN := &VCardProperty{
		Name:       "n",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      []interface{}{"Perreault", "Simon", "", "", []interface{}{"ing. jr", "M.Sc."}},
	}

	expectedFlatN := []string{
		"Perreault",
		"Simon",
		"",
		"",
		"ing. jr",
		"M.Sc.",
	}

	if !reflect.DeepEqual(j.Get("n")[0], expectedN) {
		t.Errorf("n field incorrect")
	}

	if !reflect.DeepEqual(j.Get("n")[0].Values(), expectedFlatN) {
		t.Errorf("n flat value incorrect")
	}

	expectedTel0 := &VCardProperty{
		Name:       "tel",
		Parameters: map[string][]string{"type": []string{"work", "voice"}, "pref": []string{"1"}},
		Type:       "uri",
		Value:      "tel:+1-418-656-9254;ext=102",
	}

	if !reflect.DeepEqual(j.Get("tel")[0], expectedTel0) {
		t.Errorf("tel[0] field incorrect")
	}
}

func TestVCardMixedDatatypes(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/mixed.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	expectedMixed := &VCardProperty{
		Name:       "mixed",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      []interface{}{"abc", true, float64(42), nil, []interface{}{"def", false, float64(43)}},
	}

	expectedFlatMixed := []string{
		"abc",
		"true",
		"42",
		"",
		"def",
		"false",
		"43",
	}

	if !reflect.DeepEqual(j.Get("mixed")[0], expectedMixed) {
		t.Errorf("mixed field incorrect")
	}

	flattened := j.Get("mixed")[0].Values()
	if !reflect.DeepEqual(flattened, expectedFlatMixed) {
		t.Errorf("mixed flat value incorrect %v", flattened)
	}
}

func TestVCardQuickAccessors(t *testing.T) {
	j, err := NewVCard(test.LoadFile("jcard/example.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	got := []string{
		j.Name(),
		j.POBox(),
		j.ExtendedAddress(),
		j.StreetAddress(),
		j.Locality(),
		j.Region(),
		j.PostalCode(),
		j.Country(),
		j.Tel(),
		j.Fax(),
		j.Email(),
		j.Org(),
	}

	expected := []string{
		"Simon Perreault",
		"",
		"Suite D2-630",
		"2875 Laurier",
		"Quebec",
		"QC",
		"G1V 2M2",
		"Canada",
		"tel:+1-418-656-9254;ext=102",
		"",
		"simon.perreault@viagenie.ca",
		"Viagenie",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got %v expected %v\n", got, expected)
	}
}
