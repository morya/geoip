package main

import (
	"context"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"

	"git.gouboyun.tv/live/protos/pb/geoippb"
	"git.gouboyun.tv/pkg/gommon/gostring"

	"git.gouboyun.tv/live/geoip/geo2litedb"
	"git.gouboyun.tv/live/geoip/pkg/db"
	"git.gouboyun.tv/live/geoip/pkg/ipsearch"
	"git.gouboyun.tv/live/geoip/pkg/iputil"
	"git.gouboyun.tv/live/geoip/pkg/model"
)

type App struct {
	geoDB *geoip2.Reader
}

func NewApp() *App {
	app := &App{}
	return app
}

func (app *App) loadGeolite2Database() error {
	app.geoDB = geo2litedb.DB
	return nil
}

func (app *App) initDB() error {
	if err := app.loadGeolite2Database(); err != nil {
		return err
	}
	ipsearch.New()
	return nil
}

func (app *App) IsInArea(ipAddr, area string) bool {
	ip := net.ParseIP(ipAddr)
	if ip == nil || ip.IsPrivate() {
		return false
	}
	cityData, err := db.FindInDatabase(app.geoDB, ip)
	if err != nil {
		return false
	}

	if cityData.Country == area {
		return true
	}
	return false
}

func (app *App) LookupIp(ip string) (*model.RspLookup, error) {
	s, _ := ipsearch.New()
	result := s.Get(ip)
	rsp := &model.RspLookup{}
	if result == "" {
		rsp.Message = "No IP found"
		rsp.Code = 100
	}

	arr := strings.Split(result, "|")
	if len(arr) != 11 {
		// 亚洲|中国|湖北| |潜江|联通|429005|China|CN|112.896866|30.421215
		rsp.Message = "bad ip position"
		rsp.Code = 101
	}

	rsp.Data.Country = arr[1]
	rsp.Data.Province = arr[2]
	rsp.Data.City = arr[3]
	if rsp.Data.City == "" {
		rsp.Data.City = arr[4]
	}

	return rsp, nil
}

func (app *App) Lookup(ctx context.Context, req *geoippb.ReqLookup, rsp *geoippb.RspLookup) error {
	var err error
	var cityData *geoippb.CityResult = &geoippb.CityResult{
		Country: "未知国家",
		City:    "未知城市",
	}

	ipAddr := req.Ip
	ip := net.ParseIP(ipAddr)
	if ip == nil || iputil.IsIntranet(ip) {
		cityData = &geoippb.CityResult{
			Country:  "本地",
			Province: "本地",
			City:     "局域网",
		}
		rsp.Result = cityData
		return nil
	}

	s, _ := ipsearch.New()
	result := s.Get(req.Ip)
	if result != "" {
		arr := strings.Split(result, "|")
		if len(arr) == 11 {
			// 亚洲|中国|湖北| |潜江|联通|429005|China|CN|112.896866|30.421215
			cityData.Country = arr[1]
			cityData.Province = arr[2]
			cityData.City = arr[3]
			if cityData.City == "" {
				cityData.City = arr[4]
			}
			rsp.Result = cityData
			return nil
		}
	}

	cityData, err = db.FindInDatabase(app.geoDB, ip)
	if err == nil {
		rsp.Result = cityData
	}
	logrus.WithError(err).Infof("city data %v", gostring.JsonEncodeString(cityData))
	return nil
}

func (app *App) LookupBatch(ctx context.Context, req *geoippb.ReqLookupBatch, rsp *geoippb.RspLookupBatch) error {
	return nil
}

func (app *App) Init() error {
	if err := app.initDB(); err != nil {
		return err
	}

	return nil
}
