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

	ip, err := getIP()
	if err != nil {
		log.Fatal(err)
	}

	db, err := geoip2.Open(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	output(record)
}

func getIP() (net.IP, error) {
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
		return nil, fmt.Errorf("invalid ip: %s", ipArg)
	}

	return ip, nil
}

func output(record *geoip2.City) {
	if record.City.GeoNameID != 0 {
		fmt.Printf("%s %s\n", record.City.Names["en"], record.City.Names["zh-CN"])
	}

	for _, s := range record.Subdivisions {
		fmt.Printf("%s %s\n", s.Names["en"], s.Names["zh-CN"])
	}

	fmt.Printf("%s %s (%s)\n", record.Country.Names["en"], record.Country.Names["zh-CN"], record.Country.IsoCode)
	fmt.Printf("TimeZone: %s\n", record.Location.TimeZone)
}
