package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/mitchellh/go-homedir"
)

func main() {
	var dbFile string
	defaultDbFile, _ := homedir.Expand("~/.config/xip/qqwry.ipdb")

	flag.StringVar(&dbFile, "db", defaultDbFile, "IP库文件路径")
	cmdUpdate := flag.Bool("u", false, "更新纯真IP库")
	cmdVersion := flag.Bool("v", false, "Print the version number of xip")

	flag.Parse()

	if *cmdVersion {
		fmt.Printf("xip: %s\n", version)
		printDbInfo(dbFile)
		os.Exit(0)
	}

	if *cmdUpdate {
		if err := download(dbFile); err != nil {
			log.Fatal(err)
		}
		printDbInfo(dbFile)
		os.Exit(0)
	}

	isFromPipe := fromPipe()

	// 如果 db 文件不存在，先下载
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Printf("下载纯真IP库 \"%s\"...\n", dbFile)
		if err := download(dbFile); err != nil {
			log.Fatal(err)
		}
		printDbInfo(dbFile)
	}

	db, err := ipdb.NewCity(dbFile)
	if err != nil {
		fmt.Printf("IP库 \"%s\" 错误，可以使用 xip -u 命令重新下载\n", dbFile)
		log.Fatal(err)
	}

	if isFromPipe {
		pipeMode(db)
		return
	}

	if flag.NArg() > 0 {
		cliMode(db, flag.Args())
		return
	}

	flag.Usage()
	os.Exit(1)
}

func geoString(db *ipdb.City, ip string) string {
	r, err := db.FindInfo(ip, "CN")
	if err != nil {
		return color.RedString("%v", err)
	}

	return colorful(r)
}

func colorful(r *ipdb.CityInfo) string {
	if r == nil {
		return color.RedString("未知")
	}

	parts := []string{}

	if r.CountryName != "" {
		parts = append(parts, color.MagentaString(r.CountryName))
	}

	if r.RegionName != "" {
		parts = append(parts, color.CyanString(r.RegionName))
	}

	if r.CityName != "" {
		parts = append(parts, color.GreenString(r.CityName))
	}

	if r.OwnerDomain != "" {
		parts = append(parts, color.BlueString(r.OwnerDomain))
	}

	if r.IspDomain != "" {
		parts = append(parts, color.YellowString(r.IspDomain))
	}

	if len(parts) == 0 {
		return color.RedString("未知")
	}

	return strings.Join(parts, " ")
}

func printDbInfo(dbFile string) {
	db, err := ipdb.NewCity(dbFile)
	if err != nil {
		return
	}
	fmt.Println("\nIP 库信息:")
	fmt.Printf("db file: %s\n", dbFile)            // database file path
	fmt.Printf("ipv4 support: %v\n", db.IsIPv4())  // check database support ip type
	fmt.Printf("ipv6 support: %v\n", db.IsIPv6())  // check database support ip type
	fmt.Printf("build time: %v\n", db.BuildTime()) // database build time
	fmt.Printf("languages: %v\n", db.Languages())  // database support language
	fmt.Printf("fields: %v\n", db.Fields())        // database support fields
}
