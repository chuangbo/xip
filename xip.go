package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/chuangbo/xip/pkg/qqwry"
	clr "github.com/logrusorgru/aurora"
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
	// can be download from https://github.com/out0fmemory/qqwry.dat
	flag.StringVar(&dbFile, "db", defaultDbFile, "纯真IP库")

	flag.Parse()

	if flag.Arg(0) == "update" {
		if err := download(dbFile); err != nil {
			log.Fatal(clr.Red(err))
		}
		os.Exit(0)
	}

	db, err := qqwry.Open(dbFile)
	if err != nil {
		fmt.Printf("纯真IP库 \"%s\" 不存在，可以使用 update 命令下载\n", dbFile)
		log.Fatal(clr.Red(err))
	}

	if fromPipe() {
		pipeMode(db)
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cliMode(db, flag.Args())
}

func geoString(db *qqwry.Reader, ip net.IP) string {
	r, err := db.Query(ip)
	if err != nil {
		return clr.Red(err).String()
	}

	return fmt.Sprintf("%s %s", clr.Cyan(r.City), clr.Magenta(r.Country))
}
