package gossip

type Message struct {
	MessageType MessageType `json:"message_type"`
	Sender      string      `json:"sender"`
	Signature   string      `json:"signature"`
	Payload     any         `json:"data"`
}
