# xip

Get ip geo information from GeoIP2 db and ipip.net free api.

```sh
$ xip 1.2.3.4
GeoIP2
1.2.3.4	Moscow 莫斯科	Moscow 	Russia 俄罗斯联邦 RU
ipip.net
1.2.3.4	APNIC.NET

$ # pipe from stdin
$ echo 1.1.1.1 | xip -
GeoIP2
1.1.1.1	Australia 澳大利亚 AU
ipip.net
1.1.1.1	CLOUDFLARE.COM

$ # host
$ xip google.com
GeoIP2
216.58.200.78	Mountain View 芒廷维尤	California 加利福尼亚州	United States 美国 US
2404:6800:4005:80f::200e	Australia 澳大利亚 AU
ipip.net
216.58.200.78	香港	中国
2404:6800:4005:809::200e	香港	中国
```

## Installation

`go get github.com/chuangbo/xip`

Download GeoLite2-City db and put it to location `/usr/local/etc/xip/GeoLite2-City/GeoLite2-City.mmdb`
