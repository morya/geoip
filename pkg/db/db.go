package db

import (
	"net"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/oschwald/geoip2-golang"

	"git.eagleplan.fun/geoip/pkg/geoippb"
)

const (
	Geolite2DbKeyCN = "zh-CN"
	Geolite2DbKeyEN = "en"
)

func FindInIp2DB(ip2db *ip2region.Ip2Region, ip string) (*geoippb.CityResult, error) {
	var city = &geoippb.CityResult{}

	search, err := ip2db.MemorySearch(ip)
	if err != nil {
		return nil, err
	}

	if search.Province == "0" || search.City == "0" {
		return nil, nil
	}

	city.Country = search.Country
	city.Province = search.Province
	city.City = search.City

	return city, nil
}

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
func FindInDatabase(ip2db *ip2region.Ip2Region, geodb *geoip2.Reader, ip net.IP) (cityData *geoippb.CityResult, err error) {
	cityData, err = FindInIp2DB(ip2db, ip.String())
	if err != nil || cityData == nil || cityData.City == cityData.Province {
		cityData, err = FindInGeolite2DB(geodb, ip)
	}
	return
}
