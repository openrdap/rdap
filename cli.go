package rdap

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/openrdap/rdap/bootstrap"
	"github.com/openrdap/rdap/bootstrap/cache"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version   = "OpenRDAP v0.0.1"
	usageText = version + `
(www.openrdap.org)

Usage: rdap [OPTIONS] DOMAIN|IP|ASN|ENTITY|NAMESERVER|RDAP-URL
  e.g. rdap example.cz
       rdap 192.0.2.0
       rdap 2001:db8::
       rdap AS2856
       rdap https://rdap.nic.cz/domain/example.cz

       rdap -f registrant -f administrative -f billing amazon.com.br
       rdap --json https://rdap.nic.cz/domain/example.cz
       rdap -s https://rdap.nic.cz -t help

Options:
  -h, --help          Show help message.
  -v, --verbose       Print verbose messages on STDERR.

  -T, --timeout=SECS  Timeout after SECS seconds (default: 30).
  -k, --insecure      Disable SSL certificate verification.

  -e, --experimental  Enable some experimental options:
                      - Use the bootstrap service https://test.rdap.net/rdap
                      - Enable object tag support

Contact Information Fetch Options:
  -f, --fetch=all     Fetch all available contact information (default).
  -f, --fetch=none    Disable additional RDAP requests for contact information.
  -f, --fetch=ROLE    Fetch additional contact information for the role
                      ROLE only. The regular WHOIS roles are:
                      registrant, administrative, billing.
Output Options:
      --text          Output WHOIS style, plain text format (default).
  -j, --json          Output JSON, pretty-printed format.
  -J, --compact       Output JSON, compact (one line) format.
  -r, --raw           Output the raw server response. Forces --fetch=none.

Advanced options (query):
  -s  --server=URL    RDAP server to query.
  -t  --type=TYPE     RDAP query type. Normally auto-detected. The types are:
                      - ip
                      - domain
                      - autnum
                      - nameserver
                      - entity
                      - help
                      - url
                      - domain-search
                      - domain-search-by-nameserver
                      - domain-search-by-nameserver-ip
                      - nameserver-search
                      - nameserver-search-by-ip
                      - entity-search
                      - entity-search-by-handle
                      The servers for domain, ip, autnum, url queries can be
                      determined automatically. Otherwise, the RDAP server
                      (--server=URL) must be specified.
      --strict-fetch  Exit with an error when a contact information fetch
                      (--fetch=) fails. By default these errors are ignored.

Advanced options (bootstrapping):
      --cache-dir=DIR Bootstrap cache directory to use. Specify empty string
                      to disable bootstrap caching. The directory is created
                      automatically as needed. (default: $HOME/.openrdap).
      --bs-url=URL    Bootstrap service URL (default: https://data.iana.org/rdap)
      --bs-ttl=SECS   Bootstrap cache time in seconds (default: 3600)

Advanced options (experiments):
      --exp=test_rdap_net  Use the bootstrap service https://test.rdap.net/rdap
      --exp=object_tag     Enable object tag support
                           (draft-hollenbeck-regext-rdap-object-tag)
`
)

const (
	experimentalBootstrapURL = "https://test.rdap.net/rdap"
)

// CLIOptions specifies options for the command line client.
type CLIOptions struct {
	// Sandbox mode disables the --cache-dir option, to prevent arbitrary writes to
	// the file system.
	//
	// This is used for https://www.openrdap.org/demo.
	Sandbox bool
}

// RunCLI runs the OpenRDAP command line client.
//
// |args| are the command line arguments to use (normally os.Args[1:]).
// |stdout| and |stderr| are the io.Writers for STDOUT/STDERR.
// |options| specifies extra options.
//
// Returns the program exit code.
func RunCLI(args []string, stdout io.Writer, stderr io.Writer, options CLIOptions) int {
	// For duration timer (in --verbose output).
	start := time.Now()

	// Setup command line arguments parser.
	app := kingpin.New("rdap", "RDAP command-line client")
	app.HelpFlag.Short('h')
	app.UsageTemplate(usageText)
	app.UsageWriter(stdout)
	app.ErrorWriter(stderr)

	// Instead of letting kingpin call os.Exit(), flag if it requests to exit
	// here.
	//
	// This lets the function be called in libraries/tests without exiting them.
	terminate := false
	app.Terminate(func(int) {
		terminate = true
	})

	// Command line options.
	verboseFlag := app.Flag("verbose", "").Short('v').Bool()
	timeoutFlag := app.Flag("timeout", "").Short('T').Default("30").Uint16()
	insecureFlag := app.Flag("insecure", "").Short('k').Bool()

	queryType := app.Flag("type", "").Short('t').String()
	fetchRolesFlag := app.Flag("fetch", "").Short('f').Strings()
	serverFlag := app.Flag("server", "").Short('s').String()

	experimentalFlag := app.Flag("experimental", "").Short('e').Bool()
	experimentsFlag := app.Flag("exp", "").Strings()

	cacheDirFlag := app.Flag("cache-dir", "").Default("default").String()
	bootstrapURLFlag := app.Flag("bs-url", "").Default("default").String()
	bootstrapTimeoutFlag := app.Flag("bs-ttl", "").Default("3600").Uint32()

	// Command line query (any remaining non-option arguments).
	queryArgs := app.Arg("", "").Strings()

	// Parse command line arguments.
	// The help messages for -h/--help are printed directly by app.Parse().
	_, err := app.Parse(args)
	if err != nil {
		printError(stderr, fmt.Sprintf("Error: %s\n\n%s", err, usageText))
		return 1
	} else if terminate {
		// Occurs when kingpin prints the --help message.
		return 1
	}

	var verbose func(text string)
	if *verboseFlag {
		verbose = func(text string) {
			fmt.Fprintf(stderr, "# %s\n", text)
		}
	} else {
		verbose = func(text string) {
		}
	}

	verbose(version)
	verbose("")

	verbose("rdap: Configuring query...")

	// Supported experimental options.
	experiments := map[string]bool{
		"test_rdap_net": false,
		"object_tag":    false,
	}

	// Enable experimental options.
	for _, e := range *experimentsFlag {
		if _, ok := experiments[e]; ok {
			experiments[e] = true
			verbose(fmt.Sprintf("rdap: Enabled experiment '%s'", e))
		} else {
			printError(stderr, fmt.Sprintf("Error: unknown experiment '%s'", e))
			return 1
		}
	}

	// Enable the -e selection of experiments?
	if *experimentalFlag {
		verbose("rdap: Enabled -e/--experiments: test_rdap_net, object_tag")
		experiments["test_rdap_net"] = true
		experiments["object_tag"] = true
	}

	// Exactly one argument is required (i.e. the domain/ip/url/etc), unless
	// we're making a help query.
	if *queryType != "help" && len(*queryArgs) == 0 {
		printError(stderr, fmt.Sprintf("Error: %s\n\n%s", "Query object required, e.g. rdap example.cz", usageText))
		return 1
	}

	// Grab the query text.
	queryText := ""
	if len(*queryArgs) > 0 {
		queryText = (*queryArgs)[0]
	}

	// Construct the request.
	var req *Request
	switch *queryType {
	case "":
		req = NewAutoRequest(queryText)
	case "help":
		req = NewHelpRequest()
	case "domain", "dns":
		req = NewDomainRequest(queryText)
	case "autnum", "as", "asn":
		autnum := strings.ToUpper(queryText)
		autnum = strings.TrimPrefix(autnum, "AS")
		result, err := strconv.ParseUint(autnum, 10, 32)

		if err != nil {
			printError(stderr, fmt.Sprintf("Invalid ASN '%s'", queryText))
			return 1
		}
		req = NewAutnumRequest(uint32(result))
	case "ip":
		ip := net.ParseIP(queryText)
		if ip == nil {
			printError(stderr, fmt.Sprintf("Invalid IP '%s'", queryText))
			return 1
		}
		req = NewIPRequest(ip)
	case "nameserver", "ns":
		req = NewNameserverRequest(queryText)
	case "entity":
		req = NewEntityRequest(queryText)
	case "url":
		fullURL, err := url.Parse(queryText)
		if err != nil {
			printError(stderr, fmt.Sprintf("Unable to parse URL '%s': %s", queryText, err))
			return 1
		}
		req = NewRawRequest(fullURL)
	case "entity-search":
		req = NewRequest(EntitySearchRequest, queryText)
	case "entity-search-by-handle":
		req = NewRequest(EntitySearchByHandleRequest, queryText)
	case "domain-search":
		req = NewRequest(DomainSearchRequest, queryText)
	case "domain-search-by-nameserver":
		req = NewRequest(DomainSearchByNameserverRequest, queryText)
	case "domain-search-by-nameserver-ip":
		req = NewRequest(DomainSearchByNameserverIPRequest, queryText)
	case "nameserver-search":
		req = NewRequest(NameserverSearchRequest, queryText)
	case "nameserver-search-by-ip":
		req = NewRequest(NameserverSearchByNameserverIPRequest, queryText)
	default:
		printError(stderr, fmt.Sprintf("Unknown query type '%s'", *queryType))
		return 1
	}

	// Determine the server.
	if req.Server != nil {
		if *serverFlag != "" {
			printError(stderr, fmt.Sprintf("--server option cannot be used with query type %s", req.Type))
			return 1
		}
	}

	// Server URL specified (--server)?
	if *serverFlag != "" {
		serverURL, err := url.Parse(*serverFlag)

		if err != nil {
			printError(stderr, fmt.Sprintf("--server error: %s", err))
			return 1
		}

		if serverURL.Scheme == "" {
			serverURL.Scheme = "http"
		}

		req = req.WithServer(serverURL)

		verbose(fmt.Sprintf("rdap: Using server '%s'", serverURL))
	}

	bs := &bootstrap.Client{}

	// Custom bootstrap cache type/directory?
	if *cacheDirFlag == "" {
		// Disk cache disabled, use memory cache.
		bs.Cache = cache.NewMemoryCache()

		verbose("rdap: Using in-memory cache")
	} else {
		dc := cache.NewDiskCache()
		if *cacheDirFlag != "default" {
			if !options.Sandbox {
				dc.Dir = *cacheDirFlag
			} else {
				verbose(fmt.Sprintf("rdap: Ignored --cache-dir option (sandbox mode enabled)"))
			}
		}

		verbose(fmt.Sprintf("rdap: Using disk cache (%s)", dc.Dir))

		created, err := dc.InitDir()
		if created {
			verbose(fmt.Sprintf("rdap: Cache dir %s mkdir'ed", dc.Dir))
		} else if err != nil {
			printError(stderr, fmt.Sprintf("rdap: Error making cache dir %s", dc.Dir))
			return 1
		}

		bs.Cache = dc
	}

	// Use experimental bootstrap service URL?
	if experiments["test_rdap_net"] && *bootstrapURLFlag == "default" {
		*bootstrapURLFlag = experimentalBootstrapURL

		verbose("rdap: Using test.rdap.net bootstrap service (test_rdap_net experiment)")
	}

	// Custom bootstrap service URL?
	if *bootstrapURLFlag != "default" {
		baseURL, err := url.Parse(*bootstrapURLFlag)
		if err != nil {
			printError(stderr, fmt.Sprintf("Bootstrap URL error: %s", err))
			return 1
		}

		bs.BaseURL = baseURL

		verbose(fmt.Sprintf("rdap: Bootstrap URL set to '%s'", baseURL))
	} else {
		verbose(fmt.Sprintf("rdap: Bootstrap URL is default '%s'", bootstrap.DefaultBaseURL))
	}

	// Custom bootstrap cache timeout?
	if bootstrapTimeoutFlag != nil {
		bs.Cache.SetTimeout(time.Duration(*bootstrapTimeoutFlag) * time.Second)

		verbose(fmt.Sprintf("rdap: Bootstrap cache TTL set to %d seconds", *bootstrapTimeoutFlag))
	}

	// Custom HTTP client. Used to disable TLS certificate verification.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecureFlag},
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{
		HTTP:      httpClient,
		Bootstrap: bs,

		Verbose:                   verbose,
		UserAgent:                 version,
		ServiceProviderExperiment: experiments["object_tag"],
	}

	if *insecureFlag {
		verbose(fmt.Sprintf("rdap: SSL certificate validation disabled"))
	}

	// Set the request timeout.
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(*timeoutFlag)*time.Second)
	defer cancelFunc()
	req = req.WithContext(ctx)

	verbose(fmt.Sprintf("rdap: Timeout is %d seconds", *timeoutFlag))

	// Run the request.
	var resp *Response
	resp, err = client.Do(req)

	verbose("")
	verbose(fmt.Sprintf("rdap: Finished in %s", time.Since(start)))

	if err != nil {
		printError(stderr, fmt.Sprintf("Error: %s", err))
		return 1
	}

	// Insert a blank line to seperate verbose messages/proper output.
	if *verboseFlag {
		fmt.Fprintln(stderr, "")
	}

	// Print the response out in text format.
	printer := &Printer{
		Writer: stdout,

		BriefLinks: true,
	}
	printer.Print(resp.Object)

	_ = fetchRolesFlag

	return 0
}

func printError(stderr io.Writer, text string) {
	fmt.Fprintf(stderr, "# %s\n", text)
}
