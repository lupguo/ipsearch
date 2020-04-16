package version

import (
	"fmt"
	"os"
)

var (
	// ipsearch 版本展示
	VerClient = "0.5.0"
)

//ShowVersion 显示版本信息
func ShowVersion(ver bool) {
	if ver == true {
		fmt.Println("ipsearch version", VerClient)
		os.Exit(0)
	}
}
