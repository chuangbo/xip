package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/chuangbo/xip/pkg/qqwry"
	"github.com/fatih/color"
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
func pipeMode(db *qqwry.Reader) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(color.Output, text, "\t", findGeoStr(db, text))
	}
}

func findGeoStr(db *qqwry.Reader, text string) string {
	ip := regEx.FindString(text)
	if ip == "" {
		return ""
	}

	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return ""
	}

	return geoString(db, ipv4)
}
