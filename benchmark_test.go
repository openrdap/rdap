// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Benchmarks for hot-path string handling and the reflection-based decoder.
// These exist as regression guards for the optimizations in escapePath,
// Printer.cleanString, VCardProperty.Values, and the chooseFields type-plan
// cache.
//
// Run with:
//
//	go test -run '^$' -bench . -benchmem ./...
package rdap

import (
	"os"
	"testing"
)

// BenchmarkDecodeDomain decodes a realistic nested domain response (entities,
// nameservers, links). It is the headline benchmark for the chooseFields
// type-plan cache, which resolves a struct type's decodable fields once instead
// of on every decode.
func BenchmarkDecodeDomain(b *testing.B) {
	blob, err := os.ReadFile("test/testdata/rdap/rdap.nic.cz/domain-example.cz.json")
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := NewDecoder(blob).Decode(); err != nil {
			b.Fatal(err)
		}
	}
}

// escapePath runs on every request path. The clean case (no byte needs
// escaping) is overwhelmingly common and should not allocate.
var (
	benchEscapePathClean = "rdap.example.com"
	benchEscapePathDirty = "xn--n3h.example/path with spaces & symbols"
)

func BenchmarkEscapePathClean(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = escapePath(benchEscapePathClean)
	}
}

func BenchmarkEscapePathDirty(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = escapePath(benchEscapePathDirty)
	}
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
	for i := 0; i < b.N; i++ {
		for _, s := range benchCleanStringInputs {
			_ = p.cleanString(s)
		}
	}
}

func BenchmarkCleanStringDirty(b *testing.B) {
	p := &Printer{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.cleanString(benchCleanStringDirty)
	}
}

// BenchmarkVCardValues guards the flattening cost of VCardProperty.Values,
// which Tel/Fax now call once per property rather than twice.
func BenchmarkVCardValues(b *testing.B) {
	p := &VCardProperty{
		Name:       "tel",
		Type:       "uri",
		Parameters: map[string][]string{"type": {"voice"}},
		Value:      "tel:+1.5551234567",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.Values()
	}
}
