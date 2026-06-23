// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/openrdap/rdap/test"
)

// Run with `go test -run TestRunCLI -update ./...` to regenerate the golden
// files after an intentional output change.
var updateGolden = flag.Bool("update", false, "update CLI golden files")

// fixtureServer serves the RDAP test fixtures over HTTP so the CLI can be
// driven end-to-end via --server without touching the network.
func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/domain/example.cz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/rdap+json")
		_, _ = w.Write(test.LoadFile("rdap/rdap.nic.cz/domain-example.cz.json"))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestRunCLIGolden drives RunCLI end-to-end for a range of commands and output
// modes, comparing exit code, stdout, and stderr against committed golden files.
// Network-backed cases point --server at an in-process fixture server and use an
// in-memory cache (--cache-dir "") so nothing touches the real filesystem.
func TestRunCLIGolden(t *testing.T) {
	srv := fixtureServer(t)

	withServer := func(args ...string) []string {
		return append(args, "--cache-dir", "", "--server", srv.URL)
	}

	cases := []struct {
		name string
		args []string
	}{
		{"domain-text", withServer("example.cz")},
		{"domain-json", withServer("--json", "example.cz")},
		{"domain-whois", withServer("--whois", "example.cz")},
		{"domain-raw", withServer("--raw", "example.cz")},
		{"version", []string{"--version"}},
		{"no-args", []string{}},
		{"unknown-type", []string{"--type", "bogus", "example.cz"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := RunCLI(tc.args, &stdout, &stderr, CLIOptions{})

			got := fmt.Sprintf("exit: %d\n--- stdout ---\n%s\n--- stderr ---\n%s",
				code, stdout.String(), stderr.String())

			assertGolden(t, tc.name, got)
		})
	}
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()

	goldenPath := filepath.Join("testdata", "cli", name+".golden")

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}

		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden (run `go test -run TestRunCLI -update` to create): %s", err)
	}

	if got != string(want) {
		t.Errorf("output mismatch for %s (run `go test -run TestRunCLI -update` to regenerate)\n--- got ---\n%s\n--- want ---\n%s",
			name, got, want)
	}
}
