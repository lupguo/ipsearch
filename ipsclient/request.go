package ipsclient

import (
	"crypto/tls"
	"encoding/json"
	errors "errors"
	"github.com/lupguo/ipsearch/ipsutil"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var taobaoUrl = "http://ip.taobao.com/service/getIpInfo.php"

// Ips 为Ip Search搜索对象
type Ips struct {
	Debug    bool
	Client   *http.Client
	Request  *http.Request
	Response *http.Response
}

//NewIps 初始化Ips客户端
func NewIps(debug bool, proxy string, timeout time.Duration) *Ips {
	ips := &Ips{
		Debug:   debug,
		Client:  makeClient(proxy, timeout),
		Request: makeRequest(),
	}
	return ips
}

// makeClient 生成一个http请求客户端
func makeClient(proxy string, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				if proxy == "" {
					return nil, nil
				}
				return url.Parse(proxy)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}
}

// makeRequest 初始化Ipsearch的请求头
func makeRequest() *http.Request {
	req, err := http.NewRequest("GET", taobaoUrl, nil)
	ipsutil.FatalOnError(err, "new ipsearch init request failed.")

	// head agent setting
	head := map[string]string{
		"User-Agent":       userAgent(),
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
	return req
}

var agents = []string{
	`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/534.55.3 (KHTML, like Gecko) Version/5.1.3 Safari/534.53.10`,
	`Mozilla/5.0 (iPad; CPU OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko ) Version/5.1 Mobile/9B176 Safari/7534.48.3`,
	`Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; de-at) AppleWebKit/533.21.1 (KHTML, like Gecko) Version/5.0.5 Safari/533.21.1`,
	`Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_7; da-dk) AppleWebKit/533.21.1 (KHTML, like Gecko) Version/5.0.5 Safari/533.21.1`,
	`Mozilla/5.0 (Windows; U; Windows NT 6.1; tr-TR) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27`,
}

// userAgent 返回随机Agent
func userAgent() string {
	return agents[rand.Intn(len(agents))]
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
func (ips *Ips) Search(ip string) (r *Result, err error) {
	ips.updateURL(ip)
	return ips.doHttpRequest()
}

// updateURL 更新请求的URL信息
func (ips *Ips) updateURL(ip string) {
	u := ips.Request.URL
	q := u.Query()
	if ip == "" {
		ip = "myip"
	}
	q.Add("ip", ip)
	u.RawQuery = q.Encode()
	ips.Request.URL = u
}

var (
	errRetryTimeout     = errors.New("HTTP request timed out and exceeded the maximum retry time")
	errRequestFailed    = errors.New("HTTP request failed, keep trying")
	errHttpStatusCode   = errors.New("HTTP request response is not 200 status code, keep trying")
	errEmptyContent     = errors.New("HTTP request response is empty, keep trying")
	errTaobaoStatusCode = errors.New("the json data status code returned from Taobao is non-zero")
)

// httpResult 执行http request通道返回的内容
type httpResult struct {
	r *Result
	e error
}

// doHttpRequest 执行HTTP请求，设定3秒超时，如果查询失败，尝试重试
func (ips *Ips) doHttpRequest() (r *Result, err error) {
	ch := make(chan httpResult)
	go func() {
		over := false
		after := time.After(30 * time.Second)
		for !over {
			select {
			case <-after:
				over = true
				ch <- httpResult{nil, errRetryTimeout}
			default:
				resp, err := ips.Client.Do(ips.Request)
				if err != nil || resp.StatusCode != http.StatusOK || resp.ContentLength == 0 {
					switch {
					case err != nil:
						err = errRequestFailed
					case resp.StatusCode != http.StatusOK:
						err = errHttpStatusCode
					case resp.ContentLength == 0:
						err = errEmptyContent
					}
					if ips.Debug {
						log.Println(resp, err)
					}
					time.Sleep(1000 * time.Millisecond)
					continue
				}
				r, err := ips.parseBody(resp)
				ch <- httpResult{r, err}
				over = true
			}
		}
		close(ch)
	}()
	chr := <-ch
	return chr.r, chr.e
}

// parseBody 解析淘宝的IP响应结果
func (ips *Ips) parseBody(resp *http.Response) (*Result, error) {
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	taobao := Taobao{
		Data: make(map[string]string),
	}
	if err := json.Unmarshal(raw, &taobao); err != nil {
		return nil, err
	}
	if taobao.Code != 0 {
		return nil, errTaobaoStatusCode
	}
	return taobao.toResult(), nil
}
