package handler

import (
	"github.com/tkstorm/ip-search/version"
	"log"
	"net/http"
	"text/template"
)

// 帮助信息
var usageFormat = `Version {{.}}
Usage:
	//search current client ip information
	curl localhost:8680/ips

	//search for target ip information  
	curl localhost:8680/ips?ip=targetIp	
`

var tpl = template.Must(template.New("verhttpd").Parse(usageFormat))

//HelpMessage 显示ishttpd的帮助信息
func HelpMessage(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Execute(w, version.VerHttpd); err != nil {
		log.Println(err)
	}
}
