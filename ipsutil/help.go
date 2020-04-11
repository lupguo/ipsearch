package ipsutil

import (
	"log"
)

//FatalOnError 快速错误响应
func FatalOnError(err error, msg string) {
	if err != nil {
		//panic(err)
		log.Fatalln(msg, err)
	}
}

