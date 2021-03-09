package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/chuangbo/xip/v2/pkg/qqwry"
	"github.com/fatih/color"
)

// xip 1.1.1.1 baidu.com
func cliMode(db *qqwry.DB, args []string) {
	var ips []net.IP
	for _, arg := range flag.Args() {
		result, err := argToIPs(arg)
		if err != nil {
			log.Fatal(err)
		}
		ips = append(ips, result...)
	}

	for _, ip := range ips {
		result := geoString(db, ip)
		fmt.Fprintf(color.Output, "%s\t%s\n", ip, result)
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
