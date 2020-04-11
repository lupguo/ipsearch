package ipshttpd

import (
	"log"
	"os"
	"os/signal"
)

// signalStop ctrl+c，停止ipshttpd服务
func signalStop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		log.Println("ipshttpd is shutdown")
	}
	os.Exit(1)
}
