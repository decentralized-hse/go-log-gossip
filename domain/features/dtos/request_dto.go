package dtos

type RequestDTO struct {
	NodeId   string
	Position int
}

func NewRequestDTO(nodeId string, position int) *RequestDTO {
	return &RequestDTO{
		Position: position,
		NodeId:   nodeId,
	}
}

func (l *RequestDTO) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"position": l.Position,
		"node_id":  l.NodeId,
	}
}
