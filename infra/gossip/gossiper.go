package gossip

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/decentralized-hse/go-log-gossip/infra/config"
	"github.com/decentralized-hse/go-log-gossip/infra/keys"
	"github.com/hashicorp/memberlist"
	"log"
	"sync"
	"time"
)

type MessageHandler func(*Message)

type Gossiper struct {
	ml   *memberlist.Memberlist
	keys *keys.PublicPrivateKeyPair

	encodedSender string
}

func Start(cfg *config.Config, keys *keys.PublicPrivateKeyPair, handler MessageHandler, ctx context.Context, wg *sync.WaitGroup) *Gossiper {
	mlConfig := memberlist.DefaultLocalConfig()
	mlConfig.Name = cfg.Gossip.SelfNodeName
	mlConfig.Delegate = &notifyDelegate{handler: handler}
	mlConfig.BindPort = cfg.Gossip.SelfNodePort

	secretKey, err := base64.StdEncoding.DecodeString(cfg.Gossip.SecretKey)
	log.Println(len(secretKey))
	if err != nil {
		log.Fatalln(err)
	}
	mlConfig.SecretKey = secretKey

	ml, err := memberlist.Create(mlConfig)
	if err != nil {
		log.Fatalln(err)
	}

	if !cfg.Gossip.IsBoostrapNode {
		_, err = ml.Join([]string{cfg.Gossip.BoostrapNodeAddr})
		if err != nil {
			log.Fatalln(err)
		}
	}

	go func() {
		<-ctx.Done()
		err = ml.Leave(time.Second)
		if err != nil {
			log.Println("err while closing gossip: ", err)
		}
		wg.Done()
	}()

	encodedSender := keys.GetPublicKey().Encode()
	return &Gossiper{
		ml:            ml,
		keys:          keys,
		encodedSender: encodedSender,
	}
}

func (g *Gossiper) BroadcastMessage(messageType MessageType, data any) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	signature, err := g.keys.GetPrivateKey().SignMessage(marshaledData)
	if err != nil {
		return err
	}

	signatureEncoded := base64.StdEncoding.EncodeToString(signature)
	senderEncoded := g.encodedSender

	message := &Message{
		MessageType: messageType,
		Sender:      senderEncoded,
		Signature:   signatureEncoded,
		Payload:     data,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// todo: use random slice of members, not all
	for _, node := range g.ml.Members() {
		if node == g.ml.LocalNode() {
			continue
		}
		go g.ml.SendReliable(node, bytes)
	}

	return nil
}
