package main

import (
	router2 "aura/router"
	"aura/subscriber"
	"github.com/gammazero/nexus/router"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
)

var routerStatusChan = make(chan string)

var logger = log.New(os.Stdout, "aura: ", 0)

var auraRouterSetup = router2.Options{
	Realms: []*router.RealmConfig{
		&router.RealmConfig{
			URI: "nexus.aura.realm.1",
			AnonymousAuth: true,
		},
		&router.RealmConfig{
			URI: "nexus.aura.realm.2",
			AnonymousAuth: true,
		},
	},
	Url:          "127.0.0.1:3131",
	AllowOrigins: []string{"localhost:*"},
	Logger: logger,
	Status: routerStatusChan,
}

var subscription = subscriber.Options{
	Realm: "nexus.aura.realm.1",
	Topic: "Default",
	Logger: logger,
	WSUrl: wsURL(auraRouterSetup.Url),
	EventHandler: func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {

	},
	Debug: true,
}



func main() {
	go router2.StartRouter(auraRouterSetup)
	for status := range routerStatusChan {
		if status == "listening" {
			go subscriber.InitializeSubscriber(subscription)
		}

		if status == "closed" {
			logger.Println("Router Shutdown")
			break
		}
	}

}

func wsURL(url string) string {
	return "ws://" + url
}
