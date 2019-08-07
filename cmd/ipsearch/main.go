// ipsearch是ip-search的命令行工具
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipsearch"
	"time"
)

var (
	ip, proxy, mode string
	debug, version  bool
	timeout         time.Duration
)

func init()  {
	// arg from cmdline
	flag.StringVar(&ip, "ip", "myip", "ip to search, myip is current ip")
	flag.StringVar(&proxy, "proxy", "", "request by proxy, using for debug")
	flag.StringVar(&mode, "mode", "text", "response content mode (json|text)")
	flag.BoolVar(&debug, "debug", false, "debug for request response content ")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set http request timeout seconds")
	flag.BoolVar(&version, "version", false, "ipsearch version")
	flag.Parse()
}

func main() {
	// handle ip search
	msg, err := func() (msg string, err error) {
		// version
		if version {
			return "ipsearch " + ipsearch.Version(), nil
		}
		// ip search
		ips := &ipsearch.Ips{
			Debug:   debug,
			Proxy:   proxy,
			Timeout: timeout,
		}
		ipsRs, err := ips.Search(ip)
		if err != nil {
			return "", fmt.Errorf("ip serach error: %s", err)
		}
		// out by json format
		return ipsRs.Message(mode)
	}()

	// output search message
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}
}
