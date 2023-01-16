package dtos

import (
	"encoding/base64"
	"github.com/decentralized-hse/go-log-gossip/domain"
)

type LogDTO struct {
	Hash     string `json:"hash"`
	Position int    `json:"position"`
	NodeId   string `json:"node_id"`
	Message  string `json:"message"`
}

func NewLogDTO(log *domain.Log) *LogDTO {
	hashEncoded := base64.StdEncoding.EncodeToString(log.Hash)
	return &LogDTO{
		Hash:     hashEncoded,
		Position: log.Position,
		NodeId:   log.NodeId,
		Message:  log.Message,
	}
}

func (l *LogDTO) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"hash":     l.Hash,
		"position": l.Position,
		"node_id":  l.NodeId,
		"message":  l.Message,
	}
}

func (l *LogDTO) ToDomain() (*domain.Log, error) {
	hashDecoded, err := base64.StdEncoding.DecodeString(l.Hash)
	if err != nil {
		return nil, err
	}

	return &domain.Log{
		Hash:     hashDecoded,
		Position: l.Position,
		NodeId:   l.NodeId,
		Message:  l.Message,
	}, nil
}
