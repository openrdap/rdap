package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app              = kingpin.New("rdap", "RDAP command-line client")
	insecureFlag     = app.Flag("insecure", "Disable SSL certificate verification").Short('k').Bool()
	advancedHelpFlag = app.Flag("help-advanced", "Show this help message & more advanced options").Short('H').Bool()

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
	app.HelpFlag.Short('h')
	app.UsageTemplate(usageText)
	command, err := app.Parse(os.Args[1:])

	if err != nil {
		fmt.Printf("Error: %s\n\n%s", err, usageText)
		os.Exit(1)
	}

	if *advancedHelpFlag {
		fmt.Printf("%s\n%s", usageText, advancedUsageText)
		return
	}

	if *insecureFlag {
		fmt.Printf("insecure flag\n")
	}

	switch command {
	default:
		fmt.Println("default")
	}
}
