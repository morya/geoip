package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"git.gouboyun.tv/live/geoip/pkg/model"
)

func httpHandle(app *App) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		time.Sleep(time.Second)
		log.Fatal(err)
	}

	e := echo.New()

	// 健康检查端点
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "geoip",
			"version": "1.0.0",
		})
	})

	// 就绪检查端点
	e.GET("/ready", func(c echo.Context) error {
		// 检查数据库是否已加载
		if app.geoDB == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
				"reason": "database not loaded",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
		})
	})

	// IP查询端点
	e.GET("/geoip/lookup/:ip", func(c echo.Context) error {
		rsp := &model.RspLookup{}
		var err error

		if httpAuthKey != "" {
			key := c.QueryParam("key")
			if key != httpAuthKey {
				rsp.Code = 1
				rsp.Message = "key not found"
				return c.JSON(http.StatusBadRequest, rsp)
			}
		}

		ip := c.Param("ip")
		if ip == "" {
			return c.JSON(http.StatusBadRequest, rsp)
		}
		rsp, err = app.LookupIp(ip)
		if err != nil {
			rsp.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, rsp)
		}

		return c.JSON(http.StatusOK, rsp)
	})

	http.Serve(l, e)
}
