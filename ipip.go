package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	clr "github.com/logrusorgru/aurora"
)

func ipipOutput(ip net.IP) {
	fmt.Print(ip)

	response, err := http.Get(fmt.Sprintf("http://freeapi.ipip.net/%s", ip))
	if err != nil {
		log.Printf("could not get response from ipip: %v", clr.Red(err))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("http error: %s", clr.Red(response.Status))
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("could not read response from ipip: %v", clr.Red(err))
		return
	}

	var r []string
	if err := json.Unmarshal(body, &r); err != nil {
		log.Printf("could not parse json %s: %v", body, clr.Red(err))
		return
	}

	if r[2] != "" {
		fmt.Printf("\t%s", clr.Cyan(r[2]))
	}

	if r[1] != "" {
		fmt.Printf("\t%s", clr.Green(r[1]))
	}

	if r[0] != "" {
		fmt.Printf("\t%s", clr.Magenta(r[0]))
	}

	if r[3] != "" {
		fmt.Printf("\t%s", clr.Yellow(r[3]))
	}

	if r[4] != "" {
		fmt.Printf("\t%s", clr.Blue(r[4]))
	}

	fmt.Println()
}
