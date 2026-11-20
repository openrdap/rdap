package rdap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/openrdap/rdap/bootstrap"
)

// Do must be safe for concurrent use on a shared Client. Run with -race.
func TestClientDoConcurrent(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/dns.json"):
			fmt.Fprintf(w, `{"version":"1.0","publication":"2024-01-01T00:00:00Z","services":[[["com"],["%s/rdap/"]]]}`, srv.URL)
		case strings.HasPrefix(r.URL.Path, "/rdap/domain/"):
			w.Header().Set("Content-Type", "application/rdap+json")
			fmt.Fprintf(w, `{"objectClassName":"domain","ldhName":"%s"}`, strings.TrimPrefix(r.URL.Path, "/rdap/domain/"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	base, _ := url.Parse(srv.URL + "/")
	client := &Client{Bootstrap: &bootstrap.Client{BaseURL: base}}

	var wg sync.WaitGroup
	errs := make(chan error, 32)
	for i := range 32 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := client.Do(NewDomainRequest(fmt.Sprintf("example%d.com", i)))
			if err != nil {
				errs <- err
				return
			}
			if d, ok := resp.Object.(*Domain); !ok || d.LDHName != fmt.Sprintf("example%d.com", i) {
				errs <- fmt.Errorf("unexpected object %#v", resp.Object)
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func TestRequestURLKeepsServer(t *testing.T) {
	server, _ := url.Parse("https://rdap.example/base/?key=1#frag")
	req := &Request{Type: DomainRequest, Query: "example.com", Server: server}
	got := req.URL().String()
	if server.RawQuery != "key=1" || server.Fragment != "frag" {
		t.Fatalf("URL() modified the Server URL: %s", server)
	}
	if want := "https://rdap.example/base/domain/example.com?key=1#frag"; got != want {
		t.Fatalf("URL() = %s, want %s", got, want)
	}
}
