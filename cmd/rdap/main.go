package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/openrdap/rdap"
	"github.com/openrdap/rdap/bootstrap/cache"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version   = "0.0.1"
	usageText = `OpenRDAP v` + version + ` (www.openrdap.org)

Usage: rdap [OPTIONS] DOMAIN|IP|ASN|ENTITY|NAMESERVER|RDAP-URL
  e.g. rdap google.cz
       rdap 192.0.2.0
       rdap 2001:db8::
       rdap AS2856
       rdap https://rdap.nic.cz/domain/example.cz

       rdap -f registrant -f administrative -f billing amazon.com.br
       rdap --json https://rdap.nic.cz/domain/example.cz
       rdap -s https://rdap.nic.cz -t help

Options:
  -h, --help          Show help message.
  -H, --help-advanced Show help message with advanced options.
  -v, --verbose       Print verbose messages on STDERR.

  -T, --timeout=SECS  Timeout after SECS seconds (default: 120).
  -k, --insecure      Disable SSL certificate verification.

  -E, --experimental  Enable experimental options:
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
`
	advancedUsageText = `Advanced options (query):
  -s  --server=URL    RDAP server to query.
  -t  --type=TYPE     RDAP query type. Normally auto-detected. The types are:
                      - ip
                      - domain
                      - autnum
                      - nameserver
                      - entity
                      - help
                      - url
                      - entity-search-by-handle
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

func main() {
	exitCode := run(os.Args[1:], os.Stdout, os.Stderr)

	os.Exit(exitCode)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	// Setup command line arguments parser.
	app := kingpin.New("rdap", "RDAP command-line client")
	app.HelpFlag.Short('h')
	app.UsageTemplate(usageText)
	app.UsageWriter(stderr)

	// Command line options.
	advancedHelpFlag := app.Flag("help-advanced", "").Short('H').Bool()
	verboseFlag := app.Flag("verbose", "").Short('v').Bool()
	insecureFlag := app.Flag("insecure", "").Short('k').Bool()

	queryType := app.Flag("type", "").Short('t').String()
	fetchRolesFlag := app.Flag("fetch", "").Short('f').Strings()
	serverFlag := app.Flag("server", "").Short('s').String()

	experimentsFlag := app.Flag("exp", "").Strings()

	cacheDirFlag := app.Flag("cache-dir", "").String()
	bootstrapURLFlag := app.Flag("bs-url", "").Default("default").String()
	bootstrapTimeoutFlag := app.Flag("bs-ttl", "").Uint16()

	// Command line query (any remaining non-option arguments).
	var queryArgs *[]string = app.Arg("", "").Strings()

	// Parse command line arguments.
	// The help messages for -h/--help are printed directly by app.Parse().
	_, err := app.Parse(args)
	if err != nil {
		printError(stderr, fmt.Sprintf("Error: %s\n\n%s", err, usageText))
		return 1
	} else if *advancedHelpFlag {
		printError(stderr, fmt.Sprintf("%s\n%s", usageText, advancedUsageText))
		return 0
	}

	// Supported experimental options.
	experiments := map[string]bool{
		"test_rdap_net": false,
		"object_tag":    false,
	}

	// Enable any experimental options required.
	for _, e := range *experimentsFlag {
		if _, ok := experiments[e]; ok {
			experiments[e] = true
		} else {
			printError(stderr, fmt.Sprintf("Error: unknown experiment '%s'", e))
			return 1
		}
	}

	// Exactly one argument is required (i.e. the domain/ip/url/etc), unless
	// we're making a help query.
	if *queryType != "help" && len(*queryArgs) == 0 {
		printError(stderr, "Query object required, e.g. rdap google.cz")
		return 1
	}

	queryText := ""
	if len(*queryArgs) > 0 {
		queryText = (*queryArgs)[0]
	}

	// Construct the query.
	var query *rdap.Query
	switch *queryType {
	case "":
		query = rdap.NewAutoQuery(queryText)
	case "help":
		query = rdap.NewHelpQuery()
	case "domain", "dns":
		query = rdap.NewDomainQuery(queryText)
	case "autnum", "as", "asn":
		autnum := strings.ToUpper(queryText)
		autnum = strings.TrimPrefix(autnum, "AS")
		result, err := strconv.ParseUint(autnum, 10, 32)

		if err != nil {
			printError(stderr, "Invalid ASN")
			return 1
		}
		query = rdap.NewAutnumQuery(uint32(result))
	case "ip":
		ip := net.ParseIP(queryText)
		if ip == nil {
			printError(stderr, "Invalid IP")
			return 1
		}
		query = rdap.NewIPQuery(ip)
	case "nameserver", "ns":
		query = rdap.NewNameserverQuery(queryText)
	case "entity":
		query = rdap.NewEntityQuery(queryText)
	case "url":
		fullURL, err := url.Parse(queryText)
		if err != nil {
			printError(stderr, fmt.Sprintf("Unable to parse URL: %s", err))
			return 1
		}
		query = rdap.NewURLQuery(fullURL)
	case "entity-search":
		query = rdap.NewSearchQuery(rdap.EntitySearch, queryText)
	case "entity-search-by-handle":
		query = rdap.NewSearchQuery(rdap.EntitySearchByHandle, queryText)
	case "domain-search":
		query = rdap.NewSearchQuery(rdap.DomainSearch, queryText)
	case "domain-search-by-nameserver":
		query = rdap.NewSearchQuery(rdap.DomainSearchByNameserver, queryText)
	case "domain-search-by-nameserver-ip":
		query = rdap.NewSearchQuery(rdap.DomainSearchByNameserverIP, queryText)
	case "nameserver-search":
		query = rdap.NewSearchQuery(rdap.NameserverSearch, queryText)
	case "nameserver-search-by-nameserver-ip":
		query = rdap.NewSearchQuery(rdap.NameserverSearchByNameserverIP, queryText)
	default:
		printError(stderr, fmt.Sprintf("Unknown query type %s", queryType))
		return 1
	}

	// Determine the query server.
	if query.HasServer() {
		if *serverFlag != "" {
			printError(stderr, fmt.Sprintf("--server option cannot be used with query type %s", query.Type()))
			return 1
		}
	}

	// Server URL specified (--server)?
	if *serverFlag != "" {
		var err error
		query, err = query.UsingServerURL(*serverFlag)

		if err != nil {
			printError(stderr, fmt.Sprintf("--server error: %s", err))
			return 1
		}
	}

	var client *rdap.Client = rdap.NewClient()

	// Print verbose messages on STDERR?
	if *verboseFlag {
		client.Verbose = func(text string) {
			fmt.Fprintf(stderr, "# %s\n", text)
		}
	}

	// Custom bootstrap cache type/directory?
	if cacheDirFlag == nil {
		// Disk cache, default location.
		client.Bootstrap.Cache = cache.NewDiskCache()
	} else {
		if *cacheDirFlag != "" {
			// Disk cache with custom directory.
			dc := cache.NewDiskCache()
			dc.Dir = *cacheDirFlag

			client.Bootstrap.Cache = dc
		} else {
			// Disk cache disabled, use default memory cache.
		}
	}

	// Custom bootstrap service URL?
	if *bootstrapURLFlag != "default" {
		baseURL, err := url.Parse(*bootstrapURLFlag)
		if err != nil {
			printError(stderr, fmt.Sprintf("Bootstrap URL error: %s", err))
			return 1
		}

		client.Bootstrap.BaseURL = baseURL
	}

	// Custom bootstrap cache timeout?
	if bootstrapTimeoutFlag != nil {
		client.Bootstrap.Cache.SetTimeout(time.Duration(*bootstrapTimeoutFlag) * time.Second)
	}

	var resp *rdap.Response
	resp, err = client.Query(query)

	if err != nil {
		printError(stderr, fmt.Sprintf("Error: %s", err))
		return 1
	}

	_ = resp
	_ = insecureFlag
	_ = queryType
	_ = verboseFlag
	_ = fetchRolesFlag
	_ = query

	return 0
}

func printError(stderr io.Writer, text string) {
	fmt.Fprintf(stderr, "# %s\n", text)
}
