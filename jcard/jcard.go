// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Type JCard represents a jCard (a JSON formatted vCard), as defined in https://tools.ietf.org/html/rfc7095.
//
// A jCard represents information about an individual or entity. It can include a name, telephone number,
// e-mail, delivery address, and other information.
//
// A jCard consists of an array of properties (e.g. "fn", "tel") describing the
// individual or entity. Properties may be repeated, e.g. to represent multiple
// telephone numbers. RFC6350 documents a set of standard properties.
//
// RFC7095 describes the jCard JSON document format, which looks like:
//   ["vcard", [
//     [
//       ["version", {}, "text", "4.0"],
//       ["fn", {}, "text", "Joe Appleseed"],
//       ["tel", {
//             "type":["work", "voice"],
//           },
//           "uri",
//           "tel:+1-555-555-1234;ext=555"
//       ],
//       ...
//     ]
//   ]
//
// This package implements a jCard decoder.
package jcard

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Type JCard represents a jCard.
type JCard struct {
	// List of jCard properties.
	Properties []*Property

	nameLookup map[string][]*Property
}

// Type Property represents a single jCard property.
//
// Each jCard property has four fields, these are:
//    Name   Parameters                  Type   Value
//    -----  --------------------------  -----  -----------------------------
//   ["tel", {"type":["work", "voice"]}, "uri", "tel:+1-555-555-1234;ext=555"]
type Property struct {
	Name       string
	Parameters map[string][]string
	Type       string

	// A property value can be a simple type (string/float64/bool/nil), or be
	// an array. Arrays can be nested, and can contain a mixture of types.
	//
	// Value is one of the following:
	//   * string
	//   * float64
	//   * bool
	//   * nil
	//   * []interface{}. Can contain a mixture of these five types.
	//
	// To retrieve the property value flattened into a []string, use Values().
	Value interface{}
}

// Values returns a simplified representation of the Property value.
//
// This is convenient for accessing simple unstructured data (e.g. "fn", "tel").
//
// The simplified []string representation is created by flattening the
// (potentially nested) Property value, and converting all values to strings.
func (p *Property) Values() []string {
	strings := make([]string, 0, 1)

	p.appendValueStrings(p.Value, &strings)

	return strings
}

func (p *Property) appendValueStrings(v interface{}, strings *[]string) {
	switch v := v.(type) {
	case nil:
		*strings = append(*strings, "")
	case bool:
		*strings = append(*strings, strconv.FormatBool(v))
	case float64:
		*strings = append(*strings, strconv.FormatFloat(v, 'e', -1, 64))
	case string:
		*strings = append(*strings, v)
	case []interface{}:
		for _, v2 := range v {
			p.appendValueStrings(v2, strings)
		}
	default:
		panic("Unknown type")
	}

}

// String returns the jCard as a multiline human readable string. For example:
//
//   jCard[
//     version (type=text, parameters=map[]): [4.0]
//     mixed (type=text, parameters=map[]): [abc true 42 <nil> [def false 43]]
//   ]
//
// This is intended for debugging only, and is not machine parsable.
func (j *JCard) String() string {
	s := make([]string, 0, len(j.Properties))

	for _, s2 := range j.Properties {
		s = append(s, s2.String())
	}

	return "jCard[\n" + strings.Join(s, "\n") + "\n]"
}

// String returns the Property as a human readable string. For example:
//
//     mixed (type=text, parameters=map[]): [abc true 42 <nil> [def false 43]]
//
// This is intended for debugging only, and is not machine parsable.
func (p *Property) String() string {
	return fmt.Sprintf("  %s (type=%s, parameters=%v): %v", p.Name, p.Type, p.Parameters, p.Value)
}

// NewJCard creates a JCard from jsonDocument.
func NewJCard(jsonDocument []byte) (*JCard, error) {
	var top []interface{}
	err := json.Unmarshal(jsonDocument, &top)

	if err != nil {
		return nil, err
	}

	if len(top) != 2 {
		return nil, jCardError("structure is not a JCard (expected len=2 top level array)")
	} else if s, ok := top[0].(string); !(ok && s == "vcard") {
		return nil, jCardError("structure is not a JCard (missing 'vcard')")
	}

	var properties []interface{}

	properties, ok := top[1].([]interface{})
	if !ok {
		return nil, jCardError("structure is not a JCard (bad properties array)")
	}

	j := &JCard{
		Properties: make([]*Property, 0, len(properties)),
		nameLookup: make(map[string][]*Property),
	}

	var p interface{}
	for _, p = range top[1].([]interface{}) {
		var a []interface{}
		var ok bool
		a, ok = p.([]interface{})

		if !ok {
			return nil, jCardError("JCard property was not an array")
		} else if len(a) < 3 {
			return nil, jCardError("JCard property too short (>=3 array elements required)")
		}

		name, ok := a[0].(string)

		if !ok {
			return nil, jCardError("JCard property name invalid")
		}

		var parameters map[string][]string
		var err error
		parameters, err = readParameters(a[1])

		if err != nil {
			return nil, err
		}

		propertyType, ok := a[2].(string)

		if !ok {
			return nil, jCardError("JCard property type invalid")
		}

		var value interface{}
		if len(a) == 4 {
			value, err = readValue(a[3], 0)
		} else {
			value, err = readValue(a[3:], 0)
		}

		if err != nil {
			return nil, err
		}

		property := &Property{
			Name:       name,
			Type:       propertyType,
			Parameters: parameters,
			Value:      value,
		}

		j.Properties = append(j.Properties, property)
		j.nameLookup[name] = append(j.nameLookup[name], property)

	}

	fmt.Printf("%v\n", j)
	return j, nil
}

// Get returns a list of the jCard Properties with Property name |name|.
func (j *JCard) Get(name string) []*Property {
	var properties []*Property

	properties, _ = j.nameLookup[name]

	return properties
}

func jCardError(e string) error {
	return fmt.Errorf("JCard error: %s", e)
}

func readParameters(p interface{}) (map[string][]string, error) {
	params := map[string][]string{}

	if _, ok := p.(map[string]interface{}); !ok {
		return nil, jCardError("JCard parameters invalid")
	}

	for k, v := range p.(map[string]interface{}) {
		if s, ok := v.(string); ok {
			params[k] = append(params[k], s)
		} else if arr, ok := v.([]interface{}); ok {
			for _, value := range arr {
				if s, ok := value.(string); ok {
					params[k] = append(params[k], s)
				}
			}
		}
	}

	return params, nil
}

func readValue(value interface{}, depth int) (interface{}, error) {
	switch value := value.(type) {
	case nil:
		return nil, nil
	case string:
		return value, nil
	case bool:
		return value, nil
	case float64:
		return value, nil
	case []interface{}:
		if depth == 3 {
			return "", jCardError("Structured value too deep")
		}

		result := make([]interface{}, 0, len(value))

		for _, v2 := range value {
			v3, err := readValue(v2, depth+1)

			if err != nil {
				return nil, err
			}

			result = append(result, v3)
		}

		return result, nil
	default:
		return nil, jCardError("Unknown JSON datatype in JCard value")
	}
}
