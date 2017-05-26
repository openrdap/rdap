// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package jcard

import (
	"reflect"
	"testing"

	"github.com/skip2/rdap/test"
)

func TestJCardErrors(t *testing.T) {
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
		j, err := NewJCard(test.LoadFile(filename))

		if j != nil || err == nil {
			t.Errorf("jCard with error unexpectedly parsed %s %v %s\n", filename, j, err)
		}
	}
}

func TestJCardExample(t *testing.T) {
	j, err := NewJCard(test.LoadFile("jcard/example.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	numProperties := 17
	if len(j.Properties) != numProperties {
		t.Errorf("Got %d properties expected %d", len(j.Properties), numProperties)
	}

	expectedVersion := &Property{
		Name:       "version",
		Parameters: make(map[string][]string),
		Type:       "text",
		Value:      "4.0",
	}

	if !reflect.DeepEqual(j.Get("version")[0], expectedVersion) {
		t.Errorf("version field incorrect")
	}

	expectedN := &Property{
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

	expectedTel0 := &Property{
		Name:       "tel",
		Parameters: map[string][]string{"type": []string{"work", "voice"}, "pref": []string{"1"}},
		Type:       "uri",
		Value:      "tel:+1-418-656-9254;ext=102",
	}

	if !reflect.DeepEqual(j.Get("tel")[0], expectedTel0) {
		t.Errorf("tel[0] field incorrect")
	}
}

func TestJCardMixedDatatypes(t *testing.T) {
	j, err := NewJCard(test.LoadFile("jcard/mixed.json"))
	if j == nil || err != nil {
		t.Errorf("jCard parse failed %v %s\n", j, err)
	}

	expectedMixed := &Property{
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
