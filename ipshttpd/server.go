package ipshttpd

import (
	"github.com/lupguo/ipsearch/config"
	"github.com/lupguo/ipsearch/ipshttpd/router"
	"github.com/lupguo/ipsearch/version"
	"log"
	"net/http"
)

func Main() {
	go signalStop()
	router.Register()
	listen := config.Get().Listen
	log.Printf("ipshttpd listen on http://%s, ipshttd version %s", listen, version.VerClient)
	log.Fatalln(http.ListenAndServe(listen, nil))
}
