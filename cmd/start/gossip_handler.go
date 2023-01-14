package start

import (
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"log"
)

func gossipHandler(message *gossip.Message) {
	switch message.MessageType {
	case gossip.Push:
		break
	case gossip.Pull:
	}
	log.Println(message)
}
