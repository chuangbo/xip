# xip

Get ip geo information from GeoIP2 db and ipip.net free api.

```sh
# quick start
$ xip 1.2.3.4
GeoIP2
1.2.3.4 Moscow 莫斯科   Moscow  Russia 俄罗斯联邦 RU
纯真IP
1.2.3.4 APNIC Debogon-prefix网络        澳大利亚

$ # pipe from stdin
$ echo 1.1.1.1 | xip -
GeoIP2
1.1.1.1	Australia 澳大利亚 AU
纯真IP
1.1.1.1 APNIC&CloudFlare公共DNS服务器   美国

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
