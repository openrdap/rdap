// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestDecodeEmpty(t *testing.T) {
	type Empty struct {
	}
	runDecodeAndCompareTest(t, &Empty{}, `
	{}
`, &Empty{})
}

func TestDecodeDecodeData(t *testing.T) {
	type XYZ struct {
		DecodeData *DecodeData

		S1 string
		S2 string `rdap:"s2Name"`
		SF string
	}

	result, ok := runDecode(t, &XYZ{}, `
	{
		"s1": "S1",
		"s2Name": "S2",
		"sF": 1.5,
		"unknown" : "value"
	}`)

	if !ok {
		return
	}

	x := result.(*XYZ)

	if x.S1 != "S1" || x.S2 != "S2" || x.SF != "1.5" {
		t.Errorf("Decode values bad %v", x)
	}

	if x.DecodeData == nil {
		t.Errorf("DecodeData not instantiated")
	} else if len(x.DecodeData.Notes("sF")) != 1 {
		t.Errorf("DecodeData notes not added")
	} else if len(x.DecodeData.Fields()) != 4 {
		t.Errorf("DecodeData Fields() bad")
	} else if len(x.DecodeData.UnknownFields()) != 1 {
		t.Errorf("DecodeData UnknownFields() bad")
	} else if !reflect.DeepEqual(x.DecodeData.Value("unknown"), "value") {
		t.Errorf("DecodeData bad Value()")
	}
}

func TestDecodeVCard(t *testing.T) {
	type XYZ struct {
		VCard *VCard
	}

	result, ok := runDecode(t, &XYZ{}, `
	{
		"vCard": [
			"vcard",
			[
				["version", {}, "text", "4.0"],
				["fn", {}, "text", "First Last"]
			]
		]
	}
	`)

	if !ok {
		return
	}

	x := result.(*XYZ)

	if x.VCard == nil {
		t.Errorf("VCard not decoded")
	} else if len(x.VCard.Properties) != 2 {
		t.Errorf("VCard properties not decoded")
	}
}

func TestDecodeSlice(t *testing.T) {
	type XYZ struct {
		S []string
	}

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"s": ["a", "b"]
	}
	`, &XYZ{
		S: []string{"a", "b"},
	})
}

func TestDecodeMap(t *testing.T) {
	type XYZ struct {
		M map[string]string
	}

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"m": {"a": "av", "b": "bv"}
	}
	`, &XYZ{
		M: map[string]string{"a": "av", "b": "bv"},
	})
}

func TestDecodeUints(t *testing.T) {
	type XYZ struct {
		A         uint8
		AOverflow uint8
		B         uint16
		C         uint32
		D         uint64

		S  uint8
		BF uint8
		BT uint8
		N  uint8
	}

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"a": 100,
		"aOverflow": 256,
		"b": 200,
		"c": 42,
		"d": 43,
		"s": "10",
		"bF": false,
		"bT": true,
		"n": null
	}
	`, &XYZ{
		A:         100,
		AOverflow: 0,
		B:         200,
		C:         42,
		D:         43,
		S:         10,
		BF:        0,
		BT:        1,
		N:         0,
	})
}

func TestDecodeInts(t *testing.T) {
	type XYZ struct {
		A          int8
		AUnderflow int8
		AOverflow  int8
		B          int16
		C          int32
		D          int64

		S  int8
		BF int8
		BT int8
		N  int8
	}

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"a": 100,
		"aUnderflow": -129,
		"aOverflow": 128,
		"b": 200,
		"c": 42,
		"d": 43,
		"s": "10",
		"bF": false,
		"bT": true,
		"n": null
	}
	`, &XYZ{
		A:          100,
		AUnderflow: 0,
		AOverflow:  0,
		B:          200,
		C:          42,
		D:          43,
		S:          10,
		BF:         0,
		BT:         1,
		N:          0,
	})
}

func TestDecodeFloat64(t *testing.T) {
	type XYZ struct {
		F    float64
		FPtr *float64

		S1 float64
		S2 float64

		BF float64
		BT float64

		N float64
	}

	fptr := 1.5

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"f": 1.5,
		"fPtr": 1.5,
		"s1": "1.5",
		"s2": "-1.5",
		"bF": false,
		"bT": true,
		"n": null
	}
	`, &XYZ{
		F:    1.5,
		FPtr: &fptr,
		S1:   1.5,
		S2:   -1.5,
		BF:   0.0,
		BT:   1.0,
		N:    0.0,
	})
}

func TestDecodeBool(t *testing.T) {
	type XYZ struct {
		B    bool
		BPtr *bool

		SF bool
		ST bool

		FF bool
		FT bool

		N bool
	}

	bptr := true

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"b": true,
		"bPtr": true,
		"sF": "false",
		"sT": "true",
		"fF": 0,
		"fT": 1,
		"n": null
	}
	`, &XYZ{
		B:    true,
		BPtr: &bptr,
		ST:   true,
		SF:   false,
		FF:   false,
		FT:   true,
		N:    false,
	})
}

func TestDecodeString(t *testing.T) {
	type XYZ struct {
		S    string
		SPtr *string

		BT string
		BF string

		F1 string
		F2 string

		N string
	}

	sptr := "sptr"

	runDecodeAndCompareTest(t, &XYZ{}, `
	{
		"s": "test", 
		"sPtr": "sptr", 
		"bT": true,
		"bF": false, 
		"f1": 1.0,
		"f2": -3.14,
		"n2": null
	}
	`, &XYZ{
		S:    "test",
		SPtr: &sptr,
		BT:   "true",
		BF:   "false",
		F1:   "1",
		F2:   "-3.14",
		N:    "",
	})
}

func runDecode(t *testing.T, target interface{}, jsonBlob string) (interface{}, bool) {
	d := NewDecoder([]byte(jsonBlob))
	d.target = target

	result, err := d.Decode()

	if err != nil {
		t.Errorf("While decoding '%s', got error: %s", jsonBlob, err)
		return result, false
	}

	return result, true
}

func runDecodeAndCompareTest(t *testing.T, target interface{}, jsonBlob string, expected interface{}) {
	result, ok := runDecode(t, target, jsonBlob)

	if !ok {
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("While decoding '%s':\nexpected %s\ngot %s",
			jsonBlob,
			spew.Sdump(expected),
			spew.Sdump(result))
	}
}
