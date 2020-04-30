package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
)

var (
	// ref https://regexr.com/38odc
	regEx = regexp.MustCompile(`\b(?:(?:2(?:[0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9])\.){3}(?:(?:2([0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9]))\b`)
)

// ref: https://stackoverflow.com/questions/43947363/detect-if-a-command-is-piped-or-not
func fromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

// traceroute baidu.com | xip
func pipeMode() {
	db, err := qqwry.Open(qqwryDB)
	if err != nil {
		log.Fatal(clr.Red(err))
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		ipInfo := appendIPInfo(db, text)
		fmt.Println(text, ipInfo)
	}
}

func appendIPInfo(db *qqwry.Reader, text string) string {
	ip := regEx.FindString(text)
	if ip == "" {
		return ""
	}

	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return ""
	}

	r, err := db.Query(ipv4)

	if err != nil {
		return fmt.Sprintf("error reading db: %v", clr.Red(err))
	}

	return fmt.Sprintf("\t%s %s", clr.Cyan(r.City), clr.Magenta(r.Country))
}
