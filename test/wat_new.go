package main

import (
	"fmt"

	"github.com/cloudfoundry/gibson"
	"github.com/cloudfoundry/yagnats"
)

func main() {
	// err := nats.Connect(&yagnats.ConnectionInfo{
	// 	Addr:     "127.0.0.1:4222",
	// 	Username: "user",
	// 	Password: "pass",
	// })
	// if err != nil {
	// 	panic("Wrong auth or something.")
	// }
	natsConn, err := yagnats.Connect([]string{"nats://user:pass@127.0.0.1:4242"})
	// _, err := yagnats.Connect([]string{"nats://user:pass@127.0.0.1:4222"})
	if err != nil {
		fmt.Println("Error: ", err)
		panic("Error connecting to nats")
	}

	// fmt.Printf("natsConn: %#v\n", natsConn)
	// fmt.Printf("\nping: %#v\n", natsConn.Ping())

	client := gibson.NewCFRouterClient("127.0.0.1", natsConn)

	// fmt.Printf("\n Client: %#v\n", client)

	client.Greet()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	panic("Error client GREET")
	// }

	err = client.Register(4567, "yuyu.vcap.me")
	fmt.Printf("\nclient.Register, err: %#v\n", err)

	// registrations := natsConn.PublishedMessages("router.register")

	// fmt.Printf("\n registrations: %#v\n", registrations)
}
