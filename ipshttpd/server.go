package ipshttpd

import (
	"github.com/lupguo/ipsearch/config"
	"github.com/lupguo/ipsearch/ipshttpd/router"
	"github.com/lupguo/ipsearch/version"
	"log"
	"net/http"
)

var listen = config.Get().Listen

func Main() {
	go signalStop()
	router.Register()
	log.Printf("ipshttpd listen on http://%s, ipshttd version %s", listen, version.VerHttpd)
	log.Fatalln(http.ListenAndServe(listen, nil))
}
