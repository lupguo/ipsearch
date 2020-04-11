// ipsearch是ip-search的命令行工具
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipsearch"
	"github.com/tkstorm/ip-search/ipsutil"
	version2 "github.com/tkstorm/ip-search/version"
	"time"
)

var (
	ip, proxy, format string
	debug, version    bool
	timeout           time.Duration
)

func init() {
	flag.StringVar(&ip, "ip", "", "the IP to be search, the default is the IP of the machine currently executing the command")
	flag.StringVar(&proxy, "proxy", "", "http proxy using for debugging, no proxy by default, eg http://127.0.0.1:8888")
	flag.StringVar(&format, "format", "text", "response message format, default is json (json|text)")
	flag.BoolVar(&debug, "debug", false, "debug for request response content ")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set http request timeout seconds")
	flag.BoolVar(&version, "version", false, "ipsearch version")
	flag.Parse()
}

func main() {
	version2.ShowVersion(version)

	// Ip Search
	ips := ipsearch.NewIps(debug, proxy, timeout)
	rs, err := ips.Search(ip)
	ipsutil.FatalOnError(err, "ip search failed.")

	// Ip Response
	msg, err := rs.Render(format)
	ipsutil.FatalOnError(err, "ip response render failed.")

	// Ip Output
	fmt.Println(msg)
}
