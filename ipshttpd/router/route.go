package ipshttpd

import (
	"github.com/tkstorm/ip-search/ipshttpd/handler"
	"net/http"
)

// registeRoute 路由注册
func registeRoute() {
	http.HandleFunc("/", handler.HelpMessage)
	http.HandleFunc("/ips", handler.Ipsearch)
}

