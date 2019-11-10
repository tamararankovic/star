package main

import (
	"fmt"
	"github.com/c12s/star/flusher/nats"
	actor "github.com/c12s/starsystem"
	"runtime"
)

func main() {
	config, err := ConfigFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	sync, err2 := nats.NewNatsSync(config.Flusher)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	star := NewStar(config, sync)
	star.Start(
		map[string]actor.Actor{
			"configs": ConfigsActor{},
			"secrets": SecretsActor{},
			"actions": ActionsActor{},
		})

	fmt.Println("Starting project star...")
	runtime.Goexit()
	star.Stop()
}
