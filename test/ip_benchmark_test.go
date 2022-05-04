package test

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"git.eagleplan.fun/geoip/pkg/db"
)

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for addr, ip := range ipSamples {
			record, _ := db.FindInDatabase(ip2DBTest, geoDBTest, net.ParseIP(ip))
			marshal, _ := json.Marshal(record)
			fmt.Println(string(marshal))
			fmt.Println(addr, ip)
		}
	}
}
