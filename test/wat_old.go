package main

import (
	"fmt"

	"github.com/cloudfoundry/gibson"
	"github.com/cloudfoundry/yagnats"
)

func main() {

	nats := yagnats.NewClient()

	err := nats.Connect(&yagnats.ConnectionInfo{
		Addr:     "127.0.0.1:4222",
		Username: "user",
		Password: "pass",
	})
	if err != nil {
		panic("Wrong auth or something.")
	}

	client := gibson.NewCFRouterClient("127.0.0.1", nats)
	fmt.Printf("\n Client: %#v\n", client)

	client.Greet()

	client.Register(4567, "test.vcap.me")
	//client.Register(4567, "localhost")
}
