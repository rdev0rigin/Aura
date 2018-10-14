package subscriber

import (
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
	"os/signal"
)

type Options struct {
	Realm string
	Topic string
	Logger *log.Logger
	WSUrl string
	EventHandler func(args wamp.List, kwargs wamp.Dict, details wamp.Dict)
	Debug bool
}


func InitializeSubscriber(setup Options) {

	clientCfg := client.Config{
		Realm: setup.Realm,
		Logger: setup.Logger,
		//Debug: true,
	}

	subscriber, err := client.ConnectNet(setup.WSUrl, clientCfg)
	if err != nil {
		setup.Logger.Fatal("Subscriber Failed: ", err)
	}
	setup.Logger.Println("Test")
	defer subscriber.Close()

	err = subscriber.Subscribe(setup.Topic, setup.EventHandler, nil)
	if err != nil {
		setup.Logger.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
	case <-subscriber.Done():
		setup.Logger.Println("Router Gone Exiting")
		return
	}

	if err = subscriber.Unsubscribe(setup.Topic); err != nil {
		setup.Logger.Println("Failed to unsubscribe", err)
	}
}
