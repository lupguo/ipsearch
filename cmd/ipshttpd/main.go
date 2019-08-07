// ipshttpd 是ip search httpd服务器，支持通过http查询ip相关信息
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipsearch"
	"log"
	"net/http"
)

var (
	listen  string
	version bool
)

func init() {
	flag.StringVar(&listen, "listen", "127.0.0.1:8680", "the listen address for ip search http server")
	flag.BoolVar(&version, "version", false, "ipsearch version")
	flag.Parse()
}

func main() {
	// version
	if version {
		fmt.Printf("ipshttpd %s", ipsearch.Version())
		return
	}

	// route register
	routeRegister()

	// server running
	log.Printf("ipshttpd(%s) listen on %s", ipsearch.Version(), listen)
	log.Fatalln(http.ListenAndServe(listen, nil))
}

// 路由注册
func routeRegister() {
	// default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r, r.RequestURI, r.Form)
		helpMsg := fmt.Sprintf(`Version %s
Usage:
	//search current client ip information
	curl localhost:8680/ips

	//search for target ip information  
	curl localhost:8680/ips?ip=targetIp	
`, ipsearch.Version())
		fmt.Fprintf(w, helpMsg)
	})

	// ip search handler
	http.HandleFunc("/ips", ipsHandler)
}

// ipsHandler 用于IP查询处理
func ipsHandler(w http.ResponseWriter, r *http.Request) {
	// do ip search
	msg, err := func() (msg string, err error) {
		// parse form get params
		r.ParseForm()
		ip := r.FormValue("ip")

		// ip search
		ips := &ipsearch.Ips{
			Debug:   false,
			Proxy:   r.FormValue("proxy"),
			Timeout: 0,
		}
		ipsRs, err := ips.Search(ip)
		if err != nil {
			return fmt.Sprintf("ip serach error: %s", err), nil
		}

		// out by json format
		msg, err = ipsRs.Message("json")
		if err != nil {
			return fmt.Sprintf("ip serach message show error: %s", err), nil
		}
		return
	}()

	// result handler
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		_, _ = fmt.Fprint(w, msg)
	}
}
