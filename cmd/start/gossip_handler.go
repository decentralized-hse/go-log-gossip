package start

import (
	"github.com/decentralized-hse/go-log-gossip/domain"
	"github.com/decentralized-hse/go-log-gossip/domain/features/logs/commands"
	"github.com/decentralized-hse/go-log-gossip/infra/gossip"
	"github.com/mehdihadeli/go-mediatr"
	"log"
)

func gossipHandler(message *gossip.Message) {
	switch message.MessageType {
	case gossip.Push:
		logRecord := domain.Log{} //domain.Log(message.Payload) TODO: Get log
		_, err := mediatr.Send[*commands.AddLogCommand, *commands.AddLogResponse](ctx, commands.NewAddLogCommand(logRecord))
		if err != nil {
			log.Fatalf("Error occured err = %v", err)
		}
		break
	case gossip.Pull:
		logPosition := 0 //domain.Log(message.Payload) TODO: Get log position
		nodeId := "todo"
		_, err := mediatr.Send[*commands.SendLogCommand, *commands.SendLogResponse](ctx,
			commands.NewSendLogCommand(message.Sender, nodeId, logPosition))
		if err != nil {
			log.Fatalf("Error occured err = %v", err)
		}
		break
	}
	log.Println(message)
}
