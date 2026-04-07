package geo2litedb

import (
	_ "embed"

	"github.com/oschwald/geoip2-golang"
)

//go:embed geolite2-city.mmdb
var dbfile []byte

var DB *geoip2.Reader

func init() {
	var err error
	DB, err = geoip2.FromBytes(dbfile)
	if err != nil {
		panic(err)
	}
}
