package db

import (
	"strings"

	"git.eagleplan.fun/geoip/pkg/geoippb"
)

// 判断省市归属
func IsInArea(area []string, cityData *geoippb.CityResult) bool {
	if len(area) == 0 {
		return false
	}
	for i := range area {
		if i == 0 {
			if strings.Index(cityData.Province, area[i]) == -1 {
				return false
			}
		} else {
			if strings.Index(cityData.City, area[i]) == -1 {
				return false
			}
		}
	}
	return true
}
