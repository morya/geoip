package db

import (
	"strings"

	"git.gouboyun.tv/live/protos/pb/geoippb"
)

// 判断省市归属
func IsInArea(area []string, cityData *geoippb.CityResult) bool {
	if len(area) == 0 {
		return false
	}
	for i := range area {
		if i == 0 {
			if !strings.Contains(cityData.Province, area[i]) {
				return false
			}
		} else {
			if !strings.Contains(cityData.City, area[i]) {
				return false
			}
		}
	}
	return true
}
