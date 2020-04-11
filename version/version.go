package version

import (
	"fmt"
	"os"
)

var (
	// ipsearch 版本控制
	VerClient = "0.4.0"
	VerHttpd  = "0.4.0"
)

//ShowVersion 显示版本信息
func ShowVersion(ver bool) {
	if ver == true {
		fmt.Println("ipsearch version", VerClient)
		fmt.Println("ipshttpd version", VerHttpd)
		os.Exit(0)
	}
}
