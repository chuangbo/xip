package qqwry

import (
	"os"
	"testing"
)

func Test_readUpdateInfo(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantKey     uint32
		wantVersion string
		wantErr     bool
	}{
		{"correct file", "testdata/copywrite.rar", 225, "纯真IP地址数据库 2021年02月25日", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.Open(tt.filename)
			got, got1, err := readUpdateInfo(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUpdateInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantKey {
				t.Errorf("getUpdateInfo() got = %v, wantKey %v", got, tt.wantKey)
			}
			if got1 != tt.wantVersion {
				t.Errorf("getUpdateInfo() got1 = %v, want %v", got1, tt.wantVersion)
			}
		})
	}
}

func TestSameVersion(t *testing.T) {
	type args struct {
		remote string
		local  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"outdated", args{"纯真IP地址数据库 2021年02月25日", "2021年02月02日IP数据"}, false},
		{"same", args{"纯真IP地址数据库 2021年02月25日", "2021年02月25日IP数据"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameVersion(tt.args.remote, tt.args.local); got != tt.want {
				t.Errorf("SameVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
