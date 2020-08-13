package webserver

import (
	"context"
	"github.com/giantliao/beatles/webserver/api"

	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/port"
	"log"
	"net/http"
	"strconv"
	"time"
)

var webserver *http.Server

func StartWebDaemon() {
	mux := http.NewServeMux()

	cfg := config.GetCBtl()

	mux.Handle(cfg.GetListMinersWebPath(), &api.BeatlesMasterProxy{})
	mux.Handle(cfg.GetpurchaseWebPath(),&api.BeatlesMasterProxy{})

	if cfg.HttpServerPort == 0{
		cfg.HttpServerPort = port.HttpPort()
	}

	addr := ":" + strconv.Itoa(cfg.HttpServerPort)

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)

	log.Println("Web Server Stopped")
}
