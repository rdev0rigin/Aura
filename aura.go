package main

import (
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/router"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
	"os/signal"
)

const address = "127.0.0.1:3131"
const topic = "hello.wamp"
var logger = log.New(os.Stdout, "aura: ", 0)
const auraRealm1 = "nexus.aura.realm1"

func main() {
	startServer()
	initializeSubscriber()
}

func startServer() {
	// router config
	config := &router.Config{
		Debug: true,
		RealmConfigs: []*router.RealmConfig{
			{
				URI: wamp.URI(auraRealm1),
				AnonymousAuth: true,
			},
		},
	}

	auraRouter, err := router.NewRouter(config, logger)

	if err != nil {
		log.Fatal(err)
	}

	defer auraRouter.Close()


	wsS := router.NewWebsocketServer(auraRouter)
	wsS.AllowOrigins([]string{"localhost:63343"})

	closer, err := wsS.ListenAndServe(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Websocket server listening on ws://%s/", address)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<- shutdown
	closer.Close()
}

func initializeSubscriber() {

	clientCfg := client.Config{
		Realm: auraRealm1,
		Logger: logger,
		Debug: true,
	}

	const wsURL = string("ws://" + address)

	subscribber, err := client.ConnectNet(wsURL, clientCfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer subscribber.Close()


	// event handler
	eventHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		logger.Println("Eventhandler triggered with topic ", topic)
		if len(args) != 0 {
			logger.Println("Message: ", args[0])
		}
	}

	err = subscribber.Subscribe(topic, eventHandler, nil)
	if err != nil {
		logger.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
	case <-subscribber.Done():
		logger.Println("Router Gone Exiting")
		return
	}

	if err = subscribber.Unsubscribe(topic); err != nil {
		logger.Println("Failed to unsubscribe", err)
	}
}