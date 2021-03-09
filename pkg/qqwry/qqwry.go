// Package qqwry implements download and query IP geo-location infomation
// facilities for the famous qqwry.dat database.
//
// Inspired from github.com/tonywubo/qqwry, with bug fixes, unit tests
// and performance improvements.
package qqwry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

// DB is qqwry database instance.
type DB struct {
	buff []byte

	start, end, total uint32
}

// Record is query result.
type Record struct {
	Country, City string
}

// Open the qqwry.dat database.
func Open(file string) (*DB, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buff, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	start := binary.LittleEndian.Uint32(buff[:4])
	end := binary.LittleEndian.Uint32(buff[4:8])

	return &DB{
		buff: buff,

		start: start,
		end:   end,
		total: (end-start)/7 + 1,
	}, nil
}

func (record *Record) String() string {
	return fmt.Sprintf("%s %s", record.Country, record.City)
}

// Query ip geo location infomation from giving net.IP.
func (db *DB) Query(ip net.IP) (*Record, error) {
	if db.buff == nil {
		return nil, fmt.Errorf("db not initialized")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("not a valid ipv4 address")
	}

	offset := db.search(binary.BigEndian.Uint32(ipv4))
	if offset <= 0 {
		return &Record{}, nil
	}

	return db.readRecord(offset), nil
}

// Version return the version record for the database, from the last record for 255.255.255.255.
func (db *DB) Version() string {
	offset := getAddrFromRecord(db.buff[db.end : db.end+7])
	return db.readRecord(offset).City
}

// Total returns the total number of records for the database
func (db *DB) Total() uint32 {
	return db.total
}

func (db *DB) readRecord(offset uint32) *Record {
	rq := &Record{}

	offset += 4
	mode := db.buff[offset]

	// full redirect
	if mode == redirectMode1 {
		offset = db.readUint32FromByte3(offset + 1)
		mode = db.buff[offset]
	}

	// country
	var country []byte
	if mode == redirectMode2 {
		off1 := db.readUint32FromByte3(offset + 1)
		country = readCString(db.buff[off1:])
		offset += 4
	} else {
		country = readCString(db.buff[offset:])
		offset += uint32(len(country)) + 1
	}

	// area
	mode = db.buff[offset]
	if mode == redirectMode2 {
		offset = db.readUint32FromByte3(offset + 1)
	}
	area := readCString(db.buff[offset:])

	// decode gbk
	enc := simplifiedchinese.GBK.NewDecoder()
	if encoded, err := enc.Bytes(country); err == nil {
		rq.Country = string(encoded)
	}
	if encoded, err := enc.Bytes(area); err == nil {
		rq.City = string(encoded)
	}
	return rq
}

func (db *DB) readUint32FromByte3(offset uint32) uint32 {
	return byte3ToUInt32(db.buff[offset : offset+3])
}

func readCString(buf []byte) []byte {
	idx := bytes.IndexByte(buf, 0)
	if idx < 0 {
		// TODO: return error
		return nil
	}
	return buf[:idx]
}

func getIPFromRecord(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[:4])
}

func getAddrFromRecord(buf []byte) uint32 {
	return byte3ToUInt32(buf[4:7])
}

func (db *DB) search(ip uint32) uint32 {
	left := uint32(0)
	right := db.total

	for right-left > 1 {
		mid := (left + right) / 2
		offset := db.start + mid*7
		cur := getIPFromRecord(db.buff[offset : offset+7])

		if ip < cur {
			right = mid
		} else {
			left = mid
		}
	}

	offset := db.start + 7*left
	ipBegin := getIPFromRecord(db.buff[offset : offset+7])

	offset = getAddrFromRecord(db.buff[offset : offset+7])
	ipEnd := getIPFromRecord(db.buff[offset : offset+7])

	if ipBegin <= ip && ip <= ipEnd {
		return offset
	}

	return 0
}

func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
