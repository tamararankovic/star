package main

import (
	"fmt"
	fPb "github.com/c12s/scheme/flusher"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"runtime"
)

const subject = "topology.regions.>"

func main() {
	natsConnection, _ := nats.Connect("0.0.0.0:4222")
	natsConnection.Subscribe(subject, func(msg *nats.Msg) {
		data := fPb.FlushPush{}
		err := proto.Unmarshal(msg.Data, &data)
		if err == nil {
			// Handle the message
			fmt.Println(data)
		}
	})

	// Keep the connection alive
	runtime.Goexit()

}
