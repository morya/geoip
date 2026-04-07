package main

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"

	"git.gouboyun.tv/pkg/gommon/mytrace"
)

func recoveryHandler(h server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		defer func() {
			e := recover()

			stack := make([]byte, 4<<10) // 4k
			length := runtime.Stack(stack, false)
			stack = stack[:length]

			funcName := mytrace.GetCallerName()
			l := logrus.WithFields(
				logrus.Fields{
					"req":   req.Body(),
					"ack":   rsp,
					"func":  funcName,
					"stack": string(stack),
				})

			md, _ := metadata.FromContext(ctx)
			if md != nil {
				l = l.WithField("meta", md)
			}

			if e != nil {
				l.Error("call handler panic:", e)
			}
		}()
		return h(ctx, req, rsp)
	}
}
