package main

import (
	"flag"
	"log"
	"strings"

	_ "github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"

	"git.eagleplan.fun/geoip/pkg/geoippb"
)

var (
	flagGeolite2DB = flag.String("db", "ip-database/geolite2-city.mmdb", "")
	flagIp2DB      = flag.String("db2", "ip-database/ip2region.db", "")
	flagLogLevel   = flag.String("loglevel", "info", "[debug,info,warn/error/none]")
	flagLogFormat  = flag.String("logformat", "json", "json/text")
)

func initLogging() {
	l, _ := logrus.ParseLevel(*flagLogLevel)
	logrus.SetLevel(l)

	switch strings.ToLower(*flagLogFormat) {
	case "text":

	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
func main() {
	flag.Parse()
	initLogging()

	var app = NewApp()
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}

	var svc = micro.NewService(micro.Name("geoip"))
	var s = svc.Server()
	geoippb.RegisterGeoipHandler(s, app)

	svc.Init()

	if err := svc.Run(); err != nil {
		logrus.WithError(err).Error("failed")
	}
}
