package main

import (
	"fmt"
	bPb "github.com/c12s/scheme/blackhole"
	"github.com/c12s/star/syncer"
	actor "github.com/c12s/starsystem"
)

//
// Star Message
//
type StarMessage struct {
	Data []*bPb.Payload
}

func (m StarMessage) Name() string {
	return "StarMessage"
}

func (m StarMessage) Params() map[string][]byte {
	return nil
}

//
// Configs Actor
//
type ConfigsActor struct {
	uploader syncer.Uploader
}

func (m ConfigsActor) Receive(msg interface{}, context *actor.ActorProp) {
	switch data := msg.(type) {
	case StarMessage:
		fmt.Println("Received Configs")
		fmt.Println(data)
	default:
		fmt.Println("Error")
	}
}

//
// Actions Actor
//
type ActionsActor struct {
	uploader syncer.Uploader
}

func (m ActionsActor) Receive(msg interface{}, context *actor.ActorProp) {
	switch data := msg.(type) {
	case StarMessage:
		fmt.Println("Received Actions")
		fmt.Println(data)
	default:
		fmt.Println("Error")
	}
}

//
// Secrets Actor
//
type SecretsActor struct {
	uploader syncer.Uploader
}

func (m SecretsActor) Receive(msg interface{}, context *actor.ActorProp) {
	switch data := msg.(type) {
	case StarMessage:
		fmt.Println("Received Secrets")
		fmt.Println(data)
	default:
		fmt.Println("Error")
	}
}
