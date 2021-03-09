package qqwry_test

import (
	"fmt"
	"log"
	"net"

	"github.com/chuangbo/xip/pkg/qqwry"
)

func Example_query() {
	db, err := qqwry.Open("testdata/qqwry.dat")
	if err != nil {
		log.Fatal(err)
	}

	ip := net.ParseIP("192.168.1.1")

	result, err := db.Query(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Country: %s, City: %s\n", result.Country, result.City)
	// Output: Country: 局域网, City: 对方和您在同一内部网
}

func Example_download() {
	// Check remote version, and get the key to decrypt
	key, remoteVersion, err := qqwry.GetUpdateInfo()

	// Check local database verion
	db, err := qqwry.Open("testdata/qqwry.dat")
	if err == nil {
		localVersion := db.Version()
		// Skip update if local version is same as remote version
		if qqwry.SameVersion(remoteVersion, localVersion) {
			log.Fatal("No need to download")
		}
	}

	// Open the download url
	contentLength, reader, err := qqwry.Download(key)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	fmt.Printf("ContentLength: %d\n", contentLength)

	// Save to local disk
	// io.Copy(f, reader)
}
