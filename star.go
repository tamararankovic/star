package main

import (
	fPb "github.com/c12s/scheme/flusher"
	"github.com/c12s/star/syncer"
	actor "github.com/c12s/starsystem"
	"github.com/golang/protobuf/proto"
	"strings"
)

type StarAgent struct {
	Conf     *Config
	f        syncer.Syncer
	system   *actor.System
	pointers []string //for faster lookup
}

func NewStar(c *Config, f syncer.Syncer) *StarAgent {
	return &StarAgent{
		Conf: c,
		f:    f,
	}
}

func (s *StarAgent) addActors(actors map[string]actor.Actor) {
	for key, a := range actors {
		ac := s.system.ActorOf(key, a)
		s.pointers = append(s.pointers, ac.Name())
	}
}

func (s *StarAgent) contains(e string) string {
	for _, a := range s.pointers {
		if strings.Contains(a, e) {
			return a
		}
	}
	return ""
}

func (s *StarAgent) Start(actors map[string]actor.Actor) {
	s.system = actor.NewSystem(s.Conf.NodeId)
	s.addActors(actors)
	s.f.Subscribe(s.Conf.RTopic, func(event []byte) {
		data := &fPb.Event{}
		err := proto.Unmarshal(event, data)
		if err == nil {
			key := s.contains(data.Kind)
			if key != "" {
				a := s.system.ActorSelection(key)
				a.Tell(StarMessage{Data: data})
			}
		}
	})
}

func (s *StarAgent) Stop() {
	s.system.Shutdown()
}
