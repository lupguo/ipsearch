package handler

import (
	"fmt"
	"github.com/lupguo/ipsearch/config"
	"github.com/lupguo/ipsearch/ipsclient"
	"log"
	"net"
	"net/http"
)

//Ipsearch 通过http查询信息
func Ipsearch(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ips := ipsclient.NewIps(config.Get().Debug, r.FormValue("proxy"), config.Get().Timeout)
	rs, err := ips.Search(getIP(r))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg, err := rs.Render("json")
	if err != nil {
		http.Error(w, fmt.Sprintf("ip serach message show error: %s", err), http.StatusInternalServerError)
		return
	}
	log.Println(msg)
	_, _ = fmt.Fprint(w, msg)
}

// getIP 获取HTTP请求信息中的ip内容
func getIP(r *http.Request) string {
	ip := r.FormValue("ip")
	if ip == "" {
		realIp := r.Header.Get("X-Real-IP")
		if realIp != "" {
			return realIp
		}
		remoteIp, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && remoteIp != "127.0.0.1" {
			return remoteIp
		}
	}
	return "myip"
}
