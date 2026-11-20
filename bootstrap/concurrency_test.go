package bootstrap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

func newDNSServer(t *testing.T) *url.URL {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/dns.json") {
			fmt.Fprintf(w, `{"version":"1.0","publication":"2024-01-01T00:00:00Z","services":[[["com"],["%s/rdap/"]]]}`, srv.URL)
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)
	base, _ := url.Parse(srv.URL + "/")
	return base
}

// Verbose may call back into the client. Before, Lookup held the lock while
// calling Verbose, so this deadlocked.
func TestLookupVerboseMayCallAccessors(t *testing.T) {
	base := newDNSServer(t)
	b := &Client{BaseURL: base}

	var calls int
	answer, err := b.LookupVerbose(&Question{RegistryType: DNS, Query: "example.com"}, func(string) {
		calls++
		_ = b.DNS()
	})
	if err != nil {
		t.Fatal(err)
	}
	if answer == nil || len(answer.URLs) == 0 {
		t.Fatalf("no answer: %#v", answer)
	}
	if calls == 0 {
		t.Fatal("verbose was not called")
	}
}

// Concurrent lookups with their own verbose callbacks share one client. Run
// with -race.
func TestLookupVerboseSharedClient(t *testing.T) {
	base := newDNSServer(t)
	b := &Client{BaseURL: base}

	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := b.LookupVerbose(&Question{RegistryType: DNS, Query: "example.com"}, func(string) {})
			if err != nil {
				t.Errorf("lookup %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
}
