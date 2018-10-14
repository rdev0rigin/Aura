package router

import (
	"github.com/gammazero/nexus/router"
	"log"
	"os"
	"os/signal"
)

type Options struct {
	Realms []*router.RealmConfig
	Url string
	AllowOrigins []string
	Debug bool
	Logger *log.Logger
}


func StartRouter(setup Options) {
	config := &router.Config{
		Debug: setup.Debug,
		RealmConfigs: setup.Realms,
	}

	auraRouter, err := router.NewRouter(config, setup.Logger)

	if err != nil {
		log.Fatal(err)
	}

	defer auraRouter.Close()

	wsS := router.NewWebsocketServer(auraRouter)

	wsS.AllowOrigins(setup.AllowOrigins)

	closer, err := wsS.ListenAndServe(setup.Url)
	if err != nil {
		log.Fatal("Ws Server Error: ", err)
	}

	log.Printf("Websocket server listening on ws://%s/", setup.Url)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<- shutdown
	closer.Close()
}
