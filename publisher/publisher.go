package publisher

import (
	"fmt"
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
	"log"
)

type Options struct {
	Realm string
	Topic string
	Logger *log.Logger
	WSUrl string
}

func InitializePublisher(options Options) {
	fmt.Println("Publishing!")
	clientCfg := client.Config{
		Realm: options.Realm,
	}
	publisher, err := client.ConnectNet(options.WSUrl, clientCfg)
	if err != nil {
		options.Logger.Fatal(err)
	}
	defer publisher.Close()

	err = publisher.Publish(options.Topic, nil, wamp.List{"Hello World, via teh Publisher!"}, nil)
	if err != nil {
		options.Logger.Fatal("Publish Failed: ", err)
	}
	options.Logger.Println("Published, ", options.Topic, "event")
}