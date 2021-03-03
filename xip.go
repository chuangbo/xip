package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/chuangbo/xip/pkg/qqwry"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

var (
	dbFile string

	db *qqwry.Reader
)

var (
	defaultDbFile, _ = homedir.Expand("~/.config/xip/qqwry.dat")
)

func main() {
	flag.StringVar(&dbFile, "db", defaultDbFile, "纯真IP库")
	v := flag.Bool("v", false, "Print the version number of xip")

	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.Arg(0) == "update" {
		if err := download(dbFile); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	isFromPipe := fromPipe()

	if flag.NArg() == 0 && !isFromPipe {
		flag.Usage()
		os.Exit(1)
	}

	db, err := qqwry.Open(dbFile)
	if err != nil {
		fmt.Printf("纯真IP库 \"%s\" 不存在，可以使用 xip update 命令下载\n", dbFile)
		log.Fatal(err)
	}

	if isFromPipe {
		pipeMode(db)
		return
	}

	cliMode(db, flag.Args())
}

func geoString(db *qqwry.Reader, ip net.IP) string {
	r, err := db.Query(ip)
	if err != nil {
		return color.RedString("%w", err)
	}

	return fmt.Sprintf("%s %s", color.CyanString(r.City), color.MagentaString(r.Country))
}
