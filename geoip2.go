package main

import (
	"fmt"
	"log"
	"net"

	clr "github.com/logrusorgru/aurora"
	"github.com/oschwald/geoip2-golang"
)

func geoip2Output(db *geoip2.Reader, ip net.IP) {
	fmt.Print(ip)

	record, err := db.City(ip)
	if err != nil {
		log.Printf("error reading db: %v", clr.Red(err))
		return
	}

	if record.City.GeoNameID != 0 {
		fmt.Printf("\t%s %s", clr.Cyan(record.City.Names["en"]), clr.Cyan(record.City.Names["zh-CN"]))
	}

	for _, s := range record.Subdivisions {
		fmt.Printf("\t%s %s", clr.Green(s.Names["en"]), clr.Green(s.Names["zh-CN"]))
	}

	if record.Country.GeoNameID != 0 {
		fmt.Printf("\t%s %s %s", clr.Magenta(record.Country.Names["en"]), clr.Magenta(record.Country.Names["zh-CN"]), record.Country.IsoCode)
	}

	fmt.Println()
}
