package ipserach

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// getClient 获取http客户端
// todo 这里有资源重复创建的问题，需要通过资源池修复
func getClient(proxy string, timeout time.Duration) (client *http.Client, err error) {
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				if strings.Trim(proxy, " ") == "" {
					return nil, err
				}
				return url.Parse(proxy)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}
	return client, nil
}

// Dispatch 针对资源进行初始化、回收刷新、调度
type Dispatch struct {
	ch ResCh
}
