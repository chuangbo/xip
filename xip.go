package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	dbFile := flag.String("db", "/usr/local/etc/xip/GeoLite2-City/GeoLite2-City.mmdb", "mmdb file")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	ips, err := getIPs()
	if err != nil {
		log.Fatal(err)
	}

	db, err := geoip2.Open(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, ip := range ips {
		output(db, ip)
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

func output(db *geoip2.Reader, ip net.IP) {
	fmt.Print(ip)

	record, err := db.City(ip)
	if err != nil {
		log.Printf("error reading db: %v", err)
		return
	}

	if record.City.GeoNameID != 0 {
		fmt.Printf("\t%s %s", record.City.Names["en"], record.City.Names["zh-CN"])
	}

	for _, s := range record.Subdivisions {
		fmt.Printf("\t%s %s", s.Names["en"], s.Names["zh-CN"])
	}

	if record.Country.GeoNameID != 0 {
		fmt.Printf("\t%s %s (%s)", record.Country.Names["en"], record.Country.Names["zh-CN"], record.Country.IsoCode)
	}

	if record.Location.TimeZone != "" {
		fmt.Printf("\tTimeZone: %s", record.Location.TimeZone)
	}

	fmt.Println()
}
