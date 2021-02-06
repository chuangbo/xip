package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
)

// xip 1.1.1.1 baidu.com
func cliMode(db *qqwry.Reader, args []string) {
	var ips []net.IP
	for _, arg := range flag.Args() {
		result, err := argToIPs(arg)
		if err != nil {
			log.Fatal(clr.Red(err))
		}
		ips = append(ips, result...)
	}

	for _, ip := range ips {
		result := geoString(db, ip)
		fmt.Printf("%s\t%s\n", ip, result)
	}
}

func argToIPs(arg string) ([]net.IP, error) {
	ip := net.ParseIP(arg)
	if ip == nil {
		// host
		return net.LookupIP(arg)
	}
	return []net.IP{ip}, nil
}