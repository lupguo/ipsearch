// ipsearch是ipsearch的命令行工具
package main

import (
	"flag"
	"fmt"
	"github.com/lupguo/ipsearch/ipsclient"
	"github.com/lupguo/ipsearch/ipshttpd"
	"github.com/lupguo/ipsearch/ipsutil"
	"github.com/lupguo/ipsearch/version"
	"time"
)

var (
	ip, proxy, format string
	listen            string
	debug, ver        bool
	timeout           time.Duration
)

func init() {
	flag.StringVar(&ip, "ip", "", "the IP to be search, the default is the IP of the machine currently executing the cmdline")
	flag.StringVar(&proxy, "proxy", "", "http proxy using for debugging, no proxy by default, eg http://127.0.0.1:8888")
	flag.StringVar(&format, "format", "text", "response message format, default is json (json|text)")
	flag.BoolVar(&debug, "debug", false, "debug for request response content ")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set http request timeout seconds")
	flag.StringVar(&listen, "listen", "", "the listen address for ip search http server, eg 127.0.0.1:6100")
	flag.BoolVar(&ver, "version", false, "ipsearch version")
	flag.Parse()
}

func main() {
	version.ShowVersion(ver)
	if listen != "" {
		ipshttpd.Main(listen)
	}

	cmdline()
}

//cmdline 处理Ipsclient
func cmdline() {
	ips := ipsclient.NewIps(debug, proxy, timeout)
	rs, err := ips.Search(ip)
	ipsutil.FatalOnError(err, "ip search failed.")

	// Ip Response
	msg, err := rs.Render(format)
	ipsutil.FatalOnError(err, "ip response render failed.")

	// Ip Output
	fmt.Println(msg)
}
