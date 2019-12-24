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

func init() {
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
	if version {
		fmt.Println("ipsearch " + ipsearch.Version())
		return
	}
	ips := ipsearch.NewIps(debug, proxy, timeout)
	ipsRs, err := ips.Search(ip)
	if err != nil {
		fmt.Printf("ipserach error: %s", err)
		return
	}
	msg, err := ipsRs.Message(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg)
}
