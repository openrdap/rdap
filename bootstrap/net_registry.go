// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
)

type NetRegistry struct {
	Networks map[int][]NetEntry

	numIPBytes int
}

type NetEntry struct {
	Net  *net.IPNet
	URLs []*url.URL
}

type netEntrySorter []NetEntry

func (a netEntrySorter) Len() int {
	return len(a)
}

func (a netEntrySorter) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a netEntrySorter) Less(i int, j int) bool {
	return bytes.Compare(a[i].Net.IP, a[j].Net.IP) <= 0
}

func NewNetRegistry(json []byte, ipVersion int) (*NetRegistry, error) {
	if ipVersion != 4 && ipVersion != 6 {
		return nil, fmt.Errorf("Unknown IP version %d", ipVersion)
	}

	var registry *registryFile
	registry, err := parse(json)

	if err != nil {
		return nil, fmt.Errorf("Error parsing net registry file: %s", err)
	}

	n := &NetRegistry{
		Networks:   map[int][]NetEntry{},
		numIPBytes: numIPBytesForVersion(ipVersion),
	}

	var cidr string
	var urls []*url.URL
	for cidr, urls = range registry.Entries {
		_, ipNet, err := net.ParseCIDR(cidr)

		if err != nil {
			continue
		} else if len(ipNet.IP) != n.numIPBytes {
			continue
		}

		size, _ := ipNet.Mask.Size()
		n.Networks[size] = append(n.Networks[size], NetEntry{Net: ipNet, URLs: urls})
	}

	for _, nets := range n.Networks {
		sort.Sort(netEntrySorter(nets))
	}

	return n, nil
}

func (n *NetRegistry) Lookup(input string) (*Result, error) {
	if !strings.ContainsAny(input, "/") {
		// Convert IP address to CIDR format, with a /32 or /128 mask.
		input = fmt.Sprintf("%s/%d", input, n.numIPBytes*8)
	}

	_, lookupNet, err := net.ParseCIDR(input)

	if err != nil {
		return nil, err
	}

	if len(lookupNet.IP) != n.numIPBytes {
		return nil, errors.New("Lookup address has wrong IP protocol")
	}

	lookupMask, _ := lookupNet.Mask.Size()

	var bestEntry string
	var bestURLs []*url.URL
	var bestMask int

	var mask int
	var nets []NetEntry
	for mask, nets = range n.Networks {
		if mask < bestMask || mask > lookupMask {
			continue
		}

		index := sort.Search(len(nets), func(i int) bool {
			net := nets[i].Net
			return net.Contains(lookupNet.IP) || bytes.Compare(net.IP, lookupNet.IP) >= 0
		})

		if index == len(nets) || !nets[index].Net.Contains(lookupNet.IP) {
			continue
		}

		bestEntry = nets[index].Net.String()
		bestMask = mask
		bestURLs = nets[index].URLs
	}

	return &Result{
		Query: input,
		Entry: bestEntry,
		URLs:  bestURLs,
	}, nil
}

func numIPBytesForVersion(ipVersion int) int {
	len := 0

	switch ipVersion {
	case 4:
		len = net.IPv4len
	case 6:
		len = net.IPv6len
	default:
		panic("Unknown IP version")
	}

	return len
}
