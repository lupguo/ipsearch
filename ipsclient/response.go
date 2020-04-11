package ipsclient

import (
	"encoding/json"
	"fmt"
)

// Result IP查询结果
// {"addr": "中国 广东 深圳 (南山区)", "network": "联通", "ip": "210.21.233.226" }
type Result struct {
	Addr    string `json:"addr"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
}

// Render IpsResult消息展示模式
func (r *Result) Render(mode string) (msg string, err error) {
	switch mode {
	case "json":
		rt, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(rt), nil
	default:
		return fmt.Sprintf("Ip: %s, Network: %s, Address: %s", r.Ip, r.Network, r.Addr), nil
	}
}

// Taobao 为从最上游获取到的IP请求结果
// {
// 		"code": 0,
// 		"data": {
// 			"ip": "210.21.233.226",
// 			"country": "中国",
// 			"area": "",
// 			"region": "广东",
// 			"city": "深圳",
// 			"county": "XX",
// 			"isp": "联通",
// 			"country_id": "CN",
// 			"area_id": "",
// 			"region_id": "440000",
// 			"city_id": "440300",
// 			"county_id": "xx",
// 			"isp_id": "100026"
// 		}
// }
//
type Taobao struct {
	Code int               `json:"code"`
	Data map[string]string `json:"data"`
}

//toResult 由淘宝的Json转换为IpsResult
func (tb *Taobao) toResult() *Result {
	return &Result{
		Addr:    fmt.Sprintf("%s %s %s", tb.Data["country"], tb.Data["region"], tb.Data["city"]),
		Network: tb.Data["isp"],
		Ip:      tb.Data["ip"],
	}
}

