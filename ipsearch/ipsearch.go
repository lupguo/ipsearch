package ipsearch

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ipsearch 版本控制
const version = "beta 0.2.0"

// version 获取版本信息
func Version() string {
	return fmt.Sprintf("ipsearch version %s", version)
}

// Source 为代理源
type Source struct {
	Name string
}

// Ips 为Ip Search搜索对象
type Ips struct {
	Debug   bool
	Proxy   string
	Timeout time.Duration
	once    sync.Once
	client  *http.Client
	source  *Source
}

//NewIps 初始化Ips客户端
func NewIps(debug bool, proxy string, timeout time.Duration) *Ips {
	ips := &Ips{Debug: debug, Proxy: proxy, Timeout: timeout}
	ips.makeClient()
	return ips
}

// IpsResult IP查询结果
// {"addr": "中国 广东 深圳 (南山区)", "network": "联通", "ip": "210.21.233.226" }
type IpsResult struct {
	Addr    string `json:"addr"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
}

// IpsOrigin 为从最上游获取到的IP请求结果
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
type IpsOrigin struct {
	Code int               `json:"code"`
	Data map[string]string `json:"data"`
}

// Search 基于Tick退避做IP搜索，返回IpResult指针，或者任何错误信息
// 原始请求头
// POST /service/getIpInfo2.php HTTP/1.1
//		Host: ip.taobao.com
//		Connection: keep-alive
//		Content-Length: 7
//		Pragma: no-cache
//		Cache-Control: no-cache
//		Accept: */*
//		Origin: http://ip.taobao.com
//		X-Requested-With: XMLHttpRequest
//		User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36
//		Content-Type: application/x-www-form-urlencoded
//		Referer: http://ip.taobao.com/ipSearch.html
//		Accept-Encoding: gzip, deflate
//		Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
//		Cookie: tracknick=%5Cu968F%5Cu60F3%5Cu98CE%5Cu66B4; t=80d8c3e412c80d30bf5fb21f598aca51; tg=0; cna=rENrE7ooBC8CAbfr/zaRO+bd; thw=cn; miid=262416151115249426; lgc=%5Cu968F%5Cu60F3%5Cu98CE%5Cu66B4; cookie2=10f99d7c9d9065643495fffbb2c8be45; v=0; _tb_token_=33857584b3e35; hng=CN%7Czh-CN%7CCNY%7C156; dnk=%5Cu968F%5Cu60F3%5Cu98CE%5Cu66B4; ali_ab=101.232.210.204.1563463590581.1; mt=ci=-1_0; SL_GWPT_Show_Hide_tmp=1; SL_wptGlobTipTmp=1; unb=2318891944; uc1=lng=zh_CN&tag=8&existShop=false&pas=0&cookie15=WqG3DMC9VAQiUQ%3D%3D&cookie21=UtASsssmeW6lpyd%2BB%2B3t&cookie16=VFC%2FuZ9az08KUQ56dCrZDlbNdA%3D%3D&cookie14=UoTaHY9FEv88hQ%3D%3D; uc3=vt3=F8dBy32hQwiVx9nWPms%3D&nk2=qE5eAz61h%2Bs%3D&id2=UUtKd1ZGvQ3HcA%3D%3D&lg2=W5iHLLyFOGW7aA%3D%3D; csg=c847892b; cookie17=UUtKd1ZGvQ3HcA%3D%3D; skt=0e8ccc07bb02f6ab; existShop=MTU2NTE0MDM0MQ%3D%3D; uc4=nk4=0%40qnXzu6NMttWwzp6U4KR8%2FE3kUA%3D%3D&id4=0%40U2lwJloSs8mO3WjIF9wnI%2BB%2B3gjB; _cc_=U%2BGCWk%2F7og%3D%3D; _l_g_=Ug%3D%3D; sg=%E6%9A%B449; _nk_=%5Cu968F%5Cu60F3%5Cu98CE%5Cu66B4; cookie1=VASso2Si03x8IXuaQzjz9fZXNH9k1jvHQlZNrwqUQCk%3D; l=cBI-nTnRvxd-3fT0BOfGqZ6T8v7T0Idf1sPPhXGi7ICPOJ5BqDYdWZFZgIL6CnGVLsM9-3okgj63BzLiGyUiQGmCqVMDkL7R.
func (ips *Ips) Search(ip string) (rs *IpsResult, err error) {
	// create request
	if strings.Trim(ip, " ") == "" {
		ip = "myip"
	}
	body := strings.NewReader("ip=" + ip)
	dstUrl := "http://ip.taobao.com/service/getIpInfo2.php"
	req, err := http.NewRequest("POST", dstUrl, body)
	if err != nil {
		return nil, err
	}
	// head agent setting
	head := map[string]string{
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36",
		"Content-Type":     "application/x-www-form-urlencoded",
		"Origin":           "http://ip.taobao.com",
		"Cache-Control":    "no-cache",
		"Referer":          "http://ip.taobao.com/ipSearch.html",
		"Connection":       "keep-alive",
		"Accept":           "*/*",
		"X-Requested-With": "XMLHttpRequest",
	}
	for k, v := range head {
		req.Header.Set(k, v)
	}

	// http 请求+重试
	over := time.Tick(5 * time.Second)
	var resp *http.Response
loop:
	for {
		select {
		case <-over:
			return nil, errors.New("request tick timeout")
		default:
			resp, err = ips.client.Do(req)
			if err != nil {
				log.Println(err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			break loop
		}
	}
	defer resp.Body.Close()

	// http status check
	if code := resp.StatusCode; code != http.StatusOK {
		return nil, errors.New("get ip failed, http response status code " + strconv.Itoa(code))
	}

	// read body content
	orignBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse content
	v := IpsOrigin{
		Data: make(map[string]string),
	}
	if err := json.Unmarshal(orignBytes, &v); err != nil {
		return nil, err
	}
	if v.Code != 0 {
		return nil, errors.New("search origin code != 0")
	}

	// ips result
	if ips.Debug {
		log.Println(v.Data)
	}
	d := v.Data
	ipsRs := &IpsResult{
		Addr:    fmt.Sprintf("%s %s %s", d["country"], d["region"], d["city"]),
		Network: d["isp"],
		Ip:      d["ip"],
	}
	return ipsRs, err
}

// makeClient 创建
func (ips *Ips) makeClient() {
	if ips.client == nil {
		ips.once.Do(func() {
			ips.client = &http.Client{
				Transport: &http.Transport{
					Proxy: func(request *http.Request) (*url.URL, error) {
						if ips.Proxy == "" {
							return nil, nil
						}
						return url.Parse(ips.Proxy)
					},
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
				Timeout: ips.Timeout,
			}
		})
	}
}

// SetProxy 针对IpSearch设置代理，主要用于本地调试
func (ips *Ips) SetProxy(proxy string) *Ips {
	ips.Proxy = proxy
	return ips
}

// SetDebug 开启调试信息
func (ips *Ips) SetDebug(debug bool) *Ips {
	ips.Debug = debug
	return ips
}

// Message 消息展示模式
func (ipsRs *IpsResult) Message(mode string) (msg string, err error) {
	switch mode {
	case "json":
		rt, err := json.Marshal(ipsRs)
		if err != nil {
			return "", err
		}
		return string(rt), nil
	default:
		return fmt.Sprintf("Ip: %s, Network: %s, Address: %s", ipsRs.Ip, ipsRs.Network, ipsRs.Addr), nil
	}
}
