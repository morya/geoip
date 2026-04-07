package main

import (
	"os"
	"strings"

	_ "github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"

	"git.gouboyun.tv/live/protos/pb/geoippb"
)

var (
	logLevel  string
	logFormat string

	httpAuthKey string
	listenPort  int
)

func initLogging() {
	l, _ := logrus.ParseLevel(logLevel)
	logrus.SetLevel(l)

	switch strings.ToLower(logFormat) {
	case "text":

	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func buildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "logformat",
			Value:       "json",
			Destination: &logFormat,
			Usage:       "json/text",
		},
		&cli.StringFlag{
			Name:        "loglevel",
			Value:       "info",
			Usage:       "[debug,info,warn/error/none]",
			Destination: &logLevel,
		},
		&cli.StringFlag{
			Name:        "key",
			Value:       "",
			Usage:       "specify http call key",
			EnvVars:     []string{"HTTP_KEY"},
			Destination: &httpAuthKey,
		},
		&cli.IntFlag{
			Name:        "port",
			EnvVars:     []string{"PORT"},
			Usage:       "port to listen",
			Destination: &listenPort,
			Value:       80,
		},
	}
}

func action(c *cli.Context) error {
	os.Args = os.Args[:1]
	app := NewApp()
	if err := app.Init(); err != nil {
		logrus.Fatal(err)
	}

	go httpHandle(app)

	svc := micro.NewService(
		micro.Name("lemon.geoip"),
		micro.WrapHandler(recoveryHandler),
	)
	s := svc.Server()
	geoippb.RegisterGeoipHandler(s, app)

	svc.Init()

	if err := svc.Run(); err != nil {
		logrus.WithError(err).Error("failed")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Before = func(c *cli.Context) error {
		initLogging()
		return nil
	}
	app.Flags = buildFlags()
	app.Action = action

	logrus.Infof("starting")
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("bye")
}
