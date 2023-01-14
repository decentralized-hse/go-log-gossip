package gossip

type MessageType string

const (
	Push MessageType = "push"
	Pull             = "pull"
)
