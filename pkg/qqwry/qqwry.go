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
	indexLength   = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

// Reader is qqwry db reader
type Reader struct {
	buff  []byte
	start uint32
	end   uint32
}

// Record is query result
type Record struct {
	Country string
	City    string
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

	return &Reader{
		buff: buff,

		start: binary.LittleEndian.Uint32(buff[:4]),
		end:   binary.LittleEndian.Uint32(buff[4:8]),
	}, nil
}

func (record *Record) String() string {
	return fmt.Sprintf("%s %s", record.Country, record.City)
}

// Query ip geo info
func (r *Reader) Query(ip net.IP) (*Record, error) {
	rq := &Record{}

	if r.buff == nil {
		return nil, fmt.Errorf("db not initialized")
	}

	var country []byte
	var area []byte

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("not a valid ipv4 address")
	}

	offset := r.binSearch(binary.BigEndian.Uint32(ipv4))
	if offset <= 0 {
		return rq, nil
	}

	mode := r.readMode(offset + 4)
	if mode == redirectMode1 {
		countryOffset := r.readUint32FromByte3(offset + 5)

		mode = r.readMode(countryOffset)
		if mode == redirectMode2 {
			c := r.readUint32FromByte3(countryOffset + 1)
			country = r.readString(c)
			countryOffset += 4
			area = r.readArea(countryOffset)

		} else {
			country = r.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
			area = r.readArea(countryOffset)
		}

	} else if mode == redirectMode2 {
		countryOffset := r.readUint32FromByte3(offset + 5)
		country = r.readString(countryOffset)
		area = r.readArea(offset + 8)
	}

	// decode gbk
	enc := simplifiedchinese.GBK.NewDecoder()
	if encoded, err := enc.Bytes(country); err == nil {
		rq.Country = string(encoded)
	}
	if encoded, err := enc.Bytes(area); err == nil {
		rq.City = string(encoded)
	}
	return rq, nil
}

func (r *Reader) readUint32FromByte3(offset uint32) uint32 {
	return byte3ToUInt32(r.buff[offset : offset+3])
}

func (r *Reader) readMode(offset uint32) byte {
	return r.buff[offset : offset+1][0]
}

func (r *Reader) readString(offset uint32) []byte {
	for end := offset; ; end++ {
		if r.buff[end] == 0 {
			return r.buff[offset:end]
		}
	}
}

func (r *Reader) readArea(offset uint32) []byte {
	mode := r.readMode(offset)
	if mode == redirectMode1 || mode == redirectMode2 {
		areaOffset := r.readUint32FromByte3(offset + 1)
		if areaOffset == 0 {
			return []byte{}
		}
		return r.readString(areaOffset)
	}
	return r.readString(offset)
}

func (r *Reader) getRecord(offset uint32) []byte {
	return r.buff[offset : offset+indexLength]
}

func getIPFromRecord(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[:4])
}

func getAddrFromRecord(buf []byte) uint32 {
	return byte3ToUInt32(buf[4:7])
}

func (r *Reader) binSearch(ip uint32) uint32 {
	start := r.start
	end := r.end

	// log.Printf("len info %v, %v ---- %v, %v", start, end, hex.EncodeToString(r.buff[:4]), hex.EncodeToString(r.buff[4:8]))
	for {
		mid := getMiddleOffset(start, end)
		buf := r.getRecord(mid)
		cur := getIPFromRecord(buf)

		// log.Printf(">> %v, %v, %v -- %v", start, mid, end, hex.EncodeToString(buf[:4]))

		if end-start == indexLength {
			// log.Printf(">> %v, %v, %v -- %v", start, mid, end, hex.EncodeToString(buf[:4]))
			offset := getAddrFromRecord(buf)
			buf = r.getRecord(mid + indexLength)
			if ip < getIPFromRecord(buf) {
				return offset
			}
			return 0
		}

		// 找到的比较大，向前移
		if cur > ip {
			end = mid
		} else if cur < ip { // 找到的比较小，向后移
			start = mid
		} else {
			return byte3ToUInt32(buf[4:7])
		}

	}
}

func getMiddleOffset(start uint32, end uint32) uint32 {
	records := (end - start) / indexLength
	return start + records/2*indexLength
}

func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
