package db

import (
	"net"

	"github.com/oschwald/geoip2-golang"

	"git.gouboyun.tv/live/protos/pb/geoippb"
)

const (
	Geolite2DbKeyCN = "zh-CN"
	Geolite2DbKeyEN = "en"
)

func tranlateGeolite2City(record *geoip2.City, cityData *geoippb.CityResult) {
	key := Geolite2DbKeyCN
	if _, ok := record.City.Names[key]; !ok {
		key = Geolite2DbKeyEN
	}
	cityData.Country = record.Country.Names[key]
	if len(record.Subdivisions) == 0 {
		cityData.Province = cityData.Country
	} else {
		cityData.Province = record.Subdivisions[0].Names[key]
	}
	cityData.City = record.City.Names[key]
}

func FindInGeolite2DB(geodb *geoip2.Reader, ip net.IP) (cityData *geoippb.CityResult, err error) {
	record, err := geodb.City(ip)
	if err != nil {
		return nil, err
	}
	cityData = &geoippb.CityResult{}
	tranlateGeolite2City(record, cityData)
	return
}

// FindInDatabase 先查ip2db 找不到则查geolite2
func FindInDatabase(geodb *geoip2.Reader, ip net.IP) (cityData *geoippb.CityResult, err error) {
	return FindInGeolite2DB(geodb, ip)
}
