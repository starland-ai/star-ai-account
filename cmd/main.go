package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "starland-account/api/http"
	config "starland-account/configs"
	"starland-account/internal/pkg/logs"

	"go.uber.org/zap"
)

func main() {
	// init config,log
	config.InitConfig()
	cfg := config.GetConfig()
	logs.InitLogging(cfg)

	var host string
	flag.StringVar(&host, "h", cfg.HTTP.Addr, "host")
	flag.Parse()
	ln, err := net.Listen("tcp", host)
	if err != nil {
		zap.S().Fatalf("net listen is err: %s", err.Error())
	}

	s, err := initApp(cfg)
	if err != nil {
		zap.S().Fatalf("dependency injection is err: %s", err.Error())
	}
	app, err := api.NewHTTPServer(cfg, s)

	go func() {
		if err := app.Listener(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//nolint:gomnd
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("shutting down service...")

	if err = app.Shutdown(); err != nil {
		log.Printf("shutting down service : %s", err.Error())
	}

	log.Print("bye")
}
