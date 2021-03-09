# xip

[![Go Reference](https://pkg.go.dev/badge/github.com/chuangbo/xip/v2/pkg/qqwry.svg)](https://pkg.go.dev/github.com/chuangbo/xip/v2/pkg/qqwry)
[![Go Report Card](https://goreportcard.com/badge/github.com/chuangbo/xip/v2)](https://goreportcard.com/report/github.com/chuangbo/xip/v2)

xip = 查 IP, Get ip geo information.

## Installation

Download from [GitHub Releases](https://github.com/chuangbo/xip/releases) and place the executable file in your PATH.

Or the go way:

```sh
$ go install github.com/chuangbo/xip/v2
```

Download QQWRY db (defaults to `~/.config/xip/qqwry.dat`)

```sh
$ xip -u
Downloading to "/Users/me/.config/xip/qqwry.dat"
332.01 KiB / 5.13 MiB [==>---------------------------------------|  59s ] 83.00 KiB/s
```

## Usage

Query a single IP address:

```sh
$ xip 1.2.3.4
1.2.3.4 APNIC Debogon-prefix网络        澳大利亚
```

Query multiple IP addresses:
```sh
$ xip 1.2.3.4 8.8.8.8
1.2.3.4 APNIC Debogon-prefix网络 澳大利亚
8.8.8.8 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器 美国
```

Query by domain name:
```sh
$ xip baidu.com
39.156.69.79    移动    北京市
220.181.38.148  电信IDC机房     北京市
```

Pipe from stdin (append first found IP info the the end of each line):

```sh
$ dig baidu.com +short | xip
220.181.38.148   电信IDC机房 北京市
39.156.69.79     移动 北京市

$ traceroute -n baidu.com | xip
traceroute to baidu.com (220.181.38.148), 30 hops max, 60 byte packets  电信IDC机房 北京市
...
 6  203.100.48.189  2.236 ms 210.48.136.205  0.857 ms 203.100.48.189  1.964 ms  中国电信CN2节
点 香港
 7  59.43.247.229  22.874 ms 59.43.181.185  1.470 ms 59.43.186.121  1.670 ms    中国电信CN2骨
干网 中国
 8  59.43.246.209  38.864 ms 59.43.246.225  39.442 ms 59.43.250.173  39.976 ms  中国电信CN2骨
干网 中国
 9  59.43.188.69  43.039 ms 59.43.188.65  41.212 ms 59.43.188.73  42.410 ms     中国电信CN2骨
干网 中国
10  59.43.132.17  51.926 ms 59.43.132.13  44.629 ms  44.821 ms  中国电信CN2骨干网 中国
11  202.97.12.177  83.995 ms 202.97.42.13  41.807 ms 202.97.12.177  80.883 ms   电信骨干网 中
国
12  36.110.247.46  40.963 ms 36.110.246.154  41.153 ms 180.149.128.118  41.589 ms       电信
北京市
13  36.110.246.65  42.199 ms * 36.110.249.58  43.591 ms         电信 北京市
14  * * 220.181.182.30  48.218 ms       电信互联网数据中心 北京市
15  * * 220.181.182.30  44.302 ms       电信互联网数据中心 北京市

```

## Documentation

Golang package `github.com/chuangbo/xip/v2/pkg/qqwry` implements download and query IP geo-location information
facilities for the famous qqwry.dat database.

Inspired from github.com/tonywubo/qqwry, with bug fixes, unit tests and performance improvements.

* [Examples](https://pkg.go.dev/github.com/chuangbo/xip/v2/pkg/qqwry#pkg-examples)
* [Package Reference](https://pkg.go.dev/github.com/chuangbo/xip/v2/pkg/qqwry)

## Reference

* [纯真IP](http://www.cz88.net/ip/) 纯真网络提供的免费离线 IP 数据库
* [tonywubo/qqwry](https://github.com/tonywubo/qqwry) 查询纯真数据库代码在此基础，修改了一些 bug，并加上了单元测试和性能优化
* [freshcn/qqwry](https://github.com/freshcn/qqwry/blob/master/download.go) 下载数据库代码在此基础修改而来
* [Nali CLI](https://github.com/SukkaW/nali-cli) nodejs 编写的功能更多的工具，xip 可以认为是 Nali 的 Golang 替代品

## Author

Chuangbo Li
