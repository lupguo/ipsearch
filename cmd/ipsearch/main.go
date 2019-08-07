// ipsearch是ip-search的命令行工具
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipserach"
	"log"
	"time"
)

var (
	ip, proxy, mode string
	debug           bool
	timeout         time.Duration
)

func main() {
	// arg from cmdline
	flag.StringVar(&ip, "ip", "myip", "ip to search, myip is current ip")
	flag.StringVar(&proxy, "proxy", "", "request by proxy, using for debug")
	flag.StringVar(&mode, "mode", "text", "response content mode (json|text)")
	flag.BoolVar(&debug, "debug", false, "debug for request response content ")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set http request timeout seconds")
	flag.Parse()

	// ip search
	ips := &ipserach.Ips{
		Debug:   debug,
		Proxy:   proxy,
		Timeout: timeout,
	}
	ipsRs, err := ips.Search(ip)
	if err != nil {
		log.Fatalf("ip serach error: %s", err)
	}

	// out by json format
	msg, err := ipsRs.Message(mode)
	if err != nil {
		log.Fatalf("ip serach message show error: %s", err)
	}

	// output msg
	fmt.Println(msg)
}
