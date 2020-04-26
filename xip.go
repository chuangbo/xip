package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	enableGeoIP := flag.Bool("geoip2", true, "enable geoip2")
	enableQQWRY := flag.Bool("qqwry", true, "enable 纯真IP数据库")
	enableIPIP := flag.Bool("ipip", false, "enable ipip.net")

	geoip2CityDB := flag.String("geoip2-city-db", "/usr/local/etc/xip/GeoLite2-City/GeoLite2-City.mmdb", "mmdb file")

	// can be download from https://github.com/out0fmemory/qqwry.dat
	qqwryDB := flag.String("qqwry-db", "/usr/local/etc/xip/qqwry.dat", "纯真IP数据库")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	ips, err := getIPs()
	if err != nil {
		log.Fatal(clr.Red(err))
	}

	if *enableGeoIP {
		db, err := geoip2.Open(*geoip2CityDB)
		if err != nil {
			log.Fatal(clr.Red(err))
		}
		defer db.Close()

		fmt.Println("GeoIP2")
		for _, ip := range ips {
			geoip2Output(db, ip)
		}
	}

	if *enableQQWRY {
		db, err := qqwry.Open(*qqwryDB)
		if err != nil {
			log.Fatal(clr.Red(err))
		}

		fmt.Println("纯真IP")
		for _, ip := range ips {
			qqwryOutput(db, ip)
		}
	}

	if *enableIPIP {
		fmt.Println("ipip.net")
		for _, ip := range ips {
			ipipOutput(ip)
		}
	}
}

func getIPs() ([]net.IP, error) {
	ipArg := flag.Arg(0)

	// read from stdin if ip arg is `-`
	if ipArg == "-" {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		ipArg = strings.Trim(text, " \n")
	}

	ip := net.ParseIP(ipArg)
	if ip == nil {
		return net.LookupIP(ipArg)
	}

	return []net.IP{ip}, nil
}
