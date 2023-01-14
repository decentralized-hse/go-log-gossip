package gossip

import (
	"encoding/base64"
	"encoding/json"
	"github.com/decentralized-hse/go-log-gossip/infra/keys"
	"log"
)

type notifyDelegate struct {
	handler                    MessageHandler
	receivedNodeIdAddrCallback ReceivedNodeIdAddrCallback
}

func (m *notifyDelegate) NodeMeta(_ int) []byte {
	return nil
}

func (m *notifyDelegate) NotifyMsg(bytes []byte) {
	log.Println("received new message, length: ", len(bytes))

	var message Message
	err := json.Unmarshal(bytes, &message)

	if err != nil {
		log.Println("failed to unmarshal message,", err, " skipping message")
		return
	}

	senderPublicKey, err := keys.DecodePublicKey(string(message.Sender))
	if err != nil {
		log.Println("failed to decode sender public key: ", err)
		return
	}

	messageToVerify, err := json.Marshal(message.Payload)
	if err != nil {
		log.Println("failed to marshal payload: ", err)
		return
	}

	signature, err := base64.StdEncoding.DecodeString(message.Signature)
	if err != nil {
		log.Println("failed to b64decode signature: ", err)
		return
	}

	err = senderPublicKey.VerifySignature(messageToVerify, signature)
	if err != nil {
		log.Println("failed to verify signature: ", err)
		return
	}

	m.receivedNodeIdAddrCallback(message.Sender, message.Meta.Addr)

	m.handler(&message)
}

func (m *notifyDelegate) GetBroadcasts(_, _ int) [][]byte {
	return nil
}

func (m *notifyDelegate) LocalState(_ bool) []byte {
	return nil
}

func (m *notifyDelegate) MergeRemoteState(_ []byte, _ bool) {}
