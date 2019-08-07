## ip-search
> ip serach 查询是基于 http://ip.taobao.com/ipSearch.html 接口代理获取的数据（有请求频率限制）
> 
> 以前使用ip.cn会出现频次限制，可以基于命令行走淘宝接口查询

### 安装
```
// 仅安装insearch命令
go get -v github.com/tkstorm/ip-search/...
go install -v github.com/tkstorm/ip-search/cmd/insearch

// 安装ipsearch命令工具，以及ipshttpd服务
go install -v github.com/tkstorm/ip-search/cmd/...
```

### ipsearch 使用
```
$ ipsearch -h
Usage of ipsearch:
  -debug
    	debug for request response content
  -ip string
    	ip to search, myip is current ip (default "myip")
  -mode string
    	response content mode (json|text) (default "text")
  -proxy string
    	request by proxy, using for debug
  -timeout duration
    	set http request timeout seconds (default 10s)

// 查看出口IP
$ ipsearch
Ip: 210.21.233.226, Network: 联通, Address: 中国 广东 深圳

// 查看指定IP，并以JSON格式输出
$ ipsearch -ip 118.144.149.206 -mode json
{"addr":"中国 北京 北京","network":"鹏博士","ip":"118.144.149.206"}
```

### ipshttpd

支持ipshttpd部署，相关请求会转发到ipshttpd查询服务器，然后将请求代理转发给淘宝查询IP信息。

```
// http服务
$ ipshttpd -listen 127.0.0.1:8087
2019/08/07 18:36:43 ip search httpd listen on 127.0.0.1:8087

// 通过curl查询
$ curl 'localhost:8680/ips?ip=117.90.252.208'
{"addr":"中国 江苏 镇江","network":"电信","ip":"117.90.252.208"}
```

### 注意
- 2019-08-07：代码目前版本还比较粗糙，会持续完善！
    - [ ] 代理问题，寻求更好用的代理
    - [ ] 程序中一些已知的Bug修复 

