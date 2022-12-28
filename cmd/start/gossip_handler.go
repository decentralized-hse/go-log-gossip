package start

import (
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"log"
)

func gossipHandler(message *gossip.Message) {
	// todo: handle message
	log.Println(message)
}
