// ipshttpd 是ip search httpd服务器
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipserach"
	"log"
	"net/http"
)


var listen string

// httpd服务
func main() {
	flag.StringVar(&listen, "listen", "127.0.0.1:8680", "the listen address for ip search http server")
	flag.Parse()
	log.Println("ip search httpd listen on " + listen)

	// route register
	routeRegister()

	// server running
	log.Fatalln(http.ListenAndServe(listen, nil))
}

// 路由注册
func routeRegister() {
	http.HandleFunc("/ips", ipSearch)
	http.HandleFunc("/ips/", ipSearch)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r.RequestURI, r.Form)
//
//		helpMsg := `
//
//`
		fmt.Fprintf(w, r.RequestURI)
	})
}

// ipSearch
func ipSearch(w http.ResponseWriter, r *http.Request) {
	// parse form get params
	r.ParseForm()
	ip := r.FormValue("ip")

	// ip search
	ips := &ipserach.Ips{
		Debug:   false,
		Proxy:   r.FormValue("proxy"),
		Timeout: 0,
	}
	ipsRs, err := ips.Search(ip)
	if err != nil {
		log.Printf("ip serach error: %s", err)
		return
	}

	// out by json format
	msg, err := ipsRs.Message("json")
	if err != nil {
		http.Error(w, fmt.Sprintf("ip serach message show error: %s", err), http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, msg)
}


