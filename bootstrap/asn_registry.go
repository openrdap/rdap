// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type ASNRegistry struct {
	ASNs []ASNRange
}

// ASNRange represents a range of AS numbers and their RDAP base URLs.
//
// An ASNRange represents a single AS number when MinASN==MaxASN.
type ASNRange struct {
	MinASN uint32 // First AS number in the range.
	MaxASN uint32 // Last AS number in the range.
	URLs   []*url.URL // RDAP base URLs.
}

// String returns "ASxxxx" for a single AS, or "ASxxxx-ASyyyy" for a range.
func (a ASNRange) String() string {
	if a.MinASN == a.MaxASN {
		return fmt.Sprintf("AS%d", a.MinASN)
	}

	return fmt.Sprintf("AS%d-AS%d", a.MinASN, a.MaxASN)
}

type asnRangeSorter []ASNRange

func (a asnRangeSorter) Len() int {
	return len(a)
}

func (a asnRangeSorter) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a asnRangeSorter) Less(i int, j int) bool {
	return a[i].MinASN < a[j].MinASN
}

// NewASNRegistry creates a queryable ASN registry from an ASN registry JSON document.
//
// The document format is specified in https://tools.ietf.org/html/rfc7484#section-5.3.
func NewASNRegistry(json []byte) (*ASNRegistry, error) {
	var registry *registryFile
	registry, err := parse(json)

	if err != nil {
		return nil, fmt.Errorf("Error parsing ASN registry: %s\n", err)
	}

	a := make([]ASNRange, 0, len(registry.Entries))

	var asn string
	var urls []*url.URL
	for asn, urls = range registry.Entries {
		minASN, maxASN, err := parseASNRange(asn)

		if err != nil {
			continue
		}

		a = append(a, ASNRange{MinASN: minASN, MaxASN: maxASN, URLs: urls})
	}

	sort.Sort(asnRangeSorter(a))

	return &ASNRegistry{
		ASNs: a,
	}, nil
}

func (a *ASNRegistry) Lookup(input string) (*Result, error) {
	var asn uint32
	asn, err := parseASN(input)

	if err != nil {
		return nil, err
	}

	index := sort.Search(len(a.ASNs), func(i int) bool {
		return asn <= a.ASNs[i].MaxASN
	})

	var entry string
	var urls []*url.URL

	if index != len(a.ASNs) && (asn >= a.ASNs[index].MinASN && asn <= a.ASNs[index].MaxASN) {
		entry = a.ASNs[index].String()
		urls = a.ASNs[index].URLs
	}

	return &Result{
		Query: string(asn),
		Entry: entry,
		URLs:  urls,
	}, nil
}

func parseASN(asn string) (uint32, error) {
	asn = strings.ToLower(asn)
	asn = strings.TrimLeft(asn, "as")
	result, err := strconv.ParseUint(asn, 10, 32)

	if err != nil {
		return 0, err
	}

	return uint32(result), nil
}

func parseASNRange(asnRange string) (uint32, uint32, error) {
	var minASN uint64
	var maxASN uint64
	var err error

	asns := strings.Split(asnRange, "-")

	if len(asns) != 1 && len(asns) != 2 {
		return 0, 0, errors.New("Malformed ASN range")
	}

	minASN, err = strconv.ParseUint(asns[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	if len(asns) == 2 {
		maxASN, err = strconv.ParseUint(asns[1], 10, 32)
		if err != nil {
			return 0, 0, err
		}
	} else {
		maxASN = minASN
	}

	if minASN > maxASN {
		minASN, maxASN = maxASN, minASN
	}

	return uint32(minASN), uint32(maxASN), nil
}
