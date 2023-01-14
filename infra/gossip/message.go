package gossip

import "github.com/decentralized-hse/go-log-gossip/domain"

type Message struct {
	MessageType MessageType   `json:"message_type"`
	Sender      domain.NodeId `json:"sender"`
	Signature   string        `json:"signature"`
	Payload     any           `json:"data"`
	Meta        struct {
		Addr string `json:"addr"`
	} `json:"meta"`
}
