package services

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type ClusterJoinListener struct {
	conn   *nats.Conn
	serf   *SerfAgent
	nodeId string
}

func NewClusterJoinListener(conn *nats.Conn, serf *SerfAgent, nodeId string) *ClusterJoinListener {
	return &ClusterJoinListener{
		conn:   conn,
		serf:   serf,
		nodeId: nodeId,
	}
}

func (l *ClusterJoinListener) Listen() {
	_, err := l.conn.Subscribe(fmt.Sprintf("%s.join", l.nodeId), func(msg *nats.Msg) {
		address := string(msg.Data)
		// todo: what happens if the node is not a member of a cluster?
		//l.serf.Leave()
		err := l.serf.Join(address)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}
}
