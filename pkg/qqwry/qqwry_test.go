package qqwry

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name      string
		args      args
		wantStart uint32
		wantEnd   uint32
		wantErr   bool
	}{
		{"correct db", args{"testdata/qqwry.dat"}, 6718931, 10417920, false},
		{"not exists", args{"testdata/not-exists"}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.file)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.start != tt.wantStart || got.end != tt.wantEnd {
				t.Errorf("Open() = { start: %v, end: %v }, want { start: %v, end: %v }",
					got.start, got.end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestDB_readRecord(t *testing.T) {
	r, _ := Open("testdata/qqwry.dat")

	type args struct {
		offset uint32
	}
	tests := []struct {
		name string
		args args
		want *Record
	}{
		{"first record", args{offset: getAddrFromRecord(r.buff[r.start : r.start+7])}, &Record{Country: "IANA", City: "保留地址"}},
		{"last record", args{offset: getAddrFromRecord(r.buff[r.end : r.end+7])}, &Record{Country: "纯真网络", City: "2021年02月02日IP数据"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.readRecord(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.readRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Query(t *testing.T) {
	db, _ := Open("testdata/qqwry.dat")

	type args struct {
		ip net.IP
	}
	tests := []struct {
		name    string
		args    args
		want    *Record
		wantErr bool
	}{
		{"192.168.1.1", args{net.IP{192, 168, 1, 1}}, &Record{Country: "局域网", City: "对方和您在同一内部网"}, false},
		{"invalid ip", args{nil}, nil, true},
		{"first record", args{net.IP{0, 0, 0, 0}}, &Record{Country: "IANA", City: "保留地址"}, false},
		{"last record", args{net.IP{255, 255, 255, 255}}, &Record{Country: "纯真网络", City: "2021年02月02日IP数据"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Query(tt.args.ip)
			if tt.wantErr {
				if err == nil {
					t.Errorf("DB.Query() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDB_readRecord(b *testing.B) {
	db, _ := Open("testdata/qqwry.dat")

	ip := net.IP{192, 168, 1, 1}
	iplong := binary.BigEndian.Uint32(ip)
	offset := db.search(iplong)

	for n := 0; n < b.N; n++ {
		db.readRecord(offset)
	}
}

func BenchmarkDB_Query(b *testing.B) {
	db, _ := Open("testdata/qqwry.dat")

	for n := 0; n < b.N; n++ {
		db.Query(net.IP{192, 168, 1, 1})
	}
}
