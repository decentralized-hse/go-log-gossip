package gossip

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/decentralized-hse/go-log-gossip/infra/config"
	"github.com/decentralized-hse/go-log-gossip/infra/keys"
	"github.com/hashicorp/memberlist"
	"log"
	"sync"
	"time"
)

type MessageHandler func(*Message)

type ReceivedNodeIdAddrCallback func(nodeId string, addr string)

type Gossiper struct {
	ml           *memberlist.Memberlist
	keys         *keys.PublicPrivateKeyPair
	nodeRegistry map[string]string

	encodedSender string
}

func Start(cfg *config.Config, keys *keys.PublicPrivateKeyPair, handler MessageHandler, ctx context.Context, wg *sync.WaitGroup) *Gossiper {
	nodeRegistry := make(map[string]string)

	log.Println("making gossip node config")
	mlConfig := memberlist.DefaultLocalConfig()
	mlConfig.Name = cfg.Gossip.SelfNodeName
	mlConfig.Delegate = &notifyDelegate{
		handler: handler,
		receivedNodeIdAddrCallback: func(nodeId string, addr string) {
			nodeRegistry[nodeId] = addr
		},
	}
	mlConfig.BindPort = cfg.Gossip.SelfNodePort

	log.Println("trying to decode cluster secret key")
	secretKey, err := base64.StdEncoding.DecodeString(cfg.Gossip.SecretKey)
	if err != nil {
		log.Fatalf("failed to decode cluster secret key: %v", err)
	}
	mlConfig.SecretKey = secretKey

	log.Println("trying to create memberlist from config")
	ml, err := memberlist.Create(mlConfig)
	if err != nil {
		log.Fatalf("failed to create memberlist from config: %v", err)
	}

	log.Println("trying to join cluster if node is not boostrap")
	if !cfg.Gossip.IsBoostrapNode {
		log.Printf("node is not boostrap, boostrap node addr: %v", cfg.Gossip.BoostrapNodeAddr)
		_, err = ml.Join([]string{cfg.Gossip.BoostrapNodeAddr})
		if err != nil {
			log.Fatalf("failed to join cluster: %v", err)
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
	log.Printf("encoded current node id: %v", encodedSender)
	return &Gossiper{
		ml:            ml,
		keys:          keys,
		encodedSender: encodedSender,
		nodeRegistry:  nodeRegistry,
	}
}

func (g *Gossiper) BroadcastMessage(messageType MessageType, data any) error {
	bytes, err := g.encodeMessage(messageType, data)

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

// Request не возвращает ответ. Ответ, если таковой подразумевается от toNode, будет приходить в MessageHandler.
// Аналогия: посылаем запрос по шине, а потом, когда-нибудь по шине придет ответ
func (g *Gossiper) Request(toNode string, messageType MessageType, data any) error {
	nodeAddr, isInRegistry := g.nodeRegistry[toNode]
	if !isInRegistry {
		return fmt.Errorf("%v not in current registry", toNode)
	}

	var mlNode *memberlist.Node
	for _, node := range g.ml.Members() {
		if node.Address() == nodeAddr {
			mlNode = node
			break
		}
	}
	if mlNode == nil {
		return fmt.Errorf("%v is unavailable", toNode)
	}

	message, err := g.encodeMessage(messageType, data)
	if err != nil {
		return err
	}

	return g.ml.SendReliable(mlNode, message)
}

func (g *Gossiper) encodeMessage(messageType MessageType, data any) ([]byte, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	signature, err := g.keys.GetPrivateKey().SignMessage(marshaledData)
	if err != nil {
		return nil, err
	}

	signatureEncoded := base64.StdEncoding.EncodeToString(signature)
	senderEncoded := g.encodedSender

	message := &Message{
		MessageType: messageType,
		Sender:      senderEncoded,
		Signature:   signatureEncoded,
		Payload:     data,
		Meta: struct {
			Addr string `json:"addr"`
		}{
			Addr: g.ml.LocalNode().Address(),
		},
	}

	return json.Marshal(message)
}
