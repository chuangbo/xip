// Package qqwry is a qqwry.dat reader, fork from https://github.com/tonywubo/qqwry
// TODO: refactor and release as a package
package qqwry

import (
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

// Reader is qqwry db reader
type Reader struct {
	buff []byte

	start, end, total uint32
}

// Record is query result
type Record struct {
	Country, City string
}

// Open qqwry db
func Open(file string) (*Reader, error) {
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

	return &Reader{
		buff: buff,

		start: start,
		end:   end,
		total: (end-start)/7 + 1,
	}, nil
}

func (record *Record) String() string {
	return fmt.Sprintf("%s %s", record.Country, record.City)
}

// Query ip geo info
func (r *Reader) Query(ip net.IP) (*Record, error) {
	if r.buff == nil {
		return nil, fmt.Errorf("db not initialized")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("not a valid ipv4 address")
	}

	offset := r.search(binary.BigEndian.Uint32(ipv4))
	if offset <= 0 {
		return &Record{}, nil
	}

	return r.readRecord(offset), nil
}

func (r *Reader) readRecord(offset uint32) *Record {
	rq := &Record{}

	offset += 4
	mode := r.buff[offset]

	// full redirect
	if mode == redirectMode1 {
		offset = r.readUint32FromByte3(offset + 1)
		mode = r.buff[offset]
	}

	// country
	var country []byte
	if mode == redirectMode2 {
		off1 := r.readUint32FromByte3(offset + 1)
		country = r.readString(off1)
		offset += 4
	} else {
		country = r.readString(offset)
		offset += uint32(len(country)) + 1
	}

	// area
	mode = r.buff[offset]
	if mode == redirectMode2 {
		offset = r.readUint32FromByte3(offset + 1)
	}
	area := r.readString(offset)

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

func (r *Reader) readUint32FromByte3(offset uint32) uint32 {
	return byte3ToUInt32(r.buff[offset : offset+3])
}

func (r *Reader) readString(offset uint32) []byte {
	for end := offset; ; end++ {
		if r.buff[end] == 0 {
			return r.buff[offset:end]
		}
	}
}

func getIPFromRecord(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[:4])
}

func getAddrFromRecord(buf []byte) uint32 {
	return byte3ToUInt32(buf[4:7])
}

func (r *Reader) search(ip uint32) uint32 {
	left := uint32(0)
	right := r.total

	for right-left > 1 {
		mid := (left + right) / 2
		offset := r.start + mid*7
		cur := getIPFromRecord(r.buff[offset : offset+7])

		if ip < cur {
			right = mid
		} else {
			left = mid
		}
	}

	offset := r.start + 7*left
	ipBegin := getIPFromRecord(r.buff[offset : offset+7])

	offset = getAddrFromRecord(r.buff[offset : offset+7])
	ipEnd := getIPFromRecord(r.buff[offset : offset+7])

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
