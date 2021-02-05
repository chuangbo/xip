# xip

Get ip geo information from multiple sources.

```sh
# quick start
$ xip 1.2.3.4
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
216.58.220.206  Google全球边缘网络      美国

```

## Installation

`go install github.com/chuangbo/xip`

Download QQWRY db (defaults to `~/.config/xip/qqwry.dat`)

```sh
$ xip update
Downloading to "/Users/me/.config/xip/qqwry.dat"
332.01 KiB / 5.13 MiB [==>---------------------------------------|  59s ] 83.00 KiB/s
```

## Reference

* [纯真IP](http://www.cz88.net/ip/)
* [qqwry.dat](https://github.com/out0fmemory/qqwry.dat)
* library fork from https://github.com/tonywubo/qqwry
