package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
	"github.com/oschwald/geoip2-golang"
)

var (
	enableGeoIP  bool
	enableQQWRY  bool
	enableIPIP   bool
	geoip2CityDB string
	qqwryDB      string
)

var (
	defaultGeoIP2CityDB, _ = homedir.Expand("~/.config/xip/GeoLite2-City/GeoLite2-City.mmdb")
	defaultQQWryDB, _      = homedir.Expand("~/.config/xip/qqwry.dat")
)

func main() {
	flag.BoolVar(&enableGeoIP, "geoip2", true, "enable geoip2")
	flag.BoolVar(&enableQQWRY, "qqwry", true, "enable 纯真IP数据库")
	flag.BoolVar(&enableIPIP, "ipip", false, "enable ipip.net")

	flag.StringVar(&geoip2CityDB, "geoip2-city-db", defaultGeoIP2CityDB, "mmdb file")

	// can be download from https://github.com/out0fmemory/qqwry.dat
	flag.StringVar(&qqwryDB, "qqwry-db", defaultQQWryDB, "纯真IP数据库")

	flag.Parse()

	if fromPipe() {
		pipeMode()
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var ips []net.IP
	for _, arg := range flag.Args() {
		result, err := getIPs(arg)
		if err != nil {
			log.Fatal(clr.Red(err))
		}
		ips = append(ips, result...)
	}

	if enableGeoIP {
		db, err := geoip2.Open(geoip2CityDB)
		if err != nil {
			log.Fatal(clr.Red(err))
		}
		defer db.Close()

		fmt.Println("GeoIP2")
		for _, ip := range ips {
			geoip2Output(db, ip)
		}
	}

	if enableQQWRY {
		db, err := qqwry.Open(qqwryDB)
		if err != nil {
			log.Fatal(clr.Red(err))
		}

		fmt.Println("纯真IP")
		for _, ip := range ips {
			qqwryOutput(db, ip)
		}
	}

	if enableIPIP {
		fmt.Println("ipip.net")
		for _, ip := range ips {
			ipipOutput(ip)
		}
	}
}

func getIPs(arg string) ([]net.IP, error) {
	ip := net.ParseIP(arg)
	if ip == nil {
		// host
		return net.LookupIP(arg)
	}
	return []net.IP{ip}, nil
}
