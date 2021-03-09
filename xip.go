package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/chuangbo/xip/v2/pkg/qqwry"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

var (
	dbFile string

	db *qqwry.DB
)

var (
	defaultDbFile, _ = homedir.Expand("~/.config/xip/qqwry.dat")
)

func main() {
	flag.StringVar(&dbFile, "db", defaultDbFile, "纯真IP库文件路径")
	cmdUpdate := flag.Bool("u", false, "更新纯真IP库")
	cmdVersion := flag.Bool("v", false, "Print the version number of xip")
	cmdDump := flag.Bool("dump", false, "Dump all the qqwry records")

	flag.Parse()

	if *cmdVersion {
		fmt.Printf("xip: %s\n", version)
		if db, err := qqwry.Open(dbFile); err == nil {
			fmt.Printf("qqwry: %s\n", db.Version())
		}
		os.Exit(0)
	}

	if *cmdUpdate {
		if err := download(dbFile); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	isFromPipe := fromPipe()

	db, err := qqwry.Open(dbFile)
	if err != nil {
		fmt.Printf("纯真IP库 \"%s\" 不存在，可以使用 xip -u 命令下载\n", dbFile)
		log.Fatal(err)
	}

	if *cmdDump {
		db.Dump()
		return
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

func geoString(db *qqwry.DB, ip net.IP) string {
	r, err := db.Query(ip)
	if err != nil {
		return color.RedString("%w", err)
	}

	return colorful(r)
}

func colorful(r *qqwry.Record) string {
	return fmt.Sprintf("%s %s", color.CyanString(r.City), color.MagentaString(r.Country))
}
