package router

import (
	"github.com/lupguo/ipsearch/ipshttpd/handler"
	"net/http"
)

// Register 路由注册
func Register() {
	http.HandleFunc("/", handler.HelpMessage)
	http.HandleFunc("/ips", handler.Ipsearch)
}

