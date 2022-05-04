package test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
)

var (
	ipSamples map[string]string
	ip2DBTest *ip2region.Ip2Region
	geoDBTest *geoip2.Reader
)

const (
	key = ""
	ak  = ""
)

func TestMain(m *testing.M) {
	err := loadGeolite2Database()
	if err != nil {
		log.Println(err)
		return
	}
	err = loadIp2Database()
	if err != nil {
		log.Println(err)
		return
	}
	ipSamples = map[string]string{
		"中国, 天津, 天津": "117.9.60.131",
		"中国, 四川, 成都": "119.4.41.224",
		"中国, 四川, 巴中": "221.10.33.100",
		"中国, 贵州, 遵义": "222.86.165.44",
		"中国, 江苏, 徐州": "180.104.46.118",
		"中国, 辽宁, 沈阳": "124.92.145.207",
		"中国, 吉林, 四平": "222.168.162.194",
		"中国, 中国,":    "103.193.195.75",
		"中国, 河北, 保定": "27.187.36.8",
		"中国, 安徽,":    "36.161.56.23",
		"中国, 香港，":    "47.75.108.7",
		"泰国, 曼谷":     "128.1.39.112",
	}
	os.Exit(m.Run())
}

func loadGeolite2Database() error {
	fobj, err := os.Open("../ip-database/geolite2-city.mmdb")
	defer fobj.Close()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(fobj)
	if err != nil {
		return errors.Wrap(err, "read geolite2 content failed")
	}

	db, err := geoip2.FromBytes(data)
	if err != nil {
		return errors.Wrap(err, "parse geolite2 db failed")
	}
	geoDBTest = db
	return nil
}

func loadIp2Database() error {
	ip2DB, err := ip2region.New("../ip-database/ip2region.db")
	if err != nil {
		return errors.Wrap(err, "open ip2db failed")
	}
	ip2DBTest = ip2DB
	return nil
}
