package http

import (
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	v1 "starland-account/api/http/v1"
	"starland-account/configs"
	"starland-account/internal/pkg/middlewares"
	"starland-account/internal/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewHTTPServer(config *configs.Config, us *service.Service) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.HTTP.ReadTimeout * time.Second,
		WriteTimeout: config.HTTP.WriteTimeout * time.Second,
	})
	app.Use(recover.New(), pprof.New(), cors.New(), requestid.New())
	prometheus := fiberprometheus.New("starland-account")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Use(logger.New(logger.Config{
		Format: fmt.Sprintf("${time} | ${ip} | ${status} | ${locals:%s} | ${latency} | ${method} | ${path} | "+
			"ResponseBody:${resBody} | Params:${queryParams} \n",
			requestid.ConfigDefault.ContextKey),
		Next: func(c *fiber.Ctx) bool {
			path := string(c.Request().URI().Path())
			if strings.Contains(path, "/media/v1/file") {
				return true
			}
			return false
		},
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Shanghai",
	}))
	app.Use(middlewares.Auth())
	r := app.Group("/")
	v1.InitAccountRouter(r, us.Account, config)
	v1.InitActivityRouter(r, us.Activity, config)
	zap.S().Infof("addr:%s", config.HTTP.Addr)
	return app, nil
}
