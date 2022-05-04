package main

import (
	"context"
	"io/ioutil"
	"net"
	"os"

	"git.eagleplan.fun/pkg/gommon/gostring"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"git.eagleplan.fun/geoip/pkg/db"
	"git.eagleplan.fun/geoip/pkg/geoippb"
	"git.eagleplan.fun/geoip/pkg/iputil"
)

type App struct {
	ip2DB *ip2region.Ip2Region
	geoDB *geoip2.Reader

	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	g := &App{
		ctx:    ctx,
		cancel: cancel,
	}
	return g
}

func (app *App) loadGeolite2Database() error {
	geoDB, err := func() (*geoip2.Reader, error) {
		fobj, err := os.Open(*flagGeolite2DB)
		defer fobj.Close()
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(fobj)
		if err != nil {
			return nil, errors.Wrap(err, "read geolite2 content failed")
		}

		db, err := geoip2.FromBytes(data)
		if err != nil {
			return nil, errors.Wrap(err, "parse geolite2 db failed")
		}

		return db, nil
	}()

	if err != nil {
		return err
	}
	app.geoDB = geoDB
	return nil
}

func (app *App) loadIp2Database() error {
	ip2DB, err := ip2region.New(*flagIp2DB)
	if err != nil {
		return errors.Wrap(err, "open ip2db failed")
	}

	app.ip2DB = ip2DB
	return nil
}

func (app *App) initDB() error {
	if err := app.loadGeolite2Database(); err != nil {
		return err
	}
	if err := app.loadIp2Database(); err != nil {
		return err
	}
	return nil
}


func (app *App) IsInArea(ipAddr, area string) bool {
	ip := net.ParseIP(ipAddr)
	if ip == nil || iputil.IsIntranet(ip) {
		return false
	}
	cityData, err := db.FindInDatabase(app.ip2DB, app.geoDB, ip)
	if err != nil {
	return false
	}

	if cityData.Country == area {
		return true
	}
	return false
}

type re struct {
	Ip []string `json:"ip"`
}

func (app *App) onLookups(c *gin.Context) {
	// ip := &re{}
	// err := c.BindJSON(ip)
	// if err != nil {
	// 	render.Failure(c, err.Error())
	// 	return
	// }
	// addr := make([]string, 0, len(ip.Ip))
	// for i := range ip.Ip {
	// 	ips := net.ParseIP(ip.Ip[i])
	// 	if ips == nil || iputil.IsIntranet(ips) {
	// 		cityData := &model.CityResult{
	// 			Country:  "本地",
	// 			Province: "本地",
	// 			City:     "局域网",
	// 		}
	// 		addr = append(addr, cityData.ToString())
	// 		continue
	// 	}
	// 	cityData, err := db.FindInDatabase(app.ip2DB, app.geoDB, ips)
	// 	if err != nil {
	// 		render.Failure(c, err.Error())
	// 		continue
	// 	}
	// 	addr = append(addr, cityData.String())
	// }
	//
	// c.JSON(http.StatusOK, addr)
	// return
}

func (app *App) Lookup(ctx context.Context, req *geoippb.ReqLookup, rsp *geoippb.RspLookup) error {
	var err error
	var cityData *geoippb.CityResult
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

	cityData, err = db.FindInDatabase(app.ip2DB, app.geoDB, ip)

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

func (app *App) cleanUp() {
	if app.geoDB != nil {
		app.geoDB.Close()
		app.geoDB = nil
	}
	if app.ip2DB != nil {
		app.ip2DB.Close()
		app.ip2DB = nil
	}
}
