// ipshttpd 是ip search httpd服务器，支持通过http查询ip相关信息
package main

import (
	"flag"
	"fmt"
	"github.com/tkstorm/ip-search/ipsearch"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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
		fmt.Println("ipshttpd " + ipsearch.Version())
		return
	}

	// route register
	routeRegister()

	// server running
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		select {
		case <-c:
			log.Println("ipshttpd is shutdown")
		}
		os.Exit(1)
	}()
	log.Printf("ipshttpd listen on %s, %s", listen, ipsearch.Version())
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

		// ip search
		ips := &ipsearch.Ips{
			Debug:   true,
			Proxy:   r.FormValue("proxy"),
			Timeout: 0,
		}

		// request search ip pick
		ip := r.FormValue("ip")
		if ip == "" {
			// host ip
			hostIp, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil && hostIp != "127.0.0.1" {
				ip = hostIp
			}
			// by proxy request
			realIp := r.Header.Get("X-Real-IP")
			if realIp != "" {
				ip = realIp
			}
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
		log.SetOutput(os.Stderr)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.SetOutput(os.Stdout)
		log.Println(msg)
		_, _ = fmt.Fprint(w, msg)
	}
}
