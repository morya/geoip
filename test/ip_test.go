package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"git.gouboyun.tv/pkg/gommon/gostring"

	"git.gouboyun.tv/live/geoip/pkg/db"
)

func TestSingle(t *testing.T) {
	ips := []string{
		"223.160.226.102",
	}
	for _, ip := range ips {
		record, err := db.FindInDatabase(geoDBTest, net.ParseIP(ip))
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("record=%v", gostring.JsonEncodeString(record))
	}
}

// 本地ip定位服务
func TestAll(t *testing.T) {
	for addr, ip := range ipSamples {
		record, _ := db.FindInDatabase(geoDBTest, net.ParseIP(ip))
		marshal, _ := json.Marshal(record)

		t.Log(string(marshal))
		t.Log(addr, ip)
	}
}

// 高德地图ip定位
func TestGaode(t *testing.T) {
	for _, ip := range ipSamples {
		url := "https://restapi.amap.com/v3/ip?ip=" + ip + "&key=" + key
		req, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		s, _ := ioutil.ReadAll(req.Body)
		fmt.Println(ip, string(s))
	}
}

// 百度地图ip定位
func TestBaidu(t *testing.T) {
	for _, ip := range ipSamples {
		url := "https://api.map.baidu.com/location/ip?ip=" + ip + "&ak=" + ak
		req, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		s, _ := ioutil.ReadAll(req.Body)
		fmt.Println(ip, string(s))
	}
}
