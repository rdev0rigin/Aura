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
	Status chan string
}


func statusCall(channel chan string, message string) {
	channel <- message
}

func StartRouter(options Options) {
	config := &router.Config{
		Debug:        options.Debug,
		RealmConfigs: options.Realms,
	}

	auraRouter, err := router.NewRouter(config, options.Logger)

	if err != nil {
		log.Fatal(err)
	}
	defer statusCall(options.Status, "closed")
	defer auraRouter.Close()

	wsS := router.NewWebsocketServer(auraRouter)

	wsS.AllowOrigins(options.AllowOrigins)

	closer, err := wsS.ListenAndServe(options.Url)
	if err != nil {
		log.Fatal("Ws Server Error: ", err)
	}

	options.Status <- "listening"

	log.Printf("Websocket server listening on ws://%s/", options.Url)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<- shutdown
	closer.Close()
}
