package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/ipipdotnet/ipdb-go"
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
func pipeMode(db *ipdb.City) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(color.Output, text, "\t", findGeoStr(db, text))
	}
}

func findGeoStr(db *ipdb.City, text string) string {
	ip := regEx.FindString(text)
	if ip == "" {
		return ""
	}

	return geoString(db, ip)
}
