package main

import (
	"fmt"
	"log"
	"net"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
)

func qqwryOutput(db *qqwry.Reader, ip net.IP) {
	fmt.Print(ip)

	r, err := db.Query(ip)
	if err != nil {
		log.Printf("error reading db: %v", clr.Red(err))
		return
	}

	if r.City != "" {
		fmt.Printf("\t%s", clr.Cyan(r.City))
	}

	if r.Country != "" {
		fmt.Printf("\t%s", clr.Magenta(r.Country))
	}

	fmt.Println()
}
