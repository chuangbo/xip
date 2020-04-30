# xip

Get ip geo information from multiple sources.

```sh
# quick start
$ xip 1.2.3.4
GeoIP2
1.2.3.4 Moscow 莫斯科   Moscow  Russia 俄罗斯联邦 RU
纯真IP
1.2.3.4 APNIC Debogon-prefix网络        澳大利亚

$ # pipe from stdin
$ dig baidu.com | xip

; <<>> DiG 9.10.6 <<>> baidu.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 16632
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 512
;; QUESTION SECTION:
;baidu.com.			IN	A

;; ANSWER SECTION:
baidu.com.		38	IN	A	220.181.38.148 	电信IDC机房 北京市
baidu.com.		38	IN	A	39.156.69.79 	移动 北京市

;; Query time: 30 msec
;; SERVER: 8.8.8.8#53(8.8.8.8) 	加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器 美国
;; WHEN: Thu Apr 30 17:46:12 CST 2020
;; MSG SIZE  rcvd: 70

$ # by host
$ xip google.com
GeoIP2
216.58.220.206  Mountain View 芒廷维尤  California 加利福尼亚州 United States 美国 US
2404:6800:4005:807::200e        Australia 澳大利亚 AU
纯真IP
216.58.220.206  Google全球边缘网络      美国

$ # enable ipip.net (limited free quote)
$ xip -ipip 1.2.3.4
GeoIP2
1.2.3.4	Moscow 莫斯科	Moscow 	Russia 俄罗斯联邦 RU
纯真IP
1.2.3.4 APNIC Debogon-prefix网络        澳大利亚
ipip.net
1.2.3.4	APNIC.NET
```

## Installation

`go get github.com/chuangbo/xip`

Download GeoLite2-City db and put it to location `/usr/local/etc/xip/GeoLite2-City/GeoLite2-City.mmdb`

Download QQWRY db and put it to location `/usr/local/etc/xip/qqwry.dat`

## Reference

* [GeoIP2](https://dev.maxmind.com/geoip/)
* [纯真IP](http://www.cz88.net/ip/)
    * [qqwry.dat](https://github.com/out0fmemory/qqwry.dat)
    * library fork from https://github.com/tonywubo/qqwry
* [ipip](https://www.ipip.net/)
