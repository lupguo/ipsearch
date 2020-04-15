// ipsearch是ipsearch的命令行工具
package main

import (
	"flag"
	"fmt"
	"github.com/lupguo/ipsearch/config"
	"github.com/lupguo/ipsearch/ipsclient"
	"github.com/lupguo/ipsearch/ipshttpd"
	"github.com/lupguo/ipsearch/ipsutil"
	"github.com/lupguo/ipsearch/version"
	"log"
)

var cfg = config.Get()

func main() {
	flag.Parse()
	if cfg.Debug {
		log.Printf("%#v\n", cfg)
	}
	version.ShowVersion(cfg.Version)
	if cfg.Listen != "" {
		ipshttpd.Main()
	}
	cmdline()
}

//cmdline 处理Ipsclient
func cmdline() {
	ips := ipsclient.NewIps(cfg.Debug, cfg.Proxy, cfg.Timeout)
	rs, err := ips.Search(cfg.Ip)
	ipsutil.FatalOnError(err, "ip search failed.")

	// Ip Response
	msg, err := rs.Render(cfg.Format)
	ipsutil.FatalOnError(err, "ip response render failed.")

	// Ip Output
	fmt.Println(msg)
}
