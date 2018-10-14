package main

import (
	"fmt"
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
const auraRealm2 = "nexus.aura.realm2"
const wsURL = string("ws://" + address)

func main() {
	go startServer(auraRealm1)
	initializeSubscriber()
}

func startServer(realm string) {

	// router config
	config := &router.Config{
		//Debug: true,
		RealmConfigs: []*router.RealmConfig{
			{
				URI: wamp.URI(realm),
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

	//initializeSubscriber()

	closer, err := wsS.ListenAndServe(address)
	if err != nil {
		log.Fatal("Ws Server Error: ", err)
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
		//Debug: true,
	}


	fmt.Println("Test")

	subscriber, err := client.ConnectNet(wsURL, clientCfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer subscriber.Close()


	// event handler
	eventHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		fmt.Println("Eventhandler triggered with topic ", topic)
		if len(args) != 0 {
			logger.Println("Message: ", args[0])
		}
		InitializePublisher()
	}

	err = subscriber.Subscribe(topic, eventHandler, nil)
	if err != nil {
		logger.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
	case <-subscriber.Done():
		logger.Println("Router Gone Exiting")
		return
	}

	if err = subscriber.Unsubscribe(topic); err != nil {
		logger.Println("Failed to unsubscribe", err)
	}
}

func InitializePublisher() {
	fmt.Println("Publishing!")
	clientCfg := client.Config{
		Realm: auraRealm2,
	}
	publisher, err := client.ConnectNet(wsURL, clientCfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer publisher.Close()

	err = publisher.Publish(topic, nil, wamp.List{"Hello World, via teh Publisher!"}, nil)
	if err != nil {
		logger.Fatal("Publish Failed: ", err)
	}
	logger.Println("Published, ", topic, "event")
}