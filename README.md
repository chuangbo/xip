# xip

```sh
$ xip 1.2.3.4
Moscow 莫斯科
Moscow
Russia 俄罗斯联邦 (RU)
TimeZone: Europe/Moscow
$ # pipe from stdin
$ echo 1.1.1.1 | xip -
```

## Installation

`go get github.com/chuangbo/xip`

Download GeoLite2-City db and put it to location `/usr/local/etc/xip/GeoLite2-City/GeoLite2-City.mmdb`
