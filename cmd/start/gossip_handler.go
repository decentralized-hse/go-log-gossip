package start

import (
	"github.com/decentralized-hse/go-log-gossip/domain/features/dtos"
	"github.com/decentralized-hse/go-log-gossip/domain/features/logs/commands"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/mehdihadeli/go-mediatr"
	"log"
)

func gossipHandler(message *gossip.Message) {
	log.Println(message)
	switch message.MessageType {
	case gossip.Push:
		logDto := dtos.LogDTO{
			Hash:     message.Payload["hash"].(string),
			Position: int(message.Payload["position"].(float64)),
			NodeId:   message.Payload["node_id"].(string),
			Message:  message.Payload["message"].(string),
		}
		logRecord, err := logDto.ToDomain()
		if err != nil {
			return
		}
		_, err = mediatr.Send[*commands.AddLogCommand, *commands.AddLogResponse](ctx, commands.NewAddLogCommand(logRecord))
		if err != nil {
			log.Fatalf("Error occured err = %v", err)
		}
		break
	case gossip.Pull:
		logPosition := int(message.Payload["position"].(float64))
		nodeId := message.Payload["node_id"].(string)
		_, err := mediatr.Send[*commands.SendLogCommand, *commands.SendLogResponse](ctx,
			commands.NewSendLogCommand(message.Sender, nodeId, logPosition))
		if err != nil {
			log.Fatalf("Error occured err = %v", err)
		}
		break
	}
	log.Println(message)
}
