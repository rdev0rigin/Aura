package main

import (
	router2 "aura/router"
	"aura/subscriber"
	"github.com/gammazero/nexus/router"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
)

var onRouterStart = make(chan string)

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
	router2.StartRouter(auraRouterSetup)
	//for {
		//subscriber.InitializeSubscriber(subscription) <- onRouterStart

	//}

}

func wsURL(url string) string {
	return "ws://" + url
}
