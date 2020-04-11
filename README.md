## ipsearch
> ip serach 查询是基于 http://ip.taobao.com/ipSearch.html 接口代理获取的数据（有请求频率限制）
>
> 以前使用ip.cn会出现频次限制，可以基于命令行走淘宝接口查询

### 安装
```
// 安装ipsearch命令工具，以及httpd服务
go get -u -v github.com/lupguo/ipsearch
```

### ipsearch 使用
```
// 命令行使用
Usage of ./ipsearch:
  -debug
    	debug for request response content
  -format string
    	response message format, default is json (json|text) (default "text")
  -ip string
    	the IP to be search, the default is the IP of the machine currently executing the cmdline
  -listen string
    	the listen address for ip search http server, eg 127.0.0.1:6100
  -proxy string
    	http proxy using for debugging, no proxy by default, eg http://127.0.0.1:8888
  -timeout duration
    	set http request timeout seconds (default 10s)
  -version
    	ipsearch version

// http服务
$ ./ipsearch -listen '0.0.0.0:6100'
2020/04/12 01:08:03 ipshttpd listen on http://0.0.0.0:6100, ipshttd version 0.4.0'

// 请求查询
$ curl localhost:6100
Version 0.4.0
Usage:
	//search current client ip information
	curl localhost:8680/ips

	//search for target ip information
	curl localhost:8680/ips?ip=targetIp

// 通过curl查询
$ curl localhost:8680/ips
{"addr":"中国 广东 深圳","network":"鹏博士","ip":"175.191.11.165"}
$ curl 'localhost:8680/ips?ip=175.190.11.16'
{"addr":"中国 辽宁 大连","network":"鹏博士","ip":"175.190.11.16"}
```

### 变更内容
- 2020-04-12：更新了目录结构，调整了重试机制
- 2019-08-07：代码目前版本还比较粗糙，会持续完善！
    - [ ] 代理问题，寻求更好用的代理
    - [x] 程序中一些已知的Bug修复 
- 2019-08-08
    - [x] 修复了客户端请求ipshttpd没有获取到正确IP的问题
    - [x] 修复了ipshttpd的handler处理
    - [x] 新增了版本展示
- 2019-08-13
    - [x] 新增Docker环境支持
